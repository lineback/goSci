package goSci

import "math"

/*
 Returns the mean of the array
*/
func Mean(x *GsArray, meanType uint) *GsArray {
	sum := Sum(x, meanType)
	switch meanType {
	case ALL: 
		return Times(sum, 1.0 / float64(len(x.data)))
	case COLS: 
		return Times(sum, 1.0 / float64(x.shape[0]))
	case ROWS:
		return Times(sum, 1.0 / float64(x.shape[1]))
	default:
		panic("Invalid mean type.")
	}
	return new(GsArray)
}
/*
 Returns the standard deviation of the array
*/
func Stdev(x *GsArray, stdevType uint) *GsArray {
	mean := Mean(x, stdevType)
	var meanArray *GsArray
	switch stdevType {
	case ALL:
		meanArray = Repmat(mean, x.shape[0], x.shape[1])
	case COLS:
		meanArray = Repmat(mean, x.shape[0], 1)
	case ROWS:
		meanArray = Repmat(mean, 1, x.shape[1])
	default:
		panic("Invalid stdev type.")
	}
	diff := Minus(x, meanArray)
	diffSquared := ElemTimes(diff, diff)
	meanDiff := Mean(diffSquared, stdevType)
	return ArrayFun(meanDiff, math.Sqrt)
}