package main

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

var tpl = template.Must(template.ParseFiles("templates/index.html"))

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/save", saveHandler)

	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, nil)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	data := r.FormValue("signature")
	if idx := strings.Index(data, ";base64,"); idx != -1 {
		raw := data[idx+8:]
		imgData, err := base64.StdEncoding.DecodeString(raw)
		if err != nil {
			http.Error(w, "Invalid data", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Content-Disposition", "attachment; filename=signature.png")
		w.Write(imgData)
		return
	}
	http.Error(w, "No signature data", http.StatusBadRequest)
}
