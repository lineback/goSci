package goSci

import "math"

type GsArray struct {
	data        []float64
	shape       []int
}

const(
	ALL = iota
        COLS
        ROWS
)
/*
 Creates a GsArray with shape defined by shape initialized to zero
*/
func Zeros(shape ... int) *GsArray {
	product := 1
	for _, value := range shape {
		product *= value
	}
	array := new(GsArray)
	array.data = make([]float64, product)
	array.shape = shape
	return array
}

/*
 Creates a GsArray with shape defined by shape intialized to ones
*/
func Ones(shape ... int) *GsArray {
	product := 1
	for _, value := range shape {
		product *= value
	}
	array := new(GsArray)
	array.data = make([]float64, product)
	array.shape = shape
	for i := 0; i < product; i++ {
		array.data[i] = 1
	} 
	return array
}

/*
 Creates a one dimensional GsArray with the range 0 to size -1
*/
func Arange(size int) *GsArray {
	array := new(GsArray)
	array.shape = make([]int, 1) 
	array.shape[0] = size
	array.data = make([]float64, size)
	for i := 0; i < size; i++ {
		array.data[i] = float64(i)
	}
	return array
}

/*
 Creates an identity matrix of size size
*/
func Eye(size int) *GsArray {
	array := Zeros(size,size)
	for i := 0; i < size; i++ {
		array.data[i*size + i] = float64(1)
	}
	return array
}

/*
 Reshapes array to new shape. The product of shape must be the same as product of old shape.
*/
func (array *GsArray) Reshape(shape ... int) {
	product := 1
	for _, value := range shape {
		product *= value
	}
	if product != len(array.data) {
		panic("ValueError: total size of new array must be unchanged")
	}
	array.shape = shape
}

/*
 Puts val into postion given by pos. 
 e.g. 
 array.Put(1.5, 4, 3)
 would put 1.5 into row 4 column 3
*/
func (array *GsArray) Put(val float64, pos ... int) {
	if len(array.shape) != len(pos){
		panic("Invalid posistion!")
	}
	if len(array.shape) == 1{
		array.data[pos[0]] = val
	}else if len(array.shape) == 2 {
		dim2 := array.shape[1]
		array.data[pos[0]*dim2 + pos[1]] = val
	}else {
		prod := make([]int, len(pos))
		prod[len(pos) - 1] = 1
		for i:=2; i <= len(pos); i++ {
			prod[len(pos) - i]  = array.shape[len(pos) - i + 1] * prod[len(pos) - i + 1]
		}
		idx := 0
		for i:=0; i < len(pos) ; i++ {
			idx += prod[i] * pos[i]
		}
		array.data[idx] = val
	}
}
/*
 returns value of value at postion pos
*/
func (array *GsArray) Get(pos ... int) float64 {
	if len(array.shape) != len(pos){
		panic("Invalid posistion!")
	}
	if len(array.shape) == 1{
		return array.data[pos[0]]
	}
	if len(array.shape) == 2 {
		dim2 := array.shape[1]
		return array.data[pos[0]*dim2 + pos[1]]
	}
	prod := make([]int, len(pos))
	prod[len(pos) - 1] = 1
	for i:=2; i <= len(pos); i++ {
		prod[len(pos) - i]  = array.shape[len(pos) - i + 1] * prod[len(pos) - i + 1]
	}
	idx := 0
	for i:=0; i < len(pos) ; i++ {
		idx += prod[i] * pos[i]
	}
	return array.data[idx]
}

/*
 Returns the mean of the array
*/
func Mean(x *GsArray) float64 {
	sum := Sum(x, ALL)
	return sum.data[0]/float64(len(x.data))
}
/*
 Returns the standard deviation of the array
*/
func Stdev(x *GsArray) float64 {
	mean := Mean(x)
	meanArray := Times(Ones(len(x.data)), mean)
	shape := x.shape
	x.Reshape(len(x.data))
	diff := Minus(x, meanArray)
	std := math.Sqrt(Dot(diff, diff)/float64(len(x.data)))
	x.Reshape(shape ...)
	return std
}
