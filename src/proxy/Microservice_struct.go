package main

import (
	"log"
	"time"
)

type Microservice struct {
	Name     string
	Instance []string
}

type MSA struct {
	Name    string
	Service []Microservice
}

type Service_HealthChecker struct {
	Name               string
	Available_instance []string
}

type HealthChecker struct {
	Service []Service_HealthChecker
}
