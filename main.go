package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"time"
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
    http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("pong"))
    })
    
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback for local dev
	}

	log.Println("Server started on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
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
func uploadToImgBB(file multipart.File) (string, error) {
	apiKey := "a9edffca025c2863cd7f441529af0a54"

	// Read file into memory
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, file)
	if err != nil {
		return "", err
	}

	// Base64 encode
	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())

	// Prepare request
	data := url.Values{}
	data.Set("key", apiKey)
	data.Set("image", encoded)

	// Send request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.PostForm("https://api.imgbb.com/1/upload", data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check response code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ImgBB upload failed: %s", string(body))
	}

	// Parse JSON
	var result struct {
		Data struct {
			URL string `json:"url"`
		} `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	return result.Data.URL, nil
}
