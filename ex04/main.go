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
	"gopkg.in/yaml.v3"
)

var (
	//go:embed docker-compose.yml
	yml []byte

	//go:embed table/mysql.ddl
	mysqlddl string

	//go:embed table/mysql.dml
	mysqldml string

	//go:embed table/mysqlclean.ddl
	mysqlclean string

	//go:embed table/pg.ddl
	pgddl string

	//go:embed table/pgclean.ddl
	pgclean string
)

func setup(ctx context.Context) (*sql.DB, *sql.DB, error) {
	pgDB, err := setupPg(ctx)
	if err != nil {
		return nil, nil, err
	}
	myDB, err := setupMySQL(ctx)
	if err != nil {
		pgDB.Close()
		return nil, nil, err
	}
	return pgDB, myDB, nil
}

func setupPg(ctx context.Context) (*sql.DB, error) {
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

	db, err := sql.Open("pgx",
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

	_, err = db.ExecContext(ctx, mysqldml)
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pgDB, myDB, err := setup(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := pgDB.Close(); err != nil {
			slog.WarnContext(ctx, "Close", "err", err)
		}
		if err := myDB.Close(); err != nil {
			slog.WarnContext(ctx, "Close", "err", err)
		}
	}()

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	switch {
	case strings.EqualFold(exname, "Ex04Tx01"):
		err = Ex04Tx01(ctx, pgDB, myDB)
	case strings.EqualFold(exname, "Ex04Tx02"):
		err = Ex04Tx02(ctx, pgDB, myDB)
	case strings.EqualFold(exname, "Ex04Xa01"):
		err = Ex04Xa01(ctx, pgDB, myDB)
	case strings.EqualFold(exname, "Ex04Xa02"):
		err = Ex04Xa02(ctx, pgDB, myDB)
	default:
		err = fmt.Errorf("unknown:%s", exname)
	}
	if err != nil {
		log.Fatal(err)
	}
}
