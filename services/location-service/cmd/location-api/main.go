package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adityjoshi/Uber-Service/services/location-service/internal/config"
	"github.com/adityjoshi/Uber-Service/services/location-service/internal/controller"
	"github.com/adityjoshi/Uber-Service/services/location-service/internal/services"
	"github.com/gin-gonic/gin"
)

func main() {
	redisCfg := config.LoadRedisFromEnv()
	rdb := config.NewRedisClient(redisCfg)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	if err := config.PingRedis(ctx, rdb); err != nil {
		cancel()
		_ = rdb.Close()
		log.Fatalf("redis connection failed: %v", err)
	}
	cancel()

	log.Println("Redis connection OK")

	svc := services.NewLocationService(rdb)
	if err := svc.Validate(); err != nil {
		_ = rdb.Close()
		log.Fatalf("service: %v", err)
	}

	handler := controller.NewLocationHandler(svc)

	engine := gin.Default()
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})
	handler.RegisterRoutes(engine)

	addr := getenvDefault("HTTP_ADDR", ":8082")
	srv := &http.Server{
		Addr:              addr,
		Handler:           engine,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("location-service listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown: %v", err)
	}
	_ = rdb.Close()
}

func getenvDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
