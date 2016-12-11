package main

import (
	"database/sql"
	"log"
	"math/rand"
	"net/http"

	"github.com/graphql-go/graphql"
	gqlhandler "github.com/graphql-go/graphql-go-handler"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

var queryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"latestPost": &graphql.Field{
			Type: graphql.String,
			//Always has to have graphql.ResolveParams and return interface{} and
			//error.  Ugh, out of date tutorials stink
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				stmt, err := db.Prepare("SELECT name FROM movies WHERE uid = ?")
				if err != nil {
					log.Fatal(err)
				}
				defer stmt.Close()

				var name string
				err = stmt.QueryRow("1").Scan(&name)
				if err != nil {
					log.Fatal(err)
				}

				return name, nil
			},
		},
		"randInt": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return rand.Intn(50), nil
			},
		},
		// "otherInts": &graphql.Fields{
		// 	"int1": &graphql.Field{
		// 		Type: graphql.Int,
		// 		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		// 			return rand.Intn(25), nil
		// 		},
		// 	},
		// 	"int2": &graphql.Field{
		// 		Type: graphql.Int,
		// 		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		// 			return rand.Intn(10), nil
		// 		},
		// 	},
		// 	"int3": &graphql.Field{
		// 		Type: graphql.Int,
		// 		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		// 			return rand.Intn(5000), nil
		// 		},
		// 	},
		// },
	},
})

var Schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: queryType,
})

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./simpleDB")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	handler := gqlhandler.New(&gqlhandler.Config{
		Schema: &Schema,
		Pretty: true,
	})

	http.Handle("/graphql", handler)

	http.ListenAndServe(":8080", nil)
}
