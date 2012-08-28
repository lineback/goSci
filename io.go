package goSci

import(
//        "io/ioutil"
        "strings"
        "errors"
        "bytes"
        "fmt"
        "os"
        "bufio"
)

/* 
 Stringer function for printing the arrays.
 Only prints arrays of dimension two or less
*/
func (array *GsArray) String() string {
	if len(array.shape) > 2 {
		return "I only print arrays with dimension less than 3."
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

/*
Loads table from file, fileName, with each element separated by delimiter, delim.  
Each line corresponds to a row in the matix.  Lines starting with '#' are considered comments and ignored
*/
func LoadTable(fileName string, delim string) (*GsArray, error){
	file, err := os.Open(fileName)
	if err != nil{
		return new(GsArray), err
	}
	defer file.Close()
	lines := make([] string, 0)
	r := bufio.NewReader(file)
	line, err := r.ReadString('\n')
	for err == nil {
		if (line[0] != '\n' && line[0] != '#'){
			lines = append(lines, line)
		}
		line, err = r.ReadString('\n')
	}
	
	cols := len(strings.Split(lines[0], delim))
	rows := len(lines)

	retArray := Zeros(rows, cols)
	
	for i,line := range lines {
		vals := strings.Split(line, delim)
		
		if len(vals) != cols{
			return retArray, errors.New("Invalid file format: number of columns is not constant.")
		}
		for j, val := range vals {
			var tempVal float64
			numRead,_ := fmt.Sscan(val, &tempVal)
			if numRead != 1 {
				return retArray, errors.New("Unable to parse file.")
			}
			retArray.Put(tempVal,i, j) 
		}
	}
	return retArray, nil
}

/*
Writes the array to a file, fileName,  as a table with the given delimiter, delim
*/
func (array *GsArray) WriteTable(fileName string, delim string) error {
	if len(array.shape) > 2 {
		return errors.New("Tables only valid for 2 dimensional arrays")
	}
	file, err := os.Create(fileName)
	if err != nil{
		return err
	}
	defer file.Close()
	
	buff := bytes.NewBufferString("")
	for i:=0; i<len(array.data); i++ {
		if i % array.shape[1] == 0 && i != 0 {
			fmt.Fprint(buff, "\n")
		} else if i != 0{
			fmt.Fprintf(buff, "%s", delim)
		}
		fmt.Fprintf(buff, "%f", array.data[i])
	}
	w := bufio.NewWriter(file)
	_, err = w.WriteString(buff.String())
	if err != nil{
		return err
	}
	if err = w.Flush(); err != nil{
		return err
	}
	return nil
}

