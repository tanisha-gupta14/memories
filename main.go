package main

import (
	"log"
	"net/http"
)

func main() {
	// Initialize DB connection
	InitDB()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	
	
	// Routes
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/signup", SignupHandler)
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/logout", LogoutHandler)
	http.HandleFunc("/all-memories", AllMemoriesHandler)
	http.HandleFunc("/my-memories", MyMemoriesHandler)
	http.HandleFunc("/create-memory", CreateMemoryHandler)
	http.HandleFunc("/edit-memory/", EditMemoryHandler)
	http.HandleFunc("/update-memory/", UpdateMemoryHandler)
	http.HandleFunc("/delete-memory/", DeleteMemoryHandler)
	http.HandleFunc("/check-login-status", CheckLoginStatusHandler)

	// Start server
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
// Define a struct to hold data for the template rendering
type MemoryPageData struct {
    IsLoggedIn bool
    Memories   []Memory
}


func IndexHandler(w http.ResponseWriter, r *http.Request) {
    username, _ := getUsernameFromSession(r)

    // Fetch public memories
    publicRows, err := DB.Query("SELECT memory_id, title, description, image_path FROM memories WHERE privacy = 'public'")
    if err != nil {
        log.Println("Error fetching public memories:", err)
        http.Error(w, "Unable to fetch public memories", http.StatusInternalServerError)
        return
    }
    defer publicRows.Close()

    var publicMemories []Memory
    for publicRows.Next() {
        var m Memory
        err := publicRows.Scan(&m.MemoryID, &m.Title, &m.Description, &m.ImagePath)
        if err != nil {
            log.Println("Error scanning public memories:", err)
            http.Error(w, "Error reading memories", http.StatusInternalServerError)
            return
        }
        publicMemories = append(publicMemories, m)
    }

    // Fetch user memories (only if logged in)
    var userMemories []Memory
    var loggedIn bool
    if username != "" {
        loggedIn = true
        userRows, err := DB.Query("SELECT memory_id, title, description, image_path, privacy FROM memories WHERE username = ?", username)
        if err != nil {
            log.Println("Error fetching user memories:", err)
            http.Error(w, "Unable to fetch user memories", http.StatusInternalServerError)
            return
        }
        defer userRows.Close()

        for userRows.Next() {
            var m Memory
            err := userRows.Scan(&m.MemoryID, &m.Title, &m.Description, &m.ImagePath, &m.Privacy)
            if err != nil {
                log.Println("Error scanning user memories:", err)
                http.Error(w, "Error reading memories", http.StatusInternalServerError)
                return
            }
            userMemories = append(userMemories, m)
        }
    }

    // Pass the data to the template
    data := struct {
        Username       string
        PublicMemories []Memory
        UserMemories   []Memory
        LoggedIn       bool
    }{
        Username:       username,
        PublicMemories: publicMemories,
        UserMemories:   userMemories,
        LoggedIn:       loggedIn,
    }

    // Render the template
    templates.ExecuteTemplate(w, "index.html", data)
}
