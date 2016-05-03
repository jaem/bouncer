package main

import (
	"fmt"
	"net/http"

	"github.com/jaem/bounce"
	"github.com/jaem/bounce/providers/local"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/jaem/nimble"
	"github.com/jaem/bounce/storage/jwt"

)

func main() {
	//trial.TestRouters()
	theMain()
}

func theMain() {

	router := mux_router()

	nim := nimble.Default()
	nim.UseFunc(middlewareA)
	nim.UseFunc(middlewareB)

	// router goes last
	nim.Use(router)
	nim.Run(":8000")
}

func mux_router() *mux.Router {
	router := mux.NewRouter()
	router.StrictSlash(true)
	router.HandleFunc("/hello", helloHandlerFunc).Methods("GET")

	bou := bounce.New(jwt.NewStorage())
	bou.Register("local", local.NewProvider())
	bou.Register("local2", local.NewProvider())

	authRoutes := mux.NewRouter()
	authRoutes.HandleFunc("/auth/{userid}/login", authHandlerFunc)
	router.PathPrefix("/auth").Handler(nimble.New().
		UseHandlerFunc(bou.Authenticate("local")).
		UseHandlerFunc(bou.Hoho).
		Use(authRoutes),
	)

	return router
}

func helloHandlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello Hax!")
	if value, ok := context.GetOk(r, "value"); ok {
		fmt.Println("from helloHandlerFunc, value is " + value.(string))
	}
	if value, ok := context.GetOk(r, "valueA"); ok {
		fmt.Println("from helloHandlerFunc, valueA is " + value.(string))
	}
	if value, ok := context.GetOk(r, "valueB"); ok {
		fmt.Println("from helloHandlerFunc, valueB is " + value.(string))
	}
}

func middlewareA(w http.ResponseWriter, r *http.Request) {
	if value, ok := context.GetOk(r, "value"); ok {
		fmt.Println("from middlewareA, value is " + value.(string))
	} else {
		fmt.Println("from middlewareA, value is nil")
	}
	context.Set(r, "value", "A")
	context.Set(r, "valueA", "A")
}

func middlewareB(w http.ResponseWriter, r *http.Request) {
	if value, ok := context.GetOk(r, "value"); ok {
		fmt.Println("from middlewareB, value is " + value.(string))
	}
	context.Set(r, "value", "B")
	context.Set(r, "valueB", "B")
}

func authHandlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Authed Hax!")
	if value, ok := context.GetOk(r, "value"); ok {
		fmt.Println("from authHandlerFunc, value is " + value.(string))
	}
	// using mux
	fmt.Println("from authHandlerFunc, the userid is " + mux.Vars(r)["userid"])
}

func middlewareAuth(w http.ResponseWriter, r *http.Request) {
	if value, ok := context.GetOk(r, "value"); ok {
		fmt.Println("from middlewareAuth, value is " + value.(string))
	}
	context.Set(r, "value", "AUTH")
}