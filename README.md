# Flashcard Backend (flashcard-be)

A backend server for managing flashcards and decks, built with Go.

## Features

- User authentication (signup, login)
- CRUD operations for decks and cards
- Study mode for cards
- RESTful API with versioning (`/v1`)
- CORS support

## Requirements

- Go 1.18+
- [Chi router](https://github.com/go-chi/chi)
- [urfave/cli](https://github.com/urfave/cli)
- A database (see migrations)

## Getting Started

### Clone the repository

```sh
git clone https://github.com/mrgThang/flashcard-be.git
cd flashcard-be
```

### Install dependencies

```sh
go mod tidy
```

### Run database migrations

```sh
go run main.go migrate --dir ./migrations
```

### Start the server

```sh
go run main.go serve
```

The server will start on the port specified in your configuration.

## API Endpoints

All endpoints are prefixed with `/v1`.

### Decks

- `GET /v1/decks` - List decks (auth required)
- `GET /v1/decks/{id}` - Get deck details (auth required)
- `POST /v1/decks` - Create a deck (auth required)
- `PUT /v1/decks` - Update a deck (auth required)

### Cards

- `GET /v1/cards` - List cards (auth required)
- `POST /v1/cards` - Create a card (auth required)
- `PUT /v1/cards` - Update a card (auth required)
- `PUT /v1/cards/study` - Study a card (auth required)

### Users

- `GET /v1/users` - Get user info (auth required)
- `POST /v1/signup` - Register a new user
- `POST /v1/login` - Login

## Configuration

Configuration is loaded via the `services.NewService()` method. Adjust as needed for your environment.

## License

MIT
