package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"backend/judge"
	"backend/server"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	godotenv.Load(".env", ".envrc")
	db := openDb()
	defer db.Close()

	handler := log.New(os.Stderr)
	logger := slog.New(handler)

	j := judge.CreateJudge(*logger)
	s := server.InitServer(server.NewServerArgs{
		Logger: logger,
		Db:     db,
		Judge:  j,
	})

	// Code with a syntax error in Go (e.g., missing closing parenthesis)
	code := "console.log('Hello World')"

	input := "" // empty input
	response, err := j.PistonRunner("js", code, input)
	if err != nil {
		// If there is a compilation error, it will be reported here
		fmt.Println("Error:", err)
		return
	}

	// If no error, output results
	fmt.Printf("Output:\n%s\n", response.Run.Stdout)
	s.Run()
}

func openDb() *sql.DB {
	dbString := os.Getenv("GOOSE_DBSTRING")
	db, err := sql.Open("sqlite3", dbString)
	if err != nil {
		panic(fmt.Sprint("Could not open db", err.Error()))
	}

	err = db.Ping()
	if err != nil {
		panic(fmt.Sprint("Could not ping db", err.Error()))
	}

	return db
}
