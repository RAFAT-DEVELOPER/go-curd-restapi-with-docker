# Building a CRUD REST API in Go using Mux, PostgreSQL, Docker, and Docker Compose

Welcome to our exciting journey of building a CRUD REST API in Go! In this tutorial, we'll harness the power of Go, Mux, PostgreSQL, Docker, and Docker Compose to create a robust API that handles Create, Read, Update, and Delete operations.

## Prerequisites

Before we begin, make sure you have the following installed on your system:

- Go installed on your machine.
- Docker and Docker Compose to manage containers.
- Basic knowledge of Go programming language.

## Setting Up the Environment

First things first, let's set up our development environment:

1. **Initialize a Go Module:**
   Create a new directory for your project and initialize a Go module.
   ```bash
   mkdir go-api
   cd go-api
   go mod init github.com/RAFAT-DEVELOPER/go-api
   ```

2. **Install Dependencies:**
   We'll need the following dependencies for our project:
   - [github.com/gorilla/mux](https://github.com/gorilla/mux): Mux router for Go.
   - [github.com/jackc/pgx/v4](https://github.com/jackc/pgx/v4): PostgreSQL driver for Go.
   Install them using `go get`:
   ```bash
   go get -u github.com/gorilla/mux
   go get -u github.com/jackc/pgx/v4
   ```

3. **Prepare Docker and Docker Compose:**
   Ensure Docker and Docker Compose are installed on your system. We'll use Docker to containerize our application and PostgreSQL database.

## Writing the Code

Now that our environment is set up, let's start building our CRUD API.

### Step 1: Define the User Struct

We'll define a `User` struct to represent our data model. It will have fields for `ID`, `Name`, and `Email`.

```go
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

### Step 2: Initialize the PostgreSQL Database

We'll use PostgreSQL as our database. Make sure you have PostgreSQL installed and running on your system.

### Step 3: Implement the CRUD Operations

We'll create handlers for Create, Read, Update, and Delete operations using Gorilla Mux.

### Step 4: Dockerize the Application

We'll create a Dockerfile to containerize our Go application.

```Dockerfile
FROM golang:1.16.3-alpine3.13

WORKDIR /app

COPY . .

RUN go get -d -v ./...
RUN go build -o api .

EXPOSE 8000

CMD ["./api"]
```

### Step 5: Setup Docker Compose

We'll define a `docker-compose.yml` file to manage our Docker containers.

```yaml
version: '3.9'

services:
  go-app:
    container_name: go-app
    build: .
    ports:
      - "8000:8000"
    depends_on:
      - go_db
    environment:
      DATABASE_URL: postgres://username:password@go_db:5432/database_name

  go_db:
    container_name: go_db
    image: postgres:latest
    ports:
      - "5432:5432"
    environment:
      POSTGRES_USER: your_username
      POSTGRES_PASSWORD: your_password
    volumes:
      - pgdata:/var/lib/postgresql/data

volumes:
  pgdata:
```

## Running the Application

To run our application, execute the following commands:

1. **Start Docker Containers:**
   ```bash
   docker-compose up -d go_db
   ```

2. **Build Go Application Image:**
   ```bash
   docker-compose build
   ```

3. **Run the Go Application Container:**
   ```bash
   docker-compose up go-app
   ```

Now, your Go API should be up and running! Access it by making requests to `localhost:8000/users`.

## Conclusion

Congratulations! You've successfully built a CRUD REST API in Go using Mux, PostgreSQL, Docker, and Docker Compose. You're now equipped with the tools to create powerful APIs and applications in Go. Keep exploring, experimenting, and building amazing things! 🚀

