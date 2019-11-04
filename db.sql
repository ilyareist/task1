--
-- PostgreSQL database dump
--

-- Dumped from database version 9.6.15
-- Dumped by pg_dump version 9.6.15

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
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner:
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner:
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: accounts; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.accounts (
    id character varying(255) NOT NULL,
    country character varying(50) NOT NULL,
    city character varying(50) NOT NULL,
    balance numeric(16,4) NOT NULL,
    currency character varying(3) NOT NULL,
    deleted boolean NOT NULL
);


ALTER TABLE public.accounts OWNER TO postgres;

--
-- Name: payments; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.payments (
    id character varying(36) NOT NULL,
    account character varying(255),
    amount numeric(16,4) NOT NULL,
    to_account character varying(255),
    from_account character varying(255),
    direction character varying(16) NOT NULL,
    deleted boolean NOT NULL
);


ALTER TABLE public.payments OWNER TO postgres;

--
-- Name: accounts accounts_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.accounts
    ADD CONSTRAINT accounts_pkey PRIMARY KEY (id);


--
-- Name: payments payments_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT payments_pkey PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--




CREATE INDEX payments_account_direction_index ON payments (account, direction);

CREATE OR REPLACE VIEW accounts_view AS
SELECT A.id,
       A.balance
           + (SELECT COALESCE(SUM(P.amount), 0)
              FROM payments AS P
              WHERE P.account = A.id
              AND P.direction='incoming')
           - (SELECT COALESCE(SUM(P.amount), 0)
              FROM payments AS P
              WHERE P.account = A.id
              AND P.direction='outgoing')
       AS balance,
       A.country,
       A.city,
       A.currency,
       A.deleted
FROM accounts AS A;
