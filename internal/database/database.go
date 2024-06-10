package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
)

type Employee struct {
	Id                  int64  `json:"id"`
	Login               string `json:"login"`
	Password            string `json:"password"`
	FirstName           string `json:"first_name"`
	Surname             string `json:"surname"`
	Patronymic          string `json:"patronymic"`
	Birthday            string `json:"birthday"`
	EnableNotifications int    `json:"enable_notifications"`
	NotifyOf            string `json:"notify_of"`
}

type NotifyOf struct {
	Login    string `json:"login"`
	NotifyOf string `json:"notify_of"`
}

type Service interface {
	GetAllEmployees() ([]Employee, error)
	GetAllEmployeesTodayBirthdays() ([]Employee, error)
	SetNotifyOf(string, string) error
	RegisterEmployee(Employee) (int64, error)
	GetEmployeeByLogin(string) (Employee, error)
}

type service struct {
	db *sql.DB
}

var (
	dburl      = os.Getenv("DB_URL")
	dbInstance *service
)

func New() Service {
	if dbInstance != nil {
		return dbInstance
	}

	db, err := sql.Open("sqlite3", dburl)
	if err != nil {
		log.Fatal(err)
	}

	dbInstance = &service{
		db: db,
	}
	return dbInstance
}

func (s *service) GetAllEmployees() ([]Employee, error) {
	stmt, err := s.db.Prepare("SELECT id, login, password, first_name, surname, patronymic, birthday FROM employees")
	if err != nil {
		return nil, fmt.Errorf("prepare statement: %v", err)
	}

	resultEmployees := []Employee{}
	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("rows not found: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		employee := Employee{}
		err = rows.Scan(&employee.Id, &employee.Login, &employee.Password, &employee.FirstName, &employee.Surname, &employee.Patronymic, &employee.Birthday)
		if err != nil {
			return nil, fmt.Errorf("failed to get rows: %v", err)
		}
		resultEmployees = append(resultEmployees, employee)
	}

	return resultEmployees, nil
}

func (s *service) GetAllEmployeesTodayBirthdays() ([]Employee, error) {
	currentTime := time.Now()
	_, m, d := currentTime.Date()

	stmt, err := s.db.Prepare(fmt.Sprintf("SELECT id, login, password, first_name, surname, patronymic, birthday FROM employees WHERE birthday LIKE '%%-%d-%d'", m, d))
	if err != nil {
		return nil, fmt.Errorf("prepare statement: %v", err)
	}

	resultEmployees := []Employee{}
	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("rows not found: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		employee := Employee{}
		err = rows.Scan(&employee.Id, &employee.Login, &employee.Password, &employee.FirstName, &employee.Surname, &employee.Patronymic, &employee.Birthday)
		if err != nil {
			return nil, fmt.Errorf("failed to get rows: %v", err)
		}
		resultEmployees = append(resultEmployees, employee)
	}

	return resultEmployees, nil
}

func (s *service) SetNotifyOf(user_to_notify, notify_of_users string) error {
	stmt, err := s.db.Prepare("UPDATE employees SET notify_of = ? WHERE login = ?")
	if err != nil {
		return fmt.Errorf("prepare statement: %v", err)
	}

	_, err = stmt.Exec(notify_of_users, user_to_notify)
	if err != nil {
		fmt.Printf("execute statement: %v", err)
		return fmt.Errorf("execute statement: %v", err)
	}

	return nil
}

func (s *service) RegisterEmployee(employee Employee) (int64, error) {
	stmt, err := s.db.Prepare("INSERT INTO employees(login, password, first_name, surname, patronymic, birthday, enable_notifications) values(?,?,?,?,?,?,?)")
	if err != nil {
		fmt.Print(err)
		return 0, fmt.Errorf("prepare statement: %v", err)
	}

	res, err := stmt.Exec(employee.Login, employee.Password, employee.FirstName, employee.Surname, employee.Patronymic, employee.Birthday, employee.EnableNotifications)
	if err != nil {
		fmt.Print(err)
		return 0, fmt.Errorf("execute statement: %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		fmt.Print(err)
		return 0, fmt.Errorf("failed to get last insert id: %v", err)
	}

	return id, nil
}

func (s *service) GetEmployeeByLogin(login string) (Employee, error) {
	stmt, err := s.db.Prepare("SELECT id, login, password, first_name, surname, patronymic, birthday, enable_notifications, notify_of FROM employees WHERE login = ?")
	if err != nil {
		return Employee{}, fmt.Errorf("prepare statement: %v", err)
	}

	var employee Employee
	row := stmt.QueryRow(login)
	err = row.Scan(&employee.Id, &employee.Login, &employee.Password, &employee.FirstName, &employee.Surname, &employee.Patronymic, &employee.Birthday, &employee.EnableNotifications, &employee.NotifyOf)
	if err != nil {
		return Employee{}, fmt.Errorf("no such user: %v", err)
	}

	return employee, nil
}
