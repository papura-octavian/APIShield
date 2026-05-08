package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"time"

	"github.com/fatih/color"
)

type CheckRequest struct {
	Method  string            `json:"method"`
	Path    string            `json:"path"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

type CheckResponse struct {
	Action string `json:"action"`
	Reason string `json:"reason"`
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func securityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkReq := CheckRequest{
			Method:  r.Method,
			Path:    r.URL.Path,
			Headers: map[string]string{},
			Body:    "",
		}

		data, err := json.Marshal(checkReq)
		if err != nil {
			log.Printf("eroare la marshal: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		resp, err := http.Post(
			"http://localhost:9000/check",
			"application/json",
			bytes.NewReader(data),
		)
		if err != nil {
			log.Printf("eroare la marshal: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		var checkResp CheckResponse
		if err := json.NewDecoder(resp.Body).Decode(&checkResp); err != nil {
			log.Printf("eroare la marshal: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		if checkResp.Action == "block" {
			http.Error(w, "blocked by APIShield: "+checkResp.Reason, http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

var boldCyan = color.New(color.FgCyan, color.Bold)
var boldGreen = color.New(color.FgGreen, color.Bold)
var boldWhite = color.New(color.FgWhite, color.Bold)
var boldYellow = color.New(color.FgYellow, color.Bold)
var boldMagenta = color.New(color.FgMagenta, color.Bold)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		boldWhite.Printf("%s ", time.Now().Format("15:04:05"))
		boldYellow.Print("[PROXY] ")
		boldMagenta.Printf("%s ", r.Method)
		boldWhite.Printf("%s\n", r.URL.Path)

		next.ServeHTTP(w, r)
	})
}

func main() {
	clearScreen()

	target, err := url.Parse("http://localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	http.Handle("/", loggingMiddleware(securityMiddleware(proxy)))

	log.Println("proxy listening on...")
	boldGreen.Print("http://localhost:8080")
	boldWhite.Print(" -> ")
	boldCyan.Println("http://localhost:8000")

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}
