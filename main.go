package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func init() {

	// zapLogger, _ := zap.NewDevelopment()
	// if !config.App.Debug {
	// 	zapLogger, _ = zap.NewProduction()
	// }

	// zap.ReplaceGlobals(zapLogger)
	// zap.S().Infow("Config logging enabled",
	// 	"Debug Mode", config.App.Debug,
	// 	"Log Config", zapLogger.Core(),
	// )

	// if config.App.Env != "development" {
	// 	sentryClient, err := sentry.NewClient(sentry.ClientOptions{
	// 		Dsn:              config.Log.SentryDSN,
	// 		Debug:            config.App.Debug,
	// 		AttachStacktrace: true,
	// 		EnableTracing:    true,
	// 		TracesSampleRate: 1.0,                   //send 100% of transactions
	// 		SampleRate:       1.0,                   // return 25% of errors
	// 		ServerName:       config.App.InstanceID, // os.Hostname()

	// 	})
	// 	if err != nil {
	// 		zap.S().Errorf("Failed to create sentry client: %s", err)
	// 	}
	// 	defer sentryClient.Flush(2 * time.Second)

	// 	core, err := zapsentry.NewCore(
	// 		zapsentry.Configuration{
	// 			Level:             zapcore.InfoLevel,
	// 			EnableBreadcrumbs: true,
	// 			BreadcrumbLevel:   zapcore.InfoLevel,
	// 			Tags: map[string]string{
	// 				"component": "live-api",
	// 			}},
	// 		zapsentry.NewSentryClientFromClient(sentryClient))

	// 	if err != nil {
	// 		zap.S().Warn("Failed to init zap", zap.Error(err))
	// 	}

	// 	zapLogger = zapsentry.AttachCoreToLogger(core, zapLogger)
	// 	zap.ReplaceGlobals(zapLogger)
	// 	zap.S().Info("Test Sentry")
	// }

}

var config *Config

func main() {
	cfg, err := Configs(".env")
	if err != nil {
		log.Fatalf("Error while loading configurations : %v", err)
	}
	config = cfg
	// HelloWorld()
	// OTP()
	// OTPWithRedis()
	// OTPWithRedisAndPostgres()
	// OTPWithRedisAndPostgresAndSMS(context.Background())
	// OTPWithRedisAndPostgresMemoized(context.Background())
	OTPWithRedisAndPostgresAndSMSMemoized(context.Background())
}

func HelloWorld() {
	app := fiber.New()

	app.Get("/hi", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}

func OTP() {
	app := fiber.New()

	app.Get("/hi", func(c *fiber.Ctx) error {
		// Generate a random 6-digit OTP code
		otpCode := fmt.Sprintf("%06d", rand.Intn(1000000))

		return c.JSON(fiber.Map{
			"code": otpCode,
		})
	})

	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}

func OTPWithRedis() {
	app := fiber.New()

	// Create context
	ctx := context.Background()

	// Connect to Redis
	redisClient, err := RegisterRedis(ctx, config)
	if err != nil {
		zap.L().Fatal("Failed to connect to Redis", zap.Error(err))
	}
	defer redisClient.Close()

	app.Get("/hi", func(c *fiber.Ctx) error {
		// Generate random OTP
		otpCode := fmt.Sprintf("%06d", rand.Intn(1000000))

		// Save to Redis with 5 minute expiration
		err := redisClient.Set(ctx, "otp:"+otpCode, otpCode, 5*time.Minute).Err()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to save OTP",
			})
		}

		return c.JSON(fiber.Map{
			"code": otpCode,
		})
	})

	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}

func OTPWithRedisAndPostgres() {
	app := fiber.New()

	// Create context
	ctx := context.Background()

	// Connect to Redis
	redisClient, err := RegisterRedis(ctx, config)
	if err != nil {
		zap.L().Fatal("Failed to connect to Redis", zap.Error(err))
	}
	defer redisClient.Close()

	// Connect to Postgres
	db, err := RegisterPostgres(ctx, config)
	if err != nil {
		zap.L().Fatal("Failed to connect to Postgres", zap.Error(err))
	}
	defer db.Close()

	app.Get("/hi", func(c *fiber.Ctx) error {
		// Get random tenant info from Postgres
		var (
			tenantID          string
			expirationMinutes int
		)
		err := db.QueryRowContext(ctx,
			"SELECT id, otp_expiration_minutes FROM tenants ORDER BY RAND() LIMIT 1",
		).Scan(&tenantID, &expirationMinutes)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get tenant info",
			})
		}

		// Generate random OTP
		otpCode := fmt.Sprintf("%06d", rand.Intn(1000000))

		// Save to Redis with tenant-specific expiration
		err = redisClient.Set(ctx,
			fmt.Sprintf("otp:%s:%s", tenantID, otpCode),
			otpCode,
			time.Duration(expirationMinutes)*time.Minute,
		).Err()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to save OTP",
			})
		}

		return c.JSON(fiber.Map{
			"code":       otpCode,
			"tenant_id":  tenantID,
			"expires_in": fmt.Sprintf("%d minutes", expirationMinutes),
		})
	})

	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}

func OTPWithRedisAndPostgresAndSMS(ctx context.Context) {
	app := fiber.New()

	// Connect to Redis
	redisClient, err := RegisterRedis(ctx, config)
	if err != nil {
		zap.L().Fatal("Failed to connect to Redis", zap.Error(err))
	}
	defer redisClient.Close()

	// Connect to Postgres
	db, err := RegisterPostgres(ctx, config)
	if err != nil {
		zap.L().Fatal("Failed to connect to Postgres", zap.Error(err))
	}
	defer db.Close()

	app.Get("/hi", func(c *fiber.Ctx) error {
		// Get random tenant info from Postgres
		var (
			tenantID          string
			expirationMinutes int
		)
		err := db.QueryRowContext(ctx,
			"SELECT id, otp_expiration_minutes FROM tenants ORDER BY RAND() LIMIT 1",
		).Scan(&tenantID, &expirationMinutes)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get tenant info",
			})
		}

		// Generate random OTP
		otpCode := fmt.Sprintf("%06d", rand.Intn(1000000))

		// Save to Redis with tenant-specific expiration
		err = redisClient.Set(ctx,
			fmt.Sprintf("otp:%s:%s", tenantID, otpCode),
			otpCode,
			time.Duration(expirationMinutes)*time.Minute,
		).Err()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to save OTP",
			})
		}

		// Simulate SMS sending delay
		time.Sleep(2 * time.Second)
		zap.L().Info("SMS sent successfully",
			zap.String("tenant_id", tenantID),
			zap.String("otp", otpCode),
		)

		return c.JSON(fiber.Map{
			"message":    "OTP sent via SMS",
			"tenant_id":  tenantID,
			"expires_in": fmt.Sprintf("%d minutes", expirationMinutes),
		})
	})

	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}

// RegisterRistretto initializes and returns a Ristretto cache instance.
func RegisterRistretto() *ristretto.Cache {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e4,     // number of keys to track frequency of (10k).
		MaxCost:     1 << 20, // maximum cost of cache (1MB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil {
		log.Fatalf("failed to create ristretto cache: %v", err)
	}
	return cache
}

// OTPWithRedisAndPostgresMemoized is like OTPWithRedisAndPostgres but uses Ristretto for memoizing tenant info from Postgres.
// OTPWithRedisAndPostgresMemoized builds the Fiber app and serves it (no return).
func OTPWithRedisAndPostgresMemoized(ctx context.Context) {
	app := fiber.New()

	// Register Redis, Postgres, and Ristretto cache
	redisClient, err := RegisterRedis(ctx, config)
	if err != nil {
		log.Fatalf("failed to register redis: %v", err)
	}
	db, err := RegisterPostgres(ctx, config)
	if err != nil {
		log.Fatalf("failed to register postgres: %v", err)
	}
	cache := RegisterRistretto()

	app.Get("/hi", func(c *fiber.Ctx) error {
		var (
			tenantID          string
			expirationMinutes int
		)

		// Try to get tenant info from cache
		cacheKey := "random_tenant"
		if val, found := cache.Get(cacheKey); found {
			tenant := val.(map[string]interface{})
			tenantID = tenant["id"].(string)
			expirationMinutes = tenant["otp_expiration_minutes"].(int)
		} else {
			// Not in cache, query Postgres
			err := db.QueryRowContext(ctx,
				"SELECT id, otp_expiration_minutes FROM tenants ORDER BY RAND() LIMIT 1",
			).Scan(&tenantID, &expirationMinutes)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to get tenant info",
				})
			}
			// Store in cache
			cache.Set(cacheKey, map[string]interface{}{
				"id":                     tenantID,
				"otp_expiration_minutes": expirationMinutes,
			}, 1)
		}

		// Generate random OTP
		otpCode := fmt.Sprintf("%06d", rand.Intn(1000000))

		// Save to Redis with tenant-specific expiration
		err := redisClient.Set(ctx,
			fmt.Sprintf("otp:%s:%s", tenantID, otpCode),
			otpCode,
			time.Duration(expirationMinutes)*time.Minute,
		).Err()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to save OTP",
			})
		}

		return c.JSON(fiber.Map{
			"message":    "OTP sent via SMS (memoized)",
			"tenant_id":  tenantID,
			"expires_in": fmt.Sprintf("%d minutes", expirationMinutes),
		})
	})

	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}

// OTPWithRedisAndPostgresMemoized is like OTPWithRedisAndPostgres but uses Ristretto for memoizing tenant info from Postgres.
func OTPWithRedisAndPostgresAndSMSMemoized(ctx context.Context) {
	// Create Fiber app
	app := fiber.New()

	// Register Redis, Postgres, and Ristretto cache
	redisClient, err := RegisterRedis(ctx, config)
	if err != nil {
		log.Fatalf("failed to register redis: %v", err)
	}
	db, err := RegisterPostgres(ctx, config)
	if err != nil {
		log.Fatalf("failed to register postgres: %v", err)
	}
	cache := RegisterRistretto()

	app.Get("/hi", func(c *fiber.Ctx) error {
		var (
			tenantID          string
			expirationMinutes int
		)

		// Try to get tenant info from cache
		cacheKey := "random_tenant"
		if val, found := cache.Get(cacheKey); found {
			tenant := val.(map[string]interface{})
			tenantID = tenant["id"].(string)
			expirationMinutes = tenant["otp_expiration_minutes"].(int)
		} else {
			// Not in cache, query Postgres
			err := db.QueryRowContext(ctx,
				"SELECT id, otp_expiration_minutes FROM tenants ORDER BY RAND() LIMIT 1",
			).Scan(&tenantID, &expirationMinutes)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to get tenant info",
				})
			}
			// Store in cache
			cache.Set(cacheKey, map[string]interface{}{
				"id":                     tenantID,
				"otp_expiration_minutes": expirationMinutes,
			}, 1)
		}

		// Generate random OTP
		otpCode := fmt.Sprintf("%06d", rand.Intn(1000000))

		// Save to Redis with tenant-specific expiration
		err := redisClient.Set(ctx,
			fmt.Sprintf("otp:%s:%s", tenantID, otpCode),
			otpCode,
			time.Duration(expirationMinutes)*time.Minute,
		).Err()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to save OTP",
			})
		}

		// Simulate SMS sending delay
		time.Sleep(2 * time.Second)
		zap.L().Info("SMS sent successfully (memoized)",
			zap.String("tenant_id", tenantID),
			zap.String("otp", otpCode),
		)

		return c.JSON(fiber.Map{
			"message":    "OTP sent via SMS (memoized)",
			"tenant_id":  tenantID,
			"expires_in": fmt.Sprintf("%d minutes", expirationMinutes),
		})
	})

	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}
