package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	api "github.com/vipshark78/microservices-course-homeworks/order/internal/api/order/v1"
	inventory "github.com/vipshark78/microservices-course-homeworks/order/internal/client/grpc/inventory/v1"
	payment "github.com/vipshark78/microservices-course-homeworks/order/internal/client/grpc/payment/v1"
	repository "github.com/vipshark78/microservices-course-homeworks/order/internal/repository/order"
	service "github.com/vipshark78/microservices-course-homeworks/order/internal/service/order"
	order_v1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/openapi/order/v1"
	inventory_v1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/proto/inventory/v1"
	payment_v1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/proto/payment/v1"
)

const (
	httpPort          = "8080"
	readHeaderTimeout = 10 * time.Second
	shutdownTimeout   = 10 * time.Second
	operationTimeout  = 10 * time.Second
	inventoryAddress  = "localhost:50051"
	paymentAddress    = "localhost:50052"
)

func main() {
	// Создаем клиентское соединение с сервисом Inventory
	inventortyClientConn, err := grpc.NewClient(inventoryAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Ошибка подключения к сервису Inventory: %v", err)
	}

	// Создаем клиентское соединение с сервисом Payment
	paymentClientConn, err := grpc.NewClient(paymentAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Ошибка подключения к сервису Payment: %v", err)
	}

	// Получаем клиентские соединения с сервисами Inventory и Payment
	inventoryServiceClient := inventory_v1.NewInventoryServiceClient(inventortyClientConn)
	paymentServiceClient := payment_v1.NewPaymentServiceClient(paymentClientConn)

	defer func() {
		if err := inventortyClientConn.Close(); err != nil {
			log.Printf("Ошибка закрытия соединения с сервисом Inventory: %v\n", err)
		}
	}()

	defer func() {
		if err := paymentClientConn.Close(); err != nil {
			log.Printf("Ошибка закрытия соединения с сервисом Payment: %v\n", err)
		}
	}()

	inventoryClient := inventory.NewInventoryClient(inventoryServiceClient)
	paymentClient := payment.NewPaymentClient(paymentServiceClient)

	// Создаем хранилище заказов
	orderRepository := repository.NewRepository()

	orderService := service.NewService(orderRepository, inventoryClient, paymentClient)

	api := api.NewAPI(orderService)
	// Создаем сервер OpenAPI
	orderServer, err := order_v1.NewServer(api)
	if err != nil {
		log.Printf("ошибка создания сервера OpenAPI: %v\n", err)
	}

	// Создаем маршрутизатор HTTP
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(10 * time.Second))

	// Монтируем обработчики OpenAPI
	router.Mount("/", orderServer)

	// Запускаем HTTP-сервер
	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           router,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	defer func() {
		if err := server.Close(); err != nil {
			log.Printf("Ошибка закрытия сервера: %v", err)
		}
	}()

	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	// Graceful shutdown
	stopChan := make(chan os.Signal, 1)

	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)

	<-stopChan

	log.Println("🛑 Завершение работы сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}
