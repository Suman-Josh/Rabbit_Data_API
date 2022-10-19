package models

import "gorm.io/gorm"

type Data struct {
	gorm.Model
	Msg string `json:"msg"`
	Key string `json:"key"`
}
