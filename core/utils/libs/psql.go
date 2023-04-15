package libs

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

const (
	// DatabaseConnectionTimeOut is the default timeout for established connections.
	DatabaseConnectionTimeOut = 1 * time.Second
)

// PgConfig set configuration
type PgConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	scheme   string
}

// NewPostgreSQLConnection create new PostgreSQL Connection.
func NewPostgreSQLConnection(ctx context.Context, logger *logrus.Logger, minPoolSize int32,
	maxPoolSize int32, conf *PgConfig) (*pgxpool.Pool, error) {

	dsn := fmt.Sprintf(
		//"postgresql://%s:%s@%s:%d/%s?statement_cache_mode=describe&sslmode=disable",
		//"postgresql://%s:%s@%s:%d/%s",
		"user=%s password=%s host=%s port=%d dbname=%s",
		conf.User, conf.Password, conf.Host, conf.Port, conf.Database,
	)
	fmt.Println(dsn)
	ctx, cancel := context.WithTimeout(ctx, DatabaseConnectionTimeOut)
	defer cancel()
	connConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		
	}
	if maxPoolSize == 0 {
		maxPoolSize = 1 // Default pool size is set to 1
	}
	if minPoolSize == 0 {
		minPoolSize = 1 // Default pool size is set to 1
	}
	if conf.scheme == "" {
		conf.scheme = "public" // Default scheme is public.
	}
	connConfig.MaxConns = maxPoolSize
	connConfig.MinConns = minPoolSize
	connConfig.ConnConfig.PreferSimpleProtocol = true

	//connConfig.ConnConfig.BuildStatementCache = nil
	password := fmt.Sprintf("%s%s", strings.Repeat("*", len(conf.Password)-3), conf.Password[len(conf.Password)-3:])
	logger.WithFields(logrus.Fields{
		"Driver":      "PostgreSQL",
		"User":        conf.User,
		"Password":    password,
		"Host":        conf.Host,
		"Port":        conf.Port,
		"Database":    conf.Database,
		"maxPoolSize": connConfig.MaxConns,
		"minPoolSize": connConfig.MinConns,
	}).Info("Database Connection Information")

	conn, err := pgxpool.ConnectConfig(ctx, connConfig)
	if err != nil {
		logger.WithError(err).Errorf("Unable to connection to database: %v", err)
		return nil, err
	}
	err = conn.Ping(ctx)
	if err != nil {
		logger.WithError(err).Errorf("Unable to connection to database: %v", err)
		return nil, err
	}
	return conn, nil
}
