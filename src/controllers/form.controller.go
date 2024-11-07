package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"log/slog"

	"github.com/LombardiDaniel/generic-data-collector-api/middlewares"
	"github.com/LombardiDaniel/generic-data-collector-api/models"
	"github.com/LombardiDaniel/generic-data-collector-api/services"
	"github.com/gin-gonic/gin"
)

type FormsController struct {
	dataStoreService services.FormsService
}

func NewFormsController(
	dataStoreService services.FormsService,
) FormsController {
	return FormsController{
		dataStoreService: dataStoreService,
	}
}

// @Summary AddEntry
// @Tags Form
// @Description Adds an entry to the db
// @Consume application/json
// @Accept json
// @Produce plain
// @Param   formPayload body 		models.Form true "user json"
// @Success 200 		{string} 	OKResponse "OK"
// @Failure 400 		{string} 	ErrorResponse "Bad Request"
// @Failure 502 		{string} 	ErrorResponse "Bad Gateway"
// @Router /v1/entries [PUT]
func (c *FormsController) AddEntry(ctx *gin.Context) {
	rCtx := ctx.Request.Context()
	var form models.Form

	if err := ctx.ShouldBind(&form); err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	fmt.Printf("form: %v\n", form)

	err := c.dataStoreService.InsertPayload(rCtx, form)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while inserting form '%s': '%s'", form.Email, err.Error()))
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	ctx.String(http.StatusOK, "OK")
}

// @Summary GetEntries
// @Security Bearer
// @Tags Form
// @Description Get entries from the db
// @Produce json
// @Param 	n 			query string false "Returns N parameters, if not passed, reutrns all"
// @Param   id path string true "identifier"
// @Success 200 {object} 	[]models.Form
// @Failure 400 {string} 	ErrorResponse "Bad Request"
// @Failure 409 {string} 	ErrorResponse "Conflict"
// @Failure 502 {string} 	ErrorResponse "Bad Gateway"
// @Router /v1/entries/{id} [GET]
func (c *FormsController) GetEntries(ctx *gin.Context) {
	rCtx := ctx.Request.Context()
	id := ctx.Param("id")
	n := ctx.Query("n")
	nInt, err := strconv.ParseInt(n, 10, 32)
	if err != nil {
		payloads, err := c.dataStoreService.Get(rCtx, id)
		if err != nil {
			ctx.String(http.StatusBadGateway, "BadGateway")
			return
		}
		ctx.JSON(http.StatusOK, payloads)
		return
	}

	nInt32 := int32(nInt)
	if nInt32 > 0 {
		payloads, err := c.dataStoreService.GetN(rCtx, id, uint32(nInt32))
		if err != nil {
			ctx.String(http.StatusBadGateway, "BadGateway")
			return
		}
		ctx.JSON(http.StatusOK, payloads)
		return
	}
}

func (c *FormsController) RegisterRoutes(rg *gin.RouterGroup, authMiddleware middlewares.AuthMiddleware) {
	r := rg.Group("/entries")

	r.PUT("/", c.AddEntry)
	r.GET("/:id", authMiddleware.Authorize(), c.GetEntries)
}
