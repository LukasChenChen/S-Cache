// Copyright 2022 Google Inc. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"net/http"
    "crypto/tls"

)


const schedulerName = "my-scheduler"

func startSchedule() {
    doneChan := make(chan struct{})
	var wg sync.WaitGroup

	wg.Add(1)
	go monitorUnscheduledPods(doneChan, &wg)
	log.Printf("monitorUnscheduledPods finish")


	wg.Add(1)
	go reconcileUnscheduledPods(30, doneChan, &wg)
	log.Printf("reconcileUnscheduledPods finish")



	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case <-signalChan:
			log.Printf("Shutdown signal received, exiting...")
			close(doneChan)
			wg.Wait()
			os.Exit(0)
		}
	}
}

func main() {
	file, err := os.OpenFile("./custom.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	log := log.New(file, "Custom Log", log.LstdFlags)

	if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

	LoadConfig("./config/my-config.json")

	log.Printf("Starting custom scheduler...")
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	doneChan := make(chan struct{})
	var wg sync.WaitGroup

	wg.Add(1)
	go monitorUnscheduledPods(doneChan, &wg)
	log.Printf("monitorUnscheduledPods finish")


	wg.Add(1)
	go reconcileUnscheduledPods(50, doneChan, &wg)
	log.Printf("reconcileUnscheduledPods finish")

	
	scheduleRequests()

	


	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case <-signalChan:
			log.Printf("Shutdown signal received, exiting...")
			close(doneChan)
			wg.Wait()
			os.Exit(0)
		}
	}
    
	file.Close()
}


func main_test() {
	// LoadRequest("./config/prediction.csv", 10000)

	// genZipf(10, 100, 2)
    LoadConfig("./config/my-config.json")
	loadTopo(config_G.TopoName)
	initGlobal()
	genReqZipf("./config/prediction.csv", 10000, "./config/prediction.csv")
	

	// for k, v := range requestsMap_G.Map { 
	// 	// log.Printf("time slot [%d] \n", k)
		
	// 		// log.Println(len(v.Requests))
		
	// }

	log.Println("map length", len(requestsMap_G.Map))
	log.Println(requestCount_G)
	getRequestsMapSum()
	// log.Println(requestsMap_G.Map)
}

// type testvalue struct {
//     value int
// }

// func (tv *testvalue) add(value int){
// 	tv.value = value
// }

// type testmap struct {
// 	Map map[int]testvalue  // <time_slot, requests>
// }

// func (tm *testmap) init(){
// 	tm.Map = make(map[int]testvalue)
// }

// //check empty case
// func (tm *testmap) add(v int, key int){ 
	
// 	tv, ok := tm.Map[key]


// 	if ok {
		
// 		(&tv).add(v)
// 		tm.Map[key] = tv
		

// 	}else{
        
// 		tv := testvalue{}
// 		(&tv).add(v)
// 		tm.Map[key] = tv
		
// 	}
    
// }

// func main() {

// 	tm := testmap{}
// 	tm.init()
// 	tm.add(999,1)
//     log.Printf("%d", tm.Map[1].value)
// 	tm.add(200,1)
// 	log.Printf("%d", tm.Map[1].value)

// }