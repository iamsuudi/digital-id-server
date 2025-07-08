--
-- PostgreSQL database dump
--

-- Dumped from database version 17.5
-- Dumped by pg_dump version 17.5

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

--
-- Name: DocumentStatus; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public."DocumentStatus" AS ENUM (
    'PENDING',
    'APPROVED',
    'REJECTED'
);


ALTER TYPE public."DocumentStatus" OWNER TO postgres;

--
-- Name: DocumentType; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public."DocumentType" AS ENUM (
    'ID',
    'PASSPORT',
    'DRIVING_LICENSE',
    'NATIONAL_ID'
);


ALTER TYPE public."DocumentType" OWNER TO postgres;

--
-- Name: Gender; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public."Gender" AS ENUM (
    'MALE',
    'FEMALE',
    'OTHER'
);


ALTER TYPE public."Gender" OWNER TO postgres;

--
-- Name: MaritalStatus; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public."MaritalStatus" AS ENUM (
    'SINGLE',
    'MARRIED',
    'DIVORCED',
    'WIDOWED'
);


ALTER TYPE public."MaritalStatus" OWNER TO postgres;

--
-- Name: PaymentMethod; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public."PaymentMethod" AS ENUM (
    'CASH',
    'MOBILE_MONEY',
    'BANK_TRANSFER'
);


ALTER TYPE public."PaymentMethod" OWNER TO postgres;

--
-- Name: PaymentStatus; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public."PaymentStatus" AS ENUM (
    'PENDING',
    'APPROVED',
    'REJECTED'
);


ALTER TYPE public."PaymentStatus" OWNER TO postgres;

--
-- Name: Religion; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public."Religion" AS ENUM (
    'CHRISTIAN',
    'MUSLIM',
    'HINDU',
    'BUDDHIST',
    'OTHER',
    'NONE'
);


ALTER TYPE public."Religion" OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: address; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.address (
    id integer NOT NULL,
    "houseNumber" text NOT NULL,
    district text NOT NULL,
    "cityId" integer NOT NULL,
    "createdAt" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "deletedAt" timestamp(3) without time zone,
    "searchVector" tsvector GENERATED ALWAYS AS (to_tsvector('english'::regconfig, (("houseNumber" || ' '::text) || district))) STORED
);


ALTER TABLE public.address OWNER TO postgres;

--
-- Name: address_cityId_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public."address_cityId_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public."address_cityId_seq" OWNER TO postgres;

--
-- Name: address_cityId_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public."address_cityId_seq" OWNED BY public.address."cityId";


--
-- Name: address_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.address_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.address_id_seq OWNER TO postgres;

--
-- Name: address_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.address_id_seq OWNED BY public.address.id;


--
-- Name: biometric; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.biometric (
    id integer NOT NULL,
    "residentId" integer NOT NULL,
    fingerprint bytea,
    "bloodType" text NOT NULL,
    face text NOT NULL,
    "createdAt" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "deletedAt" timestamp(3) without time zone,
    "searchVector" tsvector GENERATED ALWAYS AS (to_tsvector('english'::regconfig, "bloodType")) STORED
);


ALTER TABLE public.biometric OWNER TO postgres;

--
-- Name: biometric_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.biometric_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.biometric_id_seq OWNER TO postgres;

--
-- Name: biometric_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.biometric_id_seq OWNED BY public.biometric.id;


--
-- Name: biometric_residentId_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public."biometric_residentId_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public."biometric_residentId_seq" OWNER TO postgres;

--
-- Name: biometric_residentId_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public."biometric_residentId_seq" OWNED BY public.biometric."residentId";


--
-- Name: city; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.city (
    id integer NOT NULL,
    name text NOT NULL,
    "regionId" integer NOT NULL,
    "createdAt" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "deletedAt" timestamp(3) without time zone,
    "searchVector" tsvector GENERATED ALWAYS AS (to_tsvector('english'::regconfig, name)) STORED
);


ALTER TABLE public.city OWNER TO postgres;

--
-- Name: city_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.city_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.city_id_seq OWNER TO postgres;

--
-- Name: city_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.city_id_seq OWNED BY public.city.id;


--
-- Name: city_regionId_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public."city_regionId_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public."city_regionId_seq" OWNER TO postgres;

--
-- Name: city_regionId_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public."city_regionId_seq" OWNED BY public.city."regionId";


--
-- Name: document; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.document (
    id integer NOT NULL,
    type public."DocumentType" NOT NULL,
    "residentId" integer NOT NULL,
    url text NOT NULL,
    status public."DocumentStatus" NOT NULL,
    number text NOT NULL,
    "createdAt" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "deletedAt" timestamp(3) without time zone,
    "searchVector" tsvector GENERATED ALWAYS AS (to_tsvector('english'::regconfig, number)) STORED
);


ALTER TABLE public.document OWNER TO postgres;

--
-- Name: document_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.document_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.document_id_seq OWNER TO postgres;

--
-- Name: document_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.document_id_seq OWNED BY public.document.id;


--
-- Name: document_residentId_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public."document_residentId_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public."document_residentId_seq" OWNER TO postgres;

--
-- Name: document_residentId_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public."document_residentId_seq" OWNED BY public.document."residentId";


--
-- Name: emergency; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.emergency (
    id integer NOT NULL,
    "residentId" integer NOT NULL,
    name text NOT NULL,
    relation text NOT NULL,
    phone character varying(20) NOT NULL,
    "createdAt" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "deletedAt" timestamp(3) without time zone,
    "searchVector" tsvector GENERATED ALWAYS AS (to_tsvector('english'::regconfig, ((name || ' '::text) || relation))) STORED
);


ALTER TABLE public.emergency OWNER TO postgres;

--
-- Name: emergency_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.emergency_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.emergency_id_seq OWNER TO postgres;

--
-- Name: emergency_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.emergency_id_seq OWNED BY public.emergency.id;


--
-- Name: emergency_residentId_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public."emergency_residentId_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public."emergency_residentId_seq" OWNER TO postgres;

--
-- Name: emergency_residentId_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public."emergency_residentId_seq" OWNED BY public.emergency."residentId";


--
-- Name: employment; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.employment (
    id integer NOT NULL,
    "residentId" integer NOT NULL,
    status text NOT NULL,
    occupation text,
    "employerName" text,
    "workAddress" text,
    "createdAt" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "deletedAt" timestamp(3) without time zone,
    "searchVector" tsvector GENERATED ALWAYS AS (to_tsvector('english'::regconfig, ((((((status || ' '::text) || COALESCE(occupation, ''::text)) || ' '::text) || COALESCE("employerName", ''::text)) || ' '::text) || COALESCE("workAddress", ''::text)))) STORED
);


ALTER TABLE public.employment OWNER TO postgres;

--
-- Name: employment_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.employment_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.employment_id_seq OWNER TO postgres;

--
-- Name: employment_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.employment_id_seq OWNED BY public.employment.id;


--
-- Name: employment_residentId_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public."employment_residentId_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public."employment_residentId_seq" OWNER TO postgres;

--
-- Name: employment_residentId_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public."employment_residentId_seq" OWNED BY public.employment."residentId";


--
-- Name: idcard; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.idcard (
    id integer NOT NULL,
    "residentId" integer NOT NULL,
    number text NOT NULL,
    "issueDate" timestamp(3) without time zone NOT NULL,
    "expiryDate" timestamp(3) without time zone NOT NULL,
    "issuePlace" text NOT NULL,
    "createdAt" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "deletedAt" timestamp(3) without time zone,
    "searchVector" tsvector GENERATED ALWAYS AS (to_tsvector('english'::regconfig, ((number || ' '::text) || "issuePlace"))) STORED
);


ALTER TABLE public.idcard OWNER TO postgres;

--
-- Name: idcard_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.idcard_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.idcard_id_seq OWNER TO postgres;

--
-- Name: idcard_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.idcard_id_seq OWNED BY public.idcard.id;


--
-- Name: idcard_residentId_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public."idcard_residentId_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public."idcard_residentId_seq" OWNER TO postgres;

--
-- Name: idcard_residentId_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public."idcard_residentId_seq" OWNED BY public.idcard."residentId";


--
-- Name: payment; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.payment (
    id integer NOT NULL,
    "residentId" integer NOT NULL,
    amount double precision NOT NULL,
    description text NOT NULL,
    status public."PaymentStatus" NOT NULL,
    reference text NOT NULL,
    method public."PaymentMethod" NOT NULL,
    "createdAt" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "deletedAt" timestamp(3) without time zone,
    "searchVector" tsvector GENERATED ALWAYS AS (to_tsvector('english'::regconfig, ((description || ' '::text) || reference))) STORED,
    CONSTRAINT payment_amount_check CHECK ((amount >= (0)::double precision))
);


ALTER TABLE public.payment OWNER TO postgres;

--
-- Name: payment_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.payment_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.payment_id_seq OWNER TO postgres;

--
-- Name: payment_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.payment_id_seq OWNED BY public.payment.id;


--
-- Name: payment_residentId_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public."payment_residentId_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public."payment_residentId_seq" OWNER TO postgres;

--
-- Name: payment_residentId_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public."payment_residentId_seq" OWNED BY public.payment."residentId";


--
-- Name: region; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.region (
    id integer NOT NULL,
    name text NOT NULL,
    "createdAt" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "deletedAt" timestamp(3) without time zone,
    "searchVector" tsvector GENERATED ALWAYS AS (to_tsvector('english'::regconfig, name)) STORED
);


ALTER TABLE public.region OWNER TO postgres;

--
-- Name: region_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.region_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.region_id_seq OWNER TO postgres;

--
-- Name: region_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.region_id_seq OWNED BY public.region.id;


--
-- Name: resident; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.resident (
    id integer NOT NULL,
    email text NOT NULL,
    "firstName" text NOT NULL,
    "secondName" text NOT NULL,
    "lastName" text NOT NULL,
    "birthDate" timestamp(3) without time zone NOT NULL,
    gender public."Gender" NOT NULL,
    phone character varying(20) NOT NULL,
    "maritalStatus" public."MaritalStatus" NOT NULL,
    religion public."Religion" NOT NULL,
    ethnicity text,
    "disabilityStatus" text,
    "educationLevel" text,
    "languagesSpoken" text NOT NULL,
    "addressId" integer NOT NULL,
    "createdAt" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "deletedAt" timestamp(3) without time zone,
    "searchVector" tsvector GENERATED ALWAYS AS (to_tsvector('english'::regconfig, (((("firstName" || ' '::text) || "secondName") || ' '::text) || "lastName"))) STORED
);


ALTER TABLE public.resident OWNER TO postgres;

--
-- Name: resident_addressId_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public."resident_addressId_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public."resident_addressId_seq" OWNER TO postgres;

--
-- Name: resident_addressId_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public."resident_addressId_seq" OWNED BY public.resident."addressId";


--
-- Name: resident_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.resident_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.resident_id_seq OWNER TO postgres;

--
-- Name: resident_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.resident_id_seq OWNED BY public.resident.id;


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO postgres;

--
-- Name: address id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.address ALTER COLUMN id SET DEFAULT nextval('public.address_id_seq'::regclass);


--
-- Name: address cityId; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.address ALTER COLUMN "cityId" SET DEFAULT nextval('public."address_cityId_seq"'::regclass);


--
-- Name: biometric id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.biometric ALTER COLUMN id SET DEFAULT nextval('public.biometric_id_seq'::regclass);


--
-- Name: biometric residentId; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.biometric ALTER COLUMN "residentId" SET DEFAULT nextval('public."biometric_residentId_seq"'::regclass);


--
-- Name: city id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.city ALTER COLUMN id SET DEFAULT nextval('public.city_id_seq'::regclass);


--
-- Name: city regionId; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.city ALTER COLUMN "regionId" SET DEFAULT nextval('public."city_regionId_seq"'::regclass);


--
-- Name: document id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.document ALTER COLUMN id SET DEFAULT nextval('public.document_id_seq'::regclass);


--
-- Name: document residentId; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.document ALTER COLUMN "residentId" SET DEFAULT nextval('public."document_residentId_seq"'::regclass);


--
-- Name: emergency id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.emergency ALTER COLUMN id SET DEFAULT nextval('public.emergency_id_seq'::regclass);


--
-- Name: emergency residentId; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.emergency ALTER COLUMN "residentId" SET DEFAULT nextval('public."emergency_residentId_seq"'::regclass);


--
-- Name: employment id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.employment ALTER COLUMN id SET DEFAULT nextval('public.employment_id_seq'::regclass);


--
-- Name: employment residentId; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.employment ALTER COLUMN "residentId" SET DEFAULT nextval('public."employment_residentId_seq"'::regclass);


--
-- Name: idcard id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.idcard ALTER COLUMN id SET DEFAULT nextval('public.idcard_id_seq'::regclass);


--
-- Name: idcard residentId; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.idcard ALTER COLUMN "residentId" SET DEFAULT nextval('public."idcard_residentId_seq"'::regclass);


--
-- Name: payment id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment ALTER COLUMN id SET DEFAULT nextval('public.payment_id_seq'::regclass);


--
-- Name: payment residentId; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment ALTER COLUMN "residentId" SET DEFAULT nextval('public."payment_residentId_seq"'::regclass);


--
-- Name: region id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.region ALTER COLUMN id SET DEFAULT nextval('public.region_id_seq'::regclass);


--
-- Name: resident id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.resident ALTER COLUMN id SET DEFAULT nextval('public.resident_id_seq'::regclass);


--
-- Name: resident addressId; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.resident ALTER COLUMN "addressId" SET DEFAULT nextval('public."resident_addressId_seq"'::regclass);


--
-- Name: address address_houseNumber_district_cityId_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.address
    ADD CONSTRAINT "address_houseNumber_district_cityId_key" UNIQUE ("houseNumber", district, "cityId");


--
-- Name: address address_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.address
    ADD CONSTRAINT address_pkey PRIMARY KEY (id);


--
-- Name: biometric biometric_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.biometric
    ADD CONSTRAINT biometric_pkey PRIMARY KEY (id);


--
-- Name: biometric biometric_residentId_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.biometric
    ADD CONSTRAINT "biometric_residentId_key" UNIQUE ("residentId");


--
-- Name: city city_name_regionId_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.city
    ADD CONSTRAINT "city_name_regionId_key" UNIQUE (name, "regionId");


--
-- Name: city city_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.city
    ADD CONSTRAINT city_pkey PRIMARY KEY (id);


--
-- Name: document document_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.document
    ADD CONSTRAINT document_pkey PRIMARY KEY (id);


--
-- Name: document document_residentId_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.document
    ADD CONSTRAINT "document_residentId_key" UNIQUE ("residentId");


--
-- Name: emergency emergency_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.emergency
    ADD CONSTRAINT emergency_pkey PRIMARY KEY (id);


--
-- Name: emergency emergency_residentId_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.emergency
    ADD CONSTRAINT "emergency_residentId_key" UNIQUE ("residentId");


--
-- Name: employment employment_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.employment
    ADD CONSTRAINT employment_pkey PRIMARY KEY (id);


--
-- Name: employment employment_residentId_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.employment
    ADD CONSTRAINT "employment_residentId_key" UNIQUE ("residentId");


--
-- Name: idcard idcard_number_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.idcard
    ADD CONSTRAINT idcard_number_key UNIQUE (number);


--
-- Name: idcard idcard_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.idcard
    ADD CONSTRAINT idcard_pkey PRIMARY KEY (id);


--
-- Name: idcard idcard_residentId_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.idcard
    ADD CONSTRAINT "idcard_residentId_key" UNIQUE ("residentId");


--
-- Name: payment payment_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment
    ADD CONSTRAINT payment_pkey PRIMARY KEY (id);


--
-- Name: payment payment_residentId_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment
    ADD CONSTRAINT "payment_residentId_key" UNIQUE ("residentId");


--
-- Name: region region_name_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.region
    ADD CONSTRAINT region_name_key UNIQUE (name);


--
-- Name: region region_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.region
    ADD CONSTRAINT region_pkey PRIMARY KEY (id);


--
-- Name: resident resident_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.resident
    ADD CONSTRAINT resident_email_key UNIQUE (email);


--
-- Name: resident resident_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.resident
    ADD CONSTRAINT resident_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: address_search_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX address_search_idx ON public.address USING gin ("searchVector");


--
-- Name: biometric_search_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX biometric_search_idx ON public.biometric USING gin ("searchVector");


--
-- Name: city_search_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX city_search_idx ON public.city USING gin ("searchVector");


--
-- Name: document_search_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX document_search_idx ON public.document USING gin ("searchVector");


--
-- Name: emergency_search_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX emergency_search_idx ON public.emergency USING gin ("searchVector");


--
-- Name: employment_search_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX employment_search_idx ON public.employment USING gin ("searchVector");


--
-- Name: idcard_search_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idcard_search_idx ON public.idcard USING gin ("searchVector");


--
-- Name: payment_search_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX payment_search_idx ON public.payment USING gin ("searchVector");


--
-- Name: region_search_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX region_search_idx ON public.region USING gin ("searchVector");


--
-- Name: resident_search_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX resident_search_idx ON public.resident USING gin ("searchVector");


--
-- Name: address address_cityId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.address
    ADD CONSTRAINT "address_cityId_fkey" FOREIGN KEY ("cityId") REFERENCES public.city(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: biometric biometric_residentId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.biometric
    ADD CONSTRAINT "biometric_residentId_fkey" FOREIGN KEY ("residentId") REFERENCES public.resident(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: city city_regionId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.city
    ADD CONSTRAINT "city_regionId_fkey" FOREIGN KEY ("regionId") REFERENCES public.region(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: document document_residentId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.document
    ADD CONSTRAINT "document_residentId_fkey" FOREIGN KEY ("residentId") REFERENCES public.resident(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: emergency emergency_residentId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.emergency
    ADD CONSTRAINT "emergency_residentId_fkey" FOREIGN KEY ("residentId") REFERENCES public.resident(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: employment employment_residentId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.employment
    ADD CONSTRAINT "employment_residentId_fkey" FOREIGN KEY ("residentId") REFERENCES public.resident(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: idcard idcard_residentId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.idcard
    ADD CONSTRAINT "idcard_residentId_fkey" FOREIGN KEY ("residentId") REFERENCES public.resident(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: payment payment_residentId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment
    ADD CONSTRAINT "payment_residentId_fkey" FOREIGN KEY ("residentId") REFERENCES public.resident(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: resident resident_addressId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.resident
    ADD CONSTRAINT "resident_addressId_fkey" FOREIGN KEY ("addressId") REFERENCES public.address(id) ON UPDATE CASCADE ON DELETE SET NULL;


--
-- PostgreSQL database dump complete
--

