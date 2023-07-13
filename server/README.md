## Setup local development

### Install tools[MacOS]

- [Docker desktop](https://www.docker.com/products/docker-desktop)
- [DBeaver](https://dbeaver.com)
- [Golang](https://golang.org/)
- [Homebrew](https://brew.sh/)
- [Migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

  ```bash
   brew install golang-migrate
   brew install node
  ```

- [DB Docs](https://dbdocs.io/docs)

  ```bash
  npm install -g dbdocs
  dbdocs login
  ```

- [DBML CLI](https://www.dbml.org/cli/#installation)

  ```bash
  npm install -g @dbml/cli
  dbml2sql --version
  ```

- [Gomock](https://github.com/golang/mock)
  ```bash
  go install github.com/golang/mock/mockgen@v1.6.0
  ```

### Setup infrastructure

- Create the network
  ```bash
  make network
  ```
- Start postgres container:
  ```bash
  make postgres
  ```
- Create database:
  ```bash
  make createdb
  ```
- Run db migration up all versions:
  ```bash
  make migrateup
  ```
- Run db migration up 1 version:
  ```bash
  make migrateup1
  ```
- Run db migration down all versions:
  ```bash
  make migratedown
  ```
- Run db migration down 1 version:
  ```bash
   make migratedown1
  ```

### Documentation

- Generate DB documentation:

  ```bash
  make db_docs
  ```

- Access the DB documentation at [this address](https://dbdocs.io/hugodfbec4a9c6/hr_monitor?schema=public&view=table_structure&table=acts).

### How to generate code

- Generate schema SQL file with DBML:

  ```bash
  make db_schema
  ```

- Create a new db migration:
  ```bash
  migrate create -ext sql -dir db/migration -seq <migration_name>
  ```

### How to run

- Run server:
  ```bash
  make server
  ```