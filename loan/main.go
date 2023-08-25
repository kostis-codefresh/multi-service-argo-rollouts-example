package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"text/template"
)

type LoanApplication struct {
	AppVersion     string
	BackendVersion string
	BackendHost    string
	BackendPort    string
	LoanAmount     int
	LoanResult     string
}

func main() {

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	loanApp := LoanApplication{}

	loanApp.AppVersion = os.Getenv("APP_VERSION")
	if len(loanApp.AppVersion) == 0 {
		loanApp.AppVersion = "dev"
	}

	loanApp.BackendHost = os.Getenv("BACKEND_HOST")
	if len(loanApp.BackendHost) == 0 {
		loanApp.BackendHost = "interest"
	}

	loanApp.BackendPort = os.Getenv("BACKEND_PORT")
	if len(loanApp.BackendPort) == 0 {
		loanApp.BackendPort = "8080"
	}

	// Allow anybody to retrieve version
	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, loanApp.AppVersion)
	})

	// Kubernetes check if app is ok
	http.HandleFunc("/health/live", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "up")
	})

	// Kubernetes check if app can serve requests
	http.HandleFunc("/health/ready", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "yes")
	})

	http.HandleFunc("/", loanApp.serveFiles)

	fmt.Printf("Frontend version %s is listening now at port %s\n", loanApp.AppVersion, port)
	err := http.ListenAndServe(":"+port, nil)
	log.Fatal(err)
}

func (loanApp *LoanApplication) serveFiles(w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path
	p := "." + upath
	if p == "./" {
		loanApp.home(w, r)
		return
	} else {
		p = filepath.Join("./static/", path.Clean(upath))
	}
	http.ServeFile(w, r, p)
}

func (loanApp *LoanApplication) home(w http.ResponseWriter, r *http.Request) {

	loanApp.handleFormSubmission(w, r)

	t, err := template.ParseFiles("./static/index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error parsing template: %v", err)
		return
	}
	err = t.Execute(w, loanApp)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error executing template: %v", err)
		return
	}
}

func (loanApp *LoanApplication) handleFormSubmission(w http.ResponseWriter, r *http.Request) {
	loanAmount := parseLoanAmount(r)
	loanApp.LoanAmount = loanAmount
	if loanAmount == 0 {
		return
	}

	quote := ""
	interestFound, err := loanApp.getInterestRate()
	if err != nil {
		log.Println("Interest error :", err)
		quote = "Could not get interest. Sorry!"
	} else {
		quote = offerQuote(loanAmount, interestFound)
	}
	loanApp.LoanResult = quote

}

func parseLoanAmount(r *http.Request) int {

	err := r.ParseForm() // Parses the request body
	if err != nil {
		return 0
	}

	loanPostParameter := r.Form.Get("loan") // x will be "" if parameter is not set

	loanAmount, err := strconv.Atoi(loanPostParameter)
	if err != nil {
		return 0
	}
	return loanAmount

}

func offerQuote(loan int, interest int) string {
	if loan <= 0 {
		return ""
	}

	total := loan * interest / 100
	return fmt.Sprintf("With rate %d%% you will pay  %d extra interest", interest, total)

}

func (loanApp *LoanApplication) getInterestRate() (rate int, err error) {
	url, err := url.Parse("http://interest:8080/api/v1/interest")
	if err != nil {
		log.Fatal(err)
	}
	url.Host = loanApp.BackendHost + ":" + loanApp.BackendPort

	resp, err := http.Get(url.String())
	if err != nil {
		log.Printf("Could not access %s, got %s\n ", url, err)
		return 0, errors.New("Could not access " + url.String())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Non-OK HTTP status:", resp.StatusCode)
		return 0, errors.New("Could not access " + url.String())
	}

	log.Printf("Response status of %s: %s\n", url, resp.Status)

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return 0, err
	}
	log.Println("Found interest rate " + buf.String())
	return strconv.Atoi(buf.String())
}
