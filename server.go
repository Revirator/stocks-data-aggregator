package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/revirator/cfd/companydb"
	"github.com/revirator/cfd/external"
	"github.com/revirator/cfd/view"
)

type Server struct {
	HostAndPort string
	Database    companydb.CompanyDatabase
}

func ServerInit(hostAndPort string, db *sql.DB) *Server {
	return &Server{
		HostAndPort: hostAndPort,
		Database:    companydb.CompanyDatabase{DB: db},
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
	company, err := server.Database.GetCompanyByTicker(ticker)
	if err != nil {
		if err.Error() == "company missing" {
			w.WriteHeader(http.StatusNotFound)
			view.Error(fmt.Sprintf("Company with ticker '%s' does not exist or is not listed on any of the US exchanges.", ticker)).Render(r.Context(), w)
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
			server.Database.UpdateCompanyFinancialsByTicker(ticker, company.Financials)
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
