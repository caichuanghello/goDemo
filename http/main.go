package main

import (
	"net/http"
)


type user struct {
	Name string
	Age int
}
func main(){


	http.HandleFunc("/",hello)
	http.ListenAndServe(":8080",nil)
}

func hello(w http.ResponseWriter,r *http.Request){
	w.Write([]byte("hello world!"))
}
