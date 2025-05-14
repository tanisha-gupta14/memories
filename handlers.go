package main

import (
	
	"errors"
	
	
	"log"
	"net/http"
	
	"strconv"
	"strings"
	
)

// Example handler for rendering the home page
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("user_id")
	isLoggedIn := err == nil

	data := struct {
		IsLoggedIn bool
	}{
		IsLoggedIn: isLoggedIn,
	}

	templates.ExecuteTemplate(w, "index.html", data)
}

// Memory struct to represent each memory
type Memory struct {
	MemoryID    int
	Title       string
	Description string
	ImagePath   string
	Privacy     string
}

// Handler for "My Memories" page
// Handler for "My Memories" page (to display memories of the logged-in user)
// Handler for "My Memories" page (to display memories of the logged-in user)
// Handler for "My Memories" page (for logged-in users)
func MyMemoriesHandler(w http.ResponseWriter, r *http.Request) {
    // Get the logged-in username from the session
    username, err := getUsernameFromSession(r)
    if err != nil {
        // User is not logged in, redirect to login page
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    // Fetch memories for the logged-in user from the database
    rows, err := DB.Query("SELECT memory_id, title, description, image_path, privacy FROM memories WHERE username = ?", username)
    if err != nil {
        http.Error(w, "Error fetching memories", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var memories []Memory
    for rows.Next() {
        var memory Memory
        if err := rows.Scan(&memory.MemoryID, &memory.Title, &memory.Description, &memory.ImagePath, &memory.Privacy); err != nil {
            http.Error(w, "Error scanning memory", http.StatusInternalServerError)
            return
        }
        memories = append(memories, memory)
    }

    // Render the template with the memories data
    err = templates.ExecuteTemplate(w, "mymemories.html", memories)
    if err != nil {
        http.Error(w, "Error rendering template", http.StatusInternalServerError)
    }
}

// Handler for "All Memories" page (public gallery)
func AllMemoriesHandler(w http.ResponseWriter, r *http.Request) {
    // Fetch public memories from the database
    rows, err := DB.Query("SELECT memory_id, title, description, image_path FROM memories WHERE privacy = 'public'")
    if err != nil {
        http.Error(w, "Error fetching public memories", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var publicMemories []Memory
    for rows.Next() {
        var memory Memory
        if err := rows.Scan(&memory.MemoryID, &memory.Title, &memory.Description, &memory.ImagePath); err != nil {
            http.Error(w, "Error scanning public memories", http.StatusInternalServerError)
            return
        }
        publicMemories = append(publicMemories, memory)
    }

    // Create a MemoryPageData struct to pass to the template
    data := MemoryPageData{
        IsLoggedIn: false, // Assuming user is not logged in for public memories page
        Memories:   publicMemories,
    }

    // Render the "All Memories" page
    templates.ExecuteTemplate(w, "allmemories.html", data)
}

// CreateMemoryHandler handles the memory creation form

func CreateMemoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username, err := getUsernameFromSession(r)
		if err != nil {
			http.Error(w, "Please log in to create a memory", http.StatusUnauthorized)
			return
		}

		title := r.FormValue("title")
		description := r.FormValue("description")
		privacy := r.FormValue("privacy")

		err = r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, "File too large. Max size is 10MB", http.StatusBadRequest)
			return
		}

		file, _, err := r.FormFile("image")
		if err != nil {
			http.Error(w, "Error retrieving image file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// Upload to Imgur instead of saving locally
		imageURL, err := uploadToImgBB(file)
		if err != nil {
			http.Error(w, "Error uploading image to Imgur", http.StatusInternalServerError)
			return
		}

		// Save the Imgur image URL to DB
		_, err = DB.Exec("INSERT INTO memories (username, title, description, image_path, privacy) VALUES (?, ?, ?, ?, ?)",
			username, title, description, imageURL, privacy)
		if err != nil {
			http.Error(w, "Error creating memory", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/mymemories", http.StatusSeeOther)
		return
	}

	templates.ExecuteTemplate(w, "create_memory.html", nil)
}



// Handler to edit an existing memory
// Handler to edit an existing memory
// Handler for editing a memory
func EditMemoryHandler(w http.ResponseWriter, r *http.Request) {
	// Extract memory_id from URL: /edit-memory/1
	idStr := strings.TrimPrefix(r.URL.Path, "/edit-memory/")
	memoryID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid memory ID", http.StatusBadRequest)
		return
	}

	// Get the username from session
	username, err := getUsernameFromSession(r)
	if err != nil {
		http.Error(w, "Please log in to edit your memory", http.StatusUnauthorized)
		return
	}

	// Fetch memory from DB
	var memory Memory
	err = DB.QueryRow(`SELECT memory_id, title, description, image_path, privacy FROM memories WHERE memory_id = ? AND username = ?`, memoryID, username).
		Scan(&memory.MemoryID, &memory.Title, &memory.Description, &memory.ImagePath, &memory.Privacy)
	if err != nil {
		http.Error(w, "Memory not found or you don't have permission", http.StatusNotFound)
		return
	}

	// Show the edit form pre-filled
	err = templates.ExecuteTemplate(w, "edit_memory.html", memory)
	if err != nil {
		log.Println("Template error:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// Handler to update an existing memory
func UpdateMemoryHandler(w http.ResponseWriter, r *http.Request) {
    idStr := strings.TrimPrefix(r.URL.Path, "/update-memory/")
    memoryID, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid memory ID", http.StatusBadRequest)
        return
    }

    // Get username from session
    username, err := getUsernameFromSession(r)
    if err != nil {
        http.Error(w, "Please log in", http.StatusUnauthorized)
        return
    }

    // Parse form
    err = r.ParseMultipartForm(10 << 20)
    if err != nil {
        http.Error(w, "Error parsing form", http.StatusBadRequest)
        return
    }

    title := r.FormValue("title")
    description := r.FormValue("description")
    privacy := r.FormValue("privacy")

    // Optional image upload (this is where we modify the code to use ImgBB)
    var imagePath string
    file, _, err := r.FormFile("image")
    if err == nil {
        // Upload image to ImgBB
        imageURL, err := uploadToImgBB(file)
        if err != nil {
            http.Error(w, "Error uploading image to ImgBB", http.StatusInternalServerError)
            return
        }
        imagePath = imageURL
    } else if err != http.ErrMissingFile {
        http.Error(w, "Error retrieving image", http.StatusInternalServerError)
        return
    }

    // Update query (only update image_path if a new image is uploaded)
    if imagePath != "" {
        _, err = DB.Exec(`UPDATE memories SET title = ?, description = ?, image_path = ?, privacy = ? WHERE memory_id = ? AND username = ?`,
            title, description, imagePath, privacy, memoryID, username)
    } else {
        // If no new image, don't update the image_path
        _, err = DB.Exec(`UPDATE memories SET title = ?, description = ?, privacy = ? WHERE memory_id = ? AND username = ?`,
            title, description, privacy, memoryID, username)
    }

    if err != nil {
        http.Error(w, "Error updating memory", http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/mymemories", http.StatusSeeOther)
}

// Handler to delete an existing memory
func DeleteMemoryHandler(w http.ResponseWriter, r *http.Request) {
	// Extract memory_id from URL like: /delete-memory/3
	idStr := strings.TrimPrefix(r.URL.Path, "/delete-memory/")
	memoryID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid memory ID", http.StatusBadRequest)
		return
	}

	// Get username from session
	username, err := getUsernameFromSession(r)
	if err != nil {
		http.Error(w, "Please log in to delete a memory", http.StatusUnauthorized)
		return
	}

	// Optionally: You can check if the memory belongs to the user before deleting
	res, err := DB.Exec("DELETE FROM memories WHERE memory_id = ? AND username = ?", memoryID, username)
	if err != nil {
		http.Error(w, "Error deleting memory", http.StatusInternalServerError)
		return
	}

	// Check if any row was deleted
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Memory not found or permission denied", http.StatusNotFound)
		return
	}

	// Redirect after successful delete
	http.Redirect(w, r, "/mymemories", http.StatusSeeOther)
}

func getUsernameFromSession(r *http.Request) (string, error) {
	// Retrieve the cookie by name (assuming the cookie is named "user_id")
	cookie, err := r.Cookie("user_id")
	if err != nil {
		return "", errors.New("unauthorized: no user logged in")
	}

	// Return the username stored in the cookie value
	return cookie.Value, nil
}
