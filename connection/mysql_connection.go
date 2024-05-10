package connection

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func GetConnection() *sql.DB {
	err := godotenv.Load("../.env.dev")
	if err != nil {
		log.Fatal(err.Error())
	}

	db_host := os.Getenv("DB_HOST")
	db_port := os.Getenv("DB_PORT")
	db_user := os.Getenv("DB_USER")
	db_password := os.Getenv("DB_PASSWORD")
	db_name := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", db_user, db_password, db_host, db_port, db_name)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err.Error())
	}

	db.SetConnMaxLifetime(3 * time.Minute)
	db.SetConnMaxIdleTime(10)
	db.SetMaxIdleConns(50)

	return db
}
