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

type Movie struct {
	Uid  int    `json:"uid"`
	Name string `json:"name"`
}

var movieType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Movie",
	Fields: graphql.Fields{
		"uid": &graphql.Field{
			Type: graphql.Int,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var queryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"latestPost": &graphql.Field{
			Type: movieType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			//Always has to have graphql.ResolveParams and return interface{} and
			//error.  Ugh, out of date tutorials stink
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var id int
				var ok bool
				if id, ok = p.Args["id"].(int); !ok {
					id = 1
				}
				stmt, err := db.Prepare("SELECT * FROM movies WHERE uid = ?")
				if err != nil {
					log.Fatal(err)
				}
				defer stmt.Close()

				var movie Movie
				err = stmt.QueryRow(id).Scan(&movie.Uid, &movie.Name)
				if err != nil {
					log.Fatal(err)
				}

				return movie, nil
			},
		},
		"randInt": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return rand.Intn(50), nil
			},
		},
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
