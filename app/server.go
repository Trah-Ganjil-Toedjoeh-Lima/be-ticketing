package app

import (
	"fmt"
	"github.com/frchandra/gmcgo/config"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Server struct {
	appConfig *config.AppConfig
	db        *gorm.DB
	router    *gin.Engine
}

func NewServer(appConfig *config.AppConfig) *Server {
	db, _ := initializeDb(appConfig)
	router := initializeRouter(appConfig)
	return &Server{
		appConfig: appConfig,
		db:        db,
		router:    router,
	}
}

func initializeDb(appConfig *config.AppConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", appConfig.DBHost, appConfig.DBUser, appConfig.DBPassword, appConfig.DBName, appConfig.DBPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("Failed on connecting to the migrator server")
	} else {
		fmt.Println("db connection established")
		fmt.Println("Using migrator " + db.Migrator().CurrentDatabase())
	}
	return db, err
}

func initializeRouter(appConfig *config.AppConfig) *gin.Engine {
	fmt.Println("Welcome to " + appConfig.AppName)
	if appConfig.IsProduction == "false" {
		gin.SetMode(gin.DebugMode)
	}
	router := gin.Default()
	initializeRoutes(router)
	return router
}

func (this *Server) Run() {
	fmt.Printf("Listening to port %s", this.appConfig.AppPort)
	err := this.router.Run(":" + this.appConfig.AppPort)
	if err != nil {
		panic("Server unable to start")
	}
}
