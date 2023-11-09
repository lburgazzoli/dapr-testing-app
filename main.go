package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/gorilla/mux"
)

var (
	stateStoreName string
	appPort        string
	daprClient     dapr.Client
	once           sync.Once
)

func init() {
	appPort = os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080"
	}

	stateStoreName = os.Getenv("STATESTORE_NAME")
	if stateStoreName == "" {
		stateStoreName = "statestore"
	}
}

type MyValues struct {
	Values []string
}

func writeHandler(w http.ResponseWriter, r *http.Request) {
	value := r.URL.Query().Get("message")
	values, _ := read(r.Context(), stateStoreName, "values")

	if values.Values == nil || len(values.Values) == 0 {
		values.Values = []string{value}
	} else {
		values.Values = append(values.Values, value)
	}

	data, err := json.Marshal(values)
	if err != nil {
		panic(err)
	}

	err = save(r.Context(), stateStoreName, "values", data)
	if err != nil {
		panic(err)
	}

	respondWithJSON(w, http.StatusOK, values)
}

func client() dapr.Client {
	once.Do(func() {
		dc, err := dapr.NewClient()
		if err != nil {
			panic(err)
		}

		daprClient = dc
	})

	return daprClient
}
func readHandler(w http.ResponseWriter, r *http.Request) {
	values, _ := read(r.Context(), stateStoreName, "values")
	respondWithJSON(w, http.StatusOK, values)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/write", writeHandler).Methods("POST")
	r.HandleFunc("/read", readHandler).Methods("GET")

	r.HandleFunc("/health/readiness", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})
	r.HandleFunc("/health/liveness", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":"+appPort, nil))

}

// Platform Provided
func save(ctx context.Context, statestore string, key string, data []byte) error {
	return client().SaveState(ctx, statestore, key, data, nil)
}

func read(ctx context.Context, statestore string, key string) (MyValues, error) {
	result, err := client().GetState(ctx, statestore, key, nil)
	if err != nil {
		return MyValues{}, err
	}
	myValues := MyValues{}
	if result.Value != nil {
		json.Unmarshal(result.Value, &myValues)
	}
	if myValues.Values == nil {
		myValues.Values = make([]string, 0)
	}

	return myValues, nil
}
