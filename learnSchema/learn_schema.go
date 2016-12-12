package learnSchema

import (
	"database/sql"
	"math/rand"

	"github.com/graphql-go/graphql"
	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB

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
		"movie": &graphql.Field{
			Type: graphql.NewList(movieType),
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
				var stmtString string = "SELECT * FROM movies"
				var movies []Movie
				if id, ok = p.Args["id"].(int); !ok {
					stmt, err := Db.Prepare(stmtString)
					if err != nil {
						return nil, err
					}
					defer stmt.Close()

					rows, err := Db.Query(stmtString)
					if err != nil {
						return nil, err
					}
					defer rows.Close()

					for rows.Next() {
						var rowMovie Movie
						err = rows.Scan(&rowMovie.Uid, &rowMovie.Name)
						if err != nil {
							return nil, err
						}
						movies = append(movies, rowMovie)
					}
					return movies, nil
				}
				stmtString += " WHERE uid = ?"
				stmt, err := Db.Prepare(stmtString)
				if err != nil {
					return nil, err
				}
				defer stmt.Close()
				var movie Movie

				err = stmt.QueryRow(id).Scan(&movie.Uid, &movie.Name)
				if err != nil {
					return nil, err
				}
				movies = append(movies, movie)

				return movies, nil
			},
		},

		// var id int
		// var ok bool
		// if id, ok = p.Args["id"].(int); !ok {
		// 	id = 1
		// } else {
		// 	stmtString += " WHERE uid = ?"
		// }
		// stmt, err := Db.Prepare(stmtString)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// defer stmt.Close()
		//
		// var movie Movie
		// err = stmt.QueryRow(id).Scan(&movie.Uid, &movie.Name)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		//
		// return movie, nil
		// },
		// },
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

func init() {
	var err error
	Db, err = sql.Open("sqlite3", "./simpleDB")
	if err != nil {
		panic(err)
	}
}
