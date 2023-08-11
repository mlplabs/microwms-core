PGDMP                          {            wmsdb %   12.13 (Ubuntu 12.13-0ubuntu0.20.04.1)     13.5 (Ubuntu 13.5-2.pgdg20.04+1) J    '           0    0    ENCODING    ENCODING        SET client_encoding = 'UTF8';
                      false            (           0    0 
   STDSTRINGS 
   STDSTRINGS     (   SET standard_conforming_strings = 'on';
                      false            )           0    0 
   SEARCHPATH 
   SEARCHPATH     8   SELECT pg_catalog.set_config('search_path', '', false);
                      false            *           1262    16486    wmsdb    DATABASE     V   CREATE DATABASE wmsdb WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE = 'C.UTF-8';
    DROP DATABASE wmsdb;
                wmsuser    false            �            1259    16487    barcodes    TABLE     �   CREATE TABLE public.barcodes (
    parent_id integer,
    name character varying(128) DEFAULT ''::character varying NOT NULL,
    barcode_type integer NOT NULL,
    id integer NOT NULL
);
    DROP TABLE public.barcodes;
       public         heap    wmsuser    false            �            1259    16491    barcodes_id_seq    SEQUENCE     �   CREATE SEQUENCE public.barcodes_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 &   DROP SEQUENCE public.barcodes_id_seq;
       public          wmsuser    false    202            +           0    0    barcodes_id_seq    SEQUENCE OWNED BY     C   ALTER SEQUENCE public.barcodes_id_seq OWNED BY public.barcodes.id;
          public          wmsuser    false    203            �            1259    16493    cells    TABLE     M  CREATE TABLE public.cells (
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
    DROP TABLE public.cells;
       public         heap    wmsuser    false            �            1259    16513    cells_id_seq    SEQUENCE     �   CREATE SEQUENCE public.cells_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 #   DROP SEQUENCE public.cells_id_seq;
       public          wmsuser    false    204            ,           0    0    cells_id_seq    SEQUENCE OWNED BY     =   ALTER SEQUENCE public.cells_id_seq OWNED BY public.cells.id;
          public          wmsuser    false    205            �            1259    16515    manufacturers    TABLE     �   CREATE TABLE public.manufacturers (
    id integer NOT NULL,
    name character varying(255) DEFAULT ''::character varying NOT NULL
);
 !   DROP TABLE public.manufacturers;
       public         heap    wmsuser    false            �            1259    16519    manufacturers_id_seq    SEQUENCE     �   CREATE SEQUENCE public.manufacturers_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 +   DROP SEQUENCE public.manufacturers_id_seq;
       public          wmsuser    false    206            -           0    0    manufacturers_id_seq    SEQUENCE OWNED BY     M   ALTER SEQUENCE public.manufacturers_id_seq OWNED BY public.manufacturers.id;
          public          wmsuser    false    207            �            1259    16521    products    TABLE       CREATE TABLE public.products (
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
    DROP TABLE public.products;
       public         heap    wmsuser    false            �            1259    16533    products_id_seq    SEQUENCE     �   CREATE SEQUENCE public.products_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 &   DROP SEQUENCE public.products_id_seq;
       public          wmsuser    false    208            .           0    0    products_id_seq    SEQUENCE OWNED BY     C   ALTER SEQUENCE public.products_id_seq OWNED BY public.products.id;
          public          wmsuser    false    209            �            1259    16535    receipt_headers    TABLE       CREATE TABLE public.receipt_headers (
    id integer NOT NULL,
    number character varying(20) DEFAULT ''::character varying NOT NULL,
    date timestamp with time zone DEFAULT '1970-01-01 03:00:00+03'::timestamp with time zone NOT NULL,
    doc_type integer DEFAULT 0 NOT NULL
);
 #   DROP TABLE public.receipt_headers;
       public         heap    wmsuser    false            �            1259    16541    receipt_docs_id_seq    SEQUENCE     �   CREATE SEQUENCE public.receipt_docs_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 *   DROP SEQUENCE public.receipt_docs_id_seq;
       public          wmsuser    false    210            /           0    0    receipt_docs_id_seq    SEQUENCE OWNED BY     N   ALTER SEQUENCE public.receipt_docs_id_seq OWNED BY public.receipt_headers.id;
          public          wmsuser    false    211            �            1259    16543    receipt_items    TABLE     �   CREATE TABLE public.receipt_items (
    parent_id integer DEFAULT 0 NOT NULL,
    row_id character varying(36) DEFAULT ''::character varying NOT NULL,
    product_id integer DEFAULT 0 NOT NULL,
    quantity integer DEFAULT 0 NOT NULL
);
 !   DROP TABLE public.receipt_items;
       public         heap    wmsuser    false            �            1259    16550    storage1    TABLE     N  CREATE TABLE public.storage1 (
    zone_id integer,
    cell_id integer,
    prod_id integer,
    quantity integer,
    doc_id integer DEFAULT 0 NOT NULL,
    doc_type smallint DEFAULT 0 NOT NULL,
    row_id character varying(36) DEFAULT ''::character varying NOT NULL,
    row_time timestamp with time zone DEFAULT now() NOT NULL
);
    DROP TABLE public.storage1;
       public         heap    wmsuser    false            �            1259    16562    user_printers    TABLE       CREATE TABLE public.user_printers (
    user_id integer DEFAULT 0 NOT NULL,
    printer_name character varying(100) DEFAULT ''::character varying NOT NULL,
    printer_instance character varying(100) DEFAULT ''::character varying NOT NULL,
    printer_type integer DEFAULT 0 NOT NULL
);
 !   DROP TABLE public.user_printers;
       public         heap    wmsuser    false            �            1259    16569    users    TABLE        CREATE TABLE public.users (
    id integer NOT NULL,
    name character varying(150) DEFAULT ''::character varying NOT NULL
);
    DROP TABLE public.users;
       public         heap    wmsuser    false            �            1259    16573    users_id_seq    SEQUENCE     �   CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 #   DROP SEQUENCE public.users_id_seq;
       public          wmsuser    false    215            0           0    0    users_id_seq    SEQUENCE OWNED BY     =   ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;
          public          wmsuser    false    216            �            1259    16575    whs    TABLE     �   CREATE TABLE public.whs (
    id integer NOT NULL,
    name character varying(50) DEFAULT ''::character varying,
    address character varying(255) DEFAULT ''::character varying NOT NULL
);
    DROP TABLE public.whs;
       public         heap    wmsuser    false            �            1259    16580 
   whs_id_seq    SEQUENCE     �   CREATE SEQUENCE public.whs_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 !   DROP SEQUENCE public.whs_id_seq;
       public          wmsuser    false    217            1           0    0 
   whs_id_seq    SEQUENCE OWNED BY     9   ALTER SEQUENCE public.whs_id_seq OWNED BY public.whs.id;
          public          wmsuser    false    218            �            1259    16582    zones    TABLE     �   CREATE TABLE public.zones (
    id integer NOT NULL,
    name character varying(50) DEFAULT ''::character varying,
    parent_id integer,
    zone_type smallint
);
    DROP TABLE public.zones;
       public         heap    wmsuser    false            �            1259    16586    zones_id_seq    SEQUENCE     �   CREATE SEQUENCE public.zones_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 #   DROP SEQUENCE public.zones_id_seq;
       public          wmsuser    false    219            2           0    0    zones_id_seq    SEQUENCE OWNED BY     =   ALTER SEQUENCE public.zones_id_seq OWNED BY public.zones.id;
          public          wmsuser    false    220            E           2604    16588    barcodes id    DEFAULT     j   ALTER TABLE ONLY public.barcodes ALTER COLUMN id SET DEFAULT nextval('public.barcodes_id_seq'::regclass);
 :   ALTER TABLE public.barcodes ALTER COLUMN id DROP DEFAULT;
       public          wmsuser    false    203    202            W           2604    16589    cells id    DEFAULT     d   ALTER TABLE ONLY public.cells ALTER COLUMN id SET DEFAULT nextval('public.cells_id_seq'::regclass);
 7   ALTER TABLE public.cells ALTER COLUMN id DROP DEFAULT;
       public          wmsuser    false    205    204            Y           2604    16590    manufacturers id    DEFAULT     t   ALTER TABLE ONLY public.manufacturers ALTER COLUMN id SET DEFAULT nextval('public.manufacturers_id_seq'::regclass);
 ?   ALTER TABLE public.manufacturers ALTER COLUMN id DROP DEFAULT;
       public          wmsuser    false    207    206            c           2604    16591    products id    DEFAULT     j   ALTER TABLE ONLY public.products ALTER COLUMN id SET DEFAULT nextval('public.products_id_seq'::regclass);
 :   ALTER TABLE public.products ALTER COLUMN id DROP DEFAULT;
       public          wmsuser    false    209    208            g           2604    16592    receipt_headers id    DEFAULT     u   ALTER TABLE ONLY public.receipt_headers ALTER COLUMN id SET DEFAULT nextval('public.receipt_docs_id_seq'::regclass);
 A   ALTER TABLE public.receipt_headers ALTER COLUMN id DROP DEFAULT;
       public          wmsuser    false    211    210            u           2604    16593    users id    DEFAULT     d   ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);
 7   ALTER TABLE public.users ALTER COLUMN id DROP DEFAULT;
       public          wmsuser    false    216    215            x           2604    16594    whs id    DEFAULT     `   ALTER TABLE ONLY public.whs ALTER COLUMN id SET DEFAULT nextval('public.whs_id_seq'::regclass);
 5   ALTER TABLE public.whs ALTER COLUMN id DROP DEFAULT;
       public          wmsuser    false    218    217            z           2604    16595    zones id    DEFAULT     d   ALTER TABLE ONLY public.zones ALTER COLUMN id SET DEFAULT nextval('public.zones_id_seq'::regclass);
 7   ALTER TABLE public.zones ALTER COLUMN id DROP DEFAULT;
       public          wmsuser    false    220    219                      0    16487    barcodes 
   TABLE DATA           E   COPY public.barcodes (parent_id, name, barcode_type, id) FROM stdin;
    public          wmsuser    false    202   [V                 0    16493    cells 
   TABLE DATA           �   COPY public.cells (id, name, whs_id, zone_id, passage_id, rack_id, floor, sz_length, sz_width, sz_height, sz_volume, sz_uf_volume, sz_weight, is_size_free, is_weight_free, not_allowed_in, not_allowed_out, is_service) FROM stdin;
    public          wmsuser    false    204   �W                 0    16515    manufacturers 
   TABLE DATA           1   COPY public.manufacturers (id, name) FROM stdin;
    public          wmsuser    false    206   �W                 0    16521    products 
   TABLE DATA           �   COPY public.products (id, name, manufacturer_id, sz_length, sz_wight, sz_height, sz_weight, sz_volume, sz_uf_volume, item_number) FROM stdin;
    public          wmsuser    false    208   Y                 0    16535    receipt_headers 
   TABLE DATA           E   COPY public.receipt_headers (id, number, date, doc_type) FROM stdin;
    public          wmsuser    false    210   R\                 0    16543    receipt_items 
   TABLE DATA           P   COPY public.receipt_items (parent_id, row_id, product_id, quantity) FROM stdin;
    public          wmsuser    false    212   �\                 0    16550    storage1 
   TABLE DATA           k   COPY public.storage1 (zone_id, cell_id, prod_id, quantity, doc_id, doc_type, row_id, row_time) FROM stdin;
    public          wmsuser    false    213   ]                 0    16562    user_printers 
   TABLE DATA           ^   COPY public.user_printers (user_id, printer_name, printer_instance, printer_type) FROM stdin;
    public          wmsuser    false    214   ^                 0    16569    users 
   TABLE DATA           )   COPY public.users (id, name) FROM stdin;
    public          wmsuser    false    215   "^       !          0    16575    whs 
   TABLE DATA           0   COPY public.whs (id, name, address) FROM stdin;
    public          wmsuser    false    217   u^       #          0    16582    zones 
   TABLE DATA           ?   COPY public.zones (id, name, parent_id, zone_type) FROM stdin;
    public          wmsuser    false    219   �^       3           0    0    barcodes_id_seq    SEQUENCE SET     >   SELECT pg_catalog.setval('public.barcodes_id_seq', 74, true);
          public          wmsuser    false    203            4           0    0    cells_id_seq    SEQUENCE SET     :   SELECT pg_catalog.setval('public.cells_id_seq', 2, true);
          public          wmsuser    false    205            5           0    0    manufacturers_id_seq    SEQUENCE SET     C   SELECT pg_catalog.setval('public.manufacturers_id_seq', 79, true);
          public          wmsuser    false    207            6           0    0    products_id_seq    SEQUENCE SET     >   SELECT pg_catalog.setval('public.products_id_seq', 72, true);
          public          wmsuser    false    209            7           0    0    receipt_docs_id_seq    SEQUENCE SET     B   SELECT pg_catalog.setval('public.receipt_docs_id_seq', 46, true);
          public          wmsuser    false    211            8           0    0    users_id_seq    SEQUENCE SET     ;   SELECT pg_catalog.setval('public.users_id_seq', 10, true);
          public          wmsuser    false    216            9           0    0 
   whs_id_seq    SEQUENCE SET     9   SELECT pg_catalog.setval('public.whs_id_seq', 12, true);
          public          wmsuser    false    218            :           0    0    zones_id_seq    SEQUENCE SET     ;   SELECT pg_catalog.setval('public.zones_id_seq', 21, true);
          public          wmsuser    false    220            |           2606    16597    barcodes barcodes_pk 
   CONSTRAINT     R   ALTER TABLE ONLY public.barcodes
    ADD CONSTRAINT barcodes_pk PRIMARY KEY (id);
 >   ALTER TABLE ONLY public.barcodes DROP CONSTRAINT barcodes_pk;
       public            wmsuser    false    202            �           2606    16599    cells cells_pk 
   CONSTRAINT     L   ALTER TABLE ONLY public.cells
    ADD CONSTRAINT cells_pk PRIMARY KEY (id);
 8   ALTER TABLE ONLY public.cells DROP CONSTRAINT cells_pk;
       public            wmsuser    false    204            �           2606    16601    manufacturers manufacturers_pk 
   CONSTRAINT     \   ALTER TABLE ONLY public.manufacturers
    ADD CONSTRAINT manufacturers_pk PRIMARY KEY (id);
 H   ALTER TABLE ONLY public.manufacturers DROP CONSTRAINT manufacturers_pk;
       public            wmsuser    false    206            �           2606    16603    products products_pk 
   CONSTRAINT     R   ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_pk PRIMARY KEY (id);
 >   ALTER TABLE ONLY public.products DROP CONSTRAINT products_pk;
       public            wmsuser    false    208            �           2606    16605    receipt_headers receipt_docs_pk 
   CONSTRAINT     ]   ALTER TABLE ONLY public.receipt_headers
    ADD CONSTRAINT receipt_docs_pk PRIMARY KEY (id);
 I   ALTER TABLE ONLY public.receipt_headers DROP CONSTRAINT receipt_docs_pk;
       public            wmsuser    false    210            �           2606    16607    receipt_items receipt_items_pk 
   CONSTRAINT     `   ALTER TABLE ONLY public.receipt_items
    ADD CONSTRAINT receipt_items_pk PRIMARY KEY (row_id);
 H   ALTER TABLE ONLY public.receipt_items DROP CONSTRAINT receipt_items_pk;
       public            wmsuser    false    212            �           2606    16609    users users_pk 
   CONSTRAINT     L   ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pk PRIMARY KEY (id);
 8   ALTER TABLE ONLY public.users DROP CONSTRAINT users_pk;
       public            wmsuser    false    215            �           2606    16611 
   whs whs_pk 
   CONSTRAINT     H   ALTER TABLE ONLY public.whs
    ADD CONSTRAINT whs_pk PRIMARY KEY (id);
 4   ALTER TABLE ONLY public.whs DROP CONSTRAINT whs_pk;
       public            wmsuser    false    217            �           2606    16613    zones zones_pk 
   CONSTRAINT     L   ALTER TABLE ONLY public.zones
    ADD CONSTRAINT zones_pk PRIMARY KEY (id);
 8   ALTER TABLE ONLY public.zones DROP CONSTRAINT zones_pk;
       public            wmsuser    false    219            }           1259    16614 /   barcodes_product_id_barcode_barcode_type_uindex    INDEX     �   CREATE UNIQUE INDEX barcodes_product_id_barcode_barcode_type_uindex ON public.barcodes USING btree (parent_id, name, barcode_type);
 C   DROP INDEX public.barcodes_product_id_barcode_barcode_type_uindex;
       public            wmsuser    false    202    202    202            ~           1259    16615    cells_id_uindex    INDEX     F   CREATE UNIQUE INDEX cells_id_uindex ON public.cells USING btree (id);
 #   DROP INDEX public.cells_id_uindex;
       public            wmsuser    false    204            �           1259    16616 ,   cells_whs_id_zone_id_passage_id_floor_uindex    INDEX     �   CREATE UNIQUE INDEX cells_whs_id_zone_id_passage_id_floor_uindex ON public.cells USING btree (whs_id, zone_id, passage_id, floor);
 @   DROP INDEX public.cells_whs_id_zone_id_passage_id_floor_uindex;
       public            wmsuser    false    204    204    204    204            �           1259    16617    receipt_items_row_id_uindex    INDEX     ^   CREATE UNIQUE INDEX receipt_items_row_id_uindex ON public.receipt_items USING btree (row_id);
 /   DROP INDEX public.receipt_items_row_id_uindex;
       public            wmsuser    false    212            �           2606    16618 %   products products_manufacturers_id_fk    FK CONSTRAINT     �   ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_manufacturers_id_fk FOREIGN KEY (manufacturer_id) REFERENCES public.manufacturers(id);
 O   ALTER TABLE ONLY public.products DROP CONSTRAINT products_manufacturers_id_fk;
       public          wmsuser    false    206    2947    208            �           2606    16623 *   receipt_items receipt_items_products_id_fk    FK CONSTRAINT     �   ALTER TABLE ONLY public.receipt_items
    ADD CONSTRAINT receipt_items_products_id_fk FOREIGN KEY (product_id) REFERENCES public.products(id);
 T   ALTER TABLE ONLY public.receipt_items DROP CONSTRAINT receipt_items_products_id_fk;
       public          wmsuser    false    208    212    2949            �           2606    16633    storage1 storage1_cells_id_fk    FK CONSTRAINT     |   ALTER TABLE ONLY public.storage1
    ADD CONSTRAINT storage1_cells_id_fk FOREIGN KEY (cell_id) REFERENCES public.cells(id);
 G   ALTER TABLE ONLY public.storage1 DROP CONSTRAINT storage1_cells_id_fk;
       public          wmsuser    false    213    204    2944               "  x�mRIr1<�+�;�_rL����ƚ�k�����@#�L,j�G�I
��L��4rH'�����Lù�;�i�9ߗ�=sV�zC{=<��9P�ҪP���E��}���y���9�a$X���ǘ��İ�D"��Y�e���D������R�Ǫv�q�}S2
rB�7��h����������e�/[���e4���[3;!�v�V�DyYn��<Q�s�12k��C�;�D�n��'&�&s�1d}�RXz9FW�;m/���Ha�Zw�G��V#�F�~}"�T���         :   x�3�,I-.Q0�A#NNc C=T2��8uA�@א���� �	i����� ^_�         0  x��P�JA<w������s�����l#x�!>P6��S^D�$_�D�<t�5dg�����{����n0��=Gb;��h��L%l����`h/1v0�)Fx����.�0�WΜ�39������.by��'Y��?~c���Z���B��qT7Xe)�iE��q���m'ho�4Z�r��g�DA3*���T5����=<���=����j�v������=��%����T���T�}>m���7�P�#����1qXk�C�#Ƙ�2�u�:O���Q^���"+7�ƛ�����S         +  x��UMOQ]��7�i��f��tk���l��A��	Q!�,� �F4��.��B��vH���#�}ScXL�m�N޻�{�m+>����5�/�㚮���S��=>'~��..�
�}[������`���6f@|@��=w��Z��%Lߟ� �3����ϩk�$Ra������_`w�<m8��'iz`;��09.�V�?i����T
:!���U"�k�J��d��6b�����آwnݛ�Uѡ��J'�!�AP�Bm�QR��o��O,�����DJ�	���O�&(hO���D��#�����tEW��Oi<�N��ˢǾ�Mi��C>�$z��x�éB-��v�i��>3dR߸�dC�ety�D+�:�E�]i���^*)n���$;)������n�,�TZތZ�&6��4P����=�"�6�V��|�"3�V��HC�`��Gژ�K��0�����}�#�q��g���Y��gl5R&,�\����c:��P����gO�i�A�Q��f�:݄�ˇ67� �������J˧�?�l 5�� C4'D��Jg����+�kd�W�J�Ұt�^��5k�4��2Apqm�Z������������i5�Xq�)��hǟɚR}�eC��3+����y��n�K�_]_���3��n׵T��ss%C ��w�o�
`�xm�E�k���z��:�jjz�s�U��z�<�ª>�d�8bR�0ˢ b)&�����3�c ���?��������o ��Դ ���%�a���`���a5��q�
�         ^   x�}�9
�0D�Z>��� �����������߽bd"��b���&F(�[|e`���0i�ʪA
��-0D��]m,��2avx�!��]2|         M   x�%���0B�3�1BR\L��#^N̛O��f4�7��>04�	�%4�Bj�]�d�t�i�oR0Y�������E
         �   x��Kj�1���O�}��F�Β���R
��@�	$F��<#KLǘ����]з�vo�j_����e?O}b�[���ݾR
� !Ě��%l�4[��"�,+Cɢ+�t����Uȫh~��0��d\�:/���LN�$ ^��A=Mr$ծ���jׁ��t�"I[���$r=��\��G�1Ig��վ#��$z� E���ڎ��
��<            x������ � �         C   x�8 ��7	Михаил
3	Админ
10	Новый юзер
\.


��M      !   %   x� ��1	Мой склад	
\.


��      #   T   x�3�0�¾{/lP���bӅ�.6_�~aׅ���F\�H��/6\�qa�=PYC.c���V�� s+���7������ k�7     