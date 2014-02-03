package play

import (
	"bufio"
	"github.com/NOX73/go-neural"
	"github.com/NOX73/go-neural/engine"
	"github.com/NOX73/go-neural/lern"
	"github.com/NOX73/go-neural/persist"
	"github.com/cheggaaa/pb"
	"log"
	"math/rand"
	"os"
	"os/signal"
)

func LangMain() {

	createLangNetwork()

	n := loadNetwork()

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
	gosamplePath = "./json/sample/sample.go"
	rbsamplePath = "./json/sample/sample2.rb"
	jssamplePath = "./json/sample/sample.js"

	sampleLen    = 200
	lerningSpeed = 0.1
)

func checkEngine(n *neural.Network) bool {

	gosample := getSampleFromFile(gosamplePath)
	rbsample := getSampleFromFile(rbsamplePath)
	jssample := getSampleFromFile(jssamplePath)

	var outgo, outrb, outjs []float64

	outgo = n.Calculate(gosample)
	outrb = n.Calculate(rbsample)
	outjs = n.Calculate(jssample)

	log.Println("Check go", outgo)
	log.Println("Check rb", outrb)
	log.Println("Check js", outjs)

	if outgo[0] < 0.9 || outrb[1] < 0.9 || outjs[2] < 0.9 {
		return false
	}

	return true
}

func testEngine(n *neural.Network) {
	gosample := getSampleFromFile(gosamplePath)
	rbsample := getSampleFromFile(rbsamplePath)
	jssample := getSampleFromFile(jssamplePath)

	log.Println(gosamplePath, n.Calculate(gosample))
	log.Println(rbsamplePath, n.Calculate(rbsample))
	log.Println(jssamplePath, n.Calculate(jssample))
}

func lernEngine(n *neural.Network) {

	gofiles := getGoFiles()
	rbfiles := getRbFiles()
	jsfiles := getJsFiles()

	count := 1000
	bar := pb.StartNew(count)

	for i := 0; i < count; i++ {
		bar.Increment()
		lernLangFile(n, gofiles[rand.Intn(len(gofiles))], []float64{1, 0, 0})
		lernLangFile(n, rbfiles[rand.Intn(len(rbfiles))], []float64{0, 1, 0})
		lernLangFile(n, jsfiles[rand.Intn(len(jsfiles))], []float64{0, 0, 1})
	}
	bar.Finish()

}

func lernLangFile(n *neural.Network, path string, out []float64) {
	//log.Println("Lerning ", path)
	sample := getSampleFromFile(path)
	lern.Lern(n, sample, out, lerningSpeed)
}

func getSampleFromFile(path string) []float64 {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	sample := make([]float64, sampleLen)
	stat, _ := file.Stat()
	size := stat.Size()

	if size > int64(sampleLen) {
		offset := rand.Int63n(size - int64(sampleLen))
		file.Seek(offset, 0)
	}

	r := bufio.NewReader(file)

	for i := 0; i < sampleLen-1; i++ {
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

func getJsFiles() []string {
	return getLinesFromFile("/tmp/jsfiles")
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
	if _, err := os.Stat(jsonFile); err == nil {
		return
	}
	n := neural.NewNetwork(sampleLen, []int{1000, 1000, 3})
	n.RandomizeSynapses()

	persist.ToFile(jsonFile, n)
}
