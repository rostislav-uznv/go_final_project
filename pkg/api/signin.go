package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type signinRequest struct {
	Password string `json:"password"`
}

type tokenClaims struct {
	PasswordHash string `json:"password_hash"`
	jwt.RegisteredClaims
}

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	var data signinRequest

	if r.Method != http.MethodPost {
		writeJson(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	password := os.Getenv("TODO_PASSWORD")
	if data.Password != password {
		writeJson(w, http.StatusUnauthorized, map[string]string{"error": "Wrong password"})
		return
	}

	claims := tokenClaims{
		PasswordHash: hashPassword(password),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(password))
	if err != nil {
		writeJson(w, http.StatusInternalServerError, map[string]string{"error": "Failed to sign token"})
		return
	}

	writeJson(w, http.StatusOK, map[string]string{"token": tokenString})

}

func validateToken(tokenString, password string) bool {

	claims := &tokenClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		return []byte(password), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return false
	}

	if !token.Valid {
		return false
	}

	return hashPassword(password) == claims.PasswordHash
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		password := os.Getenv("TODO_PASSWORD")
		if len(password) > 0 {
			var token string
			cookie, err := r.Cookie("token")
			if err == nil {
				token = cookie.Value
			}
			var valid bool
			valid = validateToken(token, password)

			if !valid {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}
