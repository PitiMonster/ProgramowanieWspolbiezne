package main

import (
    "fmt"
    "math/rand"
    "sync"
    "time"
	"flag"
)

var wg sync.WaitGroup // 1
var nodesNum int
var addEdgesNum int
var graph map[int][]Node


type Bundle struct {
	id int
	nodes []Node
}

type Observer struct {
	message chan string
	quit chan bool 
}

func (o *Observer) Observe() {
	for {
		select {
		case msg := <-o.message:
			fmt.Println(msg)
		case <- o.quit:
			wg.Done() 
			return 
		}
	}
}

type Receiver struct {
	v                int
	in               chan Bundle
	quit             chan bool
	observer         chan string
	packages         []Bundle
	maxPackagesCount int
}

type Node struct {
	v        int
	in       chan Bundle
	observer chan string
	bundles []Bundle
}

func (rcv *Receiver) Receive() {
	for {
		p, _ := <-rcv.in
		rcv.observer <- fmt.Sprintf("Pakiet %d został odebrany", p.id)
		rcv.packages = append(rcv.packages, p)
		if len(rcv.packages) == rcv.maxPackagesCount {
			rcv.quit <- false
			return
		}
	}
}

func genGraph(n, d int) []Node {
	s1 := rand.NewSource(time.Now().UnixNano())
    r1 := rand.New(s1)
	graph = make(map[int][]Node)
	nodes := make([]Node, 0)

	// generating n nodes
	for i := 0; i < n; i++ {
		node := Node{v: i}
		chan_in := make(chan Bundle)
		node.in = chan_in
		nodes = append(nodes, node)
	}

	// generating basic structure of directed graph
	for i := 0; i < n-1; i++ {
		graph[i] = []Node{nodes[i+1]}
	}

	// generating random shortcuts
	for i := 0; i < d; i++ {
		var receiver int;
		var sender int
		
		for {
			sender = r1.Intn(n-1) // last one cannot be a sender so n-1
		 	receiver = r1.Intn(n-sender-1)+sender+1
			 if(!inSlice(graph[sender], nodes[receiver])) {
				 break
			 }
		}
		graph[sender] = append(graph[sender], nodes[receiver])
	}
	return nodes
}

// check if slice contain val
func inSlice(slice []Node, val Node) bool {
	for _, item := range slice {
		if item.v == val.v {
			return true
		}
	}
	return false
}

// displaying given directed graph in nice format
func printGraph(graph map[int][]Node) {
	for edgeNum, receivers := range graph {
		fmt.Print(edgeNum, " :  ")
		for _, receiver := range receivers {
			fmt.Print(receiver.v,",")
		}
		fmt.Print("\n")
	}
}

func (n *Node) Listen() {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	for {
		time.Sleep(time.Duration(r.Float64() * float64(time.Second))) // czekamy losową ilość czasu 
		select {
		case p := <-n.in:
			
			n.observer <- fmt.Sprintf("pakiet %d jest na wierzchołku %d", p.id, n.v)
			p.nodes = append(p.nodes, *n) // add current node to bundle's visited nodes
			n.bundles = append(n.bundles, p) // add current bundle to node's handled bundles
			time.Sleep(time.Duration(r.Float64() * float64(time.Second)))
				receiverId := r.Intn(len(graph[n.v]))
				receiver := graph[n.v][receiverId]
				receiver.in <- p
		default:
			time.Sleep(time.Duration(r.Float64() * float64(time.Second)))
		}
	}
}


func main() {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	// create flags to get initial values from command line
	nodesNumPtr := flag.Int("nodes", 5, "graph's nodes amount")
	addEdgesNumPtr := flag.Int("edges", 2, "graph's additional edges")
	packagesCounterPtr := flag.Int("bundles", 10, "amount of bundles to send")

	flag.Parse()

	packagesCounter := *packagesCounterPtr
	nodesNum = *nodesNumPtr
	addEdgesNum = *addEdgesNumPtr
    nodes := genGraph(nodesNum, addEdgesNum)

	fmt.Printf("Drukowanie grafu:\n\n")
    printGraph(graph)

	fmt.Printf("\n\nDzialanie programu:\n\n")

	q := make(chan bool)
	
	// crating observer and his in_channel
	o := make(chan string)
	obs := Observer{o, q}

	go obs.Observe()

	wg.Add(1) // waiting for observer to finish his work

	// creating receiver, his channel and his node equivalent
	rec_chan := make(chan Bundle)
	rec := Receiver{v: 3, in: rec_chan, maxPackagesCount: packagesCounter, observer: o, quit: q} // create reciever
	nodeReceiver := Node{v: nodesNum, in: rec.in} // receiver of node type 

	graph[nodesNum-1] = append(graph[nodesNum-1], nodeReceiver)
	
	go rec.Receive()


	// run nodes
	for i := 0; i < nodesNum; i++ {
		nodes[i].observer = obs.message
		go nodes[i].Listen()
	}

	// send packages
	for i := 1; i <= packagesCounter; i++ {
		nodes[0].in <- Bundle{id: i}
		time.Sleep(time.Duration(r.Float64() * float64(time.Second)))

	}

	wg.Wait()

	fmt.Printf("\n\n\nZestawienie danych każdego pakietu\n\n\n")

	for _, v := range rec.packages {
		fmt.Printf("Wierchołki pakietu %d:\n", v.id)
		
		for k1, v1 := range v.nodes {
			fmt.Printf("Wierzchołek numer %d to wierzchołek o id %d\n", k1, v1.v)
		}
		fmt.Printf("\n\n\n")
	}

	fmt.Printf("\n\n\nZestawienie danych każdego wierzchołka\n\n\n")

	for _, v := range nodes {
		fmt.Printf("Pakiety wierzchołka %d:\n", v.v)
		for k1, v1 := range v.bundles {
			fmt.Printf("Pakiet numer %d to pakiet o id %d\n", k1, v1.id)
		}
		fmt.Printf("\n\n\n")
	}

    fmt.Println("Main finished")


}