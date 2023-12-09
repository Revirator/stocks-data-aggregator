package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
)

type Server struct {
	HostAndPort string
	Database    *Database
}

func ServerInit(hostAndPort string, database *Database) *Server {
	return &Server{
		HostAndPort: hostAndPort,
		Database:    database,
	}
}

func (server *Server) Run() {
	router := mux.NewRouter()
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	router.HandleFunc("/", server.showHomePage).Methods("GET")
	router.HandleFunc("/companies/{ticker}", server.showCompanyPage).Methods("GET")

	router.MethodNotAllowedHandler = CustomErrorHandler(
		http.StatusMethodNotAllowed,
		"Method not allowed.",
	)
	router.NotFoundHandler = CustomErrorHandler(
		http.StatusNotFound,
		"Page not found.",
	)

	log.Printf("Server started on %s", server.HostAndPort)
	http.ListenAndServe(server.HostAndPort, router)
}

func (server *Server) showHomePage(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		WriteHTML(writer, http.StatusInternalServerError, "error.html", "Could not parse form data! Please try again.")
		return
	}

	input := request.Form.Get("searchBox")
	if input == "" {
		WriteHTML(writer, http.StatusOK, "index.html", nil)
		return
	}

	input = strings.ToUpper(input)
	http.Redirect(writer, request, fmt.Sprintf("/companies/%s", input), http.StatusSeeOther)
}

func (server *Server) showCompanyPage(writer http.ResponseWriter, request *http.Request) {
	ticker := strings.ToUpper(mux.Vars(request)["ticker"])
	company, err := server.getCompanyByTicker(ticker)
	if err != nil {
		WriteHTML(writer, err.StatusCode, "error.html", err.Message)
		return
	}
	WriteHTML(writer, http.StatusOK, "company.html", company)
}

func (server *Server) getCompanyByTicker(ticker string) (*Company, *ServerError) {
	company, err := server.Database.GetCompanyByTicker(ticker)
	if err != nil {
		return nil, err
	}

	if company.Financials == nil {
		log.Printf("Requesting data from Edgar for company with ticker '%s'", ticker)
		facts := GetFinancialFactsForCompanyGivenCIK(company.CIK)
		if facts != nil {
			company.Financials = MapFinancialFactsToFinancialMetrics(facts)
			server.Database.UpdateCompanyFinancialsByTicker(ticker, company.Financials)
		}
	}

	return company, nil
}

func WriteHTML(writer http.ResponseWriter, statusCode int, templateName string, value any) error {
	writer.Header().Add("Content-Type", "text/html")
	writer.WriteHeader(statusCode)
	template := template.Must(template.ParseFiles(fmt.Sprintf("./static/%s", templateName)))
	return template.Execute(writer, value)
}
