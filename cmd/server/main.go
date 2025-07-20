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
	refs := server.NewServerRefs(config)
	http.HandleFunc("/register", register.Handler(refs))
	http.HandleFunc("/login", login.Handler(refs))
	http.HandleFunc("/get", get.Handler(refs))
	http.HandleFunc("/set", set.Handler(refs))
	fmt.Println("listening on port 8080")
	if e := http.ListenAndServe("localhost:8080", nil); e != nil {
		fmt.Println(e)
		return
	}
}
