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

// CloseDatabase : closes database
func CloseDatabase() {
	db.Close()
}

// execute a query/statement that selects no data
func executeStatement(queryString string, args ...interface{}) {
	fmt.Println(queryString)
	_, err := db.Exec(queryString, args...)
	if err != nil {
		log.Fatal(err)
	}
}

// query single row in DB
func selectOne(queryString string, queryArgs ...interface{}) *sql.Row {
	row := db.QueryRow(queryString, queryArgs...)
	return row
}

// query multiple rows in DB
func selectMany(queryString string, queryArgs ...interface{}) *sql.Rows {
	rows, err := db.Query(queryString, queryArgs...)
	if err != nil {
		log.Fatal(err)
	}
	return rows
}

// map multiple user rows into slice of Users
func multipleRowsToUsers(rows *sql.Rows) []User {
	var users []User
	var u User
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&u.ID, &u.Name, &u.Role, &u.Email, &u.Phone, &u.Created, &u.Updated)
		fmt.Println(u)
		users = append(users, u)
	}
	rows.Close()
	return users
}

// map single user row into User struct
func singleRowToUser(row *sql.Row) User {
	var u User
	row.Scan(&u.ID, &u.Name, &u.Role, &u.Email, &u.Phone, &u.Created, &u.Updated)
	fmt.Println(u)
	return u
}

// map multiple shift rows into slice of Shifts
func multipleRowsToShifts(rows *sql.Rows) []Shift {
	var shifts []Shift
	var s Shift
	for rows.Next() {
		rows.Scan(&s.ID, &s.Manager, &s.Employee, &s.Break, &s.Start, &s.End, &s.Created, &s.Updated)
		fmt.Println(s)
		shifts = append(shifts, s)
	}
	rows.Close()
	return shifts
}

// map single shift row into Shift struct
func singleRowToShift(row *sql.Row) Shift {
	var s Shift
	row.Scan(&s.ID, &s.Manager, &s.Employee, &s.Break, &s.Start, &s.End, &s.Created, &s.Updated)
	fmt.Println(s)
	return s
}

// map multiple rows into slice of Roster
func multipleRowsToRoster(rows *sql.Rows) []Roster {
	var roster []Roster
	var r Roster
	for rows.Next() {
		rows.Scan(&r.ID, &r.Manager, &r.Employee, &r.Break, &r.Start, &r.End, &r.Created, &r.Updated, &r.Name, &r.Phone, &r.Email)
		fmt.Println(r)
		roster = append(roster, r)
	}
	rows.Close()
	return roster
}

// // EMPLOYEE user stories:

// As an employee, I want to know when I am working, by being able to see all of the shifts assigned to me:
func getShiftsByEmployee(id int64) []Shift {
	shiftRows := selectMany("SELECT * FROM shifts WHERE employee_id=$1", id)
	return multipleRowsToShifts(shiftRows)
}

// As an employee, I want to know who I am working with, by being able to see the employees that are working during the same time period as me:
func getEmployeeRostersByDateRange(start string, end string) []Roster {
	rosterRows := selectMany("SELECT shifts.* AS shift, users.name, users.email, users.phone FROM shifts FULL JOIN users ON shifts.employee_id=users.id WHERE end_time > $1 AND start_time < $2", start, end)
	return multipleRowsToRoster(rosterRows)
}

// As an employee, I want to be able to contact my managers, by seeing manager contact information for my shifts:
func getManagerRostersByDateRange(id int64) []Roster {
	rosterRows := selectMany("SELECT shifts.* AS shift, users.name, users.email, users.phone FROM shifts FULL JOIN users ON shifts.manager_id=users.id WHERE employee_id=$1", id)
	return multipleRowsToRoster(rosterRows)
}

// As an employee, I want to know how much I worked, by being able to get a summary of hours worked for each week:
func getShiftsByEmployeeInDateRange(id int64, start string, end string) []Shift {
	shiftRows := selectMany("SELECT * FROM shifts WHERE employee_id=$1 AND start_time > $2 AND end_time < $3", id, start, end)
	return multipleRowsToShifts(shiftRows)
}

// // MANAGER user stories:

// As a manager, I want to see the schedule, by listing shifts within a specific time period:
func getShiftsByDateRange(start string, end string) []Shift {
	shiftRows := selectMany("SELECT * FROM shifts WHERE start_time > $1 AND end_time < $2", start, end)
	return multipleRowsToShifts(shiftRows)
}

// As a manager, I want to schedule my employees, by creating shifts for any employee:
func createShift(shift Shift) Shift {
	queryString := "INSERT INTO shifts(manager_id, break, start_time, end_time) VALUES($1, $2, $3, $4);"
	executeStatement(queryString, shift.Manager, shift.Break.Float64, shift.Start, shift.End)
	return getNewestShift()
}

// As a manager, I want to be able to assign a shift, by changing the employee that will work a shift:
func scheduleEmployee(shift Shift) Shift {
	queryString := "UPDATE shifts SET employee_id=$1, updated_at=now() WHERE id=$2;"
	executeStatement(queryString, shift.Employee.Int64, shift.ID)
	return getShiftByID(shift.ID)
}

// As a manager, I want to be able to change a shift, by updating the time details:
func editShiftTime(shift Shift) Shift {
	queryString := "UPDATE shifts SET start_time=$1, end_time=$2, updated_at=now() WHERE id=$3;"
	executeStatement(queryString, shift.Start, shift.End, shift.ID)
	return getShiftByID(shift.ID)
}

// As a manager, I want to contact an employee, by seeing employee details.
func getAllEmployees() []User {
	employeeRows := selectMany("SELECT * FROM users WHERE role='employee'")
	return multipleRowsToUsers(employeeRows)
}

// As a manager, I want to contact an employee, by seeing employee details:
func getEmployeeByID(id int64) User {
	employeeRow := selectOne("SELECT * FROM users WHERE role='employee' AND id=$1", id)
	return singleRowToUser(employeeRow)
}

// returns most recently updated shift after PUT:
func getShiftByID(id int64) Shift {
	shiftRow := selectOne("SELECT * FROM shifts WHERE id=$1", id)
	return singleRowToShift(shiftRow)
}

// returns most recently created shift after POST:
func getNewestShift() Shift {
	shiftRow := selectOne("SELECT * FROM shifts WHERE id = (SELECT MAX(id) FROM shifts")
	return singleRowToShift(shiftRow)
}
