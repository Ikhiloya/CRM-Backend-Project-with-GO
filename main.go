package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Customer struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Contacted bool   `json:"contacted"`
}

// dummy database of customers
var customers = []Customer{
	{
		Id:        1,
		Name:      "John Doe",
		Role:      "Project Manager",
		Email:     "john.doe@gmail.com",
		Phone:     "2347060443321",
		Contacted: true,
	},
	{
		Id:        2,
		Name:      "Charles Darwin",
		Role:      "Solutions Architect",
		Email:     "charles.darwin@gmail.com",
		Phone:     "2347060443421",
		Contacted: false,
	},
	{
		Id:        3,
		Name:      "Nick Tes",
		Role:      "Software Engineer",
		Email:     "nick.tes@gmail.com",
		Phone:     "2348960443421",
		Contacted: true,
	},
}

func getCustomers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
}

func getCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//convert string to int
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	fmt.Println("id ", id)
	if err != nil {
		//executes if there is any error
		w.WriteHeader(http.StatusInternalServerError)
	} else {

		//loop through customers and find a matching id

		for _, customer := range customers {
			if customer.Id == id {
				//customer exists
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(customer)
			}
		}

		//customer does not exist
		w.WriteHeader(http.StatusNotFound)
	}
}

func deleteElement(slice []Customer, index int) []Customer {
	return append(slice[:index], slice[index+1:]...)
}

func deleteCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//convert string to int
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	fmt.Println("id ", id)
	if err != nil {
		//executes if there is any error
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		//loop through customers and find a matching id

		for i, customer := range customers {
			if customer.Id == id {
				//customer exists
				customers = deleteElement(customers, i)
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(customers)
			}
		}

		//customer does not exist
		w.WriteHeader(http.StatusNotFound)
	}
}

func updateCustomer(w http.ResponseWriter, r *http.Request) {
	// Set the appropriate Content-Type in the response header
	w.Header().Set("Content-Type", "application/json")

	// Create (but not yet assign values to) for the new entry
	var newCustomer Customer // customer == nil

	// Read the HTTP request body
	reqBody, _ := ioutil.ReadAll(r.Body)
	// Encode the request body into a Golang value so that we can work with the data
	json.Unmarshal(reqBody, &newCustomer)

	//check if id of new customer exists
	fmt.Println("parsed id ", newCustomer)

	for index, customer := range customers {
		if customer.Id == newCustomer.Id {
			//customer exists, update customer
			customer.Name = newCustomer.Name
			customer.Role = newCustomer.Role
			customer.Email = newCustomer.Email
			customer.Phone = newCustomer.Phone
			customer.Contacted = newCustomer.Contacted

			customers[index] = customer
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(customer)
		}
	}

	//customer does not exist
	w.WriteHeader(http.StatusNotFound)
}

func addCustomer(w http.ResponseWriter, r *http.Request) {
	// Set the appropriate Content-Type in the response header
	w.Header().Set("Content-Type", "application/json")

	// Create (but not yet assign values to) for the new entry
	var newCustomer Customer // customer == nil

	// Read the HTTP request body
	reqBody, _ := ioutil.ReadAll(r.Body)
	// Encode the request body into a Golang value so that we can work with the data
	json.Unmarshal(reqBody, &newCustomer)

	//check if id of new customer exists
	fmt.Println("parsed id ", newCustomer)
	var customerExist = customerExist(newCustomer.Id)

	if !customerExist {
		customers = append(customers, newCustomer)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(customers)
	} else {
		w.WriteHeader(http.StatusConflict)
	}
}

func customerExist(id int) bool {
	for _, customer := range customers {
		if customer.Id == id {
			return true
		}
	}
	return false
}

func main() {
	// Instantiate a new router
	router := mux.NewRouter()
	// Register handler functions to the same path -- but with different methods
	// E.g., only a GET request to /dictionary can invoke the "getDictionary" handler function
	router.HandleFunc("/customers", getCustomers).Methods("GET")
	router.HandleFunc("/customers/{id}", getCustomer).Methods("GET")
	router.HandleFunc("/customers/{id}", deleteCustomer).Methods("DELETE")
	router.HandleFunc("/customers", addCustomer).Methods("POST")
	router.HandleFunc("/customers", updateCustomer).Methods("PUT")

	//  router.HandleFunc("/dictionary", create).Methods("POST")

	fmt.Println("Server is starting on port 3000...")
	// Pass the customer router into ListenAndServe
	http.ListenAndServe(":3000", router)

}
