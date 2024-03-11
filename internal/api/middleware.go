package api

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shahin-bayat/go-api/internal/model"
	"github.com/shahin-bayat/go-api/internal/store"
	"github.com/shahin-bayat/go-api/internal/util"
)

func withJwtAuth(handleFunc http.HandlerFunc, s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("x-jwt-token")
		token, err := validateJwt(tokenString)
		if err != nil || !token.Valid {
			util.WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "permission denied"})
			return
		}

		userId, err := util.GetId(r)
		if err != nil {
			util.WriteJSON(w, http.StatusUnauthorized, ApiError{Error: err.Error()})
			return
		}

		account, err := s.GetAccountById(userId)
		if err != nil {
			// improvement: create a custom error func
			util.WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "permission denied"})
			// you can still log the error to some logging service (elastic) to knw=ow what went wrong
			return
		}

		// cast the claims to jwt.MapClaims
		claims := token.Claims.(jwt.MapClaims)

		if account.IBAN != claims["iban"] {
			util.WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "permission denied"})
			return
		}

		handleFunc(w, r)
	}
}

func validateJwt(tokenString string) (*jwt.Token, error) {
	// https://pkg.go.dev/github.com/golang-jwt/jwt/v5#section-readme
	// you should use export JWT_SECRET=your_secret in your terminal
	secret := os.Getenv("JWT_SECRET")
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

}

func createJwt(account *model.Account) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt": jwt.NewNumericDate(time.Unix(1516239022, 0)),
		"iban":      account.IBAN,
	}
	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))

}
