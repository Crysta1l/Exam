package main

import (
	"aroma_dom/calculation"
	"context"
	"html/template"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestApplicationWithPostgres(t *testing.T) {
	connection := os.Getenv("TEST_DATABASE_URL")
	if connection == "" {
		t.Skip("Для проверки с PostgreSQL задайте TEST_DATABASE_URL")
	}
	ctx := context.Background()
	admin, err := pgxpool.New(ctx, connection)
	if err != nil {
		t.Fatal(err)
	}
	defer admin.Close()
	schema := "exam_test_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	quoted := pgx.Identifier{schema}.Sanitize()
	if _, err := admin.Exec(ctx, "CREATE SCHEMA "+quoted); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if _, err := admin.Exec(ctx, "DROP SCHEMA "+quoted+" CASCADE"); err != nil {
			t.Error(err)
		}
	}()
	config, err := pgxpool.ParseConfig(connection)
	if err != nil {
		t.Fatal(err)
	}
	config.ConnConfig.RuntimeParams["search_path"] = schema
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	_, err = pool.Exec(ctx, `
  CREATE TABLE material_types (id SERIAL PRIMARY KEY, material_type TEXT NOT NULL UNIQUE, loss_percent NUMERIC(5,2) NOT NULL);
  CREATE TABLE product_types (id SERIAL PRIMARY KEY, product_type TEXT NOT NULL UNIQUE, coefficient NUMERIC(5,2) NOT NULL);
  CREATE TABLE materials_import (
   id SERIAL PRIMARY KEY, material_name VARCHAR(100) NOT NULL UNIQUE,
   material_type_id INTEGER NOT NULL REFERENCES material_types(id),
   price NUMERIC(12,2) NOT NULL, amount NUMERIC(12,2) NOT NULL CHECK (amount >= 0),
   min_amount NUMERIC(12,2) NOT NULL CHECK (min_amount >= 0),
   items_per_package NUMERIC(12,2) NOT NULL CHECK (items_per_package > 0), items_type VARCHAR(10) NOT NULL);
  CREATE TABLE products_import (id SERIAL PRIMARY KEY, article_number TEXT UNIQUE,
   product_type_id INTEGER REFERENCES product_types(id), product_name TEXT, min_partner_price NUMERIC(12,2));
  CREATE TABLE materials_products (id SERIAL PRIMARY KEY,
   material_name_id INTEGER REFERENCES materials_import(id),
   product_name_id INTEGER REFERENCES products_import(id), material_per_unit NUMERIC(12,4),
   UNIQUE(material_name_id, product_name_id));
  INSERT INTO material_types (material_type, loss_percent) VALUES ('Воск', 10);
  INSERT INTO product_types (product_type, coefficient) VALUES ('Свеча', 1);
  INSERT INTO materials_import (material_name, material_type_id, price, amount, min_amount, items_per_package, items_type)
   VALUES ('Воск тест', 1, 890, 120, 40, 20, 'кг'), ('Без продукции', 1, 0, 0, 0, 1, 'кг');
  INSERT INTO products_import (article_number, product_type_id, product_name, min_partner_price)
   VALUES ('T-1', 1, 'Свеча 1', 100), ('T-2', 1, 'Свеча 2', 200);
  INSERT INTO materials_products (material_name_id, product_name_id, material_per_unit)
   VALUES (1, 1, 0.18), (1, 2, 0.22);
 `)
	if err != nil {
		t.Fatal(err)
	}
	templates, err := template.ParseFiles("materials.html", "material_form.html", "products.html", "error.html")
	if err != nil {
		t.Fatal(err)
	}
	app := &App{Pool: pool, Templates: templates}
	get := func(handler http.HandlerFunc, path string, want int, text string) {
		t.Helper()
		response := httptest.NewRecorder()
		handler(response, httptest.NewRequest("GET", path, nil))
		if response.Code != want || !strings.Contains(response.Body.String(), text) {
			t.Fatalf("%s: статус %d, ожидалось %d с %q; %s", path, response.Code, want, text, response.Body.String())
		}
	}
	get(app.ListMaterials, "/", 200, "0.40 кг")
	get(app.ListMaterials, "/", 200, "0.00 кг")
	get(app.NewMaterial, "/materials/new", 200, "Воск")
	get(app.EditMaterial, "/materials/edit?id=1", 200, "Воск тест")
	get(app.EditMaterial, "/materials/edit?id=99999", 404, "Материал не найден")
	get(app.MaterialProducts, "/materials/products?id=1", 200, "Свеча 1")
	get(app.MaterialProducts, "/materials/products?id=2", 200, "пока не используется")
	get(app.MaterialProducts, "/materials/products?id=99999", 404, "Материал не найден")
	get(app.MaterialProducts, "/materials/products?id=1&calculate=1&product_type_id=1&raw_amount=100&parameter1=2&parameter2=5", 200, "9 шт.")
	get(app.MaterialProducts, "/materials/products?id=1&calculate=1&product_type_id=999&raw_amount=100&parameter1=2&parameter2=5", 400, "Результат метода: -1")
	if got := calculation.CalculateProduction(ctx, pool, 1, 999, 100, 2, 5); got != -1 {
		t.Fatal(got)
	}
	if got := calculation.CalculateProduction(ctx, pool, 1, 1, 0, 2, 5); got != 0 {
		t.Fatal(got)
	}

	post := func(values map[string]string, want int, text string) *httptest.ResponseRecorder {
		t.Helper()
		response := httptest.NewRecorder()
		app.SaveMaterial(response, formRequest(values))
		if response.Code != want || (text != "" && !strings.Contains(response.Body.String(), text)) {
			t.Fatalf("save: %d, ожидалось %d с %q; %s", response.Code, want, text, response.Body.String())
		}
		return response
	}
	values := validForm()
	response := post(values, 303, "")
	if response.Header().Get("Location") != "/?saved=1" {
		t.Fatal(response.Header())
	}
	var id int
	if err := pool.QueryRow(ctx, "SELECT id FROM materials_import WHERE material_name=$1", values["material_name"]).Scan(&id); err != nil {
		t.Fatal(err)
	}
	values["id"] = strconv.Itoa(id)
	values["material_name"] = "Изменённый материал"
	values["price"] = "0"
	values["amount"] = "12.50"
	values["items_per_package"] = "3"
	post(values, 303, "")
	updated, err := GetMaterialByID(ctx, pool, id)
	if err != nil || updated.Price != 0 || updated.ItemsPerPackage != 3 || updated.Amount != 12.5 {
		t.Fatalf("%+v %v", updated, err)
	}
	values["id"] = "0"
	post(values, 400, "уже существует")
	values["material_name"] = "Сохранить введённое имя"
	values["min_amount"] = "-1"
	post(values, 400, values["material_name"])
	var count int
	if err := pool.QueryRow(ctx, "SELECT count(*) FROM materials_import").Scan(&count); err != nil || count != 3 {
		t.Fatalf("count %d: %v", count, err)
	}
	values = validForm()
	values["material_type_id"] = "99999"
	post(values, 400, "тип материала не существует")
	values = validForm()
	values["id"] = "99999"
	post(values, 404, "Материал не найден")
	values = validForm()
	values["price"] = "NaN"
	post(values, 400, "Цена")
	values = validForm()
	values["items_per_package"] = "0"
	post(values, 400, "больше нуля")
}
