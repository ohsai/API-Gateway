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

var Region_HC_ptr *Region_HealthChecker

type Region_HealthChecker struct {
	Name             string
	Available_Zone   []string
	index_roundrobin int
}

func healthchecker_init() {
	Region_HC_ptr = &Region_HealthChecker{Name: "Region_HealthChecker", Available_Zone: nil, index_roundrobin: 0}
	go active_check(10)
}
func active_check(interval_sec int) error {
	for {
		var refresh_available_instance []string = nil
		for _, cur_instance := range Region_ptr.Available_Zone {
			err := ping(cur_instance)
			//if no respond, make it unhealthy
			if err != nil {
				log.Println("port scan on ", cur_instance, " finished with error : ", err)
				//os.Exit(3)
			} else {
				refresh_available_instance = append(refresh_available_instance, cur_instance)
			}
		}
		Region_HC_ptr.Available_Zone = refresh_available_instance
		//Region_HC_print(Region_HC_ptr)
		time.Sleep(time.Duration(interval_sec) * time.Second)
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
func Region_HC_print(RHC_in *Region_HealthChecker) {
	log.Println("ELB$Structure of ", RHC_in.Name, " :")
	for _, cur_instance := range RHC_in.Available_Zone {
		fmt.Println("    ", cur_instance)
	}
}
