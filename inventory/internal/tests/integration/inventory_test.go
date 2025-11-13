//go:build integration

package integration

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	sharedErrors "github.com/vipshark78/microservices-course-homeworks/shared/pkg/errors"
	inventoryV1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/proto/inventory/v1"
)

var _ = Describe("InventoryService", func() {
	var (
		ctx             context.Context
		cancel          context.CancelFunc
		inventoryClient inventoryV1.InventoryServiceClient
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(suiteCtx)

		// Создаём gRPC клиент
		conn, err := grpc.NewClient(
			env.App.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		Expect(err).ToNot(HaveOccurred(), "ожидали успешное подключение к gRPC приложению")

		inventoryClient = inventoryV1.NewInventoryServiceClient(conn)
	})

	AfterEach(func() {
		// Чистим коллекцию после теста
		err := env.ClearPartsCollection(ctx)
		Expect(err).ToNot(HaveOccurred(), "ожидали успешную очистку коллекции")

		cancel()
	})

	Describe("GetPart", func() {
		var partUUID string

		BeforeEach(func() {
			// Вставляем тестовую деталь
			var err error
			partUUID, err = env.InsertTestPart(ctx)
			Expect(err).ToNot(HaveOccurred(), "ожидали успешную вставку тестовой детали в MongoDB")
		})

		It("должен успешно возвращать деталь по UUID", func() {
			resp, err := inventoryClient.GetPart(ctx, &inventoryV1.GetPartRequest{
				Uuid: partUUID,
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetPart()).ToNot(BeNil())
			Expect(resp.GetPart().GetUuid()).To(Equal(partUUID))
			Expect(resp.GetPart().GetCategory()).To(BeNumerically(">", 0))
			Expect(resp.GetPart().GetCreatedAt()).ToNot(BeNil())
			Expect(resp.GetPart().GetDescription()).ToNot(BeEmpty())
			Expect(resp.GetPart().GetUpdatedAt()).ToNot(BeNil())
			Expect(resp.GetPart().GetDimensions()).ToNot(BeNil())
			Expect(resp.GetPart().GetManufacturer()).ToNot(BeNil())
			Expect(resp.GetPart().GetMetadata()).ToNot(BeNil())
			Expect(resp.GetPart().GetName()).ToNot(BeEmpty())
			Expect(resp.GetPart().GetPrice()).To(BeNumerically(">", 0))
			Expect(resp.GetPart().GetStockQuantity()).To(BeNumerically(">", 0))
			Expect(resp.Part.GetTags()).ToNot(BeNil())
		})

		It("Должна вернуться ошибка, передан nil request", func() {
			resp, err := inventoryClient.GetPart(ctx, nil)
			Expect(err).To(HaveOccurred())
			grpcStatus, ok := status.FromError(err)
			Expect(ok).To(BeTrue())
			Expect(grpcStatus.Code()).To(Equal(sharedErrors.ErrorCodeToGRPCCode(sharedErrors.InvalidArgumentErrCode)))
			Expect(resp).To(BeNil())
		})

		It("Должна вернуться ошибка, передан пустой UUID", func() {
			resp, err := inventoryClient.GetPart(ctx, &inventoryV1.GetPartRequest{})
			Expect(err).To(HaveOccurred())
			grpcStatus, ok := status.FromError(err)
			Expect(ok).To(BeTrue())
			Expect(grpcStatus.Code()).To(Equal(sharedErrors.ErrorCodeToGRPCCode(sharedErrors.InvalidArgumentErrCode)))
			Expect(resp).To(BeNil())
		})

		It("Должна вернуться ошибка, передан несуществующий UUID", func() {
			resp, err := inventoryClient.GetPart(ctx, &inventoryV1.GetPartRequest{
				Uuid: gofakeit.UUID(),
			})
			Expect(err).To(HaveOccurred())
			grpcStatus, ok := status.FromError(err)
			Expect(ok).To(BeTrue())
			Expect(grpcStatus.Code()).To(Equal(sharedErrors.ErrorCodeToGRPCCode(sharedErrors.NotFoundErrCode)))
			Expect(resp).To(BeNil())
		})
	})

	Describe("ListParts", func() {
		var partUUIDs []string

		BeforeEach(func() {
			partUUIDs = make([]string, 0, 5)
			// Вставляем тестовые детали
			for range 5 {
				partUUID, err := env.InsertTestPart(ctx)
				Expect(err).ToNot(HaveOccurred(), "ожидали успешную вставку тестовой детали в MongoDB")
				partUUIDs = append(partUUIDs, partUUID)
			}
		})

		It("Должны вернуться все детали, передан пустой фильтр", func() {
			resp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
				Filter: &inventoryV1.PartsFilter{},
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.Parts).ToNot(BeEmpty())
			Expect(len(resp.Parts)).To(BeNumerically("==", 5))
		})

		It("Должна вернуться 1 деталь, передан фильтр с 1 существующим UUID", func() {
			resp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
				Filter: &inventoryV1.PartsFilter{},
			})
			GinkgoWriter.Printf("Response parts: %v\n", resp.Parts)

			resp, err = inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
				Filter: &inventoryV1.PartsFilter{
					Uuids: []string{partUUIDs[0]},
				},
			},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp).ToNot(BeNil())
			GinkgoWriter.Printf("Filtered UUID: %s\n", partUUIDs[0])
			GinkgoWriter.Printf("all UUIDs: %v\n", partUUIDs)
			GinkgoWriter.Printf("Response parts count: %d\n", len(resp.Parts))
			Expect(resp.Parts).ToNot(BeEmpty())
			Expect(len(resp.Parts)).To(BeNumerically("==", 1))
			Expect(resp.Parts[0].GetUuid()).To(Equal(partUUIDs[0]))
		})

		It("Должен вернуться пустой список деталей, передан фильтр с 1 несуществующим UUID", func() {
			resp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
				Filter: &inventoryV1.PartsFilter{
					Uuids: []string{gofakeit.UUID()},
				},
			},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.Parts).To(BeEmpty())
		})

		It("Должна вернуться ошибка, передан nil фильтр", func() {
			resp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
				Filter: nil,
			})
			Expect(err).To(HaveOccurred())
			grpcStatus, ok := status.FromError(err)
			Expect(ok).To(BeTrue())
			Expect(grpcStatus.Code()).To(Equal(sharedErrors.ErrorCodeToGRPCCode(sharedErrors.InvalidArgumentErrCode)))
			Expect(resp).To(BeNil())
		})
	})
})
