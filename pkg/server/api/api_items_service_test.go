package api_test

import (
	"context"
	"log/slog"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ya-breeze/diary.be/pkg/config"
	"github.com/ya-breeze/diary.be/pkg/database"
	"github.com/ya-breeze/diary.be/pkg/database/models"
	"github.com/ya-breeze/diary.be/pkg/generated/goserver"
	"github.com/ya-breeze/diary.be/pkg/server/api"
	"github.com/ya-breeze/diary.be/pkg/server/common"
)

// Helper function to create context with user ID for items tests
func createContextWithUserIDForItems(userID string) context.Context {
	ctx := context.Background()
	return context.WithValue(ctx, common.UserIDKey, userID)
}

var _ = Describe("ItemsAPIService", func() {
	var (
		service  goserver.ItemsAPIService
		logger   *slog.Logger
		storage  database.Storage
		ctx      context.Context
		userID   string
		testDate string
	)

	// Create context outside of BeforeEach to avoid fatcontext linting issue
	userID = "test-user-id"
	testDate = "2024-01-15"
	ctx = createContextWithUserIDForItems(userID)

	BeforeEach(func() {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		cfg := &config.Config{
			DBPath: ":memory:",
		}
		storage = database.NewStorage(logger, cfg)
		Expect(storage.Open()).To(Succeed())

		service = api.NewItemsAPIService(logger, storage)
	})

	AfterEach(func() {
		storage.Close()
	})

	Describe("GetItems", func() {
		Context("when no user ID in context", func() {
			It("should return 401 unauthorized", func() {
				emptyCtx := context.Background()
				response, err := service.GetItems(emptyCtx, testDate)
				Expect(err).ToNot(HaveOccurred())
				Expect(response.Code).To(Equal(401))
			})
		})

		Context("when item does not exist", func() {
			It("should return empty item with 200 status", func() {
				response, err := service.GetItems(ctx, testDate)
				Expect(err).ToNot(HaveOccurred())
				Expect(response.Code).To(Equal(200))

				itemsResponse, ok := response.Body.(goserver.ItemsResponse)
				Expect(ok).To(BeTrue())
				Expect(itemsResponse.Date).To(Equal(testDate))
				Expect(itemsResponse.Title).To(Equal(""))
				Expect(itemsResponse.Body).To(Equal(""))
				Expect(itemsResponse.Tags).To(Equal([]string{}))
				Expect(itemsResponse.PreviousDate).To(BeNil())
				Expect(itemsResponse.NextDate).To(BeNil())
			})
		})

		Context("when item exists", func() {
			BeforeEach(func() {
				// Create a test item
				testItem := &models.Item{
					UserID: userID,
					Date:   testDate,
					Title:  "Test Title",
					Body:   "Test Body Content",
					Tags:   models.StringList{"tag1", "tag2"},
				}
				Expect(storage.PutItem(userID, testItem)).To(Succeed())
			})

			It("should return the item with 200 status", func() {
				response, err := service.GetItems(ctx, testDate)
				Expect(err).ToNot(HaveOccurred())
				Expect(response.Code).To(Equal(200))

				itemsResponse, ok := response.Body.(goserver.ItemsResponse)
				Expect(ok).To(BeTrue())
				Expect(itemsResponse.Date).To(Equal(testDate))
				Expect(itemsResponse.Title).To(Equal("Test Title"))
				Expect(itemsResponse.Body).To(Equal("Test Body Content"))
				Expect(itemsResponse.Tags).To(Equal([]string{"tag1", "tag2"}))
			})
		})

		Context("when previous and next items exist", func() {
			BeforeEach(func() {
				// Create previous item
				prevItem := &models.Item{
					UserID: userID,
					Date:   "2024-01-14",
					Title:  "Previous Item",
					Body:   "Previous content",
				}
				Expect(storage.PutItem(userID, prevItem)).To(Succeed())

				// Create current item
				currentItem := &models.Item{
					UserID: userID,
					Date:   testDate,
					Title:  "Current Item",
					Body:   "Current content",
				}
				Expect(storage.PutItem(userID, currentItem)).To(Succeed())

				// Create next item
				nextItem := &models.Item{
					UserID: userID,
					Date:   "2024-01-16",
					Title:  "Next Item",
					Body:   "Next content",
				}
				Expect(storage.PutItem(userID, nextItem)).To(Succeed())
			})

			It("should include previous and next dates", func() {
				response, err := service.GetItems(ctx, testDate)
				Expect(err).ToNot(HaveOccurred())
				Expect(response.Code).To(Equal(200))

				itemsResponse, ok := response.Body.(goserver.ItemsResponse)
				Expect(ok).To(BeTrue())
				Expect(itemsResponse.Date).To(Equal(testDate))
				Expect(itemsResponse.PreviousDate).ToNot(BeNil())
				Expect(*itemsResponse.PreviousDate).To(Equal("2024-01-14"))
				Expect(itemsResponse.NextDate).ToNot(BeNil())
				Expect(*itemsResponse.NextDate).To(Equal("2024-01-16"))
			})
		})
	})
})
