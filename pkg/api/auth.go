package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const secretKey string = "SecretKey"

func signinHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		log.Println("Method not allowed")
		writeError(w, "Required POST method", http.StatusMethodNotAllowed)
		return
	}

	pass := os.Getenv("TODO_PASSWORD")

	var requ struct {
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&requ)
	if err != nil {
		log.Printf("Invalid JSON: %w\n", err.Error())
		writeError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if requ.Password != pass {
		log.Println("Invalid password")
		writeError(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	claims := jwt.MapClaims{
		"exp": time.Now().Add(8 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Println("Failed generate token")
		writeError(w, "Failed generate token", http.StatusNotImplemented)
		return
	}

	response := map[string]string{
		"token": signedToken,
	}
	log.Println(response)
	w.Header().Set("Content-Type", "application/json")
	writeJson(w, response)

}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// смотрим наличие пароля
		pass := os.Getenv("TODO_PASSWORD")

		if len(pass) > 0 {
			var jwtCookie string // JWT-токен из куки
			// получаем куку
			cookie, err := r.Cookie("token")
			if err == nil {
				jwtCookie = cookie.Value
			}
			var valid bool
			// здесь код для валидации и проверки JWT-токена
			// ...
			token, err := jwt.Parse(jwtCookie, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					log.Println("incorrect signature method")
					return nil, fmt.Errorf("incorrect signature method")
				}
				return []byte(secretKey), nil
			})

			valid = token.Valid

			if err != nil || !valid {
				// возвращаем ошибку авторизации 401

				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}
