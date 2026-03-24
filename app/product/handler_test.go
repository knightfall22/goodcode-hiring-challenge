package product

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/mytheresa/go-hiring-challenge/models/mocks"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestProductGetHandler(t *testing.T) {
	mockProductPrice, _ := decimal.NewFromString("199.99")
	mockVariantPrice, _ := decimal.NewFromString("209.99")

	mockProduct := models.Product{
		Code:  "PROD-999",
		Price: mockProductPrice,
		Category: []models.Category{
			{Name: "Dresses"},
			{Name: "Evening"},
		},
		Variants: []models.Variant{
			{Name: "Red", SKU: "SKU-RED-1", Price: mockVariantPrice},
		},
	}

	tests := []struct {
		name           string
		codeParam      string
		setupMock      func(m *mocks.DataStore)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:      "Success - Product Found",
			codeParam: "PROD-999",
			setupMock: func(m *mocks.DataStore) {
				m.On("GetProduct", "PROD-999").Return(&mockProduct, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"code":"PROD-999"`,
		},
		{
			name:      "Success - Verify Variant and Category Mapping",
			codeParam: "PROD-999",
			setupMock: func(m *mocks.DataStore) {
				m.On("GetProduct", "PROD-999").Return(&mockProduct, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"sku":"SKU-RED-1"`,
		},
		{
			name:      "Error - Product Not Found",
			codeParam: "PROD-UNKNOWN",
			setupMock: func(m *mocks.DataStore) {
				m.On("GetProduct", "PROD-UNKNOWN").Return(nil, errors.New("record not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "Product not found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mocks.DataStore)
			tc.setupMock(mockRepo)

			logger := logrus.New()
			logger.SetFormatter(&logrus.JSONFormatter{})
			handler := NewProductHandler(mockRepo, logger)

			targetURL := "/catalog/" + tc.codeParam
			req := httptest.NewRequest(http.MethodGet, targetURL, nil)

			req.SetPathValue("code", tc.codeParam)

			w := httptest.NewRecorder()

			handler.HandleGet(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tc.expectedBody)

			mockRepo.AssertExpectations(t)
		})
	}
}
