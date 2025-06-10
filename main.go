package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/urfave/cli/v2"

	"github.com/mrgThang/flashcard-be/logger"
	"github.com/mrgThang/flashcard-be/middlewares"
	"github.com/mrgThang/flashcard-be/services"
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
	if err := logger.Init(); err != nil {
		panic(err)
	}
	service := services.NewService()

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	v1 := chi.NewRouter()

	// Deck routes
	v1.Get("/decks", middlewares.AuthMiddleware(service, service.GetDecksHandler))
	v1.Get("/decks/{id}", middlewares.AuthMiddleware(service, service.GetDetailDeckHandler))
	v1.Post("/decks", middlewares.AuthMiddleware(service, service.CreateDeckHandler))
	v1.Put("/decks", middlewares.AuthMiddleware(service, service.UpdateDeckHandler))

	// Card routes
	v1.Get("/cards", middlewares.AuthMiddleware(service, service.GetCardsHandler))
	v1.Post("/cards", middlewares.AuthMiddleware(service, service.CreateCardHandler))
	v1.Put("/cards", middlewares.AuthMiddleware(service, service.UpdateCardHandler))
	v1.Put("/cards/study", middlewares.AuthMiddleware(service, service.StudyCardHandler))

	// User routes
	v1.Get("/users", middlewares.AuthMiddleware(service, service.GetUserHandler))

	v1.Post("/signup", service.SignupHandler)
	v1.Post("/login", service.LoginHandler)

	// create prefix v1 for all routes
	r.Mount("/v1", v1)

	fmt.Println(fmt.Sprintf("Server is running at port %s", service.Config.Port))
	return http.ListenAndServe(fmt.Sprintf(":%s", service.Config.Port), r)
}

func runMigrate(migrationsDir string) {
	fmt.Println("Running migrations...")
	services.RunMigrations(migrationsDir)
}
