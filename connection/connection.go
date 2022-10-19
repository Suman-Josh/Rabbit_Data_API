package connection

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	Host     = "localhost"
	User     = "postgres"
	Password = "root"
	Name     = "postgres"
	Port     = "5432"
)

func SetupDB() (*gorm.DB, error) {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		Host,
		Port,
		User,
		Name,
		Password,
	)

	db, e := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if e != nil {
		return nil, e
	}
	return db, nil
}
