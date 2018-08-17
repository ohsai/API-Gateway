package main

import (
	"bufio"
	"bytes"
	"github.com/go-redis/redis"
	"log"
	"net/http"
	"net/http/httputil"
)

var redis_client *redis.Client

//const cache_or_not = true

func redis_init(PORT string) {
	redis_client = redis.NewClient(&redis.Options{
		//Addr:     "localhost:6379",
		Addr: Config_ptr.Redis_address,
		//Password: "",
		Password: Config_ptr.Redis_password,
		//DB:       0,
		DB: Config_ptr.Redis_db,
	})
	pong, err := redis_client.Ping().Result()
	if err == nil {
		log.Println("redis client for port ", PORT, " is ready ", pong) //Needless
	}
}

func routing_filter(req *http.Request) (*http.Response, error) {

	//if request method is get and it exists in redis cache already, bring it to response
	aim := req.Header.Get("Service")

	if req.Method == "GET" && uri_head(aim) != "auth" && Config_ptr.Cache_or_not {
		cached_response, err := redis_client.Get(aim).Result()
		if err == redis.Nil {
			//just forward it
		} else if err != nil {
			//post cache error, but forward it anyway
			log.Println("CacheSearchError$", err.Error())
			//return nil, err
		} else {
			//bring it to response
			r := bufio.NewReader(bytes.NewReader([]byte(cached_response)))
			proxyRes, err := http.ReadResponse(r, nil)
			if err != nil {
				log.Println("CacheValueFormatError$", err.Error())

				//return nil, err
			} else {
				//Pass headers crucial for proxy
				proxyRes.Header.Add("Service", req.Header.Get("Service"))
				return proxyRes, err
			}
		}

	}

	//Else, just forward it
	client := &http.Client{}
	proxyRes, err := client.Do(req)
	if err != nil {
		log.Println("Unable to create forward response! : ", err.Error())
	}
	//and cache it

	if req.Method == "GET" && uri_head(aim) != "auth" && Config_ptr.Cache_or_not {
		cache_in, err := httputil.DumpResponse(proxyRes, Config_ptr.Cache_or_not)
		if err != nil {
			log.Println("CacheDumpError$", err.Error())
			return nil, err
		}
		err = redis_client.Set(aim, string(cache_in), 0).Err()
		if err != nil {
			log.Println("CacheStoreError$", err.Error())
			//failed to cache
			return nil, err
		}
	}

	//Pass headers crucial for proxy
	proxyRes.Header.Add("Service", req.Header.Get("Service"))
	return proxyRes, err
}
