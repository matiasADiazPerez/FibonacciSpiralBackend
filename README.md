# Fibonacci Spiral Backend

This repository contains the implementation of the Fibonacci Spiral Matrix challenge. The API offers the following features:

- **User CRUD Operations:** Complete CRUD functionality for managing users.
- **Login and JWT Authorization:** Secure endpoints with login and JWT-based authorization.
- **Spiral Fibonacci Matrix Calculation:** Generates a matrix of size c x r (columns x rows), with customizable rows and columns (up to 100). The matrix is filled with Fibonacci numbers arranged in a spiral pattern.

The API is built using Golang, employing go-chi as the router and gorm as the database client. PostgreSQL serves as the database engine. Both the API and the database are containerized using Docker.

## Usage

1. Create a `.env` file with the following configurations:
```
POSTGRES_HOST_AUTH_METHOD=trust
POSTGRES_USER={Some user}
POSTGRES_PASSWORD={a password}
POSTGRES_HOST=spiral_db
JWT_SECRET={a secret}
```
2. To deploy the containers and serve the API at `localhost:8080`, use the command `make app-deploy`
3. To run unit tests, execute: `make test`
4. To check the code for linting issues and view errors and warnings related to code style: `make lint`


For detailed API usage instructions, refer to the API documentation served in `http://localhost:8080/docs/`.

## Module Diagrams

### API Module Diagram
![Untitled Diagram drawio](https://github.com/matiasADiazPerez/FibonacciSpiralBackend/assets/130945302/a35c1093-f646-434d-a1e4-d12b56f750d1)


### API Infrastructure Diagram
![Untitled Diagram drawio(1)](https://github.com/matiasADiazPerez/FibonacciSpiralBackend/assets/130945302/06348072-4788-4a17-a4ee-9d5e550aae0f)


