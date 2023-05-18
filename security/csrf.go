package security

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"math"
	"net/http"
)

func GenerateCsrfToken(w http.ResponseWriter, _ *http.Request) (string, error) {

	buff := make([]byte, int(math.Ceil(float64(64)/2)))
	_, err := rand.Read(buff)
	if err != nil {
		log.Println("Error al crear un búfer aleatorio para el token csrf")
		log.Println(err)
		return "", err
	}
	str := hex.EncodeToString(buff)
	token := str[:64]

	cookie := &http.Cookie{
		Name:     "csrf_token",
		Value:    token,
		Path:     "/",
		MaxAge:   1800,
		HttpOnly: true,
		Secure:   true,
	}

	http.SetCookie(w, cookie)

	return token, nil
}

func VerifyCsrfToken(r *http.Request) (bool, error) {
	cookie, err := r.Cookie("csrf_token")
	if err != nil {
		log.Println("Error al obtener la cookie de csrf_token")
		log.Println(err)
		return false, err
	}

	token := r.FormValue("csrf_token")

	if cookie.Value == token {
		return true, nil
	}

	return false, nil
}
