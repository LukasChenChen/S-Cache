package main

import(
	"encoding/csv"
	"fmt"
    "os"
	"io"
	"log"
	"strconv"
    "math"
    "errors"
    "sort"
	// "reflect"
    "encoding/json"
    "io/ioutil"
    // "github.com/openacid/low/mathext/zipf"
)



//Load the end points of k8s
func LoadEndPoints(fileName string){
    jsonFile, err := os.Open(fileName)
    // if we os.Open returns an error then handle it
    if err != nil {
        log.Println(err)
    }

    log.Printf("Successfully Opened EndPoints file....")
    // defer the closing of our jsonFile so that we can parse it later on
    defer jsonFile.Close()

    byteValue, _ := ioutil.ReadAll(jsonFile)

    var result EndPoints
    json.Unmarshal(byteValue, &result)
    endPoints_G = result

    jsonFile.Close()

}

func LoadConfig(fileName string){
    // Open our jsonFile
    jsonFile, err := os.Open(fileName)
    // if we os.Open returns an error then handle it
    if err != nil {
        log.Printf("cannot find config file " + fileName)
        log.Println(err)
        
    }
    log.Printf("Successfully Opened config file")
    // defer the closing of our jsonFile so that we can parse it later on
    defer jsonFile.Close()

    byteValue, _ := ioutil.ReadAll(jsonFile)

    var result Config
    json.Unmarshal(byteValue, &result)
    config_G = result
    
    log.Println(result.TopoName)
    log.Println(result.RequestFile)
    log.Println(result.SchedulerName)
    log.Println(result.LatencyPara)
    log.Println(result.MemCap)
    log.Println(result.NodeNum)
    log.Println(result.Beta)
    log.Println(result.SlotNum)
    log.Println(result.ReduFactor)

    jsonFile.Close()
}


func LoadRequest(fileName string, factor float64) RequestFile{
    // Load a csv file.
    csvFile, err := os.Open(fileName)
	if err != nil {
        log.Println("cannot find request file " + fileName)
		log.Println(err)
	}
	fmt.Println("Successfully Opened request file")
	defer csvFile.Close()
    
    csvLines, err := csv.NewReader(csvFile).ReadAll()
    if err != nil {
        log.Println(err)
    }    

    data := RequestFile{
    }
    for index, line := range csvLines {
        switch index {
        case 0:
            for _, value := range line {
                value, err := strconv.ParseFloat(value, 64)
                if err != nil {
                    panic(err)
                }
                data.Time = append(data.Time, int(value))
            }
        case 1:
            for _, value := range line {
                value, err := strconv.ParseFloat(value, 64)
                if err != nil {
                    panic(err)
                }
                data.App1 = append(data.App1, int(value/factor))
            }
        case 3:
            for _, value := range line {
                value, err := strconv.ParseFloat(value, 64)
                if err != nil {
                    panic(err)
                }
                data.App2 = append(data.App2, int(value/factor))
            }
        case 5:
            for _, value := range line {
                value, err := strconv.ParseFloat(value, 64)
                if err != nil {
                    panic(err)
                }
                data.App3 = append(data.App3, int(value/factor))
            }
        case 7:
            for _, value := range line {
                value, err := strconv.ParseFloat(value, 64)
                if err != nil {
                    panic(err)
                }
                data.App4 = append(data.App4, int(value/factor))
            }
        }
    }
    
    // fmt.Println(data)
    

    csvFile.Close()

    return data
}

//create a new func
func createFunc(ID int, Type int) (Function, error){

    f := Function{}

    switch Type {
    case 1:
        f := container_1
        return f, nil
    case 2:
        f := container_2
        return f, nil
    case 3:
        f := container_3
        return f, nil
    case 4:
        f := container_4
        return f, nil
    // case 5:
    //     f := container_5
    //     return f, nil
    // case 6:
    //     f := container_6
    //     return f, nil
    // case 7:
    //     f := container_7
    //     return f, nil
    // case 8:
    //     f := container_8
    //     return f, nil
    // case 9:
    //     f := container_9
    //     return f, nil
    // case 10:
    //     f := container_10
    //     return f, nil
    // case 11:
    //     f := container_11
    //     return f, nil
    // case 12:
    //     f := container_12
    //     return f, nil
    // case 13:
    //     f := container_13
    //     return f, nil
    // case 14:
    //     f := container_14
    //     return f, nil
    // case 15:
    //     f := container_15
    //     return f, nil
    // case 16:
    //     f := container_16
    //     return f, nil
    // case 17:
    //     f := container_17
    //     return f, nil
    // case 18:
    //     f := container_18
    //     return f, nil
    // case 19:
    //     f := container_19
    //     return f, nil
    // case 20:
    //     f := container_20
    //     return f, nil
    }

    return f, errors.New("Fail to getFunc....")

}

//generate request from .csv
func genRequest(ID int, FunctionID int, IngressID int, ArriveTime int) (Request, error){

    r := Request{}
    
    f, err := createFunc(ID, FunctionID) 

    if err != nil{
        return r, errors.New("Fail to getFunc....")
    }
    
    ingress, err := getPhyNode(IngressID)

    if err != nil{
        return r, errors.New("Fail to genRequest....")
    }

    r  = Request{ 
        ID: ID,
        Function:  f,
        Ingress: ingress,
        ArriveTime: ArriveTime,
        Served: false,
        LatencyPara: config_G.LatencyPara,
        IsColdStart: true,
    }

    return r, nil
}

//Phynode ID starts from 0
func loadTopo(topoFile string){

    var phynode PhyNode

    f, err := os.Open(topoFile)
    if err != nil {
        log.Println("cannot find topo file " + topoFile)
        log.Fatal(err)
    }
    log.Printf("Successfully Opened topo file")
    // remember to close the file at the end of the program
    defer f.Close()

    // read csv values using csv.Reader
    csvReader := csv.NewReader(f)

	for {
        rec, err := csvReader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Fatal(err)
        }

		if rec[0] == ""{
			continue
		}
		    

		phynode.ID, err = strconv.Atoi(rec[0])
		phynode.Lat, err= strconv.ParseFloat(rec[1], 64)
		phynode.Long, err= strconv.ParseFloat(rec[2], 64)
        phynode.Mem = config_G.MemCap
		topo_G.Nodes = append(topo_G.Nodes, phynode)

    }
    
	// log.Printf("load physical nodes ........\n")
	// for i := 0; i < len(topo_G.Nodes); i++ {
	// 	log.Printf("%d\n", topo_G.Nodes[i].ID)
	// 	log.Printf("%f\n", topo_G.Nodes[i].Lat)
	// 	log.Printf("%f\n", topo_G.Nodes[i].Long)
	// }
}

func getPhyNode(IngressID int) (PhyNode, error){

    var node PhyNode
    for i := 0; i < len(topo_G.Nodes); i++ {
        if topo_G.Nodes[i].ID == IngressID{
            return topo_G.Nodes[i], nil
        }
    }
    log.Printf("Physical Node does not exist........")
    return node, errors.New("Non-exist phyNode id")
}

func distance(PhyNode1 PhyNode, PhyNode2 PhyNode,unit ...string) float64 {	
    
    // log.Printf( "start distance()....")
    lat1 := PhyNode1.Lat
    lng1 := PhyNode1.Long
    lat2 := PhyNode2.Lat
    lng2 := PhyNode2.Long

    //if same coordinate, dist is 0.
    if lat1 == lat2 && lng1 == lng2{
        dist := float64(0)
        return dist
    }

	radlat1 := float64(math.Pi * lat1 / 180)
	radlat2 := float64(math.Pi * lat2 / 180)
	
	theta := float64(lng1 - lng2)
	radtheta := float64(math.Pi * theta / 180)
	
	dist := math.Sin(radlat1) * math.Sin(radlat2) + math.Cos(radlat1) * math.Cos(radlat2) * math.Cos(radtheta);
	if dist > 1 {
		dist = 1
	}
	
	dist = math.Acos(dist)
	dist = dist * 180 / math.Pi
	dist = dist * 60 * 1.1515
	
	if len(unit) > 0 {
		if unit[0] == "K" {
			dist = dist * 1.609344
		} else if unit[0] == "N" {
			dist = dist * 0.8684
		}
	}
    
	// log.Printf( "distance is .... %f", dist)
	return dist
}

//sort distance in ascending order
func sortPhyNodes(requestPtr *Request) DistSlice{
    // log.Printf( "start sortPhyNodes.... for ingress %d", requestPtr.Ingress.ID)
    // log.Printf( "topo_G.Nodes length %d", len(topo_G.Nodes))
    var ds DistSlice

    for i := 0; i < len(topo_G.Nodes); i++ {
    //    log.Printf( "node index %d", i)
       dist := distance(requestPtr.Ingress, topo_G.Nodes[i], "K")
    //    log.Printf( "input slice")

       var d Distance
       d.PhyNodeID = i
       d.distance = dist
       ds.Slice = append(ds.Slice, d)
    //    ds.Slice[i].PhyNodeID = i
    //    ds.Slice[i].distance = dist        
    }
    
    sort.Slice(ds.Slice, func(i, j int) bool {
        return ds.Slice[i].distance < ds.Slice[j].distance
    })

    // log.Printf( "sorted distance is %v", ds)

    return ds
}

func showPhyNodesMem(){
    log.Printf( "start to show phynode mem .... ")
    for i :=0; i < len(topo_G.Nodes); i++{
        log.Printf("Node %d remaining mem %f", i, topo_G.Nodes[i].Mem, "MB")
    }

}

func sum(array []int) int {
    sum := 0
    for idx := 0; idx < len(array); idx++ {
        sum += array[idx]
        
    }
    fmt.Printf(" sum: %d\n", sum)
    return sum
}

func genZipf(nodeNum int, maxNum int, beta float64) []int{
    // y = C * x**(1-s),  x ∈ [a, b)
    nodeNum = nodeNum+1
    a := float64(1)
    b := float64(nodeNum) // node number 1 - 10
    if beta == 1{
        beta = beta + 0.01
    }
    s := float64(beta)

    z := New(a, b, s)
    sampleCnt := float64(maxNum) //max value of our function

    sample := make([]int, nodeNum) // sample index 0 to 10, but 0 will never happen
    for u := float64(0); u < 1; u += 1 / sampleCnt {
        x := int(z.Float64(u))
        if x < 0 || x > nodeNum{
            continue
        }
        //what if non-exist
        if len(sample) > x{
            sample[x]++
        }
    }
    sample = sample[1:]

    //guarantee 10 elements in our case
    if len(sample) > nodeNum -1 {
        sample = sample[0:(nodeNum - 1)]
    }
    // fmt.Println("genZipf:", sample)
    // sum(sample)
    return sample
}

func genRequests ( sample []int, funcID int, timeSlot int) {
    //i should be 0 - 9, which correspond to the node id
    // log.Printf("sample size %d \n", len(sample))
    for i :=0; i < len(sample); i++{
        // log.Printf("sample value %d \n", sample[i])
        for j := 0; j < sample[i]; j++{
            r, _ := genRequest(requestCount_G, funcID, i, timeSlot)
            requestsMap_G.add(r, timeSlot)
            requestCount_G += 1
            // getRequestsMapSum()
        }
    }
    // log.Printf("add times %d", add_count_G)
    // log.Printf("append times %d", append_count_G)
}
// change the num by factor
func genReqZipf(filename string, factor float64){
    rFile := LoadRequest(filename, factor)

    // for i := 0; i < len(rFile.App1); i++{
    for i := 0; i < 1; i++{
        sample := genZipf(config_G.NodeNum, rFile.App1[i], config_G.Beta)
        genRequests(sample, 1, i)
        // log.Printf("sample size [%d] \n", len(sample))
        // getRequestsMapSum()
    }

    for i := 0; i < len(rFile.App2); i++{
        sample := genZipf(config_G.NodeNum, rFile.App2[i], config_G.Beta)
        genRequests(sample, 2, i)
    }

    for i := 0; i < len(rFile.App3); i++{
        sample := genZipf(config_G.NodeNum, rFile.App3[i], config_G.Beta)
        genRequests(sample, 3, i)
    }

    for i := 0; i < len(rFile.App4); i++{
        sample := genZipf(config_G.NodeNum, rFile.App4[i], config_G.Beta)
        genRequests(sample, 4, i)
    }

}

//sum of the request num in the map
func getRequestsMapSum() int{
    sum := 0
    for _, v := range (requestsMap_G.Map) { 
        sum += len(v.Requests)
	}
    log.Printf("total Request [%d] \n", sum)
    return sum
}

// //get top 4 most popular app name
// func getApp(fileName string, factor float64) ([]int, []int, []int, []int){

//     // Load a csv file.
//     csvFile, err := os.Open(fileName)
// 	if err != nil {
//         log.Println("cannot find getApp file " + fileName)
// 		log.Println(err)
// 	}
// 	fmt.Println("Successfully Opened getApp file")
// 	defer csvFile.Close()
    
//     csvLines, err := csv.NewReader(csvFile).ReadAll()
//     if err != nil {
//         log.Println(err)
//     }    
//     app1 := []int{}
//     app2 := []int{}
//     app3 := []int{}
//     app4 := []int{}

//     for index, line := range csvLines {
//         switch index {
//         case 1:
//             for i, value := range line {
//                 if ( i == 0 ){
//                     continue
//                 }
//                 if (i > 31){
//                     break
//                 }
//                 value, err := strconv.ParseFloat(value, 64)
//                 if err != nil {
//                     panic(err)
//                 }
//                 app1 = append(app1, int(value/factor))
//             }
//         case 2:
//             for i, value := range line {
//                 if ( i == 0 ){
//                     continue
//                 }
//                 if (i > 31){
//                     break
//                 }

//                 value, err := strconv.ParseFloat(value, 64)
//                 if err != nil {
//                     panic(err)
//                 }
//                 app2 = append(app2, int(value/factor))
//             }
//         case 3:
//             for i, value := range line {
//                 if ( i == 0 ){
//                     continue
//                 }
//                 if (i > 31){
//                     break
//                 }

//                 value, err := strconv.ParseFloat(value, 64)
//                 if err != nil {
//                     panic(err)
//                 }
//                 app3 = append(app3, int(value/factor))
//             }
//         case 4:
//             for i, value := range line {
//                 if ( i == 0 ){
//                     continue
//                 }
//                 if (i > 31){
//                     break
//                 }

//                 value, err := strconv.ParseFloat(value, 64)
//                 if err != nil {
//                     panic(err)
//                 }
//                 app4 = append(app4, int(value/factor))
//             }
//         }
//     }

//     return app1, app2, app3, app4
// }
// //the histogramLen is 30 minutes in our case
// func genHistogram() {
//     histogramLen := 30
//     //pre warm window is 0-5%
//     //keep alive is 5%-99%
    
//     histogram1 := make(map[int]int)
//     histogram2 := make(map[int]int)
//     histogram3 := make(map[int]int)
//     histogram4 := make(map[int]int)


//     app1, app2, app3, app4 := getApp("./config/histogram.csv", 10000)
//     sum1, sum2, sum3, sum4 := 0, 0, 0, 0

//     for i := 0; i < len(app1); i++{
//         if (i > histogramLen){
//             break
//         }
//         histogram1[i] = app1[i]
//         sum1 += sum1
//     }

//     threshold1head := 0.05* float64(sum1)
//     threshold1tail := 0.99*float64(sum1)
    
//     time1head, time1tail, tmp1 := 0, 0, 0
//     for i := 0; i < histogramLen; i++{
        
//         tmp1 += histogram1[i]
//         if float64(tmp1) > threshold1head{
//             time1head = i 
//         }

//         if float64(tmp1) > threshold1tail{
//             time1tail = i-1 
//         }
//     }

//     for i := 0; i < len(app2); i++{
//         if (i > histogramLen){
//             break
//         }
//         histogram2[i] = app2[i]

//         sum2 += sum2
//     }

//     threshold2head := 0.05* float64(sum1)
//     threshold2tail := 0.99* float64(sum1)
    
//     time2head, time2tail, tmp2 := 0, 0, 0
//     for i := 0; i < histogramLen; i++{
        
//         tmp2 += histogram2[i]
//         if float64(tmp2) > threshold2head{
//             time2head = i 
//         }

//         if float64(tmp2) > threshold2tail{
//             time2tail = i-1 
//         }
//     }

//     for i := 0; i < len(app3); i++{
//         if (i > histogramLen){
//             break
//         }
//         histogram3[i] = app3[i]
//         sum3 += sum3
//     }

//     threshold3head := 0.05* float64(sum1)
//     threshold3tail := 0.99* float64(sum1)
    
//     time3head, time3tail, tmp3 := 0, 0, 0
//     for i := 0; i < histogramLen; i++{
        
//         tmp3 += histogram3[i]
//         if float64(tmp3) > threshold3head{
//             time3head = i 
//         }

//         if float64(tmp3) > threshold3tail{
//             time3tail = i-1 
//         }
//     }

//     for i := 0; i < len(app4); i++{
//         if (i > histogramLen){
//             break
//         }
//         histogram4[i] = app4[i]
//         sum4 += sum4
//     }

//     threshold4head := 0.05* float64(sum1)
//     threshold4tail := 0.99* float64(sum1)
    
//     time4head, time4tail, tmp4 := 0, 0, 0
//     for i := 0; i < histogramLen; i++{
        
//         tmp4 += histogram4[i]
//         if float64(tmp4) > threshold4head{
//             time4head = i 
//         }

//         if float64(tmp4) > threshold4tail{
//             time4tail = i-1 
//         }
//     }

    
//     log.Printf("app1 head %d", time1head, " tail %d", time1tail)
//     log.Printf("app2 head %d", time2head, " tail %d", time2tail)
//     log.Printf("app3 head %d", time3head, " tail %d", time3tail)
//     log.Printf("app4 head %d", time4head, " tail %d", time4tail)
// }