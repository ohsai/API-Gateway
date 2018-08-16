package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Microservice struct {
	Name     string
	Instance []string
}

type MSA struct {
	Name    string
	Service []Microservice
}

var MSA_ptr *MSA

func yamlDecoder(MSAyamlpath string) (*MSA, error) {
	msa_out := new(MSA)
	yaml_file, yaml_open_err := ioutil.ReadFile(MSAyamlpath)
	if yaml_open_err != nil {
		log.Println("error while opening yaml", yaml_open_err.Error())
		return msa_out, yaml_open_err
	}
	yaml_decode_err := yaml.Unmarshal(yaml_file, &msa_out)
	if yaml_decode_err != nil {
		log.Println("error while unmarshal on yaml", yaml_decode_err.Error())
		return msa_out, yaml_decode_err
	}
	//MSA_print(msa_out)
	//log.Printf("---MSA : \n%+v\n", msa_out)
	return msa_out, nil
}
func MSA_print(MSA_in *MSA) {
	log.Println("PXY$Structure of ", MSA_in.Name, " :")
	for _, cur_service := range MSA_in.Service {
		fmt.Println("  ", cur_service.Name)
		for _, cur_instance := range cur_service.Instance {
			fmt.Println("    ", cur_instance)
		}
	}
}
