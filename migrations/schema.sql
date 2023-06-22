-- schema.sql

CREATE TABLE issuers
(
    id              SERIAL primary key NOT null,
    company_name    VARCHAR(150)       NOT NULL,
    available_funds DECIMAL(10, 2)     not null DEFAULT 0
);

INSERT INTO issuers(company_name, available_funds)
VALUES ('Gigabyte LTA', 0),
       ('Amazon INC', 3000),
       ('Apple INC', 1000);

CREATE TABLE investors
(
    id              SERIAL primary key NOT null,
    name            VARCHAR(150)       NOT NULL,
    available_funds DECIMAL(10, 2) DEFAULT 0
);

INSERT INTO investors(name, available_funds)
VALUES ('John Doe', 0),
       ('Robert McKenzie', 1000000),
       ('Margaret Hill', 5000);

CREATE TABLE invoices
(
    id           VARCHAR(50)    NOT NULL,
    due_date     VARCHAR(150)   NOT NULL,
    asking_price DECIMAL(10, 2) NOT NULL,
    status       VARCHAR(150)   NOT NULL,
    items        JSONB          NOT NULL,
    created_at   VARCHAR(150)   NOT NULL,
    issuer_id    INT            NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (issuer_id) REFERENCES issuers (id) ON UPDATE cascade
);

CREATE TABLE bids
(
    id          VARCHAR(50)    not NULL,
    investor_id INT            not NULL,
    invoice_id  VARCHAR(50)    not NULL,
    bid_amount  DECIMAL(10, 2) not NULL,
    created_at  VARCHAR(150)   not null,
    primary key (id),
    FOREIGN KEY (investor_id) REFERENCES investors (id) ON UPDATE cascade,
    FOREIGN KEY (invoice_id) REFERENCES invoices (id) ON UPDATE cascade
);

CREATE TABLE trades
(
    id            VARCHAR(50)  not null,
    invoice_id    VARCHAR(50)  not NULL,
    investors_ids INT ARRAY not NULL,
    trade_status  VARCHAR(50)  not NULL,
    created_at    VARCHAR(150) not null,
    updated_at    VARCHAR(150),
    primary key (id),
    FOREIGN KEY (invoice_id) REFERENCES invoices (id)
);
