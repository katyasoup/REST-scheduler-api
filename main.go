package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

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

func getAllUsers() []User {
	return getUsers("SELECT * FROM public.users")
}

func getUserByID(id int64) User {
	return getUser(fmt.Sprintf("SELECT * FROM public.users WHERE id=%d", id))
}

func getAllEmployees() []User {
	return getUsers("SELECT * FROM public.users WHERE role='employee'")
}

func getEmployeeByID(id int64) User {
	return getUser(fmt.Sprintf("SELECT * FROM public.users WHERE role='employee' AND id=%d", id))
}

func getAllManagers() []User {
	return getUsers("SELECT * FROM public.users WHERE role='manager'")
}

func getManagerByID(id int64) User {
	return getUser(fmt.Sprintf("SELECT * FROM public.users WHERE role='manager' AND id=%d", id))
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

func getAllShifts() []Shift {
	return getShifts("SELECT * FROM public.shifts")
}

func getShiftByID(id int64) Shift {
	return getShift(fmt.Sprintf("SELECT * FROM public.shifts WHERE id=%d", id))
}

func getMyShifts(id int64) []Shift {
	return getShifts(fmt.Sprintf("SELECT * FROM public.shifts WHERE employee_id=%d", id))
}

func getSchedule(start string, end string) []Shift {
	return getShifts(fmt.Sprintf("SELECT * FROM public.shifts WHERE start_time>'%s' AND end_time<'%s'", start, end))
}

//
func createShift(shift Shift) Shift {
	queryString := fmt.Sprintf("INSERT INTO public.shifts(manager_id, break, start_time, end_time) VALUES(%d, %f, '%s', '%s');",
		shift.Manager, shift.Break, shift.Start, shift.End)
	fmt.Println(queryString)
	rows, err := db.Query(queryString)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	rows.Close()
	// ideally will return most recently created shift; below code doesn't quite work
	return getShift(fmt.Sprintf("SELECT * FROM public.shifts WHERE id=MAX"))
}

func scheduleEmployee(shift Shift) Shift {
	queryString := fmt.Sprintf("UPDATE shifts SET employee_id =%d, updated_at=now() WHERE id = %d;",
		shift.Employee.Int64, shift.ID)
	fmt.Println(queryString)
	rows, err := db.Query(queryString)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	rows.Close()
	// ideally will return most recently updated shift; below code doesn't quite work
	return getShift(fmt.Sprintf("SELECT * FROM public.shifts WHERE updated_at=MAX"))
}

func editShiftTime(shift Shift) Shift {
	queryString := fmt.Sprintf("UPDATE shifts SET start_time='%s', end_time='%s', updated_at=now() WHERE id=%d;",
		shift.Start, shift.End, shift.ID)
	fmt.Println(queryString)
	rows, err := db.Query(queryString)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	rows.Close()
	// ideally will return most recently updated shift; below code doesn't quite work
	return getShift(fmt.Sprintf("SELECT * FROM public.shifts WHERE updated_at=MAX"))
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
	ID       int64     `json:"id"`
	Manager  int64     `json:"manager"`
	Employee NullInt64 `json:"employee"`
	Break    float64   `json:"break"`
	Start    string    `json:"startTime"`
	End      string    `json:"endTime"`
	Created  string    `json:"createdAt"`
	Updated  string    `json:"updatedAt"`
}

func stringToInt64(str string) int64 {
	id, err := strconv.ParseInt(str, 0, 64)
	if err != nil {
		panic(err)
	}
	return id
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

	// SHIFT routes:

	// get all shifts
	routes.GET("/shifts", func(c *gin.Context) {
		results := getAllShifts()
		c.JSON(200, results)
	})

	// get single shift by id
	routes.GET("/shifts/:id", func(c *gin.Context) {
		id := stringToInt64(c.Param("id"))
		result := getShiftByID(id)
		c.JSON(200, result)
	})

	// get all shifts for single employee
	routes.GET("/myshifts/:id", func(c *gin.Context) {
		id := stringToInt64(c.Param("id"))
		results := getMyShifts(id)
		c.JSON(200, results)
	})

	// get all shifts for date range
	routes.GET("/schedule/:start/:end", func(c *gin.Context) {
		start := c.Params.ByName("start")
		end := c.Params.ByName("end")
		results := getSchedule(start, end)
		c.JSON(200, results)
	})

	// add new shift
	routes.POST("/shifts", func(c *gin.Context) {
		var shift Shift
		c.BindJSON(&shift)
		result := createShift(shift)
		if err != nil {
			c.JSON(500, gin.H{"Error": err})
		} else {
			c.JSON(201, gin.H{"success": result})
		}
	})

	// change shift time
	routes.PUT("/shifts", func(c *gin.Context) {
		var shift Shift
		c.BindJSON(&shift)
		result := editShiftTime(shift)
		if err != nil {
			c.JSON(500, gin.H{"Error": err})
		} else {
			c.JSON(201, gin.H{"success": result})
		}
	})

	// add employee to shift
	routes.PUT("/shifts/assign", func(c *gin.Context) {
		var shift Shift
		c.BindJSON(&shift)
		result := scheduleEmployee(shift)
		if err != nil {
			c.JSON(500, gin.H{"Error": err})
		} else {
			c.JSON(201, gin.H{"success": result})
		}
	})

	// USER routes:

	// get all users
	routes.GET("/users", func(c *gin.Context) {
		results := getAllUsers()
		c.JSON(200, results)
	})

	// get single user by id
	routes.GET("/users/:id", func(c *gin.Context) {
		id := stringToInt64(c.Param("id"))
		result := getUserByID(id)
		c.JSON(200, result)
	})

	// get all users with role of employee
	routes.GET("/employees", func(c *gin.Context) {
		results := getAllEmployees()
		c.JSON(200, results)
	})

	// get single employee by id
	routes.GET("/employees/:id", func(c *gin.Context) {
		id := stringToInt64(c.Param("id"))
		result := getEmployeeByID(id)
		c.JSON(200, result)
	})

	// get all users with role of manager
	routes.GET("/managers", func(c *gin.Context) {
		results := getAllManagers()
		c.JSON(200, results)
	})

	// get single manager by id
	routes.GET("/managers/:id", func(c *gin.Context) {
		id := stringToInt64(c.Param("id"))
		result := getManagerByID(id)
		c.JSON(200, result)
	})

	routes.Run()
}
