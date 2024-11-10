package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/LombardiDaniel/generic-data-collector-api/controllers"
	"github.com/LombardiDaniel/generic-data-collector-api/docs"
	"github.com/LombardiDaniel/generic-data-collector-api/middlewares"
	"github.com/LombardiDaniel/generic-data-collector-api/services"
	"github.com/LombardiDaniel/generic-data-collector-api/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	router *gin.Engine

	formsCol *mongo.Collection

	// Services
	authService  services.AuthService
	formsService services.FormsService

	// Gin Controllers
	formsController controllers.FormsController

	authMiddleware middlewares.AuthMiddleware

	mongoClient *mongo.Client
	ctx         context.Context
)

func init() {
	ctx = context.TODO()

	utils.InitSlogger()

	mongoConn := options.Client().ApplyURI(
		utils.GetEnvVarDefault("MONGO_URI", "mongodb://localhost:27017"),
	)

	mongoClient, err := mongo.Connect(ctx, mongoConn)
	if err != nil {
		slog.Error(err.Error())
	}

	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		slog.Error(err.Error())
	}

	formsDb := mongoClient.Database("formsDb")
	formsCol = formsDb.Collection("forms")

	_, err = formsCol.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "id", Value: -1}},
	})
	if err != nil {
		slog.Error(fmt.Sprintf("Could not create idx: %s", err.Error()))
		return
	}

	authTokens := strings.Split(os.Getenv("AUTH_TOKENS"), ",")

	// Services
	formsService = services.NewFormsServiceMongoImpl(formsCol)
	authService = services.NewAuthServiceImpl(authTokens)

	// Middleware
	authMiddleware = middlewares.NewAuthMiddleware(authService)

	// Controllers
	formsController = controllers.NewFormsController(formsService)

	router = gin.Default()
	router.SetTrustedProxies([]string{"*"})

	docs.SwaggerInfo.Title = "Generic Forms API"
	docs.SwaggerInfo.Description = "Generic Forms API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = ""

	if os.Getenv("GIN_MODE") == "release" {
		docs.SwaggerInfo.Host = os.Getenv("HOST")
		docs.SwaggerInfo.Schemes = []string{"https"}
	} else {
		docs.SwaggerInfo.Host = "localhost:8080"
		docs.SwaggerInfo.Schemes = []string{"http"}
	}

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// @securityDefinitions.apiKey Bearer
// @in header
// @name Authorization
// @description "Type 'Bearer $TOKEN' to correctly set the API Key"
func main() {
	defer mongoClient.Disconnect(ctx)

	// corsCfg := cors.Config{
	// 	AllowAllOrigins:  true,
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
	// 	AllowHeaders:     []string{"Authorization"},
	// 	AllowCredentials: true,
	// 	MaxAge:           12 * time.Hour,
	// }

	corsCfg := cors.DefaultConfig()
	corsCfg.AllowAllOrigins = true
	corsCfg.AddAllowHeaders("Authorization")
	corsCfg.AllowCredentials = false
	corsCfg.MaxAge = 12 * time.Hour

	slog.Info(fmt.Sprintf("corsCfg: %+v\n", corsCfg))

	router.Use(cors.New(corsCfg))

	router.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "OK")
	})

	basePath := router.Group("/v1")
	formsController.RegisterRoutes(basePath, authMiddleware)

	slog.Error(router.Run(":8080").Error())
}
