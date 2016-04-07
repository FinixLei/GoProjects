package main

import (
    _ "github.com/cityweather/docs"
    _ "github.com/cityweather/routers"
    _ "github.com/mattn/go-sqlite3"
    
    "time"
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/orm"
)

func init() {
    orm.RegisterDataBase("default", "sqlite3", "./weather.db")
    orm.RunSyncdb("default", false, true)
    orm.DefaultTimeLoc = time.UTC
}

func main() {
    beego.Run()
}
