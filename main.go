package main

import (
	"gfh.com/web/routes"
	"gfh.com/web/utils"
	"log"
	"net/http"
)

func main() {
	utils.LoadTemplates("templates/*.html")
	r := routes.NewRouter()
	http.Handle("/", r)
	log.Println("running on http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}
