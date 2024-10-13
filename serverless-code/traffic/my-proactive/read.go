package main


import (
    "log"
    "io/ioutil"
	"encoding/json"
)




type RequestSet struct{

}

//TODO: implement log debug

func readConfig(configFile string) error {

	content, err := ioutil.ReadFile(configFile)
    if err != nil {
        log.Fatal("Error when opening file: ", err)
		return err
    }

	var config Config
    err = json.Unmarshal(content, &config)
    if err != nil {
        log.Fatal("Error during Unmarshal(): ", err)
		return err
    }

	return err
	
}

func readRequest(){

}

//read .csv file 
func readTopo(){

}
