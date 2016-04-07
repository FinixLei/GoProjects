package main

import (
    "testing"
    "net/http"
)

const sampleCityName string = "Shanghai"


func reportFailure(t *testing.T, respCode int, err error) {
    if respCode != http.StatusOK || err != nil {
        t.Errorf("Test Faield: respCode = %d, err = %v", respCode, err)
    }
}

func Test_DeleteOneCity(t *testing.T) {
    respCode, err := DeleteOneCity(sampleCityName)
    reportFailure(t, respCode, err)
}

func Test_AddOneCity(t *testing.T) {
    respCode, err := AddOneCity(sampleCityName)
    if respCode != http.StatusCreated || err != nil {   // 201
        t.Errorf("Test Failed when adding %s for the first time: respCode = %d, err = %v", sampleCityName, respCode, err)
    }
    
    respCode, err = AddOneCity(sampleCityName)
    if respCode != http.StatusConflict || err != nil {   // 409
        t.Errorf("Test Failed when adding %s for the second time: respCode = %d, err = %v", sampleCityName, respCode, err)
    }
}


func Test_GetAllCities(t *testing.T) {
    allCities, respCode, err := GetAllCities()
    reportFailure(t, respCode, err)
    
    found := false
    for _,v := range(allCities) {
        if v == sampleCityName {
            found = true
            break
        }
    }
    if found == false {
        t.Errorf("Test Faield due to no expected city")
    }
}


func Test_GetOneCityWeather(t *testing.T) {
    result, respCode, err := GetOneCityWeather(sampleCityName)
    reportFailure(t, respCode, err)
    
    if result == nil || result.Weather == nil || len(result.Weather) == 0 {
        t.Errorf("Test Failed: returned result = %v", result)
    }
}
