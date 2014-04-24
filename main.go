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
	db, err := postgres.Open("", "")
	if err != nil {
		log.Fatal(err)
	}

	m := migrate.NewMigrations()
	m.Add(1,
		`CREATE TABLE hits (count int)`,
		`INSERT INTO hits (count) VALUES (0)`,
	)
	if err := m.Migrate(db.DB); err != nil {
		log.Fatal(err)
	}

	stmt, err := db.Prepare("UPDATE hits SET count = count+1 RETURNING count")
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
		fmt.Fprintf(w, "Hello from Go on Flynn: port=%s hits=%d", port, count)
	})
	fmt.Println("hitcounter listening on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
