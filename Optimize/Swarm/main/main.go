package main

import( 
	"Swarm"
	"math"
	"fmt"
)
func max(x, y float64) float64 {
	if x >= y {
		return x
	}
	return y
}

func fitness(x []float64) float64{ 
	mdist := math.Sqrt(2*100.0*100)/2.0
	pdist := math.Sqrt((x[0] -20.0)*(x[0] -20.0) + (x[1]-7.0)*(x[1]-7.0))
	return 100.0*(1-pdist/mdist)
}

func otherFitness(x []float64) float64{
	mdist := math.Sqrt(2*100.0*100)/2.0
	pdist := math.Sqrt((x[0] -20.0)*(x[0] -20.0) + (x[1]-7.0)*(x[1]-7.0))
	ndist := math.Sqrt((x[0] +20.0)*(x[0] +20.0) + (x[1]+7.0)*(x[1]+7.0))
	return 9.0*max(0, 10 - pdist*pdist) + 10*(1-pdist/mdist) + 70*(1-ndist/mdist)
}

func main(){
	bestSol :=  make([]float64, 2)
	var bestFit float64

	bounds := make([]Swarm.Bound, 2)
	bounds[0] = Swarm.Bound{-50.0,50.0}
	bounds[1] = Swarm.Bound{-50.0,50.0}
	
	theSwarm := Swarm.NewSwarm(240, 2, otherFitness, bounds)
	theSwarm.Init()
	//theSwarm.SetRadius(10)
	for i:=0; i < 100; i++ {
		theSwarm.Update()
	}
	
	bestSol, bestFit = theSwarm.GetBest()
	fmt.Println(bestFit)
	fmt.Println(bestSol[0])
	fmt.Println(bestSol[1])
}