package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"time"

	"github.com/fatih/color"
)

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
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
	http.Handle("/", loggingMiddleware(proxy))

	log.Println("proxy listening on...")
	boldGreen.Print("http://localhost:8080")
	boldWhite.Print(" -> ")
	boldCyan.Println("http://localhost:8000")

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}
