package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

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

	formsCol  *mongo.Collection
	usersCol  *mongo.Collection
	tokensCol *mongo.Collection

	// Services
	authService  services.AuthService
	formsService services.FormService

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

	formsDb := mongoClient.Database("formsdb")

	formsCol = formsDb.Collection("forms")
	usersCol = formsDb.Collection("users")
	tokensCol = formsDb.Collection("tokens")

	_, err = formsCol.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "id", Value: -1}},
	})
	if err != nil {
		slog.Error(fmt.Sprintf("Could not create idx: %s", err.Error()))
		return
	}

	uniqueIdx := true
	_, err = usersCol.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "username", Value: -1}},
		Options: &options.IndexOptions{
			Unique: &uniqueIdx,
		},
	})
	if err != nil {
		slog.Error(fmt.Sprintf("Could not create idx: %s", err.Error()))
		return
	}

	_, err = tokensCol.Indexes().CreateMany(
		ctx,
		[]mongo.IndexModel{
			{
				Keys: bson.D{{Key: "username", Value: -1}},
				Options: &options.IndexOptions{
					Unique: &uniqueIdx,
				},
			}, {
				Keys: bson.D{{Key: "token", Value: -1}},
				Options: &options.IndexOptions{
					Unique: &uniqueIdx,
				},
			},
		},
	)
	if err != nil {
		slog.Error(fmt.Sprintf("Could not create idx: %s", err.Error()))
		return
	}

	// Services
	formsService = services.NewFormServiceMongoImpl(formsCol)
	authService = services.NewAuthServiceImpl(tokensCol)

	// Middleware
	authMiddleware = middlewares.NewAuthMiddleware(authService)

	// Controllers
	formsController = controllers.NewFormsController(formsService)

	router = gin.Default()
	router.SetTrustedProxies([]string{"*"})

	corsCfg := cors.DefaultConfig()
	corsCfg.AllowAllOrigins = true
	corsCfg.AddAllowHeaders("Authorization")

	slog.Info(fmt.Sprintf("corsCfg: %+v\n", corsCfg))

	router.Use(cors.New(corsCfg))

	docs.SwaggerInfo.Title = "Generic Forms API"
	docs.SwaggerInfo.Description = "Generic Forms API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = ""

	if os.Getenv("GIN_MODE") == "release" {
		docs.SwaggerInfo.Host = os.Getenv("SWAGGER_HOST")
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

	router.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "OK")
	})

	basePath := router.Group("/v1")
	formsController.RegisterRoutes(basePath, authMiddleware)

	slog.Error(router.Run(":8080").Error())
}
