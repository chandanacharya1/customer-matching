package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"html/template"
	"net/http"
	"time"
)

var jwtKey = []byte("secret_key")

var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func Index(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./www/pages/login.html")
	err := t.Execute(w, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	var credentials Credentials
	/*	err := json.NewDecoder(r.Body).Decode(&credentials)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	*/
	credentials.Username = r.FormValue("username")
	credentials.Password = r.FormValue("password")
	expectedPassword, ok := users[credentials.Username]

	if !ok || expectedPassword != credentials.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(time.Minute * 5)

	claims := &Claims{
		Username: credentials.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w,
		&http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})
	http.Redirect(w, r, "/homepage", http.StatusSeeOther)
}

func Validate(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tokenStr := cookie.Value

		claims := &Claims{}

		tkn, err := jwt.ParseWithClaims(tokenStr, claims,
			func(t *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next(w, r)
	})
}
