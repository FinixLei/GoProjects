package models

import (
    "fmt"
    "time"
    "net/http"
    "io/ioutil"
    
    "github.com/astaxie/beego/orm"
    "github.com/bitly/go-simplejson"
    
    _ "github.com/mattn/go-sqlite3"
)

const weatherTable string = "city_weather"
const timeoutSet int64 = 3600
const OpenWeatherURL string = "http://api.openweathermap.org/data/2.5/weather"
const AppID string = "f87dfd3af38ed44f157296b7150caacc"

type CityName struct {
    Name    string
}

type CityWeather struct {
    Id          int                         // primary key, auto increment
    Name        string    `orm:"unique;"`   // city name
    Summary     string                      // main in weather
    Description string                      // description in weather
    Icon        string                      // icon in weather
    Wid         int                         // id in weather
    TimeStamp   int64                       // timestamp when updating
}

func init() {
    orm.RegisterModel(new(CityWeather))
}

func AddOneCity(cn *CityName) (responseCode int) {
    cw := new(CityWeather)
    cw.Name = cn.Name
    cw.Wid         = -1
    cw.TimeStamp   = 0
    fmt.Println(cw)
    
    o := orm.NewOrm()
    o.Using("main")
    _, err := o.Insert(cw)
    
    responseCode = 201
    if err != nil {
        if err.Error() == "UNIQUE constraint failed: city_weather.name" {
            responseCode = 409    // conflicted
        } else {
            responseCode = 500    // server error
        }
    }
    
    return responseCode
}

func GetAllCities() []string {
    allCities := []string{}     // dynamic array
    
    o := orm.NewOrm()
    o.Using("main")
    qs := o.QueryTable(weatherTable)
    
    var lists []orm.ParamsList
    num, err := qs.ValuesList(&lists, "name")
    if err == nil {
        fmt.Printf("Result Nums: %d\n", num)
        for _, row := range lists {
            fmt.Println(row[0])
            allCities = append(allCities, row[0].(string))
        }
    }
    
    return allCities
}

func GetOneCity(city string) (cw CityWeather, err error) {
    o := orm.NewOrm()
    o.Using("main")
    qs := o.QueryTable(weatherTable)
    
    err = qs.Filter("name", city).One(&cw)
    if err != nil {
        cw = CityWeather{Id: -1}
        return cw, err
    }
    
    currentTime := time.Now().UTC().UnixNano()
    diffSeconds := (currentTime - cw.TimeStamp) / 1e9
    
    fmt.Printf("Diff seconds = %d\n", diffSeconds)
    
    if diffSeconds > timeoutSet || cw.Wid == -1 {  // Older than one hour or the first get, then need to update database
        client := &http.Client{}
        url := OpenWeatherURL + "?q=" + city + "&APPID=" + AppID
        reqest, err := http.NewRequest("GET", url, nil)
        if err != nil {
            panic(err)
            fmt.Println("Error happened when calling openweather")
        }
        
        response, respErr := client.Do(reqest)
        defer response.Body.Close()
        
        if respErr != nil {
            fmt.Printf("Response Error: %s\n", respErr)
        } else {   // Get Response from openweather!!
            body, err := ioutil.ReadAll(response.Body)
            if err != nil {
                panic(err.Error())
            }
            
            js, err := simplejson.NewJson(body)
            if err != nil {
                panic(err.Error())
            }
            
            weather, ok := js.CheckGet("weather")
            if ok {
                fmt.Println(weather)
                desc, _ := weather.GetIndex(0).Get("description").String()
                icon, _ := weather.GetIndex(0).Get("icon").String()
                id, _   := weather.GetIndex(0).Get("id").Int()
                wtr, _  := weather.GetIndex(0).Get("main").String()
                
                num, err := qs.Filter("name", city).Update(orm.Params{
                    "description": desc,
                    "summary": wtr,     // "main" field
                    "wid": id,
                    "icon": icon, 
                    "time_stamp": currentTime, 
                })
                fmt.Printf("num = %d\n", num)
                if err != nil {
                    fmt.Println(err)
                    panic(err)
                }
                
                err = qs.Filter("name", city).One(&cw)  // get cw after updating
                if err != nil {
                    cw = CityWeather{Id: -1}
                    return cw, err
                }
            }
        }
    } 
    
    return cw, err
}


func Delete(city string) {
    o := orm.NewOrm()
    o.Using("main")
    qs := o.QueryTable(weatherTable)
    
    num, err := qs.Filter("name", city).Delete()
    if err == nil {
        fmt.Println(num)
    } else {
        fmt.Println("Error happened in Delete()......")
    }
}
