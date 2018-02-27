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

// OpenDatabase : opens connection to database
func OpenDatabase() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s "+
		"dbname=%s sslmode=disable",
		host, port, user, "eggplant", dbname)
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	//defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to the database!!!")
}

// CloseDatabase : closes DB
func CloseDatabase() {
	db.Close()
}

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

// retrieve results of scanUsers (multiple)
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

// retrieve results of scanUser (single)
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

// retrieve results of scanShifts (multiple)
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

// retrieve results of scanShift (single)
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

// consolidate methods for PUT and POST routes
func execute(queryString string, args ...interface{}) {
	fmt.Println(queryString)
	_, err := db.Exec(queryString, args...)
	if err != nil {
		log.Fatal(err)
	}
}

// not needed for user stories but needed to return validation after create and update shift
func getShiftByID(id int64) Shift {
	return getShift(fmt.Sprintf("SELECT * FROM public.shifts WHERE id=%d", id))
}

// // EMPLOYEE user stories:

// As an employee, I want to know when I am working, by being able to see all of the shifts assigned to me:
func getMyShifts(id int64) []Shift {
	return getShifts(fmt.Sprintf("SELECT * FROM public.shifts WHERE employee_id=%d", id))
}

// As an employee, I want to know who I am working with, by being able to see the employees that are working during the same time period as me:
func getCoworkers(start string, end string) []Roster {
	return getRoster(fmt.Sprintf("SELECT shifts.* AS shift, users.name, users.email, users.phone FROM shifts FULL JOIN users ON shifts.employee_id=users.id WHERE end_time > '%s' AND start_time < '%s'", start, end))
}

// As an employee, I want to be able to contact my managers, by seeing manager contact information for my shifts:
func getMyManagers(id int64) []Roster {
	return getRoster(fmt.Sprintf("SELECT shifts.* AS shift, users.name, users.email, users.phone FROM shifts FULL JOIN users ON shifts.manager_id=users.id WHERE employee_id=%d", id))
}

// As an employee, I want to know how much I worked, by being able to get a summary of hours worked for each week:
func getMyHours(id int64, start string, end string) []Shift {
	return getShifts(fmt.Sprintf("SELECT * FROM public.shifts WHERE employee_id=%d AND start_time >='%s' AND end_time < '%s'", id, start, end))
}

// // MANAGER user stories:

// As a manager, I want to see the schedule, by listing shifts within a specific time period:
func getSchedule(start string, end string) []Shift {
	return getShifts(fmt.Sprintf("SELECT * FROM public.shifts WHERE start_time>'%s' AND end_time<'%s'", start, end))
}

// As a manager, I want to schedule my employees, by creating shifts for any employee:
func createShift(shift Shift) Shift {
	queryString := "INSERT INTO public.shifts(manager_id, break, start_time, end_time) VALUES($1, $2, $3, $4);"
	execute(queryString, shift.Manager, shift.Break.Float64, shift.Start, shift.End)
	return getShift(fmt.Sprintf("SELECT * FROM public.shifts WHERE id=MAX"))
}

// As a manager, I want to be able to assign a shift, by changing the employee that will work a shift:
func scheduleEmployee(shift Shift) Shift {
	queryString := "UPDATE shifts SET employee_id=$1, updated_at=now() WHERE id=$2;"
	execute(queryString, shift.Employee.Int64, shift.ID)
	return getShiftByID(shift.ID)
}

// As a manager, I want to be able to change a shift, by updating the time details:
func editShiftTime(shift Shift) Shift {
	queryString := "UPDATE shifts SET start_time=$1, end_time=$2, updated_at=now() WHERE id=$3;"
	execute(queryString, shift.Start, shift.End, shift.ID)
	return getShiftByID(shift.ID)
}

// As a manager, I want to contact an employee, by seeing employee details.
func getAllEmployees() []User {
	return getUsers("SELECT * FROM public.users WHERE role='employee'")
}

// As a manager, I want to contact an employee, by seeing employee details:
func getEmployeeByID(id int64) User {
	return getUser(fmt.Sprintf("SELECT * FROM public.users WHERE role='employee' AND id=%d", id))
}

// // not explicitly needed for user stories
// func getAllUsers() []User {
// 	return getUsers("SELECT * FROM public.users")
// }

// // not explicitly needed for user stories
// func getUserByID(id int64) User {
// 	return getUser(fmt.Sprintf("SELECT * FROM public.users WHERE id=%d", id))
// }

// // not explicitly needed for user stories
// func getAllManagers() []User {
// 	return getUsers("SELECT * FROM public.users WHERE role='manager'")
// }

// // not explicitly needed for user stories
// func getManagerByID(id int64) User {
// 	return getUser(fmt.Sprintf("SELECT * FROM public.users WHERE role='manager' AND id=%d", id))
// }

// // not explicitly needed for user stories
// func getAllShifts() []Shift {
// 	return getShifts("SELECT * FROM public.shifts")
// }
