package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/fatih/color"
)

var boldGreen = color.New(color.FgGreen, color.Bold)
var boldRed = color.New(color.FgRed, color.Bold)

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path)
	fmt.Fprintln(w, "hello")
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
