package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"go.uber.org/fx"
	"mem.com/mem-go/db"
	h "mem.com/mem-go/handlers"
)

func main() {
	fx.New(
		fx.Provide(
			NewDBConn,
			NewHTTPServer,
			NewDBQ,

			fx.Annotate(
				h.NewMuxRouter,
				fx.ParamTags(`group:"routes"`),
			),
			h.AsRoute(h.NewBooksHandler),
			h.AsRoute(h.NewGetBookHandler),
			h.AsRoute(h.NewEntryHandler),
			h.AsRoute(h.NewDeleteEntryHandler),
		),
		fx.Invoke(func(*http.Server) {}),
	).Run()
}

func NewDBQ(conn *sql.DB) *db.Queries {
	return db.New(conn)
}

func NewDBConn(lc fx.Lifecycle) *sql.DB {
	conn, err := sql.Open(os.Getenv("DB_DRIVER"), os.Getenv("DB_STRING"))
	if err != nil {
		panic(err)
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return conn.Close()
		},
	})

	return conn
}

func NewHTTPServer(lc fx.Lifecycle, r *mux.Router) *http.Server {
	srv := &http.Server{Addr: ":8080", Handler: r}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			fmt.Println("Starting HTTP server at", srv.Addr)

			go srv.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
	return srv
}
