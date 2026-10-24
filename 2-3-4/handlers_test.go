package main

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func validForm() map[string]string {
	return map[string]string{
		"id": "0", "material_name": "Тестовый материал", "material_type_id": "1",
		"items_type": "кг", "price": "10.50", "amount": "5", "min_amount": "1", "items_per_package": "2",
	}
}

func TestParseMaterial(t *testing.T) {
	material, message := parseMaterial(validForm())
	if message != "" || material.Price != 10.5 || material.ItemsPerPackage != 2 {
		t.Fatalf("%+v %s", material, message)
	}
	for _, test := range []struct{ key, value string }{
		{"id", ""}, {"id", "-1"}, {"material_name", " "}, {"material_name", strings.Repeat("я", 101)},
		{"items_type", ""}, {"material_type_id", "x"}, {"price", "-1"}, {"price", "NaN"},
		{"amount", "Inf"}, {"min_amount", "-0.01"}, {"items_per_package", "0"},
		{"price", "1.234"}, {"price", "10000000000"}, {"amount", "wrong"},
	} {
		t.Run(test.key+"="+test.value, func(t *testing.T) {
			values := validForm()
			values[test.key] = test.value
			if _, message := parseMaterial(values); message == "" {
				t.Fatal("ожидалась ошибка")
			}
		})
	}
	values := validForm()
	values["price"] = "0"
	values["min_amount"] = "0"
	if _, message := parseMaterial(values); message != "" {
		t.Fatal(message)
	}
	values["price"] = "12,50"
	if material, message := parseMaterial(values); message != "" || material.Price != 12.5 {
		t.Fatal(message)
	}
}

func TestTemplates(t *testing.T) {
	templates, err := template.ParseFiles("materials.html", "material_form.html", "products.html", "error.html")
	if err != nil {
		t.Fatal(err)
	}
	app := &App{Templates: templates}
	zero := 0
	tests := []struct {
		name string
		data any
	}{
		{"materials.html", struct {
			Materials []Material
			Saved     bool
		}{[]Material{{Id: 1, MaterialName: "<script>x</script>"}}, false}},
		{"material_form.html", MaterialFormData{Material: Material{MaterialTypeId: 1}, Types: []MaterialType{{1, "Воск"}}, Values: validForm(), Error: "Исправьте цену"}},
		{"products.html", ProductsPageData{Material: Material{Id: 1}, ProductTypes: []MaterialType{{1, "Свеча"}}, Result: &zero, Values: map[string]string{}}},
		{"error.html", struct{ Title, Message string }{"Ошибка", "Нет материала"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			app.render(response, 200, test.name, test.data)
			if response.Code != 200 {
				t.Fatal(response.Body.String())
			}
			if strings.Contains(response.Body.String(), "<script>x</script>") {
				t.Fatal("данные не экранированы")
			}
		})
	}
}

func TestInvalidMaterialID(t *testing.T) {
	templates, err := template.ParseFiles("error.html")
	if err != nil {
		t.Fatal(err)
	}
	app := &App{Templates: templates}
	for _, path := range []string{"/materials/edit", "/materials/edit?id=bad", "/materials/edit?id=-1"} {
		response := httptest.NewRecorder()
		app.EditMaterial(response, httptest.NewRequest(http.MethodGet, path, nil))
		if response.Code != 400 {
			t.Fatalf("%s: %d", path, response.Code)
		}
	}
}

func formRequest(values map[string]string) *http.Request {
	form := url.Values{}
	for key, value := range values {
		form.Set(key, value)
	}
	request := httptest.NewRequest(http.MethodPost, "/materials/save", strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return request
}
