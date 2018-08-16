package main

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Region_Structure struct {
	Cur_region_code string
	Regions         []Region
}
type Region struct {
	Region_code    string
	Available_Zone []string
}

var RS_ptr *Region_Structure

func yamlDecoderRS(RSyamlpath string) (*Region_Structure, error) {
	rs_out := new(Region_Structure)
	yaml_file, yaml_open_err := ioutil.ReadFile(RSyamlpath)
	if yaml_open_err != nil {
		log.Println("error while opening yaml", yaml_open_err.Error())
		return rs_out, yaml_open_err
	}
	yaml_decode_err := yaml.Unmarshal(yaml_file, &rs_out)
	if yaml_decode_err != nil {
		log.Println("error while unmarshal on yaml", yaml_decode_err.Error())
		return rs_out, yaml_decode_err
	}
	RS_print(rs_out)
	return rs_out, nil
}
func RS_print(RS_in *Region_Structure) {
	log.Println("PXY$ Region ELB :")
	for _, cur_region := range RS_in.Regions {
		fmt.Println("  ", cur_region.Region_code)
		for _, cur_AZ := range cur_region.Available_Zone {
			fmt.Println("    ", cur_AZ)
		}
	}
	fmt.Println("cur region code : ", RS_in.Cur_region_code)
}
func region2AZlist(region_code_in string) ([]string, error) { //select one of region
	var AZ_list []string
	//requested_service := uri_head(uri_input)
	available_zone_flag := false
	for _, cur_region := range RS_ptr.Regions {
		if region_code_in == cur_region.Region_code {
			AZ_list = cur_region.Available_Zone
			available_zone_flag = true
			break
		}
	}
	if available_zone_flag == false {
		return nil,
			errors.New(RESOURCE_NONEXISTENT_ERROR + ERROR_STRING_SEPARATOR +
				"No region proxy exists for particular locale")
	} else if len(AZ_list) == 0 {
		return nil,
			errors.New(NO_AVAILABLE_INSTANCE_ERROR + ERROR_STRING_SEPARATOR +
				"No Available Zone instance for particular locale")
	}
	//log.Println(inst_list)
	return AZ_list, nil
}
