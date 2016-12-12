package main

import (
	"net/http"

	gqlhandler "github.com/graphql-go/graphql-go-handler"
	"github.com/mcquackers/learn-graphql/learnSchema"
)

func main() {
	defer learnSchema.Db.Close()

	handler := gqlhandler.New(&gqlhandler.Config{
		Schema: &learnSchema.Schema,
		Pretty: true,
	})

	http.Handle("/graphql", handler)

	http.ListenAndServe(":8080", nil)
}
