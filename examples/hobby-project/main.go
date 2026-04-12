package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	queries := db.New(db)
	ctx := context.Background()

	fmt.Println("=== Hobby Project Example ===")
	fmt.Println()

	runExample("Creating a user", func() error {
		user, createErr := queries.CreateUser(ctx, db.CreateUserParams{
			ID:       "user-001",
			Email:    "john.doe@example.com",
			FullName: "John Doe",
		})
		if createErr != nil {
			return createErr
		}
		fmt.Printf("   Created user: %s <%s>\n", user.FullName, user.Email)
		return nil
	})

	runExample("Getting user by ID", func() error {
		user, getErr := queries.GetUserByID(ctx, "user-001")
		if getErr != nil {
			return getErr
		}
		fmt.Printf("   Found user: %s <%s>\n", user.FullName, user.Email)
		return nil
	})

	runExample("Listing all users", func() error {
		users, listErr := queries.ListUsers(ctx, db.ListUsersParams{
			Limit:  int32(10),
			Offset: int32(0),
		})
		if listErr != nil {
			return listErr
		}
		fmt.Printf("   Found %d users:\n", len(users))
		for _, u := range users {
			fmt.Printf("   - %s <%s>\n", u.FullName, u.Email)
		}
		return nil
	})

	runExample("Updating user", func() error {
		user, updateErr := queries.UpdateUser(ctx, db.UpdateUserParams{
			ID:       "user-001",
			FullName: "John Updated",
		})
		if updateErr != nil {
			return updateErr
		}
		fmt.Printf("   Updated user: %s <%s>\n", user.FullName, user.Email)
		return nil
	})

	runExample("Searching users", func() error {
		users, searchErr := queries.SearchUsersByEmail(ctx, db.SearchUsersByEmailParams{
			Email: "%john%",
		})
		if searchErr != nil {
			return searchErr
		}
		fmt.Printf("   Found %d users matching 'john':\n", len(users))
		for _, u := range users {
			fmt.Printf("   - %s <%s>\n", u.FullName, u.Email)
		}
		return nil
	})

	runExample("Deleting user", func() error {
		return queries.DeleteUser(ctx, "user-001")
	})

	fmt.Println("=== All Examples Completed ===")
	fmt.Println()
	fmt.Println("Database file: ./database.db")
	fmt.Println("Generated code: ./internal/db/")
	fmt.Println()
	fmt.Println("To regenerate code:")
	fmt.Println("  sqlc generate")
}

func runExample(name string, fn func() error) {
	fmt.Printf("%s...\n", name)
	if err := fn(); err != nil {
		log.Printf("   Error: %v", err)
	}
	fmt.Println()
}
