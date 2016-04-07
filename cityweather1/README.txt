* 语言 Go  
* 框架 随意  
* 后端数据库 Redis或者SQLite，只需要一种即可  
  
## API 列表  
* POST /location  
* GET /location  
* GET /location/{name}  
* DELETE /location/{name}  
  
### POST /location  
增加支持的城市，如果已存在于数据库，返回409，否则返回201  
  
例子1   
  
POST /location  
{ "name": "Shanghai" }  
  
201 Created  
  
例子2  
  
POST /location  
{ "name": "Shanghai" }  
  
201 Created  
  
POST /location  
{ "name": "Shanghai" }  
  
409 Conflicted  
  
### GET /location 返回数据库中的所有城市  
  
例子3   
  
GET /location  
  
200 OK  
[]  
  
例子4  
  
POST /location  
{ "name": "Shanghai" }  
  
201 Created  
  
POST /location  
{ "name": "Beijing" }  
  
201 Created  
  
GET /location  
  
200 OK  
["Shanghai", "Beijing"]  
  
  
### GET /location/{name} 查询openweathermap.com，返回结果，因为天气数据更新不频繁，可缓存在数据库中，保留1个小时  
不需要考虑查询openweathermap.com返回错误的情况  
  
例子5  
  
GET /location/Shanghai  
  
200 OK  
{  
    "weather": [  
        {  
            "description": "few clouds",  
            "icon": "02d",  
            "id": 801,  
            "main": "Clouds"  
        }  
    ]  
}  
  
### DELETE /location/{name}  
  
例子6  
  
DELETE /location/Shanghai  
  
200 OK  

--------------------------------

curl "api.openweathermap.org/data/2.5/weather?q=Shanghai&APPID=xxxxxxxxxxxxxxxxxxxxxxxxxx"  
  
{"coord":{"lon":121.46,"lat":31.22},"weather":[{"id":801,"main":"Clouds","description":"few clouds","icon":"02d"}],"base":"cmc stations","main":{"temp":286.15,"pressure":1019,"humidity":71,"temp_min":286.15,"temp_max":286.15},"wind":{"speed":7,"deg":140},"clouds":{"all":20},"dt":1458608400,"sys":{"type":1,"id":7452,"message":0.0091,"country":"CN","sunrise":1458597323,"sunset":1458641219},"id":1796236,"name":"Shanghai","cod":200}  
  
## 参数  
* q: 城市名  
* APPID: xxxxxxxxxxxxxxxxxxxxxxxxxx 是预先申请的ID