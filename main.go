package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var (
	db *sql.DB
	counter int
	mu sync.Mutex
	activeConnections = make(map[string]chan bool) // Potential memory leak - connections not cleaned up
	connectionsMu sync.RWMutex
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	// Initialize database connection
	var err error
	db, err = sql.Open("postgres", "user=postgres password=password dbname=testdb sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Create Gin router
	r := gin.Default()

	// Routes
	r.GET("/users/:id", getUserHandler)
	r.POST("/users", createUserHandler)
	r.GET("/counter", getCounterHandler)
	r.POST("/increment", incrementHandler)
	r.GET("/stream/:id", streamHandler) // Fixed memory leak - proper cleanup
	r.DELETE("/stream/:id", stopStreamHandler) // Helper endpoint to stop streams

	// Start server
	fmt.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// Fixed: SQL Injection vulnerability - now using parameterized queries
func getUserHandler(c *gin.Context) {
	userIDStr := c.Param("id")
	
	// Validate and convert user ID to integer
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}
	
	// SECURE: Using parameterized query to prevent SQL injection
	query := "SELECT id, name, email FROM users WHERE id = $1"
	
	var user User
	err = db.QueryRow(query, userID).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}
	
	c.JSON(http.StatusOK, user)
}

func createUserHandler(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := "INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id"
	err := db.QueryRow(query, user.Name, user.Email).Scan(&user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Fixed: Race condition - counter is now thread-safe using mutex
func getCounterHandler(c *gin.Context) {
	// Protected read with mutex
	mu.Lock()
	currentCounter := counter
	mu.Unlock()
	
	c.JSON(http.StatusOK, gin.H{"counter": currentCounter})
}

func incrementHandler(c *gin.Context) {
	// Simulating work
	time.Sleep(10 * time.Millisecond)
	
	// Fixed: Race condition - using mutex to protect counter access
	mu.Lock()
	counter++
	currentCounter := counter
	mu.Unlock()
	
	c.JSON(http.StatusOK, gin.H{"counter": currentCounter})
}

// Fixed: Memory leak - proper cleanup of goroutines and channels
func streamHandler(c *gin.Context) {
	clientID := c.Param("id")
	
	// Create a context with timeout for proper cleanup
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()
	
	// Create a channel for this connection with proper cleanup
	ch := make(chan bool, 1) // Buffered to prevent blocking
	connectionsMu.Lock()
	activeConnections[clientID] = ch
	connectionsMu.Unlock()
	
	// Cleanup function to remove connection and close channel
	cleanup := func() {
		connectionsMu.Lock()
		if existingCh, exists := activeConnections[clientID]; exists {
			close(existingCh)
			delete(activeConnections, clientID)
		}
		connectionsMu.Unlock()
	}
	
	// Start a goroutine with proper cleanup
	go func() {
		defer cleanup() // Ensure cleanup happens when goroutine exits
		
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop() // Prevent ticker leak
		
		for {
			select {
			case <-ctx.Done():
				// Context cancelled or timeout reached
				return
			case <-ch:
				// Channel closed, exit gracefully
				return
			case <-ticker.C:
				// Simulate sending data
				fmt.Printf("Sending data to client %s\n", clientID)
			}
		}
	}()
	
	// Set up cleanup when client disconnects
	go func() {
		<-c.Request.Context().Done()
		cleanup()
	}()
	
	c.JSON(http.StatusOK, gin.H{"message": "Stream started", "client_id": clientID})
}

// Helper function to stop a specific stream (demonstrates proper cleanup)
func stopStreamHandler(c *gin.Context) {
	clientID := c.Param("id")
	
	connectionsMu.Lock()
	if ch, exists := activeConnections[clientID]; exists {
		close(ch)
		delete(activeConnections, clientID)
	}
	connectionsMu.Unlock()
	
	c.JSON(http.StatusOK, gin.H{"message": "Stream stopped", "client_id": clientID})
}

