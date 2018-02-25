package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "Katie"
	dbname = "wiw-challenge"
)

// Global DB handles
var db *sql.DB
var err error

// NullInt64 : allow for null value in employee field for shifts
type NullInt64 sql.NullInt64

// Scan : check for null value and set Valid bool
func (ni *NullInt64) Scan(value interface{}) error {
	var i sql.NullInt64
	if err := i.Scan(value); err != nil {
		return err
	}

	// if nil then make Valid false
	if reflect.TypeOf(value) == nil {
		*ni = NullInt64{i.Int64, false}
	} else {
		*ni = NullInt64{i.Int64, true}
	}
	return nil
}

// MarshalJSON : allow for EITHER null or populated value in employee field for shifts
func (ni *NullInt64) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Int64)
}

// UnmarshalJSON for NullInt64
func (ni *NullInt64) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ni.Int64)
	ni.Valid = (err == nil)
	return err
}

// multiple user rows into slice of Users
func scanUsers(rows *sql.Rows) []User {
	var users []User
	var u User
	for rows.Next() {
		rows.Scan(&u.ID, &u.Name, &u.Role, &u.Email, &u.Phone, &u.Created, &u.Updated)
		fmt.Println(u)
		users = append(users, u)
	}
	return users
}

func getUsers(queryString string) []User {
	rows, err := db.Query(queryString)
	defer rows.Close()
	results := scanUsers(rows)
	if err != nil {
		log.Fatal(err)
	}
	rows.Close()
	return results
}

// multiple shift rows into slice of Shifts
func scanShifts(rows *sql.Rows) []Shift {
	var shifts []Shift
	var s Shift
	for rows.Next() {
		rows.Scan(&s.ID, &s.Manager, &s.Employee, &s.Break, &s.Start, &s.End, &s.Created, &s.Updated)
		fmt.Println(s)
		shifts = append(shifts, s)
	}
	return shifts
}

func getShifts(queryString string) []Shift {
	rows, err := db.Query(queryString)
	defer rows.Close()
	results := scanShifts(rows)
	if err != nil {
		log.Fatal(err)
	}
	rows.Close()
	return results
}

// User : for retrieval of user rows from db
type User struct {
	ID      int64
	Name    string
	Role    string
	Email   string
	Phone   string
	Created string
	Updated string
}

// Shift : for retrieval of shift rows from db
type Shift struct {
	ID       int64
	Manager  int64
	Employee NullInt64
	Break    float64
	Start    string
	End      string
	Created  string
	Updated  string
}

func main() {
	// db setup
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"dbname=%s sslmode=disable",
		host, port, user, dbname)
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to the database!!!")

	routes := gin.Default()

	// get all shifts
	routes.GET("/Shifts", func(c *gin.Context) {
		results := getShifts("SELECT * FROM public.shifts")
		c.JSON(200, results)
	})

	// get all users
	routes.GET("/Users", func(c *gin.Context) {
		results := getUsers("SELECT * FROM public.users")
		c.JSON(200, results)
	})

	// get all users with role of employee
	routes.GET("/Employees", func(c *gin.Context) {
		results := getUsers("SELECT * FROM public.users WHERE role='employee'")
		c.JSON(200, results)
	})

	// get single EMPLOYEE by id
	routes.GET("/Employees/:id", func(c *gin.Context) {
		var u User
		idParam := c.Params.ByName("id")
		queryString := fmt.Sprintf("SELECT * FROM public.users WHERE role='employee' AND id=%s", idParam)
		row := db.QueryRow(queryString)

		row.Scan(&u.ID, &u.Name, &u.Role, &u.Email, &u.Phone, &u.Created, &u.Updated)
		fmt.Println(u)
		c.JSON(200, u)
	})

	routes.Run()
}
