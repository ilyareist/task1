--
-- PostgreSQL database dump
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
-- Name: accounts accounts_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.accounts
    ADD CONSTRAINT accounts_pkey PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
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
-- Name: payments payments_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT payments_pkey PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--


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
