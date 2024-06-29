package main

import (
    "database/sql"
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "log"
    "os"
    "mini-project-sanbercode/controllers"
    "mini-project-sanbercode/database"

    _ "github.com/lib/pq"
)

var (
    DB  *sql.DB
    err error
)

// Fungsi untuk membuat database jika tidak ada
func createDatabaseIfNotExists(dbName, user, password, host, port string) error {
    connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable", host, port, user, password)
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return fmt.Errorf("failed to connect to PostgreSQL: %v", err)
    }
    defer db.Close()

    // Periksa apakah database ada
    var exists bool
    query := fmt.Sprintf("SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = '%s')", dbName)
    err = db.QueryRow(query).Scan(&exists)
    if err != nil {
        return fmt.Errorf("failed to check if database exists: %v", err)
    }

    // Buat database jika tidak ada
    if !exists {
        _, err := db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
        if err != nil {
            return fmt.Errorf("failed to create database: %v", err)
        }
        fmt.Println("Database created successfully")
    } else {
        fmt.Println("Database already exists")
    }

    return nil
}

func main() {
    // ENV Configuration
    err = godotenv.Load("config/.env")
    if err != nil {
        log.Fatalf("failed to load environment file: %v", err)
    }

    dbName := os.Getenv("DB_NAME")
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")

    // Periksa dan buat database jika tidak ada
    err = createDatabaseIfNotExists(dbName, dbUser, dbPassword, dbHost, dbPort)
    if err != nil {
        log.Fatalf("failed to create database: %v", err)
    }

    psqlInfo := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        dbHost, dbPort, dbUser, dbPassword, dbName,
    )

    DB, err = sql.Open("postgres", psqlInfo)
    if err != nil {
        log.Fatalf("failed to open database: %v", err)
    }

    err = DB.Ping()
    if err != nil {
        log.Fatalf("DB Connection Failed: %v", err)
    } else {
        fmt.Println("DB Connection Success")
    }

    database.DbMigrate(DB)
    defer DB.Close()

    // Router GIN
    router := gin.Default()
    router.GET("/persons", controllers.GetAllPerson)
    router.POST("/persons", controllers.InsertPerson)
    router.PUT("/persons/:id", controllers.UpdatePerson)
    router.DELETE("/persons/:id", controllers.DeletePerson)

    router.Run("localhost:8080")
}
