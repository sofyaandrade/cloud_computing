# Plataforma de Treinos Online

## Descrição da aplicação

A Plataforma de Treinos Online é uma aplicação web simples para cadastro e consulta de exercícios físicos.

A aplicação permite:

- cadastrar exercícios com nome, categoria e descrição;
- consultar os exercícios cadastrados em uma tabela;
- armazenar os dados em um banco PostgreSQL executado em container separado.

## Tecnologias utilizadas

- Go / Golang
- PostgreSQL
- Docker
- Docker Compose
- HTML/CSS

## Arquitetura utilizada

O projeto utiliza uma arquitetura multicontainer:

- um container para a aplicação em Go;
- um container para o banco de dados PostgreSQL;
- uma rede Docker para comunicação entre os containers;
- um volume Docker para persistência dos dados do banco.

A aplicação se comunica com o banco usando o nome do serviço `db`, definido no `docker-compose.yml`.

## Procedimento golang

Caso o backend não rode verifcar se na pasta /app existem os seguintes arquivos:
go.mod
go.sum

Caso o go.mod ou ambos estejam faltando rodar os seguintes comando:

```bash
go mod init plataforma-treinos-online

go mod tidy

```

Caso seja apenas o go.sum rodar 

```bash
go mod tidy
```

## Portas utilizadas

| Serviço | Porta no computador | Porta no container |
|--------|----------------------|--------------------|
| Aplicação Go | 8080 | 8080 |
| PostgreSQL | 5432 | 5432 |

## Variáveis de ambiente

| Variável | Valor usado | Descrição |
|---------|-------------|-----------|
| APP_PORT | 8080 | Porta da aplicação |
| DB_HOST | db | Nome do serviço do banco no Docker Compose |
| DB_PORT | 5432 | Porta do PostgreSQL |
| DB_USER | postgres | Usuário do banco |
| DB_PASSWORD | postgres | Senha do banco |
| DB_NAME | treinosdb | Nome do banco de dados |

## Como executar com Docker 

```bash
docker compose up 
```

Depois acesse no navegador:

```text
http://localhost:8080
```

## Evidências

Apasta Evidencias contém os prints referentes ao funcionamento do Docker e da Aplicação de forma local
