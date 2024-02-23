# The Crypto Service Api built on Go



It is a microservice that collect data from several crypto data providers using its API.

This microservice uses:

* gin-gonic/gin package to start and serve HTTP server
* nhooyr.io/websocket package to manage websocket connection
* go-sql-driver/mysql package to work with mysql database
* lib/pq package to work with postgresql database
* go-redis/redis package for saving the active session and restoring it when the puller is restarted


## Build app

```bash
$ go build -o ccd .
````

## Run app
You should previously export some environment variables:

```bash
export CCDC_DATAPROVIDER=cryptocompare
export CCDC_DATABASEURL=postgres://username:password@127.0.0.1:5432/dbname?sslmode=disable
export CCDC_APIKEY=put you api key here
export CCDC_SESSIONSTORE=redis // or "db", default value is "db"
export REDIS_URL=redis://:redis_password@127.0.0.1:6379/0 // only when "redis" session store selected
```

if you want use **huobi** as data provider export this:
```bash
export CCDC_DATAPROVIDER=huobi
```

If you use **mysql** db, you should export something like this:
```bash
export CCDC_DATABASEURL=mysql://username:password@tcp(localhost:3306)/dbname
``` 

And run application:
```bash
$ ./ccd -debug
```

The default port is 8080, you can test the application in a browser or with curl:

```bash
$ curl 127.0.0.1:8080/v1/service/ping
```

You can choose a different port and run more than one copy of **ccd** on your local host. For example:

```bash
$ ./ccd -port 8081
``` 

You also can specify some setting before run application: 
```bash
$ ./ccd -h
ccd is a microservice that collect data from several crypto data providers cryprocompare using its API.

Usage of ccd:
  -dataprovider string
        use selected data provider ("cryptocompare", "huobi") (default "cryptocompare")
  -debug
        run the program in debug mode
  -h    display help
  -port string
        set specify port (default ":8080")
  -session string
        set session store "db" or "redis" (default "db")  
  -timeout int
        how long to wait for a response from the api server before sending data from the cache (default 1000)
```

List of the implemented endpoints:
* **/healthz** [GET]   _check node status_
* **/v1/collect/add** [GET] _add new worker to collect data for the selected pair_
* **/v1/collect/remove** [GET] _stop and remove worker and collecting data for the selected pair_
* **/v1/collect/status** [GET] _show info about running workers_
* **/v1/collect/update** [GET]  _update pulling interval for the selected pair_
* **/v1/symbols/add** [GET] _add new currency symbol to the db_
* **/v1/symbols/update** [GET]  _update currency symbol in the db_
* **/v1/symbols/remove** [GET] _remove currency symbol in the db_
* **/v1/price** [POST, GET] _get actual (or cached if dataprovider is unavailable) info for the selected pair_
* **/v1/ws** [GET] _websocket connection url, when you connected, try to send request like {"fsym":"BTC","tsym":"USD"}_
* **/v1/ws/subscribe** [POST, GET] _subscribe to collect data for the selected pair_
* **/v1/ws/unsubscribe** [POST, GET] _unsubscribe to stop collect data for the selected pair_
* **/v1/symbols** [POST, PUT, DELETE] _add, update, delete currency symbol_
* **/v1/collect** [POST, PUT, DELETE] _add, update, delete worker to collect data_

Example getting a GET request for getting actual info about selected pair:

```bash
$ curl "http://localhost:8080/v1/price?fsym=ETH&tsym=JPY"
```

Example of sending a POST request to add a new worker:

```bash
$ curl -X POST -H "Content-Type: application/json" -d '{ "fsym": "BTC", "tsym": "USD", "interval": 60}' "http://localhost:8080/v1/collect"
```

Example of sending a GET request to remove worker:

```bash
$ curl "http://localhost:8080/v1/collect/remove?fsym=BTC&tsym=USD&interval=60"
```

Example of sending a GET request to subscribe wss channel:

```bash
$ curl "http://localhost:8080/v1/ws/subscribe?fsym=BTC&tsym=USD"
```

### I still need to build the part where the API can be tested , I'm working on it
