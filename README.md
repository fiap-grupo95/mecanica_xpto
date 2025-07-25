# mecanica_xpto

## Descrição
O projeto "mecanica_xpto" é um sistema integrado de atendimento e execução de serviços, desenvolvido em Go (Golang) e utilizando o framework Gin para a construção de APIs RESTful. O sistema é projetado para gerenciar atendimentos, serviços e clientes de forma eficiente, com uma arquitetura modular que facilita a manutenção e escalabilidade.

Event Storming: https://miro.com/app/board/uXjVIgU2y2I=/

## Pre-requisitos
- Golang 1.24 ou superior
- Docker
- Docker Compose

## Instalação
1. Clone o repositório:
   ```bash
   git clone
2. Navegue até o diretório do projeto:
   ```bash
   cd mecanica_xpto
   ```
3. Crie um arquivo `.env` baseado no `.env.example`:
   ```bash 
   mv .env.example .env
   ```
4. Inicie os serviços com Docker Compose:
   ```bash
    docker-compose up -d
    ```
5. Acesse a aplicação em `http://localhost:8080`

6. Para parar os serviços, use:
   ```bash
   docker-compose down
   ```
## Testes
Para executar os testes, utilize o seguinte comando:

```bash
go test ./... -v
```