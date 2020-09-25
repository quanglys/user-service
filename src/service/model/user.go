package model

type (
	UserID int
	Status string
	Gender string
)

const (
	StatusActive   Status = "ACTIVE"
	StatusInactive Status = "INACTIVE"
)

func (s Status) IsValid() bool {
	return s == StatusActive || s == StatusInactive
}

const (
	Female Gender = "FEMALE"
	Male   Gender = "MALE"
)

func (g Gender) IsValid() bool {
	return g == Female || g == Male
}

type User struct {
	ID     UserID `gorm:"column:id" json:"id"`
	Name   string `gorm:"column:name" json:"name"`
	Gender Gender `gorm:"column:gender" json:"gender"`
	Status *Status `gorm:"column:status;default:null" json:"status"`
}
