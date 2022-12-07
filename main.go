package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
	"log"
	"msrd-products/db"
	_ "msrd-products/docs"
	"msrd-products/kafka/consumers"
	"msrd-products/routes"
	"msrd-products/utils"
	"os"
)

func main() {
	if os.Getenv("ENV") != "prod" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	dbContext, err := db.BuildDbContext(os.Getenv("DB_CONNECTION_STRING"), os.Getenv("DATABASE"), func(config *db.DbContextConfig) {})

	if err != nil {
		log.Fatal("Error loading database")
	}

	if os.Getenv("APP_MODE") == "STOCKS_CONSUMER" {
		consumers.LaunchProductStockRecordsConsumer(dbContext)
		return
	}

	if os.Getenv("APP_MODE") == "DOCUMENT_STATUS_CONSUMER" {
		consumers.LaunchDocumentStatusConsumer(dbContext)
		return
	}

	app := fiber.New()
	app.Use(recover.New())
	app.Use(func(c *fiber.Ctx) error {
		utils.SetLocal(c, "db_context", dbContext)
		return c.Next()
	})
	app.Use(cors.New())
	app.Use(logger.New())

	app.Get("/swagger/*", swagger.HandlerDefault)

	authMiddleware := jwtware.New(jwtware.Config{
		TokenLookup: "header:Authorization",
		AuthScheme:  "Bearer",
		KeySetURLs: []string{
			os.Getenv("JWSK_URL"),
		},
	})
	app.Use(authMiddleware)

	routes.SetupRoutes(app)

	log.Fatal(app.Listen(os.Getenv("HOST")))
}
