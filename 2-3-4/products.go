package main

import (
	"aroma_dom/calculation"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Product struct {
	Id             int
	Article        string
	Name           string
	Type           string
	MaterialAmount float64
}

type ProductsPageData struct {
	Material     Material
	Products     []Product
	ProductTypes []MaterialType
	Values       map[string]string
	Result       *int
	Error        string
}

func GetProductsByMaterial(ctx context.Context, pool *pgxpool.Pool, materialID int) ([]Product, error) {
	rows, err := pool.Query(ctx, `SELECT pi.id, pi.article_number, pi.product_name, pt.product_type, mp.material_per_unit
		FROM materials_products mp
		JOIN products_import pi ON pi.id = mp.product_name_id
		JOIN product_types pt ON pt.id = pi.product_type_id
		WHERE mp.material_name_id = $1 ORDER BY pi.article_number`, materialID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.Id, &product.Article, &product.Name, &product.Type, &product.MaterialAmount); err != nil {
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
		app.showError(w, 500, "Ошибка загрузки", "Не удалось получить продукцию. Повторите действие после проверки БД.")
		return
	}
	rows, err := app.Pool.Query(r.Context(), "SELECT id, product_type FROM product_types ORDER BY id")
	if err != nil {
		app.showError(w, 500, "Ошибка загрузки", "Не удалось получить типы продукции. Обновите страницу.")
		return
	}
	var types []MaterialType
	for rows.Next() {
		var productType MaterialType
		if err := rows.Scan(&productType.Id, &productType.Name); err != nil {
			rows.Close()
			app.showError(w, 500, "Ошибка загрузки", "Не удалось прочитать типы продукции. Обновите страницу.")
			return
		}
		types = append(types, productType)
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		app.showError(w, 500, "Ошибка загрузки", "Не удалось прочитать типы продукции. Обновите страницу.")
		return
	}
	data := ProductsPageData{Material: material, Products: products, ProductTypes: types,
		Values: map[string]string{"raw_amount": "0", "parameter1": "1", "parameter2": "1"}}
	status := 200
	if r.URL.Query().Get("calculate") == "1" {
		for _, key := range []string{"product_type_id", "raw_amount", "parameter1", "parameter2"} {
			data.Values[key] = r.URL.Query().Get(key)
		}
		productType, typeErr := strconv.Atoi(data.Values["product_type_id"])
		rawAmount, amountErr := strconv.Atoi(data.Values["raw_amount"])
		p1, p1Err := strconv.ParseFloat(strings.ReplaceAll(data.Values["parameter1"], ",", "."), 64)
		p2, p2Err := strconv.ParseFloat(strings.ReplaceAll(data.Values["parameter2"], ",", "."), 64)
		result := -1
		if typeErr == nil && amountErr == nil && p1Err == nil && p2Err == nil {
			result = calculation.CalculateProduction(r.Context(), app.Pool, productType, material.MaterialTypeId, rawAmount, p1, p2)
		}
		if result == -1 {
			data.Error = "Выберите существующий тип продукции, целое неотрицательное количество сырья и два положительных параметра. Результат метода: -1. При корректных параметрах проверьте доступность БД."
			status = 400
		} else {
			data.Result = &result
		}
	}
	app.render(w, status, "products.html", data)
}
