package main

import (
	"context"
	"log"
	"net/http"
	"text/template"

	"github.com/jackc/pgx/v5/pgxpool"
)

//  go get github.com/jackc/pgx/v5/pgxpool

type Material struct {
	Id           int
	MaterialName string
	Price        float64
	Amount       float64
}

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

	rows, err := pool.Query(ctx, "SELECT id,material_name,price,amount
		FROM materials_import WHERE m.material_type_id = t.id" 
	)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var materials []Material

	for rows.Next() {
		var material Material

		err := rows.Scan(
			&material.Id,
			&material.MaterialName,
			&material.Price,
			&material.Amount,
		)

		if err != nil {
			log.Fatal(err)
		}

		//log.Println(material)
		materials = append(materials, material)
	}
	log.Println("Materials uploaded", len(materials))

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	// frontend
	tmpl, err := template.ParseFiles("materials.html")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.Execute(w, materials); err != nil {
			log.Println(err)
		}
	})

	log.Println("http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
