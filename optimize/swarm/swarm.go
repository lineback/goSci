package swarm

import (
	"math/rand"
	"math"
	"time"
)
		
type Bound struct {
	Lower   float64
	Upper   float64
}

type Swarm struct {
	dim           int        //dimension of the system
	np            int        //number of particles
	radius        int        //radius for ring lattice
	v_max         float64    //maximum particle velocity
	position      []float64  //current position of particles
	velocity      []float64  //current velocity
	fitness       []float64  //current fitness
	best_pos      []float64  //best position
	best_fit      []float64  //best fitness
	lbest         []int      //local best in lattice
	gbest         float64    //global best
	gbestIdx      int        //index of the global best
	inertia       float64    //inertia
	cog           float64    //cognitive parameter
	soc           float64    //social parameter
	fitFunc       func([]float64) float64 //fitness function
	ring          bool       //use ring lattice
	bounds        []Bound    //bounds of search
}

func NewSwarm(np int, dim int, fit func([]float64) float64, bounds []Bound) *Swarm {
	swarm := new(Swarm)
	swarm.dim = dim
	swarm.np = np
	swarm.radius = 0
	swarm.v_max = 1.0
	swarm.position = make([]float64, dim*np)
	swarm.velocity = make([]float64, dim*np)
	swarm.fitness = make([]float64, np)
	swarm.best_pos = make([]float64, dim*np)
	swarm.best_fit = make([]float64, dim*np)
	swarm.lbest = make([]int, np)
	swarm.gbest = float64(math.Inf(-1))
	swarm.gbestIdx = 0
	swarm.inertia = 1.0
	swarm.cog = 2.0
	swarm.soc = 2.0
	swarm.ring = false 
	swarm.fitFunc = fit
	swarm.bounds = bounds
	now := time.Now()
	rand.Seed(now.Unix())
	return swarm
}

func (swarm *Swarm) Init(){
	
	var boundSize float64
	var lower float64

	for i := 0; i < swarm.np; i++ {
		for j := 0; j < swarm.dim; j++ {
			boundSize = (swarm.bounds[j]).Upper - (swarm.bounds[j]).Lower
			lower = (swarm.bounds[j]).Lower
			swarm.position[i*swarm.dim+j] = rand.Float64()*boundSize + lower
			swarm.best_pos[i*swarm.dim+j] = swarm.position[i*swarm.dim+j]
			swarm.velocity[i*swarm.dim+j] = rand.Float64()
		}
		
		swarm.lbest[i] = i
	}
}

func (swarm *Swarm) calcFitness(){
	var start, end int
	for i := 0; i < swarm.np; i++ {
		start = i*swarm.dim
		end = start + swarm.dim
		swarm.fitness[i] = swarm.fitFunc(swarm.position[start:end])
		if swarm.fitness[i] > swarm.gbest {
			swarm.gbest = swarm.fitness[i]
			swarm.gbestIdx = i
		}
		if swarm.fitness[i] > swarm.best_fit[i] {
			swarm.best_fit[i] = swarm.fitness[i]
			copy(swarm.best_pos[start:end], swarm.position[start:end])
		}
	}
	if swarm.ring {
		swarm.setLocalBest()
	}
}

func (swarm *Swarm) setLocalBest(){	
	var tempIdx int
	var tempFit float64

	for i := 0; i < swarm.np; i++ {
		tempFit = swarm.best_fit[swarm.lbest[i]]
		for j := 0; j < swarm.radius; j++ {
			tempIdx = ((j%swarm.np) + swarm.np)%swarm.np
			if swarm.best_fit[tempIdx] > tempFit {
				tempFit = swarm.best_fit[tempIdx]
				swarm.lbest[i] = tempIdx
			}
		}
	}
		
}

func (swarm *Swarm) updateVelocity() {
	var rand1, rand2, cogTerm, socTerm, mag float64
	
	for i := 0; i < swarm.np; i++ {
		rand1 = rand.Float64()
		rand2 = rand.Float64()
		for j := 0; j < swarm.dim; j++ {
			cogTerm = swarm.cog*rand1*(swarm.best_pos[i*swarm.dim+j]-swarm.position[i*swarm.dim+j])
			if swarm.ring {
				socTerm = swarm.best_pos[swarm.lbest[i]*swarm.dim+j]
				socTerm -= swarm.position[i*swarm.dim+j]
				socTerm *= (swarm.soc*rand2)
			}else {
				socTerm = swarm.best_pos[swarm.gbestIdx*swarm.dim+j]
				socTerm -= swarm.position[i*swarm.dim+j]
				socTerm *= (swarm.soc*rand2)
			}
			swarm.velocity[i*swarm.dim+j] = swarm.inertia*swarm.velocity[i*swarm.dim+j] + cogTerm + socTerm
		}
		mag = swarm.magVel(i)
		if mag > swarm.v_max {
			for j := 0; j < swarm.dim; j++ {
				swarm.velocity[i*swarm.dim+j] = swarm.velocity[i*swarm.dim+j]*swarm.v_max/mag
			}
		}
	}
}

func (swarm *Swarm) magVel(idx int) float64{
	var mag float64
	mag = 0.0
	for i:=0; i<swarm.dim; i++ {
		mag += swarm.velocity[idx*swarm.dim +i]*swarm.velocity[idx*swarm.dim +i]
	}
	return mag
}

func (swarm *Swarm) updatePosition(){
	for i:=0; i < len(swarm.position); i++ {
		swarm.position[i] += swarm.velocity[i]
	}
}

func (swarm *Swarm) Update(){
	swarm.calcFitness()
	swarm.updateVelocity()
	swarm.updatePosition()
}

func (swarm *Swarm) GetBest() ([]float64, float64) {
	start := swarm.dim*swarm.gbestIdx
	end := start + swarm.dim
	return swarm.best_pos[start:end], swarm.gbest
}

func (swarm *Swarm) Optimize(iterations int) ([]float64, float64) {
	for i := 0; i < iterations; i++ {
		swarm.Update()
	}
	return swarm.GetBest()
}

func (swarm *Swarm) SetRadius(radius int){
	swarm.radius = radius
	swarm.ring = true
}

func (swarm *Swarm) SetCog(cog float64){
	swarm.cog = cog
}

func (swarm *Swarm) SetSoc(soc float64){
	swarm.soc = soc
}

func (swarm *Swarm) SetVmax(v_max float64){
	swarm.v_max = v_max
}
