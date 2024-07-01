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
	CurrentRole    string
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

	loanApp.readCurrentConfiguration()

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
	} else if p == "./diagram.svg" {
		loanApp.showDiagram(w, r)
		return
	} else {
		p = filepath.Join("./static/", path.Clean(upath))
	}
	http.ServeFile(w, r, p)
}

func (loanApp *LoanApplication) findBackendVersion() {
	version, err := loanApp.callBackend("version")
	if err != nil {
		log.Println("Interest error :", err)
		version = "unknown"
	}

	loanApp.BackendVersion = version
}

func (loanApp *LoanApplication) home(w http.ResponseWriter, r *http.Request) {

	loanApp.findBackendVersion()
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

func (loanApp *LoanApplication) showDiagram(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("./static/diagram.svg")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error parsing template: %v", err)
		return
	}

	type versions struct {
		FV string
		BV string
	}

	versionsFound := versions{}
	versionsFound.FV = loanApp.AppVersion
	versionsFound.BV = loanApp.BackendVersion

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Accept-Ranges", "bytes")

	err = t.Execute(w, versionsFound)
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
	interestFound, err := loanApp.callBackend("api/v1/interest")
	if err != nil {
		log.Println("Interest error :", err)
		quote = "Could not get interest. Sorry!"
	} else {
		log.Println("Found interest rate " + interestFound)
		interestConverted, _ := strconv.Atoi(interestFound)
		quote = offerQuote(loanAmount, interestConverted)
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

func (loanApp *LoanApplication) callBackend(path string) (result string, err error) {

	backendUrl := url.URL{
		Scheme: "http",
		Host:   loanApp.BackendHost + ":" + loanApp.BackendPort,
		Path:   path,
	}

	url := backendUrl.String()
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Could not access %s, got %s\n ", url, err)
		return "", errors.New("Could not access " + url)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Non-OK HTTP status:", resp.StatusCode)
		return "", errors.New("Could not access " + url)
	}

	log.Printf("Response status of %s: %s\n", url, resp.Status)

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
