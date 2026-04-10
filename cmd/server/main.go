package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/asino-nelson/safiri-logistics/internal/auth"
	"github.com/asino-nelson/safiri-logistics/internal/config"
	"github.com/asino-nelson/safiri-logistics/internal/database"
	"github.com/asino-nelson/safiri-logistics/internal/dispatch"
	"github.com/asino-nelson/safiri-logistics/internal/driver"
	"github.com/asino-nelson/safiri-logistics/internal/middleware"
	"github.com/asino-nelson/safiri-logistics/internal/order"
	"github.com/asino-nelson/safiri-logistics/internal/payment"
	"github.com/asino-nelson/safiri-logistics/internal/pricing"
	"github.com/asino-nelson/safiri-logistics/internal/tracking"
	"github.com/asino-nelson/safiri-logistics/internal/user"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pool, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	userRepository := user.NewRepository(pool)
	userService := user.NewService(userRepository)
	tokenManager := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTTokenLifetime)
	authService := auth.NewService(userRepository, auth.BcryptHasher{}, tokenManager)
	authHandler := auth.NewHandler(authService, userService)
	driverRepository := driver.NewRepository(pool)
	driverService := driver.NewService(driverRepository)
	driverHandler := driver.NewHandler(driverService)
	loadRepository := order.NewRepository(pool)
	dispatchService := dispatch.NewService(loadRepository, driverService, nil)
	dispatchHandler := dispatch.NewHandler(dispatchService)
	pricingService := pricing.NewService()
	loadService := order.NewService(loadRepository, driverService, dispatchService, pricingService)
	loadHandler := order.NewHandler(loadService)
	paymentRepository := payment.NewRepository(pool)
	paymentService := payment.NewService(paymentRepository, loadRepository, payment.NewSandboxMPESAClient())
	paymentHandler := payment.NewHandler(paymentService)
	trackingRepository := tracking.NewRepository(pool)
	trackingService := tracking.NewService(trackingRepository, loadRepository, nil)
	trackingHandler := tracking.NewHandler(trackingService)
	authMiddleware := middleware.Authenticate(tokenManager)

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := router.Group("/api/v1")
	authHandler.RegisterRoutes(api.Group("/auth"), authMiddleware)
	dispatchHandler.RegisterRoutes(api, authMiddleware)
	driverHandler.RegisterRoutes(api, authMiddleware)
	loadHandler.RegisterRoutes(api, authMiddleware)
	paymentHandler.RegisterRoutes(api, authMiddleware)
	trackingHandler.RegisterRoutes(api, authMiddleware)

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("server shutdown error: %v", err)
		}
	}()

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
