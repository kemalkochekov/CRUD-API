# CRUD API

This project implements a CRUD (Create, Read, Update, Delete) API server for database and the HTTP protocol server by using Go. 

## Table of Contents
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
  - [Get](#get)
  - [Create](#create)
  - [Delete](#delete)
  - [Update](#update)
  - [Api Documentation](#api-documentation)
- [Linting and Code Quality](#linting-and-code-quality)
  - [Linting Installation](#linting-installation)
  - [Linting Usage](#linting-usage)
- [Testing](#testing)
- [Contributing](#contributing)
- [License](#license)


## Features

- Get: Retrieve data from the database based on the provided ID.
- Create: Add new data to the database.
- Delete: Remove data from the database based on the provided ID.
- Update: Update existing data in the database based on the provided ID.

## Prerequisites

Before running this application, ensure that you have the following prerequisites installed:

- Go: [Install Go](https://go.dev/doc/install/)
- Docker: [Install Docker](https://docs.docker.com/get-docker/)
- Docker Compose: [Install Docker Compose](https://docs.docker.com/compose/install/)

## Installation

1. Clone the repository:
  ```bash
    https://github.com/kemalkochekov/CRUD-API.git
  ```

2. Navigate to the project directory:
  ```
    cd CRUD-API
  ```
3. Build the Docker image:
  ```
    docker-compose build
  ```

## Usage
1. Start the Docker containers:
  ```
    docker-compose up
  ```
2. The application will be accessible at:
  ```
    localhost:9000
  ```
## API Endpoints

### Get

- Method: GET
- Endpoint: /entity
- Query Parameter: id=[entity_id]

Retrieve data from the database based on the provided ID.

- Response:
  - 200 OK: Returns the data in the response body.
  - 400 Bad Request: If the `id` query parameter is missing.
  - 404 Not Found: If the provided ID does not exist in the database.
  - 500 Internal Server Error: If there is an internal server error.

### Create

- Method: POST
- Endpoint: /entity
- Request Body: JSON payload containing the ID and data.

Add new data to the database.

- Response:
  - 200 OK: If the request is successful.
  - 500 Internal Server Error: If there is an internal server error.

### Delete

- Method: DELETE
- Endpoint: /entity
- Query Parameter: id=[entity_id]

Remove data from the database based on the provided ID.

- Response:
  - 200 OK: If the request is successful.
  - 404 Not Found: If the provided ID does not exist in the database.
  - 500 Internal Server Error: If there is an internal server error.

### Update

- Method: PUT
- Endpoint: /entity
- Request Body: JSON payload containing the ID and updated data.

Update existing data in the database based on the provided ID.

- Response:
  - 200 OK: If the request is successful.
  - 404 Not Found: If the provided ID does not exist in the database.
  - 500 Internal Server Error: If there is an internal server error.

### Api Documentation

For detailed API documentation, including examples, request/response structures, and authentication details, please refer to the

<a href="https://documenter.getpostman.com/view/31073105/2s9YeN2oV9" target="_blank">
    <img alt="View API Doc Button" src="https://github.com/kemalkochekov/Go-Backend-CRUD-Api-Server/assets/85355663/e5cc7ad1-a31f-4c0d-b4b7-c4ab6e69f5a7" width="200" height="60"/>
</a>

## Linting and Code Quality

This project maintains code quality using `golangci-lint`, a fast and customizable Go linter. `golangci-lint` checks for various issues, ensures code consistency, and enforces best practices, helping maintain a clean and standardized codebase.

### Linting Installation

To install `golangci-lint`, you can use `brew`:

```bash
  brew install golangci-lint
```

### Linting Usage
1. Configuration: 

After installing golangci-lint, create or use a personal configuration file (e.g., .golangci.yml) to define specific linting rules and settings:
```bash
  golangci-lint run --config=.golangci.yml
```
This command initializes linting based on the specified configuration file.

2. Run the linter:

Once configuration is completed, you can execute the following command at the root directory of your project to run golangci-lint:

```bash
  golangci-lint run
```
This command performs linting checks on your entire project and provides a detailed report highlighting any issues or violations found.

3. Customize Linting Rules:

You can customize the linting rules by modifying the .golangci.yml file.

For more information on using golangci-lint, refer to the golangci-lint documentation.


## Testing

- To run unit tests, use the following command:
  ```bash
    go test ./... -cover
  ```
- To run integration tests, use the following command:
  ```bash
    go test -tags integration ./... -cover
  ```

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request.

## License

This project is licensed under the [MIT License](LICENSE).
