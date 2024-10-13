package main

import(
    "log"
    "encoding/json"
	"net/http"
    "net/url"
    "bytes"
    "io/ioutil"
    "errors"
    "time"
)


var activeFunctions_G ActiveFunctions //add will auto add freq by 1

//the cache that stores all the cached function
var cacheMap_G CacheMap //add will auto remove freq by 1

var config_G Config

var functionfreq_G FunctionFreq  //add when push to  activeFunctions_G, minus when push to cacheMap_G

var timeInterval int //time slot 1 ....1440

var topo_G Topology

var clock_G float64 //update at the begining of each round

var count_G int // a count for the total number of containers, also used to name new containers

var endPoints_G EndPoints

var requestCount_G int

var requestsMap_G RequestsMap

// log.SetLevel(log.DebugLevel)
//add deploynode to request

func initGlobal(){
    count_G = 0
    //always init the maps first.
    functionfreq_G.init()
    activeFunctions_G.init()
    cacheMap_G.init()
    requestsMap_G.init()
    clock_G = 1
    requestCount_G = 0
}

//assign a request to the current phynode with cached container
func placeToCurrent(requestPtr *Request, Function Function, i int){ 

    // log.Debug( "start placeToCurrent....")
    //put it in the active list
    activeFunctions_G.add(Function, requestPtr.Ingress.ID)

	//delete it from the cached list
    cacheMap_G.delete(requestPtr.Ingress.ID, i)

	requestPtr.update(Function, requestPtr.Ingress, false)


}

//update topo memory
func updateTopo(operator string, phyNodeID int, mem float64){
    if operator == "add"{
        topo_G.Nodes[phyNodeID].Mem = topo_G.Nodes[phyNodeID].Mem + mem
    }else if operator == "minus" {
        topo_G.Nodes[phyNodeID].Mem = topo_G.Nodes[phyNodeID].Mem - mem
    }
}

//create new containers
func createToCurrent(requestPtr *Request, i int){
    // log.Debug( "start createToCurrent....")

    requestPtr.Function.activePriority()//when use this request, update the priority first

    succFlag := true

    var f Function
    
    //if memory not sufficient and this request prioirty is higher than the cached one
    if requestPtr.Function.Size > topo_G.Nodes[requestPtr.Ingress.ID].Mem{
        //clear the cachesmap
        for{

                p2 := cacheMap_G.getLowestPriority(requestPtr.Ingress.ID)

                if p2 < requestPtr.Function.Priority{

                    f, succFlag = cacheMap_G.deleteLowFunction(requestPtr.Ingress.ID)//delete the lower priority function
                    if succFlag == true{

                        //If success, terminate the container
                        termContainers(f)

                        updateTopo("add", requestPtr.Ingress.ID, f.Size)

                    }
                    if requestPtr.Function.Size <= topo_G.Nodes[requestPtr.Ingress.ID].Mem{
                        // call knative to create new containers
                        createContainers(requestPtr)

                        //put it in the active list
                        activeFunctions_G.add(requestPtr.Function, requestPtr.Ingress.ID)

                        // functionfreq_G.add(requestPtr.Function.Type)
                        updateTopo("minus", requestPtr.Ingress.ID, requestPtr.Function.Size)

                        // requestPtr.update(requestPtr.Function, requestPtr.Ingress, true)

                        return
                    }
                    
                    //this means no more space to delete
                    if succFlag == false {
                        return
                    }
                }else{
                    return
                }
           }

    }else{
        createContainers(requestPtr)

        //put it in the active list
        activeFunctions_G.add(requestPtr.Function, requestPtr.Ingress.ID)

        // functionfreq_G.add(requestPtr.Function.Type)
        updateTopo("minus", requestPtr.Ingress.ID, requestPtr.Function.Size)

        // requestPtr.update(requestPtr.Function, requestPtr.Ingress, true)

        return
    }
}

//deploy to a neighbour cached function
func placeToNeighbour(requestPtr *Request, Function Function, i int, phyNodeID int){
    // log.Debug( "start placeToNeighbour....")

    //put it in the active list
    activeFunctions_G.add(Function, phyNodeID)

	//delete it from the cached list
    cacheMap_G.delete(phyNodeID, i)
    
    deployNode, err := getPhyNode(phyNodeID)

    if err != nil{
        return
    }

	requestPtr.update(Function, deployNode, false)
}


func  deployToNeighbour(ds DistSlice, requestPtr *Request) (bool){
    //try assigin to neighbour cached function otherwise create new one on current node.
    
    succFlag := false

    //ignore the first element, it is the current node
    for j := 1; j < len(ds.Slice); j= j+1{
        if (ds.Slice[j].distance * config_G.LatencyPara < requestPtr.Function.ColdStartTime){
            f, i := cacheMap_G.getIdleFunction(ds.Slice[j].PhyNodeID, requestPtr.Function.Type)
         
            if (i == -1){
                //no such function on the neighbour node, find next node
                continue
            }else{
                //place to the local node
                placeToNeighbour(requestPtr, f, i, ds.Slice[j].PhyNodeID)

                succFlag = true
            }
        }

    }

    return succFlag

}


//remove from active list, remove from cache
func createContainers(requestPtr *Request){

    // log.Debug("Starting createContainers()......")

    count_G = count_G + 1

    requestPtr.Function.init(requestPtr.Ingress) // create new request only one case, on the current node

	service, err := createService(requestPtr.Function)

    requestPtr.update(requestPtr.Function, requestPtr.Ingress, true)

	if err != nil {
		log.Printf("fail to createContainers")
	}

    err = postService(service)
    log.Println("containercount %d", count_G)
}


//remove from active list, remove from cache
func termContainers(f Function) error{
    // log.Debug("Starting termContainers......")

    err := deleteService(f)

	if err != nil {
		log.Printf("fail to termContainers")
		return err
	}

	return nil
}

//one interval finished, update priority, add to cacheMap or terminate it
func updateCache(){
    // log.Debug("Starting updateCache......")

    //at the end of each interval, move active containers to cached containers
    //get all the function
    
    for phyNodeID, _ := range activeFunctions_G.Map{
        for funcType, _ := range activeFunctions_G.Map[phyNodeID].Functions{
            for k := 0; k < len(activeFunctions_G.Map[phyNodeID].Functions[funcType].Slice); k++{
                  f := activeFunctions_G.Map[phyNodeID].Functions[funcType].Slice[k]
                  cacheMap_G.add(f)
                  activeFunctions_G.delete(phyNodeID, funcType, k)
            }
        }
    }
    
    cacheMap_G.sort()

    // clear the active function list
    activeFunctions_G = ActiveFunctions{}

}

//deploy one request
func deployRequest(requestPtr *Request){
    
    //sort the physical nodes based on their distance to the request
    ds := sortPhyNodes(requestPtr)
    
    // log.Debug("Starting getIdleFunction......")

    f, i := cacheMap_G.getIdleFunction(requestPtr.Ingress.ID, requestPtr.Function.Type)

    if i == -1{
        //no such function on the current node, find a neighbour node
        if deployToNeighbour(ds, requestPtr) == false{
            createToCurrent(requestPtr, i)
        }

    }else{
        //place to the local node cache functuin
        placeToCurrent(requestPtr, f, i)
    }

}

//deploy request of one time interval
func deployRequests(requests []Request){

    for i := 0; i < len(requests); i++{

        requestPtr := &(requests[i])

        deployRequest(requestPtr)
    }
}

//request finished, send it using json
func sendResult(r Request) error{
    log.Printf("sendResult......")
    var b []byte
	body := bytes.NewBuffer(b)
	err := json.NewEncoder(body).Encode(r)
	if err != nil {
		return err
	}

    request := &http.Request{
		Body:          ioutil.NopCloser(body),
		ContentLength: int64(body.Len()),
		Header:        make(http.Header),
		Method:        http.MethodPost,
		URL: &url.URL{
			Host:   trafficGenEndpoint,
			// Path:   trafficGenEndpoint,
			Scheme: "http",	
            //port
		},
	}
	request.Header.Set("Content-Type", "application/json")

    resp, err := http.DefaultClient.Do(request)

    if err != nil {
		return err
	}
	if resp.StatusCode != 200 {

        log.Printf("sendResult: Unexpected HTTP status code", resp.Status)
		return errors.New("sendResult: Unexpected HTTP status code" + resp.Status)
	}

    // log.Debug("sendResult: finish....")
	return nil
}

func sendResults(requests []Request) {
    for i := 0; i < len(requests); i++{
        sendResult(requests[i])
    }
}

//deploy all request
func scheduleRequests(){



    LoadConfig("./config/my-config.json")
	loadTopo(config_G.TopoName)
	initGlobal()
    //the factor to reduce the request amolunt is 10000
	genReqZipf(config_G.RequestFile, config_G.ReduFactor, config_G.PredRequestFile)

    // r1, err := genRequest(1, 1, 1, 10)

    // r2, err := genRequest(2, 2, 1, 10)

    // r3, err := genRequest(3, 3, 2, 10)

    // r4, err := genRequest(4, 4, 2, 10)

	// if err != nil{

	// }

    showPhyNodesMem()
    
    //set up a message broker, send result message to it using json,
    //modify the ingress controller
    // for timeInterval := 1; timeInterval < 2; timeInterval++{
        // requests := loadRequest(timeInterval)
        //  var requests []Request

        // requests = append(requests, r1)
        // requests = append(requests, r2)
        // requests = append(requests, r3)
        // requests = append(requests, r4)
        //TODO: iterate over request_G FOR ONE SLOT FIRST
    for i := 0; i < config_G.SlotNum; i++{
        requests, found := requestsMap_G.get(i)

        if found == false{
            log.Printf("cannot find time slot %d", i)
            break
        }
    
        deployRequests(requests)

        // startSchedule()

        
        //invoke sending http request
        sendResults(requests)

        updateCache()

        log.Println("result node id", requests[0].Function.PhyNode.ID)

        //wait for the container up
        time.Sleep(60 * time.Second)
        log.Printf("wait 60 seconds .....")
    }

    getRequestsMapSum()
       
    // }
    
    showPhyNodesMem()

    activeFunctions_G.show()
    cacheMap_G.show()

    activeFunctions_G.showPriority()
    cacheMap_G.showPriority()

    log.Printf("All algorithm finished............")


}