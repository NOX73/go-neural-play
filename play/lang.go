package play

import (
	"github.com/NOX73/go-neural"
	"github.com/NOX73/go-neural/engine"
	"github.com/NOX73/go-neural/persist"
)

func LangMain() {

	//createLangNetwork()
	lernEngine()

}

func lernEngine() {
	e := loadEngine()
	e.Start()

}

func loadEngine() engine.Engine {
	n := loadNetwork()

	e := engine.New(n)
	e.Start()

	return e
}

func createLangNetwork() {
	n := neural.NewNetwork(300, []int{300, 300, 2})
	n.RandomizeSynapses()

	persist.ToFile(jsonFile, n)
}
