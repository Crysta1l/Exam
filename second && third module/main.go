package main

import (
	"context"
	"html/template"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

//  go get github.com/jackc/pgx/v5/pgxpool

func main() {
	connStr := "postgres://postgres:admin@localhost:5432/aroma_dom"

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, connStr)

	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	err = pool.Ping(ctx)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connection established")

	// frontend
	tmpl, err := template.ParseFiles("materials.html")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		materials, err := GetInitialMaterials(r.Context(), pool)
		if err != nil {
			log.Println(err)
			http.Error(w, "Cant load materials", http.StatusInternalServerError)
			return
		}
		if err := tmpl.Execute(w, materials); err != nil {
			log.Println(err)
		}
	})

	// form
	formTmpt, err := template.ParseFiles("material_form.html")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/materials/new", func(w http.ResponseWriter, r *http.Request) {
		types, err := GetMaterialType(r.Context(), pool)
		if err != nil {
			log.Println(err)
			http.Error(w, "cant load materials", 500)
			return
		}

		data := MaterialFormData{
			Types: types,
		}

		if err := formTmpt.Execute(w, data); err != nil {
			log.Println(err)
		}
	})

	// save

	app := &App{
		Pool:      pool,
		Templates: formTmpt,
	}

	http.HandleFunc("/materials/save", app.SaveMaterial)

	http.HandleFunc("/materials/edit", app.EditMaterial)

	http.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Println("http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
