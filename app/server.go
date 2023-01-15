package app

import (
	"fmt"
	"github.com/frchandra/gmcgo/config"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	Config   *config.AppConfig
	Database *gorm.DB
	Router   *gin.Engine
}

func NewServer(appConfig *config.AppConfig) *Server {
	db, _ := initializeDb(appConfig)
	router := InitializeRouter(appConfig)
	return &Server{
		Config:   appConfig,
		Database: db,
		Router:   router,
	}
}

func initializeDb(appConfig *config.AppConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", appConfig.DBHost, appConfig.DBUser, appConfig.DBPassword, appConfig.DBName, appConfig.DBPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed on connecting to the database server")
	} else {
		fmt.Println("Database connection established")
		fmt.Println("Using database " + db.Migrator().CurrentDatabase())
	}
	return db, err
}

func InitializeRouter(appConfig *config.AppConfig) *gin.Engine {
	fmt.Println("Welcome to " + appConfig.AppName)
	if appConfig.IsProduction == "false" {
		gin.SetMode(gin.DebugMode)
	}
	router := gin.Default()
	initializeRoutes(router)
	return router
}

func (s *Server) Run() {
	fmt.Printf("Listening to port %s", s.Config.AppPort)
	err := s.Router.Run(":" + s.Config.AppPort)
	if err != nil {
		panic("Server unable to start")
	}
}
