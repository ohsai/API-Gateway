package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	if len(os.Args) != 5 {
		panic("Usage : (program name) <port> <yaml path MSA structure> 	<yaml path Region structure> <yaml path Configuration> 	")
	}
	proxy_init(os.Args[1:2][0], os.Args[2:3][0], os.Args[3:4][0], os.Args[4:5][0])
}
func proxy_init(PORT string, MSAyamlpath string, RSyamlpath string, Configyamlpath string) {
	config_manager_init(MSAyamlpath, RSyamlpath, Configyamlpath)
	healthchecker_init()
	redis_init(PORT)
	createListener(PORT)

}

var Config_ptr *Configuration

func config_manager_init(MSAyamlpath string, RSyamlpath string, Configyamlpath string) error {
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
	Config_ptr, err = yamlDecoderConfig(Configyamlpath)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	go config_manager(MSAyamlpath, RSyamlpath, Configyamlpath)
	return nil
}
func config_manager(MSAyamlpath string, RSyamlpath string, Configyamlpath string) {
	var msamodtime time.Time
	var rsmodtime time.Time
	var configmodtime time.Time
	for {

		msastat, err := os.Stat(MSAyamlpath)
		if err == nil {
			if !(msastat.ModTime().Equal(msamodtime)) {
				MSA_ptr_temp, err2 := yamlDecoder(MSAyamlpath)
				if err2 == nil {
					MSA_ptr = MSA_ptr_temp
					msamodtime = msastat.ModTime()
					log.Println("Successfully modified msa setting")
				} else {
					log.Println("Error while decoding msa yaml")
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
					log.Println("Error while decoding rs yaml")
				}
			}
		}
		configstat, err := os.Stat(Configyamlpath)
		if err == nil {
			if !(configstat.ModTime().Equal(configmodtime)) {
				Config_ptr_temp, err2 := yamlDecoderConfig(Configyamlpath)
				if err2 == nil {
					Config_ptr = Config_ptr_temp
					configmodtime = configstat.ModTime()
					log.Println("Successfully modified config setting")
				} else {
					log.Println("Error while decoding config yaml")
				}
			}
		}
		log.Printf("%+v", Config_ptr)

		time.Sleep(time.Duration(5) * time.Second)
	}
}

func createListener(PORT string) {
	http.Handle("/", new(proxyHandler))
	http.ListenAndServe(":"+PORT, nil)
}

type proxyHandler struct {
	http.Handler
}

func (h *proxyHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	proxyReq, prefilter_err := prefilter(req)
	if prefilter_err != nil {
		filter_error_handler("pre_filter", w, prefilter_err)
		return
	}
	proxyRes, routing_err := routing_filter(proxyReq)
	if routing_err != nil {
		filter_error_handler("routing_filter", w, routing_err)
		return
	}
	post_err := post_filter(proxyRes, w)
	if post_err != nil {
		filter_error_handler("post_filter", w, post_err)
		return
	}
}

//This should comply to src/resource/config.yaml
type Configuration struct {
	cache_or_not   bool
	redis_address  string
	redis_password string `
	redis_db       int `yaml:"redis_db"`
	MSA_healthcheck_interval int `yaml:"MSA_healthcheck_interval"`
	config_check_interval int `yaml:"config_check_interval"`
	load_balancer_policy string `yaml:"load_balancer_policy"`
	RS_print_or_not bool `yaml:"RS_print_or_not"`
	MSA_print_or_not bool `yaml:"MSA_print_or_not"`
	auth_key string `yaml:"auth_key"`
}

func yamlDecoderConfig(Configyamlpath string) (*Configuration, error) {
	config_out := new(Configuration)
	yaml_file, yaml_open_err := ioutil.ReadFile(Configyamlpath)
	if yaml_open_err != nil {
		log.Println("error while opening yaml", yaml_open_err.Error())
		return config_out, yaml_open_err
	}
	yaml_decode_err := yaml.Unmarshal(yaml_file, &config_out)
	if yaml_decode_err != nil {
		log.Println("error while unmarshal on yaml", yaml_decode_err.Error())
		return config_out, yaml_decode_err
	}
	return config_out, nil
}
