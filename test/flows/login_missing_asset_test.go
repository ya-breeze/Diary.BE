package flows_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ya-breeze/diary.be/pkg/auth"
	"github.com/ya-breeze/diary.be/pkg/config"
	"github.com/ya-breeze/diary.be/pkg/database"
	"github.com/ya-breeze/diary.be/pkg/generated/goclient"
	"github.com/ya-breeze/diary.be/pkg/server"
)

var _ = Describe("Login and Missing Asset Flow", func() {
	var (
		logger     *slog.Logger
		cfg        *config.Config
		storage    database.Storage
		serverAddr string
		apiClient  *goclient.APIClient
		ctx        context.Context
		cancel     context.CancelFunc
		testEmail  string
		testPass   string
		tempDir    string
	)

	BeforeEach(func() {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

		// Create temporary directory for assets
		var err error
		tempDir, err = os.MkdirTemp("", "flow_test_assets")
		Expect(err).NotTo(HaveOccurred())

		cfg = &config.Config{
			Port:      0, // Use random available port
			DBPath:    ":memory:",
			AssetPath: tempDir,
			Issuer:    "test-issuer",
			JWTSecret: "test-secret-key-for-jwt-tokens",
		}

		storage = database.NewStorage(logger, cfg)
		Expect(storage.Open()).To(Succeed())

		// Create test user
		testEmail = "test@example.com"
		testPass = "testpassword123"

		hashedPassBytes, err := auth.HashPassword([]byte(testPass))
		Expect(err).ToNot(HaveOccurred())
		hashedPass := base64.StdEncoding.EncodeToString(hashedPassBytes)

		_, err = storage.CreateUser(testEmail, hashedPass)
		Expect(err).ToNot(HaveOccurred())

		// Start test server
		ctx, cancel = context.WithCancel(context.Background())
		addr, _, err := server.Serve(ctx, logger, storage, cfg)
		Expect(err).ToNot(HaveOccurred())

		serverAddr = fmt.Sprintf("http://localhost:%d", addr.(*net.TCPAddr).Port)
		logger.Info("Test server started", "address", serverAddr)

		// Wait for server to be ready by polling the authorize endpoint
		Eventually(func() bool {
			resp, err := http.Post(serverAddr+"/v1/authorize", "application/json", nil)
			if err != nil {
				return false
			}
			defer resp.Body.Close()
			// We expect 400 (bad request) because we're not sending valid JSON,
			// but this means the server is up and responding
			return resp.StatusCode == http.StatusBadRequest
		}, "5s", "100ms").Should(BeTrue())

		// Create API client
		clientConfig := goclient.NewConfiguration()
		clientConfig.Servers = goclient.ServerConfigurations{
			{
				URL:         serverAddr,
				Description: "Test server",
			},
		}

		apiClient = goclient.NewAPIClient(clientConfig)
	})

	AfterEach(func() {
		if cancel != nil {
			cancel()
		}
		if storage != nil {
			storage.Close()
		}
		if tempDir != "" {
			os.RemoveAll(tempDir)
		}
	})

	Describe("Authentication and Asset Access Flow", func() {
		Context("when user logs in successfully", func() {
			It("should authenticate and then receive 404 for missing asset", func() {
				// Step 1: Login via API
				authData := goclient.AuthData{
					Email:    testEmail,
					Password: testPass,
				}

				authResponse, httpResponse, err := apiClient.AuthAPI.Authorize(context.Background()).AuthData(authData).Execute()
				Expect(err).ToNot(HaveOccurred())
				Expect(httpResponse.StatusCode).To(Equal(http.StatusOK))
				Expect(authResponse.Token).ToNot(BeEmpty())

				// Step 2: Configure client with JWT token for subsequent requests
				clientConfig := apiClient.GetConfig()
				clientConfig.AddDefaultHeader("Authorization", "Bearer "+authResponse.Token)

				// Step 3: Try to get a missing asset
				missingAssetPath := "nonexistent/missing-image.jpg"

				_, httpResponse, err = apiClient.AssetsAPI.GetAsset(context.Background()).Path(missingAssetPath).Execute()

				// We expect this to fail with 404
				Expect(err).To(HaveOccurred())
				Expect(httpResponse.StatusCode).To(Equal(http.StatusNotFound))

				// Verify the error is a GenericOpenAPIError with 404 status
				if openAPIErr, ok := err.(*goclient.GenericOpenAPIError); ok {
					Expect(openAPIErr.Error()).To(ContainSubstring("404"))
				}
			})
		})

		Context("when user tries to access asset without authentication", func() {
			It("should receive 401 unauthorized", func() {
				// Try to get an asset without authentication
				missingAssetPath := "some-asset.jpg"

				_, httpResponse, err := apiClient.AssetsAPI.GetAsset(context.Background()).Path(missingAssetPath).Execute()

				// We expect this to fail with 401
				Expect(err).To(HaveOccurred())
				Expect(httpResponse.StatusCode).To(Equal(http.StatusUnauthorized))
			})
		})

		Context("when user provides invalid credentials", func() {
			It("should receive 401 authentication failed", func() {
				// Try to login with invalid credentials
				authData := goclient.AuthData{
					Email:    testEmail,
					Password: "wrongpassword",
				}

				_, httpResponse, err := apiClient.AuthAPI.Authorize(context.Background()).AuthData(authData).Execute()

				// We expect this to fail with 401
				Expect(err).To(HaveOccurred())
				Expect(httpResponse.StatusCode).To(Equal(http.StatusUnauthorized))
			})
		})

		Context("when user logs in and fetches an existing asset", func() {
			It("should successfully retrieve the asset", func() {
				// Step 1: Login via API
				authData := goclient.AuthData{
					Email:    testEmail,
					Password: testPass,
				}

				authResponse, httpResponse, err := apiClient.AuthAPI.Authorize(context.Background()).AuthData(authData).Execute()
				Expect(err).ToNot(HaveOccurred())
				Expect(httpResponse.StatusCode).To(Equal(http.StatusOK))
				Expect(authResponse.Token).ToNot(BeEmpty())

				// Step 2: Configure client with JWT token for subsequent requests
				clientConfig := apiClient.GetConfig()
				clientConfig.AddDefaultHeader("Authorization", "Bearer "+authResponse.Token)

				// Step 3: Create a test asset file in the user's directory
				// First, we need to get the user ID from the JWT token to create the correct directory structure
				userID, err := auth.CheckJWT(authResponse.Token, cfg.Issuer, cfg.JWTSecret)
				Expect(err).ToNot(HaveOccurred())

				userAssetDir := filepath.Join(tempDir, userID)
				err = os.MkdirAll(userAssetDir, 0755)
				Expect(err).ToNot(HaveOccurred())

				testAssetPath := "images/photos/test-photo.jpg"
				testAssetFullPath := filepath.Join(userAssetDir, testAssetPath)
				testAssetDir := filepath.Dir(testAssetFullPath)
				err = os.MkdirAll(testAssetDir, 0755)
				Expect(err).ToNot(HaveOccurred())

				testAssetContent := []byte("fake image content for testing")
				err = os.WriteFile(testAssetFullPath, testAssetContent, 0644)
				Expect(err).ToNot(HaveOccurred())

				// Step 4: Fetch the existing asset
				assetFile, httpResponse, err := apiClient.AssetsAPI.GetAsset(context.Background()).Path(testAssetPath).Execute()

				// We expect this to succeed with 200
				Expect(err).ToNot(HaveOccurred())
				Expect(httpResponse.StatusCode).To(Equal(http.StatusOK))
				Expect(assetFile).ToNot(BeNil())

				// Step 5: Verify the content (optional - read and compare)
				defer assetFile.Close()
				retrievedContent, err := io.ReadAll(assetFile)
				Expect(err).ToNot(HaveOccurred())
				Expect(retrievedContent).To(Equal(testAssetContent))
			})
		})
	})
})
