# task-manager-auth-service

## Description

This is a golang authentication microservice with a gRPC server related to a Task Management API. This microservice provides authentication functionalities such as register, login and validate token.

## Folder structure

```
├───api
│   └───pb
├───cmd
├───internal
│   ├───db
│   ├───domain
│   ├───errors
│   ├───handlers
│   ├───repositories
│   ├───service
│   └───utils
└───migrations
```

- api - contains proto definitions
- cmd - contains the main application entrypoint
- internal - contains application core components
    - db - database connection
    - domain - business logic/models of the application
    - service - use case layer of the application (where the business logic is written)
    - repositories - where the db CRUD operations are
    - handlers - handlers/controllers of the application
    - utils - application utils
    - errors - custom error package
- migrations - contains init.sql script

## Prerequisites
- Golang installed in your system (this project use go version 1.19)
- Docker installed in your system
- Make installed in your system

# Getting started

1. Clone the repository

```
$ git clone https://github.com/niluwats/task-manager-auth-service.git
```

2. Navigate to the project directory

```
$ cd task-manager-auth-service
```

3. Build and run application using make

```
$ make up_build
```

## Dependencies
- [jwt-go](https://github.com/golang-jwt/jwt) - generate and validate JWT
- [crypto](https://github.com/golang/crypto) - hashing passwords
- [gorm](https://github.com/go-gorm/gorm) - gorm for object relational mapping in golang
- [db driver](https://github.com/go-gorm/postgres) - postgresql driver for gorm

```
$ go get github.com/golang-jwt/jwt
$ go get golang.org/x/crypto
$ go get go get -u gorm.io/gorm
$ go get gorm.io/driver/postgres
```
