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
	router.HandleFunc("/", server.homePage).Methods("GET")
	router.HandleFunc("/stocks/{ticker}", server.stockPage).Methods("GET")

	router.MethodNotAllowedHandler = CustomerErrorHandler(
		http.StatusMethodNotAllowed,
		"Method not allowed.",
	)
	router.NotFoundHandler = CustomerErrorHandler(
		http.StatusNotFound,
		"Page not found.",
	)

	log.Printf("Server started on %s", server.HostAndPort)
	http.ListenAndServe(server.HostAndPort, router)
}

func (server *Server) homePage(writer http.ResponseWriter, request *http.Request) {
	WriteHTML(writer, http.StatusOK, "index.html", nil)
}

func (server *Server) stockPage(writer http.ResponseWriter, request *http.Request) {
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
	template, err := template.ParseFiles(fmt.Sprintf("./templates/%s", templateName))
	if err != nil {
		return err
	}
	return template.Execute(writer, value)
}
