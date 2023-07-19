package main

import (
	"net/http"

	"github.com/taufiqoo/technical-test-berat/controllers/usercontroller"
)

func main() {

	http.HandleFunc("/", usercontroller.Index)
	http.HandleFunc("/user", usercontroller.Index)
	http.HandleFunc("/user/index", usercontroller.Index)
	http.HandleFunc("/user/add", usercontroller.Add)
	http.HandleFunc("/user/edit", usercontroller.Edit)
	http.HandleFunc("/user/delete", usercontroller.Delete)
	http.HandleFunc("/user/show", usercontroller.Show)

	http.ListenAndServe(":3030", nil)
}
