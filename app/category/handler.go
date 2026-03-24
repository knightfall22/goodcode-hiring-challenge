package category

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/sirupsen/logrus"
)

type Response struct {
	Categories []Category `json:"categories"`
}

type Category struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type CategoryHandler struct {
	repo   models.DataStore
	logger *logrus.Logger
}

func NewCatalogHandler(r models.DataStore, log *logrus.Logger) *CategoryHandler {
	return &CategoryHandler{
		repo:   r,
		logger: log.WithField("module", "Category").Logger,
	}
}

func (h *CategoryHandler) HandlePost(w http.ResponseWriter, r *http.Request) {
	useLogger := h.logger.WithField("category", "HandlePost").Logger

	var addCategory models.AddCategory

	decoder := json.NewDecoder(r.Body)

	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&addCategory); err != nil {
		useLogger.WithError(err).Error("failed to decode input")
		api.ErrorResponse(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	prodExists, err := h.repo.CheckProductExists(addCategory.ProductID)
	if err != nil {
		useLogger.WithError(err).Error("failed to check if product exists")
		api.ErrorResponse(w, http.StatusInternalServerError, "internal failure")
		return
	}

	if !prodExists {
		api.ErrorResponse(w, http.StatusNotFound, "product not found")
		return
	}

	err = h.repo.AddCategory(addCategory)
	if err != nil {
		useLogger.WithError(err).Error("failed to add category")
		api.ErrorResponse(w, http.StatusInternalServerError, "internal failure")
		return
	}

	api.OKResponse(w, http.StatusCreated, api.ApiResponse[any]{
		Success: true,
		Message: "Category added",
		Data:    nil,
	})

}

func (h *CategoryHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	useLogger := h.logger.WithField("category", "HandleGet").Logger

	filter, err := h.parseCategoryFilters(r)
	if err != nil {
		useLogger.WithError(err).Error("failed to parse filters")
		api.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.repo.GetAllCategories(filter)
	if err != nil {
		useLogger.WithError(err).Error("failed to retrieve all categories")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Map response
	categories := make([]Category, len(res.Categories))
	for i, c := range res.Categories {
		categories[i] = Category{
			Name: c.Name,
			Code: c.Code,
		}
	}

	response := Response{
		Categories: categories,
	}

	api.OKResponse(w, http.StatusOK, api.ApiResponse[Response]{
		Success: true,
		Message: "Categories found",
		Data:    response,

		Page:    filter.Page,
		Limit:   filter.Limit,
		Count:   len(res.Categories),
		HasNext: (filter.Page * filter.Limit) < res.Total,
		Total:   res.Total,
	})
}

func (h *CategoryHandler) parseCategoryFilters(r *http.Request) (*models.GetCategoryFilter, error) {
	query := r.URL.Query()

	limit := 10
	page := 1

	if limitStr := query.Get("limit"); limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil {
			return nil, fmt.Errorf("invalid 'limit' parameter")
		}
		limit = parsedLimit
	}

	if limit > 100 || limit < 1 {
		//reset limit to default when it exceeds this bound
		limit = 10
	}

	if pageStr := query.Get("page"); pageStr != "" {
		parsedPage, err := strconv.Atoi(pageStr)
		if err != nil {
			return nil, fmt.Errorf("invalid 'limit' parameter")
		}
		page = parsedPage
	}

	return &models.GetCategoryFilter{
		Limit:  limit,
		Page:   page,
		Offset: (page - 1) * limit,
	}, nil
}
