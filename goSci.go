package goSci

/*
  #cgo LDFLAGS: -lblas
  #include <atlas/cblas.h>
*/
import 	"C"
import(	
	 "fmt"
         "unsafe"
         "bytes"
)
import "math"

type GsArray struct {
	data        []float64
	shape       []int
}

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
		fmt.Println("ValueError: total size of new array must be unchanged2")
	}
	array.shape = shape
}


func (array *GsArray) PrintShape(){
	for _,value := range array.shape {
		fmt.Println(value)
	}
}

func (array *GsArray) Print() {
	for _, value := range array.data {
		fmt.Println(value)
	}
}
/*
 Puts val into postion given by pos. 
 e.g. 
 array.Put(1.5, 4, 3)
 would put 1.5 into row 4 column 3
*/
func (array *GsArray) Put(val float64, pos []int) {
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
func (array *GsArray) Get(pos []int) float64 {
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
 Returns the dot product of two one dimensional arrays
 Panics if either of the arrays is not one dimensional or if the lengths are not equal
*/
func Dot(x, y *GsArray) float64 {
	if len(x.shape) > 2  || len(y.shape) > 2 {
		panic("Invalid dimension for dot product!!")
	}
	if len(x.shape) == 2 && len(y.shape) == 2 {
		if x.shape[0] != 1 && x.shape[1] != 1 {
			if y.shape[0] != 1 && y.shape[1] != 1 {
				panic("Invalid dimension for dot product!!")
			}
		}
	}else if len(x.shape) == 1 && len(y.shape) == 2 {
		if (y.shape[0] != 1) && (y.shape[1] != 1) {
			panic("Invalid dimension for dot product!!")
		}
	} else if (len(x.shape) == 2 && len(y.shape) == 1) {
		if (x.shape[0] != 1) && (x.shape[1] != 1) {
			panic("Invalid dimension for dot product!!")
		}
	}
	if len(x.data) != len(y.data){
		panic("Vectors must be of the same lenght!!")
	}
	c_N:= C.int(len(x.data))
	c_x := (*C.double)(unsafe.Pointer(&x.data[0]))
	c_incX := C.int(1)
	c_y := (*C.double)(unsafe.Pointer(&y.data[0]))
	c_incY := C.int(1)
	
	return float64(C.cblas_ddot(c_N, c_x, c_incX, c_y, c_incY))
}
/*
 Returns the matrix multiplication of two GsArrays.
 Panics if the the dimensions of the arrays are invalid for matrix multiplication
*/
func MatMult(x, y *GsArray) *GsArray {
	if len(x.shape) != 2 || len(y.shape) != 2 {
		panic("Arrays must have dimension 2 for matrix multiply")
	}
	if x.shape[1] != y.shape[0] {
		panic("Invalid dimensions for matrix multiply")
	}
	c_M := C.int(x.shape[0])
	c_N := C.int(y.shape[1])
	c_K := C.int(x.shape[1])
	c_Lda := c_K
	c_Ldb := c_M
	c_Ldc := c_N
	c_alpha := C.double(1.0)
	c_beta := C.double(1.0)
	
	c_x := (*C.double)(unsafe.Pointer(&x.data[0]))
	c_y := (*C.double)(unsafe.Pointer(&y.data[0]))
	z := Zeros(x.shape[0],y.shape[1])
	c_z := (*C.double)(unsafe.Pointer(&(z.data[0])))
	
	C.cblas_dgemm(101, 111, 111, c_M, c_N, c_K, c_alpha, c_x, c_Lda, c_y, c_Ldb, c_beta, c_z, c_Ldc)
	
	return z
}
/*
 Returns the addtion of two matrices.
 Panics if the dimensions are incorrect
*/
func Add(x, y *GsArray) *GsArray {
	if len(x.data) != len(y.data){
		panic("Arrays must have the same shape!!")
	}
	if len(x.shape) != len(y.shape){
		panic("Arrays must have the same shape!!")
	}
	for i := 0; i < len(x.shape); i++{
		if x.shape[i] != y.shape[i]{
			panic("Arrays must have the same shape!!")
		}
	}
	result := Zeros(x.shape ...)
	for i := 0; i < len(x.data) ; i++ {
		result.data[i] = x.data[i] + y.data[i]
	}
	return result
}
/*
 Returns the x - y.
 Panics if the dimensions are incorrect
*/
func Minus(x, y *GsArray) *GsArray {
	if len(x.data) != len(y.data){
		panic("Arrays must have the same shape!!")
	}
	if len(x.shape) != len(y.shape){
		panic("Arrays must have the same shape!!")
	}
	for i := 0; i < len(x.shape); i++{
		if x.shape[i] != y.shape[i]{
			panic("Arrays must have the same shape!!")
		}
	}
	result := Zeros(x.shape ...)
	for i := 0; i < len(x.data) ; i++ {
		result.data[i] = x.data[i] - y.data[i]
	}
	return result
}
/*
 Returns a*x where a is a scalar and x is an array
*/
func ScalarMult(x *GsArray, a float64) *GsArray {
	result := Zeros(x.shape ...)
	for i := 0; i < len(x.data) ; i++ {
		result.data[i] = x.data[i]*a 
	}
	return result
}
/*
 Returns the sum of all of the elements in the array
*/
func Sum(x *GsArray) float64 {
	sum := float64(0)
	for _,val := range x.data {
		sum += val
	}
	return sum
}
/*
 Returns the mean of the array
*/
func Mean(x *GsArray) float64 {
	sum := Sum(x)
	return sum/float64(len(x.data))
}
/*
 Returns the standard deviation of the array
*/
func Stdev(x *GsArray) float64 {
	mean := Mean(x)
	meanArray := ScalarMult(Ones(len(x.data)), mean)
	shape := x.shape
	x.Reshape(len(x.data))
	diff := Minus(x, meanArray)
	std := math.Sqrt(Dot(diff, diff)/float64(len(x.data)))
	x.Reshape(shape ...)
	return std
}
/*
 Stringer function for printing the arrays.
 Only prints arrays of dimension two or less
*/
func (array *GsArray) String() string {
	if len(array.shape) > 2 {
		return "I only print arrays with dimension less than 2."
	}
	buff := bytes.NewBufferString("")
	for i:=0; i<len(array.data); i++ {
		if i % array.shape[1] == 0 && i != 0 {
			fmt.Fprint(buff, "\n")
		}
		fmt.Fprintf(buff, "%f ", array.data[i])
	}
	return buff.String()
}