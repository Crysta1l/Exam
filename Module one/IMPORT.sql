CREATE TABLE product_types (
	id SERIAL PRIMARY KEY,
	product_type VARCHAR(100) NOT NULL UNIQUE ,
	coefficient NUMERIC(5,2) NOT NULL
);

CREATE TABLE material_types(
	id SERIAL PRIMARY KEY,
	material_type VARCHAR(100) NOT NULL UNIQUE,
	loss_percent NUMERIC(5,2) NOT NULL
);

CREATE TABLE materials_import(
	id SERIAL PRIMARY KEY,
	material_name VARCHAR(100) NOT NULL UNIQUE,
	material_type_id INTEGER REFERENCES material_types(id) NOT NULL,

	price NUMERIC(12,2) NOT NULL,
	amount NUMERIC(12,2) NOT NULL CHECK(amount >= 0),
	min_amount NUMERIC(12,2) NOT NULL CHECK(min_amount >= 0),
	items_per_package  NUMERIC(12,2) NOT NULL CHECK(items_per_package > 0),
	items_type VARCHAR(10) NOT NULL
	
	
);


CREATE TABLE products_import(
	id SERIAL PRIMARY KEY,
	article_number VARCHAR(50) NOT NULL UNIQUE,
	product_type_id INTEGER REFERENCES product_types(id) NOT NULL,
	product_name VARCHAR(100) NOT NULL,
	min_partner_price NUMERIC(12,2) NOT NULL CHECK(min_partner_price > 0)
);


CREATE TABLE materials_products(
	id SERIAL PRIMARY KEY,
	material_name_id INTEGER REFERENCES materials_import(id) NOT NULL,
	product_name_id INTEGER REFERENCES products_import(id) NOT NULL,
	material_per_unit NUMERIC(12,4) NOT NULL CHECK(material_per_unit > 0),
	UNIQUE(material_name_id, product_name_id)
);




-- XXX

CREATE TEMP TABLE temp_product_types(
	product_type TEXT,
	coefficient TEXT
);

COPY temp_product_types
FROM 'D:\IMPORTS\product_types.csv'
WITH (FORMAT CSV, HEADER TRUE, DELIMITER ';', ENCODING 'UTF8')


INSERT INTO product_types (product_type, coefficient)
SELECT product_type, 
	REPLACE(coefficient, ',','.')::NUMERIC
FROM temp_product_types


SELECT * FROM product_types

-- XXX

CREATE TEMP TABLE temp_material_types(
	material_type TEXT,
	loss_percent TEXT
);


COPY temp_material_types
FROM 'D:\IMPORTS\Material_type.csv'
WITH (FORMAT CSV, HEADER TRUE, DELIMITER ';', ENCODING 'UTF-8');


INSERT INTO material_types (material_type, loss_percent)
SELECT material_type,
	REPLACE(REPLACE(loss_percent, ',','.'),'%','')::NUMERIC
FROM temp_material_types;

SELECT * FROM material_types

-- XXX


CREATE TEMP TABLE temp_materials_import(
	material_name TEXT,
	material_type TEXT,
	price TEXT,
	amount TEXT,
	min_amount TEXT,
	items_per_package  TEXT,
	items_type TEXT
);

COPY temp_materials_import
FROM 'D:\IMPORTS\3_Materials_import.csv'
WITH (FORMAT CSV, HEADER TRUE, DELIMITER ';', ENCODING 'UTF-8');

INSERT INTO materials_import(material_name, material_type_id, price, amount, min_amount, items_per_package, items_type)
SELECT 
	m.material_name,
	t.id,
	REPLACE(m.price, ',','.')::NUMERIC,
	REPLACE(m.amount, ',','.')::NUMERIC,
	REPLACE(m.min_amount, ',','.')::NUMERIC,
	REPLACE(m.items_per_package, ',','.')::NUMERIC,
	m.items_type
FROM temp_materials_import m
JOIN material_types t ON m.material_type = t.material_type;

SELECT * FROM materials_import 

--- XXX

CREATE TEMP TABLE temp_products_import(
	article_number TEXT,
	product_type TEXT,
	product_name TEXT,
	min_partner_price TEXT
);


COPY temp_products_import
FROM 'D:\IMPORTS\4_Products_import.csv'
WITH (FORMAT CSV, HEADER TRUE, DELIMITER ';', ENCODING 'UTF-8');

INSERT INTO products_import(
	article_number,
	product_type_id,
	product_name,
	min_partner_price)
SELECT
	m.article_number,
	pt.id,
	product_name,
	min_partner_price::NUMERIC
FROM temp_products_import m
JOIN product_types pt ON m.product_type = pt.product_type;


SELECT * FROM products_import 
--- XXX

CREATE TEMP TABLE temp_material_products(
	material_name TEXT,
	product_name TEXT,
	material_per_unit TEXT
);

COPY temp_material_products
FROM 'D:\IMPORTS\5_Material_products_import.csv'
WITH (FORMAT CSV, HEADER TRUE, DELIMITER ';', ENCODING 'UTF-8');

INSERT INTO materials_products(
	material_name_id,
	product_name_id,
	material_per_unit)
SELECT
	m.id,
	p.id,
	REPLACE(tmp.material_per_unit, ',','.')::NUMERIC
FROM temp_material_products tmp
JOIN materials_import m ON m.material_name  = tmp.material_name
JOIN products_import p ON p.product_name = tmp.product_name;


SELECT * FROM materials_products mp 
--- XXX















