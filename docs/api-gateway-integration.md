# API Gateway Integration with EKS

## Architecture Overview

```
Internet → API Gateway → Lambda (auth) + VPC Link → ALB/NLB → EKS API → PostgreSQL
```

## Integration Strategy

### Option 1: VPC Link with Network Load Balancer (Recommended)

1. **Deploy API with Internal NLB**
   - Use `app-service-loadbalancer.yaml` for internal NLB
   - API Gateway connects via VPC Link to internal load balancer

2. **API Gateway Configuration**
   ```yaml
   # In your serverless repository's serverless.yml
   functions:
     api-proxy:
       handler: proxy.handler
       events:
         - http:
             path: /{proxy+}
             method: ANY
             integration: http-proxy
             request:
               uri: http://${env:EKS_LOAD_BALANCER_URL}/{proxy}
   ```

### Option 2: Direct Service Integration

1. **VPC Link Setup**
   ```yaml
   # serverless.yml addition
   resources:
     Resources:
       VpcLink:
         Type: AWS::ApiGateway::VpcLink
         Properties:
           Name: eks-vpc-link
           TargetArns:
             - ${env:EKS_NLB_ARN}
   ```

## Implementation Steps

### Step 1: Add VPC Link Module to Serverless Repository

Create `serverless/modules/vpc-link/main.tf`:
```hcl
# VPC Link for API Gateway → EKS integration
resource "aws_apigatewayv2_vpc_link" "eks_link" {
  name               = "eks-integration-link"
  security_group_ids = [aws_security_group.vpc_link.id]
  subnet_ids         = var.private_subnet_ids

  tags = {
    Name = "EKS-API-Gateway-Link"
  }
}

resource "aws_security_group" "vpc_link" {
  name_prefix = "vpc-link-"
  vpc_id      = var.vpc_id

  egress {
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = ["10.0.0.0/16"]
  }
}
```

### Step 2: Update Serverless Configuration

Add to your serverless repository:

```yaml
# serverless.yml
custom:
  eksLoadBalancerUrl: ${env:EKS_LOAD_BALANCER_URL, 'localhost:8080'}

functions:
  # Existing auth functions...

  # API Proxy for non-auth routes
  apiProxy:
    handler: src/handlers/proxy.handler
    events:
      - http:
          path: /api/{proxy+}
          method: ANY
          integration: http-proxy
          request:
            uri: http://${self:custom.eksLoadBalancerUrl}/{proxy}
            headers:
              X-Forwarded-For: true
              X-Forwarded-Proto: true

resources:
  Resources:
    # VPC Link for internal communication
    ApiVpcLink:
      Type: AWS::ApiGateway::VpcLink
      Properties:
        Name: ${self:service}-eks-link
        Description: VPC Link for EKS API communication
        TargetArns:
          - ${env:EKS_NLB_ARN}
```

### Step 2: Create Proxy Handler

```javascript
// src/handlers/proxy.js
exports.handler = async (event) => {
  const { path, httpMethod, headers, body, queryStringParameters } = event;

  // Forward request to EKS API
  const response = await fetch(`http://${process.env.EKS_LOAD_BALANCER_URL}${path}`, {
    method: httpMethod,
    headers: {
      ...headers,
      'Content-Type': 'application/json',
    },
    body: body
  });

  const responseBody = await response.text();

  return {
    statusCode: response.status,
    headers: {
      'Content-Type': 'application/json',
      'Access-Control-Allow-Origin': '*'
    },
    body: responseBody
  };
};
```

### Step 3: Environment Configuration

Add to your deployment pipeline:

```bash
# Get NLB ARN after EKS deployment
NLB_ARN=$(aws elbv2 describe-load-balancers \
  --names go-web-api-nlb \
  --query 'LoadBalancers[0].LoadBalancerArn' \
  --output text)

# Get NLB DNS name
NLB_DNS=$(aws elbv2 describe-load-balancers \
  --names go-web-api-nlb \
  --query 'LoadBalancers[0].DNSName' \
  --output text)

# Set environment variables for serverless deployment
export EKS_NLB_ARN=$NLB_ARN
export EKS_LOAD_BALANCER_URL=$NLB_DNS:8080
```

## Route Planning

### Authenticated Routes (via Lambda)
- `/customer/register` → Lambda (Cognito integration)
- `/admin/login` → Lambda (Cognito integration)
- `/admin/register` → Lambda (Cognito integration)

### API Routes (via VPC Link)
- `/api/product/*` → EKS API
- `/api/order/*` → EKS API
- `/api/payment/*` → EKS API
- `/ping` → EKS API (health check)

### Combined Authentication Flow
1. User authenticates via API Gateway → Lambda → Cognito
2. JWT token returned to client
3. Subsequent API calls include JWT in Authorization header
4. API Gateway validates JWT and forwards to EKS API
5. EKS API validates JWT and processes request

## Security Considerations

1. **Network Security**
   - Use internal NLB (not internet-facing)
   - VPC Link ensures private communication
   - Security groups restrict access

2. **Authentication**
   - JWT validation in both API Gateway and EKS API
   - Cognito integration for user management
   - Admin role validation

3. **CORS Configuration**
   - Configure CORS in API Gateway
   - Ensure proper headers in EKS API responses

## Monitoring & Logging

1. **CloudWatch Integration**
   - API Gateway logs
   - Lambda function logs
   - EKS pod logs via CloudWatch Container Insights

2. **Health Checks**
   - API Gateway health check → `/ping`
   - EKS readiness/liveness probes
   - ALB/NLB health checks