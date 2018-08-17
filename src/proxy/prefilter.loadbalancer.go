package main

import (
	"errors"
	"math/rand"
	"proxy/mycrypt"
)

func load_balance(instance_list []string, policy string, load_balancer_info []string) (string, error) {
	var instance_chosen string = ""
	var chosen_index int
	for i := 0; i < len(instance_list); i++ { // retry
		//Implement round robin / weighted_round_robin / random
		if policy == "round_robin" { //Assume that service for uri exists
			requested_service := uri_head(load_balancer_info[1])
			for j, cur_service := range HealthChecker_ptr.Service {
				if requested_service == cur_service.Name {
					chosen_index = HealthChecker_ptr.Service[j].index_roundrobin % len(instance_list)
					HealthChecker_ptr.Service[j].index_roundrobin += 1
					break
				}
			}

		} else if policy == "ip_hash" {
			chosen_index = mycrypt.String_modhash(load_balancer_info[0], len(instance_list))

			//need url
		} else if policy == "url_hash" {
			chosen_index = mycrypt.String_modhash(
				load_balancer_info[0]+load_balancer_info[1], len(instance_list))
			//need url
		} else { //random
			chosen_index = rand.Intn(len(instance_list))
		}
		instance_chosen = instance_list[chosen_index]
		//Implement one_more_check, retry
		if err := ping(instance_chosen); err == nil {
			return instance_chosen, nil
		}
	}
	return "", errors.New(NO_AVAILABLE_INSTANCE_ERROR + ERROR_STRING_SEPARATOR +
		"instance failed before HealthChecker check")
}
