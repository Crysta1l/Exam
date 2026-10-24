package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Product struct {
	Name           string
	MaterialAmount float64
}

type ProductsPageData struct {
	Material Material
	Products []Product
}

func GetProductsByMaterial(ctx context.Context, pool *pgxpool.Pool, materialID int) ([]Product, error) {
	rows, err := pool.Query(ctx, `
        SELECT pi.product_name, mp.material_per_unit
        FROM materials_products mp
        JOIN products_import pi ON pi.id = mp.product_name_id
        WHERE mp.material_name_id = $1
        ORDER BY pi.product_name
    `, materialID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.Name, &product.MaterialAmount); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, rows.Err()
}

func (app *App) MaterialProducts(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id <= 0 {
		app.showError(w, 400, "Неверный номер материала", "Выберите материал из списка и нажмите «Продукция».")
		return
	}

	material, err := GetMaterialByID(r.Context(), app.Pool, id)
	if errors.Is(err, pgx.ErrNoRows) {
		app.showError(w, 404, "Материал не найден", "Вернитесь к списку и выберите существующий материал.")
		return
	}
	if err != nil {
		log.Println(err)
		app.showError(w, 500, "Ошибка загрузки", "Не удалось получить материал. Проверьте подключение к БД.")
		return
	}

	products, err := GetProductsByMaterial(r.Context(), app.Pool, id)
	if err != nil {
		log.Println(err)
		app.showError(w, 500, "Ошибка загрузки", "Не удалось получить продукцию. Проверьте подключение к БД и повторите действие.")
		return
	}

	app.render(w, 200, "products.html", ProductsPageData{
		Material: material,
		Products: products,
	})
}
