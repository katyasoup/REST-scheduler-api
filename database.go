package main

import (
	"database/sql"
	"fmt"
	"log"
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

// db setup
// func OpenDatabase() {
// 	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
// 		"dbname=%s sslmode=disable",
// 		host, port, user, dbname)
// 	db, err = sql.Open("postgres", psqlInfo)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer db.Close()
// 	err = db.Ping()
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("Successfully connected to the database!!!")
// }

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

// map multiple roster rows into slice of Roster
func scanRoster(rows *sql.Rows) []Roster {
	var roster []Roster
	var r Roster
	for rows.Next() {
		rows.Scan(&r.ID, &r.Manager, &r.Employee, &r.Break, &r.Start, &r.End, &r.Created, &r.Updated, &r.Name, &r.Phone, &r.Email)
		fmt.Println(r)
		roster = append(roster, r)
	}
	return roster
}

// retrieve results of scanRoster
func getRoster(queryString string) []Roster {
	rows, err := db.Query(queryString)
	defer rows.Close()
	results := scanRoster(rows)
	if err != nil {
		log.Fatal(err)
	}
	rows.Close()
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

func getCoworkers(start string, end string) []Roster {
	return getRoster(fmt.Sprintf("SELECT shifts.* AS shift, users.name, users.email, users.phone FROM shifts FULL JOIN users ON shifts.employee_id=users.id WHERE end_time > '%s' AND start_time < '%s'", start, end))
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
