package play

import (
	"bufio"
	"github.com/NOX73/go-neural"
	"github.com/NOX73/go-neural/engine"
	"github.com/NOX73/go-neural/lern"
	"github.com/NOX73/go-neural/persist"
	"log"
	"os"
	"os/signal"
)

func LangMain() {
	n := loadNetwork()

	//createLangNetwork()
	testEngine(n)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

lernLoop:
	for !checkEngine(n) {

		lernEngine(n)

		select {
		case <-c:
			log.Println("Interrupt !")
			break lernLoop
		default:
		}

	}
	testEngine(n)

	saveNetwork(n)
}

var (
	//gosamplePath = "/Users/nox73/kaize/docker/archive/archive_test.go"
	//rbsamplePath = "/Users/nox73/kaize/sfk/lib/sfk/generators/app/templates/controllers/root.rb"
	gosamplePath = "./json/sample/sample.go"
	rbsamplePath = "./json/sample/sample.rb"
)

func checkEngine(n *neural.Network) bool {
	gosample := getSampleFromFile(gosamplePath)
	rbsample := getSampleFromFile(rbsamplePath)

	var out []float64

	out = n.Calculate(gosample)
	if out[0] < 0.8 || out[1] > 0.2 {
		log.Println("CheckFailed", gosamplePath, out)
		return false
	}

	out = n.Calculate(rbsample)
	if out[1] < 0.8 || out[0] > 0.2 {
		log.Println("CheckFailed", rbsamplePath, out)
		return false
	}

	return true
}

func testEngine(n *neural.Network) {
	gosample := getSampleFromFile(gosamplePath)
	rbsample := getSampleFromFile(rbsamplePath)

	log.Println(gosamplePath, n.Calculate(gosample))
	log.Println(rbsamplePath, n.Calculate(rbsample))
}

func lernEngine(n *neural.Network) {

	gofiles := getGoFiles()
	rbfiles := getRbFiles()

	//min := 5
	min := len(gofiles)

	if min > len(rbfiles) {
		min = len(rbfiles)
	}

	for i := 0; i < min; i++ {
		//log.Println("Iteration #", i)

		lernLangFile(n, gofiles[i], []float64{1, 0})
		//lernLangFile(n, rbfiles[i], []float64{0, 1})
	}

}

func lernLangFile(n *neural.Network, path string, out []float64) {
	//log.Println("Lerning ", path)
	sample := getSampleFromFile(path)
	lern.Lern(n, sample, out, 0.1)
}

func getSampleFromFile(path string) []float64 {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	sample := make([]float64, 300)

	r := bufio.NewReader(file)

	for i := 0; i < 299; i++ {
		b, err := r.ReadByte()
		if err != nil {
			break
		}
		sample[i] = float64(b)
	}

	return sample
}

func loadEngine() engine.Engine {
	n := loadNetwork()

	e := engine.New(n)
	e.Start()
	return e
}

func getGoFiles() []string {
	return getLinesFromFile("/tmp/gofiles")
}

func getRbFiles() []string {
	return getLinesFromFile("/tmp/rbfiles")
}

func getLinesFromFile(path string) []string {
	lines := make([]string, 0)
	file, err := os.Open(path)

	defer file.Close()

	if err != nil {
		log.Fatal(err)
	}

	r := bufio.NewReader(file)

	for {
		line, err := r.ReadBytes(10)
		if err != nil {
			break
		}
		lines = append(lines, string(line[:len(line)-1]))
	}

	return lines
}

func createLangNetwork() {
	n := neural.NewNetwork(300, []int{300, 300, 2})
	n.RandomizeSynapses()

	persist.ToFile(jsonFile, n)
}
