package database

import (
	"database/sql"
	"fmt"
	"time"
)

type Repository struct {
	db *sql.DB
}

type Employee struct {
	Id      int
	Name    string
	Workday time.Time
}

func NewRepository(conn *DBConn) *Repository {
	return &Repository{db: conn.db}
}

func (r *Repository) GettingGroupsId() (int64, error) {
	var groupids int64
	query := "SELECT groupchat_id FROM groupid"

	err := r.db.QueryRow(query).Scan(&groupids)
	if err != nil {
		fmt.Printf("error during retrieving groupsid from database %v", err)
	}
	return groupids, nil
}

func (r *Repository) InsertGroupChatId(groupId int64) error {

	query := "INSERT INTO groupid(groupchat_id) VALUES ($1)"
	_, err := r.db.Exec(query, groupId)
	if err != nil {
		return fmt.Errorf("could not insert into groupchatId: %v", err)
	}
	return nil
}

func (r *Repository) GetDutyEmployeeData() (Employee, error) {
	query := `
    WITH last_employee AS (
        SELECT workday, ename
FROM employees
ORDER BY
    CASE
        WHEN workday = CURRENT_DATE THEN 1
        ELSE 2
    END,
    workday
LIMIT 1
    )
    UPDATE employees
    SET workday = CURRENT_DATE
    WHERE ename = (SELECT ename FROM last_employee)
    RETURNING id, ename, workday;
    `

	row := r.db.QueryRow(query)

	var employee Employee
	err := row.Scan(&employee.Id, &employee.Name, &employee.Workday)
	if err != nil {
		return employee, fmt.Errorf("error fetching duty employee: %v", err)
	}

	return employee, nil
}

func (r *Repository) GetNextDutyEmployee() (string, error) {
	query := `
   WITH today_employee AS (
    SELECT id, ename, workday
    FROM employees
    WHERE workday = CURRENT_DATE
    LIMIT 1
),
next_employee AS (
    SELECT id, ename, workday
    FROM employees
    WHERE id = (SELECT id FROM today_employee) + 1
),
update_today_employee AS (
    UPDATE employees
    SET workday = (SELECT workday FROM next_employee)
    WHERE id = (SELECT id FROM today_employee)
    RETURNING id, ename, workday
),
update_next_employee AS (
    UPDATE employees
    SET workday = CURRENT_DATE
    WHERE id = (SELECT id FROM next_employee)
    RETURNING id, ename, workday
)
SELECT ename FROM update_next_employee;
    `
	row := r.db.QueryRow(query)

	var employee Employee
	err := row.Scan(&employee.Name)
	if err != nil {
		return employee.Name, fmt.Errorf("error fetching next duty employee: %v", err)
	}

	return employee.Name, nil
}

func (r *Repository) GetPreviousDutyEmployee() (string, error) {
	query := `
	WITH today_employee AS (
    SELECT id, ename, workday
    FROM employees
    WHERE workday = CURRENT_DATE
    LIMIT 1
),
previous_employee AS (
    SELECT id, ename, workday
    FROM employees
    WHERE id = (SELECT id FROM today_employee) - 1
),
update_today_employee AS (
    UPDATE employees
    SET workday = (SELECT workday FROM previous_employee)
    WHERE id = (SELECT id FROM today_employee)
    RETURNING id, ename, workday
),
update_previous_employee AS (
    UPDATE employees
    SET workday = CURRENT_DATE
    WHERE id = (SELECT id FROM previous_employee)
    RETURNING id, ename, workday
)
SELECT ename FROM update_previous_employee;
	`
	row := r.db.QueryRow(query)
	var employee Employee
	err := row.Scan(&employee.Name)
	if err != nil {
		return employee.Name, fmt.Errorf("error fetching previous duty employee: %v", err)
	}
	return employee.Name, nil
}
