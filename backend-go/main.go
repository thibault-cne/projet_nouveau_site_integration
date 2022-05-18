package main

import (
	"backend-go/controllers"
	"backend-go/models"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type todo struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

var todoList = []todo{
	{ID: "1", Text: "Todo 1"},
	{ID: "2", Text: "Todo 2"},
	{ID: "3", Text: "Todo 3"},
}

func getTodoList(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, todoList)
}

func getTodoById(c *gin.Context) {
	id := c.Param("id")
	for _, item := range todoList {
		if item.ID == id {
			c.IndentedJSON(http.StatusOK, item)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
}

func load_env() (string, string) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	domain := os.Getenv("APP_DOMAIN")
	port := os.Getenv("APP_PORT")

	return domain, port
}

func init_database() {
	fmt.Println("Initializing database")

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Stars{})
	db.AutoMigrate(&models.DailyGame{})
	db.AutoMigrate(&models.Notifications{})
	db.AutoMigrate(&models.Calendar{})
	db.AutoMigrate(&models.Challenge{})
	db.AutoMigrate(&models.Suggestion{})
	db.AutoMigrate(&models.Tnder{})
}

func main() {
	domain, port := load_env()
	init_database()

	fmt.Println("Starting server on domain " + domain)
	fmt.Println("Starting server on port " + port)

	router := gin.Default()
	// Config CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:8080"
		},
		MaxAge: 12 * time.Hour,
	}))
	basepath := router.Group("/api/v1")
	controllers.Register_login_routes(basepath)

	router.GET("/todos", getTodoList)
	router.GET("/todos/:id", getTodoById)
	router.Run(domain + ":" + port)

}
