package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

type Response struct {
	Status  string
	Message string
}

func doResponse(w http.ResponseWriter, status int, statusMessage string, message string) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	response := Response{statusMessage, message}
	jsonResponse, err := Marshal(response)
	if err != nil {
		log.Info().
			Msg("Error marshalling response: " + string(err.Error()))
		return
	}
	io.WriteString(w, string(jsonResponse))
	log.Info().
		Msg("Response:" + string(jsonResponse))
}

func getHealth(w http.ResponseWriter, r *http.Request) {
	log.Info().
		Msg("got /health request\n")
	doResponse(w, http.StatusOK, "ok", "Healthy")
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func postSubmit(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		doResponse(w, http.StatusOK, "ok", "fileName")
	} else {
		doResponse(w, http.StatusMethodNotAllowed, "error", "Method not allowed")
		log.Info().
			Msg("Method not allowed")
	}
}

func main() {
	log.Info().
		Msg("Starting server...\n")

	viewPath := "static/"
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info().Msg("Requested path: " + r.URL.Path)

		// hack, if requested url is / then point towards /index
		if r.URL.Path == "/" {
			r.URL.Path = "/submit"
		}

		requestedPath := strings.TrimLeft(filepath.Clean(r.URL.Path), "/")
		// If an extension is provided, just use the path
		if strings.Contains(r.URL.Path, ".") {
			http.ServeFile(w, r, requestedPath)
		} else {
			// Otherwise append .html to get the correct file
			filename := fmt.Sprintf("%s/%s.html", viewPath, requestedPath)
			http.ServeFile(w, r, filename)
		}
	}))
	http.HandleFunc("/health", getHealth)
	http.HandleFunc("/api/submit", postSubmit)

	err := http.ListenAndServe(":3333", nil)

	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
		return
	}

	log.Info().
		Msg("Server started on port 3333\n")
}
