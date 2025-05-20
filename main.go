package main

import (
	"log"
	"time"

	"brick/data"
	"brick/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/sqlite3" // For session storage
)

var store *session.Store

func main() {
	// Initialize Database
	data.InitDB("./brick.db") // This will create brick.db in the project root
	defer data.CloseDB()

	// Initialize session store
	storage := sqlite3.New(sqlite3.Config{
		Database: "./brick_sessions.db", // Separate DB for sessions
		Table:    "sessions",
	})
	store = session.New(session.Config{
		Storage: storage,
		CookieHTTPOnly: true,
		CookieSameSite: "Lax",
		Expiration:     24 * time.Hour, // Sessions expire after 24 hours
		// In production, set CookieSecure to true if using HTTPS
	})

	app := fiber.New()

	// Add CORS middleware to allow requests from the Svelte dev server
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:8080,http://localhost:5173", // Allow Svelte dev (Vite/Rollup)
		AllowCredentials: true, // Important for sessions/cookies
		AllowHeaders:     "Origin, Content-Type, Accept",
	}))

	app.Static("/", "./web/ui/public") // Serve Svelte app

	// API routes for user management
	api := app.Group("/api")
	api.Post("/register", registerUserHandler) // Updated to registerUserHandler
	api.Post("/login", loginUserHandler)    // Updated to loginUserHandler
	api.Post("/logout", logoutUserHandler)  // Updated to logoutUserHandler
	api.Get("/me", getCurrentUserHandler) // New endpoint to get current user
	// api.Put("/account", updateUserAccount) // Placeholder, will need auth middleware

	// Test endpoint
	api.Get("/hello", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello from the Go backend!",
		})
	})

	log.Fatal(app.Listen(":3000"))
}

// Registration request structure
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func registerUserHandler(c *fiber.Ctx) error {
	req := new(RegisterRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Basic validation (in a real app, use a validation library)
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username, email, and password are required",
		})
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Println("Error hashing password:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not process registration",
		})
	}

	user, err := data.AddUser(req.Username, req.Email, hashedPassword)
	if err != nil {
		if _, ok := err.(*data.UserExistsError); ok {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		log.Println("Error adding user:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not create user",
		})
	}

	// Don't send password hash back to client
	// Create a user response struct if you need more control over the output
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"user_id": user.ID,
		"username": user.Username,
		"email": user.Email,
	})
}

// Login request structure
type LoginRequest struct {
	Username string `json:"username"` // Can be username or email
	Password string `json:"password"`
}

func loginUserHandler(c *fiber.Ctx) error {
	req := new(LoginRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Username and password are required"})
	}

	user, err := data.GetUserByUsername(req.Username)
	if err != nil {
		// Try fetching by email if username not found
		user, err = data.GetUserByEmail(req.Username)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
		}
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	sess, err := store.Get(c)
	if err != nil {
		log.Println("Session error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create session"})
	}

	sess.Set("userID", user.ID)
	if err := sess.Save(); err != nil {
		log.Println("Session save error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not save session"})
	}

	return c.JSON(fiber.Map{
		"message": "Logged in successfully",
		"user": fiber.Map{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

func logoutUserHandler(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		// If no session, effectively logged out
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Logged out"})
	}

	if err := sess.Destroy(); err != nil {
		log.Println("Session destroy error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not log out"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Logged out successfully"})
}

func getCurrentUserHandler(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Session error"})
	}

	userID := sess.Get("userID")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Not authenticated"}) // Or return null/empty user
	}

	uid, ok := userID.(int64)
	if !ok {
	    // This case should ideally not happen if userID is always stored as int64
	    log.Printf("Session userID is not int64: %T\n", userID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid session data"})
	}

	user, err := data.GetUserByID(uid)
	if err != nil {
		// User might have been deleted after session was created
		sess.Destroy() // Clean up invalid session
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not found or invalid session"})
	}

	return c.JSON(fiber.Map{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
	})
}

// Placeholder handlers - Implement actual logic
// func updateUserAccount(c *fiber.Ctx) error { // Placeholder - requires auth middleware
// 	return c.SendString("Update user account placeholder")
// }
