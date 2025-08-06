CREATE SCHEMA db_mecanica_xpto;

CREATE TABLE IF NOT EXISTS db_mecanica_xpto.tb_customer (
    id BIGSERIAL PRIMARY KEY,
    cpf_cnpj VARCHAR(20) UNIQUE NOT NULL,
    fullname VARCHAR(255) NOT NULL,
    phone_number VARCHAR(20),
    user_id INTEGER NOT NULL
);

-- Adiciona um índice na coluna document para otimizar buscas
CREATE INDEX IF NOT EXISTS idx_customer_document ON db_mecanica_xpto.tb_customer (cpf_cnpj);

-- Exemplo de inserção de dados na tabela
INSERT INTO db_mecanica_xpto.tb_customer (cpf_cnpj, fullname, phone_number, user_id)
    VALUES ('123.456.789-00', 'João da Silva', '(11) 98765-4321', 1);

CREATE TABLE IF NOT EXISTS db_mecanica_xpto.tb_vehicle (
    id BIGSERIAL PRIMARY KEY,
    plate VARCHAR(20) UNIQUE NOT NULL,
    customer_id INTEGER NOT NULL,
    model VARCHAR(100) NOT NULL,
    brand VARCHAR(100) NOT NULL,
    year INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    FOREIGN KEY (customer_id) REFERENCES db_mecanica_xpto.tb_customer(id) ON DELETE CASCADE
);

-- Adiciona um índice na coluna customer_id para otimizar buscas
-- e consultas de JOIN.
CREATE INDEX IF NOT EXISTS idx_vehicle_customer_id ON db_mecanica_xpto.tb_vehicle (customer_id);

-- Exemplo de inserção de dados na tabela
INSERT INTO db_mecanica_xpto.tb_vehicle (plate, customer_id, model, brand, year)
    VALUES ('ABC-1234', 1, 'Corsa', 'Chevrolet', 2005);
