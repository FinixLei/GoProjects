package routers

import (
    "github.com/cityweather/controllers"
    "github.com/astaxie/beego"
)

func init() {
    beego.Router("/location/?:name", &controllers.CityWeatherController{})
}
