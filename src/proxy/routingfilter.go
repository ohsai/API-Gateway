package main

import (
	"bufio"
	"bytes"
	"github.com/go-redis/redis"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

var redis_client *redis.Client
var cache_duration time.Duration = time.Minute //default duration

func redis_init(PORT string) {
	redis_client = redis.NewClient(&redis.Options{
		Addr:     Config_ptr.Redis_address,
		Password: Config_ptr.Redis_password,
		DB:       Config_ptr.Redis_db,
	})
	_, err := redis_client.Ping().Result()
	if err == nil {
		log.Println("redis client for port ", PORT, " is ready ")
	}
	cache_duration, _ = time.ParseDuration("30s")
}

func routing_filter(req *http.Request) (*http.Response, error) {
	aim := req.Header.Get("Service")
	if req.Method == "GET" && uri_head(aim) != "auth" && Config_ptr.Cache_or_not {
		cached_response, err := redis_client.Get(aim).Result()
		if err == redis.Nil { //cache miss
		} else if err != nil {
			log.Println("CacheSearchError$", err.Error())
		} else { //cache hit
			r := bufio.NewReader(bytes.NewReader([]byte(cached_response)))
			proxyRes, err := http.ReadResponse(r, nil)
			if err != nil {
				log.Println("CacheValueFormatError$", err.Error())
			} else {
				proxyRes.Header.Add("Service", req.Header.Get("Service"))
				return proxyRes, err
			}
		}
	}
	client := &http.Client{}
	proxyRes, err := client.Do(req)
	if err != nil {
		log.Println("Unable to create forward response! : ", err.Error())
	}
	if req.Method == "GET" && uri_head(aim) != "auth" && Config_ptr.Cache_or_not { //cache it
		cache_in, err := httputil.DumpResponse(proxyRes, Config_ptr.Cache_or_not)
		if err != nil {
			log.Println("CacheDumpError$", err.Error())
			return nil, err
		}
		err = redis_client.Set(aim, string(cache_in), cache_duration).Err()
		if err != nil {
			log.Println("CacheStoreError$", err.Error())
			return nil, err
		}
	}
	proxyRes.Header.Add("Service", req.Header.Get("Service"))
	return proxyRes, err
}
