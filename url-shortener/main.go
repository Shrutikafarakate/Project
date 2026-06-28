package main

import (
	"log"
	"net/http"
	"url-shortener/db"
	"url-shortener/handlers"
	"url-shortener/models"
)

func main() {
	// Initialize the database connection
	database := db.InitDB()
	defer database.Close()

	// Delete any expired URLs on server start
	if err := models.DeleteExpiredURLs(database); err != nil {
		log.Println("Error deleting expired URLs:", err)
	}

	// Serve static files (CSS, JS, etc.)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Route handlers
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)
	http.HandleFunc("/shorten", handlers.ShortenURLHandler)
	http.HandleFunc("/u/", handlers.RedirectHandler)
	http.HandleFunc("/about", handlers.AboutHandler)


	// Start the HTTP server
	log.Println("Server running at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
