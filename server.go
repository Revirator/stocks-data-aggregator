package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/revirator/cfd/external"
	"github.com/revirator/cfd/model"
	"github.com/revirator/cfd/view"
)

type Server struct {
	HostAndPort string
	Database    *sql.DB
}

func ServerInit(hostAndPort string, db *sql.DB) *Server {
	return &Server{
		HostAndPort: hostAndPort,
		Database:    db,
	}
}

func (server *Server) Run() {
	router := mux.NewRouter()
	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
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

func (server *Server) showHomePage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		view.Error("Could not parse form data! Please try again.").Render(r.Context(), w)
		return
	}

	input := r.Form.Get("searchBox")
	if input == "" {
		view.Index().Render(r.Context(), w)
		return
	}

	input = strings.ToUpper(input)
	http.Redirect(w, r, fmt.Sprintf("/companies/%s", input), http.StatusSeeOther)
}

func (server *Server) showCompanyPage(w http.ResponseWriter, r *http.Request) {
	ticker := strings.ToUpper(mux.Vars(r)["ticker"])
	company, err := getCompanyByTicker(ticker, server.Database)
	if err != nil {
		if err.Error() == "company missing" {
			errorMessage := "Company with ticker '%s' does not exist or is not listed on any of the US exchanges."
			w.WriteHeader(http.StatusNotFound)
			view.Error(fmt.Sprintf(errorMessage, ticker)).Render(r.Context(), w)
		} else {
			log.Printf("Error when retrieving information for company with ticker '%s'. Root cause:\n%s", ticker, err)
			w.WriteHeader(http.StatusInternalServerError)
			view.Error("Something went wrong! Please try again later.").Render(r.Context(), w)
		}
		return
	}

	if company.Financials == nil {
		log.Printf("Requesting data from 'sec.gov' for company with ticker '%s'", ticker)
		facts := external.GetFinancialFactsForCompanyGivenCIK(company.CIK)
		if facts != nil {
			company.Financials = MapFinancialFactsToFinancialMetrics(facts)
			updateCompanyFinancialsByTicker(ticker, company.Financials, server.Database)
		}
	}

	stockPrice, percentageChange := 0.0, 0.0
	log.Printf("Requesting metadata from 'finance.yahoo.com' for company with ticker '%s'", ticker)
	metadata := external.GetCompanyMetadataGivenTicker(ticker)
	if metadata != nil {
		stockPrice = metadata.RegularMarketPrice
		percentageChange = 100 * (stockPrice - metadata.ChartPreviousClose) / metadata.ChartPreviousClose
	}

	view.Company(company, stockPrice, percentageChange).Render(r.Context(), w)
}

func getCompanyByTicker(ticker string, db *sql.DB) (*model.Company, error) {
	query := "SELECT * FROM COMPANIES WHERE TICKER = $1"
	rows, err := db.Query(query, ticker)

	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, errors.New("company missing")
	}

	var exchange, financialData sql.NullString
	company := &model.Company{}
	err = rows.Scan(
		&company.Ticker,
		&company.CIK,
		&company.Name,
		&exchange,
		&financialData,
	)
	if err != nil {
		return nil, err
	}

	if financialData.Valid {
		err = json.Unmarshal([]byte(financialData.String), &company.Financials)
		if err != nil {
			return nil, err
		}
	}

	company.Exchange = &exchange.String
	return company, nil
}

func updateCompanyFinancialsByTicker(ticker string, financials map[string]*model.FinancialMetric, db *sql.DB) {
	financialData, err := json.Marshal(financials)
	if err != nil {
		log.Println(err)
		return
	}

	query := "UPDATE COMPANIES SET FINANCIALS = $1 WHERE TICKER = $2"
	_, err = db.Exec(query, financialData, ticker)
	if err != nil {
		log.Println(err)
	}
}
