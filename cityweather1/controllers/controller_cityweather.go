package controllers

import (
    "encoding/json"
    
    "github.com/cityweather/models"
    "github.com/astaxie/beego"
)

// Operations about cityweather
type CityWeatherController struct {
    beego.Controller
}

// @router /location [post]
// Body: {"name": "SomeCity"}
// Return: "201 Created"　or "409 Conflicted"
func (o *CityWeatherController) Post() {
    var cn models.CityName
    json.Unmarshal(o.Ctx.Input.RequestBody, &cn)
    responseCode := models.AddOneCity(&cn)
    
    o.Ctx.Output.Status = responseCode
}


// @router /location/?:name [get]
func (o *CityWeatherController) Get() {
    name := o.Ctx.Input.Param(":name")
    
    if name != "" {
        cw, err := models.GetOneCity(name)
        if err != nil {
            o.Data["json"] = err.Error()
        } else {
            o.Data["json"] = cw
        }
    } else {    // name is empty, then Get all cities' names
        cities := models.GetAllCities()
        o.Data["json"] = cities
    }
    o.ServeJSON()
}


// @router /location/:name [delete]
// Return: always 200 OK 
func (o *CityWeatherController) Delete() {
    name := o.Ctx.Input.Param(":name")
    models.Delete(name)
    
    o.Data["json"] = "Delete success!"
    o.ServeJSON()
}
