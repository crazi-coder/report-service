package views

import (
	"context"
	"net/http"
	"strconv"

	"github.com/crazi-coder/report-service/controller"
	"github.com/crazi-coder/report-service/core/utils"
	"github.com/crazi-coder/report-service/core/utils/helpers"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ReportView interface {
	Register(ctx context.Context) error // register filter urls
	Run(ctx *gin.Context)
	Download(ctx *gin.Context)
}

type reportView struct {
	controller controller.ReportController
	routeGroup *gin.RouterGroup
	logger     *logrus.Logger
}

func NewReportView(controller controller.ReportController,
	routeGroup *gin.RouterGroup, logger *logrus.Logger) ReportView {
	return &reportView{controller: controller, routeGroup: routeGroup, logger: logger}
}

// Register registers a API endpoint
func (r *reportView) Register(ctx context.Context) error {
	r.routeGroup.POST("/run", r.Run)
	r.routeGroup.GET("/download", r.Download)
	return nil
}

func (r *reportView) Run(ctx *gin.Context) {

	ctx.AbortWithStatusJSON(http.StatusOK, "")
}

func (r *reportView) Download(ctx *gin.Context) {
	vctx := context.WithValue(ctx.Request.Context(), utils.CtxSchema, ctx.Value(utils.CtxSchema))
	u := ctx.Value(utils.CtxUserID).(string)
	resp := helpers.NewResponse()

	userID, err := strconv.ParseInt(u, 10, 64)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusExpectationFailed, resp.Error(helpers.ErrCodeServerError, "Unknown User", err))
		return
	}

	p, err := r.controller.Download(vctx, userID, controller.Report{})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusExpectationFailed, resp.Error(helpers.ErrCodeServerError, "Process failed", err))
		return
	}
	ctx.AbortWithStatusJSON(http.StatusOK, p)
}
