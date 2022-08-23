package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

type Order struct {
	Id    string
	Title string
	Desc  string
	Price string
}

var ctx = context.Background()

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "HomePage!")
	fmt.Println("homePage")
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/getAllOrder", returnAllOrders)
	myRouter.HandleFunc("/order/{id}", returnSingleOrder)
	log.Fatal(http.ListenAndServe(":8000", myRouter))
}

func returnSingleOrder(w http.ResponseWriter, r *http.Request) {
	/*vars := mux.Vars(r)
	key := vars["id"]

	for _, order := range orders {
		if order.Id == key {
			json.NewEncoder(w).Encode(order)
		}
	}*/
}

func (u *Order) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, u); err != nil {
		return err
	}
	return nil
}

func returnAllOrders(w http.ResponseWriter, r *http.Request) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	err := client.Ping(context.Background()).Err()

	if err != nil {
		// Sleep for 3 seconds and wait for Redis to initialize
		time.Sleep(3 * time.Second)
		err := client.Ping(context.Background()).Err()
		if err != nil {
			panic(err)
		}
	}

	subscriber := client.Subscribe(ctx, "send-user-data")
	channel := subscriber.Channel()

	for msg := range channel {
		u := &Order{}
		// Unmarshal the data into the user
		err := u.UnmarshalBinary([]byte(msg.Payload))
		if err != nil {
			panic(err)
		}

		fmt.Println(u)
	}

	//order := Order{}

	//for {
	//msg, err := subscriber.ReceiveMessage(ctx)
	//fmt.Println(msg)
	//if err != nil {
	//	panic(err)
	//}
	//
	//if err := json.Unmarshal([]byte(msg.Payload), &order); err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println("Received message from " + msg.Channel + " channel.")
	//fmt.Printf("%+v\n", order)
	//}
}

func main() {
	handleRequests()
}
