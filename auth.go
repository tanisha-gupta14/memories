package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Render the signup page
		templates.ExecuteTemplate(w, "signup.html", nil)
		return
	}

	// Handle POST request for signup
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Check if username and password are provided
	if username == "" || password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	// Insert the user into the database
	_, err = DB.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", username, string(hashedPassword))
	if err != nil {
		http.Error(w, "Error creating account", http.StatusInternalServerError)
		return
	}

	// Redirect to login page after successful signup
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Render the login page
		templates.ExecuteTemplate(w, "login.html", nil)
		return
	}

	// Handle POST request for login
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Check if username and password are provided
	if username == "" || password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Query the database for the user by username
	var storedPasswordHash string
	err := DB.QueryRow("SELECT password_hash FROM users WHERE username = ?", username).Scan(&storedPasswordHash)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Compare the hashed password with the stored password
	err = bcrypt.CompareHashAndPassword([]byte(storedPasswordHash), []byte(password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Redirect to My Memories page after successful login
	// After successful login, set the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "user_id",
		Value:    username, // You can use user_id or username, depending on your application
		HttpOnly: true,     // To prevent JavaScript access
		Path:     "/",
		MaxAge:   3600,     // 1 hour expiration
	})

	http.Redirect(w, r, "/mymemories", http.StatusSeeOther)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Clear the "user_id" cookie (same name used in LoginHandler)
	http.SetCookie(w, &http.Cookie{
		Name:     "user_id", 
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	// Redirect to the homepage after logout
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func CheckLoginStatusHandler(w http.ResponseWriter, r *http.Request) {
    cookie, err := r.Cookie("user_id")
    isLoggedIn := err == nil && cookie.Value != ""

    // Respond with login status as JSON
    response := struct {
        LoggedIn bool `json:"loggedIn"`
    }{
        LoggedIn: isLoggedIn,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
