package GA
import (
	"math/rand"
)
type GA struct {
	length     int       //length of the bit strings
	num        int       //number of individuals
	cp         float64   //cross over probability
	mp         float64   //mutation probability
	population []bool    //population of bit strings
	fitness    []float64 //fitness for each individual
	fitFunc    func([]bool) float64 //fitness function
}

func NewGA(length int, num int, fitFunc func([]bool) float64) (*GA){
	ga := new(GA)
	ga.length = length
	ga.num = num
	ga.cp = 0.6
	ga.mp = 1.0/float64(num)
	ga.fitFunc = fitFunc
	ga.population = make([]bool,length*num)
	ga.fitness = make([]float64, num)
	return ga
} 

func (ga *GA) Init(seed int64) {
	rand.Seed(seed)
	for i:=0; i < len(ga.population); i++ {
		if rand.Intn(2) == 1 {
			ga.population[i] = true
		}
	}
	ga.CalcFitness()
}

func (ga *GA) CalcFitness() {
	var start, end int
	for i:=0; i < ga.num; i++ {
		start = i*ga.length
		end = start + ga.length
		ga.fitness[i] = ga.fitFunc(ga.population[start:end])
	}
}

func (ga *GA) SetCrossover(cp float64) {
	ga.cp = cp
}

func (ga *GA) SetMutation(mp float64) {
	ga.mp = mp
}