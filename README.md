# Todo List API

This project is a RESTful API built using Go that allows users to manage their to-do list. The API supports user authentication, CRUD operations for to-do items, error handling, and security measures.

## Goals

The skills you will learn from this project include:

- User authentication
- Schema design and databases
- RESTful API design
- CRUD operations
- Error handling
- Security

## Requirements

You are required to develop a RESTful API with the following endpoints:

1. User registration to create a new user
2. Login endpoint to authenticate the user and generate a token
3. CRUD operations for managing the to-do list
4. Implement user authentication to allow only authorized users to access the to-do list
5. Implement error handling and security measures
6. Use a database to store the user and to-do list data (you can use any database of your choice)
7. Implement proper data validation
8. Implement pagination and filtering for the to-do list

## Endpoints

### User Registration

Register a new user using the following request:

```
POST /register
{
  "name": "John Doe",
  "email": "john@doe.com",
  "password": "password"
}
```

Response:

```
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
}
```

### User Login

Authenticate the user using the following request:

```
POST /login
{
  "email": "john@doe.com",
  "password": "password"
}
```

Response:

```
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
}
```

### Create a To-Do Item

Create a new to-do item using the following request:

```
POST /todos
{
  "title": "Buy groceries",
  "description": "Buy milk, eggs, and bread"
}
```

Response:

```
{
  "id": 1,
  "title": "Buy groceries",
  "description": "Buy milk, eggs, and bread"
}
```

### Update a To-Do Item

Update an existing to-do item using the following request:

```
PUT /todos/1
{
  "title": "Buy groceries",
  "description": "Buy milk, eggs, bread, and cheese"
}
```

Response:

```
{
  "id": 1,
  "title": "Buy groceries",
  "description": "Buy milk, eggs, bread, and cheese"
}
```

### Delete a To-Do Item

Delete an existing to-do item using the following request:

```
DELETE /todos/1
```

Response:

Status code: 204

### Get To-Do Items

Get the list of to-do items using the following request:

```
GET /todos?page=1&limit=10
```

Response:

```
{
  "data": [
    {
      "id": 1,
      "title": "Buy groceries",
      "description": "Buy milk, eggs, bread"
    },
    {
      "id": 2,
      "title": "Pay bills",
      "description": "Pay electricity and water bills"
    }
  ],
  "page": 1,
  "limit": 10,
  "total": 2
}
```

## Bonus Features

- Implement filtering and sorting for the to-do list
- Implement unit tests for the API
- Implement rate limiting and throttling for the API
- Implement refresh token mechanism for the authentication

## Installation

1. Clone the repository.
2. Install the dependencies using `go mod tidy`.
3. Set up the database connection.
4. Start the server using `go run main.go`.

## License

This project is licensed under the MIT License.

https://roadmap.sh/projects/todo-list-api