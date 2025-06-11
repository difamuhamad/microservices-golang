package cmd

import (
	"fmt"
	"net/http"
	"time"
	"user-service/common/response"
	"user-service/config"
	"user-service/constants"
	"user-service/controllers"
	"user-service/database/seeders"
	"user-service/domain/models"
	"user-service/middlewares"
	"user-service/repositories"
	"user-service/routes"
	"user-service/services"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var command = *&cobra.Command{
	Use:   "serve",
	Short: "Start the server",
	Run: func(c *cobra.Command, args []string) {

		//	Get config from .env file
		_ = godotenv.Load()

		//	Config setup
		config.Init()

		//	Connect to DB
		db, err := config.InitDatabase()
		if err != nil {
			panic(err)
		}

		//	Set TimeZone
		loc, err := time.LoadLocation("Asia/Jakarta")
		if err != nil {
			panic(err)
		}
		time.Local = loc

		//	GORM will automaticaly create new table if empty
		err = db.AutoMigrate(
			&models.Role{},
			&models.User{},
		)
		if err != nil {
			panic(err)
		}

		//	Fill db with seed (default data)
		seeders.NewSeederRegistry(db).Run()

		// Dependency Injection (DI)
		repository := repositories.NewRepositoryRegistry(db) //inject db to repo
		service := services.NewServiceRegistry(repository)
		controller := controllers.NewControllerRegistry(service)

		// Setup gin router
		router := gin.Default()
		router.Use(middlewares.HandlePanic())
		router.NoRoute(func(c *gin.Context) {
			c.JSON(http.StatusNotFound, response.Response{
				Status:  constants.Error,
				Message: fmt.Sprintf("Path %s", http.StatusText(http.StatusNotFound)),
			})
		})
		router.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, response.Response{
				Status:  constants.Success,
				Message: "Welcome to User Service",
			})
		})
		router.Use(func(c *gin.Context) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // CORS
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			c.Next()
		})

		// Limit the request *per second
		lmt := tollbooth.NewLimiter(
			config.Config.RateLimiterMaxRequest,
			&limiter.ExpirableOptions{
				DefaultExpirationTTL: time.Duration(config.Config.RateLimiterTimeSecond) * time.Second,
			})
		router.Use(middlewares.RateLimiter(lmt))

		// Register all endpoint routes
		group := router.Group("/api/v1")
		route := routes.NewRouteRegistry(controller, group)
		route.Serve()

		port := fmt.Sprintf(":%d", config.Config.Port)
		router.Run(port)

	},
}

func Run() {
	err := command.Execute()
	if err != nil {
		panic(err) // panic is a robust strategies to prevent server crashes when errors occur

	}
}
