package main

// Deals with health checker

import (
	"fmt"
	"log"
	"net/url"
	"os/exec"
	"strings"
	"time"
)

var HealthChecker_ptr *HealthChecker

type Service_HealthChecker struct {
	Name               string
	Available_instance []string
	index_roundrobin   int
}

type HealthChecker struct {
	Name    string
	Service []Service_HealthChecker
}

func healthchecker_init() {
	HealthChecker_ptr = &HealthChecker{Name: "HealthChecker", Service: make([]Service_HealthChecker, 0)}
	for _, cur_service := range MSA_ptr.Service {
		HealthChecker_ptr.Service = append(HealthChecker_ptr.Service,
			Service_HealthChecker{Name: cur_service.Name,
				Available_instance: nil,
				index_roundrobin:   0})
	}
	go active_check(Config_ptr.MSA_healthcheck_interval)
}
func active_check(interval_sec int) error {
	for {
		for i, cur_service := range MSA_ptr.Service {
			//refresh available instance list
			var refresh_available_instance []string = nil
			for _, cur_instance := range cur_service.Instance {
				err := ping(cur_instance)
				//if no respond, make it unhealthy
				if err != nil {
					log.Println("Unable to connect ", cur_instance, " due to error : ", err)
				} else {
					refresh_available_instance = append(refresh_available_instance, cur_instance)
				}
			}
			HealthChecker_ptr.Service[i].Available_instance = refresh_available_instance
		}
		if Config_ptr.MSA_print_or_not {
			HealthChecker_print(HealthChecker_ptr)
		}
		time.Sleep(time.Duration(interval_sec) * time.Second)
	}
	return nil
}

// Netcat specific port
func ping(cur_instance string) error {
	cur_instance_url, err := url.ParseRequestURI(cur_instance)
	if err != nil {
		return err
	}
	cur_instance_host_port := strings.Split(cur_instance_url.Host, ":")
	cmd := exec.Command("nc", "-zv", cur_instance_host_port[0], cur_instance_host_port[1])
	err = cmd.Run()
	return err
}

//Pretty printer
func HealthChecker_print(HC_in *HealthChecker) {
	fmt.Println("PXY$Status [", HC_in.Name, "] :")
	for _, cur_service := range HC_in.Service {
		fmt.Println("  ", cur_service.Name)
		for _, cur_instance := range cur_service.Available_instance {
			fmt.Println("    ", cur_instance)
		}
	}
}
