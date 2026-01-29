--
-- PostgreSQL database dump
--

\restrict yZDJ0rbbsQXI9RW9vJwhRRL8W11ZscdWCTCG655q4pA0kwD9ocpQwI9MU7ddeEQ

-- Dumped from database version 18.1
-- Dumped by pg_dump version 18.1

-- Started on 2026-01-29 17:40:53

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
-- TOC entry 238 (class 1259 OID 20658)
-- Name: inventories; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.inventories (
    id bigint NOT NULL,
    image character varying(500),
    name character varying(255) NOT NULL,
    category character varying(100) DEFAULT 'uncategorized'::character varying,
    quantity integer DEFAULT 0 NOT NULL,
    status character varying(20) DEFAULT 'active'::character varying NOT NULL,
    retail_price numeric(10,2) DEFAULT 0 NOT NULL,
    unit character varying(50) NOT NULL,
    min_stock integer DEFAULT 5 NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp with time zone
);


ALTER TABLE public.inventories OWNER TO postgres;

--
-- TOC entry 237 (class 1259 OID 20657)
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
-- Dependencies: 237
-- Name: inventories_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.inventories_id_seq OWNED BY public.inventories.id;


--
-- TOC entry 232 (class 1259 OID 20458)
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
-- TOC entry 231 (class 1259 OID 20457)
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
-- Dependencies: 231
-- Name: notifications_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.notifications_id_seq OWNED BY public.notifications.id;


--
-- TOC entry 228 (class 1259 OID 20417)
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
-- TOC entry 227 (class 1259 OID 20416)
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
-- Dependencies: 227
-- Name: order_items_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.order_items_id_seq OWNED BY public.order_items.id;


--
-- TOC entry 226 (class 1259 OID 20389)
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
-- TOC entry 225 (class 1259 OID 20388)
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
-- Dependencies: 225
-- Name: orders_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.orders_id_seq OWNED BY public.orders.id;


--
-- TOC entry 224 (class 1259 OID 20363)
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
-- TOC entry 223 (class 1259 OID 20362)
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
-- Dependencies: 223
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
-- TOC entry 230 (class 1259 OID 20440)
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
-- TOC entry 229 (class 1259 OID 20439)
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
-- Dependencies: 229
-- Name: reservations_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.reservations_id_seq OWNED BY public.reservations.id;


--
-- TOC entry 236 (class 1259 OID 20634)
-- Name: staff; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.staff (
    id bigint NOT NULL,
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
    deleted_at timestamp with time zone
);


ALTER TABLE public.staff OWNER TO postgres;

--
-- TOC entry 235 (class 1259 OID 20633)
-- Name: staff_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.staff_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.staff_id_seq OWNER TO postgres;

--
-- TOC entry 5173 (class 0 OID 0)
-- Dependencies: 235
-- Name: staff_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.staff_id_seq OWNED BY public.staff.id;


--
-- TOC entry 240 (class 1259 OID 20684)
-- Name: tables; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.tables (
    id bigint NOT NULL,
    number character varying(10) NOT NULL,
    capacity bigint NOT NULL,
    status character varying(20) DEFAULT 'available'::character varying NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp with time zone
);


ALTER TABLE public.tables OWNER TO postgres;

--
-- TOC entry 239 (class 1259 OID 20683)
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
-- Dependencies: 239
-- Name: tables_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.tables_id_seq OWNED BY public.tables.id;


--
-- TOC entry 234 (class 1259 OID 20592)
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
-- TOC entry 233 (class 1259 OID 20591)
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
-- Dependencies: 233
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- TOC entry 4906 (class 2604 OID 20322)
-- Name: categories id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.categories ALTER COLUMN id SET DEFAULT nextval('public.categories_id_seq'::regclass);


--
-- TOC entry 4942 (class 2604 OID 20661)
-- Name: inventories id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.inventories ALTER COLUMN id SET DEFAULT nextval('public.inventories_id_seq'::regclass);


--
-- TOC entry 4929 (class 2604 OID 20461)
-- Name: notifications id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.notifications ALTER COLUMN id SET DEFAULT nextval('public.notifications_id_seq'::regclass);


--
-- TOC entry 4922 (class 2604 OID 20420)
-- Name: order_items id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.order_items ALTER COLUMN id SET DEFAULT nextval('public.order_items_id_seq'::regclass);


--
-- TOC entry 4916 (class 2604 OID 20392)
-- Name: orders id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders ALTER COLUMN id SET DEFAULT nextval('public.orders_id_seq'::regclass);


--
-- TOC entry 4913 (class 2604 OID 20366)
-- Name: payment_methods id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment_methods ALTER COLUMN id SET DEFAULT nextval('public.payment_methods_id_seq'::regclass);


--
-- TOC entry 4909 (class 2604 OID 20333)
-- Name: products id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.products ALTER COLUMN id SET DEFAULT nextval('public.products_id_seq'::regclass);


--
-- TOC entry 4925 (class 2604 OID 20443)
-- Name: reservations id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.reservations ALTER COLUMN id SET DEFAULT nextval('public.reservations_id_seq'::regclass);


--
-- TOC entry 4937 (class 2604 OID 20637)
-- Name: staff id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.staff ALTER COLUMN id SET DEFAULT nextval('public.staff_id_seq'::regclass);


--
-- TOC entry 4950 (class 2604 OID 20687)
-- Name: tables id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tables ALTER COLUMN id SET DEFAULT nextval('public.tables_id_seq'::regclass);


--
-- TOC entry 4933 (class 2604 OID 20595)
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
-- TOC entry 5157 (class 0 OID 20658)
-- Dependencies: 238
-- Data for Name: inventories; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.inventories (id, image, name, category, quantity, status, retail_price, unit, min_stock, created_at, updated_at, deleted_at) FROM stdin;
1		Coca Cola 1L	beverage	150	active	15.50	litre	50	2026-01-29 17:35:11.216331+07	2026-01-29 17:35:11.216331+07	\N
2		Sprite 1L	beverage	120	active	14.00	litre	50	2026-01-29 17:35:11.216331+07	2026-01-29 17:35:11.216331+07	\N
3		Pepsi 1L	beverage	45	active	15.00	litre	50	2026-01-29 17:35:11.216331+07	2026-01-29 17:35:11.216331+07	\N
4		Mineral Water 1.5L	beverage	200	active	5.00	litre	100	2026-01-29 17:35:11.216331+07	2026-01-29 17:35:11.216331+07	\N
5		Orange Juice 1L	beverage	30	active	25.00	litre	40	2026-01-29 17:35:11.216331+07	2026-01-29 17:35:11.216331+07	\N
\.


--
-- TOC entry 5151 (class 0 OID 20458)
-- Dependencies: 232
-- Data for Name: notifications; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.notifications (id, title, message, is_read, type, created_at, updated_at, deleted_at) FROM stdin;
1	Stok Menipis	Stok Minyak Goreng tersisa 5 liter	f	alert	2026-01-25 19:40:13.097468+07	2026-01-25 19:40:13.097468+07	\N
2	Order Baru	Meja T02 melakukan pemesanan	f	order	2026-01-25 19:40:13.097468+07	2026-01-25 19:40:13.097468+07	\N
\.


--
-- TOC entry 5147 (class 0 OID 20417)
-- Dependencies: 228
-- Data for Name: order_items; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.order_items (id, order_id, product_id, quantity, price, subtotal, created_at, updated_at, deleted_at) FROM stdin;
3	2	2	1	30000.00	30000.00	2026-01-25 19:40:13.097468+07	2026-01-25 19:40:13.097468+07	\N
4	3	1	2	25000.00	50000.00	2026-01-29 17:10:23.429676+07	2026-01-29 17:10:23.429676+07	\N
5	3	2	1	30000.00	30000.00	2026-01-29 17:10:23.429676+07	2026-01-29 17:10:23.429676+07	\N
1	1	1	1	25000.00	25000.00	2026-01-25 19:40:13.097468+07	2026-01-25 19:40:13.097468+07	2026-01-29 17:10:30.850901+07
2	1	2	1	30000.00	30000.00	2026-01-25 19:40:13.097468+07	2026-01-25 19:40:13.097468+07	2026-01-29 17:10:30.850901+07
6	1	1	3	25000.00	75000.00	2026-01-29 17:10:30.852065+07	2026-01-29 17:10:30.852065+07	2026-01-29 17:11:56.218749+07
8	1	1	3	25000.00	75000.00	2026-01-29 17:11:56.219254+07	2026-01-29 17:11:56.219254+07	2026-01-29 17:12:02.545871+07
9	1	1	3	25000.00	75000.00	2026-01-29 17:12:02.546436+07	2026-01-29 17:12:02.546436+07	2026-01-29 17:12:41.477858+07
10	1	1	3	25000.00	75000.00	2026-01-29 17:12:41.478385+07	2026-01-29 17:12:41.478385+07	2026-01-29 17:12:59.077236+07
11	1	1	4	25000.00	100000.00	2026-01-29 17:12:59.077799+07	2026-01-29 17:12:59.077799+07	2026-01-29 17:18:05.114192+07
12	1	1	4	25000.00	100000.00	2026-01-29 17:18:05.116216+07	2026-01-29 17:18:05.116216+07	\N
\.


--
-- TOC entry 5145 (class 0 OID 20389)
-- Dependencies: 226
-- Data for Name: orders; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.orders (id, user_id, table_id, payment_method_id, customer_name, total_amount, tax, status, created_at, updated_at, deleted_at) FROM stdin;
2	3	2	2	Customer B	30000.00	3000.00	paid	2026-01-25 19:40:13.097468+07	2026-01-25 19:40:13.097468+07	\N
3	1	1	1	John Doe	85500.00	5500.00	pending	2026-01-29 17:10:23.376272+07	2026-01-29 17:10:23.376272+07	\N
1	3	1	2	Jane Doe	100000.00	5500.00	paid	2026-01-23 19:40:13.097468+07	2026-01-29 17:18:05.112011+07	\N
\.


--
-- TOC entry 5143 (class 0 OID 20363)
-- Dependencies: 224
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
-- TOC entry 5149 (class 0 OID 20440)
-- Dependencies: 230
-- Data for Name: reservations; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.reservations (id, customer_name, customer_phone, table_id, reservation_time, status, created_at, updated_at, deleted_at) FROM stdin;
\.


--
-- TOC entry 5155 (class 0 OID 20634)
-- Dependencies: 236
-- Data for Name: staff; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.staff (id, full_name, email, role, phone_number, salary, date_of_birth, shift_start_timing, shift_end_timing, address, additional_details, created_at, updated_at, deleted_at) FROM stdin;
1	John Doe	john.doe@example.com	manager	081234567890	5000000.00	\N			Jakarta, Indonesia		2026-01-29 17:35:11.221448	2026-01-29 17:35:11.221448	\N
2	Jane Smith	jane.smith@example.com	cashier	081234567891	3500000.00	\N			Bandung, Indonesia		2026-01-29 17:35:11.221448	2026-01-29 17:35:11.221448	\N
3	Bob Wilson	bob.wilson@example.com	staff	081234567892	4000000.00	\N			Surabaya, Indonesia		2026-01-29 17:35:11.221448	2026-01-29 17:35:11.221448	\N
\.


--
-- TOC entry 5159 (class 0 OID 20684)
-- Dependencies: 240
-- Data for Name: tables; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.tables (id, number, capacity, status, created_at, updated_at, deleted_at) FROM stdin;
1	T01	4	available	2026-01-29 17:35:11.22516	2026-01-29 17:35:11.22516	\N
2	T02	2	available	2026-01-29 17:35:11.22516	2026-01-29 17:35:11.22516	\N
3	T03	6	available	2026-01-29 17:35:11.22516	2026-01-29 17:35:11.22516	\N
4	T04	4	available	2026-01-29 17:35:11.22516	2026-01-29 17:35:11.22516	\N
5	T05	2	available	2026-01-29 17:35:11.22516	2026-01-29 17:35:11.22516	\N
\.


--
-- TOC entry 5153 (class 0 OID 20592)
-- Dependencies: 234
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
-- Dependencies: 237
-- Name: inventories_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.inventories_id_seq', 5, true);


--
-- TOC entry 5178 (class 0 OID 0)
-- Dependencies: 231
-- Name: notifications_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.notifications_id_seq', 2, true);


--
-- TOC entry 5179 (class 0 OID 0)
-- Dependencies: 227
-- Name: order_items_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.order_items_id_seq', 12, true);


--
-- TOC entry 5180 (class 0 OID 0)
-- Dependencies: 225
-- Name: orders_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.orders_id_seq', 3, true);


--
-- TOC entry 5181 (class 0 OID 0)
-- Dependencies: 223
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
-- Dependencies: 229
-- Name: reservations_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.reservations_id_seq', 1, false);


--
-- TOC entry 5184 (class 0 OID 0)
-- Dependencies: 235
-- Name: staff_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.staff_id_seq', 3, true);


--
-- TOC entry 5185 (class 0 OID 0)
-- Dependencies: 239
-- Name: tables_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.tables_id_seq', 5, true);


--
-- TOC entry 5186 (class 0 OID 0)
-- Dependencies: 233
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
-- TOC entry 4981 (class 2606 OID 20681)
-- Name: inventories inventories_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.inventories
    ADD CONSTRAINT inventories_pkey PRIMARY KEY (id);


--
-- TOC entry 4967 (class 2606 OID 20471)
-- Name: notifications notifications_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.notifications
    ADD CONSTRAINT notifications_pkey PRIMARY KEY (id);


--
-- TOC entry 4963 (class 2606 OID 20428)
-- Name: order_items order_items_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT order_items_pkey PRIMARY KEY (id);


--
-- TOC entry 4961 (class 2606 OID 20400)
-- Name: orders orders_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_pkey PRIMARY KEY (id);


--
-- TOC entry 4959 (class 2606 OID 20372)
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
-- TOC entry 4965 (class 2606 OID 20451)
-- Name: reservations reservations_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.reservations
    ADD CONSTRAINT reservations_pkey PRIMARY KEY (id);


--
-- TOC entry 4976 (class 2606 OID 20651)
-- Name: staff staff_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.staff
    ADD CONSTRAINT staff_pkey PRIMARY KEY (id);


--
-- TOC entry 4984 (class 2606 OID 20698)
-- Name: tables tables_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tables
    ADD CONSTRAINT tables_pkey PRIMARY KEY (id);


--
-- TOC entry 4978 (class 2606 OID 20653)
-- Name: staff uni_staff_email; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.staff
    ADD CONSTRAINT uni_staff_email UNIQUE (email);


--
-- TOC entry 4986 (class 2606 OID 20700)
-- Name: tables uni_tables_number; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tables
    ADD CONSTRAINT uni_tables_number UNIQUE (number);


--
-- TOC entry 4969 (class 2606 OID 20607)
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- TOC entry 4971 (class 2606 OID 20605)
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- TOC entry 4979 (class 1259 OID 20682)
-- Name: idx_inventories_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_inventories_deleted_at ON public.inventories USING btree (deleted_at);


--
-- TOC entry 4972 (class 1259 OID 20654)
-- Name: idx_staff_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_staff_deleted_at ON public.staff USING btree (deleted_at);


--
-- TOC entry 4973 (class 1259 OID 20656)
-- Name: idx_staff_email; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_staff_email ON public.staff USING btree (email);


--
-- TOC entry 4974 (class 1259 OID 20655)
-- Name: idx_staff_role; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_staff_role ON public.staff USING btree (role);


--
-- TOC entry 4982 (class 1259 OID 20701)
-- Name: idx_tables_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_tables_deleted_at ON public.tables USING btree (deleted_at);


--
-- TOC entry 4989 (class 2606 OID 20429)
-- Name: order_items order_items_order_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT order_items_order_id_fkey FOREIGN KEY (order_id) REFERENCES public.orders(id);


--
-- TOC entry 4990 (class 2606 OID 20434)
-- Name: order_items order_items_product_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT order_items_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.products(id);


--
-- TOC entry 4988 (class 2606 OID 20411)
-- Name: orders orders_payment_method_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_payment_method_id_fkey FOREIGN KEY (payment_method_id) REFERENCES public.payment_methods(id);


--
-- TOC entry 4987 (class 2606 OID 20344)
-- Name: products products_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.categories(id);


-- Completed on 2026-01-29 17:40:53

--
-- PostgreSQL database dump complete
--

\unrestrict yZDJ0rbbsQXI9RW9vJwhRRL8W11ZscdWCTCG655q4pA0kwD9ocpQwI9MU7ddeEQ

