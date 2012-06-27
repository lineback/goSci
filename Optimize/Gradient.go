package Optimize

type FuncGrad interface {
	Gradient([]float64) []float64
	Feval([]float64) float64
}

func GradientDescent(minObj FuncGrad, start []float64, maxIter int, lr float64) (float64, []float64) {
	currBest := make([]float64, len(start))
	copy(currBest, start)
	
	for i := 0; i < maxIter; i++ {
		grad := minObj.Gradient(currBest)
		for idx, val := range grad {
			currBest[idx] -= lr*val
		}
	}
	return minObj.Feval(currBest), currBest
}