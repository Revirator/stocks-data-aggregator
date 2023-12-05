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
	router.HandleFunc("/stocks/{ticker}", server.showStockPage).Methods("GET")

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

	http.Redirect(writer, request, fmt.Sprintf("/stocks/%s", input), http.StatusSeeOther)
}

func (server *Server) showStockPage(writer http.ResponseWriter, request *http.Request) {
	ticker := strings.ToUpper(mux.Vars(request)["ticker"])
	stock, err := server.getStockByTicker(ticker)
	if err != nil {
		WriteHTML(writer, err.StatusCode, "error.html", err.Message)
		return
	}
	WriteHTML(writer, http.StatusOK, "stock.html", stock)
}

func (server *Server) getStockByTicker(ticker string) (*Stock, *ServerError) {
	stock, err := server.Database.GetStockByTicker(ticker)
	if err != nil {
		return nil, err
	}

	if stock.Financials == nil {
		log.Printf("Requesting data from Edgar for stock with ticker '%s'", ticker)
		facts := GetFinancialFactsForCompanyGivenCIK(stock.CIK)
		if facts != nil {
			stock.Financials = MapFinancialFactsToFinancialMetrics(facts)
			server.Database.UpdateStockFinancialsByTicker(ticker, stock.Financials)
		}
	}

	return stock, nil
}

func WriteHTML(writer http.ResponseWriter, statusCode int, templateName string, value any) error {
	writer.Header().Add("Content-Type", "text/html")
	writer.WriteHeader(statusCode)
	template := template.Must(template.ParseFiles(fmt.Sprintf("./static/%s", templateName)))
	return template.Execute(writer, value)
}
