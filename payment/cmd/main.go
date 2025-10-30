package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	api "github.com/vipshark78/microservices-course-homeworks/payment/internal/api/payment/v1"
	"github.com/vipshark78/microservices-course-homeworks/payment/internal/interceptor"
	service "github.com/vipshark78/microservices-course-homeworks/payment/internal/service/payment"
	payment_v1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/proto/payment/v1"
)

const grpcPort = 50052

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}
	defer func() {
		if err := lis.Close(); err != nil {
			log.Printf("failed to close listener: %v\n", err)
		}
	}()

	// Создаем gRPC сервер с интерцептором логирования
	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc.UnaryServerInterceptor(interceptor.LoggerInterceptor()),
		),
	)

	// Регистрируем наш сервис
	paymentService := service.NewPaymentService()
	api := api.NewApi(paymentService)
	payment_v1.RegisterPaymentServiceServer(s, api)

	// Включаем рефлексию для отладки
	reflection.Register(s)

	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}
