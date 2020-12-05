# Learning to write first server in go lang with example to build simple bank application

We are going to use 
- go lang as programming language
- postgres as our datasource
- docker as our container engine

To generate the schema diagram I have use https://dbdiagram.io/

Prequisite 
- install docker
- install go
- golang-migrate [https://github.com/golang-migrate/migrate]

To install postgres using docker follow below steps:
- pull image for postgres using below command
`docker pull postgres:13-alpine`
- run dokcer container using above image
`docker run --name postgres13 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSOWRD=secret -d postgres:13-alpine`
- exec in to container to check if it working
`docker exec -it postgres13 psql -U root`
- to check the logs for the postgres container
`docker logs -f postgres13`

To simplify the local development please use Makefile commands