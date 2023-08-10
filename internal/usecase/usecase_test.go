package usecase

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/nurtidev/kmf/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) Save(ctx context.Context, currency *model.Currency) error {
	args := m.Called(ctx, currency)
	return args.Error(0)
}

func (m *MockRepo) FindByDateAndCode(ctx context.Context, date time.Time, code string) ([]*model.Currency, error) {
	args := m.Called(ctx, date, code)
	return args.Get(0).([]*model.Currency), args.Error(1)
}

func TestUseCase_SaveCurrenciesAsync(t *testing.T) {
	mockRepo := new(MockRepo)
	useCase := UseCase{Repo: mockRepo}

	currencies := []*model.Currency{
		{Title: "Currency 1", Code: "USD", Value: 123.45, Date: time.Now()},
		{Title: "Currency 2", Code: "EUR", Value: 456.78, Date: time.Now()},
	}

	mockRepo.On("Save", mock.Anything, currencies[0]).Return(nil).Once()
	mockRepo.On("Save", mock.Anything, currencies[1]).Return(errors.New("error saving")).Once()

	err := useCase.SaveCurrenciesAsync(context.Background(), currencies)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUseCase_FetchRates(t *testing.T) {
	// Создание mock HTTP-сервера
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			<xml>
				<items>
					<item>
						<title>USD</title>
						<description>123.45</description>
						<!-- Остальные поля ... -->
					</item>
				</items>
			</xml>
		`))
	}))
	defer server.Close()

	useCase := UseCase{}

	rates, err := useCase.FetchRates(time.Now())
	assert.NoError(t, err)
	assert.Len(t, rates, 1)
	assert.Equal(t, "USD", rates[0].Title)
	assert.Equal(t, 123.45, rates[0].Description)
}
