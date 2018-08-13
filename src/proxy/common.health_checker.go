package main

import (
	"fmt"
	"log"
	"net/url"
	//"os"
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
		HealthChecker_ptr.Service = append(HealthChecker_ptr.Service, Service_HealthChecker{Name: cur_service.Name, Available_instance: nil, index_roundrobin: 0})
	}
	go active_check(10)
}
func active_check(interval_sec int) error {
	for {
		for i, cur_service := range MSA_ptr.Service {
			var refresh_available_instance []string = nil
			//refresh available instance list
			//HealthChecker_ptr.Service[i].Available_instance = nil
			for _, cur_instance := range cur_service.Instance {
				err := ping(cur_instance)
				//if no respond, make it unhealthy
				if err != nil {
					log.Println("port scan on ", cur_instance, " finished with error : ", err)
					//os.Exit(3)
				} else {
					//HealthChecker_ptr.Service[i].Available_instance = append(HealthChecker_ptr.Service[i].Available_instance, cur_instance)
					refresh_available_instance = append(refresh_available_instance, cur_instance)
				}
			}
			HealthChecker_ptr.Service[i].Available_instance = refresh_available_instance
		}
		time.Sleep(time.Duration(interval_sec) * time.Second)
		//HealthChecker_print(HealthChecker_ptr)
	}
	return nil
}
func ping(cur_instance string) error {
	cur_instance_url, err := url.ParseRequestURI(cur_instance)
	if err != nil {
		return err
	}
	cur_instance_host_port := strings.Split(cur_instance_url.Host, ":")
	//ping:
	cmd := exec.Command("nc", "-zv", cur_instance_host_port[0], cur_instance_host_port[1])
	err = cmd.Run()
	return err
}
func HealthChecker_print(HC_in *HealthChecker) {
	log.Println("Structure of ", HC_in.Name, " :")
	for _, cur_service := range HC_in.Service {
		fmt.Println("  ", cur_service.Name)
		for _, cur_instance := range cur_service.Available_instance {
			fmt.Println("    ", cur_instance)
		}
	}
}
