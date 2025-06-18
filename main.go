package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Deba00407/basic-web-server/controllers"
	"github.com/Deba00407/basic-web-server/database"
	schemamodels "github.com/Deba00407/basic-web-server/schema-models"
	"golang.org/x/crypto/bcrypt"
)

type FormData struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email-address"`
	Password string `json:"password"`
}

type APIMessage struct {
	Message    string
	statusCode int
}

func homePageHandler(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	if req.Method != "GET" {
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	urlPath := req.URL.Path

	if urlPath != "/" {
		http.Error(res, "404 Not Found", http.StatusNotFound)
		return
	}

	http.ServeFile(res, req, "./static/index.html")
}

func serverForm(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "./static/Form.html")
}

func (formData FormData) handleDatabaseUpload() *APIMessage {
	newUser := schemamodels.User{
		Name:     formData.Name,
		Username: formData.Username,
		Email:    formData.Email,
		Password: formData.Password,
	}

	_, err := controllers.RegisterUser(newUser)
	if err != nil {
		log.Println("User could not be created")
		return &APIMessage{Message: err.Error(), statusCode: http.StatusInternalServerError}
	}

	log.Println("User creation successful")
	return &APIMessage{Message: "User created successfully", statusCode: http.StatusCreated}
}

func formhandler(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	if req.Method != "POST" {
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := req.ParseForm(); err != nil {
		http.Error(res, "Error parsing form", 400)
	}

	formData := &FormData{
		Name:     req.FormValue("name"),
		Email:    req.FormValue("email"),
		Username: req.FormValue("username"),
		Password: req.FormValue("password"),
	}

	hashed_Password, err := bcrypt.GenerateFromPassword([]byte(formData.Password), bcrypt.DefaultCost)
	if err != nil {
		panic("Error hashing user password")
	}

	formData.Password = string(hashed_Password)

	response := formData.handleDatabaseUpload()

	// show a custom error page to the user
	errTemp, err := template.ParseFiles("./templates/error-template.html")
	if err != nil {
		panic(err)
	}

	res.WriteHeader(response.statusCode)

	// show an error page in-case an error occurs
	if response.statusCode > 210 {
		errTemp.Execute(res, response.Message)
		return
	}

	res.Write([]byte(response.Message))
}

func listAllRegisteredUsers(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	// check for right method
	if req.Method != "GET" {
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	registeredUser, err := controllers.GetAllRegisteredUsers()
	if err != nil {
		// show error-page
		errTemp, _ := template.ParseFiles("./templates/error-template.html")
		res.WriteHeader(http.StatusInternalServerError)
		errTemp.Execute(res, err.Error())
		return
	}

	// If there was no error
	usersTemp, _ := template.ParseFiles("./templates/registered-users.html")
	res.WriteHeader(http.StatusFound)
	usersTemp.Execute(res, registeredUser)
}

func main() {
	http.HandleFunc("/", homePageHandler)
	http.HandleFunc("/form", serverForm)
	http.HandleFunc("/register", formhandler)
	http.HandleFunc("/registered-users", listAllRegisteredUsers)

	server := &http.Server{
		Addr:    ":5001",
		Handler: nil,
	}

	// Make a DB connection
	database.MakeConnectionToDB()

	func() {
		if database.Collection == nil {
			log.Fatal("Database connection failed")
		}
	}()

	go func() {
		log.Println("Server running on PORT: 5001")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for shutdown signal
	gracefulShutDown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutDown, syscall.SIGINT, syscall.SIGTERM)

	<-gracefulShutDown

	log.Println("Shutdown signal received...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

	log.Println("Server shutdown gracefully ✨")
}

// func (formData FormData) handleRegisteredUsers() *APIMessage {
// 	var users []FormData
// 	var file *os.File
// 	var err error

// 	filePath := "./registeredUsers.json"

// 	// Check if file exists
// 	if _, err := os.Stat(filePath); err == nil {
// 		existingUsers, err := os.ReadFile(filePath)
// 		if err != nil {
// 			panic(fmt.Sprintf("Error occurred while reading JSON file: %v", err))
// 		}

// 		if len(existingUsers) > 0 {
// 			if err := json.Unmarshal(existingUsers, &users); err != nil {
// 				panic(fmt.Sprintf("Error unmarshalling JSON form data: %v", err))
// 			}

// 			// check for unique users
// 			for i := 0; i < len(users); i++ {
// 				if users[i].Username == formData.Username || users[i].Email == formData.Email {
// 					return &APIMessage{Message: "User already exists", statusCode: 400}
// 				}
// 			}
// 		}

// 		// Open file for writing (overwrite mode)
// 		file, err = os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0644)
// 		if err != nil {
// 			panic(fmt.Sprintf("Error opening file for writing: %v", err))
// 		}

// 	} else {
// 		// File does not exist → Create new file
// 		file, err = os.Create(filePath)
// 		if err != nil {
// 			panic(fmt.Sprintf("Error creating file: %v", err))
// 		}
// 	}

// 	// Always defer close
// 	defer func() {
// 		if err := file.Close(); err != nil {
// 			panic(fmt.Sprintf("Error closing file: %v", err))
// 		}
// 	}()

// 	// Append new user data
// 	users = append(users, formData)

// 	// Convert to JSON
// 	jsonData, err := json.MarshalIndent(users, "", "\t")
// 	if err != nil {
// 		panic(fmt.Sprintf("Error converting form data to JSON: %v", err))
// 	}

// 	// Write to file
// 	bytesWritten, err := file.Write(jsonData)
// 	if err != nil {
// 		panic(fmt.Sprintf("Error writing JSON to file: %v", err))
// 	}

// 	fmt.Printf("Data updated successfully, bytes written: %d\n", bytesWritten)

// 	return &APIMessage{Message: "User registered successfully", statusCode: 201}
// }
