package handler

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/nurtidev/kmf/internal/model"
	"github.com/nurtidev/kmf/internal/usecase"
	"log"
	"net/http"
	"time"
)

type Handler struct {
	UseCase *usecase.UseCase
}

func (h *Handler) SaveCurrency(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dateStr := vars["date"]
	date, err := time.Parse("02.01.2006", dateStr)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	rates, err := h.UseCase.FetchRates(date)
	if err != nil {
		log.Printf("Failed to fetch rates: %v\n", err)
		http.Error(w, "Failed to fetch rates", http.StatusBadRequest)
		return
	}

	currencies := convertRatesToCurrencies(rates, date)
	if err = h.UseCase.SaveCurrenciesAsync(context.TODO(), currencies); err != nil {
		log.Printf("Failed to save currency: %v\n", err)
		http.Error(w, "Failed to save currency", http.StatusInternalServerError)
		return
	}

	response := map[string]bool{"success": true}
	if err = json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode response: %v\n", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *Handler) GetCurrency(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dateStr := vars["date"]
	code := vars["code"]

	date, err := time.Parse("02.01.2006", dateStr)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	currencies, err := h.UseCase.FindCurrencyByDateAndCode(context.TODO(), date, code)
	if err != nil {
		log.Printf("Failed to retrieve currency: %v", err)
		http.Error(w, "Failed to retrieve currency", http.StatusNotFound)
		return
	}

	if err = json.NewEncoder(w).Encode(currencies); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusNotFound)
		return
	}
}

func convertRatesToCurrencies(rates []model.Rate, date time.Time) []*model.Currency {
	currencies := make([]*model.Currency, len(rates))
	for i, rate := range rates {
		currencies[i] = &model.Currency{
			Title: rate.FullName,
			Code:  rate.Title,
			Value: rate.Description,
			Date:  date,
		}
	}
	return currencies
}
