package main

import (
    "context"
    "database/sql"
    "fmt"
    "log"

    _ "github.com/mattn/go-sqlite3"
    "github.com/LarsArtmann/SQLC-Wizzard/examples/hobby-project/internal/db"
)

func main() {
    // Open database connection
    // SQLite creates file automatically if it doesn't exist
    db, err := sql.Open("sqlite3", "./database.db")
    if err != nil {
        log.Fatalf("Failed to open database: %v", err)
    }
    defer db.Close()

    // Use generated queries
    queries := db.New(db)
    ctx := context.Background()

    fmt.Println("=== Hobby Project Example ===")
    fmt.Println()

    // Example 1: Create a user
    fmt.Println("1. Creating a user...")
    user, err := queries.CreateUser(ctx, db.CreateUserParams{
        ID:       "user-001",
        Email:    "john.doe@example.com",
        FullName: "John Doe",
    })
    if err != nil {
        log.Printf("Error creating user: %v", err)
    } else {
        fmt.Printf("   Created user: %s <%s>\n", user.FullName, user.Email)
    }
    fmt.Println()

    // Example 2: Get user by ID
    fmt.Println("2. Getting user by ID...")
    user, err = queries.GetUserByID(ctx, "user-001")
    if err != nil {
        log.Printf("Error getting user: %v", err)
    } else {
        fmt.Printf("   Found user: %s <%s>\n", user.FullName, user.Email)
    }
    fmt.Println()

    // Example 3: List all users
    fmt.Println("3. Listing all users...")
    limit := int32(10)
    offset := int32(0)
    users, err := queries.ListUsers(ctx, db.ListUsersParams{
        Limit:  limit,
        Offset: offset,
    })
    if err != nil {
        log.Printf("Error listing users: %v", err)
    } else {
        fmt.Printf("   Found %d users:\n", len(users))
        for _, u := range users {
            fmt.Printf("   - %s <%s>\n", u.FullName, u.Email)
        }
    }
    fmt.Println()

    // Example 4: Update user
    fmt.Println("4. Updating user...")
    user, err = queries.UpdateUser(ctx, db.UpdateUserParams{
        ID:       "user-001",
        FullName: "John Updated",
    })
    if err != nil {
        log.Printf("Error updating user: %v", err)
    } else {
        fmt.Printf("   Updated user: %s <%s>\n", user.FullName, user.Email)
    }
    fmt.Println()

    // Example 5: Search users
    fmt.Println("5. Searching users...")
    users, err = queries.SearchUsersByEmail(ctx, db.SearchUsersByEmailParams{
        Email: "%john%",
    })
    if err != nil {
        log.Printf("Error searching users: %v", err)
    } else {
        fmt.Printf("   Found %d users matching 'john':\n", len(users))
        for _, u := range users {
            fmt.Printf("   - %s <%s>\n", u.FullName, u.Email)
        }
    }
    fmt.Println()

    // Example 6: Delete user
    fmt.Println("6. Deleting user...")
    err = queries.DeleteUser(ctx, "user-001")
    if err != nil {
        log.Printf("Error deleting user: %v", err)
    } else {
        fmt.Println("   Deleted user-001")
    }
    fmt.Println()

    fmt.Println("=== All Examples Completed ===")
    fmt.Println()
    fmt.Println("Database file: ./database.db")
    fmt.Println("Generated code: ./internal/db/")
    fmt.Println()
    fmt.Println("To regenerate code:")
    fmt.Println("  sqlc generate")
}
