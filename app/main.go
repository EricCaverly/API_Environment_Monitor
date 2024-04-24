package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/MichaelS11/go-dht"
)

type env_data struct {
    Humidity float64 `json:"humidity"`
    Temp float64 `json:"temperature"`
}

var monitor *dht.DHT

func main() {
    // Initialize the DHT object
    var err error
    monitor, err = initialize_monitor("GPIO19", dht.Celsius)
    if err != nil {
        log.Fatal("Error creating device: ", err)
    }

    // Add handlers (endpoints)
    http.HandleFunc("/env_data", env_handler)

    // Server info
    addr := ":8080"
    log.Printf("Server listening on %s", addr)

    // Start HTTP server
    err = http.ListenAndServe(addr, nil)
    if err != nil {
        log.Fatal(err)
    }

}


// Provides environment data
func env_handler(w http.ResponseWriter, r *http.Request) {
    // Get Data from DHT22
    var data env_data
    var err error
    data.Humidity, data.Temp, err = monitor.ReadRetry(11)
    
    // Check for read error
    if err != nil {
        http.Error(w, fmt.Sprintf("Error: %e", err), http.StatusInternalServerError)
        return
    }

    // Return data and log the interaction
    jsonData, err := json.Marshal(data)
    fmt.Fprintf(w, "%s", string(jsonData))
    log.Printf("[/env_data] Server responded to client - %s", r.RemoteAddr)
}


// Given some settings, return a dht.DHT object
func initialize_monitor(pin string, unit dht.TemperatureUnit) (*dht.DHT, error) {
    err := dht.HostInit()
    if err != nil {
        return nil, err
    }

    dht, err := dht.NewDHT(pin, unit, "")
    if err != nil {
        return nil, err
    }

    return dht, nil
}



