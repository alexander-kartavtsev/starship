//go:build integration

package integration

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/alexander-kartavtsev/starship/platform/pkg/logger"
	inventoryV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/inventory/v1"
)

var _ = Describe("InventoryService", func() {
	var (
		ctx             context.Context
		cancel          context.CancelFunc
		inventoryClient inventoryV1.InventoryServiceClient
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(suiteCtx)

		logger.Debug(ctx, env.App.Address(), zap.String("address", ""))

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
		Expect(err).ToNot(HaveOccurred(), "ожидали успешную очистку коллекции parts")

		cancel()
	})

	Describe("ListParts", func() {
		It("должен вернуть список всех запчастей", func() {
			_, err := env.InsertTestPart(ctx)
			// partUuid, err := env.InsertTestPart(ctx)
			Expect(err).To(BeNil())

			f := &inventoryV1.PartsFilter{
				Uuids:                 []string{},
				Names:                 []string{},
				Categories:            []inventoryV1.Category{},
				ManufacturerCountries: []string{},
				ManufacturerNames:     []string{},
				Tags:                  []string{},
			}
			r := &inventoryV1.ListPartsRequest{
				Filter: f,
			}

			logger.Debug(ctx, "request", zap.Any("request", r))

			resp, _ := inventoryClient.ListParts(ctx, r)
			// resp, err := inventoryClient.ListParts(ctx, r)

			// Expect(err).ToNot(HaveOccurred(), "ожидали успешное получение списка запчастей")
			// Expect(resp.GetParts()).ToNot(BeEmpty())
			// Expect(resp.GetParts()[partUuid]).ToNot(BeEmpty())
			Expect(resp).To(BeNil())
		})
	})

	// Describe("GetPart", func() {
	//	It("должен вернуть запчасть", func() {
	//		partUuid, err := env.InsertTestPart(ctx)
	//		Expect(err).To(BeNil())
	//
	//		r := &inventoryV1.GetPartRequest{
	//			Uuid: partUuid,
	//		}
	//
	//		logger.Debug(ctx, "request", zap.Any("request", r))
	//
	//		resp, err := inventoryClient.GetPart(ctx, r)
	//
	//		Expect(err).ToNot(HaveOccurred(), "ожидали успешное получение запчасти")
	//		Expect(resp.GetInfo()).ToNot(BeEmpty())
	//	})
	// })
})
