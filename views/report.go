package views

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/crazi-coder/report-service/controller"
	"github.com/crazi-coder/report-service/core/utils"
	"github.com/crazi-coder/report-service/core/utils/helpers"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
)

type requestContext struct {
	requestUserID int64
	requestSchema string
}

type ReportView interface {
	Register(ctx context.Context) error // register filter urls
	PhotoType(ctx *gin.Context)
	Download(ctx *gin.Context)
	Store(ctx *gin.Context)
	Category(ctx *gin.Context)
	Users(ctx *gin.Context)
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
	r.routeGroup.GET("/photo-types", r.PhotoType)
	r.routeGroup.GET("/stores", r.Store)
	r.routeGroup.GET("/stores/channel", r.StoreBrand)
	r.routeGroup.GET("/stores/brand", r.StoreBrand)
	r.routeGroup.GET("/categories", r.Category)
	r.routeGroup.GET("/users", r.Users)
	r.routeGroup.GET("/photos/sessions", r.PhotoSession)
	r.routeGroup.GET("/runs/downloads", r.Download)
	return nil
}

func (r *reportView) validate(ctx *gin.Context) (requestContext, error) {
	rCtx := requestContext{}
	rCtx.requestSchema = ctx.Value(utils.CtxSchema).(string)
	u := ctx.Value(utils.CtxUserID).(string)

	requestUserID, err := strconv.ParseInt(u, 10, 64)
	rCtx.requestUserID = requestUserID
	return rCtx, err
}

func (r *reportView) PhotoType(ctx *gin.Context) {
	resp := helpers.NewResponse()
	rCtx, err := r.validate(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusExpectationFailed, resp.Error(helpers.ErrCodeServerError, "Unknown User", err))
		return
	}

	p, err := r.controller.PhotoTypes(ctx.Request.Context(), rCtx.requestSchema, rCtx.requestUserID, controller.Request{})

	ctx.AbortWithStatusJSON(http.StatusOK, p)
}

func (r *reportView) Download(ctx *gin.Context) {
	resp := helpers.NewResponse()
	rCtx, err := r.validate(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusExpectationFailed, resp.Error(helpers.ErrCodeServerError, "Unknown User", err))
		return
	}

	p, err := r.controller.Download(ctx.Request.Context(), rCtx.requestSchema, rCtx.requestUserID, controller.Request{})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusExpectationFailed, resp.Error(helpers.ErrCodeServerError, "Process failed", err))
		return
	}
	ctx.AbortWithStatusJSON(http.StatusOK, p)
}

func (r *reportView) Store(ctx *gin.Context) {
	resp := helpers.NewResponse()
	rCtx, err := r.validate(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusExpectationFailed, resp.Error(helpers.ErrCodeServerError, "Unknown User", err))
		return
	}

	storeBrandStr := ctx.Query("store_brand_list")
	storeBrandStrList := strings.Split(storeBrandStr, ",")
	if len(storeBrandStrList) < 1 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp.Error(helpers.ErrCodeServerError, "Expected to pass store chanel", err))
		return
	}
	var storeBrandList []int
	req := controller.Request{}
	for _, e := range storeBrandStrList {
		i, err := strconv.ParseInt(e, 10, 64)
		if err == nil {
			storeBrandList = append(storeBrandList, int(i))
		}
	}

	storeChannelStr := ctx.Query("store_channel_list")
	storeChannelStrList := strings.Split(storeChannelStr, ",")

	if len(storeChannelStrList) == 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp.Error(helpers.ErrCodeServerError, "Expected to pass store chanel", err))
		return
	}
	var storeChannelList []int
	for _, element := range storeChannelStrList {
		i, err := strconv.ParseInt(element, 10, 64)
		if err == nil {
			storeChannelList = append(storeChannelList, int(i))
		}
	}

	req.StoreBrand = storeBrandList
	req.StoreChannel = storeChannelList
	if len(storeBrandList) > 0 || len(storeChannelList) > 0 {
		s, _ := r.controller.Store(ctx.Request.Context(), rCtx.requestSchema, rCtx.requestUserID, req)
		ctx.AbortWithStatusJSON(http.StatusOK, s)
		return
	}
	ctx.AbortWithStatusJSON(http.StatusNoContent, "")
}

func (r *reportView) StoreBrand(ctx *gin.Context) {
	resp := helpers.NewResponse()

	rCtx, err := r.validate(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusExpectationFailed, resp.Error(helpers.ErrCodeServerError, "Unknown User", err))
		return
	}
	storeChannelStr := ctx.Query("store_channel_list")
	storeChannelStrList := strings.Split(storeChannelStr, ",")
	if len(storeChannelStrList) == 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, resp.Error(helpers.ErrCodeServerError, "Expected to pass store chanel", err))
		return
	}
	req := controller.Request{}
	for _, e := range storeChannelStrList {
		i, err := strconv.ParseInt(e, 10, 64)
		if err == nil {
			req.StoreChannel = append(req.StoreChannel, int(i))
		}
	}
	s, _ := r.controller.StoreBrand(ctx.Request.Context(), rCtx.requestSchema, rCtx.requestUserID, req)
	ctx.AbortWithStatusJSON(http.StatusOK, s)
}

func (r *reportView) Category(ctx *gin.Context) {
	resp := helpers.NewResponse()
	rCtx, err := r.validate(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusExpectationFailed, resp.Error(helpers.ErrCodeServerError, "Unknown User", err))
		return
	}

	p, err := r.controller.Category(ctx.Request.Context(), rCtx.requestSchema, rCtx.requestUserID, controller.Request{})

	ctx.AbortWithStatusJSON(http.StatusOK, p)
}

func (r *reportView) Users(ctx *gin.Context) {

	resp := helpers.NewResponse()
	rCtx, err := r.validate(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusExpectationFailed, resp.Error(helpers.ErrCodeServerError, "Unknown User", err))
		return
	}

	p, err := r.controller.Users(ctx.Request.Context(), rCtx.requestSchema, rCtx.requestUserID, controller.Request{})

	ctx.AbortWithStatusJSON(http.StatusOK, p)
}

func (r *reportView) PhotoSession(ctx *gin.Context) {
	resp := helpers.NewResponse()
	rCtx, err := r.validate(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusExpectationFailed, resp.Error(helpers.ErrCodeServerError, "Unknown User", err))
		return
	}

	storeStr := ctx.Query("store_list")
	storeBrandStr := ctx.Query("store_brand_list")
	storeChannelStr := ctx.Query("store_channel_list")
	categoryStr := ctx.Query("category_list")
	photoTypeStr := ctx.Query("photo_type_list")
	visitedFrom := ctx.Query("visited_from")
	visitedTo := ctx.Query("visited_to")
	pageSize := ctx.DefaultQuery("page_size", "100")
	pageNumber := ctx.DefaultQuery("page", "1")

	req := controller.Request{}

	req.StoreList(storeStr)
	req.StoreBrandList(storeBrandStr)
	req.StoreChannelList(storeChannelStr)
	req.CategoryList(categoryStr)
	req.PhotoTypeList(photoTypeStr)
	req.SetPageNumber(pageNumber)
	req.SetPageSize(pageSize)

	// Convert the string representation of timestamp to a date object
	if visitedFrom != "" {
		from, err := strconv.ParseInt(visitedFrom, 10, 64)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusExpectationFailed, resp.Error(helpers.ErrCodeStatusBadRequest, "Wring from date", err))
			return
		}
		req.VisitedFrom = time.Unix(from, 0)
	}
	if visitedTo != "" {
		to, err := strconv.ParseInt(visitedTo, 10, 64)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusExpectationFailed, resp.Error(helpers.ErrCodeStatusBadRequest, "Wring to date", err))
			return
		}
		req.VisitedTo = time.Unix(to, 0)
	}
	p, err := r.controller.PhotoSessions(ctx.Request.Context(), rCtx.requestSchema, rCtx.requestUserID, ctx.Request.RequestURI, req)
	switch err {
	case nil:
		ctx.AbortWithStatusJSON(http.StatusOK, p)
	case pgx.ErrNoRows:
		ctx.AbortWithStatusJSON(http.StatusNoContent,
			resp.Error(helpers.ErrCodeDataNotFound, controller.DataNotFound, err),
		)
	case helpers.ErrPageLimitExceeded:
		ctx.AbortWithStatusJSON(http.StatusBadRequest,
			resp.Error(helpers.ErrPageLimitExceededError, controller.InvalidPageNumber, err),
		)
	default:
		ctx.AbortWithStatusJSON(http.StatusExpectationFailed,
			resp.Error(helpers.ErrCodeServerError, controller.Unrecognized, err),
		)
		r.logger.WithError(err).Error("Error retrieving Photo Session details")
	}
	return
}
