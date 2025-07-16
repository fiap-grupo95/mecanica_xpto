#!/bin/bash

# Sobe o container MySQL
docker compose up -d db

# Aguarda o MySQL iniciar
echo "Aguardando o MySQL iniciar..."

# Cria tabela exemplo e insere um item
docker exec -i mysql_db mysql -uxpto_user -pxpto_password mecanica_xpto <<EOF
CREATE TABLE IF NOT EXISTS exemplo (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nome VARCHAR(100) NOT NULL
);

INSERT INTO exemplo (nome) VALUES ('Item exemplo');

SELECT * FROM exemplo;
EOF

echo "Tabela criada, item inserido e select realizado!"