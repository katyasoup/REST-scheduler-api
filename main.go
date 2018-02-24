package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "Katie"
	dbname = "wiw-challenge"
)

type User struct {
	ID         int32
	Name       string
	Role       string
	Email      string
	Phone      string
	Created_at string
	Updated_at string
}

func main() {
	// db setup
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"dbname=%s sslmode=disable",
		host, port, user, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to the database!!!")

	var (
		id        int32
		name      string
		role      string
		email     string
		phone     string
		createdAt string
		updatedAt string
		results   []User
	)

	rows, err := db.Query("SELECT * FROM public.users")
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &name, &role, &email, &phone, &createdAt, &updatedAt)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, name, role, email, phone, createdAt, updatedAt)
		results = append(results, User{id, name, role, email, phone, createdAt, updatedAt})
	}

	r := gin.Default()
	r.GET("/Users", func(c *gin.Context) {
		c.JSON(200, results)
	})
	r.Run()
}
