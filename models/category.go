package models

//Represent product category.
type Category struct {
	//Todo: Change this to id
	ID        uint   `gorm:"primaryKey"`
	Code      string `gorm:"uniqueIndex;not null"`
	Name      string `gorm:"not null"`
	ProductID uint   `gorm:"not null"`
}

func (v *Category) TableName() string {
	return "categories"
}
