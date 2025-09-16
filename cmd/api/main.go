package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"

	"github.com/projeto-estudos/api-golang/database"
	_ "github.com/projeto-estudos/api-golang/docs"
	admincontroller "github.com/projeto-estudos/api-golang/internal/admin/controller"
	adminmodel "github.com/projeto-estudos/api-golang/internal/admin/dto"
	admindatasource "github.com/projeto-estudos/api-golang/internal/admin/external/datasource"
	adminhandler "github.com/projeto-estudos/api-golang/internal/admin/handler"
	authcontroller "github.com/projeto-estudos/api-golang/internal/auth/controller"
	"github.com/projeto-estudos/api-golang/internal/auth/external"
	customercontroller "github.com/projeto-estudos/api-golang/internal/customer/controller"
	customermodel "github.com/projeto-estudos/api-golang/internal/customer/dto"
	customerdatasource "github.com/projeto-estudos/api-golang/internal/customer/external/datasource"
	customerhandler "github.com/projeto-estudos/api-golang/internal/customer/handler"
	"github.com/projeto-estudos/api-golang/internal/http/middleware"
	ordercontroller "github.com/projeto-estudos/api-golang/internal/order/controller"
	ordermodel "github.com/projeto-estudos/api-golang/internal/order/dto"
	orderdatasource "github.com/projeto-estudos/api-golang/internal/order/external/datasource"
	ordergateway "github.com/projeto-estudos/api-golang/internal/order/gateway"
	orderservicegateway "github.com/projeto-estudos/api-golang/internal/order/gateway/services"
	orderhandler "github.com/projeto-estudos/api-golang/internal/order/handler"
	orderusecases "github.com/projeto-estudos/api-golang/internal/order/usecases"
	paymentcontroller "github.com/projeto-estudos/api-golang/internal/payment/controllers"
	paymentmodel "github.com/projeto-estudos/api-golang/internal/payment/dto"
	paymentdatasource "github.com/projeto-estudos/api-golang/internal/payment/external/datasource"
	paymentgateway "github.com/projeto-estudos/api-golang/internal/payment/gateway"
	paymentservicegateway "github.com/projeto-estudos/api-golang/internal/payment/gateway/services"
	paymenthandler "github.com/projeto-estudos/api-golang/internal/payment/handlers"
	paymentusecases "github.com/projeto-estudos/api-golang/internal/payment/usecases"
	productcontroller "github.com/projeto-estudos/api-golang/internal/product/controller"
	productmodel "github.com/projeto-estudos/api-golang/internal/product/dto"
	productdatasource "github.com/projeto-estudos/api-golang/internal/product/external/datasource"
	productgateway "github.com/projeto-estudos/api-golang/internal/product/gateway"
	productservicegateway "github.com/projeto-estudos/api-golang/internal/product/gateway/services"
	producthandler "github.com/projeto-estudos/api-golang/internal/product/handler"
	productusecases "github.com/projeto-estudos/api-golang/internal/product/usecases"
	productordermodel "github.com/projeto-estudos/api-golang/internal/productorder/dto"
	productorderdatasource "github.com/projeto-estudos/api-golang/internal/productorder/external/datasource"
	productordergateway "github.com/projeto-estudos/api-golang/internal/productorder/gateway"
	productorderservicegateway "github.com/projeto-estudos/api-golang/internal/productorder/gateway/services"
	productorderusecases "github.com/projeto-estudos/api-golang/internal/productorder/usecases"
	qrcodeprovider "github.com/projeto-estudos/api-golang/internal/qrcodeproviders/gateways"
)

// @title           GoLunch
// @version         1.0
// @description     REST API to facilitate order management in a snack bar.
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {

	// UNCOMMENT TO RUN ONLY THE DATABASE IN DOCKER
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Erro ao carregar o .env")
	// }

	r := gin.Default()
	loadYAML()

	db := database.NewPostgresDatabase().GetDb()

	if err := db.AutoMigrate(
		&customermodel.CustomerDAO{},
		&adminmodel.AdminDAO{},
		&productmodel.ProductDAO{},
		&ordermodel.OrderDAO{},
		&productordermodel.ProductOrderDAO{},
		&paymentmodel.PaymentDAO{},
	); err != nil {
		log.Fatalf("Erro ao migrar o banco: %v", err)
	}

	// Serve static files
	uploadDir := os.Getenv("UPLOAD_DIR")

	// Jwt service for generate and validate tokens (CLEANARCH)
	jwtGateway := external.NewJWTService(os.Getenv("SECRET_KEY"), 24*time.Hour)
	authController := authcontroller.New(jwtGateway)

	// Customer
	customerDatasource := customerdatasource.New(db)
	customerController := customercontroller.Build(customerDatasource, authController)
	customerHandler := customerhandler.New(customerController)

	// Product
	productDataSource := productdatasource.New(db)
	productController := productcontroller.Build(productDataSource)
	productHandler := producthandler.New(productController)

	// Admin
	adminDatasource := admindatasource.New(db)
	adminController := admincontroller.Build(adminDatasource, authController)
	adminHandler := adminhandler.New(adminController)

	// Product Order
	productOrderDataSource := productorderdatasource.New(db)
	productOrderGateway := productordergateway.Build(productOrderDataSource)
	productOrderUseCase := productorderusecases.Build(*productOrderGateway)

	// Payment
	paymentDataSource := paymentdatasource.New(db)
	paymentGateway := paymentgateway.Build(paymentDataSource)

	// QR Code Client
	qrCodeClient := qrcodeprovider.New()

	// Order Data Source and Gateway
	orderDataSource := orderdatasource.New(db)
	orderGateway := ordergateway.Build(orderDataSource)

	// Common Gateways
	productGateway := productgateway.Build(productDataSource)
	productUseCase := productusecases.Build(*productGateway)

	productServiceGateway := productservicegateway.NewProductServiceGateway(productUseCase)
	productOrderServiceGatewayForOrder, productOrderServiceGatewayForPayment := productorderservicegateway.NewProductOrderServiceGateway(productOrderUseCase)

	// Creating payment use case without orderService (to avoid circular dependency)
	paymentUseCaseWithoutOrder := paymentusecases.Build(paymentGateway, qrCodeClient, productServiceGateway, productOrderServiceGatewayForPayment, nil)
	paymentServiceGateway := paymentservicegateway.NewPaymentServiceGateway(paymentUseCaseWithoutOrder)

	// Creating orderUseCase with productService and productOrderService (to avoid circular dependency)
	orderUseCase := orderusecases.Build(orderGateway, productServiceGateway, productOrderServiceGatewayForOrder, paymentServiceGateway)

	// Creating orderServiceGateway with orderUseCase
	orderServiceGateway := orderservicegateway.NewOrderServiceGateway(orderUseCase)

	// Creating payment use case with orderServiceGateway
	paymentUseCase := paymentusecases.Build(paymentGateway, qrCodeClient, productServiceGateway, productOrderServiceGatewayForPayment, orderServiceGateway)

	// Order Controller and Handler
	orderController := ordercontroller.Build(orderUseCase)
	orderHandler := orderhandler.New(orderController)

	// Payment
	paymentController := paymentcontroller.Build(paymentUseCase)
	paymentHandler := paymenthandler.New(paymentController)

	// Default Routes
	r.GET("/ping", ping)
	r.GET("/swagger/*any", ginswagger.WrapHandler(swaggerfiles.Handler))
	r.Static("/uploads", uploadDir)

	// Public Routes (login/register)
	r.GET("/customer/identify/:cpf", customerHandler.Identify)
	r.GET("/customer/anonymous", customerHandler.Anonymous)
	r.POST("/customer/register", customerHandler.Create)
	r.POST("/admin/register", adminHandler.Register)
	r.POST("/admin/login", adminHandler.Login)

	// Webhook for Mercado Pago
	r.POST("/webhook/payment/check", paymentHandler.CheckPayment)

	// Authenticated Group
	authenticated := r.Group("/")
	authenticated.Use(middleware.AuthMiddleware(authController))

	// Routes for regular authenticated users
	// Product
	authenticated.GET("/product/categories", productHandler.ListCategories)
	authenticated.GET("/product", productHandler.GetAllByCategory)

	// Order
	authenticated.POST("/order", orderHandler.Create)
	authenticated.GET("/order", middleware.AdminOnly(), orderHandler.GetAll)
	authenticated.PUT("/order/:id", middleware.AdminOnly(), orderHandler.Update)
	authenticated.GET("/order/panel", middleware.AdminOnly(), orderHandler.GetPanel)

	// Group for admin users inside authenticated group
	adminRoutes := authenticated.Group("/product")
	adminRoutes.Use(middleware.AdminOnly())
	adminRoutes.POST("/image/upload", productHandler.UploadImage)
	adminRoutes.POST("/", productHandler.Create)
	adminRoutes.PUT("/:id", productHandler.Update)
	adminRoutes.DELETE("/:id", productHandler.Delete)

	r.Run(":8080")
}

func loadYAML() {
	viper.SetConfigName("default")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./conf/environment")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("error reading yaml config: %v", err)
	}
}

// Ping godoc
// @Summary      Answers with "pong"
// @Description  Health Check
// @Tags         Ping
// @Accept       json
// @Produce      json
// @Success      200 {object}  PongResponse
// @Router       /ping [get]
func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

type PongResponse struct {
	Message string `json:"message"`
}
