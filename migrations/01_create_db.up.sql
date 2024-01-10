CREATE TABLE
    users (
        user_id uuid PRIMARY KEY,
        email varchar(50) NOT NULL UNIQUE,
        password varchar(100) NOT NULL,
        name varchar(50) NOT NULL,
        funds int NOT NULL,
        role varchar(10) NOT NULL
    );

CREATE TABLE
    invoices (
        invoice_id uuid PRIMARY KEY,
        issuer_pk uuid,
        name varchar(50) NOT NULL,
        price int NOT NULL,
        funds int,
        status VARCHAR(10),
        expire_date date,
        FOREIGN KEY (issuer_pk) REFERENCES users(user_id)
    );

CREATE TABLE
    invoice_records (
        invoice_pk uuid,
        investor_pk uuid,
        invested_funds INT NOT NULL,
        FOREIGN KEY (invoice_pk) REFERENCES invoices(invoice_id),
        FOREIGN KEY (investor_pk) REFERENCES users(user_id)
    );