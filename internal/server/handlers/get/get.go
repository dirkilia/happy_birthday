package get

import (
	"encoding/json"
	"happy_birthday/internal/database"
	jwtlib "happy_birthday/internal/lib/jwt"
	"log/slog"
	"net/http"
)

// Get all employees
func GetAllEmployeesHandler(log *slog.Logger, db database.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := db.GetAllEmployees()
		if err != nil {
			log.Error("Can't get employees list. Err: %v", err)
		}

		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Error("error handling JSON marshal. Err: %v", err)
		}

		_, _ = w.Write(jsonResp)
	}
}

// Get all today employees birthdays
func GetAllEmployeesTodayBirthdaysHandler(log *slog.Logger, db database.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := db.GetAllEmployeesTodayBirthdays()
		if err != nil {
			log.Error("Can't get employees list. Err: %v", err)
		}

		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Error("error handling JSON marshal. Err: %v", err)
		}

		_, _ = w.Write(jsonResp)
	}
}

// Get current user info
func GetCurrentUserInfo(log *slog.Logger, db database.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tknStr := c.Value

		claims := jwtlib.GetClaimsFromToken(tknStr)
		for key, val := range claims {
			if key == "login" {
				employee_db, err := db.GetEmployeeByLogin(val.(string))
				if err != nil {
					log.Error("can't get user info. Err: %v", err)
				}

				jsonResp, err := json.Marshal(employee_db)
				if err != nil {
					log.Error("error handling JSON marshal. Err: %v", err)
				}

				_, _ = w.Write(jsonResp)
			}
		}

	}
}
