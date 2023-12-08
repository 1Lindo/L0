package main

import (
	"L0/config"
	"L0/models"
	"L0/repo"
	"L0/service/event"
	"L0/service/provider"
	context "context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	log "log"
	"net/http"
	"time"
)

const DATA = `{
  "order_uid": "b563feb7b2b84b6test",
  "track_number": "WBILMTESTTRACK",
  "entry": "WBIL",
  "delivery": {
    "name": "Test Testov",
    "phone": "+9720000000",
    "zip": "2639809",
    "city": "Kiryat Mozkin",
    "address": "Ploshad Mira 15",
    "region": "Kraiot",
    "email": "test@gmail.com"
  },
  "payment": {
    "transaction": "b563feb7b2b84b6test",
    "request_id": "",
    "currency": "USD",
    "provider": "wbpay",
    "amount": 1817,
    "payment_dt": 1637907727,
    "bank": "alpha",
    "delivery_cost": 1500,
    "goods_total": 317,
    "custom_fee": 0
  },
  "items": [
    {
      "chrt_id": 9934930,
      "track_number": "WBILMTESTTRACK",
      "price": 453,
      "rid": "ab4219087a764ae0btest",
      "name": "Mascaras",
      "sale": 30,
      "size": "0",
      "total_price": 317,
      "nm_id": 2389212,
      "brand": "Vivienne Sabo",
      "status": 202
    }
  ],
  "locale": "en",
  "internal_signature": "",
  "customer_id": "test",
  "delivery_service": "meest",
  "shardkey": "9",
  "sm_id": 99,
  "date_created": "2021-11-26T06:22:19Z",
  "oof_shard": "1"
}`

var (
	ev    event.IEvent
	cfg   *config.Config
	db    repo.IDBRepo
	cache repo.ICacheRepo
	ctx   = context.Background()
)

func newRouter(repoDB repo.IDBRepo) (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc("/orders/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		order, err := repoDB.GetOrderByID(ctx, id)
		if err != nil {
			log.Println(err)
			w.Write([]byte("error in service"))
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode(order)
	}).Methods("GET", "OPTIONS")

	return
}

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	cfg = config.GetConfig()
	fmt.Println(ctx)

	// db connect
	dbConn, err := provider.ConnectDB(cfg.DataBase.DbType + cfg.DataBase.DbUser + cfg.DataBase.DbUserPassword + cfg.DataBase.DbPort + cfg.DataBase.DbName)
	if err != nil {
		log.Fatalln(err)
	}
	defer dbConn.Close()
	db = repo.NewDBRepo(dbConn)
	fmt.Printf("Data Base connection established to %v%v%v%v%v \n", cfg.DataBase.DbType, cfg.DataBase.DbUser, cfg.DataBase.DbUserPassword, cfg.DataBase.DbPort, cfg.DataBase.DbName)

	//cache connect
	cacheConn, err := provider.ConnectCache(cfg.Cache.Addr)
	if err != nil {
		log.Fatalln(err, "cacheConn err")
	}
	log.Println(cacheConn, config.GetConfig())
	defer cacheConn.Close()
	status, err := cacheConn.Ping(ctx).Result()
	if err != nil {
		log.Fatalln("Redis connection was refused")
	}
	fmt.Println(status)
	cache = repo.NewCacheRepo(cacheConn, db)
	cache.GetOrders(ctx)
	//nats connect and chan initialisation
	ev = event.NewEvent()
	if err := ev.Connect(cfg.Nats.URL); err != nil {
		log.Fatalln(err)
	}
	defer ev.Close()
	fmt.Printf("NATS connection established to %v at ", cfg.Nats.URL)
	ch := make(chan *models.EventOrder)
	if err := ev.Subscribe("updates", ch); err != nil {
		log.Fatalln(err)
	}

	go func(ctx context.Context) {
		for {
			select {
			case <-time.After(time.Second * 1):
				var order models.EventOrder
				json.Unmarshal([]byte(DATA), &order)
				err := ev.Publish("updates", &order)
				if err != nil {
					return
				}
			case orderData := <-ch:
				if err := cache.InsertOrder(ctx, orderData); err != nil {
					log.Fatalln(err)
				}
			}
		}
	}(ctx)

	router := newRouter(db)
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
