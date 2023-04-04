package controller

import (
	"context"
	"time"

	"github.com/crazi-coder/report-service/core/utils"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"

	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
)

type ReportController interface {
	Run(ctx context.Context, report Report) (*Report, error)
	Download(ctx context.Context, userID int64, report Report) ([]*Download, error)
}

type reportController struct {
	conn    *pgxpool.Pool
	ctx     context.Context
	logger  *logrus.Logger
	tenant  exp.IdentifierExpression
	dialect goqu.DialectWrapper
}

func NewReportController(ctx context.Context, logger *logrus.Logger, conn *pgxpool.Pool) ReportController {

	return &reportController{ctx: ctx, logger: logger, conn: conn, dialect: goqu.Dialect("postgres")}
}

func (r *reportController) Run(ctx context.Context, report Report) (*Report, error) {

	return nil, nil
}

func (r *reportController) Download(ctx context.Context, userID int64, report Report) ([]*Download, error) {
	schema := ctx.Value(utils.CtxSchema).(string)
	tblDownloadReport := goqu.S(schema).Table("download_report")
	tblReportModelMap := goqu.S(schema).Table("report_model_map")
	tblReportType := goqu.S(schema).Table("report_type")
	nq := r.dialect.From(tblDownloadReport).Select(
		"download_report.id", "download_report.status", "report_type.name", "download_report.created", "download_report.modified",
	).InnerJoin(
		tblReportModelMap, goqu.On(goqu.Ex{
			"report_model_map.id": goqu.I("download_report.report_map_id"),
		}),
	).InnerJoin(
		tblReportType, goqu.On(goqu.Ex{
			"report_model_map.report_type_id": goqu.I("report_type.id"),
		}),
	).Prepared(true)
	sql, args, err := nq.ToSQL()
	if err != nil {
		return nil, err
	}

	r.logger.WithFields(logrus.Fields{"query": sql, "params": args}).Debug("Running ...")
	res, err := r.conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	results := []*Download{}
	for res.Next() {
		var (
			created  time.Time
			modified time.Time
		)

		d := Download{}
		err := res.Scan(&d.ID, &d.Status, &d.ReportName, &created, &modified)
		if err != nil {
			return nil, err
		}
		d.Created = created.Format(time.RFC3339)
		d.Modified = modified.Format(time.RFC3339)
		results = append(results, &d)
	}
	return results, nil
}
