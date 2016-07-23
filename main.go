package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/flynn/flynn/pkg/postgres"
)

func main() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)

	db := postgres.Wait(nil, nil)

	m := postgres.NewMigrations()
	m.Add(1, "CREATE SEQUENCE hits")
	if err := m.Migrate(db); err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		var count int
		if err := db.QueryRow("SELECT nextval('hits')").Scan(&count); err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
		fmt.Fprintf(w, "Hello from Flynn on port %s from container %s\nHits = %d\n", port, os.Getenv("HOSTNAME"), count)
	})
	fmt.Println("hitcounter listening on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
