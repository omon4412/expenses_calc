package models

type Category struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"category_id"`
	Name        string `gorm:"not null" json:"name"`
	Description string `gorm:"" json:"description"`
	OwnerId     uint   `gorm:"foreignKey:UserID" json:"-"`
}

var DefaultCategories = []Category{
	{Name: "Еда", Description: "Расходы на еду", OwnerId: 0},
	{Name: "Транспорт", Description: "Расходы на транспорт", OwnerId: 0},
	{Name: "Развлечения", Description: "Кино, рестораны и другие развлечения", OwnerId: 0},
	{Name: "Здоровье", Description: "Расходы на здоровье, медицинские услуги", OwnerId: 0},
}
