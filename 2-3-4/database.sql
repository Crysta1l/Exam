--
-- PostgreSQL database dump
--

-- Выполнять на пустой БД. Содержит структуру и исходные данные, CSV не требуются.

-- Dumped from database version 18.4
-- Dumped by pg_dump version 18.4

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: material_types; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.material_types (
    id integer NOT NULL,
    material_type character varying(100) NOT NULL,
    loss_percent numeric(5,2) NOT NULL,
    CONSTRAINT material_types_loss_percent_check CHECK (((loss_percent >= (0)::numeric) AND (loss_percent <= (100)::numeric)))
);


--
-- Name: material_types_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.material_types_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: material_types_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.material_types_id_seq OWNED BY public.material_types.id;


--
-- Name: materials_import; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.materials_import (
    id integer NOT NULL,
    material_name character varying(100) NOT NULL,
    material_type_id integer NOT NULL,
    price numeric(12,2) NOT NULL,
    amount numeric(12,2) NOT NULL,
    min_amount numeric(12,2) NOT NULL,
    items_per_package numeric(12,2) NOT NULL,
    items_type character varying(10) NOT NULL,
    CONSTRAINT materials_import_amount_check CHECK ((amount >= (0)::numeric)),
    CONSTRAINT materials_import_items_per_package_check CHECK ((items_per_package > (0)::numeric)),
    CONSTRAINT materials_import_min_amount_check CHECK ((min_amount >= (0)::numeric)),
    CONSTRAINT materials_import_price_check CHECK ((price >= (0)::numeric))
);


--
-- Name: materials_import_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.materials_import_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: materials_import_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.materials_import_id_seq OWNED BY public.materials_import.id;


--
-- Name: materials_products; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.materials_products (
    id integer NOT NULL,
    material_name_id integer NOT NULL,
    product_name_id integer NOT NULL,
    material_per_unit numeric(12,4) NOT NULL,
    CONSTRAINT materials_products_material_per_unit_check CHECK ((material_per_unit > (0)::numeric))
);


--
-- Name: materials_products_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.materials_products_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: materials_products_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.materials_products_id_seq OWNED BY public.materials_products.id;


--
-- Name: product_types; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.product_types (
    id integer NOT NULL,
    product_type character varying(100) NOT NULL,
    coefficient numeric(5,2) NOT NULL,
    CONSTRAINT product_types_coefficient_check CHECK ((coefficient > (0)::numeric))
);


--
-- Name: product_types_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.product_types_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: product_types_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.product_types_id_seq OWNED BY public.product_types.id;


--
-- Name: products_import; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.products_import (
    id integer NOT NULL,
    article_number character varying(50) NOT NULL,
    product_type_id integer NOT NULL,
    product_name character varying(100) NOT NULL,
    min_partner_price numeric(12,2) NOT NULL,
    CONSTRAINT products_import_min_partner_price_check CHECK ((min_partner_price > (0)::numeric))
);


--
-- Name: products_import_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.products_import_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: products_import_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.products_import_id_seq OWNED BY public.products_import.id;


--
-- Name: material_types id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.material_types ALTER COLUMN id SET DEFAULT nextval('public.material_types_id_seq'::regclass);


--
-- Name: materials_import id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.materials_import ALTER COLUMN id SET DEFAULT nextval('public.materials_import_id_seq'::regclass);


--
-- Name: materials_products id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.materials_products ALTER COLUMN id SET DEFAULT nextval('public.materials_products_id_seq'::regclass);


--
-- Name: product_types id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_types ALTER COLUMN id SET DEFAULT nextval('public.product_types_id_seq'::regclass);


--
-- Name: products_import id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.products_import ALTER COLUMN id SET DEFAULT nextval('public.products_import_id_seq'::regclass);


--
-- Data for Name: material_types; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.material_types VALUES (1, 'Соевый воск', 0.12);
INSERT INTO public.material_types VALUES (2, 'Пчелиный воск', 0.18);
INSERT INTO public.material_types VALUES (3, 'Ароматическое масло', 0.25);
INSERT INTO public.material_types VALUES (4, 'Краситель', 0.08);
INSERT INTO public.material_types VALUES (5, 'Фитиль', 0.03);
INSERT INTO public.material_types VALUES (6, 'Стеклянная тара', 0.04);
INSERT INTO public.material_types VALUES (7, 'Керамическая тара', 0.06);
INSERT INTO public.material_types VALUES (8, 'Декор', 0.10);
INSERT INTO public.material_types VALUES (9, 'Упаковка', 0.05);
INSERT INTO public.material_types VALUES (10, 'Этикетка', 0.02);


--
-- Data for Name: materials_import; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.materials_import VALUES (1, 'Соевый воск премиум', 1, 890.00, 120.00, 40.00, 20.00, 'кг');
INSERT INTO public.materials_import VALUES (2, 'Соевый воск EcoBlend', 1, 760.00, 95.00, 30.00, 25.00, 'кг');
INSERT INTO public.materials_import VALUES (3, 'Пчелиный воск натуральный', 2, 1450.00, 48.00, 15.00, 10.00, 'кг');
INSERT INTO public.materials_import VALUES (4, 'Масло лаванды', 3, 3200.00, 18.50, 6.00, 1.00, 'л');
INSERT INTO public.materials_import VALUES (5, 'Масло ванили', 3, 2950.00, 22.00, 8.00, 1.00, 'л');
INSERT INTO public.materials_import VALUES (6, 'Масло цитрусовое', 3, 2870.00, 16.70, 5.00, 1.00, 'л');
INSERT INTO public.materials_import VALUES (7, 'Краситель янтарный', 4, 640.00, 6.20, 2.00, 0.50, 'кг');
INSERT INTO public.materials_import VALUES (8, 'Краситель розовый', 4, 615.00, 5.10, 1.50, 0.50, 'кг');
INSERT INTO public.materials_import VALUES (9, 'Фитиль хлопковый 8 см', 5, 4.50, 2500.00, 800.00, 500.00, 'шт');
INSERT INTO public.materials_import VALUES (10, 'Фитиль хлопковый 12 см', 5, 5.80, 1800.00, 600.00, 400.00, 'шт');
INSERT INTO public.materials_import VALUES (11, 'Стакан прозрачный 200 мл', 6, 28.00, 950.00, 300.00, 100.00, 'шт');
INSERT INTO public.materials_import VALUES (12, 'Стакан матовый 250 мл', 6, 34.00, 620.00, 200.00, 60.00, 'шт');
INSERT INTO public.materials_import VALUES (13, 'Чаша керамическая белая', 7, 72.00, 280.00, 80.00, 40.00, 'шт');
INSERT INTO public.materials_import VALUES (14, 'Сухоцветы лаванды', 8, 520.00, 4.80, 1.50, 0.30, 'кг');
INSERT INTO public.materials_import VALUES (15, 'Лента сатиновая бежевая', 9, 12.00, 340.00, 120.00, 50.00, 'м');
INSERT INTO public.materials_import VALUES (16, 'Коробка подарочная малая', 9, 18.00, 420.00, 150.00, 50.00, 'шт');
INSERT INTO public.materials_import VALUES (17, 'Этикетка фирменная круглая', 10, 3.20, 1800.00, 500.00, 250.00, 'шт');
INSERT INTO public.materials_import VALUES (18, 'Наклейка предупреждающая', 10, 2.10, 1600.00, 400.00, 200.00, 'шт');


--
-- Data for Name: materials_products; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.materials_products VALUES (1, 1, 1, 0.1800);
INSERT INTO public.materials_products VALUES (2, 4, 1, 0.0090);
INSERT INTO public.materials_products VALUES (3, 9, 1, 1.0000);
INSERT INTO public.materials_products VALUES (4, 11, 1, 1.0000);
INSERT INTO public.materials_products VALUES (5, 17, 1, 1.0000);
INSERT INTO public.materials_products VALUES (6, 1, 2, 0.1800);
INSERT INTO public.materials_products VALUES (7, 5, 2, 0.0100);
INSERT INTO public.materials_products VALUES (8, 9, 2, 1.0000);
INSERT INTO public.materials_products VALUES (9, 11, 2, 1.0000);
INSERT INTO public.materials_products VALUES (10, 17, 2, 1.0000);
INSERT INTO public.materials_products VALUES (11, 2, 3, 0.1750);
INSERT INTO public.materials_products VALUES (12, 6, 3, 0.0110);
INSERT INTO public.materials_products VALUES (13, 9, 3, 1.0000);
INSERT INTO public.materials_products VALUES (14, 11, 3, 1.0000);
INSERT INTO public.materials_products VALUES (15, 17, 3, 1.0000);
INSERT INTO public.materials_products VALUES (16, 1, 4, 0.2200);
INSERT INTO public.materials_products VALUES (17, 5, 4, 0.0060);
INSERT INTO public.materials_products VALUES (18, 10, 4, 1.0000);
INSERT INTO public.materials_products VALUES (19, 12, 4, 1.0000);
INSERT INTO public.materials_products VALUES (20, 17, 4, 1.0000);
INSERT INTO public.materials_products VALUES (21, 1, 5, 0.2250);
INSERT INTO public.materials_products VALUES (22, 4, 5, 0.0040);
INSERT INTO public.materials_products VALUES (23, 10, 5, 1.0000);
INSERT INTO public.materials_products VALUES (24, 12, 5, 1.0000);
INSERT INTO public.materials_products VALUES (25, 17, 5, 1.0000);
INSERT INTO public.materials_products VALUES (26, 3, 6, 0.2800);
INSERT INTO public.materials_products VALUES (27, 7, 6, 0.0040);
INSERT INTO public.materials_products VALUES (28, 10, 6, 1.0000);
INSERT INTO public.materials_products VALUES (29, 13, 6, 1.0000);
INSERT INTO public.materials_products VALUES (30, 17, 6, 1.0000);
INSERT INTO public.materials_products VALUES (31, 3, 7, 0.2850);
INSERT INTO public.materials_products VALUES (32, 5, 7, 0.0030);
INSERT INTO public.materials_products VALUES (33, 10, 7, 1.0000);
INSERT INTO public.materials_products VALUES (34, 13, 7, 1.0000);
INSERT INTO public.materials_products VALUES (35, 17, 7, 1.0000);
INSERT INTO public.materials_products VALUES (36, 2, 8, 0.0700);
INSERT INTO public.materials_products VALUES (37, 4, 8, 0.0030);
INSERT INTO public.materials_products VALUES (38, 9, 8, 1.0000);
INSERT INTO public.materials_products VALUES (39, 11, 8, 1.0000);
INSERT INTO public.materials_products VALUES (40, 2, 9, 0.0700);
INSERT INTO public.materials_products VALUES (41, 5, 9, 0.0030);
INSERT INTO public.materials_products VALUES (42, 9, 9, 1.0000);
INSERT INTO public.materials_products VALUES (43, 11, 9, 1.0000);
INSERT INTO public.materials_products VALUES (44, 2, 10, 0.0700);
INSERT INTO public.materials_products VALUES (45, 6, 10, 0.0030);
INSERT INTO public.materials_products VALUES (46, 9, 10, 1.0000);
INSERT INTO public.materials_products VALUES (47, 11, 10, 1.0000);
INSERT INTO public.materials_products VALUES (48, 1, 11, 0.1800);
INSERT INTO public.materials_products VALUES (49, 4, 11, 0.0090);
INSERT INTO public.materials_products VALUES (50, 9, 11, 1.0000);
INSERT INTO public.materials_products VALUES (51, 11, 11, 1.0000);
INSERT INTO public.materials_products VALUES (52, 15, 11, 1.2000);
INSERT INTO public.materials_products VALUES (53, 16, 11, 1.0000);
INSERT INTO public.materials_products VALUES (54, 17, 11, 1.0000);
INSERT INTO public.materials_products VALUES (55, 1, 12, 0.2250);
INSERT INTO public.materials_products VALUES (56, 5, 12, 0.0100);
INSERT INTO public.materials_products VALUES (57, 10, 12, 1.0000);
INSERT INTO public.materials_products VALUES (58, 12, 12, 1.0000);
INSERT INTO public.materials_products VALUES (59, 15, 12, 1.5000);
INSERT INTO public.materials_products VALUES (60, 16, 12, 1.0000);
INSERT INTO public.materials_products VALUES (61, 17, 12, 1.0000);
INSERT INTO public.materials_products VALUES (62, 1, 13, 0.2100);
INSERT INTO public.materials_products VALUES (63, 5, 13, 0.0080);
INSERT INTO public.materials_products VALUES (64, 7, 13, 0.0030);
INSERT INTO public.materials_products VALUES (65, 10, 13, 1.0000);
INSERT INTO public.materials_products VALUES (66, 12, 13, 1.0000);
INSERT INTO public.materials_products VALUES (67, 18, 13, 1.0000);
INSERT INTO public.materials_products VALUES (68, 1, 14, 0.2150);
INSERT INTO public.materials_products VALUES (69, 6, 14, 0.0040);
INSERT INTO public.materials_products VALUES (70, 7, 14, 0.0020);
INSERT INTO public.materials_products VALUES (71, 10, 14, 1.0000);
INSERT INTO public.materials_products VALUES (72, 12, 14, 1.0000);
INSERT INTO public.materials_products VALUES (73, 18, 14, 1.0000);


--
-- Data for Name: product_types; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.product_types VALUES (1, 'Классическая свеча', 1.00);
INSERT INTO public.product_types VALUES (2, 'Премиальная свеча', 1.15);
INSERT INTO public.product_types VALUES (3, 'Подарочный набор', 1.30);
INSERT INTO public.product_types VALUES (4, 'Интерьерная свеча', 1.20);
INSERT INTO public.product_types VALUES (5, 'Мини-свеча', 0.80);
INSERT INTO public.product_types VALUES (6, 'Сезонная коллекция', 1.25);


--
-- Data for Name: products_import; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.products_import VALUES (1, 'AR-001', 1, 'Свеча «Лавандовый вечер»', 980.00);
INSERT INTO public.products_import VALUES (2, 'AR-002', 1, 'Свеча «Ванильный крем»', 950.00);
INSERT INTO public.products_import VALUES (3, 'AR-003', 1, 'Свеча «Цитрусовый сад»', 970.00);
INSERT INTO public.products_import VALUES (4, 'AR-004', 2, 'Свеча «Белый мускус»', 1290.00);
INSERT INTO public.products_import VALUES (5, 'AR-005', 2, 'Свеча «Дикая орхидея»', 1360.00);
INSERT INTO public.products_import VALUES (6, 'AR-006', 4, 'Свеча «Теплый янтарь»', 1490.00);
INSERT INTO public.products_import VALUES (7, 'AR-007', 4, 'Свеча «Сандаловое дерево»', 1520.00);
INSERT INTO public.products_import VALUES (8, 'AR-008', 5, 'Мини-свеча «Лаванда»', 420.00);
INSERT INTO public.products_import VALUES (9, 'AR-009', 5, 'Мини-свеча «Ваниль»', 410.00);
INSERT INTO public.products_import VALUES (10, 'AR-010', 5, 'Мини-свеча «Цитрус»', 430.00);
INSERT INTO public.products_import VALUES (11, 'AR-011', 3, 'Набор «Уютный дом»', 2190.00);
INSERT INTO public.products_import VALUES (12, 'AR-012', 3, 'Набор «Теплый вечер»', 2340.00);
INSERT INTO public.products_import VALUES (13, 'AR-013', 6, 'Свеча «Зимние специи»', 1380.00);
INSERT INTO public.products_import VALUES (14, 'AR-014', 6, 'Свеча «Хвойный лес»', 1410.00);


--
-- Name: material_types_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.material_types_id_seq', 10, true);


--
-- Name: materials_import_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.materials_import_id_seq', 18, true);


--
-- Name: materials_products_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.materials_products_id_seq', 73, true);


--
-- Name: product_types_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.product_types_id_seq', 6, true);


--
-- Name: products_import_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.products_import_id_seq', 14, true);


--
-- Name: material_types material_types_material_type_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.material_types
    ADD CONSTRAINT material_types_material_type_key UNIQUE (material_type);


--
-- Name: material_types material_types_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.material_types
    ADD CONSTRAINT material_types_pkey PRIMARY KEY (id);


--
-- Name: materials_import materials_import_material_name_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.materials_import
    ADD CONSTRAINT materials_import_material_name_key UNIQUE (material_name);


--
-- Name: materials_import materials_import_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.materials_import
    ADD CONSTRAINT materials_import_pkey PRIMARY KEY (id);


--
-- Name: materials_products materials_products_material_name_id_product_name_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.materials_products
    ADD CONSTRAINT materials_products_material_name_id_product_name_id_key UNIQUE (material_name_id, product_name_id);


--
-- Name: materials_products materials_products_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.materials_products
    ADD CONSTRAINT materials_products_pkey PRIMARY KEY (id);


--
-- Name: product_types product_types_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_types
    ADD CONSTRAINT product_types_pkey PRIMARY KEY (id);


--
-- Name: product_types product_types_product_type_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_types
    ADD CONSTRAINT product_types_product_type_key UNIQUE (product_type);


--
-- Name: products_import products_import_article_number_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.products_import
    ADD CONSTRAINT products_import_article_number_key UNIQUE (article_number);


--
-- Name: products_import products_import_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.products_import
    ADD CONSTRAINT products_import_pkey PRIMARY KEY (id);


--
-- Name: materials_import materials_import_material_type_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.materials_import
    ADD CONSTRAINT materials_import_material_type_id_fkey FOREIGN KEY (material_type_id) REFERENCES public.material_types(id);


--
-- Name: materials_products materials_products_material_name_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.materials_products
    ADD CONSTRAINT materials_products_material_name_id_fkey FOREIGN KEY (material_name_id) REFERENCES public.materials_import(id);


--
-- Name: materials_products materials_products_product_name_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.materials_products
    ADD CONSTRAINT materials_products_product_name_id_fkey FOREIGN KEY (product_name_id) REFERENCES public.products_import(id);


--
-- Name: products_import products_import_product_type_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.products_import
    ADD CONSTRAINT products_import_product_type_id_fkey FOREIGN KEY (product_type_id) REFERENCES public.product_types(id);


--
-- PostgreSQL database dump complete
--

SET search_path TO public;

