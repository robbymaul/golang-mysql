package testing

import (
	"context"
	"log"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/robbymaul/golang-mysql.git/connection"
	"github.com/robbymaul/golang-mysql.git/model"
	"github.com/stretchr/testify/assert"
)

func TestSQL(t *testing.T) {
	tests := []struct {
		Test string
		Sql  string
		User model.User
	}{
		{
			Test: "INSERT",
			Sql:  "INSERT INTO users (username, password) values (?,?)",
			User: model.User{
				Username: "robby",
				Password: "robby",
			},
		},
		{
			Test: "UPDATE",
			Sql:  "UPDATE users SET username=?, password=? where username =?",
			User: model.User{
				Username: "robby",
				Password: "robby",
			},
		},
		{
			Test: "SELECT",
			Sql:  "SELECT username, password from users WHERE username = ?",
			User: model.User{
				Username: "robby",
			},
		},
		{
			Test: "DELETE",
			Sql:  "DELETE FROM users",
		},
	}

	for _, test := range tests {
		t.Run(test.Test, func(t *testing.T) {
			db := connection.GetConnection()
			defer db.Close()

			ctx := context.Background()

			query := test.Sql

			if query == "INSERT INTO users (username, password) values (?,?)" {
				tx, err := db.Begin()
				if err != nil {
					log.Fatal(err.Error())
				}

				_, err = tx.ExecContext(ctx, query, test.User.Username, test.User.Password)
				if err != nil {
					log.Fatal(err.Error())
				}

				assert.Nil(t, err)
				err = tx.Commit()
				if err != nil {
					log.Fatal(err.Error())
				}

			} else if query == "UPDATE users SET username=?, password=? where username =?" {
				tx, err := db.Begin()
				if err != nil {
					log.Fatal(err.Error())
				}

				tx.ExecContext(ctx, query, test.User.Username, test.User.Password, test.User.Username)

				err = tx.Commit()
				if err != nil {
					log.Fatal(err.Error())
				}

				assert.Nil(t, err)
			} else if query == "SELECT username, password FROM users WHERE username = ?" {
				var user model.User
				rows, err := db.QueryContext(ctx, query, test.User.Username)

				if err != nil {
					log.Fatal(err.Error())
				}

				if rows.Next() {
					err := rows.Scan(&user.Username, &user.Password)
					if err != nil {
						log.Fatal(err.Error())
					}
				}

				assert.Nil(t, err)
				assert.Equal(t, "robby", user.Username)
				defer rows.Close()
			} else if query == "DELETE FROM users" {
				tx, err := db.Begin()
				if err != nil {
					log.Fatal(err.Error())
				}

				_, err = tx.ExecContext(ctx, query)
				if err != nil {
					log.Fatal(err.Error())
				}

				err = tx.Commit()

				assert.Nil(t, err)
			}
		})
	}
}
