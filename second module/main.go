package main

import (
	"context"
	"html/template"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

//  go get github.com/jackc/pgx/v5/pgxpool

type Material struct {
	Id           int
	MaterialName string
	MaterialType string
	Price        float64
	Amount       float64
	MinAmount    float64
	ItemsType    string
	RequiredCost float64
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

	rows, err := pool.Query(ctx, `
		SELECT
			mi.id,
			mi.material_name,
			mt.material_type,
			mi.price,
			mi.amount,
			mi.min_amount,
			mi.items_type,
			ROUND(
				COALESCE(SUM(mp.material_per_unit), 0) * mi.price,
				2
			) AS required_cost
		FROM materials_import mi
		JOIN material_types mt ON mi.material_type_id  = mt.id 
		LEFT JOIN materials_products mp ON mp.material_name_id = mi.id
		GROUP BY mi.id, mt.material_type 
		ORDER BY mi.id;
	`)
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
			&material.MaterialType,
			&material.Price,
			&material.Amount,
			&material.MinAmount,
			&material.ItemsType,
			&material.RequiredCost,
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
