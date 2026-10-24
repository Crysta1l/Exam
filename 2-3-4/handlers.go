package main

import (
	"bytes"
	"errors"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Pool      *pgxpool.Pool
	Templates *template.Template
}

func (app *App) render(w http.ResponseWriter, status int, name string, data any) {
	var buffer bytes.Buffer
	if err := app.Templates.ExecuteTemplate(&buffer, name, data); err != nil {
		log.Println(err)
		http.Error(w, "Ошибка отображения страницы. Вернитесь назад и повторите действие.", 500)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write(buffer.Bytes())
}

func (app *App) showError(w http.ResponseWriter, status int, title, message string) {
	app.render(w, status, "error.html", struct{ Title, Message string }{title, message})
}

func (app *App) ListMaterials(w http.ResponseWriter, r *http.Request) {
	materials, err := GetInitialMaterials(r.Context(), app.Pool)
	if err != nil {
		log.Println(err)
		app.showError(w, 500, "Ошибка загрузки", "Не удалось получить материалы. Проверьте подключение к БД и обновите страницу.")
		return
	}
	app.render(w, 200, "materials.html", struct {
		Materials []Material
		Saved     bool
	}{materials, r.URL.Query().Get("saved") == "1"})
}

func (app *App) NewMaterial(w http.ResponseWriter, r *http.Request) {
	app.showForm(w, r, MaterialFormData{Values: formValues(Material{})}, 200)
}

func (app *App) showForm(w http.ResponseWriter, r *http.Request, data MaterialFormData, status int) {
	types, err := GetMaterialType(r.Context(), app.Pool)
	if err != nil {
		log.Println(err)
		app.showError(w, 500, "Ошибка загрузки", "Не удалось загрузить типы материалов. Проверьте подключение к БД и повторите действие.")
		return
	}
	data.Types = types
	app.render(w, status, "material_form.html", data)
}

func (app *App) EditMaterial(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id <= 0 {
		app.showError(w, 400, "Неверный номер материала", "Откройте список материалов и выберите материал для редактирования.")
		return
	}
	material, err := GetMaterialByID(r.Context(), app.Pool, id)
	if errors.Is(err, pgx.ErrNoRows) {
		app.showError(w, 404, "Материал не найден", "Такого материала уже нет. Вернитесь к списку и выберите существующий материал.")
		return
	}
	if err != nil {
		log.Println(err)
		app.showError(w, 500, "Ошибка загрузки", "Не удалось загрузить материал. Проверьте подключение к БД и повторите действие.")
		return
	}
	app.showForm(w, r, MaterialFormData{Material: material, Values: formValues(material)}, 200)
}

func parseMaterial(values map[string]string) (Material, string) {
	var material Material
	var messages []string
	material.Id, _ = strconv.Atoi(values["id"])
	if id, err := strconv.Atoi(values["id"]); err != nil || id < 0 {
		messages = append(messages, "Неверный номер материала. Откройте форму заново из списка.")
	}
	material.MaterialName = strings.TrimSpace(values["material_name"])
	material.ItemsType = strings.TrimSpace(values["items_type"])
	if length := utf8.RuneCountInString(material.MaterialName); length < 1 || length > 100 {
		messages = append(messages, "Введите название длиной от 1 до 100 символов.")
	}
	if length := utf8.RuneCountInString(material.ItemsType); length < 1 || length > 10 {
		messages = append(messages, "Введите единицу измерения длиной от 1 до 10 символов, например кг или шт.")
	}
	material.MaterialTypeId, _ = strconv.Atoi(values["material_type_id"])
	if material.MaterialTypeId <= 0 {
		messages = append(messages, "Выберите тип материала из списка.")
	}
	fields := []struct {
		key, label string
		target     *float64
	}{
		{"price", "Цена", &material.Price},
		{"amount", "Остаток", &material.Amount},
		{"min_amount", "Минимальное количество", &material.MinAmount},
		{"items_per_package", "Количество в упаковке", &material.ItemsPerPackage},
	}
	for _, field := range fields {
		text := strings.ReplaceAll(strings.TrimSpace(values[field.key]), ",", ".")
		value, err := strconv.ParseFloat(text, 64)
		if err != nil || math.IsNaN(value) || math.IsInf(value, 0) || value < 0 || value > 9999999999.99 {
			messages = append(messages, field.label+": введите число от 0 до 9999999999.99.")
			continue
		}
		if dot := strings.IndexByte(text, '.'); strings.ContainsAny(text, "eE") || (dot >= 0 && len(text)-dot-1 > 2) {
			messages = append(messages, field.label+": укажите не более двух знаков после запятой.")
			continue
		}
		if field.key == "items_per_package" && value <= 0 {
			messages = append(messages, "Количество в упаковке должно быть больше нуля.")
		}
		*field.target = value
	}
	return material, strings.Join(messages, "\n")
}

func (app *App) SaveMaterial(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 64<<10)
	if err := r.ParseForm(); err != nil {
		app.showError(w, 400, "Ошибка формы", "Не удалось прочитать форму. Откройте её заново и повторите ввод.")
		return
	}
	values := make(map[string]string)
	for key := range r.PostForm {
		values[key] = r.PostFormValue(key)
	}
	material, message := parseMaterial(values)
	if message != "" {
		app.showForm(w, r, MaterialFormData{Material: material, Values: values, Error: message}, 400)
		return
	}
	args := []any{material.MaterialName, material.MaterialTypeId, material.Price,
		material.Amount, material.MinAmount, material.ItemsPerPackage, material.ItemsType}
	query := `INSERT INTO materials_import
  (material_name, material_type_id, price, amount, min_amount, items_per_package, items_type)
  VALUES ($1, $2, $3, $4, $5, $6, $7)`
	if material.Id != 0 {
		query = `UPDATE materials_import SET material_name=$1, material_type_id=$2,
   price=$3, amount=$4, min_amount=$5, items_per_package=$6, items_type=$7 WHERE id=$8`
		args = append(args, material.Id)
	}
	result, err := app.Pool.Exec(r.Context(), query, args...)
	if err != nil {
		log.Println(err)
		message = "Не удалось сохранить материал. Проверьте подключение к БД и повторите сохранение."
		status := 500
		var dbError *pgconn.PgError
		if errors.As(err, &dbError) {
			switch dbError.Code {
			case "23505":
				message = "Материал с таким названием уже существует. Укажите другое название."
				status = 400
			case "23503":
				message = "Выбранный тип материала не существует. Выберите тип из актуального списка."
				status = 400
			}
		}
		app.showForm(w, r, MaterialFormData{Material: material, Values: values, Error: message}, status)
		return
	}
	if result.RowsAffected() == 0 {
		app.showError(w, 404, "Материал не найден", "Материал удалён или не существует. Вернитесь к списку материалов.")
		return
	}
	http.Redirect(w, r, "/?saved=1", http.StatusSeeOther)
}
