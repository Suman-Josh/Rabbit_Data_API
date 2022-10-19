package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/suman/task/connection"
	"github.com/suman/task/models"
	"github.com/suman/task/rabbit"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

func publishToDB(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var data models.Data
	json.NewDecoder(r.Body).Decode(&data)
	DB.Create(&data)
	json.NewEncoder(w).Encode(data)
	rabbit.Messages <- data.Msg

}

func main() {
	rabbit.SetupRabbitMQ()
	defer rabbit.Conn.Close()
	defer rabbit.Ch.Close()
	DB, err = connection.SetupDB()
	if err != nil {
		log.Panic(err)
		return
	}
	fmt.Println("Connected to database")

	DB.AutoMigrate(&models.Data{})
	go rabbit.SendMessage()
	go rabbit.ReadMessage("taskQueue1")
	go rabbit.ReadMessage("taskQueue2")
	http.HandleFunc("/", publishToDB)
	http.ListenAndServe(":8000", nil)
}
