package db

import (
	"fmt"
	"log"
	"os"

	"social_graph_api/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Read env variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	fmt.Println("✅ Connected to the database successfully!")

	DB = database

	err = database.AutoMigrate(&models.User{}, &models.Connection{})
	if err != nil {
		log.Fatal("❌ Failed to migrate User model: ", err)
	}

	fmt.Println("✅ User model migrated successfully!")

	var count int64
	DB.Model(&models.User{}).Count(&count)
	if count == 0 {
		users := []models.User{
			{Name: "Alice", Email: "alice@example.com"},
			{Name: "Bob", Email: "bob@example.com"},
			{Name: "Charlie", Email: "charlie@example.com"},
			{Name: "Lily", Email: "lily@example.com"},
			{Name: "Claire", Email: "claire@example.com"},
			{Name: "Andreas", Email: "andreas@example.com"},
			{Name: "Nyomi", Email: "nyomi@example.com"},
			{Name: "Hazel", Email: "hazel@example.com"},
			{Name: "Baker", Email: "baker@example.com"},
			{Name: "Callum", Email: "callum@example.com"},
			{Name: "Ivy", Email: "ivy@example.com"},
			{Name: "Cassie", Email: "cassie@example.com"},
			{Name: "Tiffany", Email: "tiffany@example.com"},
			{Name: "Sasha", Email: "sasha@example.com"},
			{Name: "Sean", Email: "sean@example.com"},
			{Name: "Jefferson", Email: "jefferson@example.com"},
		}
		if err := DB.Create(&users).Error; err != nil {
			log.Fatal("❌ Failed to seed users: ", err)
		} else {
			fmt.Println("✅ Seeded initial users!")
		}
	}

	var dbUsers []models.User
	DB.Find(&dbUsers)

	emailToUser := make(map[string]models.User)
	for _, user := range dbUsers {
		emailToUser[user.Email] = user
	}

	connectPairs := [][2]string{
		{"alice@example.com", "bob@example.com"},
		{"alice@example.com", "charlie@example.com"},
		{"bob@example.com", "lily@example.com"},
		{"charlie@example.com", "claire@example.com"},
		{"claire@example.com", "hazel@example.com"},
		{"baker@example.com", "callum@example.com"},
		{"cassie@example.com", "ivy@example.com"},
	}

	for _, pair := range connectPairs {
		u1 := emailToUser[pair[0]]
		u2 := emailToUser[pair[1]]

		if u1.ID != 0 && u2.ID != 0 {
			err := DB.Model(&u1).Association("Connections").Append(&u2)
			if err != nil {
				log.Printf("❌ Failed to connect %s and %s: %v", u1.Email, u2.Email, err)
			}
			err = DB.Model(&u2).Association("Connections").Append(&u1)
			if err != nil {
				log.Printf("❌ Failed to connect %s and %s: %v", u2.Email, u1.Email, err)
			}
		}
	}

}
