package goSci


/*
  #cgo LDFLAGS: -lblas
  #include <atlas/cblas.h>
*/
import 	"C"
import(	
         "unsafe"
)

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
func Plus(x, y *GsArray) *GsArray {
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
 Accumulates the addtion of two matrices.
 Panics if the dimensions are incorrect
*/
func (x *GsArray) PlusEquals(y *GsArray) {
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
	for i := 0; i < len(x.data) ; i++ {
		x.data[i] = x.data[i] + y.data[i]
	}
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
 Accumulates the difference of two matricies.
 Panics if the dimensions are incorrect
*/
func (x *GsArray) MinusEquals(y *GsArray) {
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
	
	for i := 0; i < len(x.data) ; i++ {
		x.data[i] = x.data[i] - y.data[i]
	}
}
/*
 Returns a*x where a is a scalar and x is an array
*/
func Times(x *GsArray, a float64) *GsArray {
	result := Zeros(x.shape ...)
	for i := 0; i < len(x.data) ; i++ {
		result.data[i] = x.data[i]*a 
	}
	return result
}

/*
 Calculates a*x and stores the result in x where a is a scalar and x is an array
*/
func (x *GsArray) TimesEquals(a float64) {
	for i := 0; i < len(x.data) ; i++ {
		x.data[i] = x.data[i]*a 
	}	
}
/*
 Returns element wise multiplication of x and y
*/
func ElemTimes(x, y *GsArray) *GsArray {
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
		result.data[i] = x.data[i] * y.data[i]
	}
	return result
}
/*
 Returns the sum of all of the elements in the array
*/
func Sum(x *GsArray, sumType uint) *GsArray{
	if sumType == ALL{
		sum := Zeros(1,1)
		for _,val := range x.data {
			sum.data[0] += val
		}
		return sum
	} else if sumType == COLS{
		if len(x.shape) < 2 {
			panic("There are no columns use goSci.ALL.")
		}
		sum := Zeros(1, x.shape[1])
		for i, val := range x.data {
			sum.data[i%x.shape[1]] += val
		}
		return sum
	} else if sumType == ROWS{
		if len(x.shape) < 2 {
			return Sum(x, ALL)
		}
		sum := Zeros(x.shape[0],1)
		for i, val := range x.data {
			sum.data[i/x.shape[1]] += val
		}
		return sum
	} else {
		panic ("Invalid sum type")
	}
	return Zeros(x.shape ...)
}
/*
Returns a rep1Xrep2 tiling of x,  x must be a vector or matrix
*/
func Repmat(x *GsArray, rep1, rep2 int) *GsArray {
	var dim1, dim2, x1, x2 int
	if len(x.shape) > 2 {
		panic("Shape of input array is invalid for Repmat")
	}
	x1 = x.shape[0]
	dim1 = x1 * rep1
	if len(x.shape) == 2 {
		x2 = x.shape[1]
		dim2 = x2 * rep2
	} else {
		x2 = 1
		dim2 = rep2
	}
	returnArray := Zeros(dim1, dim2)
	for i := 0; i < dim1; i++ {
		origIdx := i%x1
		for j := 0; j < dim2; j++ {
			returnArray.Put(x.Get(origIdx, j%x2), i,j)
		}
	}
	return returnArray
}
/*
 Applies someFunc to every element of an array and returns the new array
*/
func ArrayFun(x *GsArray, someFunc func(float64) float64) *GsArray {
	returnArray := Zeros(x.shape ...)
	for i, val := range x.data {
		returnArray.data[i] = someFunc(val)
	}
	return returnArray
}