-- Only tb_user, tb_customer, tb_vehicle. We are creating the tables by migration.

CREATE TABLE IF NOT EXISTS tb_customer (
    id BIGSERIAL PRIMARY KEY,
    cpf_cnpj VARCHAR(20) UNIQUE NOT NULL,
    fullname VARCHAR(255) NOT NULL,
    phone_number VARCHAR(20),
    user_id INTEGER NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_customer_document ON tb_customer (cpf_cnpj);

INSERT INTO tb_customer (cpf_cnpj, fullname, phone_number, user_id)
    VALUES ('19134950869', 'João da Silva', '(11) 98765-4321', 1);

INSERT INTO tb_customer (cpf_cnpj, fullname, phone_number, user_id)
VALUES ('93684508000137', 'Joana da Silva', '(11) 99775-4829', 2);

CREATE TABLE IF NOT EXISTS tb_vehicle (
    id BIGSERIAL PRIMARY KEY,
    plate VARCHAR(20) UNIQUE NOT NULL,
    customer_id INTEGER NOT NULL,
    model VARCHAR(100) NOT NULL,
    brand VARCHAR(100) NOT NULL,
    year INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    FOREIGN KEY (customer_id) REFERENCES tb_customer(id) ON DELETE CASCADE
);


CREATE INDEX IF NOT EXISTS idx_vehicle_customer_id ON tb_vehicle (customer_id);

INSERT INTO tb_vehicle (plate, customer_id, model, brand, year)
    VALUES ('EXD6183', 1, 'Corsa', 'Chevrolet', 2005);

INSERT INTO tb_vehicle (plate, customer_id, model, brand, year)
VALUES ('AAA0X00', 2, 'Onix', 'Chevrolet', 2014);

CREATE TABLE IF NOT EXISTS tb_service (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create join table for many-to-many relationship between services and service orders
CREATE TABLE IF NOT EXISTS tb_service_service_orders (
    service_id INTEGER REFERENCES tb_service(id),
    service_order_id INTEGER,
    PRIMARY KEY (service_id, service_order_id)
);

-- Add some sample data
INSERT INTO tb_service (name, description, price)
VALUES
    ('Oil Change', 'Complete oil change service with filter replacement', 150.00),
    ('Brake Service', 'Brake pad replacement and system check', 300.00);
