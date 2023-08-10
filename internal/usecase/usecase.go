package usecase

import (
	"context"
	"encoding/xml"
	"fmt"
	"github.com/nurtidev/kmf/internal/model"
	"github.com/nurtidev/kmf/internal/repo"
	"log"
	"net/http"
	"strconv"
	"time"
)

type UseCase struct {
	Repo repo.Repository
}

func (u *UseCase) SaveCurrenciesAsync(ctx context.Context, currencies []*model.Currency) error {
	done := make(chan struct{}) // Канал для сигнала завершения горутин
	errChan := make(chan error) // Канал для ошибок

	for _, currency := range currencies {
		go func(c *model.Currency) {
			err := u.Repo.Save(ctx, c)
			if err != nil {
				log.Printf("Error saving currency: %s", err)
				errChan <- err // Отправляем ошибку в канал
			}
			done <- struct{}{} // Отправляем сигнал о завершении в канал
		}(currency)
	}

	// Ждем завершения всех горутин
	for range currencies {
		select {
		case <-done:
		case err := <-errChan:
			return err
		}
	}

	return nil
}

func (u *UseCase) FindCurrencyByDateAndCode(ctx context.Context, date time.Time, code string) ([]*model.Currency, error) {
	return u.Repo.FindByDateAndCode(ctx, date, code)
}

func (u *UseCase) FetchRates(date time.Time) ([]model.Rate, error) {
	url := fmt.Sprintf("https://nationalbank.kz/rss/get_rates.cfm?fdate=%s", date.Format("02.01.2006"))

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var xmlRates model.XMLRates
	if err := xml.NewDecoder(resp.Body).Decode(&xmlRates); err != nil {
		return nil, err
	}

	// Преобразуем XML структуру в нашу структуру Rate
	rates := make([]model.Rate, len(xmlRates.Items))
	for i, xmlRate := range xmlRates.Items {
		description, err := strconv.ParseFloat(xmlRate.Description, 64)
		if err != nil {
			return nil, err
		}
		rates[i] = model.Rate{
			FullName:    xmlRate.FullName,
			Title:       xmlRate.Title,
			Description: description,
			Quant:       xmlRate.Quant,
			Index:       xmlRate.Index,
			Change:      xmlRate.Change,
		}
	}

	return rates, nil
}
