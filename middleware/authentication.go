package middleware

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

var SECRET = []byte("super-secret-key")
var api_key = "test"

func CreaateJWT() (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["exp"] = time.Now().Add(time.Hour).Unix()

	tokenStr, err := token.SignedString(SECRET)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	return tokenStr, nil
}

func ValidateJWT(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(t *jwt.Token) (interface{}, error) {
				_, ok := t.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("not authorized"))
				}
				return SECRET, nil
			})

			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("not authorized:" + err.Error()))
			}
			if token.Valid {
				next(w, r)
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("not authorized"))
		}
	})
}

func GetJWT(w http.ResponseWriter, r *http.Request) {
	if r.Header["Access"] != nil {
		if r.Header["Access"][0] == api_key {
			token, err := CreaateJWT()
			if err != nil {
				return
			}
			fmt.Fprintf(w, token)
		}
	}
}
