package models

import "gorm.io/gorm"

type Tenant struct {
	gorm.Model
	Id int64 `gorm:"primaryKey;autoIncrement" json:"ID,omitempty"`

	Name string `json:"name,omitempty"`
}
