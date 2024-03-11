# Go CRUD API

This is a simple CRUD API (simple Bank account) using Go and Postgres. It allows you to create, read, update and delete accounts, also utilizes jwt to authenticate users. .

The folder structure I used inspired by [golang documentation](https://go.dev/doc/modules/layout)

The folder structure is as follows:

```
go-crud-api
├── cmd
│   └── api-server
│       └── main.go
│       └── config
│           └── config.go
│       └── seed
│           └── seed.go
├── internal
│   └── api
│       └── api.go
│       └── handler.go
│       └── middleware.go
│       └── server.go
│   └── model
│       └── types.go
│       └── types_test.go
├── store
│   └── postgres.go
│   └── store.go
├── util
│   └── util.go
```
