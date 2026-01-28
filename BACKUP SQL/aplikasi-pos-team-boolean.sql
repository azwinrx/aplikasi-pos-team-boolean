--
-- PostgreSQL database dump
--

\restrict 36NvCmSxAjwAn74mFsHdeSLpUDKTxtg7WblEed1Pu4XP0VQJN59nYbxZFUFi1SU

-- Dumped from database version 18.1
-- Dumped by pg_dump version 18.1

-- Started on 2026-01-29 02:46:04

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
-- TOC entry 220 (class 1259 OID 20319)
-- Name: categories; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.categories (
    id bigint NOT NULL,
    name character varying(100) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);


ALTER TABLE public.categories OWNER TO postgres;

--
-- TOC entry 219 (class 1259 OID 20318)
-- Name: categories_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.categories_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.categories_id_seq OWNER TO postgres;

--
-- TOC entry 5165 (class 0 OID 0)
-- Dependencies: 219
-- Name: categories_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.categories_id_seq OWNED BY public.categories.id;


--
-- TOC entry 228 (class 1259 OID 20374)
-- Name: inventories; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.inventories (
    id bigint NOT NULL,
    name character varying(255) NOT NULL,
    quantity integer DEFAULT 0 NOT NULL,
    unit character varying(50) NOT NULL,
    min_stock integer DEFAULT 5,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,
    category character varying(100) DEFAULT 'uncategorized'::character varying,
    retail_price numeric(10,2) DEFAULT 0 NOT NULL,
    status character varying(20) DEFAULT 'active'::character varying NOT NULL,
    image character varying(500)
);


ALTER TABLE public.inventories OWNER TO postgres;

--
-- TOC entry 227 (class 1259 OID 20373)
-- Name: inventories_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.inventories_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.inventories_id_seq OWNER TO postgres;

--
-- TOC entry 5166 (class 0 OID 0)
-- Dependencies: 227
-- Name: inventories_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.inventories_id_seq OWNED BY public.inventories.id;


--
-- TOC entry 236 (class 1259 OID 20458)
-- Name: notifications; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.notifications (
    id bigint NOT NULL,
    title character varying(100) NOT NULL,
    message text NOT NULL,
    is_read boolean DEFAULT false,
    type character varying(50),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);


ALTER TABLE public.notifications OWNER TO postgres;

--
-- TOC entry 235 (class 1259 OID 20457)
-- Name: notifications_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.notifications_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.notifications_id_seq OWNER TO postgres;

--
-- TOC entry 5167 (class 0 OID 0)
-- Dependencies: 235
-- Name: notifications_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.notifications_id_seq OWNED BY public.notifications.id;


--
-- TOC entry 232 (class 1259 OID 20417)
-- Name: order_items; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.order_items (
    id bigint NOT NULL,
    order_id bigint,
    product_id bigint,
    quantity integer NOT NULL,
    price numeric(10,2) NOT NULL,
    subtotal numeric(15,2) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);


ALTER TABLE public.order_items OWNER TO postgres;

--
-- TOC entry 231 (class 1259 OID 20416)
-- Name: order_items_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.order_items_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.order_items_id_seq OWNER TO postgres;

--
-- TOC entry 5168 (class 0 OID 0)
-- Dependencies: 231
-- Name: order_items_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.order_items_id_seq OWNED BY public.order_items.id;


--
-- TOC entry 230 (class 1259 OID 20389)
-- Name: orders; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.orders (
    id bigint NOT NULL,
    user_id bigint,
    table_id bigint,
    payment_method_id bigint,
    customer_name character varying(100),
    total_amount numeric(15,2) DEFAULT 0,
    tax numeric(15,2) DEFAULT 0,
    status character varying(20) DEFAULT 'pending'::character varying,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);


ALTER TABLE public.orders OWNER TO postgres;

--
-- TOC entry 229 (class 1259 OID 20388)
-- Name: orders_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.orders_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.orders_id_seq OWNER TO postgres;

--
-- TOC entry 5169 (class 0 OID 0)
-- Dependencies: 229
-- Name: orders_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.orders_id_seq OWNED BY public.orders.id;


--
-- TOC entry 226 (class 1259 OID 20363)
-- Name: payment_methods; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.payment_methods (
    id bigint NOT NULL,
    name character varying(50) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);


ALTER TABLE public.payment_methods OWNER TO postgres;

--
-- TOC entry 225 (class 1259 OID 20362)
-- Name: payment_methods_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.payment_methods_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.payment_methods_id_seq OWNER TO postgres;

--
-- TOC entry 5170 (class 0 OID 0)
-- Dependencies: 225
-- Name: payment_methods_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.payment_methods_id_seq OWNED BY public.payment_methods.id;


--
-- TOC entry 222 (class 1259 OID 20330)
-- Name: products; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.products (
    id bigint NOT NULL,
    category_id bigint,
    name character varying(150) NOT NULL,
    description text,
    price numeric(10,2) NOT NULL,
    image_url text,
    is_available boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);


ALTER TABLE public.products OWNER TO postgres;

--
-- TOC entry 221 (class 1259 OID 20329)
-- Name: products_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.products_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.products_id_seq OWNER TO postgres;

--
-- TOC entry 5171 (class 0 OID 0)
-- Dependencies: 221
-- Name: products_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.products_id_seq OWNED BY public.products.id;


--
-- TOC entry 234 (class 1259 OID 20440)
-- Name: reservations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.reservations (
    id bigint NOT NULL,
    customer_name character varying(100) NOT NULL,
    customer_phone character varying(20),
    table_id bigint,
    reservation_time timestamp with time zone NOT NULL,
    status character varying(20) DEFAULT 'pending'::character varying,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);


ALTER TABLE public.reservations OWNER TO postgres;

--
-- TOC entry 233 (class 1259 OID 20439)
-- Name: reservations_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.reservations_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.reservations_id_seq OWNER TO postgres;

--
-- TOC entry 5172 (class 0 OID 0)
-- Dependencies: 233
-- Name: reservations_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.reservations_id_seq OWNED BY public.reservations.id;


--
-- TOC entry 238 (class 1259 OID 20544)
-- Name: staff; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.staff (
    id integer NOT NULL,
    full_name character varying(100) NOT NULL,
    email character varying(100) NOT NULL,
    role character varying(20) DEFAULT 'staff'::character varying NOT NULL,
    phone_number character varying(20),
    salary numeric(15,2) DEFAULT 0,
    date_of_birth date,
    shift_start_timing character varying(10),
    shift_end_timing character varying(10),
    address text,
    additional_details text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp without time zone
);


ALTER TABLE public.staff OWNER TO postgres;

--
-- TOC entry 237 (class 1259 OID 20543)
-- Name: staff_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.staff_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.staff_id_seq OWNER TO postgres;

--
-- TOC entry 5173 (class 0 OID 0)
-- Dependencies: 237
-- Name: staff_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.staff_id_seq OWNED BY public.staff.id;


--
-- TOC entry 224 (class 1259 OID 20350)
-- Name: tables; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.tables (
    id bigint NOT NULL,
    number character varying(10) NOT NULL,
    capacity integer NOT NULL,
    status character varying(20) DEFAULT 'available'::character varying,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);


ALTER TABLE public.tables OWNER TO postgres;

--
-- TOC entry 223 (class 1259 OID 20349)
-- Name: tables_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.tables_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.tables_id_seq OWNER TO postgres;

--
-- TOC entry 5174 (class 0 OID 0)
-- Dependencies: 223
-- Name: tables_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.tables_id_seq OWNED BY public.tables.id;


--
-- TOC entry 240 (class 1259 OID 20592)
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id bigint NOT NULL,
    name character varying(100) NOT NULL,
    email character varying(100) NOT NULL,
    password character varying(255) NOT NULL,
    role character varying(20) DEFAULT 'staff'::character varying NOT NULL,
    otp character varying(6),
    otp_expiration timestamp with time zone,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);


ALTER TABLE public.users OWNER TO postgres;

--
-- TOC entry 239 (class 1259 OID 20591)
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO postgres;

--
-- TOC entry 5175 (class 0 OID 0)
-- Dependencies: 239
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- TOC entry 4906 (class 2604 OID 20322)
-- Name: categories id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.categories ALTER COLUMN id SET DEFAULT nextval('public.categories_id_seq'::regclass);


--
-- TOC entry 4920 (class 2604 OID 20377)
-- Name: inventories id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.inventories ALTER COLUMN id SET DEFAULT nextval('public.inventories_id_seq'::regclass);


--
-- TOC entry 4941 (class 2604 OID 20461)
-- Name: notifications id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.notifications ALTER COLUMN id SET DEFAULT nextval('public.notifications_id_seq'::regclass);


--
-- TOC entry 4934 (class 2604 OID 20420)
-- Name: order_items id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.order_items ALTER COLUMN id SET DEFAULT nextval('public.order_items_id_seq'::regclass);


--
-- TOC entry 4928 (class 2604 OID 20392)
-- Name: orders id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders ALTER COLUMN id SET DEFAULT nextval('public.orders_id_seq'::regclass);


--
-- TOC entry 4917 (class 2604 OID 20366)
-- Name: payment_methods id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment_methods ALTER COLUMN id SET DEFAULT nextval('public.payment_methods_id_seq'::regclass);


--
-- TOC entry 4909 (class 2604 OID 20333)
-- Name: products id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.products ALTER COLUMN id SET DEFAULT nextval('public.products_id_seq'::regclass);


--
-- TOC entry 4937 (class 2604 OID 20443)
-- Name: reservations id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.reservations ALTER COLUMN id SET DEFAULT nextval('public.reservations_id_seq'::regclass);


--
-- TOC entry 4945 (class 2604 OID 20547)
-- Name: staff id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.staff ALTER COLUMN id SET DEFAULT nextval('public.staff_id_seq'::regclass);


--
-- TOC entry 4913 (class 2604 OID 20353)
-- Name: tables id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tables ALTER COLUMN id SET DEFAULT nextval('public.tables_id_seq'::regclass);


--
-- TOC entry 4950 (class 2604 OID 20595)
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- TOC entry 5139 (class 0 OID 20319)
-- Dependencies: 220
-- Data for Name: categories; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.categories (id, name, created_at, updated_at, deleted_at) FROM stdin;
1	Makanan Berat	2026-01-25 19:40:13.097468+07	2026-01-25 19:40:13.097468+07	\N
2	Minuman	2026-01-25 19:40:13.097468+07	2026-01-25 19:40:13.097468+07	\N
3	Snack	2026-01-25 19:40:13.097468+07	2026-01-25 19:40:13.097468+07	\N
4	Dessert	2026-01-25 19:40:13.097468+07	2026-01-25 19:40:13.097468+07	\N
\.


--
-- TOC entry 5147 (class 0 OID 20374)
-- Dependencies: 228
-- Data for Name: inventories; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.inventories (id, name, quantity, unit, min_stock, created_at, updated_at, deleted_at, category, retail_price, status, image) FROM stdin;
1	Coca Cola 1L	150	litre	50	2026-01-28 23:12:02.707281+07	2026-01-28 23:12:02.707281+07	\N	beverage	15.50	active	https://images.unsplash.com/photo-1554866585-cd94860890b7
2	Sprite 1L	120	litre	50	2026-01-28 23:12:02.707281+07	2026-01-28 23:12:02.707281+07	\N	beverage	14.00	active	https://images.unsplash.com/photo-1629203851122-3726ecdf080e
3	Pepsi 1L	45	litre	50	2026-01-28 23:12:02.707281+07	2026-01-28 23:12:02.707281+07	\N	beverage	15.00	active	https://images.unsplash.com/photo-1622483767028-3f66f32aef97
4	Mineral Water 1.5L	200	litre	100	2026-01-28 23:12:02.707281+07	2026-01-28 23:12:02.707281+07	\N	beverage	5.00	active	https://images.unsplash.com/photo-1560512823-829485b8bf24
5	Orange Juice 1L	80	litre	40	2026-01-28 23:12:02.707281+07	2026-01-28 23:12:02.707281+07	\N	beverage	25.00	active	https://images.unsplash.com/photo-1600271886742-f049cd451bba
6	Apple Juice 1L	75	litre	40	2026-01-28 23:12:02.707281+07	2026-01-28 23:12:02.707281+07	\N	beverage	28.00	active	https://images.unsplash.com/photo-1603569283847-aa295f0d016a
7	Green Tea 500ml	100	litre	60	2026-01-28 23:12:02.707281+07	2026-01-28 23:12:02.707281+07	\N	beverage	12.00	active	https://images.unsplash.com/photo-1546173159-315724a31696
8	Iced Coffee 250ml	90	litre	50	2026-01-28 23:12:02.707281+07	2026-01-28 23:12:02.707281+07	\N	beverage	18.00	active	https://images.unsplash.com/photo-1564890369478-c89ca6d9cde9
9	Energy Drink 250ml	60	litre	40	2026-01-28 23:12:02.707281+07	2026-01-28 23:12:02.707281+07	\N	beverage	22.00	active	https://images.unsplash.com/photo-1572490122747-3968b75cc699
10	Lemonade 1L	55	litre	50	2026-01-28 23:12:02.707281+07	2026-01-28 23:12:02.707281+07	\N	beverage	16.00	active	https://images.unsplash.com/photo-1625772452859-1c03d5bf1137
11	Chocolate Milk 500ml	70	litre	40	2026-01-28 23:12:02.707281+07	2026-01-28 23:12:02.707281+07	\N	beverage	20.00	active	https://images.unsplash.com/photo-1568702846914-96b305d2aaeb
12	Strawberry Smoothie 500ml	40	litre	30	2026-01-28 23:12:02.707281+07	2026-01-28 23:12:02.707281+07	\N	beverage	30.00	active	https://images.unsplash.com/photo-1523294587484-bae6cc870010
13	Potato Chips 100g	180	pcs	100	2026-01-28 23:12:02.707281+07	2026-01-28 23:12:02.707281+07	\N	snack	8.50	active	https://images.unsplash.com/photo-1610970881699-44a5587cabec
14	Chocolate Bar 50g	200	pcs	150	2026-01-28 23:12:02.707281+07	2026-01-28 23:12:02.707281+07	\N	snack	6.00	active	https://images.unsplash.com/photo-1613919113640-c65cf1d676f5
15	Cookies Pack 200g	120	pcs	80	2026-01-28 23:12:02.707281+07	2026-01-28 23:12:02.707281+07	\N	snack	12.00	active	https://images.unsplash.com/photo-1621939514649-280e2ee25f60
16	Instant Noodles	250	pcs	150	2026-01-28 23:12:02.707281+07	2026-01-28 23:12:02.707281+07	\N	food	5.50	active	https://images.unsplash.com/photo-1599490659213-e2b9527bd087
17	Canned Tuna 185g	100	pcs	50	2026-01-28 23:12:02.707281+07	2026-01-28 23:12:02.707281+07	\N	food	18.00	active	https://images.unsplash.com/photo-1588137378633-dea1336ce1e2
18	Rice 5kg	80	kg	40	2026-01-28 23:12:02.707281+07	2026-01-28 23:12:02.707281+07	\N	food	45.00	active	https://images.unsplash.com/photo-1589367920969-ab8e050bbb04
19	Cooking Oil 1L	90	litre	50	2026-01-28 23:12:02.707281+07	2026-01-28 23:12:02.707281+07	\N	food	28.00	active	https://images.unsplash.com/photo-1586201375761-83865001e31c
20	Sugar 1kg	110	kg	60	2026-01-28 23:12:02.707281+07	2026-01-28 23:12:02.707281+07	\N	food	12.00	active	https://images.unsplash.com/photo-1563379091339-03b21ab4a4f8
21	Coca Cola	100		5	2026-01-28 23:16:03.495595+07	2026-01-28 23:16:03.495595+07	\N	beverage	15.50	active	https://example.com/image.jpg
\.


--
-- TOC entry 5155 (class 0 OID 20458)
-- Dependencies: 236
-- Data for Name: notifications; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.notifications (id, title, message, is_read, type, created_at, updated_at, deleted_at) FROM stdin;
1	Stok Menipis	Stok Minyak Goreng tersisa 5 liter	f	alert	2026-01-25 19:40:13.097468+07	2026-01-25 19:40:13.097468+07	\N
2	Order Baru	Meja T02 melakukan pemesanan	f	order	2026-01-25 19:40:13.097468+07	2026-01-25 19:40:13.097468+07	\N
\.


--
-- TOC entry 5151 (class 0 OID 20417)
-- Dependencies: 232
-- Data for Name: order_items; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.order_items (id, order_id, product_id, quantity, price, subtotal, created_at, updated_at, deleted_at) FROM stdin;
1	1	1	1	25000.00	25000.00	2026-01-25 19:40:13.097468+07	2026-01-25 19:40:13.097468+07	\N
2	1	2	1	30000.00	30000.00	2026-01-25 19:40:13.097468+07	2026-01-25 19:40:13.097468+07	\N
3	2	2	1	30000.00	30000.00	2026-01-25 19:40:13.097468+07	2026-01-25 19:40:13.097468+07	\N
\.


--
-- TOC entry 5149 (class 0 OID 20389)
-- Dependencies: 230
-- Data for Name: orders; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.orders (id, user_id, table_id, payment_method_id, customer_name, total_amount, tax, status, created_at, updated_at, deleted_at) FROM stdin;
1	3	1	1	Customer A	55000.00	5500.00	paid	2026-01-23 19:40:13.097468+07	2026-01-25 19:40:13.097468+07	\N
2	3	2	2	Customer B	30000.00	3000.00	paid	2026-01-25 19:40:13.097468+07	2026-01-25 19:40:13.097468+07	\N
\.


--
-- TOC entry 5145 (class 0 OID 20363)
-- Dependencies: 226
-- Data for Name: payment_methods; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.payment_methods (id, name, created_at, updated_at, deleted_at) FROM stdin;
1	Cash	2026-01-25 19:40:13.097468+07	2026-01-25 19:40:13.097468+07	\N
2	QRIS	2026-01-25 19:40:13.097468+07	2026-01-25 19:40:13.097468+07	\N
3	Debit Card	2026-01-25 19:40:13.097468+07	2026-01-25 19:40:13.097468+07	\N
\.


--
-- TOC entry 5141 (class 0 OID 20330)
-- Dependencies: 222
-- Data for Name: products; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.products (id, category_id, name, description, price, image_url, is_available, created_at, updated_at, deleted_at) FROM stdin;
1	1	Nasi Goreng Spesial	Nasi goreng dengan telur dan ayam	25000.00	\N	t	2025-12-16 19:40:13.097468+07	2026-01-25 19:40:13.097468+07	\N
2	1	Ayam Bakar Madu	Ayam bakar oles madu	30000.00	\N	t	2026-01-20 19:40:13.097468+07	2026-01-25 19:40:13.097468+07	\N
3	2	Es Teh Manis	Teh manis dingin segar	5000.00	\N	t	2025-12-16 19:40:13.097468+07	2026-01-25 19:40:13.097468+07	\N
4	2	Kopi Susu Gula Aren	Kopi kekinian	18000.00	\N	t	2026-01-23 19:40:13.097468+07	2026-01-25 19:40:13.097468+07	\N
5	3	Kentang Goreng	French fries original	15000.00	\N	t	2025-12-16 19:40:13.097468+07	2026-01-25 19:40:13.097468+07	\N
\.


--
-- TOC entry 5153 (class 0 OID 20440)
-- Dependencies: 234
-- Data for Name: reservations; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.reservations (id, customer_name, customer_phone, table_id, reservation_time, status, created_at, updated_at, deleted_at) FROM stdin;
\.


--
-- TOC entry 5157 (class 0 OID 20544)
-- Dependencies: 238
-- Data for Name: staff; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.staff (id, full_name, email, role, phone_number, salary, date_of_birth, shift_start_timing, shift_end_timing, address, additional_details, created_at, updated_at, deleted_at) FROM stdin;
1	Budi Santoso	budi.santoso@pos.com	admin	081234567890	12000000.00	1985-03-15	08:00:00	17:00:00	Jl. Sudirman No. 123, Jakarta Pusat, DKI Jakarta 10110	System administrator dengan full access ke semua fitur. Bertanggung jawab atas keamanan sistem.	2026-01-28 23:54:05.774921	2026-01-28 23:54:05.774921	\N
2	Siti Nurhaliza	siti.nurhaliza@pos.com	manager	081234567891	10000000.00	1988-07-20	08:00:00	17:00:00	Jl. Gatot Subroto No. 45, Jakarta Selatan, DKI Jakarta 12190	Store manager dengan pengalaman 8 tahun. Mengelola operasional harian dan staff.	2026-01-28 23:54:05.774921	2026-01-28 23:54:05.774921	\N
3	Andi Wijaya	andi.wijaya@pos.com	cashier	081234567892	5500000.00	1995-01-10	07:00:00	15:00:00	Jl. Kebon Jeruk No. 78, Jakarta Barat, DKI Jakarta 11530	Kasir shift pagi. Cepat dan akurat dalam bertransaksi. Pelayanan customer excellent.	2026-01-28 23:54:05.774921	2026-01-28 23:54:05.774921	\N
4	Dewi Lestari	dewi.lestari@pos.com	cashier	081234567893	5500000.00	1997-05-25	07:00:00	15:00:00	Jl. Raya Bogor No. 234, Jakarta Timur, DKI Jakarta 13770	Kasir shift pagi. Teliti dan ramah. Customer satisfaction rating tinggi.	2026-01-28 23:54:05.774921	2026-01-28 23:54:05.774921	\N
5	Rudi Hartono	rudi.hartono@pos.com	cashier	081234567894	5500000.00	1996-11-30	14:00:00	22:00:00	Jl. Pasar Minggu No. 56, Jakarta Selatan, DKI Jakarta 12520	Kasir shift sore. Berpengalaman handling rush hour. Fast worker.	2026-01-28 23:54:05.774921	2026-01-28 23:54:05.774921	\N
6	Linda Wijayanti	linda.wijayanti@pos.com	cashier	081234567895	5500000.00	1998-08-17	14:00:00	22:00:00	Jl. Cempaka Putih No. 89, Jakarta Pusat, DKI Jakarta 10510	Kasir shift sore. Multitasking dan dapat handle customer dengan baik.	2026-01-28 23:54:05.774921	2026-01-28 23:54:05.774921	\N
7	Ahmad Fauzi	ahmad.fauzi@pos.com	staff	081234567896	4500000.00	1992-04-12	09:00:00	18:00:00	Jl. Ciputat Raya No. 45, Tangerang Selatan, Banten 15412	Staff inventory. Bertanggung jawab atas stock opname dan pengelolaan gudang.	2026-01-28 23:54:05.774921	2026-01-28 23:54:05.774921	\N
8	Yuni Kartika	yuni.kartika@pos.com	staff	081234567897	4500000.00	1994-09-08	09:00:00	18:00:00	Jl. Kebayoran Lama No. 67, Jakarta Selatan, DKI Jakarta 12220	Staff general. Handling customer service dan inventory support.	2026-01-28 23:54:05.774921	2026-01-28 23:54:05.774921	\N
9	Joko Susilo	joko.susilo@pos.com	staff	081234567898	4000000.00	1990-12-25	06:00:00	14:00:00	Jl. Pemuda No. 123, Bekasi, Jawa Barat 17141	Staff cleaning dan maintenance. Menjaga kebersihan dan perawatan toko.	2026-01-28 23:54:05.774921	2026-01-28 23:54:05.774921	\N
10	Mega Putri	mega.putri@pos.com	supervisor	081234567899	7500000.00	1991-06-14	10:00:00	19:00:00	Jl. Tebet Raya No. 34, Jakarta Selatan, DKI Jakarta 12810	Supervisor operasional. Mengawasi kinerja staff dan quality control.	2026-01-28 23:54:05.774921	2026-01-28 23:54:05.774921	\N
11	Alice Johnson Updated	alice.johnson@example.com	supervisor	081234567893	4000000.00	1995-05-15	09:00:00	17:00:00	Semarang, Indonesia	Promoted to supervisor	2026-01-29 02:09:44.801369	2026-01-29 02:14:29.078698	2026-01-29 02:15:14.929281
\.


--
-- TOC entry 5143 (class 0 OID 20350)
-- Dependencies: 224
-- Data for Name: tables; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.tables (id, number, capacity, status, created_at, updated_at, deleted_at) FROM stdin;
1	T01	4	available	2026-01-25 19:40:13.097468+07	2026-01-25 19:40:13.097468+07	\N
2	T02	2	occupied	2026-01-25 19:40:13.097468+07	2026-01-25 19:40:13.097468+07	\N
3	T03	6	reserved	2026-01-25 19:40:13.097468+07	2026-01-25 19:40:13.097468+07	\N
4	T04	4	available	2026-01-25 19:40:13.097468+07	2026-01-25 19:40:13.097468+07	\N
5	T05	2	available	2026-01-25 19:40:13.097468+07	2026-01-25 19:40:13.097468+07	\N
\.


--
-- TOC entry 5159 (class 0 OID 20592)
-- Dependencies: 240
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, name, email, password, role, otp, otp_expiration, created_at, updated_at, deleted_at) FROM stdin;
\.


--
-- TOC entry 5176 (class 0 OID 0)
-- Dependencies: 219
-- Name: categories_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.categories_id_seq', 4, true);


--
-- TOC entry 5177 (class 0 OID 0)
-- Dependencies: 227
-- Name: inventories_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.inventories_id_seq', 21, true);


--
-- TOC entry 5178 (class 0 OID 0)
-- Dependencies: 235
-- Name: notifications_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.notifications_id_seq', 2, true);


--
-- TOC entry 5179 (class 0 OID 0)
-- Dependencies: 231
-- Name: order_items_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.order_items_id_seq', 3, true);


--
-- TOC entry 5180 (class 0 OID 0)
-- Dependencies: 229
-- Name: orders_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.orders_id_seq', 2, true);


--
-- TOC entry 5181 (class 0 OID 0)
-- Dependencies: 225
-- Name: payment_methods_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.payment_methods_id_seq', 3, true);


--
-- TOC entry 5182 (class 0 OID 0)
-- Dependencies: 221
-- Name: products_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.products_id_seq', 5, true);


--
-- TOC entry 5183 (class 0 OID 0)
-- Dependencies: 233
-- Name: reservations_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.reservations_id_seq', 1, false);


--
-- TOC entry 5184 (class 0 OID 0)
-- Dependencies: 237
-- Name: staff_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.staff_id_seq', 11, true);


--
-- TOC entry 5185 (class 0 OID 0)
-- Dependencies: 223
-- Name: tables_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.tables_id_seq', 5, true);


--
-- TOC entry 5186 (class 0 OID 0)
-- Dependencies: 239
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_id_seq', 1, false);


--
-- TOC entry 4955 (class 2606 OID 20328)
-- Name: categories categories_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_pkey PRIMARY KEY (id);


--
-- TOC entry 4964 (class 2606 OID 20387)
-- Name: inventories inventories_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.inventories
    ADD CONSTRAINT inventories_pkey PRIMARY KEY (id);


--
-- TOC entry 4972 (class 2606 OID 20471)
-- Name: notifications notifications_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.notifications
    ADD CONSTRAINT notifications_pkey PRIMARY KEY (id);


--
-- TOC entry 4968 (class 2606 OID 20428)
-- Name: order_items order_items_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT order_items_pkey PRIMARY KEY (id);


--
-- TOC entry 4966 (class 2606 OID 20400)
-- Name: orders orders_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_pkey PRIMARY KEY (id);


--
-- TOC entry 4961 (class 2606 OID 20372)
-- Name: payment_methods payment_methods_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment_methods
    ADD CONSTRAINT payment_methods_pkey PRIMARY KEY (id);


--
-- TOC entry 4957 (class 2606 OID 20343)
-- Name: products products_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_pkey PRIMARY KEY (id);


--
-- TOC entry 4970 (class 2606 OID 20451)
-- Name: reservations reservations_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.reservations
    ADD CONSTRAINT reservations_pkey PRIMARY KEY (id);


--
-- TOC entry 4978 (class 2606 OID 20563)
-- Name: staff staff_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.staff
    ADD CONSTRAINT staff_email_key UNIQUE (email);


--
-- TOC entry 4980 (class 2606 OID 20561)
-- Name: staff staff_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.staff
    ADD CONSTRAINT staff_pkey PRIMARY KEY (id);


--
-- TOC entry 4959 (class 2606 OID 20361)
-- Name: tables tables_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tables
    ADD CONSTRAINT tables_pkey PRIMARY KEY (id);


--
-- TOC entry 4982 (class 2606 OID 20607)
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- TOC entry 4984 (class 2606 OID 20605)
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- TOC entry 4962 (class 1259 OID 20482)
-- Name: idx_inventories_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_inventories_deleted_at ON public.inventories USING btree (deleted_at);


--
-- TOC entry 4973 (class 1259 OID 20567)
-- Name: idx_staff_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_staff_deleted_at ON public.staff USING btree (deleted_at);


--
-- TOC entry 4974 (class 1259 OID 20564)
-- Name: idx_staff_email; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_staff_email ON public.staff USING btree (email);


--
-- TOC entry 4975 (class 1259 OID 20566)
-- Name: idx_staff_full_name; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_staff_full_name ON public.staff USING btree (full_name);


--
-- TOC entry 4976 (class 1259 OID 20565)
-- Name: idx_staff_role; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_staff_role ON public.staff USING btree (role);


--
-- TOC entry 4988 (class 2606 OID 20429)
-- Name: order_items order_items_order_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT order_items_order_id_fkey FOREIGN KEY (order_id) REFERENCES public.orders(id);


--
-- TOC entry 4989 (class 2606 OID 20434)
-- Name: order_items order_items_product_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT order_items_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.products(id);


--
-- TOC entry 4986 (class 2606 OID 20411)
-- Name: orders orders_payment_method_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_payment_method_id_fkey FOREIGN KEY (payment_method_id) REFERENCES public.payment_methods(id);


--
-- TOC entry 4987 (class 2606 OID 20406)
-- Name: orders orders_table_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_table_id_fkey FOREIGN KEY (table_id) REFERENCES public.tables(id);


--
-- TOC entry 4985 (class 2606 OID 20344)
-- Name: products products_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.categories(id);


--
-- TOC entry 4990 (class 2606 OID 20452)
-- Name: reservations reservations_table_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.reservations
    ADD CONSTRAINT reservations_table_id_fkey FOREIGN KEY (table_id) REFERENCES public.tables(id);


-- Completed on 2026-01-29 02:46:04

--
-- PostgreSQL database dump complete
--

\unrestrict 36NvCmSxAjwAn74mFsHdeSLpUDKTxtg7WblEed1Pu4XP0VQJN59nYbxZFUFi1SU

