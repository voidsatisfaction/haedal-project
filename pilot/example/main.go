package main

import (
	"fmt"

	haedal "github.com/voidsatisfaction/haedal-project/pilot"
)

func main() {
	s := haedal.NewServer()

	s.HandleFunc("GET", "/", func(c *haedal.Context) {
		fmt.Fprintln(c.ResponseWriter, "welcome! ")
	})

	s.HandleFunc("GET", "/about", func(c *haedal.Context) {
		fmt.Fprintln(c.ResponseWriter, "About! ")
	})

	s.HandleFunc("GET", "/users/:id", func(c *haedal.Context) {
		if c.Params["id"] == "0" {
			panic("id is zero")
		}
		fmt.Fprintf(c.ResponseWriter, "retrieve user %v\n", c.Params["id"])
	})

	s.HandleFunc("GET", "/users/:user_id/addresses/:address_id", func(c *haedal.Context) {
		fmt.Fprintf(c.ResponseWriter, "retrieve user%v's address %v\n", c.Params["user_id"], c.Params["address_id"])
	})

	s.HandleFunc("POST", "/users", func(c *haedal.Context) {
		fmt.Fprintln(c.ResponseWriter, "create user ")
	})

	s.HandleFunc("POST", "/users/:user_id/addresses", func(c *haedal.Context) {
		fmt.Fprintf(c.ResponseWriter, "create user%v's address ", c.Params["user_id"])
	})

	s.Run(":8080")
}
