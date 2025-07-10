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
-- Name: document_status; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.document_status AS ENUM (
    'PENDING',
    'APPROVED',
    'REJECTED'
);


ALTER TYPE public.document_status OWNER TO postgres;

--
-- Name: document_type; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.document_type AS ENUM (
    'ID',
    'PASSPORT',
    'DRIVING_LICENSE',
    'NATIONAL_ID'
);


ALTER TYPE public.document_type OWNER TO postgres;

--
-- Name: gender; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.gender AS ENUM (
    'MALE',
    'FEMALE',
    'OTHER'
);


ALTER TYPE public.gender OWNER TO postgres;

--
-- Name: marital_status; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.marital_status AS ENUM (
    'SINGLE',
    'MARRIED',
    'DIVORCED',
    'WIDOWED'
);


ALTER TYPE public.marital_status OWNER TO postgres;

--
-- Name: payment_method; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.payment_method AS ENUM (
    'CASH',
    'MOBILE_MONEY',
    'BANK_TRANSFER'
);


ALTER TYPE public.payment_method OWNER TO postgres;

--
-- Name: payment_status; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.payment_status AS ENUM (
    'PENDING',
    'APPROVED',
    'REJECTED'
);


ALTER TYPE public.payment_status OWNER TO postgres;

--
-- Name: religion; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.religion AS ENUM (
    'CHRISTIAN',
    'MUSLIM',
    'HINDU',
    'BUDDHIST',
    'OTHER',
    'NONE'
);


ALTER TYPE public.religion OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: address; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.address (
    id uuid NOT NULL,
    house_number text NOT NULL,
    district text NOT NULL,
    city_id uuid NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp without time zone,
    search_vector tsvector GENERATED ALWAYS AS (to_tsvector('english'::regconfig, ((house_number || ' '::text) || district))) STORED
);


ALTER TABLE public.address OWNER TO postgres;

--
-- Name: biometric; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.biometric (
    id uuid NOT NULL,
    resident_id uuid NOT NULL,
    fingerprint bytea,
    blood_type text NOT NULL,
    face text NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp without time zone,
    search_vector tsvector GENERATED ALWAYS AS (to_tsvector('english'::regconfig, blood_type)) STORED
);


ALTER TABLE public.biometric OWNER TO postgres;

--
-- Name: city; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.city (
    id uuid NOT NULL,
    name text NOT NULL,
    region_id uuid NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp without time zone,
    search_vector tsvector GENERATED ALWAYS AS (to_tsvector('english'::regconfig, name)) STORED
);


ALTER TABLE public.city OWNER TO postgres;

--
-- Name: document; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.document (
    id uuid NOT NULL,
    type public.document_type NOT NULL,
    resident_id uuid NOT NULL,
    url text NOT NULL,
    status public.document_status NOT NULL,
    number text NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp without time zone,
    search_vector tsvector GENERATED ALWAYS AS (to_tsvector('english'::regconfig, number)) STORED
);


ALTER TABLE public.document OWNER TO postgres;

--
-- Name: emergency; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.emergency (
    id uuid NOT NULL,
    resident_id uuid NOT NULL,
    name text NOT NULL,
    relation text NOT NULL,
    phone character varying(20) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp without time zone,
    search_vector tsvector GENERATED ALWAYS AS (to_tsvector('english'::regconfig, ((name || ' '::text) || relation))) STORED
);


ALTER TABLE public.emergency OWNER TO postgres;

--
-- Name: employment; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.employment (
    id uuid NOT NULL,
    resident_id uuid NOT NULL,
    status text NOT NULL,
    occupation text,
    employer_name text,
    work_address text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp without time zone,
    search_vector tsvector GENERATED ALWAYS AS (to_tsvector('english'::regconfig, ((((((status || ' '::text) || COALESCE(occupation, ''::text)) || ' '::text) || COALESCE(employer_name, ''::text)) || ' '::text) || COALESCE(work_address, ''::text)))) STORED
);


ALTER TABLE public.employment OWNER TO postgres;

--
-- Name: idcard; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.idcard (
    id uuid NOT NULL,
    resident_id uuid NOT NULL,
    number text NOT NULL,
    issue_date timestamp without time zone NOT NULL,
    expiry_date timestamp without time zone NOT NULL,
    issue_place text NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp without time zone,
    search_vector tsvector GENERATED ALWAYS AS (to_tsvector('english'::regconfig, ((number || ' '::text) || issue_place))) STORED
);


ALTER TABLE public.idcard OWNER TO postgres;

--
-- Name: payment; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.payment (
    id uuid NOT NULL,
    resident_id uuid NOT NULL,
    amount numeric NOT NULL,
    description text NOT NULL,
    status public.payment_status NOT NULL,
    reference text NOT NULL,
    method public.payment_method NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp without time zone,
    search_vector tsvector GENERATED ALWAYS AS (to_tsvector('english'::regconfig, ((description || ' '::text) || reference))) STORED,
    CONSTRAINT payment_amount_check CHECK ((amount >= (0)::numeric))
);


ALTER TABLE public.payment OWNER TO postgres;

--
-- Name: refresh_tokens; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.refresh_tokens (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    token text NOT NULL,
    expires_at timestamp without time zone NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.refresh_tokens OWNER TO postgres;

--
-- Name: region; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.region (
    id uuid NOT NULL,
    name text NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp without time zone,
    search_vector tsvector GENERATED ALWAYS AS (to_tsvector('english'::regconfig, name)) STORED
);


ALTER TABLE public.region OWNER TO postgres;

--
-- Name: resident; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.resident (
    id uuid NOT NULL,
    email text NOT NULL,
    first_name text NOT NULL,
    second_name text NOT NULL,
    last_name text NOT NULL,
    birth_date timestamp without time zone NOT NULL,
    gender public.gender NOT NULL,
    phone character varying(20) NOT NULL,
    marital_status public.marital_status NOT NULL,
    religion public.religion NOT NULL,
    ethnicity text,
    disability_status text,
    education_level text,
    languages_spoken text NOT NULL,
    address_id uuid,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp without time zone,
    search_vector tsvector GENERATED ALWAYS AS (to_tsvector('english'::regconfig, ((((first_name || ' '::text) || second_name) || ' '::text) || last_name))) STORED,
    CONSTRAINT resident_email_check CHECK ((email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'::text)),
    CONSTRAINT resident_phone_check CHECK (((phone)::text ~ '^\+?[1-9]\d{1,14}$'::text))
);


ALTER TABLE public.resident OWNER TO postgres;

--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO postgres;

--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id uuid NOT NULL,
    first_name text NOT NULL,
    second_name text NOT NULL,
    last_name text NOT NULL,
    email text NOT NULL,
    phone text NOT NULL,
    password text NOT NULL,
    role text NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp without time zone,
    CONSTRAINT users_email_check CHECK ((email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'::text)),
    CONSTRAINT users_phone_check CHECK ((phone ~ '^\+?[1-9]\d{1,14}$'::text)),
    CONSTRAINT users_role_check CHECK ((role = ANY (ARRAY['SUPERADMIN'::text, 'MANAGER'::text, 'ENCODER'::text, 'ADMIN'::text, 'CASHIER'::text])))
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: address address_house_number_district_city_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.address
    ADD CONSTRAINT address_house_number_district_city_id_key UNIQUE (house_number, district, city_id);


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
-- Name: biometric biometric_resident_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.biometric
    ADD CONSTRAINT biometric_resident_id_key UNIQUE (resident_id);


--
-- Name: city city_name_region_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.city
    ADD CONSTRAINT city_name_region_id_key UNIQUE (name, region_id);


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
-- Name: document document_resident_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.document
    ADD CONSTRAINT document_resident_id_key UNIQUE (resident_id);


--
-- Name: emergency emergency_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.emergency
    ADD CONSTRAINT emergency_pkey PRIMARY KEY (id);


--
-- Name: emergency emergency_resident_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.emergency
    ADD CONSTRAINT emergency_resident_id_key UNIQUE (resident_id);


--
-- Name: employment employment_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.employment
    ADD CONSTRAINT employment_pkey PRIMARY KEY (id);


--
-- Name: employment employment_resident_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.employment
    ADD CONSTRAINT employment_resident_id_key UNIQUE (resident_id);


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
-- Name: idcard idcard_resident_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.idcard
    ADD CONSTRAINT idcard_resident_id_key UNIQUE (resident_id);


--
-- Name: payment payment_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment
    ADD CONSTRAINT payment_pkey PRIMARY KEY (id);


--
-- Name: payment payment_resident_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment
    ADD CONSTRAINT payment_resident_id_key UNIQUE (resident_id);


--
-- Name: refresh_tokens refresh_tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT refresh_tokens_pkey PRIMARY KEY (id);


--
-- Name: refresh_tokens refresh_tokens_token_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT refresh_tokens_token_key UNIQUE (token);


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
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: address_search_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX address_search_idx ON public.address USING gin (search_vector);


--
-- Name: biometric_search_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX biometric_search_idx ON public.biometric USING gin (search_vector);


--
-- Name: city_search_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX city_search_idx ON public.city USING gin (search_vector);


--
-- Name: document_search_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX document_search_idx ON public.document USING gin (search_vector);


--
-- Name: emergency_search_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX emergency_search_idx ON public.emergency USING gin (search_vector);


--
-- Name: employment_search_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX employment_search_idx ON public.employment USING gin (search_vector);


--
-- Name: idcard_search_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idcard_search_idx ON public.idcard USING gin (search_vector);


--
-- Name: payment_search_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX payment_search_idx ON public.payment USING gin (search_vector);


--
-- Name: region_search_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX region_search_idx ON public.region USING gin (search_vector);


--
-- Name: resident_search_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX resident_search_idx ON public.resident USING gin (search_vector);


--
-- Name: address address_city_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.address
    ADD CONSTRAINT address_city_id_fkey FOREIGN KEY (city_id) REFERENCES public.city(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: biometric biometric_resident_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.biometric
    ADD CONSTRAINT biometric_resident_id_fkey FOREIGN KEY (resident_id) REFERENCES public.resident(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: city city_region_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.city
    ADD CONSTRAINT city_region_id_fkey FOREIGN KEY (region_id) REFERENCES public.region(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: document document_resident_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.document
    ADD CONSTRAINT document_resident_id_fkey FOREIGN KEY (resident_id) REFERENCES public.resident(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: emergency emergency_resident_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.emergency
    ADD CONSTRAINT emergency_resident_id_fkey FOREIGN KEY (resident_id) REFERENCES public.resident(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: employment employment_resident_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.employment
    ADD CONSTRAINT employment_resident_id_fkey FOREIGN KEY (resident_id) REFERENCES public.resident(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: idcard idcard_resident_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.idcard
    ADD CONSTRAINT idcard_resident_id_fkey FOREIGN KEY (resident_id) REFERENCES public.resident(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: payment payment_resident_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment
    ADD CONSTRAINT payment_resident_id_fkey FOREIGN KEY (resident_id) REFERENCES public.resident(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: refresh_tokens refresh_tokens_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT refresh_tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: resident resident_address_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.resident
    ADD CONSTRAINT resident_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.address(id) ON UPDATE CASCADE ON DELETE SET NULL;


--
-- PostgreSQL database dump complete
--

