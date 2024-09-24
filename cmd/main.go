package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"jordi.tfg.rewrite/internal/data"
)

var secretKeyJwt = []byte("your-secret-key")

type Config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  time.Duration
	}
}

type Application struct {
	Config
	logger *slog.Logger
	Models data.Models
}

func HealthCheck(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(200)
	fmt.Fprint(w, "Welcome!\n everything is fine")
}

func main() {
	adrr := flag.String("addr", ":4400", "server route")
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := OpenDB()
	if err != nil {
		fmt.Print(err)
		logger.Error("Error opening db..")
		os.Exit(1)
	}

	defer db.Close()

	app := Application{
		Models: data.NewModels(db),
	}

	srv := &http.Server{
		Addr:         *adrr,
		Handler:      app.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(("Server running.."))
}

func OpenDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", "postgres://jordisalazarbadia@localhost/crowdfun?sslmode=disable")
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)

	if err != nil {
		db.Close()
		return nil, err

	}

	return db, nil

}
