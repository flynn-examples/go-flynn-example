package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/flynn/go-discoverd"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func main() {
	ss, err := discoverd.NewServiceSet("mongo")
	if err != nil {
		log.Fatal(err)
	}
	mongo := <-ss.Watch(true)
	ss.Close()

	sess, err := mgo.Dial(mongo.Addr)
	if err != nil {
		log.Fatal(err)
	}
	coll := sess.DB("hits").C("counter")

	port := os.Getenv("PORT")

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fail := func(err error) {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
		mongoCount := struct {
			Hits int `bson:"hits"`
		}{}

		if _, err := coll.UpsertId("hits", bson.M{"$inc": bson.M{"hits": 1}}); err != nil {
			fail(err)
		}
		if err := coll.FindId("hits").One(&mongoCount); err != nil {
			fail(err)
		}

		fmt.Fprintf(w, "Hello from Go + MongoDB on Flynn: port=%s hits=%d\n", port, mongoCount.Hits)
	})
	fmt.Println("hitcounter listening on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
