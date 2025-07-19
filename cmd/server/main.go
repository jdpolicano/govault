package main

import (
	"fmt"
	"net/http"

	"github.com/jdpolicano/govault/internal/server"
	"github.com/jdpolicano/govault/internal/server/routes/get"
	"github.com/jdpolicano/govault/internal/server/routes/login"
	"github.com/jdpolicano/govault/internal/server/routes/register"
	"github.com/jdpolicano/govault/internal/server/routes/set"
)

func main() {
	config := server.DefaultConfig()
	context := server.NewContext(config)
	http.HandleFunc("/register", register.Handler(context))
	http.HandleFunc("/login", login.Handler(context))
	http.HandleFunc("/get", get.Handler(context))
	http.HandleFunc("/set", set.Handler(context))
	fmt.Println("listening on port 8080")
	if e := http.ListenAndServe("localhost:8080", nil); e != nil {
		fmt.Println(e)
		return
	}
}
