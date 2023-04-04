package core

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/crazi-coder/report-service/controller"
	"github.com/crazi-coder/report-service/core/middleware"
	"github.com/crazi-coder/report-service/core/utils/helpers"
	"github.com/crazi-coder/report-service/core/utils/libs"
	"github.com/crazi-coder/report-service/views"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Server is the serve Config Object
type Server interface {
	Start(context.Context) error
}

type server struct {
	host   string
	port   string
	route  *gin.Engine
	logger *logrus.Logger
}

// NewServer creates a new BiddanoAPIServer instance
func NewServer(ctx context.Context, host string, port string) Server {
	// if ldflags.Environment == "production" {
	// 	gin.SetMode(gin.ReleaseMode)
	// }
	s := server{host: host, port: port, route: gin.Default()}
	logger := logrus.StandardLogger()
	// // To initialize Sentry's handler, you need to initialize Sentry itself beforehand
	// if ldflags.Environment == "production" {

	// 	if err := sentry.Init(sentry.ClientOptions{
	// 		Dsn: "https://5510c809df3144a28ef58037bff26bbb@o1074929.ingest.sentry.io/6747083",
	// 		// Set TracesSampleRate to 1.0 to capture 100%
	// 		// of transactions for performance monitoring.
	// 		// We recommend adjusting this value in production,
	// 		TracesSampleRate: 1.0,
	// 		Release:          ldflags.Version + "-" + ldflags.BuildDate,
	// 		Environment:      ldflags.Environment,
	// 	}); err != nil {
	// 		fmt.Printf("Sentry initialization failed: %v\n", err)
	// 	}

	logger.SetLevel(logrus.DebugLevel)

	// 	//Recovery middleware recovers from any panics and writes a 500 if there was one.
	// 	s.route.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
	// 		e := helpers.NewResponse()
	// 		c.AbortWithStatusJSON(http.StatusInternalServerError,
	// 			e.Error(http.StatusInternalServerError, "Internal Server Error", nil))

	// 	}))

	// } else {
	s.route.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		e := helpers.NewResponse()
		if er, ok := recovered.(string); ok {

			c.AbortWithStatusJSON(http.StatusInternalServerError,
				e.Error(http.StatusInternalServerError, "Server error", errors.New(er)))

		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError,
				e.Error(http.StatusInternalServerError, "Server error", nil))
		}
	}))
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.TextFormatter{})
	// }

	s.logger = logger
	s.route.Use(middleware.CORSMiddleware())
	return &s
}

func (s server) Start(ctx context.Context) error {
	// Postgres Information

	addrs := fmt.Sprintf("%s:%s", s.host, s.port)

	psqlPort, _ := strconv.Atoi(helpers.GetEnv("POSTGRES_PORT", "5432"))
	psqlConf := libs.PgConfig{
		Host:     helpers.GetEnv("POSTGRES_HOST", "34.124.157.33"),
		Port:     psqlPort,
		Database: helpers.GetEnv("POSTGRES_DB", "infiviz_dev_db"),
		User:     helpers.GetEnv("POSTGRES_USER", "postgres"),
		Password: helpers.GetEnv("POSTGRES_PASSWORD", "postgres"),
	}
	psql, err := libs.NewPostgreSQLConnection(ctx, s.logger, 1, 10, &psqlConf)
	if err != nil {
		s.logger.WithError(err).Error("Failed to create postgres connection")
		return err
	}
	defer psql.Close() // Close connection before stopping the server.

	// After the connection has been established, enable the jwtAuthMiddleware
	s.route.Use(middleware.AuthMiddleware(psql, s.logger))

	v1 := s.route.Group("/v1/report")
	authCtl := controller.NewReportController(ctx, s.logger, psql)
	v := views.NewReportView(authCtl, v1, s.logger)
	v.Register(ctx)

	srv := &http.Server{
		Addr:           addrs,
		Handler:        s.route,
		ReadTimeout:    1 * time.Second,
		WriteTimeout:   1 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err = srv.ListenAndServe()
	if err != nil {
		s.logger.Error(err)
	}
	return nil
}
