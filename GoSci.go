package GoSci

import "fmt"

type GsArray struct {
	data        []float64
	shape       []int
}

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

func (array *GsArray) Reshape(shape ... int) {
	product := 1
	for _, value := range shape {
		product *= value
	}
	if product != len(array.data) {
		fmt.Println("ValueError: total size of new array must be unchanged2")
	}
	array.shape = shape
}

func (array *GsArray) Print() {
	for _, value := range array.data {
		fmt.Println(value)
	}
}

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
