/*
 Tools for a generic genetic algorithm in which the genes can be encoded as bit strings
*/
package genetic
import (
	"math/rand"
	"math"
)
// encapsulates the ga population
type GA struct {
	length     int       //length of the bit strings
	num        int       //number of individuals
	cp         float64   //cross over probability
	mp         float64   //mutation probability
	population []bool    //population of bit strings
	fitness    []float64 //fitness for each individual
	cumSum     []float64 //cummulative sum for crossover
	fitFunc    func([]bool) float64 //fitness function
}
/* 
 Creates a new ga with where:
           Each individual is a bit string of length length.
           The population consists of num individuals.
           The fitness of each individual is calculated using fitFunc  
*/
func NewGA(length int, num int, fitFunc func([]bool) float64) (*GA){
	ga := new(GA)
	ga.length = length
	ga.num = num
	ga.cp = 0.6
	ga.mp = 1.0/float64(num)
	ga.cumSum = make([]float64, num)
	ga.fitFunc = fitFunc
	ga.population = make([]bool,length*num)
	ga.fitness = make([]float64, num)
	
	return ga
} 
/*
 Initializes the value using an integer seed so runs can be reproduced
*/
func (ga *GA) Init(seed int64) {
	rand.Seed(seed)
	for i:=0; i < len(ga.population); i++ {
		if rand.Intn(2) == 1 {
			ga.population[i] = true
		}
	}
	ga.calcFitness()
}

func (ga *GA) calcFitness() {
	var start, end int
	for i:=0; i < ga.num; i++ {
		start = i*ga.length
		end = start + ga.length
		ga.fitness[i] = ga.fitFunc(ga.population[start:end])
		if i == 0 {
			ga.cumSum[i] = ga.fitness[i]
		}else {
			ga.cumSum[i] = ga.fitness[i] + ga.cumSum[i-1]
		}
	}
	for i := 0; i < ga.length; i++ {
		ga.cumSum[i] = ga.cumSum[i] / ga.cumSum[ga.length -1]
	}
}
/*
 Sets crossover probability
*/
func (ga *GA) SetCrossover(cp float64) {
	ga.cp = cp
}
/*
 Sets mutation probability
*/
func (ga *GA) SetMutation(mp float64) {
	ga.mp = mp
}

func (ga *GA) findNearest(val float64) int{
	max := ga.length
	min := 0
	mid := min
	for ; max >= min ; {
		mid := (min + max) / 2
		if ga.cumSum[mid] < val {
			min = mid + 1
		}else {
			max = mid - 1
		}
	}
	if mid == 0 {
		return 0
	} else if mid == ga.length - 1 {
		return ga.length -1
	} else {
		minIdx := mid -1
		minVal := math.Abs(ga.cumSum[minIdx])
		if math.Abs(ga.cumSum[mid]) < minVal{
			minVal = ga.cumSum[mid]
			minIdx = mid
		} 
		if math.Abs(ga.cumSum[mid+1]) < minVal{
			minIdx = mid + 1
		}
		return minIdx		
	}	
	return 0
}

func (ga *GA) crossover(idxa int , idxb int) ([]bool, []bool){
	newA := make([]bool, ga.length)
	newB := make([]bool, ga.length)
	crsovrpt := rand.Intn(ga.length - 1) + 1

	copy(newA[0:crsovrpt], ga.population[idxa*ga.length:idxa*ga.length + crsovrpt]);
	copy(newA[crsovrpt:ga.length], ga.population[idxb*ga.length + crsovrpt:(idxb+1)*ga.length])
	copy(newB[0:crsovrpt], ga.population[idxb*ga.length:idxb*ga.length + crsovrpt]);
	copy(newB[crsovrpt:ga.length], ga.population[idxa*ga.length + crsovrpt:(idxa+1)*ga.length])

	return newA, newB
}
/*
 Creates the next generation from the current population
*/
func (ga *GA) CreateNextGen() {
	ga.calcFitness()
	
	nextPop := make([]bool, len(ga.population))

	for i := 0; i < ga.num / 2; i ++ {
		randA := rand.Float64()
		randB := rand.Float64()
		
		indA := make([]bool, ga.length)
		indB := make([]bool, ga.length)

		idxA := ga.findNearest(randA)
		idxB := ga.findNearest(randB)
		
		if rand.Float64() < ga.cp {
			indA, indB = ga.crossover(idxA, idxB)
		} else {
			copy(indA, ga.population[idxA*ga.length:(idxA+1)*ga.length])
			copy(indB, ga.population[idxB*ga.length:(idxB+1)*ga.length])
		}
		
		for j := 0; j < ga.length ; j++ {
			if rand.Float64() < ga.mp {
				indA[j] = !indA[j]
			}
			if rand.Float64() < ga.mp {
				indB[j] = !indB[j]
			}
		}
		copy(nextPop[i*2*ga.length:(i*2+1)*ga.length], indA)
		copy(nextPop[(i*2+1)*ga.length:(i*2+2)*ga.length], indB)
	}
	copy(ga.population, nextPop)
}
/*
 Returns the best individual from the population, and its fitness
*/
func (ga *GA) GetBest() (float64, []bool) {
	ga.calcFitness()
	bestIdx := 0
	best := math.Inf(1)
	for idx, val := range ga.fitness {
		if val < best {
			best = val
			bestIdx = idx
		}
	}
	return best, ga.population[bestIdx*ga.length:(bestIdx+1)*ga.length]
}
/*
 Returns the average fitness of the population
*/
func (ga *GA) GetAveFitness() float64 {
	retVal := float64(0)
	for _,val := range ga.fitness {
		retVal += val
	}
	retVal = retVal / float64(ga.num)
	return retVal
}