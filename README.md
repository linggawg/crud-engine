[![SonarCloud](https://sonarcloud.io/images/project_badges/sonarcloud-white.svg)](https://sonarcloud.io/summary/new_code?id=Nazyli_curd-engine) <br />
[![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=Nazyli_curd-engine&metric=code_smells)](https://sonarcloud.io/summary/new_code?id=Nazyli_curd-engine)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=Nazyli_curd-engine&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=Nazyli_curd-engine)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=Nazyli_curd-engine&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=Nazyli_curd-engine)
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=Nazyli_curd-engine&metric=bugs)](https://sonarcloud.io/summary/new_code?id=Nazyli_curd-engine)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=Nazyli_curd-engine&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=Nazyli_curd-engine)
[![Duplicated Lines (%)](https://sonarcloud.io/api/project_badges/measure?project=Nazyli_curd-engine&metric=duplicated_lines_density)](https://sonarcloud.io/summary/new_code?id=Nazyli_curd-engine)

# crud-engine
Crud Engine using Go

## Prerequisite
- Docker Engine Latest
- Docker Compose Latest 
- Golang
- RDBMS (PostgreSQL)

## How to Run ?
- edit .env file to setting connection database, and place at root project path
```sh
DB_DIALECT=postgres
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=crud_go
DB_HOST=YOUR_IP_HOST/SERVER
DB_PORT=5432
```
- and run docker-compose up -d at root project path
```sh
docker-compose up -d
```
## SwaggerUI
- Generate Swagger Specification and SwaggerUI
```sh
$ swag init -g main.go --output docs
```
- Documentation API in path
```sh
{{base_url}}/swagger/index.html
```
