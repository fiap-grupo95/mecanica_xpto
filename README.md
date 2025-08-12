
# mecanica_xpto

## Descrição

O projeto "mecanica_xpto" é um sistema integrado de atendimento e execução de serviços, desenvolvido em Go (Golang) e utilizando o framework Gin para a construção de APIs RESTful. O sistema é projetado para gerenciar atendimentos, serviços e clientes de forma eficiente, com uma arquitetura modular que facilita a manutenção e escalabilidade.

Event Storming: <https://miro.com/app/board/uXjVIgU2y2I=/>

## Pré-requisitos

- Golang 1.24 ou superior  
- Docker  
- Docker Compose  
- Make (para facilitar comandos via Makefile)

## Instalação

1. Clone o repositório:

   ```bash
   git clone <url-do-repositorio>
   ```

2. Navegue até o diretório do projeto:

   ```bash
   cd mecanica_xpto
   ```

3. Inicialize o ambiente, que vai:
   - copiar o arquivo `.env.example` para `.env` (sem sobrescrever se já existir)
   - instalar o Swag CLI (se necessário)
   - gerar a documentação Swagger
   - subir os containers Docker (app e banco)
   - aguardar banco e app ficarem prontos
   - executar migrations e seeds

   Para isso, execute:

   ```bash
   make init
   ```

4. A aplicação estará disponível em `http://localhost:8080`

## Comandos úteis

- Para subir os containers (build + background):

  ```bash
  make up
  ```

- Para parar e remover os containers:

  ```bash
  make down
  ```

- Para acompanhar os logs do container da aplicação:

  ```bash
  make logs
  ```

- Para gerar a documentação Swagger (a partir dos comentários no código):

  ```bash
  make swag-generate
  ```

- Para rodar a aplicação localmente (fora do Docker) junto com a geração da documentação Swagger:

  ```bash
  make swag-run
  ```

## Testes

Para executar os testes automatizados, utilize o seguinte comando:

```bash
go test ./... -v
```

## Documentação da API

### Swagger

Primeiro, instale o Swagger CLI:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

Depois, baixe as dependências do projeto:

```bash
go mod tidy
```

Finalmente, gere a documentação Swagger:

```bash
swag init -g cmd/api/main.go
```

Após rodar esses comandos, você poderá acessar a documentação da API em:  
`http://localhost:8080/swagger/index.html` enquanto a aplicação estiver rodando.
