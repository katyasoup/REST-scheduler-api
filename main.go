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

// Global DB handles
var db *sql.DB
var err error

// map multiple user rows into slice of Users
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

// map single user row into User struct
func scanUser(row *sql.Row) User {
	var u User
	row.Scan(&u.ID, &u.Name, &u.Role, &u.Email, &u.Phone, &u.Created, &u.Updated)
	fmt.Println(u)
	return u
}

// retrieve results of scanUsers
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

// retrieve results of scanUser
func getUser(queryString string) User {
	row := db.QueryRow(queryString)
	results := scanUser(row)
	if err != nil {
		log.Fatal(err)
	}
	return results
}

// map multiple shift rows into slice of Shifts
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

// map single shift row into Shift struct
func scanShift(row *sql.Row) Shift {
	var s Shift
	row.Scan(&s.ID, &s.Manager, &s.Employee, &s.Break, &s.Start, &s.End, &s.Created, &s.Updated)
	fmt.Println(s)
	return s
}

// retrieve results of scanShifts
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

// retrieve results of scanShift
func getShift(queryString string) Shift {
	row := db.QueryRow(queryString)
	results := scanShift(row)
	if err != nil {
		log.Fatal(err)
	}
	return results
}

// User : for storing retrieval of user rows from db
type User struct {
	ID      int64
	Name    string
	Role    string
	Email   string
	Phone   string
	Created string
	Updated string
}

// Shift : for storing retrieval of shift rows from db
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

	// get SHIFTS:

	// get all shifts
	routes.GET("/shifts", func(c *gin.Context) {
		results := getShifts("SELECT * FROM public.shifts")
		c.JSON(200, results)
	})

	// get single shift by id
	routes.GET("/shifts/:id", func(c *gin.Context) {
		idParam := c.Params.ByName("id")
		queryString := fmt.Sprintf("SELECT * FROM public.shifts WHERE id=%s", idParam)
		result := getShift(queryString)
		c.JSON(200, result)
	})

	// get all shifts for single employee
	routes.GET("/myshifts/:id", func(c *gin.Context) {
		idParam := c.Params.ByName("id")
		queryString := fmt.Sprintf("SELECT * FROM public.shifts WHERE employee_id=%s", idParam)
		results := getShifts(queryString)
		c.JSON(200, results)
	})

	// get all shifts for date range
	routes.GET("/schedule/:start/:end", func(c *gin.Context) {
		startParam := c.Params.ByName("start")
		endParam := c.Params.ByName("end")
		queryString := fmt.Sprintf("SELECT * FROM public.shifts WHERE start_time>%s AND end_time<%s", startParam, endParam)
		results := getShifts(queryString)
		c.JSON(200, results)
	})

	// get USERS:

	// get all users
	routes.GET("/users", func(c *gin.Context) {
		results := getUsers("SELECT * FROM public.users")
		c.JSON(200, results)
	})

	// get single user by id
	routes.GET("/users/:id", func(c *gin.Context) {
		idParam := c.Params.ByName("id")
		queryString := fmt.Sprintf("SELECT * FROM public.users WHERE id=%s", idParam)
		result := getUser(queryString)
		c.JSON(200, result)
	})

	// get all users with role of employee
	routes.GET("/employees", func(c *gin.Context) {
		results := getUsers("SELECT * FROM public.users WHERE role='employee'")
		c.JSON(200, results)
	})

	// get single employee by id
	routes.GET("/employees/:id", func(c *gin.Context) {
		idParam := c.Params.ByName("id")
		queryString := fmt.Sprintf("SELECT * FROM public.users WHERE role='employee' AND id=%s", idParam)
		result := getUser(queryString)
		c.JSON(200, result)
	})

	// get all users with role of manager
	routes.GET("/managers", func(c *gin.Context) {
		results := getUsers("SELECT * FROM public.users WHERE role='manager'")
		c.JSON(200, results)
	})

	// get single manager by id
	routes.GET("/managers/:id", func(c *gin.Context) {
		idParam := c.Params.ByName("id")
		queryString := fmt.Sprintf("SELECT * FROM public.users WHERE role='manager' AND id=%s", idParam)
		result := getUser(queryString)
		c.JSON(200, result)
	})

	routes.Run()
}
