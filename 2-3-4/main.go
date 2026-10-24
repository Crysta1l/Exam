package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:admin@localhost:5432/aroma_dom"
	}
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	err = pool.Ping(ctx)
	cancel()
	if err != nil {
		log.Fatal("Не удалось подключиться к БД: ", err)
	}

	templates, err := template.ParseFiles("materials.html", "material_form.html", "products.html", "error.html")
	if err != nil {
		log.Fatal(err)
	}
	app := &App{Pool: pool, Templates: templates}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", app.ListMaterials)
	mux.HandleFunc("GET /materials/new", app.NewMaterial)
	mux.HandleFunc("GET /materials/edit", app.EditMaterial)
	mux.HandleFunc("POST /materials/save", app.SaveMaterial)
	mux.HandleFunc("GET /materials/products", app.MaterialProducts)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		app.showError(w, 404, "Страница не найдена", "Вернитесь к списку материалов и выберите нужное действие.")
	})
	address := os.Getenv("HTTP_ADDR")
	if address == "" {
		address = "localhost:8080"
	}
	server := &http.Server{Addr: address, Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	log.Printf("Откройте http://%s", address)
	log.Fatal(server.ListenAndServe())
}
