# Go CRUD API

This is a simple CRUD API (simple Bank account) using Go and Postgres. It allows you to create, read, update and delete accounts, also utilizes jwt to authenticate users. The initial setup was done using [this tutorial](https://www.youtube.com/watch?v=9VRvEzXgImM&t=2s).

As it was not continued by the mentor, I decided to continue it and make it better.

The folder structure I used inspired by [golang documentation](https://go.dev/doc/modules/layout)

I tried to better structure the code and add a config feature to it.

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

I am not an experience golang developer. Please feel free to use it and contribute to it. I am open to suggestions and improvements.
