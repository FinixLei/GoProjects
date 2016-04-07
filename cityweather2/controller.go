package main

import (
    "io"
    "fmt"
    "net/http"
    "strings"
    "encoding/json"
    
    _ "github.com/mattn/go-sqlite3"
)

// POST /location
func postLocationHandler(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    body := make([]byte, r.ContentLength)
    r.Body.Read(body)
    
    var city CityName
    err := json.Unmarshal(body, &city)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        io.WriteString(w, err.Error())
        return
    }

    status, err := AddOneCity(city.Name)
    
    w.WriteHeader(status)
    if err != nil {
        io.WriteString(w, err.Error())
    }
}


// GET /location
func getAllLocationsHandler(w http.ResponseWriter, r *http.Request) {
    cities, respCode, err := GetAllCities()
    w.WriteHeader(respCode)
    
    if err == nil {
        citiesStr := "["
        for i, city := range cities {
            if i > 0 {
                citiesStr += (", " + city)
            } else {
                citiesStr += city
            }
        }
        citiesStr += "]"
        
        io.WriteString(w, citiesStr)
    } else {
      io.WriteString(w, err.Error())
    }
}


// DELETE /location/{name}
func deleteCityHandler(w http.ResponseWriter, r *http.Request, city string) {
    respCode, err := DeleteOneCity(city)
    
    w.WriteHeader(respCode)
    if err != nil {
        io.WriteString(w, err.Error())
    }
}


// GET /location/{name}
func getCityWeatherHandler(w http.ResponseWriter, r *http.Request, city string) {
    result, respCode, err := GetOneCityWeather(city)
    resp, err := json.Marshal(result)
    
    w.WriteHeader(respCode)
    if err == nil {
        w.Write(resp)
    } else {
        io.WriteString(w, err.Error())
    }
}


func topHandler(w http.ResponseWriter, r *http.Request) {
    items := strings.Split(r.URL.Path, "/")
    
    if (len(items) > 4 || len(items) <=1) {
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprintf(w, "404 Not Found: %s", r.URL.Path)
        return
    }
    
    loc := "location"
    firstPlace := strings.ToLower(items[1])
    
    if firstPlace == loc {
        if (r.Method == http.MethodPost && len(items) == 2) {  // POST /location
            postLocationHandler(w, r)
            
        } else if (r.Method == http.MethodGet && (len(items) == 2 || (len(items) == 3 && items[2] == ""))) {  // GET /location
            getAllLocationsHandler(w, r)
            
        } else if (r.Method == http.MethodGet && (len(items) == 3 || (len(items) == 4 && items[3] == ""))) {  // GET /location/{name}
            getCityWeatherHandler(w, r, items[2])
        
        } else if (r.Method == http.MethodDelete && len(items) == 3) {  // DELETE /location/{name}
            deleteCityHandler(w, r, items[2])
        
        } else {
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprintf(w, "404 Not Found: %s", r.URL.Path)
        }
    } else {
        w.WriteHeader(http.StatusNotFound)
            fmt.Fprintf(w, "404 Not Found: %s", r.URL.Path)
    }
}
