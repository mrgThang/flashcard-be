package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/mrgThang/flashcard-be/services"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "flashcard-be",
		Usage: "Flashcard backend server",
		Commands: []*cli.Command{
			{
				Name:  "serve",
				Usage: "Run the HTTP server",
				Action: func(c *cli.Context) error {
					runServer()
					return nil
				},
			},
			{
				Name:  "migrate",
				Usage: "Run database migrations",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "dir",
						Aliases: []string{"d"},
						Value:   "./migrations",
						Usage:   "Directory containing migration files",
					},
				},
				Action: func(c *cli.Context) error {
					dir := c.String("dir")
					runMigrate(dir)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func runServer() {
	service := services.NewService()

	r := chi.NewRouter()

	// Deck routes
	r.Get("/decks", service.GetDecksHandler)
	r.Post("/decks", service.CreateDeckHandler)
	r.Put("/decks", service.UpdateDeckHandler)

	// Card routes
	r.Get("/cards", service.GetCardsHandler)
	r.Post("/cards", service.CreateCardHandler)
	r.Put("/cards", service.UpdateCardHandler)

	// User routes
	r.Get("/users", service.GetUserHandler)
	r.Post("/users", service.CreateUserHandler)
	r.Put("/users", service.UpdateUserHandler)

	http.ListenAndServe(":8080", r)
}

func runMigrate(migrationsDir string) {
	fmt.Println("Running migrations...")
	services.RunMigrations(migrationsDir)
}
