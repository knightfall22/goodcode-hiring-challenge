package models

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ProductList struct {
	Products      []Product
	TotalProducts int
}

type CategoryList struct {
	Categories []Category
	Total      int
}

type GetProductsFilter struct {
	Limit    int
	Offset   int
	Category string
	Price    float64
	Page     int
}

type AddCategory struct {
	Name      string `json:"name"`
	Code      string `json:"code"`
	ProductID uint   `json:"product_id"`
}

type GetCategoryFilter struct {
	Limit  int
	Offset int
	Page   int
}

//go:generate mockery --name DataStore
type DataStore interface {
	GetAllProducts(query *GetProductsFilter) (*ProductList, error)
	GetProduct(code string) (*Product, error)
	GetAllCategories(query *GetCategoryFilter) (*CategoryList, error)
	CheckProductExists(id uint) (bool, error)
	AddCategory(category AddCategory) error
}

type ProductsRepository struct {
	db *gorm.DB
}

func NewProductsRepository(db *gorm.DB) *ProductsRepository {
	return &ProductsRepository{
		db: db,
	}
}

func (r *ProductsRepository) GetAllProducts(query *GetProductsFilter) (*ProductList, error) {
	var products []Product
	var total int64

	q := r.db.Model(&Product{})
	if query.Category != "" {
		q = q.Joins("JOIN categories ON categories.product_id = products.id").
			Where("categories.name = ?", query.Category)
	}

	if query.Price > 0 {
		q = q.Where("products.price < ?", query.Price)
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, err
	}

	if err := q.
		Limit(query.Limit).Offset(query.Offset).Preload("Variants").
		Preload("Category").Find(&products).Error; err != nil {
		return nil, err
	}
	return &ProductList{Products: products, TotalProducts: int(total)}, nil
}

func (r *ProductsRepository) GetProduct(code string) (*Product, error) {
	var product Product

	if err := r.db.Preload("Variants").
		Preload("Category").
		Where("code = ?", code).
		First(&product).Error; err != nil {
		return nil, err
	}

	for i := range product.Variants {
		if decimal.Zero.Equal(product.Variants[i].Price) {
			product.Variants[i].Price = product.Price
		}
	}
	return &product, nil
}

func (r *ProductsRepository) GetAllCategories(query *GetCategoryFilter) (*CategoryList, error) {
	var categories []Category
	var total int64

	q := r.db.Model(&Category{})

	if err := q.Count(&total).Error; err != nil {
		return nil, err
	}

	if err := q.Limit(query.Limit).Offset(query.Offset).Find(&categories).Error; err != nil {
		return nil, err
	}

	return &CategoryList{Categories: categories}, nil
}

func (r *ProductsRepository) CheckProductExists(id uint) (bool, error) {
	var product Product
	if err := r.db.Where("id = ?", id).First(&product).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (r *ProductsRepository) AddCategory(input AddCategory) error {
	newCategory := Category{
		ProductID: uint(input.ProductID),
		Name:      input.Name,
		Code:      input.Code,
	}

	// 2. Pass the POINTER of the model to GORM
	return r.db.Create(&newCategory).Error
}
