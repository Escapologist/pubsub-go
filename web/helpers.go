package web

import (
	"fmt"
	"net/http"
)

func setFlashMessage(w http.ResponseWriter, message string) {
	http.SetCookie(w, &http.Cookie{
		Name:   "flash",
		Value:  message,
		Path:   "/",
		MaxAge: 60,
	})
}

func getFlashMessage(w http.ResponseWriter, r *http.Request) string {
	cookie, err := r.Cookie("flash")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return ""
	}
	flashMessage := cookie.Value
	http.SetCookie(w, &http.Cookie{
		Name:   "flash",
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	})
	return flashMessage
}
