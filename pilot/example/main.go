package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"time"

	haedal "github.com/voidsatisfaction/haedal-project/pilot"
)

type User struct {
	ID        string
	AddressID string
}

func main() {
	s := haedal.NewServer()

	s.Use(haedal.AuthHandler)

	s.HandleFunc("GET", "/", func(c *haedal.Context) {
		c.RenderTemplate("pilot/example/public/index.html", map[string]interface{}{"time": time.Now()})
	})

	s.HandleFunc("GET", "/about", func(c *haedal.Context) {
		fmt.Fprintln(c.ResponseWriter, "About! ")
	})

	s.HandleFunc("GET", "/users/:id", func(c *haedal.Context) {
		u := User{ID: c.Params["id"].(string)}
		c.RenderXML(u)
	})

	s.HandleFunc("GET", "/users/:user_id/addresses/:address_id", func(c *haedal.Context) {
		u := User{c.Params["user_id"].(string), c.Params["address_id"].(string)}
		c.RenderJSON(u)
	})

	s.HandleFunc("POST", "/users", func(c *haedal.Context) {
		fmt.Fprintln(c.ResponseWriter, "create user ")
	})

	s.HandleFunc("POST", "/users/:user_id/addresses", func(c *haedal.Context) {
		fmt.Fprintf(c.ResponseWriter, "create user%v's address ", c.Params["user_id"])
	})

	s.HandleFunc("GET", "/login", func(c *haedal.Context) {
		c.RenderTemplate("pilot/example/public/login.html", map[string]interface{}{"message": "it needs login"})
	})

	s.HandleFunc("POST", "/login", func(c *haedal.Context) {
		fmt.Printf("%+v", c.Params)
		if CheckLogin(c.Params["username"].(string), c.Params["password"].(string)) {
			http.SetCookie(c.ResponseWriter, &http.Cookie{
				Name:  "X_AUTH",
				Value: Sign(haedal.VerifyMessage),
				Path:  "/",
			})
			c.Redirect("/")
			return
		}
		c.RenderTemplate("/public/login.html", map[string]interface{}{"message": "it is not correct id and password"})
	})
	s.Run(":8080")
}

// TODO: should be encapsulized
func CheckLogin(username, password string) bool {
	const (
		USERNAME = "tester"
		PASSWORD = "12345"
	)

	return username == USERNAME && password == PASSWORD
}

// TODO: should be encapsulized
func Sign(message string) string {
	secretKey := []byte("golang-book-secret-key2")
	if len(secretKey) == 0 {
		return ""
	}
	mac := hmac.New(sha1.New, secretKey)
	io.WriteString(mac, message)
	return hex.EncodeToString(mac.Sum(nil))
}
