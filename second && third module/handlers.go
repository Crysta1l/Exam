package main

import (
	"context"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Pool      *pgxpool.Pool
	Templates *template.Template
}

func (app *App) SaveMaterial(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Cannot read form", 400)
		return
	}

	var material Material
	material.MaterialName = strings.TrimSpace(r.PostFormValue("material_name"))
	material.ItemsType = strings.TrimSpace(r.PostFormValue("items_type"))

	id, err := strconv.Atoi(r.PostFormValue("id"))
	if err != nil || id < 0 {
		http.Error(w, "Invalid material ID", 400)
		return
	}
	material.Id = id
	log.Printf("Получен id: %q", r.PostFormValue("id"))
	typeID, err := strconv.Atoi(r.PostFormValue("material_type_id"))
	if err != nil || typeID <= 0 {
		http.Error(w, "Choose a material type", 400)
		return
	}
	material.MaterialTypeId = typeID

	numbers := map[string]*float64{
		"price":             &material.Price,
		"amount":            &material.Amount,
		"min_amount":        &material.MinAmount,
		"items_per_package": &material.ItemsPerPackage,
	}

	for name, destination := range numbers {
		value, err := strconv.ParseFloat(r.PostFormValue(name), 64)
		if err != nil || math.IsNaN(value) || math.IsInf(value, 0) || value < 0 {
			http.Error(w, "Invalid number: "+name, 400)
			return
		}
		*destination = value
	}

	if material.MaterialName == "" || material.ItemsType == "" {
		http.Error(w, "Name and unit are required", 400)
		return
	}
	if material.ItemsPerPackage <= 0 {
		http.Error(w, "Quantity per package must be greater than zero", 400)
		return
	}

	args := []any{
		material.MaterialName, material.MaterialTypeId,
		material.Price, material.Amount, material.MinAmount,
		material.ItemsPerPackage, material.ItemsType,
	}

	query := `
        INSERT INTO materials_import (
            material_name, material_type_id, price, amount,
            min_amount, items_per_package, items_type
        ) VALUES ($1, $2, $3, $4, $5, $6, $7)
    `

	if material.Id != 0 {
		query = `
            UPDATE materials_import
            SET material_name=$1, material_type_id=$2, price=$3,
                amount=$4, min_amount=$5, items_per_package=$6,
                items_type=$7
            WHERE id=$8
        `
		args = append(args, material.Id)
	}

	result, err := app.Pool.Exec(r.Context(), query, args...)
	if err != nil {
		log.Println(err)
		http.Error(w, "Cannot save material. Check the values and name uniqueness.", 500)
		return
	}
	if result.RowsAffected() == 0 {
		http.Error(w, "Material not found", 404)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func GetMaterialByID(
	ctx context.Context,
	pool *pgxpool.Pool,
	id int,
) (Material, error) {
	var material Material

	err := pool.QueryRow(ctx, `
		SELECT 
			id, 
			material_name,
			material_type_id,
			price,
			amount,
			min_amount,
			items_per_package,
			items_type

		FROM materials_import
		WHERE id = $1

			`, id).Scan(
		&material.Id,
		&material.MaterialName,
		&material.MaterialTypeId,
		&material.Price,
		&material.Amount,
		&material.MinAmount,
		&material.ItemsPerPackage,
		&material.ItemsType,
	)

	if err != nil {
		return Material{}, err
	}

	return material, nil
}

func (app *App) EditMaterial(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "failed to get id", http.StatusBadRequest)
		return
	}

	material, err := GetMaterialByID(r.Context(), app.Pool, id)
	if err != nil {
		log.Println("GetMaterialById", err)
		http.Error(w, "cannot load material", http.StatusInternalServerError)
		return
	}

	types, err := GetMaterialType(r.Context(), app.Pool)
	if err != nil {
		http.Error(w, "cannot load material types", http.StatusInternalServerError)
		return
	}

	data := MaterialFormData{
		Material: material,
		Types:    types,
	}

	if err := app.Templates.Execute(w, data); err != nil {
		log.Println(err)
	}
}
