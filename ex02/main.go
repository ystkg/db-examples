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

	_ "github.com/jackc/pgx/v5/stdlib"
	"gopkg.in/yaml.v3"
)

var (
	//go:embed docker-compose.yml
	yml []byte

	//go:embed table/pg.ddl
	pgddl string

	//go:embed table/pgclean.ddl
	pgclean string
)

func setup(ctx context.Context) (*sql.DB, error) {
	return setupPg(ctx)
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

func stats(db *sql.DB, msg string) {
	stats := db.Stats()
	slog.Info(fmt.Sprintf("%-6s", msg), "Open", stats.OpenConnections, "InUse", stats.InUse, "Idle", stats.Idle)
}

func statsIdName(db *sql.DB, msg string, id int32, name string) {
	stats := db.Stats()
	slog.Info(fmt.Sprintf("%-6s", msg), "Open", stats.OpenConnections, "InUse", stats.InUse, "Idle", stats.Idle, "id", id, "name", name)
}

func statsErr(db *sql.DB, msg string, err error) {
	stats := db.Stats()
	slog.Info(fmt.Sprintf("%-6s", msg), "Open", stats.OpenConnections, "InUse", stats.InUse, "Idle", stats.Idle, "err", err)
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("no name")
	}
	exname := os.Args[1]

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := setup(ctx)
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
	case strings.EqualFold(exname, "Ex0201"):
		err = Ex0201(ctx, db)
	case strings.EqualFold(exname, "Ex0202"):
		err = Ex0202(ctx, db)
	case strings.EqualFold(exname, "Ex0203"):
		err = Ex0203(ctx, db)
	case strings.EqualFold(exname, "Ex0204"):
		err = Ex0204(ctx, db)
	case strings.EqualFold(exname, "Ex0205"):
		err = Ex0205(ctx, db)
	case strings.EqualFold(exname, "Ex0206"):
		err = Ex0206(ctx, db)
	case strings.EqualFold(exname, "Ex0207"):
		err = Ex0207(ctx, db)
	case strings.EqualFold(exname, "Ex0208"):
		err = Ex0208(ctx, db)
	case strings.EqualFold(exname, "Ex0209"):
		err = Ex0209(ctx, db)
	case strings.EqualFold(exname, "Ex0210"):
		err = Ex0210(ctx, db)
	case strings.EqualFold(exname, "Ex0211"):
		err = Ex0211(ctx, db)
	default:
		err = fmt.Errorf("unknown:%s", exname)
	}
	if err != nil {
		log.Fatal(err)
	}
}
