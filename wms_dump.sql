--
-- PostgreSQL database dump
--

-- Dumped from database version 13.1
-- Dumped by pg_dump version 13.5 (Ubuntu 13.5-2.pgdg20.04+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
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
-- Name: barcodes; Type: TABLE; Schema: public; Owner: devuser
--

CREATE TABLE public.barcodes (
    parent_id integer,
    name character varying(128) DEFAULT ''::character varying NOT NULL,
    barcode_type integer NOT NULL,
    id integer NOT NULL
);


ALTER TABLE public.barcodes OWNER TO devuser;

--
-- Name: barcodes_id_seq; Type: SEQUENCE; Schema: public; Owner: devuser
--

CREATE SEQUENCE public.barcodes_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.barcodes_id_seq OWNER TO devuser;

--
-- Name: barcodes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: devuser
--

ALTER SEQUENCE public.barcodes_id_seq OWNED BY public.barcodes.id;


--
-- Name: cells; Type: TABLE; Schema: public; Owner: devuser
--

CREATE TABLE public.cells (
    id integer NOT NULL,
    name character varying(50) DEFAULT ''::character varying NOT NULL,
    whs_id integer DEFAULT 0 NOT NULL,
    zone_id integer DEFAULT 0 NOT NULL,
    passage_id integer DEFAULT 0 NOT NULL,
    rack_id integer DEFAULT 0 NOT NULL,
    floor integer DEFAULT 0 NOT NULL,
    sz_length integer DEFAULT 0 NOT NULL,
    sz_width integer DEFAULT 0 NOT NULL,
    sz_height integer DEFAULT 0 NOT NULL,
    sz_volume numeric(8,3) DEFAULT 0 NOT NULL,
    sz_uf_volume numeric(8,3) DEFAULT 0 NOT NULL,
    sz_weight numeric(8,3) DEFAULT 0 NOT NULL,
    is_size_free boolean DEFAULT false NOT NULL,
    is_weight_free boolean DEFAULT false NOT NULL,
    not_allowed_in boolean DEFAULT false NOT NULL,
    not_allowed_out boolean DEFAULT false NOT NULL,
    is_service boolean DEFAULT false NOT NULL
);


ALTER TABLE public.cells OWNER TO devuser;

--
-- Name: cells_id_seq; Type: SEQUENCE; Schema: public; Owner: devuser
--

CREATE SEQUENCE public.cells_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.cells_id_seq OWNER TO devuser;

--
-- Name: cells_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: devuser
--

ALTER SEQUENCE public.cells_id_seq OWNED BY public.cells.id;


--
-- Name: manufacturers; Type: TABLE; Schema: public; Owner: devuser
--

CREATE TABLE public.manufacturers (
    id integer NOT NULL,
    name character varying(255) DEFAULT ''::character varying NOT NULL
);


ALTER TABLE public.manufacturers OWNER TO devuser;

--
-- Name: manufacturers_id_seq; Type: SEQUENCE; Schema: public; Owner: devuser
--

CREATE SEQUENCE public.manufacturers_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.manufacturers_id_seq OWNER TO devuser;

--
-- Name: manufacturers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: devuser
--

ALTER SEQUENCE public.manufacturers_id_seq OWNED BY public.manufacturers.id;


--
-- Name: products; Type: TABLE; Schema: public; Owner: devuser
--

CREATE TABLE public.products (
    id integer NOT NULL,
    name character varying(255) DEFAULT ''::character varying NOT NULL,
    manufacturer_id integer DEFAULT 0,
    sz_length integer DEFAULT 0 NOT NULL,
    sz_wight integer DEFAULT 0 NOT NULL,
    sz_height integer DEFAULT 0 NOT NULL,
    sz_weight numeric(8,3) DEFAULT 0 NOT NULL,
    sz_volume numeric(8,3) DEFAULT 0 NOT NULL,
    sz_uf_volume numeric(8,3) DEFAULT 0 NOT NULL,
    item_number character varying(50) DEFAULT ''::character varying NOT NULL
);


ALTER TABLE public.products OWNER TO devuser;

--
-- Name: products_id_seq; Type: SEQUENCE; Schema: public; Owner: devuser
--

CREATE SEQUENCE public.products_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.products_id_seq OWNER TO devuser;

--
-- Name: products_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: devuser
--

ALTER SEQUENCE public.products_id_seq OWNED BY public.products.id;


--
-- Name: receipt_headers; Type: TABLE; Schema: public; Owner: devuser
--

CREATE TABLE public.receipt_headers (
    id integer NOT NULL,
    number character varying(20) DEFAULT ''::character varying NOT NULL,
    date timestamp with time zone DEFAULT '1970-01-01 00:00:00+00'::timestamp with time zone NOT NULL,
    doc_type integer DEFAULT 0 NOT NULL
);


ALTER TABLE public.receipt_headers OWNER TO devuser;

--
-- Name: receipt_docs_id_seq; Type: SEQUENCE; Schema: public; Owner: devuser
--

CREATE SEQUENCE public.receipt_docs_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.receipt_docs_id_seq OWNER TO devuser;

--
-- Name: receipt_docs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: devuser
--

ALTER SEQUENCE public.receipt_docs_id_seq OWNED BY public.receipt_headers.id;


--
-- Name: receipt_items; Type: TABLE; Schema: public; Owner: devuser
--

CREATE TABLE public.receipt_items (
    parent_id integer DEFAULT 0 NOT NULL,
    row_id character varying(36) DEFAULT ''::character varying NOT NULL,
    product_id integer DEFAULT 0 NOT NULL,
    quantity integer DEFAULT 0 NOT NULL
);


ALTER TABLE public.receipt_items OWNER TO devuser;

--
-- Name: storage1; Type: TABLE; Schema: public; Owner: devuser
--

CREATE TABLE public.storage1 (
    zone_id integer,
    cell_id integer,
    prod_id integer,
    quantity integer,
    doc_id integer DEFAULT 0 NOT NULL,
    doc_type smallint DEFAULT 0 NOT NULL,
    row_id character varying(36) DEFAULT ''::character varying NOT NULL
);


ALTER TABLE public.storage1 OWNER TO devuser;

--
-- Name: storage10; Type: TABLE; Schema: public; Owner: devuser
--

CREATE TABLE public.storage10 (
    doc_id integer DEFAULT 0 NOT NULL,
    doc_type smallint DEFAULT 0 NOT NULL,
    row_id character varying(36) DEFAULT ''::character varying NOT NULL,
    zone_id integer,
    cell_id integer,
    prod_id integer,
    quantity integer
);


ALTER TABLE public.storage10 OWNER TO devuser;

--
-- Name: user_printers; Type: TABLE; Schema: public; Owner: devuser
--

CREATE TABLE public.user_printers (
    user_id integer DEFAULT 0 NOT NULL,
    printer_name character varying(100) DEFAULT ''::character varying NOT NULL,
    printer_instance character varying(100) DEFAULT ''::character varying NOT NULL,
    printer_type integer DEFAULT 0 NOT NULL
);


ALTER TABLE public.user_printers OWNER TO devuser;

--
-- Name: users; Type: TABLE; Schema: public; Owner: devuser
--

CREATE TABLE public.users (
    id integer NOT NULL,
    name character varying(150) DEFAULT ''::character varying NOT NULL
);


ALTER TABLE public.users OWNER TO devuser;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: devuser
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO devuser;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: devuser
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: whs; Type: TABLE; Schema: public; Owner: devuser
--

CREATE TABLE public.whs (
    id integer NOT NULL,
    name character varying(50) DEFAULT ''::character varying,
    address character varying(255) DEFAULT ''::character varying NOT NULL
);


ALTER TABLE public.whs OWNER TO devuser;

--
-- Name: whs_id_seq; Type: SEQUENCE; Schema: public; Owner: devuser
--

CREATE SEQUENCE public.whs_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.whs_id_seq OWNER TO devuser;

--
-- Name: whs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: devuser
--

ALTER SEQUENCE public.whs_id_seq OWNED BY public.whs.id;


--
-- Name: zones; Type: TABLE; Schema: public; Owner: devuser
--

CREATE TABLE public.zones (
    id integer NOT NULL,
    name character varying(50) DEFAULT ''::character varying,
    parent_id integer,
    zone_type smallint
);


ALTER TABLE public.zones OWNER TO devuser;

--
-- Name: zones_id_seq; Type: SEQUENCE; Schema: public; Owner: devuser
--

CREATE SEQUENCE public.zones_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.zones_id_seq OWNER TO devuser;

--
-- Name: zones_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: devuser
--

ALTER SEQUENCE public.zones_id_seq OWNED BY public.zones.id;


--
-- Name: barcodes id; Type: DEFAULT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.barcodes ALTER COLUMN id SET DEFAULT nextval('public.barcodes_id_seq'::regclass);


--
-- Name: cells id; Type: DEFAULT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.cells ALTER COLUMN id SET DEFAULT nextval('public.cells_id_seq'::regclass);


--
-- Name: manufacturers id; Type: DEFAULT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.manufacturers ALTER COLUMN id SET DEFAULT nextval('public.manufacturers_id_seq'::regclass);


--
-- Name: products id; Type: DEFAULT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.products ALTER COLUMN id SET DEFAULT nextval('public.products_id_seq'::regclass);


--
-- Name: receipt_headers id; Type: DEFAULT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.receipt_headers ALTER COLUMN id SET DEFAULT nextval('public.receipt_docs_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: whs id; Type: DEFAULT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.whs ALTER COLUMN id SET DEFAULT nextval('public.whs_id_seq'::regclass);


--
-- Name: zones id; Type: DEFAULT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.zones ALTER COLUMN id SET DEFAULT nextval('public.zones_id_seq'::regclass);


--
-- Data for Name: barcodes; Type: TABLE DATA; Schema: public; Owner: devuser
--

COPY public.barcodes (parent_id, name, barcode_type, id) FROM stdin;
18	110 123456sw78901	3	1
18	110 0045f36782901	2	2
11	123456sw78901	3	3
11	0045f36782901	2	4
9	12345678901	3	5
9	0045678901	2	6
17	110 123456sw78901	3	15
17	110 0045f36782901	2	16
10	123456s78901	3	17
10	0045f3678901	2	18
2	1234567890	3	19
20	123456sw78901	3	20
20	0045f36782901	2	21
13	123456sw78901	3	26
13	0045f36782901	2	27
15	12345611278901	3	28
15	0045f36112782901	2	29
22	sdfgsdfgew5435645	0	34
16	110 123456sw78901	3	35
16	110 0045f36782901	2	36
16	110 004	0	37
23	345234523452345	0	39
23	34563456	0	40
25	2543452345	0	41
27	52345	0	49
27	3425345234	0	50
41	345234523452345	0	51
43	34563456345654	0	52
60		0	53
14	123456sw78901	3	54
14	0045f36782901	2	55
14	45235423452345	0	56
\.


--
-- Data for Name: cells; Type: TABLE DATA; Schema: public; Owner: devuser
--

COPY public.cells (id, name, whs_id, zone_id, passage_id, rack_id, floor, sz_length, sz_width, sz_height, sz_volume, sz_uf_volume, sz_weight, is_size_free, is_weight_free, not_allowed_in, not_allowed_out, is_service) FROM stdin;
1	test 1	1	1	2	0	3	0	0	0	0.000	0.000	0.000	f	f	f	f	f
2	1-1-1-0-1	99	1	1	0	1	0	0	0	0.000	0.000	0.000	f	f	f	f	f
\.


--
-- Data for Name: manufacturers; Type: TABLE DATA; Schema: public; Owner: devuser
--

COPY public.manufacturers (id, name) FROM stdin;
63	Аквафор
64	PROFILUX
65	
66	Новый производитель тестового товара
69	какой-то производитель
70	Новый производитель2
71	DEKO2
72	PROFILUXX
13	Shenzhen Xunlong Software
14	OUIO
15	SanDisk
16	FNIRSI
19	РОСТЕРМ
20	UNIEL
21	СВАРИС
12	DEKO
\.


--
-- Data for Name: products; Type: TABLE DATA; Schema: public; Owner: devuser
--

COPY public.products (id, name, manufacturer_id, sz_length, sz_wight, sz_height, sz_weight, sz_volume, sz_uf_volume, item_number) FROM stdin;
44	Набор картриджей Аквафор К5-К4-К8 Универсал Н для жесткой и железистой воды	63	0	0	0	0.000	0.000	0.000	
52	Покрытие декоративное Profilux 7 кг цвет белый	64	0	0	0	0.000	0.000	0.000	
14	Осциллограф FNIRSI-1013D	16	0	0	0	0.000	0.000	0.000	1013D
63	Набор отверток DKMT65 065-0223	71	0	0	0	0.000	0.000	0.000	
61	Тестовый товар	69	0	0	0	0.000	0.000	0.000	2342567
64	Тестовый товар2	72	0	0	0	0.000	0.000	0.000	987654321
65	Тестовый товар2	64	0	0	0	0.000	0.000	0.000	
9	Набор отверток DKMT65 065-0223	12	0	0	0	0.000	0.000	0.000	
19	Редуктор давления РДСГ 1-1.2	19	0	0	0	0.000	0.000	0.000	
2	Патрон керамический Uniel GU4/GU5.3	20	0	0	0	0.000	0.000	0.000	
20	Припой Сварис ПОС-61, D1 мм, катушка, с канифолью 50 г	21	0	0	0	0.000	0.000	0.000	
13	Карта памяти Micro SD 64Gb class 10	15	0	0	0	0.000	0.000	0.000	64G10
15	Карта памяти Micro SD 32Gb class 10	14	0	0	0	0.000	0.000	0.000	32G10
16	Orange Pi Zero 512MB H3	13	0	0	0	0.000	0.000	0.000	
\.


--
-- Data for Name: receipt_headers; Type: TABLE DATA; Schema: public; Owner: devuser
--

COPY public.receipt_headers (id, number, date, doc_type) FROM stdin;
23		2023-02-20 00:00:00+00	2
24		2023-02-20 00:00:00+00	2
25		2023-02-21 00:00:00+00	2
26		2023-02-21 00:00:00+00	2
\.


--
-- Data for Name: receipt_items; Type: TABLE DATA; Schema: public; Owner: devuser
--

COPY public.receipt_items (parent_id, row_id, product_id, quantity) FROM stdin;
23	23.1	63	7
24	24.1	63	3
25	25.1	65	3
26	26.1	61	5
\.


--
-- Data for Name: storage1; Type: TABLE DATA; Schema: public; Owner: devuser
--

COPY public.storage1 (zone_id, cell_id, prod_id, quantity, doc_id, doc_type, row_id) FROM stdin;
0	2	32	100	0	0	
0	2	34	40	0	0	
0	2	34	40	0	0	
0	2	34	40	0	0	
0	2	34	40	0	0	
0	2	34	40	0	0	
0	2	32	100	0	0	
0	2	32	-180	0	0	
0	2	32	100	0	0	
0	2	34	40	0	0	
0	2	34	40	0	0	
0	2	34	40	0	0	
0	2	34	40	0	0	
0	2	34	40	0	0	
0	2	32	100	0	0	
0	2	32	-180	0	0	
1	2	57	6	20	1	20.1
1	2	58	10	21	1	21.1
1	2	59	15	22	1	22.1
1	2	63	7	23	1	23.1
1	2	63	3	24	1	24.1
1	2	65	3	25	1	25.1
1	2	61	5	26	1	26.1
\.


--
-- Data for Name: storage10; Type: TABLE DATA; Schema: public; Owner: devuser
--

COPY public.storage10 (doc_id, doc_type, row_id, zone_id, cell_id, prod_id, quantity) FROM stdin;
\.


--
-- Data for Name: user_printers; Type: TABLE DATA; Schema: public; Owner: devuser
--

COPY public.user_printers (user_id, printer_name, printer_instance, printer_type) FROM stdin;
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: devuser
--

COPY public.users (id, name) FROM stdin;
6	Mike
7	Михаил
3	Админ
10	Новый юзер
\.


--
-- Data for Name: whs; Type: TABLE DATA; Schema: public; Owner: devuser
--

COPY public.whs (id, name, address) FROM stdin;
1	Мой склад 1	
10	Мой склад 2	
\.


--
-- Data for Name: zones; Type: TABLE DATA; Schema: public; Owner: devuser
--

COPY public.zones (id, name, parent_id, zone_type) FROM stdin;
1	Зона приемки	6	1
2	Зона отгрузки	6	2
3	Зона хранения	6	0
4	Зона приемки	1	1
5	Зона отгрузки	1	2
6	Зона хранения	1	0
13	Зона приемки	10	1
14	Зона отгрузки	10	2
15	Зона хранения	10	0
\.


--
-- Name: barcodes_id_seq; Type: SEQUENCE SET; Schema: public; Owner: devuser
--

SELECT pg_catalog.setval('public.barcodes_id_seq', 56, true);


--
-- Name: cells_id_seq; Type: SEQUENCE SET; Schema: public; Owner: devuser
--

SELECT pg_catalog.setval('public.cells_id_seq', 2, true);


--
-- Name: manufacturers_id_seq; Type: SEQUENCE SET; Schema: public; Owner: devuser
--

SELECT pg_catalog.setval('public.manufacturers_id_seq', 72, true);


--
-- Name: products_id_seq; Type: SEQUENCE SET; Schema: public; Owner: devuser
--

SELECT pg_catalog.setval('public.products_id_seq', 65, true);


--
-- Name: receipt_docs_id_seq; Type: SEQUENCE SET; Schema: public; Owner: devuser
--

SELECT pg_catalog.setval('public.receipt_docs_id_seq', 35, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: devuser
--

SELECT pg_catalog.setval('public.users_id_seq', 10, true);


--
-- Name: whs_id_seq; Type: SEQUENCE SET; Schema: public; Owner: devuser
--

SELECT pg_catalog.setval('public.whs_id_seq', 10, true);


--
-- Name: zones_id_seq; Type: SEQUENCE SET; Schema: public; Owner: devuser
--

SELECT pg_catalog.setval('public.zones_id_seq', 15, true);


--
-- Name: barcodes barcodes_pk; Type: CONSTRAINT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.barcodes
    ADD CONSTRAINT barcodes_pk PRIMARY KEY (id);


--
-- Name: cells cells_pk; Type: CONSTRAINT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.cells
    ADD CONSTRAINT cells_pk PRIMARY KEY (id);


--
-- Name: manufacturers manufacturers_pk; Type: CONSTRAINT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.manufacturers
    ADD CONSTRAINT manufacturers_pk PRIMARY KEY (id);


--
-- Name: products products_pk; Type: CONSTRAINT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_pk PRIMARY KEY (id);


--
-- Name: receipt_headers receipt_docs_pk; Type: CONSTRAINT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.receipt_headers
    ADD CONSTRAINT receipt_docs_pk PRIMARY KEY (id);


--
-- Name: receipt_items receipt_items_pk; Type: CONSTRAINT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.receipt_items
    ADD CONSTRAINT receipt_items_pk PRIMARY KEY (row_id);


--
-- Name: users users_pk; Type: CONSTRAINT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pk PRIMARY KEY (id);


--
-- Name: whs whs_pk; Type: CONSTRAINT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.whs
    ADD CONSTRAINT whs_pk PRIMARY KEY (id);


--
-- Name: zones zones_pk; Type: CONSTRAINT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.zones
    ADD CONSTRAINT zones_pk PRIMARY KEY (id);


--
-- Name: barcodes_product_id_barcode_barcode_type_uindex; Type: INDEX; Schema: public; Owner: devuser
--

CREATE UNIQUE INDEX barcodes_product_id_barcode_barcode_type_uindex ON public.barcodes USING btree (parent_id, name, barcode_type);


--
-- Name: cells_id_uindex; Type: INDEX; Schema: public; Owner: devuser
--

CREATE UNIQUE INDEX cells_id_uindex ON public.cells USING btree (id);


--
-- Name: cells_whs_id_zone_id_passage_id_floor_uindex; Type: INDEX; Schema: public; Owner: devuser
--

CREATE UNIQUE INDEX cells_whs_id_zone_id_passage_id_floor_uindex ON public.cells USING btree (whs_id, zone_id, passage_id, floor);


--
-- Name: receipt_items_row_id_uindex; Type: INDEX; Schema: public; Owner: devuser
--

CREATE UNIQUE INDEX receipt_items_row_id_uindex ON public.receipt_items USING btree (row_id);


--
-- Name: products products_manufacturers_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_manufacturers_id_fk FOREIGN KEY (manufacturer_id) REFERENCES public.manufacturers(id);


--
-- Name: receipt_items receipt_items_products_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.receipt_items
    ADD CONSTRAINT receipt_items_products_id_fk FOREIGN KEY (product_id) REFERENCES public.products(id);


--
-- Name: storage10 storage10_cells_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.storage10
    ADD CONSTRAINT storage10_cells_id_fk FOREIGN KEY (cell_id) REFERENCES public.cells(id);


--
-- Name: storage1 storage1_cells_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.storage1
    ADD CONSTRAINT storage1_cells_id_fk FOREIGN KEY (cell_id) REFERENCES public.cells(id);


--
-- PostgreSQL database dump complete
--

