package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type Configuration struct { //should comply to config input file
	Cache_or_not             bool   `json:"Cache_or_not"`
	Redis_address            string `json:"Redis_address"`
	Redis_password           string `json:"Redis_password"`
	Redis_db                 int    `json:"Redis_db"`
	MSA_healthcheck_interval int    `json:"MSA_healthcheck_interval"`
	Config_check_interval    int    `json:"Config_check_interval"`
	Load_balancer_policy     string `json:"Load_balancer_policy"`
	RS_print_or_not          bool   `json:"RS_print_or_not"`
	MSA_print_or_not         bool   `json:"MSA_print_or_not"`
	Auth_key                 string `json:"Auth_key"`
}

func jsonDecoderConfig(Configjsonpath string) (*Configuration, error) { //read file
	config_out := new(Configuration)
	json_file, json_open_err := ioutil.ReadFile(Configjsonpath)
	if json_open_err != nil {
		log.Println("error while opening json", json_open_err.Error())
		return config_out, json_open_err
	}
	json_decode_err := json.Unmarshal(json_file, &config_out)
	if json_decode_err != nil {
		log.Println("error while unmarshal on json", json_decode_err.Error())
		return config_out, json_decode_err
	}
	return config_out, nil
}

var Config_ptr *Configuration

func config_manager_init(MSAyamlpath string, RSyamlpath string, Configjsonpath string) error { //initial read
	var err error
	MSA_ptr, err = yamlDecoder(MSAyamlpath)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	RS_ptr, err = yamlDecoderRS(RSyamlpath)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	Config_ptr, err = jsonDecoderConfig(Configjsonpath)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	go config_manager(MSAyamlpath, RSyamlpath, Configjsonpath) //active check
	return nil
}
func config_manager(MSAyamlpath string, RSyamlpath string, Configjsonpath string) {
	msastatinit, _ := os.Stat(MSAyamlpath)
	var msamodtime time.Time = msastatinit.ModTime()
	rsstatinit, _ := os.Stat(RSyamlpath)
	var rsmodtime time.Time = rsstatinit.ModTime()
	configstatinit, _ := os.Stat(Configjsonpath)
	var configmodtime time.Time = configstatinit.ModTime()
	for {
		time.Sleep(time.Duration(Config_ptr.Config_check_interval) * time.Second)
		//check modification
		msastat, err := os.Stat(MSAyamlpath)
		if err == nil {
			if !(msastat.ModTime().Equal(msamodtime)) {
				MSA_ptr_temp, err2 := yamlDecoder(MSAyamlpath)
				if err2 == nil {
					MSA_ptr = MSA_ptr_temp
					msamodtime = msastat.ModTime()
					log.Println("Successfully modified msa setting")
				} else {
					log.Println("Error while decoding msa yaml", err2.Error())
				}

			}
		}
		rsstat, err := os.Stat(RSyamlpath)
		if err == nil {
			if !(rsstat.ModTime().Equal(rsmodtime)) {
				RS_ptr_temp, err2 := yamlDecoderRS(RSyamlpath)
				if err2 == nil {
					RS_ptr = RS_ptr_temp
					rsmodtime = rsstat.ModTime()
					log.Println("Successfully modified rs setting")
				} else {
					log.Println("Error while decoding rs yaml", err2.Error())
				}
			}
		}
		configstat, err := os.Stat(Configjsonpath)
		if err == nil {
			if !(configstat.ModTime().Equal(configmodtime)) {
				Config_ptr_temp, err2 := jsonDecoderConfig(Configjsonpath)
				if err2 == nil {
					Config_ptr = Config_ptr_temp
					configmodtime = configstat.ModTime()
					log.Println("Successfully modified config setting")
				} else {
					log.Println("Error while decoding", err2.Error())
				}
			}
		}
	}
}
