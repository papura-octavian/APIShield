package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/fatih/color"
)

var boldGreen = color.New(color.FgGreen, color.Bold)
var boldRed = color.New(color.FgRed, color.Bold)

func clearScreen() {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default: // linux, darwin (mac)
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

type Response struct {
	Method  string `json:"method"`
	Path    string `json:"path"`
	Message string `json:"message"`
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path)

	w.Header().Set("Content-Type", "application/json")

	resp := Response{Method: r.Method, Path: r.URL.Path, Message: "hello from backend!"}

	json.NewEncoder(w).Encode(resp)

}

func main() {
	clearScreen()

	http.HandleFunc("/", helloHandler)

	log.Println("Listening on...")
	boldGreen.Println("-> http://localhost:8000")

	err := http.ListenAndServe(":8000", nil)

	if err != nil {
		log.Fatal(err)
	}
}
