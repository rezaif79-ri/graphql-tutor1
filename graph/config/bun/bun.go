package bun

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var BunDB *bun.DB

func OpenBunDBConn() {
	dbHost := os.Getenv("DBHOST")
	dbPort := os.Getenv("DBPORT")
	dbUser := os.Getenv("DBUSER")
	dbPass := os.Getenv("DBPASS")
	dbName := os.Getenv("DBNAME")

	pgconn := pgdriver.NewConnector(
		pgdriver.WithNetwork("tcp"),
		pgdriver.WithAddr(dbHost+":"+dbPort),
		pgdriver.WithUser(dbUser),
		pgdriver.WithPassword(dbPass),
		pgdriver.WithDatabase(dbName),
		pgdriver.WithTimeout(60*time.Second),
		pgdriver.WithDialTimeout(60*time.Second),
		pgdriver.WithReadTimeout(60*time.Second),
		pgdriver.WithWriteTimeout(60*time.Second),
		pgdriver.WithInsecure(true),
	)
	sqldb := sql.OpenDB(pgconn)
	db := bun.NewDB(sqldb, pgdialect.New())
	fmt.Println("ENV Mode: ", os.Getenv("ENV"))
	if !strings.EqualFold(os.Getenv("ENV"), "PRODUCTION") {
		db.AddQueryHook(initBunLogger())
	}

	BunDB = db
}

func CloseBunDBConn() error {
	return BunDB.Close()
}

type bunLogger struct {
	writer io.Writer
}

func initBunLogger() *bunLogger {
	return &bunLogger{
		writer: os.Stderr,
	}
}

func (h bunLogger) BeforeQuery(c context.Context, event *bun.QueryEvent) context.Context {
	return c
}

func (h bunLogger) AfterQuery(c context.Context, event *bun.QueryEvent) {
	now := time.Now()
	dur := now.Sub(event.StartTime)

	args := []interface{}{
		"[bun]",
		now.Format(" 15:04:05.000 "),
		formatOperation(event),
		fmt.Sprintf(" %10s ", dur.Round(time.Microsecond)),
		event.Query,
	}

	if event.Err != nil {
		typ := reflect.TypeOf(event.Err).String()
		args = append(args,
			"\t",
			color.New(color.BgRed).Sprintf(" %s ", typ+": "+event.Err.Error()),
		)
	}

	fmt.Fprintln(h.writer, args...)
	// q := event.Query
	// fmt.Println(string(q))
}

func formatOperation(event *bun.QueryEvent) string {
	operation := event.Operation()
	return operationColor(operation).Sprintf(" %-16s ", operation)
}

func operationColor(operation string) *color.Color {
	switch operation {
	case "SELECT":
		return color.New(color.BgGreen, color.FgHiWhite)
	case "INSERT":
		return color.New(color.BgBlue, color.FgHiWhite)
	case "UPDATE":
		return color.New(color.BgYellow, color.FgHiBlack)
	case "DELETE":
		return color.New(color.BgMagenta, color.FgHiWhite)
	default:
		return color.New(color.BgWhite, color.FgHiBlack)
	}
}
