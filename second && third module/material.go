package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Material struct {
	Id              int
	MaterialName    string
	MaterialType    string
	Price           float64
	Amount          float64
	MinAmount       float64
	ItemsType       string
	RequiredCost    float64
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
	Erorr    string
}

func GetInitialMaterials(
	ctx context.Context,
	pool *pgxpool.Pool,
) ([]Material, error) {

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
		return nil, err
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
			return nil, err
		}

		//log.Println(material)
		materials = append(materials, material)
	}
	log.Println("Materials uploaded", len(materials))

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return materials, nil

}

func GetMaterialType(
	ctx context.Context,
	pool *pgxpool.Pool,
) ([]MaterialType, error) {

	rows, err := pool.Query(ctx, `
		SELECT
			id,
			material_type
		FROM material_types
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var materialsTypes []MaterialType

	for rows.Next() {
		var materialType MaterialType

		err := rows.Scan(
			&materialType.Id,
			&materialType.Name,
		)
		if err != nil {
			return nil, err
		}
		materialsTypes = append(materialsTypes, materialType)
	}

	log.Println("MaterialTypes uploaded")

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return materialsTypes, nil
}
