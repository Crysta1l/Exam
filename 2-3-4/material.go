package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"strconv"
)

type Material struct {
	Id              int
	MaterialName    string
	MaterialType    string
	Price           float64
	Amount          float64
	MinAmount       float64
	ItemsType       string
	RequiredAmount  float64
	MaterialTypeId  int
	ItemsPerPackage float64
}

type MaterialType struct {
	Id   int
	Name string
}

type MaterialFormData struct {
	Material Material
	Types    []MaterialType
	Error    string
	Values   map[string]string
}

func formValues(material Material) map[string]string {
	number := func(value float64) string { return strconv.FormatFloat(value, 'f', -1, 64) }
	return map[string]string{
		"id":                strconv.Itoa(material.Id),
		"material_name":     material.MaterialName,
		"material_type_id":  strconv.Itoa(material.MaterialTypeId),
		"items_type":        material.ItemsType,
		"price":             number(material.Price),
		"amount":            number(material.Amount),
		"min_amount":        number(material.MinAmount),
		"items_per_package": number(material.ItemsPerPackage),
	}
}

func GetInitialMaterials(ctx context.Context, pool *pgxpool.Pool) ([]Material, error) {
	rows, err := pool.Query(ctx, `
  SELECT mi.id, mi.material_name, mt.material_type, mi.price, mi.amount,
         mi.min_amount, mi.items_type, ROUND(COALESCE(SUM(mp.material_per_unit), 0), 2)
  FROM materials_import mi
  JOIN material_types mt ON mi.material_type_id = mt.id
  LEFT JOIN materials_products mp ON mp.material_name_id = mi.id
  GROUP BY mi.id, mt.material_type
  ORDER BY mi.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var materials []Material
	for rows.Next() {
		var material Material
		if err := rows.Scan(&material.Id, &material.MaterialName, &material.MaterialType,
			&material.Price, &material.Amount, &material.MinAmount, &material.ItemsType,
			&material.RequiredAmount); err != nil {
			return nil, err
		}
		materials = append(materials, material)
	}
	return materials, rows.Err()
}

func GetMaterialType(ctx context.Context, pool *pgxpool.Pool) ([]MaterialType, error) {
	rows, err := pool.Query(ctx, "SELECT id, material_type FROM material_types ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var types []MaterialType
	for rows.Next() {
		var materialType MaterialType
		if err := rows.Scan(&materialType.Id, &materialType.Name); err != nil {
			return nil, err
		}
		types = append(types, materialType)
	}
	return types, rows.Err()
}

func GetMaterialByID(ctx context.Context, pool *pgxpool.Pool, id int) (Material, error) {
	var material Material
	err := pool.QueryRow(ctx, `
  SELECT mi.id, mi.material_name, mi.material_type_id, mi.price, mi.amount,
         mi.min_amount, mi.items_per_package, mi.items_type, mt.material_type
  FROM materials_import mi
  JOIN material_types mt ON mt.id = mi.material_type_id
  WHERE mi.id = $1`, id).Scan(
		&material.Id, &material.MaterialName, &material.MaterialTypeId, &material.Price,
		&material.Amount, &material.MinAmount, &material.ItemsPerPackage, &material.ItemsType,
		&material.MaterialType)
	return material, err
}
