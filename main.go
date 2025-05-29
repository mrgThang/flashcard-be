package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/mrgThang/flashcard-be/logger"
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
					return runServer()
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

func runServer() error {
	logger.Init()
	service := services.NewService()

	r := chi.NewRouter()

	v1 := chi.NewRouter()

	// Deck routes
	v1.Get("/decks", service.GetDecksHandler)
	v1.Post("/decks", service.CreateDeckHandler)
	v1.Put("/decks", service.UpdateDeckHandler)

	// Card routes
	v1.Get("/cards", service.GetCardsHandler)
	v1.Post("/cards", service.CreateCardHandler)
	v1.Put("/cards", service.UpdateCardHandler)

	// User routes
	v1.Get("/users", service.GetUserHandler)
	v1.Post("/users", service.CreateUserHandler)
	v1.Put("/users", service.UpdateUserHandler)

	// create prefix v1 for all routes
	r.Mount("/v1", v1)

	fmt.Println(fmt.Sprintf("Server is running at port %s", service.Config.Port))
	return http.ListenAndServe(fmt.Sprintf(":%s", service.Config.Port), r)
}

func runMigrate(migrationsDir string) {
	fmt.Println("Running migrations...")
	services.RunMigrations(migrationsDir)
}
