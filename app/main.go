package main

import (
	"log"
	"net/http"

	"github.com/MichaelS11/go-dht"
	"github.com/gin-gonic/gin"
)

type env_data struct {
    Humidity float64 `json:"humidity"`
    Temp float64 `json:"temperature"`
}

var monitor *dht.DHT

func main() {
    // Initialize the DHT object
    var err error
    //monitor, err = initialize_monitor("GPIO19", dht.Celsius)
    monitor, err = initialize_monitor("GPIO4", dht.Celsius)
    if err != nil {
        log.Fatal("Error creating device: ", err)
    }

    // Server info
    addr := ":8080"
    log.Printf("Server listening on %s", addr)

    // Setup gin HTTP
    router := gin.Default()
    router.LoadHTMLGlob("./templates/*.html")
    router.GET("/", env_handler)

    // Start HTTP server
    err = router.RunTLS(addr, "./certs/server.crt", "./certs/server.key")
    if err != nil {
        log.Fatal(err)
    }

}


// Provides environment data
func env_handler(c *gin.Context) {
    // Get Data from DHT22
    var data env_data
    var err error
    data.Humidity, data.Temp, err = monitor.ReadRetry(11)
    
    // Check for read error
    if err != nil {
        c.String(http.StatusInternalServerError, "Error: %s", err.Error())
        return
    }

    // Return the HTML 
    c.HTML(http.StatusOK, "index.html", gin.H{"temp": data.Temp, "humidity": data.Humidity})
    return
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



