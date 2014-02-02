package play

import (
	"encoding/json"
	"fmt"
	"github.com/NOX73/go-neural"
	"github.com/NOX73/go-neural/lern"
	"github.com/NOX73/go-neural/persist"
	"io/ioutil"
	"log"
)

const (
	jsonFile   = "json/lang.json"
	speed      = 0.1
	sampleFile = "json/sample/"
)

var sampleNames = []string{"plus", "minus", "multiple", "divide"}

type Sample struct {
	In  []float64
	Out []float64
}

func DiscernmentMain() {
	//createNetwork()
	n := loadNetwork()

	//testNetwork(n)
	//lernNetwork(n)
	testNetwork(n)

	saveNetwork(n)
}

func testNetwork(n *neural.Network) {
	log.Println("--------------------------------------")
	//sampleNames := []string{"plus", "plus2", "minus", "multiple", "multiple2", "divide"}
	sampleNames := []string{"minus3"}

	log.Println("Sample \t Evaluation \t Result")
	for _, sampleName := range sampleNames {
		sample := loadSample(sampleName)
		res := n.Calculate(sample.In)
		e := lern.Evaluation(n, sample.In, sample.Out)
		log.Printf("%s \t %.3f \t\t %v", sampleName, e, res)
	}
}
func lernNetwork(n *neural.Network) {
	samples := make([]*Sample, 0, 10)

	samples = append(samples, loadSample("plus"))
	samples = append(samples, loadSample("minus"))
	samples = append(samples, loadSample("multiple"))
	samples = append(samples, loadSample("divide"))

	for i := 0; i < 10000; i++ {

		for _, s := range samples {
			lern.Lern(n, s.In, s.Out, speed)
		}

	}

}

func loadSample(name string) *Sample {

	s := &Sample{}

	fileName := fmt.Sprint(sampleFile, name, ".json")
	b, _ := ioutil.ReadFile(fileName)
	json.Unmarshal([]byte(b), s)

	return s
}

func loadNetwork() *neural.Network {
	return persist.FromFile(jsonFile)
}

func saveNetwork(n *neural.Network) {
	persist.ToFile(jsonFile, n)
}

func createNetwork() {

	n := neural.NewNetwork(9, []int{9, 9, 4})
	n.RandomizeSynapses()

	persist.ToFile(jsonFile, n)
}
