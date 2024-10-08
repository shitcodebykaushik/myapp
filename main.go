// main.go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// User struct to hold user data.
type User struct {
    Phone    string `json:"phone"`
    Password string `json:"password"`
    Username string `json:"username"`
}

// Database and Collection variables
var client *mongo.Client
var usersCollection *mongo.Collection

// InitializeMongoDB initializes the MongoDB client and users collection.
func InitializeMongoDB() {
    var err error
    // Set client options with your MongoDB connection string
    clientOptions := options.Client().ApplyURI("mongodb+srv://algorithmunloack:Adii77XmGVj99WZM@db.gaa7h.mongodb.net/?retryWrites=true&w=majority&appName=DB")

    // Connect to MongoDB
    client, err = mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        log.Fatalf("Failed to connect to MongoDB: %v\n", err)
    }

    // Check the connection
    err = client.Ping(context.TODO(), nil)
    if err != nil {
        log.Fatalf("Failed to ping MongoDB: %v\n", err)
    }

    fmt.Println("Connected to MongoDB!")
    usersCollection = client.Database("your_database_name").Collection("users") // Replace with your database name
}

// SignUp handler for user sign-up.
func SignUp(c *gin.Context) {
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Insert user into the database
    _, err := usersCollection.InsertOne(context.TODO(), user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "User already exists or database error."})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "User created successfully."})
}

// SignIn handler for user sign-in.
func SignIn(c *gin.Context) {
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Check if the user exists and password matches
    var foundUser User
    err := usersCollection.FindOne(context.TODO(), bson.M{"phone": user.Phone, "password": user.Password}).Decode(&foundUser)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid phone number or password."})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Sign-in successful", "username": foundUser.Username})
}

func main() {
    InitializeMongoDB() // Initialize the MongoDB client

    r := gin.Default()
    r.POST("/signup", SignUp)
    r.POST("/login", SignIn)
    r.Run() // listen and serve on 0.0.0.0:8080
}
