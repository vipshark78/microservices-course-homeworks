package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

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
	// Создаем хранилище заказов
	orderStorage := NewOrderStorage()

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
	inventoryClient := inventory_v1.NewInventoryServiceClient(inventortyClientConn)
	paymentClient := payment_v1.NewPaymentServiceClient(paymentClientConn)

	defer func() {
		if err := inventortyClientConn.Close(); err != nil {
			log.Printf("Ошибка закрытия соединения с сервисом Inventory: %v\n", err)
		}
		if err := paymentClientConn.Close(); err != nil {
			log.Printf("Ошибка закрытия соединения с сервисом Payment: %v\n", err)
		}
	}()

	// Создаем обработчик для операций с заказами
	orderHandler := NewOrderHandler(orderStorage, inventoryClient, paymentClient)

	// Создаем сервер OpenAPI
	orderServer, err := order_v1.NewServer(orderHandler)
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

// OrderStorage представляет хранилище заказов
type OrderStorage struct {
	orders map[string]*order_v1.OrderDto
	mu     sync.RWMutex
}

// NewOrderStorage создает новое хранилище заказов.
func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]*order_v1.OrderDto),
	}
}

// Insert добавляет новый заказ в хранилище.
func (o *OrderStorage) Insert(userUUID uuid.UUID, partsUUIDs []uuid.UUID, price float64) (uuid.UUID, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	orderUUID := uuid.New()
	order := &order_v1.OrderDto{
		UserUUID:        userUUID,
		OrderUUID:       orderUUID,
		PartUuids:       partsUUIDs,
		TotalPrice:      price,
		TransactionUUID: order_v1.OptNilUUID{},
		PaymentMethod:   order_v1.OptNilPaymentMethod{},
		Status:          order_v1.OrderStatusPENDINGPAYMENT,
	}

	o.orders[order.OrderUUID.String()] = order
	return orderUUID, nil
}

// Read возвращает заказ по его идентификатору.
func (o *OrderStorage) Read(uuid uuid.UUID) *order_v1.OrderDto {
	o.mu.RLock()
	defer o.mu.RUnlock()

	order, ok := o.orders[uuid.String()]
	if !ok {
		return nil
	}
	return order
}

// Update обновляет информацию о заказе.
func (o *OrderStorage) Update(order *order_v1.OrderDto) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.orders[order.OrderUUID.String()] = order
	return nil
}

// OrderHandler представляет обработчик для операций с заказами.
type OrderHandler struct {
	storage                *OrderStorage
	inventoryServiceClient inventory_v1.InventoryServiceClient
	paymentServiceClient   payment_v1.PaymentServiceClient
}

// NewOrderHandler создает новый экземпляр OrderHandler.
func NewOrderHandler(storage *OrderStorage, inventoryServiceClient inventory_v1.InventoryServiceClient, paymentServiceClient payment_v1.PaymentServiceClient) *OrderHandler {
	return &OrderHandler{
		inventoryServiceClient: inventoryServiceClient,
		paymentServiceClient:   paymentServiceClient,
		storage:                storage,
	}
}

// validateCreateOrderRequest выполняет валидацию запроса на создание заказа.
func (o *OrderHandler) validateCreateOrderRequest(req *order_v1.CreateOrderRequest) error {
	if err := req.Validate(); err != nil {
		return fmt.Errorf("ошибка валидации запроса: %w", err)
	}

	if err := uuid.Validate(req.UserUUID.String()); err != nil {
		return fmt.Errorf("ошибка валидации UUID пользователя: %w", err)
	}
	return nil
}

// convertUUIDToSliceString преобразует слайс UUID в слайс строк.
func (o *OrderHandler) convertUUIDToSliceString(uuids []uuid.UUID) ([]string, error) {
	strUuids := make([]string, 0, len(uuids))
	for _, UUID := range uuids {
		stringUUID := UUID.String()
		if err := uuid.Validate(stringUUID); err != nil {
			return nil, fmt.Errorf("ошибка валидации UUID: %w", err)
		}
		strUuids = append(strUuids, stringUUID)
	}
	return strUuids, nil
}

// priceCalculate рассчитывает общую стоимость заказа.
func (o *OrderHandler) priceCalculate(reqUuids []string, parts []*inventory_v1.Part) (float64, error) {
	var totalPrice float64

	for _, uuid := range reqUuids {
		exist := false
		for _, part := range parts {
			if part.Uuid == uuid {
				totalPrice += part.Price
				exist = true
				break
			}
		}
		if !exist {
			return 0, errors.New("указанного UUID детали не существует")
		}
	}
	return totalPrice, nil
}

// CreateOrder создает новый заказ.
func (o *OrderHandler) CreateOrder(ctx context.Context, req *order_v1.CreateOrderRequest) (order_v1.CreateOrderRes, error) {
	if err := o.validateCreateOrderRequest(req); err != nil {
		return &order_v1.BadRequestError{Code: 400, Message: "Bad Request"}, err
	}

	reqUuids, err := o.convertUUIDToSliceString(req.PartUuids)
	if err != nil {
		return &order_v1.BadRequestError{Code: 400, Message: "Bad Request"}, err
	}

	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()

	listPartsResp, err := o.inventoryServiceClient.ListParts(ctx, &inventory_v1.ListPartsRequest{Filter: &inventory_v1.PartsFilter{
		Uuids: reqUuids,
	}})
	if err != nil {
		return &order_v1.InternalServerError{Code: 500, Message: "Internal Server Error"}, err
	}

	totalPrice, err := o.priceCalculate(reqUuids, listPartsResp.GetParts())
	if err != nil {
		return &order_v1.BadRequestError{Code: 400, Message: "Bad Request"}, err
	}

	orderUUID, err := o.storage.Insert(req.UserUUID, req.PartUuids, totalPrice)
	if err != nil {
		return &order_v1.InternalServerError{Code: 500, Message: "Internal Server Error"}, err
	}

	return &order_v1.CreateOrderResponse{OrderUUID: orderUUID, TotalPrice: totalPrice}, nil
}

// OrderByUUID получает информацию о заказе по его идентификатору.
func (o *OrderHandler) OrderByUUID(ctx context.Context, params order_v1.OrderByUUIDParams) (order_v1.OrderByUUIDRes, error) {
	if err := uuid.Validate(params.OrderUUID.String()); err != nil {
		return &order_v1.BadRequestError{Code: 400, Message: "Bad Request"}, err
	}

	order := o.storage.Read(params.OrderUUID)

	if order == nil {
		return &order_v1.BadRequestError{Code: 404, Message: "Order Not Found"}, nil
	}

	return &order_v1.GetOrderResponse{AllOf: order_v1.NewOptOrderDto(*order)}, nil
}

// OrderCancel отменяет заказ.
func (o *OrderHandler) OrderCancel(ctx context.Context, params order_v1.OrderCancelParams) (order_v1.OrderCancelRes, error) {
	if err := uuid.Validate(params.OrderUUID.String()); err != nil {
		return &order_v1.BadRequestError{Code: 400, Message: "Bad Request"}, err
	}

	order := o.storage.Read(params.OrderUUID)
	if order == nil {
		return &order_v1.BadRequestError{Code: 404, Message: "Order Not Found"}, nil
	}

	if order.Status == order_v1.OrderStatusPAID {
		return &order_v1.ConflictError{Code: 409, Message: "Order has already paid and cannot be cancelled"}, nil
	}

	order.Status = order_v1.OrderStatusCANCELLED
	if err := o.storage.Update(order); err != nil {
		return &order_v1.BadRequestError{Code: 500, Message: "Internal Server Error"}, err
	}
	return &order_v1.OrderCancelNoContent{}, nil
}

// validateOrderPayRequest выполняет валидацию запроса на оплату заказа.
func (o *OrderHandler) validateOrderPayRequest(req *order_v1.PayOrderRequest, params order_v1.OrderPayParams) error {
	if err := req.Validate(); err != nil {
		return fmt.Errorf("ошибка валидации запроса: %w", err)
	}

	if err := uuid.Validate(params.OrderUUID.String()); err != nil {
		return fmt.Errorf("ошибка валидации UUID заказа: %w", err)
	}

	if req.PaymentMethod.IsNull() {
		return fmt.Errorf("не указан метод оплаты")
	}
	return nil
}

// OrderPay оплачивает заказ.
func (o *OrderHandler) OrderPay(ctx context.Context, req *order_v1.PayOrderRequest, params order_v1.OrderPayParams) (order_v1.OrderPayRes, error) {
	if err := o.validateOrderPayRequest(req, params); err != nil {
		return &order_v1.BadRequestError{Code: 400, Message: "Bad Request"}, err
	}

	order := o.storage.Read(params.OrderUUID)
	if order == nil {
		return &order_v1.BadRequestError{Code: 404, Message: "Order Not Found"}, nil
	}

	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()

	paymentMethod, ok := req.PaymentMethod.Get()
	if !ok {
		return &order_v1.BadRequestError{Code: 400, Message: "Bad Request"}, fmt.Errorf("е указан метод оплаты")
	}

	paymentMethodBytes, err := paymentMethod.MarshalText()
	if err != nil {
		return &order_v1.InternalServerError{Code: 500, Message: "Internal Server Error"}, err
	}

	payOrderResp, err := o.paymentServiceClient.PayOrder(ctx, &payment_v1.PayOrderRequest{
		OrderUuid:     params.OrderUUID.String(),
		UserUuid:      order.UserUUID.String(),
		PaymentMethod: payment_v1.PaymentMethod(payment_v1.PaymentMethod_value[string(paymentMethodBytes)]),
	})
	if err != nil {
		return &order_v1.InternalServerError{Code: 500, Message: "Internal Server Error"}, err
	}

	transactionUuid, err := uuid.Parse(payOrderResp.TransactionUuid)
	if err != nil {
		return &order_v1.InternalServerError{Code: 500, Message: "Internal Server Error"}, err
	}

	order.Status = order_v1.OrderStatusPAID
	order.TransactionUUID = order_v1.NewOptNilUUID(transactionUuid)
	order.PaymentMethod = order_v1.NewOptNilPaymentMethod(paymentMethod)

	if err := o.storage.Update(order); err != nil {
		return &order_v1.InternalServerError{Code: 500, Message: "Internal Server Error"}, err
	}

	return &order_v1.PayOrderResponse{TransactionUUID: transactionUuid}, nil
}

// NewError создает ошибку с указанным статус кодом и сообщением.
func (o *OrderHandler) NewError(ctx context.Context, err error) *order_v1.GenericErrorStatusCode {
	return &order_v1.GenericErrorStatusCode{
		StatusCode: 500,
		Response: order_v1.GenericError{
			Message: "Internal Server Error",
			Code:    500,
		},
	}
}
