//A scaler read function request from LSTM Prediction,
//push new pods to unscheduled queue.


package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"net/http/httputil"
	"strconv"
	// "go.uber.org/zap"
)

var (
	
)

//create a pod based on name
func createPod(function Function) (Pod, error) {
	log.Printf("Start to create pods in scaler .......")
	var err error
	if function.Name ==""{
		err = errors.New("createPod: empty functionname" )
	}

	var port Port
	
	port.ContainerPort = function.Port

	var ports []Port
	ports = append(ports, port)

	container := Container {
		Name:     function.Name,       
		Image:    function.Image,
		Ports:    ports,
	}

	containers := make([]Container, 1)
	containers[0] = container

	pod := Pod{
		Kind:          "Pod",
		Spec: PodSpec{
			Containers:   containers,
			SchedulerName: function.Scheduler,
			
		},
	}

	return pod, err
}


//post the pod to api server
func postPod(pod Pod) error {
	log.Println("Starting scaler postPod()......")

	podList, err := getPods()
	if err != nil {
		return err
	}

	podList.Items = append(podList.Items, pod)


	var b []byte
	body := bytes.NewBuffer(b)
	err = json.NewEncoder(body).Encode(&podList)
	if err != nil {
		return err
	}

	


	request := &http.Request{
		Body:          ioutil.NopCloser(body),
		ContentLength: int64(body.Len()),
		Header:        make(http.Header),
		Method:        http.MethodPost,
		URL: &url.URL{
			Host:   apiHost,
			Path:   podsEndpoint,
			Scheme: "https",	
		},
	}
	request.Header.Set("Content-Type", "application/json")
	
	bearer := "Bearer " + token
	request.Header.Add("Authorization", bearer)
	resp, err := http.DefaultClient.Do(request)

	if err != nil {
		log.Println("Starting scaler postPod() failed!")
		return err
	}
	if resp.StatusCode != 201 {
		log.Println("postPod: Unexpected HTTP status code" + resp.Status)
		return errors.New("postPod: Unexpected HTTP status code" + resp.Status)
	}
	return nil
}

func createService (function Function) (Service, error) {
	// log.Printf("Start to createService in scaler .......")
	var err error
	if function.Name ==""{
		err = errors.New("createService: empty functionname" )
	}
	var port Port
	
	port.ContainerPort = function.Port

	var ports []Port
	ports = append(ports, port)
    
	var singleEnv SingleEnv
	singleEnv.Name = "phynodeid"
	singleEnv.Value = strconv.Itoa(function.PhyNode.ID)
	var env []SingleEnv
	env = append(env, singleEnv)

	container := Container {
		Name:     function.Name,       
		Image:    function.Image,
		Ports:    ports,
		Env:      env,
	}

	containers := make([]Container, 1)
	containers[0] = container 


	service := Service{
		ApiVersion:     function.ApiVersion,
		Kind:           "Service",
		Metadata:       Metadata{Name: function.Name},
		Spec: ServiceSpec{
				Template{
					TempSpec{
					SchedulerName:  function.Scheduler,
				    Containers: containers,
					},
				},
		},
	}

	return service, err
}

//post the pod to api server
func getService() (*DeploymentList, error ){
	var deploymentList DeploymentList

	request := &http.Request{
		Header:        make(http.Header),
		Method:        http.MethodGet,
		URL: &url.URL{
			Host:   apiHost,
			Path:   knativeSvcEndpoint,
			Scheme: "https",	
		},
	}
	request.Header.Set("Content-Type", "application/json")
	
	bearer := "Bearer " + token
	request.Header.Add("Authorization", bearer)
	resp, err := http.DefaultClient.Do(request)

	if err != nil {
		log.Println("Starting getService failed!")
		return nil, err
	}

	if resp.StatusCode != 200 {
		log.Println("getService: Unexpected HTTP status code" + resp.Status)
		return nil, errors.New("getService: Unexpected HTTP status code" + resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&deploymentList)


	if err != nil {
		log.Println("getService failed!")
		return nil, err
	}

	log.Println("getService succeed!")
	return &deploymentList, nil

}



//post the pod to api server
func postService(service Service) error {

	var b []byte
	body := bytes.NewBuffer(b)
	err := json.NewEncoder(body).Encode(service)
	if err != nil {
		return err
	}

	request := &http.Request{
		Body:          ioutil.NopCloser(body),
		ContentLength: int64(body.Len()),
		Header:        make(http.Header),
		Method:        http.MethodPost,
		URL: &url.URL{
			Host:   apiHost,
			Path:   knativeSvcEndpoint,
			Scheme: "https",	
		},
	}
	request.Header.Set("Content-Type", "application/json")
	
	bearer := "Bearer " + token
	request.Header.Add("Authorization", bearer)
	resp, err := http.DefaultClient.Do(request)

	if err != nil {
		log.Printf("Starting scaler postService() failed!")
		return err
	}

	b, err  = httputil.DumpResponse(resp, true)
	if err != nil {
		log.Println(err)
	}

	// log.Printf(string(b))

	if resp.StatusCode != 201   {
		log.Println("postService: Unexpected HTTP status code" + resp.Status)
		log.Println("Service Name:" + service.Metadata.Name)
		return errors.New("postService: Unexpected HTTP status code" + resp.Status)
	}
	
	// log.Printf("postService: service is posted " + service.Metadata.Name)
	
	return nil

}

//delete a pod based on name
// curl $KUBE_API/apis/apps/v1/namespaces/default/deployments/sleep \
//   --cacert ~/.minikube/ca.crt \
//   --cert ~/.minikube/profiles/cluster1/client.crt \
//   --key ~/.minikube/profiles/cluster1/client.key \
//   -X DELETE
func deleteService(function Function) error {
	log.Println("Starting scaler deleteService() ......")
	serviceName := function.Name

	servicePath := knativeSvcEndpoint + "/" + serviceName

	request := &http.Request{
		Header:        make(http.Header),
		Method:        http.MethodDelete,
		URL: &url.URL{
			Host:   apiHost,
			Path:   servicePath,
			Scheme: "https",	
		},
	}
	request.Header.Set("Content-Type", "application/json")
	
	bearer := "Bearer " + token
	request.Header.Add("Authorization", bearer)
	resp, err := http.DefaultClient.Do(request)

	if err != nil {
		log.Println("Starting scaler deleteService() failed! " + serviceName)
		return err
	}
	if resp.StatusCode != 200 {
		log.Println("deleteService: Unexpected HTTP status code" + resp.Status)
		return errors.New("deleteService: Unexpected HTTP status code" + resp.Status)
	}

	if resp.StatusCode == 200 {
		log.Println("deleteService: service is deleted: " + function.Name)
	}

	// log.Printf("Finish scaler deleteService() ......")
	return nil
}



func post() error{
	log.Printf("Starting scaler post()......")

	//getService()
    //TODO:
	service, err := createService(hello_worlddelete)

	if err != nil {
		log.Printf("fail to create service")
		return err
	}



	err = postService(service)

	deleteService(hello_world)

	if err != nil {
		log.Printf("fail to post service")
		return err
	}

	return nil
}

