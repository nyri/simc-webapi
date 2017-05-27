package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	port, err := strconv.ParseInt(os.Getenv("PORT"), 10, 32)
	if err != nil {
		port = 8080
	}
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func getOutputFilePath(r *http.Request) string {

	dir, file := filepath.Split("." + r.URL.Path)

	for _, acceptHeader := range r.Header["Accept"] {
		acceptTypes := strings.Split(acceptHeader, ",")
		for _, acceptType := range acceptTypes {
			switch acceptType {
			case "application/json":
				return filepath.Join(dir, file+".json")
			case "text/html":
				return filepath.Join(dir, file+".html")
			case "text/plain":
				return filepath.Join(dir, file+".json")
			}
		}
	}

	return filepath.Join(dir, file+".html")
}

func handler(w http.ResponseWriter, r *http.Request) {

	dir, file := filepath.Split("." + r.URL.Path)

	switch r.Method {
	case http.MethodHead:
	case http.MethodGet:

		path := filepath.Join(dir, file+".simc")

		if _, err := os.Stat(path); os.IsNotExist(err) {
			http.Error(w, "404 page not found", http.StatusNotFound)
			return
		}

		filePath := getOutputFilePath(r)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			w.WriteHeader(http.StatusAccepted)
			return
		}

		http.ServeFile(w, r, filePath)

	case http.MethodPut:

		path := filepath.Join(dir, file+".simc")

		if _, err := os.Stat(path); err == nil {
			http.Error(w, "409 conflict", http.StatusConflict)
			return
		}

		f, err := os.Create(path)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer f.Close()

		if _, err := io.Copy(f, r.Body); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		f.Close()

		cmd := exec.Command("simc", fmt.Sprintf("%s.simc", file), fmt.Sprintf("json=%s.json", file), fmt.Sprintf("html=%s.html", file), fmt.Sprintf("output=%s.txt", file))
		cmd.Dir = dir

		if err := cmd.Start(); err != nil {
			os.Remove(path)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}
