package controller

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"

	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
)

type ReportController interface {
	Run(ctx context.Context, schema string, userID int64, request Request) (*PhotoSession, error)
	Download(ctx context.Context, schema string, userID int64, request Request) ([]*Download, error)
	StoreChannel(ctx context.Context, schema string, userID int64, request Request) ([]*StoreChannel, error)
	StoreBrand(ctx context.Context, schema string, userID int64, request Request) ([]*StoreBrand, error)
	Store(ctx context.Context, schema string, userID int64, request Request) ([]*Store, error)
	Category(ctx context.Context, schema string, userID int64, request Request) ([]*Category, error)
	Users(ctx context.Context, schema string, userID int64, request Request) ([]*User, error)
	PhotoTypes(ctx context.Context, schema string, userID int64, request Request) ([]*PhotoType, error)
	PhotoSessions(ctx context.Context, schema string, userID int64, request Request) ([]*PhotoSession, error)
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

func (r *reportController) Run(ctx context.Context, schema string, userID int64, request Request) (*PhotoSession, error) {

	return nil, nil
}

func (r *reportController) Download(ctx context.Context, schema string, userID int64, request Request) ([]*Download, error) {

	tblDownloadReport := goqu.S(schema).Table("download_report")
	tblReportModelMap := goqu.S(schema).Table("report_model_map")
	tblReportType := goqu.S(schema).Table("report_type")
	nq := r.dialect.From(tblDownloadReport).Select(
		"download_report.id", "download_report.status", "report_type.name",
		"download_report.created", "download_report.modified",
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

func (r *reportController) Store(ctx context.Context, schema string, userID int64, request Request) ([]*Store, error) {
	tblStore := goqu.S(schema).Table("store_store")
	nq := r.dialect.From(tblStore).Select("id", "title").Where(
		goqu.Ex{
			"is_active": true,
		},
	)
	if len(request.StoreBrand) > 0 {
		nq = nq.Where(
			goqu.Ex{"store_brand_id": request.StoreBrand},
		)
	}
	if len(request.StoreChannel) > 0 {
		fmt.Println("........", request.StoreChannel)
		nq = nq.Where(
			goqu.Ex{"store_type_id": request.StoreChannel},
		)
	}
	nq = nq.Prepared(true)
	q, args, err := nq.ToSQL()
	if err != nil {
		return nil, err
	}
	fmt.Println(q)
	res, err := r.conn.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	storeList := []*Store{}
	for res.Next() {
		store := Store{}
		res.Scan(&store.ID, &store.Name)
		storeList = append(storeList, &store)

	}
	fmt.Println(storeList)
	return storeList, nil
}

func (r *reportController) StoreBrand(ctx context.Context, schema string, userID int64, request Request) ([]*StoreBrand, error) {
	tblStore := goqu.S(schema).Table("store_storebrand")
	nq := r.dialect.From(tblStore).Select("id", "title").Where(goqu.Ex{"is_active": true}).Limit(30)
	q, args, err := nq.ToSQL()
	if err != nil {
		return nil, err
	}
	res, err := r.conn.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	storeBrandList := []*StoreBrand{}
	for res.Next() {
		store := StoreBrand{}
		res.Scan(&store.ID, &store.Name)
		storeBrandList = append(storeBrandList, &store)

	}
	return storeBrandList, nil
}

func (r *reportController) StoreChannel(ctx context.Context, schema string, userID int64, request Request) ([]*StoreChannel, error) {
	tblStore := goqu.S(schema).Table("store_storetype")
	nq := r.dialect.From(tblStore).Select("id", "title").Where(goqu.Ex{"is_active": true}).Limit(30)
	q, args, err := nq.ToSQL()
	if err != nil {
		return nil, err
	}
	res, err := r.conn.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	storeChannelList := []*StoreChannel{}
	for res.Next() {
		store := StoreChannel{}
		res.Scan(&store.ID, &store.Name)
		storeChannelList = append(storeChannelList, &store)

	}
	return storeChannelList, nil
}

func (r *reportController) Category(ctx context.Context, schema string, userID int64, request Request) ([]*Category, error) {
	tblCategory := goqu.S(schema).Table("common_category")
	nq := r.dialect.From(tblCategory).Select("id", "title").Where(goqu.Ex{"is_active": true})
	q, args, err := nq.ToSQL()
	if err != nil {
		return nil, err
	}
	res, err := r.conn.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	categoryList := []*Category{}
	for res.Next() {
		c := Category{}
		res.Scan(&c.ID, &c.Name)
		categoryList = append(categoryList, &c)
	}
	return categoryList, nil
}

func (r *reportController) Users(ctx context.Context, schema string, userID int64, request Request) ([]*User, error) {
	tblUser := goqu.S(schema).Table("auth_user")
	nq := r.dialect.From(tblUser).Select("id", "username").Where(goqu.Ex{"is_active": true}).Limit(30)
	q, args, err := nq.ToSQL()
	if err != nil {
		return nil, err
	}
	res, err := r.conn.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	userList := []*User{}
	for res.Next() {
		c := User{}
		res.Scan(&c.ID, &c.Name)
		userList = append(userList, &c)
	}
	return userList, nil
}

func (r *reportController) PhotoSessions(ctx context.Context, schema string, userID int64, request Request) ([]*PhotoSession, error) {

	tblPhotoSession := goqu.S(schema).Table("photo_photosession")
	tblStore := goqu.S(schema).Table("store_store")
	tblUser := goqu.S(schema).Table("auth_user")
	tblCategory := goqu.S(schema).Table("common_category")

	nq := r.dialect.From(tblPhotoSession).Select(
		"photo_photosession.session_id", "photo_photosession.photo_count", "store_store.id", "store_store.title",
		"auth_user.id", "auth_user.username", "common_category.id", "common_category.name",
		"photo_photosession.created_on", "photo_photosession.visit_timestamp",
	).LeftJoin(
		tblStore, goqu.On(goqu.Ex{
			"store_store.id": goqu.I("photo_photosession.store_id"),
		}),
	).LeftJoin(
		tblUser, goqu.On(goqu.Ex{
			"auth_user.id": goqu.I("photo_photosession.user_id"),
		}),
	).RightJoin( // TODO: Change to Left Join for multi-category use case
		tblCategory, goqu.On(goqu.Ex{
			"common_category.id": goqu.I("photo_photosession.user_id"),
		}),
	).Limit(100).Prepared(true)

	q, args, err := nq.ToSQL()
	if err != nil {
		return nil, err
	}
	res, err := r.conn.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	results := []*PhotoSession{}
	for res.Next() {
		p := PhotoSession{}
		s := Store{}
		u := User{}
		c := Category{}

		var (
			created time.Time
			visited sql.NullTime
		)

		err := res.Scan(&p.ID, &p.PhotoCount, &s.ID, &s.Name, &u.ID, &u.Name, &c.ID, &c.Name, &created, &visited)
		if err != nil {
			panic(err)
		}

		p.CreatedAt = created.Format(time.RFC822)

		if visited.Valid {
			p.VisitedOn = visited.Time.Format(time.RFC822)
		}

		p.Store = s
		p.PhotoTakenBy = u
		p.Category = c
		results = append(results, &p)
	}
	return results, nil
}

func (r *reportController) PhotoTypes(ctx context.Context, schema string, userID int64, request Request) ([]*PhotoType, error) {
	tblPhotoType := goqu.S(schema).Table("common_phototype")
	nq := r.dialect.From(tblPhotoType).Select("id", "title")
	q, args, err := nq.ToSQL()
	if err != nil {
		return nil, err
	}
	res, err := r.conn.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	photoTypeList := []*PhotoType{}
	for res.Next() {
		c := PhotoType{}
		res.Scan(&c.ID, &c.Name)
		photoTypeList = append(photoTypeList, &c)
	}
	return photoTypeList, nil
}