package product

import (
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/sirupsen/logrus"
)

type ProductHandler struct {
	repo   models.DataStore
	logger *logrus.Logger
}

type Response struct {
	Product Product `json:"product"`
}

type Product struct {
	Code     string    `json:"code"`
	Price    float64   `json:"price"`
	Category []string  `json:"category"`
	Variant  []Variant `json:"variant"`
}

type Variant struct {
	Name  string  `json:"name"`
	SKU   string  `json:"sku"`
	Price float64 `json:"price"`
}

func NewProductHandler(r models.DataStore, log *logrus.Logger) *ProductHandler {
	return &ProductHandler{
		repo:   r,
		logger: log.WithField("module", "Product").Logger,
	}
}

func (h *ProductHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	useLogger := h.logger.WithField("product", "HandleGet").Logger
	code := r.PathValue("code")

	product, err := h.repo.GetProduct(code)
	if err != nil {
		useLogger.WithError(err).Error("cannot find product")
		api.ErrorResponse(w, http.StatusNotFound, "Product not found")
		return
	}

	// Map response
	var variants []Variant
	for _, v := range product.Variants {
		variants = append(variants, Variant{
			Name:  v.Name,
			SKU:   v.SKU,
			Price: v.Price.InexactFloat64(),
		})
	}

	categories := make([]string, len(product.Category))
	for i, c := range product.Category {
		categories[i] = c.Name
	}
	response := Response{
		Product{
			Code:     product.Code,
			Price:    product.Price.InexactFloat64(),
			Category: categories,
			Variant:  variants,
		},
	}
	api.OKResponse(w, http.StatusOK, api.ApiResponse[Response]{
		Success: true,
		Message: "Product found",
		Data:    response,
	})

}
