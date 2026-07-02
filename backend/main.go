package main

import (
	"blogsite/handlers"
	"blogsite/internal/database"
	"blogsite/middleware"
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"

	firebase "firebase.google.com/go/v4"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func initFirebase() *firebase.App {
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	return app
}

func main() {
	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT string not available")
	}

	dbConn, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal("Cannot connect to DB:", err)
	}

	queries := database.New(dbConn)

	// 1. Initialize Firebase before setting up the router
	firebaseApp := initFirebase()

	userHandler := &handlers.UserHandler{Queries: queries}
	postHandler := &handlers.PostHandler{Queries: queries}

	router := chi.NewRouter()

	// 2. Your existing CORS setup (Leave this exactly as is)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// 3. Public Routes (No authentication required)
	router.Get("/posts", postHandler.ListPosts)

	// 4. Protected Routes (Requires Firebase Token)
	router.Group(func(r chi.Router) {
		// Intercept these specific routes with your Auth Middleware
		r.Use(middleware.AuthMiddleware(firebaseApp))

		r.Post("/api/register", userHandler.RegisterUser)
		r.Get("/api/me", userHandler.GetMe)
		r.Get("/api/me/posts", postHandler.ListMyPosts)

		r.Get("/api/posts", postHandler.ListPosts)
		r.Post("/api/posts", postHandler.CreatePost)
		r.Get("/api/posts/{id}", postHandler.GetPost)
		r.Put("/api/posts/{id}", postHandler.UpdatePost)
		r.Delete("/api/posts/{id}", postHandler.DeletePost)

	})

	// 5. Your existing Server setup
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server starting on Port %v", portString)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
