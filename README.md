
This project implements a CRUD (Create, Read, Update, Delete) API server using Go, a database, and the HTTP protocol. It allows you to perform basic CRUD operations on a specific entity.

## Table of Contents
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
  - [GetByID](#getbyid)
  - [Create](#create)
  - [Delete](#delete)
  - [Update](#update)
- [Testing](#testing)
- [Contributing](#contributing)
- [License](#license)

## Features

- GetByID: Retrieve data from the database based on the provided ID.
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
    https://github.com/kemalkochekov/Go-Backend-CRUD-Api-Server.git
  ```

2. Navigate to the project directory:
  ```
    cd Go-Backend-CRUD-Api-Server
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
    http://127.0.0.1:9000
  ```
## API Endpoints

### GetByID

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

For detailed API documentation, including examples, request/response structures, and authentication details, please refer to the

<a href="https://documenter.getpostman.com/view/31073105/2s9YeN2oV9" target="_blank">
    <img alt="View API Doc Button" src="https://github.com/kemalkochekov/Go-Backend-CRUD-Api-Server/assets/85355663/e5cc7ad1-a31f-4c0d-b4b7-c4ab6e69f5a7" width="200" height="60"/>
</a>

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
