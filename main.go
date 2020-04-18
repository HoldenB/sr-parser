package main

import "container/list"

////////////////////////////////////////////////////////////////

/********* Parse Stack Item *********/
type parseStackItem struct {
	grammarSymbol, stateSymbol string
}

func (item parseStackItem) String()  string {
	return item.grammarSymbol + item.stateSymbol
}


/********* Parse Stack *********/

// using a struct w/ one field to help simulate call by ref
// needed to obtain the "stack" side-effects w/ push/pop
type parseStack struct {
	stack *list.List
}

func newParseStack() parseStack {
	pStack := parseStack{}
	pStack.stack = list.New()
	return pStack
}




func main() {
	//reader := bufio.NewScanner(os.Stdin)
	//fmt.Println("Enter a file to parse. Must be in this directory.")
	//file := ""
	//for reader.Scan() {
	//	if _, err := os.Stat(reader.Text()); !os.IsNotExist(err) {
	//		file = reader.Text()
	//		break
	//	} else {
	//		fmt.Println("File not in directory")
	//	}
	//}
	//
	//filepath := file
	//
	//filebuffer, err := ioutil.ReadFile(filepath)
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//
	//inputdata := string(filebuffer)
	//data := bufio.NewScanner(strings.NewReader(inputdata))
	//data.Split(bufio.ScanRunes)

	inputArray := []string{"id", "+", "id", "*", "id"}


}
