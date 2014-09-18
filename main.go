package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/flynn/go-flynn/migrate"
	"github.com/flynn/go-flynn/postgres"
)

func main() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)

	db, err := postgres.Open("", "")
	if err != nil {
		log.Fatal(err)
	}

	m := migrate.NewMigrations()
	m.Add(1, "CREATE SEQUENCE hits")
	if err := m.Migrate(db.DB); err != nil {
		log.Fatal(err)
	}

	stmt, err := db.Prepare("SELECT nextval('hits')")
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		var count int
		if err := stmt.QueryRow().Scan(&count); err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
		fmt.Fprintf(w, "Hello from Go+PostgreSQL on Flynn: port=%s hits=%d container=%s\n", port, count, os.Getenv("HOSTNAME"))
	})
	fmt.Println("hitcounter listening on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
