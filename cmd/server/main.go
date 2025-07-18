package main

import (
	"fmt"
	"net/http"

	"github.com/jdpolicano/govault/internal/server"
	"github.com/jdpolicano/govault/internal/server/routes"
)

func main() {
	config := server.DefaultConfig()
	context := server.NewContext(config)
	http.HandleFunc("/register", routes.GetRegisterRoute(context))
	http.HandleFunc("/login", routes.GetLoginRoute(context))
	fmt.Println("listening on port 8080")
	if e := http.ListenAndServe("localhost:8080", nil); e != nil {
		fmt.Println(e)
		return
	}
}
