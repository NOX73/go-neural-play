package play

import (
	"github.com/NOX73/go-neural"
	"github.com/NOX73/go-neural/learn"
	//"github.com/NOX73/go-neural/persist"
	"log"
	"math"
	"math/rand"
	"time"
)

func MulriplyMain() {

	network := neural.NewNetwork(2, []int{5, 5, 1})
	network.RandomizeSynapses()

	ch := make(chan float64, 300)

	go func() {
		tick := time.Tick(5 * time.Second)
		acount := 1000
		arr := make([]float64, acount)
		index := 0

		for {
			select {
			case v := <-ch:
				arr[index] = v
				index++
				if index == len(arr) {
					index = 0
				}
			case <-tick:
				var sum float64 = 0
				for _, val := range arr {
					sum += val
				}

				//log.Println(index)
				//log.Println(arr)
				log.Printf("Avarege error: %f", sum/float64(acount))

			}
		}

	}()

	count := 100000000000
	for i := 0; i < count; i++ {
		test := []float64{rand.Float64(), rand.Float64()}
		result := network.Calculate(test)[0]

		ch <- math.Abs(test[0]*test[1] - result)

		log.Printf("Error value: %f * %f = %f", test[0], test[1], test[0]*test[1])

		learn.Learn(network, test, []float64{test[0] * test[1]}, 0.1)
	}

}
