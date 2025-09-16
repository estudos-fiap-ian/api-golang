#!/bin/bash

# Setup script for GitHub repository secrets
# This script helps you configure the necessary secrets for the EKS deployment workflow

echo "🔧 GitHub Repository Secrets Setup for EKS Deployment"
echo "=================================================="

# Check if gh CLI is installed
if ! command -v gh &> /dev/null; then
    echo "❌ GitHub CLI (gh) is not installed. Please install it first:"
    echo "   brew install gh"
    echo "   or visit: https://cli.github.com/"
    exit 1
fi

# Check if user is authenticated
if ! gh auth status &> /dev/null; then
    echo "🔐 Please authenticate with GitHub first:"
    echo "   gh auth login"
    exit 1
fi

echo "📋 Required secrets for EKS deployment:"
echo ""

# Function to set secret
set_secret() {
    local secret_name=$1
    local secret_description=$2

    echo "Setting up: $secret_name"
    echo "Description: $secret_description"
    echo -n "Enter value for $secret_name: "
    read -s secret_value
    echo ""

    if [ -n "$secret_value" ]; then
        echo "$secret_value" | gh secret set "$secret_name"
        echo "✅ Secret $secret_name has been set"
    else
        echo "⚠️  Skipping $secret_name (empty value)"
    fi
    echo ""
}

# Set required secrets
set_secret "AWS_ACCESS_KEY_ID" "AWS Access Key ID with EKS and ECR permissions"
set_secret "AWS_SECRET_ACCESS_KEY" "AWS Secret Access Key"

echo "🎯 Optional: Update environment variables in the workflow file if needed:"
echo "   - AWS_REGION (default: us-east-1)"
echo "   - EKS_CLUSTER_NAME (default: fiap-cluster)"
echo "   - ECR_REPOSITORY (default: golunch-api)"
echo ""

echo "📝 Next steps:"
echo "   1. Ensure your EKS cluster exists (from terraform-infra repository)"
echo "   2. Commit and push to trigger the deployment workflow"
echo "      (ECR repository will be created automatically)"
echo ""

echo "✅ Setup complete! Your repository is ready for EKS deployment."