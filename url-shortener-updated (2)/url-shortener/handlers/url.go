package handlers

import (
	"net/http"
	"time"

	"url-shortener/models"
	"url-shortener/utils"
)

// HomeHandler displays the home page with the user's shortened URLs
// HomeHandler displays the home page with the user's shortened URLs
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	session := utils.GetSession(w, r)
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Fetch URLs associated with the user
	urls, err := models.GetUserURLs(userID)
	if err != nil {
		http.Error(w, "Failed to retrieve URLs", http.StatusInternalServerError)
		return
	}

	// Convert slice to map for template compatibility
	urlMap := make(map[string]models.URL)
	for _, url := range urls {
		urlMap[url.ShortCode] = url
	}

	data := struct {
		Username string
		URLs     map[string]models.URL
	}{
		Username: session.Values["username"].(string),
		URLs:     urlMap,
	}

	utils.RenderTemplate(w, "home.html", data)
}


// ShortenURLHandler creates a shortened URL and stores it
func ShortenURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Get session and validate user ID
	session := utils.GetSession(w, r)
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Get the original URL and expiry date
	originalURL := r.FormValue("original_url") // Ensure this matches the input name in your HTML
	expiryStr := r.FormValue("expiry")
	if originalURL == "" {
		http.Error(w, "Original URL is required", http.StatusBadRequest)
		return
	}
	
	// Parse expiry date
	expiry, err := time.Parse("2006-01-02T15:04", expiryStr)
	if err != nil {
		http.Error(w, "Invalid expiry date format", http.StatusBadRequest)
		return
	}

	// Generate short code for the URL
	shortCode := utils.GenerateShortCode()

	// Save the URL to the database (or in-memory storage)
	err = models.SaveURL(userID, originalURL, shortCode, expiry)
	if err != nil {
		http.Error(w, "Failed to shorten URL", http.StatusInternalServerError)
		return
	}

	// Optionally, you could redirect to a different page (like the user's homepage with the shortened URLs)
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}


// RedirectHandler handles the redirection for shortened URLs
func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	shortCode := r.URL.Path[len("/u/"):]

	url, err := models.GetURLByCode(shortCode)
	if err != nil || url == nil {
		http.NotFound(w, r)
		return
	}

	if time.Now().After(url.Expiry) {
		http.Error(w, "This URL has expired", http.StatusGone)
		return
	}

	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}
