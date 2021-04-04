package main

import (
	"fmt"
	"math/rand"
	"time"
)

// type package struct {
// 	id int
// }

// type node struct {
// 	id int
// 	receive chan package
// 	send chan package
// 	connections []node
// }

// add function to node struct to choose random receiver of package from connections slice

func genGraph(n, d int) map[int][]int {
	s1 := rand.NewSource(time.Now().UnixNano())
    r1 := rand.New(s1)
	mapa := make(map[int][]int)

	// generating basic structure of directed graph
	for i := 0; i < n-1; i++ {
		mapa[i] = []int{i+1}
	}

	// generating random shortcuts
	for i := 0; i < d; i++ {
		var receiver int;
		var sender int
		
		for {
			sender = r1.Intn(n-1) // last one cannot be a sender so n-1
		 	receiver = r1.Intn(n-sender-1)+sender+1
			 if(!inSlice(mapa[sender], receiver)) {
				 break
			 }
		}
		mapa[sender] = append(mapa[sender], receiver)
	}
	return mapa
}

// check if slice contain val
func inSlice(slice []int, val int) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

// displaying given directed graph in nice format
func printGraph(graph map[int][]int) {
	for edgeNum, receivers := range graph {
		fmt.Print(edgeNum, " :  ")
		for _, receiver := range receivers {
			fmt.Print(receiver,",")
		}
		fmt.Print("\n")
	}
}

func main() {


	graph := genGraph(5,3)
	fmt.Println(graph)
	printGraph(graph)
}