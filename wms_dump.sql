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
-- Name: barcodes; Type: TABLE; Schema: public; Owner: wmsuser
--

CREATE TABLE public.barcodes (
    parent_id integer,
    name character varying(128) DEFAULT ''::character varying NOT NULL,
    barcode_type integer NOT NULL,
    id integer NOT NULL
);


ALTER TABLE public.barcodes OWNER TO wmsuser;

--
-- Name: barcodes_id_seq; Type: SEQUENCE; Schema: public; Owner: wmsuser
--

CREATE SEQUENCE public.barcodes_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.barcodes_id_seq OWNER TO wmsuser;

--
-- Name: barcodes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: wmsuser
--

ALTER SEQUENCE public.barcodes_id_seq OWNED BY public.barcodes.id;


--
-- Name: cells; Type: TABLE; Schema: public; Owner: wmsuser
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


ALTER TABLE public.cells OWNER TO wmsuser;

--
-- Name: cells_id_seq; Type: SEQUENCE; Schema: public; Owner: wmsuser
--

CREATE SEQUENCE public.cells_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.cells_id_seq OWNER TO wmsuser;

--
-- Name: cells_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: wmsuser
--

ALTER SEQUENCE public.cells_id_seq OWNED BY public.cells.id;


--
-- Name: manufacturers; Type: TABLE; Schema: public; Owner: wmsuser
--

CREATE TABLE public.manufacturers (
    id integer NOT NULL,
    name character varying(255) DEFAULT ''::character varying NOT NULL
);


ALTER TABLE public.manufacturers OWNER TO wmsuser;

--
-- Name: manufacturers_id_seq; Type: SEQUENCE; Schema: public; Owner: wmsuser
--

CREATE SEQUENCE public.manufacturers_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.manufacturers_id_seq OWNER TO wmsuser;

--
-- Name: manufacturers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: wmsuser
--

ALTER SEQUENCE public.manufacturers_id_seq OWNED BY public.manufacturers.id;


--
-- Name: products; Type: TABLE; Schema: public; Owner: wmsuser
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


ALTER TABLE public.products OWNER TO wmsuser;

--
-- Name: products_id_seq; Type: SEQUENCE; Schema: public; Owner: wmsuser
--

CREATE SEQUENCE public.products_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.products_id_seq OWNER TO wmsuser;

--
-- Name: products_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: wmsuser
--

ALTER SEQUENCE public.products_id_seq OWNED BY public.products.id;


--
-- Name: receipt_headers; Type: TABLE; Schema: public; Owner: wmsuser
--

CREATE TABLE public.receipt_headers (
    id integer NOT NULL,
    number character varying(20) DEFAULT ''::character varying NOT NULL,
    date timestamp with time zone DEFAULT '1970-01-01 00:00:00+00'::timestamp with time zone NOT NULL,
    doc_type integer DEFAULT 0 NOT NULL
);


ALTER TABLE public.receipt_headers OWNER TO wmsuser;

--
-- Name: receipt_docs_id_seq; Type: SEQUENCE; Schema: public; Owner: wmsuser
--

CREATE SEQUENCE public.receipt_docs_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.receipt_docs_id_seq OWNER TO wmsuser;

--
-- Name: receipt_docs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: wmsuser
--

ALTER SEQUENCE public.receipt_docs_id_seq OWNED BY public.receipt_headers.id;


--
-- Name: receipt_items; Type: TABLE; Schema: public; Owner: wmsuser
--

CREATE TABLE public.receipt_items (
    parent_id integer DEFAULT 0 NOT NULL,
    row_id character varying(36) DEFAULT ''::character varying NOT NULL,
    product_id integer DEFAULT 0 NOT NULL,
    quantity integer DEFAULT 0 NOT NULL
);


ALTER TABLE public.receipt_items OWNER TO wmsuser;

--
-- Name: storage1; Type: TABLE; Schema: public; Owner: wmsuser
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


ALTER TABLE public.storage1 OWNER TO wmsuser;

--
-- Name: user_printers; Type: TABLE; Schema: public; Owner: wmsuser
--

CREATE TABLE public.user_printers (
    user_id integer DEFAULT 0 NOT NULL,
    printer_name character varying(100) DEFAULT ''::character varying NOT NULL,
    printer_instance character varying(100) DEFAULT ''::character varying NOT NULL,
    printer_type integer DEFAULT 0 NOT NULL
);


ALTER TABLE public.user_printers OWNER TO wmsuser;

--
-- Name: users; Type: TABLE; Schema: public; Owner: wmsuser
--

CREATE TABLE public.users (
    id integer NOT NULL,
    name character varying(150) DEFAULT ''::character varying NOT NULL
);


ALTER TABLE public.users OWNER TO wmsuser;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: wmsuser
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO wmsuser;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: wmsuser
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: whs; Type: TABLE; Schema: public; Owner: wmsuser
--

CREATE TABLE public.whs (
    id integer NOT NULL,
    name character varying(50) DEFAULT ''::character varying,
    address character varying(255) DEFAULT ''::character varying NOT NULL
);


ALTER TABLE public.whs OWNER TO wmsuser;

--
-- Name: whs_id_seq; Type: SEQUENCE; Schema: public; Owner: wmsuser
--

CREATE SEQUENCE public.whs_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.whs_id_seq OWNER TO wmsuser;

--
-- Name: whs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: wmsuser
--

ALTER SEQUENCE public.whs_id_seq OWNED BY public.whs.id;


--
-- Name: zones; Type: TABLE; Schema: public; Owner: wmsuser
--

CREATE TABLE public.zones (
    id integer NOT NULL,
    name character varying(50) DEFAULT ''::character varying,
    parent_id integer,
    zone_type smallint
);


ALTER TABLE public.zones OWNER TO wmsuser;

--
-- Name: zones_id_seq; Type: SEQUENCE; Schema: public; Owner: wmsuser
--

CREATE SEQUENCE public.zones_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.zones_id_seq OWNER TO wmsuser;

--
-- Name: zones_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: wmsuser
--

ALTER SEQUENCE public.zones_id_seq OWNED BY public.zones.id;


--
-- Name: barcodes id; Type: DEFAULT; Schema: public; Owner: wmsuser
--

ALTER TABLE ONLY public.barcodes ALTER COLUMN id SET DEFAULT nextval('public.barcodes_id_seq'::regclass);


--
-- Name: cells id; Type: DEFAULT; Schema: public; Owner: wmsuser
--

ALTER TABLE ONLY public.cells ALTER COLUMN id SET DEFAULT nextval('public.cells_id_seq'::regclass);


--
-- Name: manufacturers id; Type: DEFAULT; Schema: public; Owner: wmsuser
--

ALTER TABLE ONLY public.manufacturers ALTER COLUMN id SET DEFAULT nextval('public.manufacturers_id_seq'::regclass);


--
-- Name: products id; Type: DEFAULT; Schema: public; Owner: wmsuser
--

ALTER TABLE ONLY public.products ALTER COLUMN id SET DEFAULT nextval('public.products_id_seq'::regclass);


--
-- Name: receipt_headers id; Type: DEFAULT; Schema: public; Owner: wmsuser
--

ALTER TABLE ONLY public.receipt_headers ALTER COLUMN id SET DEFAULT nextval('public.receipt_docs_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: wmsuser
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: whs id; Type: DEFAULT; Schema: public; Owner: wmsuser
--

ALTER TABLE ONLY public.whs ALTER COLUMN id SET DEFAULT nextval('public.whs_id_seq'::regclass);


--
-- Name: zones id; Type: DEFAULT; Schema: public; Owner: wmsuser
--

ALTER TABLE ONLY public.zones ALTER COLUMN id SET DEFAULT nextval('public.zones_id_seq'::regclass);


--
-- Data for Name: barcodes; Type: TABLE DATA; Schema: public; Owner: wmsuser
--

COPY public.barcodes (parent_id, name, barcode_type, id) FROM stdin;
\.


--
-- Data for Name: cells; Type: TABLE DATA; Schema: public; Owner: wmsuser
--

COPY public.cells (id, name, whs_id, zone_id, passage_id, rack_id, floor, sz_length, sz_width, sz_height, sz_volume, sz_uf_volume, sz_weight, is_size_free, is_weight_free, not_allowed_in, not_allowed_out, is_service) FROM stdin;
\.


--
-- Data for Name: manufacturers; Type: TABLE DATA; Schema: public; Owner: wmsuser
--

COPY public.manufacturers (id, name) FROM stdin;
\.


--
-- Data for Name: products; Type: TABLE DATA; Schema: public; Owner: wmsuser
--

COPY public.products (id, name, manufacturer_id, sz_length, sz_wight, sz_height, sz_weight, sz_volume, sz_uf_volume, item_number) FROM stdin;
\.


--
-- Data for Name: receipt_headers; Type: TABLE DATA; Schema: public; Owner: wmsuser
--

COPY public.receipt_headers (id, number, date, doc_type) FROM stdin;
\.


--
-- Data for Name: receipt_items; Type: TABLE DATA; Schema: public; Owner: wmsuser
--

COPY public.receipt_items (parent_id, row_id, product_id, quantity) FROM stdin;
\.


--
-- Data for Name: storage1; Type: TABLE DATA; Schema: public; Owner: wmsuser
--

COPY public.storage1 (zone_id, cell_id, prod_id, quantity, doc_id, doc_type, row_id) FROM stdin;
\.


--
-- Data for Name: user_printers; Type: TABLE DATA; Schema: public; Owner: wmsuser
--

COPY public.user_printers (user_id, printer_name, printer_instance, printer_type) FROM stdin;
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: wmsuser
--

COPY public.users (id, name) FROM stdin;
\.


--
-- Data for Name: whs; Type: TABLE DATA; Schema: public; Owner: wmsuser
--

COPY public.whs (id, name, address) FROM stdin;
1	Мой склад	
\.


--
-- Data for Name: zones; Type: TABLE DATA; Schema: public; Owner: wmsuser
--

COPY public.zones (id, name, parent_id, zone_type) FROM stdin;
4	Зона приемки	1	1
5	Зона отгрузки	1	2
6	Зона хранения	1	0
\.


--
-- Name: barcodes_id_seq; Type: SEQUENCE SET; Schema: public; Owner: wmsuser
--

SELECT pg_catalog.setval('public.barcodes_id_seq', 1, true);


--
-- Name: cells_id_seq; Type: SEQUENCE SET; Schema: public; Owner: wmsuser
--

SELECT pg_catalog.setval('public.cells_id_seq', 1, true);


--
-- Name: manufacturers_id_seq; Type: SEQUENCE SET; Schema: public; Owner: wmsuser
--

SELECT pg_catalog.setval('public.manufacturers_id_seq', 1, true);


--
-- Name: products_id_seq; Type: SEQUENCE SET; Schema: public; Owner: wmsuser
--

SELECT pg_catalog.setval('public.products_id_seq', 1, true);


--
-- Name: receipt_docs_id_seq; Type: SEQUENCE SET; Schema: public; Owner: wmsuser
--

SELECT pg_catalog.setval('public.receipt_docs_id_seq', 1, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: wmsuser
--

SELECT pg_catalog.setval('public.users_id_seq', 1, true);


--
-- Name: whs_id_seq; Type: SEQUENCE SET; Schema: public; Owner: wmsuser
--

SELECT pg_catalog.setval('public.whs_id_seq', 10, true);


--
-- Name: zones_id_seq; Type: SEQUENCE SET; Schema: public; Owner: wmsuser
--

SELECT pg_catalog.setval('public.zones_id_seq', 15, true);


--
-- Name: barcodes barcodes_pk; Type: CONSTRAINT; Schema: public; Owner: wmsuser
--

ALTER TABLE ONLY public.barcodes
    ADD CONSTRAINT barcodes_pk PRIMARY KEY (id);


--
-- Name: cells cells_pk; Type: CONSTRAINT; Schema: public; Owner: wmsuser
--

ALTER TABLE ONLY public.cells
    ADD CONSTRAINT cells_pk PRIMARY KEY (id);


--
-- Name: manufacturers manufacturers_pk; Type: CONSTRAINT; Schema: public; Owner: wmsuser
--

ALTER TABLE ONLY public.manufacturers
    ADD CONSTRAINT manufacturers_pk PRIMARY KEY (id);


--
-- Name: products products_pk; Type: CONSTRAINT; Schema: public; Owner: wmsuser
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_pk PRIMARY KEY (id);


--
-- Name: receipt_headers receipt_docs_pk; Type: CONSTRAINT; Schema: public; Owner: wmsuser
--

ALTER TABLE ONLY public.receipt_headers
    ADD CONSTRAINT receipt_docs_pk PRIMARY KEY (id);


--
-- Name: receipt_items receipt_items_pk; Type: CONSTRAINT; Schema: public; Owner: wmsuser
--

ALTER TABLE ONLY public.receipt_items
    ADD CONSTRAINT receipt_items_pk PRIMARY KEY (row_id);


--
-- Name: users users_pk; Type: CONSTRAINT; Schema: public; Owner: wmsuser
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pk PRIMARY KEY (id);


--
-- Name: whs whs_pk; Type: CONSTRAINT; Schema: public; Owner: wmsuser
--

ALTER TABLE ONLY public.whs
    ADD CONSTRAINT whs_pk PRIMARY KEY (id);


--
-- Name: zones zones_pk; Type: CONSTRAINT; Schema: public; Owner: wmsuser
--

ALTER TABLE ONLY public.zones
    ADD CONSTRAINT zones_pk PRIMARY KEY (id);


--
-- Name: barcodes_product_id_barcode_barcode_type_uindex; Type: INDEX; Schema: public; Owner: wmsuser
--

CREATE UNIQUE INDEX barcodes_product_id_barcode_barcode_type_uindex ON public.barcodes USING btree (parent_id, name, barcode_type);


--
-- Name: cells_id_uindex; Type: INDEX; Schema: public; Owner: wmsuser
--

CREATE UNIQUE INDEX cells_id_uindex ON public.cells USING btree (id);


--
-- Name: cells_whs_id_zone_id_passage_id_floor_uindex; Type: INDEX; Schema: public; Owner: wmsuser
--

CREATE UNIQUE INDEX cells_whs_id_zone_id_passage_id_floor_uindex ON public.cells USING btree (whs_id, zone_id, passage_id, floor);


--
-- Name: receipt_items_row_id_uindex; Type: INDEX; Schema: public; Owner: wmsuser
--

CREATE UNIQUE INDEX receipt_items_row_id_uindex ON public.receipt_items USING btree (row_id);


--
-- Name: products products_manufacturers_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: wmsuser
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_manufacturers_id_fk FOREIGN KEY (manufacturer_id) REFERENCES public.manufacturers(id);


--
-- Name: receipt_items receipt_items_products_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: wmsuser
--

ALTER TABLE ONLY public.receipt_items
    ADD CONSTRAINT receipt_items_products_id_fk FOREIGN KEY (product_id) REFERENCES public.products(id);


--
-- Name: storage1 storage1_cells_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: wmsuser
--

ALTER TABLE ONLY public.storage1
    ADD CONSTRAINT storage1_cells_id_fk FOREIGN KEY (cell_id) REFERENCES public.cells(id);


--
-- PostgreSQL database dump complete
--

