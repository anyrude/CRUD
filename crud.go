package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
    "fmt"
	"net/http"
)

type Person struct {
	ID        string   `json:"id,omitempty"`
	Firstname string   `json:"firstname,omitempty"`
	Lastname  string   `json:"lastname,omitempty"`
	Address   *Address `json:"address,omitempty"`
}
type Address struct {
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
}

var people []Person

var database string = "yoyoyo"


var db *mgo.Database

const (
	COLLECTION = "peop"
)

func GetPeople(w http.ResponseWriter, r *http.Request) {
	err := db.C(COLLECTION).Find(bson.M{}).All(&people)
	if err !=nil {
		 fmt.Println("failed to get people",err)
        return
	}
	json.NewEncoder(w).Encode(people)
}
func GetPerson(w http.ResponseWriter, r *http.Request) {
	//defer r.Body.Close()
	params := mux.Vars(r)
	var person Person
	person.ID = params["id"]
	//for _, item := range people {
	//	if item.ID == params["id"] {
			err := db.C(COLLECTION).Find(bson.M{"id" : person.ID }).One(&person)
			//err := db.C(COLLECTION).Find(bson.M{}).One(&item)
			json.NewEncoder(w).Encode(&person)
			if err!=nil {
				fmt.Println("failed to get person doc",err)
				return
						}
			
	//	return
	//	}
	//}
}
func CreatePerson(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	params := mux.Vars(r)
	var person Person

	err := json.NewDecoder(r.Body).Decode(&person)
	 if err!=nil{
	    fmt.Println("failed to decode payload",err)
        return
	}
	person.ID = params["id"]
	err=db.C(COLLECTION).Insert(&person)

	//defer session.Close()

	 if err!=nil{
	    fmt.Println("failed to insert person doc",err)
        return
	}
	json.NewEncoder(w).Encode(person)	
}

func UpdatePerson(w http.ResponseWriter, r *http.Request) {
defer r.Body.Close()
	var person Person
	if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
		fmt.Println("failed to decode payload", err)
		return
	}
	
    if err := db.C(COLLECTION).Update(bson.M{"id": person.ID}, &person); err!=nil{
	    fmt.Println("failed to update person because : ",err)
        return
	}

	json.NewEncoder(w).Encode(person)	



}

func DeletePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range people {
		if item.ID == params["id"] {
			err := db.C(COLLECTION).Remove(&item)
			json.NewEncoder(w).Encode(item)
			if err!=nil {
				fmt.Println("failed to delete person doc",err)
				return
			}
			break
		}

	}

	
}

// our main function
func main() {
    fmt.Println("basic curd api service  running on port : 8040 ")

    session, err := mgo.Dial("mongodb://localhost:27017/" + database)
	if err != nil {
		fmt.Println(err)
	}

 

   // fmt.Println("The session: ",session)
	db = session.DB(database)
	



	router := mux.NewRouter()
	router.HandleFunc("/people", GetPeople).Methods("GET")
	router.HandleFunc("/people/{id}", GetPerson).Methods("GET")
	router.HandleFunc("/people/{id}", CreatePerson).Methods("POST")
	router.HandleFunc("/people/{id}", DeletePerson).Methods("DELETE")
	router.HandleFunc("/people", UpdatePerson).Methods("PUT")
	log.Fatal(http.ListenAndServe(":8040", router))
}
