package main

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
	"gopkg.in/yaml.v3"
)

var (
	//go:embed docker-compose.yml
	yml []byte

	//go:embed table/mysql.ddl
	mysqlddl string

	//go:embed table/mysqlclean.ddl
	mysqlclean string

	//go:embed table/pg.ddl
	pgddl string

	//go:embed table/pgclean.ddl
	pgclean string
)

func setup(ctx context.Context, exname, pgDriver string) (*sql.DB, error) {
	if len(exname) < 5 {
		return nil, fmt.Errorf("unknown:%s", exname)
	}
	name := strings.ToUpper(exname[4:])
	pgDriverName := "pgx"
	if strings.EqualFold(pgDriver, "pq") {
		pgDriverName = "postgres"
	}
	switch {
	case strings.HasPrefix(name, "PG"):
		return setupPg(ctx, pgDriverName)
	case strings.HasPrefix(name, "MYSQL"):
		return setupMySQL(ctx)
	}
	return nil, fmt.Errorf("unknown:%s", exname)
}

func setupPg(ctx context.Context, driverName string) (*sql.DB, error) {
	conf := struct {
		Services struct {
			Postgres struct {
				Ports       []string
				Environment struct {
					PostgresPassword string `yaml:"POSTGRES_PASSWORD"`
				}
			}
		}
	}{}
	if err := yaml.Unmarshal(yml, &conf); err != nil {
		return nil, err
	}

	db, err := sql.Open(driverName,
		fmt.Sprintf("postgres://postgres:%s@localhost:5432/postgres?sslmode=disable&TimeZone=Asia/Tokyo",
			conf.Services.Postgres.Environment.PostgresPassword,
		),
	)
	if err != nil {
		return nil, err
	}

	_, err = db.ExecContext(ctx, pgclean)
	if err != nil {
		return nil, err
	}

	_, err = db.ExecContext(ctx, pgddl)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func setupMySQL(ctx context.Context) (*sql.DB, error) {
	conf := struct {
		Services struct {
			Mysql struct {
				Environment struct {
					MysqlRootPassword string `yaml:"MYSQL_ROOT_PASSWORD"`
					MysqlDatabase     string `yaml:"MYSQL_DATABASE"`
				}
			}
		}
	}{}
	if err := yaml.Unmarshal(yml, &conf); err != nil {
		return nil, err
	}

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return nil, err
	}

	conn, err := mysql.NewConnector(&mysql.Config{
		Addr:      "localhost:3306",
		DBName:    conf.Services.Mysql.Environment.MysqlDatabase,
		User:      "root",
		Passwd:    conf.Services.Mysql.Environment.MysqlRootPassword,
		ParseTime: true,
		Loc:       loc,
	})
	if err != nil {
		return nil, err
	}
	db := sql.OpenDB(conn)

	_, err = db.ExecContext(ctx, mysqlclean)
	if err != nil {
		return nil, err
	}

	_, err = db.ExecContext(ctx, mysqlddl)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("no name")
	}
	exname := os.Args[1]
	pgDriver := ""
	if 3 <= len(os.Args) {
		pgDriver = os.Args[2]
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := setup(ctx, exname, pgDriver)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			slog.WarnContext(ctx, "Close", "err", err)
		}
	}()

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	switch {
	case strings.EqualFold(exname, "Ex01MySQL01"):
		err = Ex01MySQL01(ctx, db)
	case strings.EqualFold(exname, "Ex01MySQL02"):
		err = Ex01MySQL02(ctx, db)
	case strings.EqualFold(exname, "Ex01MySQL03"):
		err = Ex01MySQL03(ctx, db)
	case strings.EqualFold(exname, "Ex01Pg01"):
		err = Ex01Pg01(ctx, db)
	case strings.EqualFold(exname, "Ex01Pg02"):
		err = Ex01Pg02(ctx, db)
	case strings.EqualFold(exname, "Ex01Pg03"):
		err = Ex01Pg03(ctx, db)
	case strings.EqualFold(exname, "Ex01Pg04"):
		err = Ex01Pg04(ctx, db)
	case strings.EqualFold(exname, "Ex01Pg05"):
		err = Ex01Pg05(ctx, db)
	case strings.EqualFold(exname, "Ex01Pg06"):
		err = Ex01Pg06(ctx, db)
	default:
		err = fmt.Errorf("unknown:%s", exname)
	}
	if err != nil {
		log.Fatal(err)
	}
}
