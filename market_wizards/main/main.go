package main

import (
	"log"
	"net/http"
	"path/filepath"
	"runtime"

	"github.com/KushamiNeko/market_wizards/handler"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	mux := http.NewServeMux()

	mux.Handle("/view/", &handler.ViewHandler{})

	mux.Handle("/plot/", &handler.PlotHandler{})

	mux.HandleFunc("/resources/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(".", r.RequestURI))
	})

	log.Fatal(http.ListenAndServe(":8080", mux))
}
