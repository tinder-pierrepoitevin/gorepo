package main

import (
	"database/sql"
	"fmt"
	"syscall"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/ssh/terminal"
)

/*
 * Task structure
 */
type Task struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	fmt.Print("Enter Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Printf("Failed to parse password!\n")
		return
	}
	fmt.Printf("\n")
	password := string(bytePassword)
	db, err := sql.Open("mysql", fmt.Sprintf("root:%s@tcp(127.0.0.1:3306)/Test", password))
	defer db.Close()
	if err != nil {
		fmt.Printf("Failed!\n")
		return
	}
	err = db.Ping()
	if err != nil {
		fmt.Printf("Failed Ping!\n")
		return
	}
	// connection established :)

	results, err := db.Query("SELECT id, name FROM Task")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	for results.Next() {
		var task Task
		// for each row, scan the result into our tag composite object
		err = results.Scan(&task.ID, &task.Name)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		// and then print out the tag's Name attribute
		fmt.Printf("task: id %d, name %s\n", task.ID, task.Name)
	}

}
