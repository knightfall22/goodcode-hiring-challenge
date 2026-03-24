package catalog

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/mytheresa/go-hiring-challenge/models/mocks"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestCatalogGetHandler(t *testing.T) {
	mockPrice, _ := decimal.NewFromString("50.00")
	mockProducts := &models.ProductList{
		Products: []models.Product{
			{
				Code:  "PROD-123",
				Price: mockPrice,
				Category: []models.Category{
					{Name: "Shoes"},
				},
			},
		},
		TotalProducts: 25,
	}

	tests := []struct {
		name           string
		targetURL      string
		setupMock      func(m *mocks.DataStore)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:      "Success - Default Parameters",
			targetURL: "/catalog",
			setupMock: func(m *mocks.DataStore) {
				expectedFilter := &models.GetProductsFilter{
					Limit: 10, Page: 1, Offset: 0, Category: "", Price: 0.0,
				}
				m.On("GetAllProducts", expectedFilter).Return(mockProducts, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"code":"PROD-123"`,
		},
		{
			name:      "Success - With Filters & Pagination",
			targetURL: "/catalog?limit=5&page=2&category=Shoes&price=100.50",
			setupMock: func(m *mocks.DataStore) {
				expectedFilter := &models.GetProductsFilter{
					Limit: 5, Page: 2, Offset: 5, Category: "Shoes", Price: 100.50,
				}
				m.On("GetAllProducts", expectedFilter).Return(mockProducts, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"limit":5`,
		},
		{
			name:           "Error - Invalid Limit String",
			targetURL:      "/catalog?limit=abc",
			setupMock:      func(m *mocks.DataStore) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid 'limit' parameter",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mocks.DataStore)
			tc.setupMock(mockRepo)

			logger := logrus.New()
			logger.SetFormatter(&logrus.JSONFormatter{})
			handler := NewCatalogHandler(mockRepo, logger)

			req := httptest.NewRequest(http.MethodGet, tc.targetURL, nil)
			w := httptest.NewRecorder()

			handler.HandleGet(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tc.expectedBody)

			mockRepo.AssertExpectations(t)
		})
	}

}
