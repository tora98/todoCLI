package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const database = "todo.db"

type Todo struct {
	Description   string
	Completed     bool
	DateCreated   string
	DateCompleted string
}

func createDatabase() error {
	connection, err := sql.Open("sqlite3", database)
	if err != nil {
		return err
	}

	defer connection.Close()

	sqlStmnt := `CREATE TABLE IF NOT EXISTS todo (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    description TEXT,
    completed BOOLEAN NOT NULL DEFAULT (0),
		date_created string,
		date_completed string
  );`
	_, err = connection.Exec(sqlStmnt)
	if err != nil {
		return err
	} else {
		fmt.Println("Database connected!")
		return nil
	}
}

func addTodo(description string) error {
	if len(description) == 0 {
		return fmt.Errorf("empty description")
	}
	connection, err := sql.Open("sqlite3", database)
	if err != nil {
		return err
	}

	defer connection.Close()

	t := time.Now()

	todo := Todo{
		Description:   description,
		Completed:     false,
		DateCreated:   t.Format(time.RFC3339),
		DateCompleted: "",
	}

	sqlStmt := `INSERT INTO todo (
		description,
		completed,
		date_created,
		date_completed
		) VALUES (?, ?, ?, ?)`

	_, err = connection.Exec(sqlStmt, todo.Description, todo.Completed, todo.DateCreated, todo.DateCompleted)
	if err != nil {
		return err
	}
	fmt.Println("Todo added successfully!")
	return nil
}

func deleteTodo(todoId int) error {
	connection, err := sql.Open("sqlite3", database)
	if err != nil {
		return err
	}

	defer connection.Close()

	sqlStmt := `DELETE FROM todo WHERE id = ?`

	result, err := connection.Exec(sqlStmt, todoId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no todo found with ID = %d", todoId)
	}

	fmt.Println("Todo deleted successfully!")
	return nil
}

func showTodo() error {
	connection, err := sql.Open("sqlite3", database)
	if err != nil {
		log.Fatal(err)
	}

	defer connection.Close()

	sqlStmt := "SELECT id, description, completed, date_created, date_completed FROM todo"

	rows, err := connection.Query(sqlStmt)
	if err != nil {
		return err
	}

	defer rows.Close()

	fmt.Println("ID | DESCRIPTION | COMPLETED | DATE CREATED | DATE COMPLETED")
	fmt.Println("-----------------------------------------------------------")

	count := 0
	for rows.Next() {
		count++
		var id int
		var description string
		var completed bool
		var date_created string
		var date_completed string

		err = rows.Scan(&id, &description, &completed, &date_created, &date_completed)
		if err != nil {
			return err
		}

		status := "❌"
		if completed {
			status = "✅"
		}

		fmt.Printf("%d | %s | %s | %s | %s\n", id, description, status, date_created, date_completed)
		fmt.Println("")
	}

	if count == 0 {
		fmt.Println("No Todos Found!")
	}

	return nil
}

func completeTodo(todoId int) error {
	connection, err := sql.Open("sqlite3", database)
	if err != nil {
		return err
	}

	defer connection.Close()

	t := time.Now()
	dateCompleted := t.Format(time.RFC3339)

	sqlStmt := `UPDATE todo SET completed = true, date_completed = ? WHERE id = ?`
	result, err := connection.Exec(sqlStmt, dateCompleted, todoId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no todo found with todoId = %d", todoId)
	}

	fmt.Println("Todo Marked as Completed!")
	return nil
}

func showHelp() {
	fmt.Println("Todo CLI Application")
	fmt.Println("-------------------")
	fmt.Println("Commands:")
	fmt.Println("  add \"description\"       - Add a new todo item")
	fmt.Println("  complete id 		  - Mark a todo as completed")
	fmt.Println("  delete id               - Delete a todo item")
	fmt.Println("  show                    - Display all todos")
	fmt.Println("  help                    - Show this help message")
	fmt.Println("")
}

func main() {
	err := createDatabase()
	if err != nil {
		log.Fatal(err)
	}

	args := os.Args
	arglen := len(args)

	if arglen == 1 {
		showHelp()
		return
	}

	command := strings.ToLower(args[1])
	if command == "show" {
		err = showTodo()
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	if arglen < 3 {
		fmt.Println("Insufficient arguments!")
		showHelp()
		return
	}

	arg := args[2]

	switch command {
	case "add":
		err = addTodo(arg)
		if err != nil {
			fmt.Println(err)
			return
		}
	case "delete":
		id, err := strconv.Atoi(arg)
		if err != nil {
			log.Fatal("Invalid ID: must be a number!")
			fmt.Println(err)
			return
		}
		err = deleteTodo(id)
		if err != nil {
			fmt.Println(err)
			return
		}
	case "complete":
		id, err := strconv.Atoi(arg)
		if err != nil {
			log.Fatal("Invalid ID: must be a number!")
		}
		err = completeTodo(id)
		if err != nil {
			fmt.Println(err)
			return
		}
	default:
		fmt.Println("Invalid command!")
		showHelp()
	}

	if err != nil {
		log.Fatal(err)
	}
}
