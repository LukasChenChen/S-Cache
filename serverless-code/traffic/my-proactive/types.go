// Copyright 2016 Google Inc. All Rights Reserved.
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

import(
	"errors"
	"sort"
	"strconv"
	"log"
)

// Event is a report of an event somewhere in the cluster.
type Event struct {
	ApiVersion     string          `json:"apiVersion,omitempty"`
	Count          int64           `json:"count,omitempty"`
	FirstTimestamp string          `json:"firstTimestamp"`
	LastTimestamp  string          `json:"lastTimestamp"`
	InvolvedObject ObjectReference `json:"involvedObject"`
	Kind           string          `json:"kind,omitempty"`
	Message        string          `json:"message,omitempty"`
	Metadata       Metadata        `json:"metadata"`
	Reason         string          `json:"reason,omitempty"`
	Source         EventSource     `json:"source,omitempty"`
	Type           string          `json:"type,omitempty"`
}

// EventSource contains information for an event.
type EventSource struct {
	Component string `json:"component,omitempty"`
	Host      string `json:"host,omitempty"`
}

// ObjectReference contains enough information to let you inspect or modify
// the referred object.
type ObjectReference struct {
	ApiVersion string `json:"apiVersion,omitempty"`
	Kind       string `json:"kind,omitempty"`
	Name       string `json:"name,omitempty"`
	Namespace  string `json:"namespace,omitempty"`
	Uid        string `json:"uid"`
}

// PodList is a list of Pods.
type PodList struct {
	ApiVersion string       `json:"apiVersion"`
	Kind       string       `json:"kind"`
	Metadata   ListMetadata `json:"metadata"`
	Items      []Pod        `json:"items"`
}

type PodWatchEvent struct {
	Type   string `json:"type"`
	Object Pod    `json:"object"`
}

type Pod struct {
	Kind     string   `json:"kind,omitempty"`
	Metadata Metadata `json:"metadata"`
	Spec     PodSpec  `json:"spec"`
}

type PodSpec struct {
	NodeName   string      `json:"nodeName"`
	Containers []Container `json:"containers"`
	SchedulerName string   `json:"schedulerName"`
}

type Container struct {
	Name      string               `json:"name"`
	Resources ResourceRequirements `json:"resources"`
	Image     string               `json:"image"`
	Ports     []Port               `json:"ports"`
	Env       []SingleEnv          `json:"env"`
}

type SingleEnv struct{
	Name  string   `json:"name"`
	Value string   `json:"value"`
}

type Port struct{
	ContainerPort  int   `json:"containerPort"`
}



type ResourceRequirements struct {
	Limits   ResourceList `json:"limits"`
	Requests ResourceList `json:"requests"`
}

type ResourceList map[string]string

type Binding struct {
	ApiVersion string   `json:"apiVersion"`
	Kind       string   `json:"kind"`
	Target     Target   `json:"target"`
	Metadata   Metadata `json:"metadata"`
}

type Target struct {
	ApiVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Name       string `json:"name"`
}

type NodeList struct {
	ApiVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Items      []Node
}

type Node struct {
	Metadata Metadata   `json:"metadata"`
	Status   NodeStatus `json:"status"`
}

type NodeStatus struct {
	Capacity    ResourceList `json:"capacity"`
	Allocatable ResourceList `json:"allocatable"`
}

type ListMetadata struct {
	ResourceVersion string `json:"resourceVersion"`
}

type Metadata struct {
	Name            string            `json:"name"`
	GenerateName    string            `json:"generateName"`
	ResourceVersion string            `json:"resourceVersion"`
	Labels          map[string]string `json:"labels"`
	Annotations     map[string]string `json:"annotations"`
	Uid             string            `json:"uid"`
}


type Request struct {
	ID           int
	Function     Function
    Ingress      PhyNode
    ArriveTime   int 
	Served       bool //succ, failed
	DeployNode   PhyNode
	LatencyPara  float64
	IsColdStart  bool
	IsNotPredicted bool
}

func (r *Request) update(Function Function, deployNode PhyNode, IsColdStart bool){
    r.Function = Function
    r.Served = true
	r.DeployNode = deployNode
	r.IsColdStart = IsColdStart
}

type PodQueue struct {
	Items      []Pod 
}

//define Function name as hello_word_1,2...n
//each function is a container in the knative
type Function struct {
	ApiVersion 	   string
	kind       	   string  //service
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
	PhyNode        PhyNode
}

//init the function, name is "type-count_G, such as 1-166"
func (f *Function) init(phyNode PhyNode){
	name := "container-" + strconv.Itoa(f.Type) + "-" + strconv.Itoa(count_G) 

	f.Name = name

	f.PhyNode = phyNode
}

func (f *Function) activePriority(){
	freq := functionfreq_G.get(f.Type)
	f.Clock = clock_G
	f.Priority = f.Clock + (float64(freq) * f.ColdStartTime)/(f.Size/1000)
}

//if cache do not use current clock
func (f *Function) cachePriority(){
	freq := functionfreq_G.get(f.Type)
	f.Priority = f.Clock + (float64(freq) * f.ColdStartTime)/(f.Size/1000) //GB unit
}

func (f *Function) showPriority(){
	log.Printf("function %s priority : %f",  f.Name, f.Priority)
}

type Service struct {
	ApiVersion    string          `json:"apiVersion"`
	Kind          string          `json:"kind"`
	Metadata      Metadata        `json:"metadata"`
	Spec          ServiceSpec     `json:"spec"`
}

type ServiceSpec struct {
    Template   Template `json:"template"`
	
}

type Template struct {

	Spec  TempSpec `json:"spec"`
}

type TempSpec struct {

	SchedulerName     string  `json:"schedulerName"`
	Containers []Container `json:"containers"`

}


type DeploymentList struct {
	ApiVersion string      `json:"apiVersion"`
	Items      []Service 
}

type Config struct {
    TopoName        string    `json:"TopoName"`
    RequestFile     string    `json:"RequestFile"`
	SchedulerName   string    `json:"SchedulerName"`
	LatencyPara     float64   `json:"LatencyPara"`
	MemCap          float64   `json:"MemCap"`
	NodeNum         int       `json:"NodeNum"`
	Beta            float64   `json:"Beta"`   //the parameter for zipf distribution
	SlotNum         int       `json:"SlotNum"`  //the time slot to run
	ReduFactor      float64   `json:"ReduFactor"`  //factor to reduce the request num
	Token           string    `json:"Token"`  //factor to reduce the request num
	PredRequestFile string    `json:"PredRequestFile"`
 }

type Topology struct {
	Nodes []PhyNode
}

type PhyNode struct {
	ID     int
	Lat    float64
	Long   float64
	Mem    float64 
}

//caPhyNodeche method
func (p *PhyNode) getMem() float64 {
    return p.Mem
}

type Cache struct {
	PhyNodeID int //the phynodeID
	FunctionList []Function  //prioirtised function list, ascending
}

func (c *Cache) show() {
	log.Printf("phynode id %d cache size is %d ", c.PhyNodeID, len(c.FunctionList))
}

func (c *Cache) showPriority() {
	for i :=0; i < len(c.FunctionList); i++{
		log.Printf("Function %s Priority %f ", c.FunctionList[i].Name, c.FunctionList[i].Priority)
	}
}

//cache method
func (c *Cache) sortList() {
    sort.SliceStable(c.FunctionList, func(i, j int) bool {
		return c.FunctionList[i].Priority < c.FunctionList[j].Priority
	})
}

func (c *Cache) find(functionName string) (int, error) {
	
    for i := 0; i< len(c.FunctionList); i++{
		if functionName == c.FunctionList[i].Name{
			return i,nil
		}
	}
	return 0, errors.New("can't find this function by name")
}

//add a function to the list, then sort it
func (c *Cache)add(f Function){
    functionfreq_G.minus(f.Type)
	f.cachePriority()
    c.FunctionList = append(c.FunctionList, f)
	c.sortList()
	//add to cache means not active any more
	
}

//delete a function from list, if not exist, will not delete, keep the order
func (c *Cache) delete(i int){
	//check if i exist
	if(i >= 0 && i < len(c.FunctionList)){
	    c.FunctionList = append(c.FunctionList[:i],  c.FunctionList[i+1:]...)
	}
	
}



//cache on phynode 0...n
type CacheMap struct {
	Caches map[int]Cache  //phynode id, cache
}


func (cm *CacheMap) init(){
	cm.Caches = make(map[int]Cache)
}

func (cm *CacheMap) show(){
	if cm == nil{
		log.Printf("CacheMap is nil")
		return
	}

	if len(cm.Caches) == 0{
		log.Printf("CacheMap is empty")
		return
	}
    
	log.Printf("show CacheMap....")
	for _, c := range cm.Caches{

		(&c).show()
	}
}

func (cm *CacheMap) showPriority(){
	if cm == nil{
		log.Printf("CacheMap is nil")
		return
	}

	if len(cm.Caches) == 0{
		log.Printf("CacheMap is empty")
		return
	}
    
	log.Printf("CacheMap Priority....")
	for _, c := range cm.Caches{

		(&c).showPriority()
	}
}



//add to cache, memory not changed
func (cm *CacheMap) add(f Function){
    c,ok := cm.Caches[f.PhyNode.ID]
	if ok {
        (&c).add(f)
		cm.Caches[f.PhyNode.ID] = c
	}else{
		//the cache does not exist yet
		c = Cache{}
		c.PhyNodeID = f.PhyNode.ID
		c.FunctionList = []Function{}
		(&c).add(f)
		cm.Caches[f.PhyNode.ID] = c
	}
}

//i is index of idle function in the cache
func (cm *CacheMap) getIdleFunction(phyNodeID int, funcType int)(f Function, i int){

    c := cm.Caches[phyNodeID]

	emptyF := Function{}

	for i := 0; i< len(c.FunctionList); i = i+1{
		if funcType == c.FunctionList[i].Type{
			return c.FunctionList[i], i
		}
	}
    
	//-1 means no such idle function
	return  emptyF, -1
}

func (cm *CacheMap) delete(phyNodeID int, funcIndex int){
	c := cm.Caches[phyNodeID]
    (&c).delete(funcIndex)
	cm.Caches[phyNodeID] = c

}

func (cm *CacheMap) getLowestPriority(phyNodeID int) float64{
	c := cm.Caches[phyNodeID]

	if len(c.FunctionList) == 0{
		return 0
	}else{
		return c.FunctionList[0].Priority
	}
}

//delete the function with lowest priority
func (cm *CacheMap) deleteLowFunction(phyNodeID int) (Function, bool) {
	c := cm.Caches[phyNodeID]

	var f Function
	
	if len((&c).FunctionList) > 0{
		f = (&c).FunctionList[0]
		(&c).delete(0)
		cm.Caches[phyNodeID] = c
		return f, true
	}else{
		return f, false
	}
}

//sort all the caches
func (cm *CacheMap) sort() {
	for i := 0; i < len(cm.Caches); i++{
		c := cm.Caches[i]
		(&c).sortList()
	}
}

//Active slice that record active functions on single node
type Functions struct {
	Type  int 
	Slice []Function
}

func (f *Functions) add(Function Function){
	
	f.Slice = append(f.Slice, Function)
}

func (f *Functions) delete(k int){
	f.Slice = append(f.Slice[:k], f.Slice[k+1:]...)
}

func (f *Functions) showPriority(){
	for i := 0; i < len(f.Slice); i++{
		log.Printf("Function ", f.Slice[i].Name, " Priority ", f.Slice[i].Priority)
	}
}

// active functions of different types
type NodeFunctions struct {
    PhyNodeID  int
    Functions  map[int]Functions  //<funcType, FunctionSlice>
}

func (nf *NodeFunctions) show(nodeID int) {

	for funcType, Functions := range nf.Functions {
        log.Printf("Node id %d funcType: %d => active container num: %d", nodeID, funcType, len(Functions.Slice))
    }
}

func (nf *NodeFunctions) showPriority() {

	for _, Functions := range nf.Functions {
		(&Functions).showPriority()
    }
}

func (nf *NodeFunctions) add(Function Function){
	if nf == nil || len(nf.Functions) == 0{
		nf.Functions = make(map[int]Functions)
		log.Println("NodeFunctions is nil, init now")
	}
    
    f,ok:= nf.Functions[Function.Type]
	if ok {
	   (&f).add(Function)
	   nf.Functions[Function.Type] = f 
	}else{
	   f = Functions{}
	   f.Type = Function.Type
	   (&f).add(Function)
	   nf.Functions[Function.Type] = f 
	}
    
}

func (nf *NodeFunctions) delete(j int, k int){
    
    f := nf.Functions[j]
    (&f).delete(k)
	nf.Functions[j] = f
}

//TODO: update priortiy
//active functions on all nodes
type ActiveFunctions struct {
	Map  map[int]NodeFunctions // < node id, functionsmap>
}

func (af *ActiveFunctions) show(){

	if af == nil{
		log.Printf("ActiveFunctions is nil")
		return
	}

	if len(af.Map) == 0{
		log.Printf("ActiveFunctions is empty")
		return
	}

    log.Printf("show ActiveFunctions....")

	for nodeID, nf := range af.Map{
		(&nf).show(nodeID)
	}
}

func (af *ActiveFunctions) showPriority(){

	if af == nil{
		log.Printf("ActiveFunctions is nil")
		return
	}

	if len(af.Map) == 0{
		log.Printf("ActiveFunctions is empty")
		return
	}

    log.Printf("show priority....")

	for _, nf := range af.Map{
		(&nf).showPriority()
	}
}

func (af *ActiveFunctions) init(){
	af.Map = make(map[int]NodeFunctions)

}

func (af *ActiveFunctions) add( Function Function, phyNodeID int){

	if af == nil || len(af.Map) == 0{
		af.Map = make(map[int]NodeFunctions)
		log.Println("ActiveFunctions is nil, init now")
	}
	
	functionfreq_G.add(Function.Type)
    Function.activePriority()
    nf, ok:= af.Map[phyNodeID]
	if ok {
	    (&nf).add(Function)
		af.Map[phyNodeID] = nf
	}else{
		nf = NodeFunctions{}
		nf.PhyNodeID = phyNodeID
		nf.Functions = make(map[int]Functions)
		(&nf).add(Function)
		af.Map[phyNodeID] = nf
	}
}

func (af *ActiveFunctions) delete(i int, j int, k int){

    nf := af.Map[i]
	(&nf).delete(j, k)
	af.Map[i] = nf
}

type Distance struct {
	PhyNodeID int 
	distance  float64 //in km
}

type DistSlice struct{
	Slice []Distance
}

type FunctionFreq struct {
	Map map[int]int  //<function type, freq>
}


func (ff *FunctionFreq) init(){  
    ff.Map = make(map[int]int)
}


func (ff *FunctionFreq) get( Type int) int{  
    freq := ff.Map[Type]
    return freq
}

func (ff *FunctionFreq) add( Type int){  
	if ff == nil{
		ff.Map = make(map[int]int)
		log.Println("FunctionFreq is nil, init now")
	}


	val, ok := ff.Map[Type]

	if ok {
		ff.Map[Type] = val + 1
	}else{
        ff.Map[Type] = 1
	}
}

func (ff *FunctionFreq) minus( Type int){  

	val, ok := ff.Map[Type]

	if ok {
		ff.Map[Type] = val - 1
	}else{
        //if nil no need to minus
	}
}

type EndPoints struct{

	apiHost             string   `json:"apiHost"`
    bindingsEndpoint    string   `json:"bindingsEndpoint"`
	eventsEndpoint      string   `json:"eventsEndpoint"`
	nodesEndpoint       string   `json:"nodesEndpoint"`
	podsEndpoint        string   `json:"podsEndpoint"`
	watchPodsEndpoint   string   `json:"watchPodsEndpoint"`
	defaultPodsEndPoint string   `json:"defaultPodsEndPoint"`
	deploymentEndpoint  string   `json:"deploymentEndpoint"`
	knativeSvcEndpoint  string   `json:"knativeSvcEndpoint"`
	token               string   `json:"token"`
	trafficGenEndpoint  string   `json:"trafficGenEndpoint"`

}

type RequestFile struct {
    Time []int
    App1 []int
    App2 []int
	App3 []int
	App4 []int
}

type Requests struct {
	Requests []Request // a slice of request to be added
}


func (rs *Requests) add(r Request){  
	
    rs.Requests = append(rs.Requests, r)
	
}

type RequestsMap struct {
	Map map[int]Requests  // <time_slot, requests>
}

func (rm *RequestsMap) init(){
	rm.Map = make(map[int]Requests)
}

//check empty case
func (rm *RequestsMap) add(r Request, timeSlot int){ 

	if rm == nil{
		rm.Map = make(map[int]Requests)
	}
	
	rs, ok := rm.Map[timeSlot]

	if ok {
		//change the rs value will not change the value in rm
		(&rs).add(r)
		//TODO: WHY NEED THIS, can pointer not change?
		rm.Map[timeSlot] = rs
		

	}else{
        
		requests := Requests{}
		requests.Requests = []Request{}
		(&requests).add(r)
		rm.Map[timeSlot] = requests
		
	}
    
}

func (rm *RequestsMap) get(timeSlot int) ([]Request, bool) {
	rs, ok := rm.Map[timeSlot]
	r := []Request{}

	if ok {
		return rs.Requests, true
	}else{
		return r,false
	}

}