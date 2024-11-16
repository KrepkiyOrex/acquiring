package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/pelletier/go-toml"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Book struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	Author    string `json:"author"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}

type Repository struct {
	DB *gorm.DB
}

type Config struct {
	Database struct {
		Host     string `toml:"host"`
		Port     string `toml:"port"`
		User     string `toml:"user"`
		Password string `toml:"password"`
		DBName   string `toml:"dbname"`
		SSLMode  string `toml:"sslmode"`
	} `toml:"database"`
}

func LoadConfig(filePath string) (*Config, error) {
	config := &Config{}
	tree, err := toml.LoadFile(filePath)
	if err != nil {
		return nil, err
	}
	err = tree.Unmarshal(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create_books", r.CreateBook)
	api.Delete("delete_book/:id", r.DeleteBook)
}

func main() {
	fmt.Println("Server started...")

	config, err := LoadConfig("config/configrep.toml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	dsn := "host=" + config.Database.Host +
		" user=" + config.Database.User +
		" password=" + config.Database.Password +
		" dbname=" + config.Database.DBName +
		" port=" + config.Database.Port +
		" sslmode=" + config.Database.SSLMode

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	err = db.AutoMigrate(&Book{})
	if err != nil {
		log.Fatalf("Could not migrate the database: %v", err)
	}

	r := Repository{DB: db}
	app := fiber.New()
	r.SetupRoutes(app)

	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func (r *Repository) CreateBook(context *fiber.Ctx) error {
	book := Book{}

	err := context.BodyParser(&book)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	err = r.DB.Create(&book).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create book"})
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "book has been added"})

	return nil
}

func (r *Repository) DeleteBook(context *fiber.Ctx) error {
	bookModel := Book{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "id cannot be empty"})
		return nil
	}

	err := r.DB.Delete(bookModel, id)
	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not delete book"})
		return err.Error
	}
	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "book delete successfully"})

	return nil
}
