--
-- PostgreSQL database dump
--

-- Dumped from database version 12.1
-- Dumped by pg_dump version 12.1

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

DROP DATABASE IF EXISTS apartments;
--
-- Name: apartments; Type: DATABASE; Schema: -; Owner: sschiz
--

CREATE DATABASE apartments WITH TEMPLATE = template0 ENCODING = 'UTF8' LC_COLLATE = 'en_US.UTF-8' LC_CTYPE = 'en_US.UTF-8';


ALTER DATABASE apartments OWNER TO sschiz;

\connect apartments

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

--
-- Name: apartments; Type: SCHEMA; Schema: -; Owner: sschiz
--

CREATE SCHEMA apartments;


ALTER SCHEMA apartments OWNER TO sschiz;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: apartment_complexes; Type: TABLE; Schema: public; Owner: sschiz
--

CREATE TABLE public.apartment_complexes (
    name name NOT NULL,
    min_apartment_number integer DEFAULT 1 NOT NULL,
    max_apartment_number integer DEFAULT 1 NOT NULL
);


ALTER TABLE public.apartment_complexes OWNER TO sschiz;

--
-- Name: apartments; Type: TABLE; Schema: public; Owner: sschiz
--

CREATE TABLE public.apartments (
    id integer NOT NULL,
    rooms integer DEFAULT 1 NOT NULL,
    area double precision NOT NULL,
    rent money NOT NULL,
    house_id integer NOT NULL,
    floor integer DEFAULT 1 NOT NULL
);


ALTER TABLE public.apartments OWNER TO sschiz;

--
-- Name: apartments_id_seq; Type: SEQUENCE; Schema: public; Owner: sschiz
--

CREATE SEQUENCE public.apartments_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.apartments_id_seq OWNER TO sschiz;

--
-- Name: apartments_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: sschiz
--

ALTER SEQUENCE public.apartments_id_seq OWNED BY public.apartments.id;


--
-- Name: houses; Type: TABLE; Schema: public; Owner: sschiz
--

CREATE TABLE public.houses (
    id integer NOT NULL,
    city text NOT NULL,
    district text NOT NULL,
    address text NOT NULL,
    corpus text DEFAULT ''::text NOT NULL,
    floors integer DEFAULT 1 NOT NULL,
    ac_name name
);


ALTER TABLE public.houses OWNER TO sschiz;

--
-- Name: houses_id_seq; Type: SEQUENCE; Schema: public; Owner: sschiz
--

CREATE SEQUENCE public.houses_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.houses_id_seq OWNER TO sschiz;

--
-- Name: houses_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: sschiz
--

ALTER SEQUENCE public.houses_id_seq OWNED BY public.houses.id;


--
-- Name: apartments id; Type: DEFAULT; Schema: public; Owner: sschiz
--

ALTER TABLE ONLY public.apartments ALTER COLUMN id SET DEFAULT nextval('public.apartments_id_seq'::regclass);


--
-- Name: houses id; Type: DEFAULT; Schema: public; Owner: sschiz
--

ALTER TABLE ONLY public.houses ALTER COLUMN id SET DEFAULT nextval('public.houses_id_seq'::regclass);


--
-- Data for Name: apartment_complexes; Type: TABLE DATA; Schema: public; Owner: sschiz
--

COPY public.apartment_complexes (name, min_apartment_number, max_apartment_number) FROM stdin;
test	1	20
\.


--
-- Data for Name: apartments; Type: TABLE DATA; Schema: public; Owner: sschiz
--

COPY public.apartments (id, rooms, area, rent, house_id, floor) FROM stdin;
9	2	100	$200.21	17	1
10	10	100.1	$20,000.10	18	1
\.


--
-- Data for Name: houses; Type: TABLE DATA; Schema: public; Owner: sschiz
--

COPY public.houses (id, city, district, address, corpus, floors, ac_name) FROM stdin;
17	nchk	ivanovo	sov	a	10	test
18	nchk	ivanovo	zh. krurovoi, 10, 18		1	\N
\.


--
-- Name: apartments_id_seq; Type: SEQUENCE SET; Schema: public; Owner: sschiz
--

SELECT pg_catalog.setval('public.apartments_id_seq', 10, true);


--
-- Name: houses_id_seq; Type: SEQUENCE SET; Schema: public; Owner: sschiz
--

SELECT pg_catalog.setval('public.houses_id_seq', 18, true);


--
-- Name: apartment_complexes apartment_complexes_pk; Type: CONSTRAINT; Schema: public; Owner: sschiz
--

ALTER TABLE ONLY public.apartment_complexes
    ADD CONSTRAINT apartment_complexes_pk PRIMARY KEY (name);


--
-- Name: apartments apartments_pk; Type: CONSTRAINT; Schema: public; Owner: sschiz
--

ALTER TABLE ONLY public.apartments
    ADD CONSTRAINT apartments_pk PRIMARY KEY (id);


--
-- Name: houses houses_pk; Type: CONSTRAINT; Schema: public; Owner: sschiz
--

ALTER TABLE ONLY public.houses
    ADD CONSTRAINT houses_pk PRIMARY KEY (id);


--
-- Name: apartment_complexes_name_uindex; Type: INDEX; Schema: public; Owner: sschiz
--

CREATE UNIQUE INDEX apartment_complexes_name_uindex ON public.apartment_complexes USING btree (name);


--
-- Name: apartments apartments_houses_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: sschiz
--

ALTER TABLE ONLY public.apartments
    ADD CONSTRAINT apartments_houses_id_fk FOREIGN KEY (house_id) REFERENCES public.houses(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: houses houses_apartment_complexes_name_fk; Type: FK CONSTRAINT; Schema: public; Owner: sschiz
--

ALTER TABLE ONLY public.houses
    ADD CONSTRAINT houses_apartment_complexes_name_fk FOREIGN KEY (ac_name) REFERENCES public.apartment_complexes(name) ON UPDATE CASCADE ON DELETE SET NULL;


--
-- PostgreSQL database dump complete
--

