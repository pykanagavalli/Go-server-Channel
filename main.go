package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type RequestBody struct {
	Ev     string `json:"ev"`
	Et     string `json:"et"`
	ID     string `json:"id"`
	Uid    string `json:"uid"`
	Mid    string `json:"mid"`
	T      string `json:"t"`
	P      string `json:"p"`
	L      string `json:"l"`
	Sc     string `json:"sc"`
	Atrk1  string `json:"atrk1"`
	Atrv1  string `json:"atrv1"`
	Atrt1  string `json:"atrt1"`
	Atrk2  string `json:"atrk2"`
	Atrv2  string `json:"atrv2"`
	Atrt2  string `json:"atrt2"`
	Uatrk1 string `json:"uatrk1"`
	Uatrv1 string `json:"uatrv1"`
	Uatrt1 string `json:"uatrt1"`
	Uatrk2 string `json:"uatrk2"`
	Uatrv2 string `json:"uatrv2"`
	Uatrt2 string `json:"uatrt2"`
	Uatrk3 string `json:"uatrk3"`
	Uatrv3 string `json:"uatrv3"`
	Uatrt3 string `json:"uatrt3"`
}

type ConvertedFormat struct {
	Event       string                `json:"event"`
	EventType   string                `json:"event_type"`
	AppID       string                `json:"app_id"`
	UserID      string                `json:"user_id"`
	MessageID   string                `json:"message_id"`
	PageTitle   string                `json:"page_title"`
	PageURL     string                `json:"page_url"`
	BrowserLang string                `json:"browser_language"`
	ScreenSize  string                `json:"screen_size"`
	Attributes  map[string]Attribute  `json:"attributes"`
	Traits      map[string]TraitValue `json:"traits"`
}

type Attribute struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}
type TraitValue struct {
	Value string `json:"value"`
	Type  string `json:"Type"`
}

func main() {

	mychan := make(chan RequestBody)

	go worker(mychan)

	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return

		}
		var req RequestBody
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Failed to decode JSON request body", http.StatusBadRequest)
			return
		}

		mychan <- req

		fmt.Printf("Received request :%v\n", req)
		successMsg := "Request Receives successfully"
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/,json")
		json.NewEncoder(w).Encode(map[string]string{"message": successMsg})

	})
	fmt.Println("server listening on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func worker(requestChannel <-chan RequestBody) {
	for {
		select {
		case reqBody := <-requestChannel:
			// Convert the request into the desired format
			converted := ConvertToConvertedFormat(reqBody)

			// Do something with the converted format, for now, just print it
			fmt.Printf("Converted format: %+v\n", converted)
		}
	}
}

func ConvertToConvertedFormat(req RequestBody) ConvertedFormat {
	return ConvertedFormat{
		Event:       req.Ev,
		EventType:   req.Et,
		AppID:       req.ID,
		UserID:      req.Uid,
		MessageID:   req.Mid,
		PageTitle:   req.T,
		PageURL:     req.P,
		BrowserLang: req.L,
		ScreenSize:  req.Sc,
		Attributes: map[string]Attribute{
			req.Atrk1: {Value: req.Atrv1, Type: req.Atrt1},
			req.Atrk2: {Value: req.Atrv2, Type: req.Atrt2},
		},
		Traits: map[string]TraitValue{
			req.Uatrk1: {Value: req.Uatrv1, Type: req.Uatrt1},
			req.Uatrk2: {Value: req.Uatrv2, Type: req.Uatrt2},
			req.Uatrk3: {Value: req.Uatrv3, Type: req.Uatrt3},
		},
	}
}
