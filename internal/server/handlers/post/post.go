package post

import (
	"encoding/json"
	"happy_birthday/internal/database"
	jwtlib "happy_birthday/internal/lib/jwt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

// Register user with `login`, `password`, `first_name`, `surname`, `patronymic`, `birthday`, `enable_notifications` fields
func RegisterUserHandler(log *slog.Logger, db database.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var employee database.Employee

		err := render.DecodeJSON(r.Body, &employee)
		if err != nil {
			log.Error("failed to decode request. Err: %v", err)
			render.JSON(w, r, Response{
				Status: "Error",
				Error:  "failed to decode request",
			})

			return
		}

		passHash, err := bcrypt.GenerateFromPassword([]byte(employee.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Error("failed to generate password. Err: %v", err)
			render.JSON(w, r, Response{
				Status: "Error",
				Error:  "failed to generate password",
			})
			return
		}
		employee.Password = string(passHash)

		_, err = db.RegisterEmployee(employee)
		if err != nil {
			log.Error("failed to register employee. Err: %v", err)
			render.JSON(w, r, Response{
				Status: "Error",
				Error:  "failed to register employee",
			})

			return
		}

		render.JSON(w, r, Response{
			Status: "OK",
		})
	}
}

// Auth user with given login and password
func AuthUserHandler(log *slog.Logger, db database.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var employee database.Employee

		err := render.DecodeJSON(r.Body, &employee)
		if err != nil {
			log.Error("failed to decode request. Err: %v", err)
			render.JSON(w, r, Response{
				Status: "Error",
				Error:  "failed to decode request",
			})

			return
		}

		employee_db, err := db.GetEmployeeByLogin(employee.Login)
		if err != nil {
			log.Error("failed to get employee pass. Err: %v", err)
			render.JSON(w, r, Response{
				Status: "Error",
				Error:  "failed to get employee pass",
			})

			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(employee_db.Password), []byte(employee.Password)); err != nil {
			log.Error("failed to compare password hashes. Err: %v", err)
			render.JSON(w, r, Response{
				Status: "Error",
				Error:  "failed to compare password hashes",
			})

			return
		}

		token, err := jwtlib.NewToken(employee, 24*time.Hour)
		if err != nil {
			log.Error("failed to generate token. Err: %v", err)
			render.JSON(w, r, Response{
				Status: "Error",
				Error:  "failed to generate token",
			})

			return
		}

		cookie := http.Cookie{
			Name:  "token",
			Path:  "/",
			Value: token,
		}

		http.SetCookie(w, &cookie)

		tdaybday, err := db.GetAllEmployeesTodayBirthdays()
		if err != nil {
			log.Error("failed to get today birthdays employees. Err: %v", err)
			render.JSON(w, r, Response{
				Status: "Error",
				Error:  "failed to get today birthdays employees",
			})

			return
		}

		if employee_db.EnableNotifications == 0 {
			_, _ = w.Write([]byte(`{"message":"Logged in successfully"}`))
		} else {
			notify_of_employees := strings.Split(employee_db.NotifyOf, ";")
			emp_to_send := findIntersection(notify_of_employees, tdaybday)
			if len(emp_to_send) > 0 {
				jsonResp, err := json.Marshal(emp_to_send)
				if err != nil {
					log.Error("error handling JSON marshal. Err: %v", err)
				}
				_, _ = w.Write([]byte(jsonResp))
			}
		}

	}
}

// Set notify_of field to current user in database
func SetNotifyOfHandler(log *slog.Logger, db database.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var notify_of_employees database.NotifyOf

		err := render.DecodeJSON(r.Body, &notify_of_employees)
		if err != nil {
			log.Error("failed to decode request. Err: %v", err)
			render.JSON(w, r, Response{
				Status: "Error",
				Error:  "failed to decode request",
			})

			return
		}

		err = db.SetNotifyOf(notify_of_employees.Login, notify_of_employees.NotifyOf)
		if err != nil {
			log.Error("failed to set notifies for employee. Err: %v", err)
			render.JSON(w, r, Response{
				Status: "Error",
				Error:  "failed to set notifies for employee",
			})

			return
		}

		render.JSON(w, r, Response{
			Status: "OK",
		})
	}

}

// Logout current user
func Logout(log *slog.Logger, db database.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Expires: time.Now(),
		})
	}
}

// Util func to find which employees birthdays user chose to be notified of
func findIntersection(notify_of_employees []string, tdaybday []database.Employee) []database.Employee {
	intersection := make([]database.Employee, 0)

	set := make(map[string]bool)

	for _, name := range notify_of_employees {
		set[name] = true
	}

	for _, emp := range tdaybday {
		login := emp.Login
		if set[login] {
			intersection = append(intersection, emp)
		}
	}

	return intersection
}
