DROP TABLE IF EXISTS additional_repair_dtos CASCADE;
DROP TABLE IF EXISTS parts_supply_service_order_dtos CASCADE;
DROP TABLE IF EXISTS service_service_order_dtos CASCADE;
DROP TABLE IF EXISTS payment_dtos CASCADE;
DROP TABLE IF EXISTS service_order_dtos CASCADE;
DROP TABLE IF EXISTS service_order_status_dtos CASCADE;
DROP TABLE IF EXISTS service_dtos CASCADE;
DROP TABLE IF EXISTS parts_supply_dtos CASCADE;
DROP TABLE IF EXISTS tb_vehicle CASCADE;
DROP TABLE IF EXISTS tb_customer CASCADE;
DROP TABLE IF EXISTS user_dtos CASCADE;
DROP TABLE IF EXISTS user_type_dtos CASCADE;
DROP TABLE if exists additional_repair_status_dtos CASCADE;

-- Base tables without dependencies

-- public.additional_repair_status_dtos definition

-- Drop table

-- DROP TABLE additional_repair_status_dtos;

CREATE TABLE additional_repair_status_dtos (
                                               id bigserial NOT NULL,
                                               description varchar(50) NOT NULL,
                                               CONSTRAINT additional_repair_status_dtos_pkey PRIMARY KEY (id)
);

-- public.user_type_dtos definition
CREATE TABLE user_type_dtos (
                                id bigserial NOT NULL,
                                "type" varchar(50) NOT NULL,
                                CONSTRAINT user_type_dtos_pkey PRIMARY KEY (id)
);

-- public.user_dtos definition
CREATE TABLE user_dtos (
                           id bigserial NOT NULL,
                           email varchar(100) NOT NULL,
                           "password" varchar(255) NOT NULL,
                           user_type_id int8 NOT NULL,
                           created_at timestamptz NULL,
                           updated_at timestamptz NULL,
                           deleted_at timestamptz NULL,
                           CONSTRAINT uni_user_dtos_email UNIQUE (email),
                           CONSTRAINT user_dtos_pkey PRIMARY KEY (id),
                           CONSTRAINT fk_user_type_dtos_users FOREIGN KEY (user_type_id) REFERENCES user_type_dtos(id)
);

-- public.service_order_status_dtos definition
CREATE TABLE service_order_status_dtos (
                                           id bigserial NOT NULL,
                                           description varchar(50) NOT NULL,
                                           CONSTRAINT service_order_status_dtos_pkey PRIMARY KEY (id)
);

-- public.parts_supply_dtos definition

-- Drop table

-- DROP TABLE parts_supply_dtos;

CREATE TABLE parts_supply_dtos (
                                   id bigserial NOT NULL,
                                   name varchar(100) NOT NULL,
                                   description text NULL,
                                   price numeric(10, 2) NOT NULL,
                                   CONSTRAINT parts_supply_dtos_pkey PRIMARY KEY (id)
);

-- public.tb_customer definition with user reference

CREATE TABLE tb_customer (
                             id BIGSERIAL PRIMARY KEY,
                             cpf_cnpj VARCHAR(20) UNIQUE NOT NULL,
                             fullname VARCHAR(255) NOT NULL,
                             phone_number VARCHAR(20),
                             user_id INTEGER NOT NULL,
                             FOREIGN KEY (user_id) REFERENCES user_dtos(id)
);

-- public.tb_vehicle definition

CREATE TABLE tb_vehicle (
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

-- public.service_dtos definition

CREATE TABLE service_dtos (
                              id BIGSERIAL PRIMARY KEY,
                              name VARCHAR(100) NOT NULL,
                              description TEXT,
                              price DECIMAL(10, 2) NOT NULL,
                              created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                              updated_at TIMESTAMP WITH TIME ZONE,
                              deleted_at TIMESTAMP WITH TIME ZONE
);

-- public.service_order_dtos definition

CREATE TABLE service_order_dtos (
                                    id bigserial NOT NULL,
                                    customer_id int8 NOT NULL,
                                    vehicle_id int8 NOT NULL,
                                    os_status_id int8 NOT NULL,
                                    estimate numeric(10, 2) NULL,
                                    started_execution_date timestamptz NULL,
                                    final_execution_date timestamptz NULL,
                                    created_at timestamptz DEFAULT NOW(),
                                    updated_at timestamptz NULL,
                                    deleted_at timestamptz NULL,
                                    CONSTRAINT service_order_dtos_pkey PRIMARY KEY (id),
                                    CONSTRAINT fk_customer_service_order FOREIGN KEY (customer_id) REFERENCES tb_customer(id),
                                    CONSTRAINT fk_vehicle_service_order FOREIGN KEY (vehicle_id) REFERENCES tb_vehicle(id),
                                    CONSTRAINT fk_status_service_order FOREIGN KEY (os_status_id) REFERENCES service_order_status_dtos(id)
);

-- public.service_service_order_dtos definition (join table for services and service orders)
CREATE TABLE service_service_order_dtos (
                                            service_id int8 NOT NULL,
                                            service_order_id int8 NOT NULL,
                                            price numeric(10, 2) NOT NULL,
                                            CONSTRAINT service_service_order_dtos_pkey PRIMARY KEY (service_id, service_order_id),
                                            CONSTRAINT fk_service_service_order FOREIGN KEY (service_id) REFERENCES service_dtos(id),
                                            CONSTRAINT fk_service_order_services FOREIGN KEY (service_order_id) REFERENCES service_order_dtos(id)
);

-- public.parts_supply_service_order_dtos definition
CREATE TABLE parts_supply_service_order_dtos (
                                                 parts_supply_id int8 NOT NULL,
                                                 service_order_id int8 NOT NULL,
                                                 price numeric(10, 2) NOT NULL,
                                                 quantity int NOT NULL DEFAULT 1,
                                                 CONSTRAINT parts_supply_service_order_dtos_pkey PRIMARY KEY (parts_supply_id, service_order_id),
                                                 CONSTRAINT fk_parts_supply_service_order FOREIGN KEY (parts_supply_id) REFERENCES parts_supply_dtos(id),
                                                 CONSTRAINT fk_service_order_parts_supply FOREIGN KEY (service_order_id) REFERENCES service_order_dtos(id)
);

-- public.payment_dtos definition
CREATE TABLE payment_dtos (
                              id bigserial NOT NULL,
                              service_order_id int8 NOT NULL,
                              payment_date timestamptz NOT NULL,
                              amount numeric(10, 2) NOT NULL,
                              CONSTRAINT payment_dtos_pkey PRIMARY KEY (id),
                              CONSTRAINT uni_payment_dtos_service_order_id UNIQUE (service_order_id),
                              CONSTRAINT fk_service_order_payment FOREIGN KEY (service_order_id) REFERENCES service_order_dtos(id)
);

-- public.additional_repair_dtos definition

-- Drop table

-- DROP TABLE additional_repair_dtos;

CREATE TABLE additional_repair_dtos (
                                        id bigserial NOT NULL,
                                        service_order_id int8 NOT NULL,
                                        service_id int8 NOT NULL,
                                        parts_supply_id int8 NOT NULL,
                                        ar_status_id int8 NOT NULL,
                                        CONSTRAINT additional_repair_dtos_pkey PRIMARY KEY (id),
                                        CONSTRAINT fk_additional_repair_status_dtos_additional_repairs FOREIGN KEY (ar_status_id) REFERENCES additional_repair_status_dtos(id),
                                        CONSTRAINT fk_parts_supply_dtos_additional_repairs FOREIGN KEY (parts_supply_id) REFERENCES parts_supply_dtos(id),
                                        CONSTRAINT fk_service_dtos_additional_repairs FOREIGN KEY (service_id) REFERENCES service_dtos(id),
                                        CONSTRAINT fk_service_order_dtos_additional_repairs FOREIGN KEY (service_order_id) REFERENCES service_order_dtos(id)
);

-- Create indexes
CREATE INDEX idx_service_order_dtos_customer_id ON service_order_dtos(customer_id);
CREATE INDEX idx_service_order_dtos_vehicle_id ON service_order_dtos(vehicle_id);
CREATE INDEX idx_service_order_dtos_status_id ON service_order_dtos(os_status_id);
CREATE INDEX idx_service_order_dtos_deleted_at ON service_order_dtos(deleted_at);

-- Insert sample data
INSERT INTO user_type_dtos ("type")
VALUES
    ('ADMIN'),
    ('CUSTOMER'),
    ('MECHANIC');

INSERT INTO user_dtos (email, password, user_type_id, created_at)
VALUES
    ('admin@xpto.com', '$2a$10$YourHashedPasswordHere', 1, NOW()),
    ('joao@email.com', '$2a$10$YourHashedPasswordHere', 2, NOW()),
    ('joana@email.com', '$2a$10$YourHashedPasswordHere', 2, NOW());

INSERT INTO tb_customer (cpf_cnpj, fullname, phone_number, user_id)
VALUES
    ('19134950869', 'João da Silva', '(11) 98765-4321', 1),
    ('93684508000137', 'Joana da Silva', '(11) 99775-4829', 2);

INSERT INTO tb_vehicle (plate, customer_id, model, brand, year)
VALUES
    ('EXD6183', 1, 'Corsa', 'Chevrolet', 2005),
    ('AAA0X00', 2, 'Onix', 'Chevrolet', 2014);

INSERT INTO service_dtos (name, description, price)
VALUES
    ('Oil Change', 'Complete oil change service with filter replacement', 150.00),
    ('Brake Service', 'Brake pad replacement and system check', 300.00);

INSERT INTO service_order_status_dtos (description)
VALUES
    ('Pending'),
    ('In Progress'),
    ('Completed'),
    ('Cancelled');