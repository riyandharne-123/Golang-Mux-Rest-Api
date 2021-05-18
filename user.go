package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/thedevsaddam/govalidator"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

const DSN = "root:@tcp(localhost:3306)/golang-rest-api?parseTime=true"

type User struct {
	gorm.Model
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
}

func InitialMigration() {
	DB, err = gorm.Open(mysql.Open(DSN), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("Cannot Connect to DB")
	}
	DB.AutoMigrate(&User{})
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var users []User
	DB.Find(&users)
	json.NewEncoder(w).Encode(users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var user User
	DB.First(&user, params["id"])
	json.NewEncoder(w).Encode(user)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user User

	//validation
	rules := govalidator.MapData{
		"firstname": []string{"required"},
		"lastname":  []string{"required"},
		"email":     []string{"required", "email"},
	}

	opts := govalidator.Options{
		Request: r,
		Data:    &user,
		Rules:   rules,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()
	err := map[string]interface{}{"validationError": e}

	if len(e) != 0 {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewDecoder(r.Body).Decode(&user)
		DB.Create(&user)
		json.NewEncoder(w).Encode(user)
	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var user User

	//validation
	rules := govalidator.MapData{
		"firstname": []string{"required"},
		"lastname":  []string{"required"},
		"email":     []string{"required", "email"},
	}

	opts := govalidator.Options{
		Request: r,
		Data:    &user,
		Rules:   rules,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()
	err := map[string]interface{}{"validationError": e}

	if len(e) != 0 {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewDecoder(r.Body).Decode(&user)
		DB.First(&user, params["id"])
		DB.Save(&user)
		json.NewEncoder(w).Encode(user)
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var user User
	DB.Delete(&user, params["id"])
	json.NewEncoder(w).Encode("The User is Deleted Successfully!")
}
