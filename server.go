package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/revirator/cfd/companydb"
	"github.com/revirator/cfd/views"
)

type Server struct {
	HostAndPort string
	Database    companydb.CompanyDatabase
}

func ServerInit(hostAndPort string, database companydb.CompanyDatabase) *Server {
	return &Server{
		HostAndPort: hostAndPort,
		Database:    database,
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
		views.Error("Could not parse form data! Please try again.").Render(r.Context(), w)
		return
	}

	input := r.Form.Get("searchBox")
	if input == "" {
		views.Index().Render(r.Context(), w)
		return
	}

	input = strings.ToUpper(input)
	http.Redirect(w, r, fmt.Sprintf("/companies/%s", input), http.StatusSeeOther)
}

func (server *Server) showCompanyPage(w http.ResponseWriter, r *http.Request) {
	ticker := strings.ToUpper(mux.Vars(r)["ticker"])
	company, err := server.Database.GetCompanyByTicker(ticker)
	if err != nil {
		if err.Error() == "stock missing" {
			w.WriteHeader(http.StatusNotFound)
			views.Error(fmt.Sprintf("Company with ticker '%s' does not exist or is not listed on any of the US exchanges.", ticker)).Render(r.Context(), w)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			views.Error("Something went wrong! Please try again later.").Render(r.Context(), w)
		}
		return
	}

	if company.Financials == nil {
		log.Printf("Requesting data from 'sec.gov' for company with ticker '%s'", ticker)
		facts := GetFinancialFactsForCompanyGivenCIK(company.CIK)
		if facts != nil {
			company.Financials = MapFinancialFactsToFinancialMetrics(facts)
			err = server.Database.UpdateCompanyFinancialsByTicker(ticker, company.Financials)
			log.Fatal(err)
		}
	}

	var stockPrice, dayMovePercentage float64
	log.Printf("Requesting metadata from 'finance.yahoo.com' for company with ticker '%s'", ticker)
	metadata := GetCompanyMetadataGivenTicker(ticker)
	if metadata != nil {
		stockPrice = metadata.RegularMarketPrice
		dayMovePercentage = calculateDayMovePercentage(metadata)
	}

	views.Company(company, stockPrice, dayMovePercentage).Render(r.Context(), w)
}

func calculateDayMovePercentage(metadata *CompanyMetadata) float64 {
	dayMovePercentage := 100 * (metadata.RegularMarketPrice - metadata.ChartPreviousClose) / metadata.ChartPreviousClose
	return math.Round(dayMovePercentage*100) / 100
}
