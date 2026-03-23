package category

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/mytheresa/go-hiring-challenge/models/mocks"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCategoryGetHandler(t *testing.T) {
	mockCategories := models.CategoryList{
		Categories: []models.Category{
			{Name: "Shoes", Code: "CAT-SHOES"},
			{Name: "Dresses", Code: "CAT-DRESSES"},
		},
		Total: 20,
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
			targetURL: "/categories",
			setupMock: func(m *mocks.DataStore) {
				expectedFilter := &models.GetCategoryFilter{
					Limit:  10,
					Page:   1,
					Offset: 0,
				}
				m.On("GetAllCategories", expectedFilter).Return(&mockCategories, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"name":"Shoes"`,
		},
		{
			name:      "Success - Custom Pagination",
			targetURL: "/categories?limit=5&page=3",
			setupMock: func(m *mocks.DataStore) {
				expectedFilter := &models.GetCategoryFilter{
					Limit:  5,
					Page:   3,
					Offset: 10,
				}
				m.On("GetAllCategories", expectedFilter).Return(&mockCategories, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"name":"Shoes"`,
		},
		{
			name:      "Success - Limit Out of Bounds (Resets to 10)",
			targetURL: "/categories?limit=150&page=1",
			setupMock: func(m *mocks.DataStore) {
				expectedFilter := &models.GetCategoryFilter{
					Limit:  10,
					Page:   1,
					Offset: 0,
				}
				m.On("GetAllCategories", expectedFilter).Return(&mockCategories, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"code":"CAT-SHOES"`,
		},
		{
			name:           "Error - Invalid Limit Parameter",
			targetURL:      "/categories?limit=abc",
			setupMock:      func(m *mocks.DataStore) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid 'limit' parameter",
		},
		{
			name:           "Error - Invalid Page Parameter",
			targetURL:      "/categories?page=xyz",
			setupMock:      func(m *mocks.DataStore) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid 'limit' parameter",
		},
		{
			name:      "Error - Repository Failure",
			targetURL: "/categories",
			setupMock: func(m *mocks.DataStore) {
				m.On("GetAllCategories", mock.Anything).Return(
					nil,
					errors.New("database connection lost"),
				)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "database connection lost",
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

func TestCategoryPostHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		setupMock      func(m *mocks.DataStore)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "Success - Category Added",
			requestBody: `{"name": "Sneakers", "code": "CAT-SNEAK", "product_id": 10}`,
			setupMock: func(m *mocks.DataStore) {
				m.On("CheckProductExists", uint(10)).Return(true, nil)

				expectedAdd := models.AddCategory{
					Name:      "Sneakers",
					Code:      "CAT-SNEAK",
					ProductID: 10,
				}
				m.On("AddCategory", expectedAdd).Return(nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `"message":"Category added"`,
		},
		{
			name:           "Error - Malformed JSON",
			requestBody:    `{"name": "Sneakers", "code": "CAT-SNEAK",`,
			setupMock:      func(m *mocks.DataStore) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid JSON body",
		},
		{
			name:           "Error - Unknown Field in JSON",
			requestBody:    `{"name": "Sneakers", "code": "CAT-SNEAK", "product_id": 10, "extra": "invalid"}`,
			setupMock:      func(m *mocks.DataStore) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid JSON body",
		},
		{
			name:        "Error - CheckProductExists Fails (DB Error)",
			requestBody: `{"name": "Sneakers", "code": "CAT-SNEAK", "product_id": 11}`,
			setupMock: func(m *mocks.DataStore) {
				m.On("CheckProductExists", uint(11)).Return(false, errors.New("database connection lost"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "internal failure",
		},
		{
			name:        "Error - Product Not Found",
			requestBody: `{"name": "Sneakers", "code": "CAT-SNEAK", "product_id": 99}`,
			setupMock: func(m *mocks.DataStore) {
				m.On("CheckProductExists", uint(99)).Return(false, nil)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "product not found",
		},
		{
			name:        "Error - AddCategory Fails (DB Error)",
			requestBody: `{"name": "Sneakers", "code": "CAT-SNEAK", "product_id": 12}`,
			setupMock: func(m *mocks.DataStore) {
				m.On("CheckProductExists", uint(12)).Return(true, nil)
				m.On("AddCategory", mock.AnythingOfType("models.AddCategory")).Return(errors.New("db write failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "internal failure",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			mockRepo := new(mocks.DataStore)
			tc.setupMock(mockRepo)

			logger := logrus.New()
			logger.SetFormatter(&logrus.JSONFormatter{})
			handler := NewCatalogHandler(mockRepo, logger)

			bodyReader := bytes.NewReader([]byte(tc.requestBody))
			req := httptest.NewRequest(http.MethodPost, "/categories", bodyReader)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			handler.HandlePost(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tc.expectedBody)

			mockRepo.AssertExpectations(t)
		})
	}
}
