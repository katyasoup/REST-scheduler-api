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

// User : for retrieval of user rows from db
type User struct {
	ID      int32
	Name    string
	Role    string
	Email   string
	Phone   string
	Created string
	Updated string
}

// Shift : for retrieval of shift rows from db
type Shift struct {
	ID       int32
	Manager  int32
	Employee NullInt64
	Break    float32
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
		shiftID        int32
		managerID      int32
		employeeID     NullInt64
		breakTime      float32
		shiftStart     string
		shiftEnd       string
		shiftCreatedAt string
		shiftUpdatedAt string
		shiftResults   []Shift
	)

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

	r := gin.Default()

	r.GET("/Shifts", func(c *gin.Context) {
		rows, err := db.Query("SELECT * FROM public.shifts")
		defer rows.Close()

		for rows.Next() {
			rows.Scan(&shiftID, &managerID, &employeeID, &breakTime, &shiftStart, &shiftEnd, &shiftCreatedAt, &shiftUpdatedAt)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(shiftID, managerID, employeeID, breakTime, shiftStart, shiftEnd, shiftCreatedAt, shiftUpdatedAt)
			shiftResults = append(shiftResults, Shift{shiftID, managerID, employeeID, breakTime, shiftStart, shiftEnd, shiftCreatedAt, shiftUpdatedAt})
		}
		c.JSON(200, shiftResults)
	})

	r.GET("/Users", func(c *gin.Context) {
		rows, err := db.Query("SELECT * FROM public.users")
		defer rows.Close()

		for rows.Next() {
			rows.Scan(&id, &name, &role, &email, &phone, &createdAt, &updatedAt)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(id, name, role, email, phone, createdAt, updatedAt)
			results = append(results, User{id, name, role, email, phone, createdAt, updatedAt})
		}
		c.JSON(200, results)

		c.JSON(200, "Hook up shifts query!")
	})

	r.GET("/Employees", func(c *gin.Context) {
		rows, err := db.Query("SELECT * FROM public.users WHERE role='employee'")
		defer rows.Close()

		for rows.Next() {
			rows.Scan(&id, &name, &role, &email, &phone, &createdAt, &updatedAt)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(id, name, role, email, phone, createdAt, updatedAt)
			results = append(results, User{id, name, role, email, phone, createdAt, updatedAt})
		}
		c.JSON(200, results)

		c.JSON(200, "Hook up shifts query!")
	})

	r.Run()
}
