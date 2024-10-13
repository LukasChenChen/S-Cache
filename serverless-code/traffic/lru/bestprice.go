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
	// "strconv"
	"errors"
	"log"
)

// func bestPrice(nodes []Node) (Node, error) {
// 	type NodePrice struct {
// 		Node  Node
// 		Price float64
// 	}

// 	var bestNodePrice *NodePrice
// 	for _, n := range nodes {
// 		price, ok := n.Metadata.Annotations["hightower.com/cost"]
// 		if !ok {
// 			continue
// 		}
// 		f, err := strconv.ParseFloat(price, 32)
// 		if err != nil {
// 			return Node{}, err
// 		}
// 		if bestNodePrice == nil {
// 			bestNodePrice = &NodePrice{n, f}
// 			continue
// 		}
// 		if f < bestNodePrice.Price {
// 			bestNodePrice.Node = n
// 			bestNodePrice.Price = f
// 		}
// 	}

// 	if bestNodePrice == nil {
// 		bestNodePrice = &NodePrice{nodes[0], 0}
// 	}
// 	return bestNodePrice.Node, nil
// }

//test version
// func selectNode(nodes []Node) (Node, error){
// 	type NodePrice struct {
// 		Node  Node
// 		Price float64
// 	}

// 	var bestNodePrice *NodePrice

// 	for _, n := range nodes {
// 		index, ok := n.Metadata.Labels["index"]

// 		if !ok {
// 			continue
// 		}
		// i, err := strconv.ParseFloat(index, 32)
		// if err != nil {
		// 	return Node{}, err
		// }
		// if i == 2 {
		// 	bestNodePrice = &NodePrice{n, i}
		// 	break
		// }




// 	return bestNodePrice, nil

// }

//TODO: select a node for the request
//select the best fit node
//func selectNode(request Request,nodes []Node) (Node, error){
//take the field from the service
func selectNode(nodes []Node, pod *Pod) (Node, error){
	log.Println("Starting selectNode......")

	type NodePrice struct {
		Node  Node
		Price float64
	}

	var bestNodePrice *NodePrice
    deployNodeID := "nil"
	for _, n := range nodes {
		index, ok := n.Metadata.Labels["index"]

		if !ok {
			continue
		}

		
		for _, j := range pod.Spec.Containers[0].Env {
			if j.Name == "phynodeid" {
				deployNodeID = j.Value
			}

			if index == deployNodeID {
				bestNodePrice = &NodePrice{n, 0}
				return bestNodePrice.Node, nil
			}
		}

	}

	if bestNodePrice == nil{
		//just for empty pointer, this result will not be used any way
		node := Node{}
		bestNodePrice = &NodePrice{node, 0}
		return bestNodePrice.Node, errors.New("Node resource exausted")
	}

	return bestNodePrice.Node, errors.New("Cannot find assigned node" + deployNodeID)

}
