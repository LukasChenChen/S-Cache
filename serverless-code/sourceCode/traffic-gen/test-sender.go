package main

import(
    "encoding/json"
	"net/http"
        "net/url"
        "bytes"
        "log"
        "io/ioutil"
        "errors"
)

var
( 
	trafficGenEndpoint = "localhost:2222"
)

type PhyNode struct {
	ID     int
	Lat    float64
	Long   float64
	Mem    float64 
}

type Request struct {
	ID           int
	Function     Function
        Ingress      PhyNode
        ArriveTime   int 
	Served       bool //succ, failed
	DeployNode   PhyNode
        LatencyPara  float64
        IsColdStart    bool
}

type Function struct {
	ApiVersion     string
	kind           string  //service
	Scheduler      string
  	Name	       string   //function name, also service name, unchanged, use this to feed traffic, this one must be unique
	Image          string
	Port           int
	Size           float64
	Type           int      //the type of the container, from 1-20 in our case
	ColdStartTime  float64
	ProcessTime    float64 
	Clock          float64
	Priority       float64
	//PhyNode        PhyNode
}


func sendResult() error{

        node_1 := PhyNode{
                ID:     1,
                Lat:    -37.813938,
                Long:   144.954587,
                Mem:    4,    
        }

        node_2 := PhyNode{
                ID:     5,
                Lat:    -37.817843,
                Long:   144.956397,
                Mem:    5,    
        }

        f := Function{
		ApiVersion:     "serving.knative.dev/v1",
		Scheduler:      "my-scheduler",
		Image:          "gcr.io/knative-samples/helloworld-go",
                Port:           8080,
		Size:           300,     //[50,300]
		Type:           19,
		ColdStartTime:  3.5,   //[1.7,4.5]
		ProcessTime:    2.5,  //[2,6.5]
		Clock:          0,
		Priority:       1,
	}

        // r := Request{
        //         ID:      999,
        //         Function: f,
        //         Ingress: node_1,
        //         DeployNode: node_2,
        //         LatencyPara: 20,
        //         Served: true,
        //         IsColdStart: true,
        // }

        r_2 := Request{
                ID:      1000,
                Function: f,
                Ingress: node_2,
                DeployNode: node_1,
                LatencyPara: 20,
                Served: true,
                IsColdStart: false,
        }
        
      
        // var requests []Request
        // requests = append(requests, r)

        // requests = append(requests, r_2)


        var b []byte
        body := bytes.NewBuffer(b)
        err := json.NewEncoder(body).Encode(r_2)
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

                log.Println("sendResult: Unexpected HTTP status code")
                return errors.New("sendResult: Unexpected HTTP status code" + resp.Status)
        }
        return nil


}

func main(){
        log.Println("sendResult start .....")
        sendResult()

}
