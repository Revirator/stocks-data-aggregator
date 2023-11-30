package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type ServerFunc func(http.ResponseWriter, *http.Request) error

type ServerError struct {
	Status int    `json:"-"`
	Error  string `json:"error"`
}

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
	router.HandleFunc("/stocks/{ticker}", serverFuncHandler(server.handleStocks))

	log.Printf("Server started on %s", server.HostAndPort)
	http.ListenAndServe(server.HostAndPort, router)
}

func (server *Server) handleStocks(writer http.ResponseWriter, request *http.Request) error {
	ticker := strings.ToUpper(mux.Vars(request)["ticker"])
	if request.Method == "GET" {
		return server.getAndUpdateStockByTicker(writer, ticker, false)
	}

	if request.Method == "POST" {
		return server.getAndUpdateStockByTicker(writer, ticker, true)
	}

	return writeJSON(
		writer,
		http.StatusMethodNotAllowed,
		ServerError{Error: "Method not allowed."},
	)
}

func (server *Server) getAndUpdateStockByTicker(writer http.ResponseWriter, ticker string, updateFinancials bool) error {
	stock, err := server.getStockFromDatabase(ticker)
	if err != nil {
		return writeJSON(writer, err.Status, err)
	}

	if stock.Financials == nil || updateFinancials {
		log.Printf("Requesting data from Edgar for stock with ticker '%s'", ticker)
		concept := GetConceptForCompanyGivenCIK(stock.CIK)
		if concept != nil {
			stock.Financials = MapConceptToFinancials(concept)
			server.Database.UpdateStockFinancialsByTicker(ticker, stock.Financials)
		}
	}

	return writeJSON(writer, http.StatusOK, stock)
}

func (server *Server) getStockFromDatabase(ticker string) (*Stock, *ServerError) {
	stock, outcome := server.Database.GetStockByTicker(ticker)
	switch outcome {
	case DATABASE_ERROR:
		return nil, &ServerError{
			Status: http.StatusInternalServerError,
			Error:  "Something went wrong! Please try again later.",
		}
	case STOCK_MISSING_ERROR:
		return nil, &ServerError{
			Status: http.StatusNotFound,
			Error:  fmt.Sprintf("Stock with ticker '%s' does not exist or is not listed on any of the US exchanges.", ticker),
		}
	default:
		return stock, nil
	}
}

func serverFuncHandler(function ServerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if err := function(writer, request); err != nil {
			log.Println(err)
			writeJSON(
				writer,
				http.StatusInternalServerError,
				ServerError{Error: "Something went wrong! Please try again later."},
			)
		}
	}
}

func writeJSON(writer http.ResponseWriter, statusCode int, value any) error {
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	return json.NewEncoder(writer).Encode(value)
}
