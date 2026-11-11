package model

type Todo struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	Title     string `json:"title" gorm:"not null"`
	Completed bool   `json:"completed" gorm:"not null;default:false"`
}
