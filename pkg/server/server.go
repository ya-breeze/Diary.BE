package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/ya-breeze/diary.be/pkg/auth"
	"github.com/ya-breeze/diary.be/pkg/config"
	"github.com/ya-breeze/diary.be/pkg/constants"
	"github.com/ya-breeze/diary.be/pkg/database"
	"github.com/ya-breeze/diary.be/pkg/database/models"
	"github.com/ya-breeze/diary.be/pkg/generated/goserver"
	"github.com/ya-breeze/diary.be/pkg/server/api"
	"github.com/ya-breeze/diary.be/pkg/server/background"
	"github.com/ya-breeze/diary.be/pkg/server/webapp"
)

func Server(logger *slog.Logger, cfg *config.Config) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	storage := database.NewStorage(logger, cfg)
	if err := storage.Open(); err != nil {
		return fmt.Errorf("failed to open storage: %w", err)
	}

	forcedImportChan := make(chan background.ForcedImport, 100)
	_, finishChan, err := Serve(ctx, logger, storage, cfg, forcedImportChan)
	if err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	// Start bank importers
	var importFinishChan <-chan struct{}
	if !cfg.DisableImporters {
		importFinishChan = background.StartBankImporters(ctx, logger, storage, forcedImportChan)
	} else {
		logger.Info("Bank importers are disabled")
	}

	// Start currency rate fetcher
	var fetchCurrenciesRatesChan <-chan struct{}
	if !cfg.DisableCurrenciesRatesFetch {
		fetchCurrenciesRatesChan = background.StartCurrenciesRatesFetcher(ctx, logger, storage)
	} else {
		logger.Info("Fetcher for currencies rates is disabled")
	}

	// Wait for an interrupt signal
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-stopChan
	logger.Info("Received signal. Shutting down server...")

	// Stop the server
	cancel()
	<-finishChan
	if !cfg.DisableImporters {
		<-importFinishChan
	}
	if !cfg.DisableCurrenciesRatesFetch {
		<-fetchCurrenciesRatesChan
	}
	return nil
}

func createControllers(logger *slog.Logger, cfg *config.Config, db database.Storage) goserver.CustomControllers {
	return goserver.CustomControllers{
		AuthAPIService:                    api.NewAuthAPIService(logger, db, cfg.JWTSecret),
		UserAPIService:                    api.NewUserAPIService(logger, db),
		AccountsAPIService:                api.NewAccountsAPIService(logger, db),
		CurrenciesAPIService:              api.NewCurrenciesAPIServicer(logger, db),
		TransactionsAPIService:            api.NewTransactionsAPIService(logger, db),
		UnprocessedTransactionsAPIService: api.NewUnprocessedTransactionsAPIServiceImpl(logger, db),
		MatchersAPIService:                api.NewMatchersAPIServiceImpl(logger, db),
		BankImportersAPIService:           api.NewBankImportersAPIServiceImpl(logger, db),
		AggregationsAPIService:            api.NewAggregationsAPIServiceImpl(logger, db),
		NotificationsAPIService:           api.NewNotificationsAPIServiceImpl(logger, db),
		ImportAPIService:                  api.NewImportAPIServiceImpl(logger, db),
		ExportAPIService:                  api.NewExportAPIServiceImpl(logger, db),
	}
}

func Serve(
	ctx context.Context, logger *slog.Logger,
	storage database.Storage, cfg *config.Config,
	forcedImports chan<- background.ForcedImport,
) (net.Addr, chan int, error) {
	commit := func() string {
		if info, ok := debug.ReadBuildInfo(); ok {
			for _, setting := range info.Settings {
				if setting.Key == "vcs.revision" {
					return setting.Value
				}
			}
		}

		return ""
	}()
	logger.Info("Built from git commit: " + commit)

	if cfg.JWTSecret == "" {
		logger.Warn("JWT secret is not set. Creating random secret...")
		cfg.JWTSecret = auth.GenerateRandomString(32)
	}

	logger.Info("Starting GeekBudget server...")

	if cfg.Users != "" {
		logger.Info("Creating users...")
		users := strings.Split(cfg.Users, ",")
		for _, user := range users {
			tokens := strings.Split(user, ":")
			if len(tokens) != 2 {
				return nil, nil, fmt.Errorf("invalid user format: %s", user)
			}

			if err := upsertUser(storage, tokens[0], tokens[1], logger, cfg.Prefill); err != nil {
				return nil, nil, fmt.Errorf("failed to update user %q: %w", tokens[0], err)
			}
		}
	} else {
		logger.Info("No users defined in configuration")
	}

	return goserver.Serve(ctx, logger, cfg,
		createControllers(logger, cfg, storage),
		[]goserver.Router{webapp.NewWebAppRouter(commit, logger, cfg, storage)},
		createMiddlewares(logger, cfg, forcedImports)...)
}

func upsertUser(storage database.Storage, username, hashedPassword string, logger *slog.Logger, prefill bool) error {
	userID, err := storage.GetUserID(username)
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return fmt.Errorf("failed to reading user from DB: %w", err)
	}
	var user *models.User
	if !errors.Is(err, database.ErrNotFound) {
		logger.Info(fmt.Sprintf("Updating password for user %q", username))

		user, err = storage.GetUser(userID)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}
		user.HashedPassword = hashedPassword
		if err = storage.PutUser(user); err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}
	} else {
		logger.Info(fmt.Sprintf("Creating user %q", username))
		user, err = storage.CreateUser(username, hashedPassword)
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		if prefill {
			err = prefillNewUser(storage, user.ID.String(), logger)
			if err != nil {
				return fmt.Errorf("failed to prefill new user: %w", err)
			}
		}
	}

	return nil
}

//nolint:funlen,cyclop,maintidx // function is not complex, it just creates many default items
func prefillNewUser(storage database.Storage, userID string, logger *slog.Logger) error {
	// Create default currencies
	currency := &goserver.CurrencyNoId{
		Name: "USD",
	}
	curUSD, err := storage.CreateCurrency(userID, currency)
	if err != nil {
		return fmt.Errorf("failed to create USD currency: %w", err)
	}

	currency = &goserver.CurrencyNoId{
		Name: "EUR",
	}
	curEUR, err := storage.CreateCurrency(userID, currency)
	if err != nil {
		return fmt.Errorf("failed to create EUR currency: %w", err)
	}

	currency = &goserver.CurrencyNoId{
		Name: "CZK",
	}
	curCZK, err := storage.CreateCurrency(userID, currency)
	if err != nil {
		return fmt.Errorf("failed to create CZK currency: %w", err)
	}

	// Create default accounts
	account := &goserver.AccountNoId{
		Name:        "Cash",
		Description: "Cash account",
		Type:        constants.AccountAsset,
		BankInfo: goserver.BankAccountInfo{
			Balances: []goserver.BankAccountInfoBalancesInner{
				{
					CurrencyId:     curCZK.Id,
					OpeningBalance: 1000,
				},
				{
					CurrencyId:     curUSD.Id,
					OpeningBalance: 1,
				},
				{
					CurrencyId:     curEUR.Id,
					OpeningBalance: 2,
				},
			},
		},
	}
	accCash, err := storage.CreateAccount(userID, account)
	if err != nil {
		return fmt.Errorf("failed to create cash account: %w", err)
	}

	account = &goserver.AccountNoId{
		Name:        "FIO bank",
		Description: "FIO bank account",
		Type:        constants.AccountAsset,
	}
	accFio, err := storage.CreateAccount(userID, account)
	if err != nil {
		return fmt.Errorf("failed to create FIO bank account: %w", err)
	}

	account = &goserver.AccountNoId{
		Name:        "Revolut",
		Description: "Revolut bank account",
		Type:        constants.AccountAsset,
	}
	accRevolut, err := storage.CreateAccount(userID, account)
	if err != nil {
		return fmt.Errorf("failed to create Revolut bank account: %w", err)
	}

	account = &goserver.AccountNoId{
		Name:        "Salary",
		Description: "Salary account",
		Type:        constants.AccountIncome,
	}
	accSalary, err := storage.CreateAccount(userID, account)
	if err != nil {
		return fmt.Errorf("failed to create income account: %w", err)
	}
	logger.Info(fmt.Sprintf("Created income account %q", accSalary.Id))

	account = &goserver.AccountNoId{
		Name: "Bank fees",
		Type: constants.AccountExpense,
	}
	accFees, err := storage.CreateAccount(userID, account)
	if err != nil {
		return fmt.Errorf("failed to create fees account: %w", err)
	}
	logger.Info(fmt.Sprintf("Created fees account %q", accFees.Id))

	account = &goserver.AccountNoId{
		Name: "🛒 Groceries",
		Type: constants.AccountExpense,
	}
	accGroceries, err := storage.CreateAccount(userID, account)
	if err != nil {
		return fmt.Errorf("failed to create groceries account: %w", err)
	}
	logger.Info(fmt.Sprintf("Created groceries account %q", accGroceries.Id))

	account = &goserver.AccountNoId{
		Name:        "Transport",
		Description: "Transport expenses account",
		Type:        constants.AccountExpense,
	}
	accTransport, err := storage.CreateAccount(userID, account)
	if err != nil {
		return fmt.Errorf("failed to create transport account: %w", err)
	}
	logger.Info(fmt.Sprintf("Created transport account %q", accTransport.Id))

	account = &goserver.AccountNoId{
		Name: "🏠 Rent",
		Type: constants.AccountExpense,
	}
	accRent, err := storage.CreateAccount(userID, account)
	if err != nil {
		return fmt.Errorf("failed to create rent account: %w", err)
	}
	logger.Info(fmt.Sprintf("Created rent account %q", accRent.Id))

	// Create default transactions
	transaction := &goserver.TransactionNoId{
		Date:        time.Now().Add(-5 * 24 * time.Hour),
		Description: "Initial state for cash",
		Tags:        []string{"initial_account_state"},
		Movements: []goserver.Movement{
			{
				AccountId:  accCash.Id,
				Amount:     1000,
				CurrencyId: curCZK.Id,
			},
			{
				AccountId:  accCash.Id,
				Amount:     1,
				CurrencyId: curUSD.Id,
			},
			{
				AccountId:  accCash.Id,
				Amount:     1,
				CurrencyId: curEUR.Id,
			},
		},
	}
	if _, err := storage.CreateTransaction(userID, transaction); err != nil {
		return fmt.Errorf("failed to create initial cash transaction: %w", err)
	}

	transaction = &goserver.TransactionNoId{
		Date:        time.Now().Add(-5 * 24 * time.Hour),
		Description: "Initial state for bank",
		Tags:        []string{"initial_account_state"},
		Movements: []goserver.Movement{
			{
				AccountId:  accFio.Id,
				Amount:     10000,
				CurrencyId: curCZK.Id,
			},
		},
	}
	if _, err := storage.CreateTransaction(userID, transaction); err != nil {
		return fmt.Errorf("failed to create initial bank transaction: %w", err)
	}

	transaction = &goserver.TransactionNoId{
		Date:        time.Now().Add(-24 * time.Hour),
		Description: "Monthly salary",
		Movements: []goserver.Movement{
			{
				AccountId:  accSalary.Id,
				Amount:     -10000,
				CurrencyId: curCZK.Id,
			},
			{
				AccountId:  accFio.Id,
				Amount:     10000,
				CurrencyId: curCZK.Id,
			},
		},
	}
	if _, err := storage.CreateTransaction(userID, transaction); err != nil {
		return fmt.Errorf("failed to create salary transaction: %w", err)
	}

	transaction = &goserver.TransactionNoId{
		Date:        time.Now().Add(-4 * time.Hour),
		Description: "Lunch",
		Movements: []goserver.Movement{
			{
				AccountId:  accGroceries.Id,
				Amount:     100,
				CurrencyId: curCZK.Id,
			},
			{
				AccountId:  accCash.Id,
				Amount:     -100,
				CurrencyId: curCZK.Id,
			},
		},
	}
	if _, err := storage.CreateTransaction(userID, transaction); err != nil {
		return fmt.Errorf("failed to create lunch transaction: %w", err)
	}

	transaction = &goserver.TransactionNoId{
		Date:        time.Now().Add(-24 * time.Hour),
		Description: "Lunch",
		Movements: []goserver.Movement{
			{
				AccountId:  accGroceries.Id,
				Amount:     150,
				CurrencyId: curCZK.Id,
			},
			{
				AccountId:  accCash.Id,
				Amount:     -150,
				CurrencyId: curCZK.Id,
			},
		},
	}
	if _, err := storage.CreateTransaction(userID, transaction); err != nil {
		return fmt.Errorf("failed to create lunch transaction: %w", err)
	}

	transaction = &goserver.TransactionNoId{
		Date:        time.Now().Add(-3 * time.Hour),
		Description: "Bus ticket",
		Movements: []goserver.Movement{
			{
				AccountId:  accTransport.Id,
				Amount:     25,
				CurrencyId: curCZK.Id,
			},
			{
				AccountId:  accCash.Id,
				Amount:     -25,
				CurrencyId: curCZK.Id,
			},
		},
	}
	if _, err := storage.CreateTransaction(userID, transaction); err != nil {
		return fmt.Errorf("failed to create bus ticket transaction: %w", err)
	}

	transaction = &goserver.TransactionNoId{
		Date:        time.Now().Add(-2 * time.Hour),
		Description: "Apartment rent",
		Movements: []goserver.Movement{
			{
				AccountId:  accRent.Id,
				Amount:     5000,
				CurrencyId: curCZK.Id,
			},
			{
				AccountId:  accFio.Id,
				Amount:     -5000,
				CurrencyId: curCZK.Id,
			},
		},
	}
	if _, err := storage.CreateTransaction(userID, transaction); err != nil {
		return fmt.Errorf("failed to create rent transaction: %w", err)
	}

	transaction = &goserver.TransactionNoId{
		Date:        time.Now().Add(-4 * time.Hour),
		Description: "USD spending",
		Movements: []goserver.Movement{
			{
				AccountId:  accGroceries.Id,
				Amount:     1,
				CurrencyId: curUSD.Id,
			},
			{
				AccountId:  accFio.Id,
				Amount:     -25,
				CurrencyId: curCZK.Id,
			},
		},
	}
	if _, err := storage.CreateTransaction(userID, transaction); err != nil {
		return fmt.Errorf("failed to create lunch transaction: %w", err)
	}

	// Bank importers
	bankImporter := &goserver.BankImporterNoId{
		Name:        "FIO Bank CZK",
		Description: "Fio banka a.s. (CZK)",
		Extra:       "token",
		AccountId:   accFio.Id,
		Type:        "fio",
		Mappings: []goserver.BankImporterNoIdMappingsInner{
			{
				FieldToMatch: "user",
				ValueToMatch: "Korolev, Ilya",
				TagToSet:     "ilya",
			},
			{
				FieldToMatch: "user",
				ValueToMatch: "Koroleva, Anzhela",
				TagToSet:     "angela",
			},
		},
	}
	if _, err := storage.CreateBankImporter(userID, bankImporter); err != nil {
		return fmt.Errorf("failed to create bank importer: %w", err)
	}

	bankImporter = &goserver.BankImporterNoId{
		Name:         "Revolut bank",
		Description:  "Revolut bank",
		AccountId:    accRevolut.Id,
		FeeAccountId: accFees.Id,
		Type:         "revolut",
	}
	if _, err := storage.CreateBankImporter(userID, bankImporter); err != nil {
		return fmt.Errorf("failed to create bank importer: %w", err)
	}

	matcher := &goserver.MatcherNoId{
		Name: "Groceries: Billa",

		OutputDescription: "Groceries: Billa",
		OutputAccountId:   accGroceries.Id,

		DescriptionRegExp: `(?i)\bBilla\b`,
	}
	if _, err := storage.CreateMatcher(userID, matcher); err != nil {
		return fmt.Errorf("failed to create matcher: %w", err)
	}

	matcher = &goserver.MatcherNoId{
		Name: "Groceries: Albert",

		OutputDescription: "Groceries: Albert",
		OutputAccountId:   accGroceries.Id,

		DescriptionRegExp: `(?i)\bAlbert\b`,
	}
	if _, err := storage.CreateMatcher(userID, matcher); err != nil {
		return fmt.Errorf("failed to create matcher: %w", err)
	}

	matcher = &goserver.MatcherNoId{
		Name: "Groceries: Makro",

		OutputDescription: "Groceries: Makro",
		OutputAccountId:   accGroceries.Id,

		DescriptionRegExp: `(?i)\bMakro 06\b`,
	}
	if _, err := storage.CreateMatcher(userID, matcher); err != nil {
		return fmt.Errorf("failed to create matcher: %w", err)
	}

	matcher = &goserver.MatcherNoId{
		Name: "Groceries: Globus",

		OutputDescription: "Groceries: Globus",
		OutputAccountId:   accGroceries.Id,

		DescriptionRegExp: `(?i)\bGlobus\b`,
	}
	if _, err := storage.CreateMatcher(userID, matcher); err != nil {
		return fmt.Errorf("failed to create matcher: %w", err)
	}

	matcher = &goserver.MatcherNoId{
		Name: "Groceries: Kaufland",

		OutputDescription: "Groceries: Kaufland",
		OutputAccountId:   accGroceries.Id,

		DescriptionRegExp: `(?i)\bKaufland\b`,
	}
	if _, err := storage.CreateMatcher(userID, matcher); err != nil {
		return fmt.Errorf("failed to create matcher: %w", err)
	}

	return nil
}

func createMiddlewares(
	logger *slog.Logger, cfg *config.Config, forcedImports chan<- background.ForcedImport,
) []mux.MiddlewareFunc {
	return []mux.MiddlewareFunc{
		AuthMiddleware(logger, cfg),
		ForcedImportMiddleware(logger, forcedImports),
	}
}
