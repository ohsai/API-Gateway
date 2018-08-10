package main

import (
	"log"
	"net/url"
	"os/exec"
	"strings"
	"time"
)

func healthchecker_init() {
	HealthChecker_ptr = &MSA{Name: "HealthChecker", Service: make([]Microservice, 0)}
	for _, cur_service := range MSA_ptr.Service {
		HealthChecker_ptr.Service = append(HealthChecker_ptr.Service, Microservice{Name: cur_service.Name, Instance: nil})
	}
	go active_check(10)
}
func active_check(interval_sec int) error {
	for {
		for i, cur_service := range MSA_ptr.Service {
			//refresh available instance list
			HealthChecker_ptr.Service[i].Instance = nil
			for _, cur_instance := range cur_service.Instance {
				err := ping(cur_instance)
				//if no respond, make it unhealthy
				if err != nil {
					log.Println("port scan on ", cur_instance, " finished with error : ", err)
				} else {
					HealthChecker_ptr.Service[i].Instance = append(HealthChecker_ptr.Service[i].Instance, cur_instance)
				}

			}
		}
		time.Sleep(time.Duration(interval_sec) * time.Second)
		MSA_print(HealthChecker_ptr)
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
