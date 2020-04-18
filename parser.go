package main

import (
	"fmt"
	"strconv"
	"strings"
)

///////////////////////////////////////////////////////////////////
type inputQueue struct {
	list []string
}

func (iq *inputQueue) createEmptyListQueue() {
	iq.list = nil
}

func (iq *inputQueue) createQueue(list []string) {
	iq.list = list
	iq.list = append(iq.list, "$")
}

func (iq *inputQueue) pop() string {
	a := iq.list[0]
	iq.list = iq.list[1:]
	return a
}

func (iq *inputQueue) first() string {
	return iq.list[0]
}

func (iq *inputQueue) toString() string {
	var sb strings.Builder
	for i := 0; i < len(iq.list); i++ {
		sb.WriteString(iq.list[i])
	}
	return sb.String()
}

///////////////////////////////////////////////////////////////////

type lhssym struct {
	sym string
}

func (lh *lhssym) constructor(a string) {
	lh.sym = a
}

func (lh lhssym) printTreeWork(indent int) {
	a := ""
	for i := 0; i < indent; i++ {
		a = a + "    "
	}
	fmt.Println(a + lh.sym)
}

func (lh lhssym) toString() string {
	return lh.sym
}

///////////////////////////////////////////////////////////////////

type listOfTrees struct {
	list []parseTree
}

func (lot *listOfTrees) constructor(tree parseTree) {
	lot.list = append(lot.list, tree)
}

func (lot listOfTrees) printTreeWork(indent int) {
	for i := 0; i < len(lot.list); i++ {
		pt := lot.list[i]
		pt.printTreeWork(indent)
	}
}

func (lot *listOfTrees) toString() string {
	a := ""
	for _, e := range lot.list {
		b := e.toString()
		a += " " + b
	}
	return a
}

///////////////////////////////////////////////////////////////////

type nonLeafTree struct {
	parent   *lhssym
	children *listOfTrees
}

func constructor(par *lhssym, lot *listOfTrees, nlt *nonLeafTree) *nonLeafTree {
	nlt.parent = par
	nlt.children = lot
	return nlt
}

func (nlt nonLeafTree) printTree() {
	nlt.printTreeWork(0)
}

func (nlt nonLeafTree) printTreeWork(indent int) {
	nlt.parent.printTreeWork(indent)
	nlt.children.printTreeWork(indent + 1)
}

func (nlt nonLeafTree) toString() string {
	return "[" + nlt.parent.toString() + " " + nlt.children.toString() + "]"
}

///////////////////////////////////////////////////////////////////

type parseStack struct {
	list []pstackEntry
}

func (ps *parseStack) pop() pstackEntry {
	n := len(ps.list) - 1
	a := ps.list[n]
	ps.list[n] = pstackEntry{"", ""}
	ps.list = ps.list[:n]
	return a
}

func (ps *parseStack) push(a pstackEntry) {
	ps.list = append(ps.list, a)
}

func (ps *parseStack) isEmpty() bool {
	if len(ps.list) == 0 {
		return true
	}
	return false
}

func (ps *parseStack) top() pstackEntry {
	return ps.list[len(ps.list)-1]
}

func (ps *parseStack) toString() string {
	a := ""
	for i := 0; i < len(ps.list); i++ {
		a = a + ps.list[i].toString()
	}
	return a
}

///////////////////////////////////////////////////////////////////

type parseTree interface {
	printTree()
	printTreeWork(indent int)
	toString() string
}

///////////////////////////////////////////////////////////////////

type pstackEntry struct {
	stateSym, grammarSym string
}

func (se *pstackEntry) setPstackEntry(a string, b string) {
	se.stateSym = a
	se.grammarSym = b
}

func (se *pstackEntry) toString() string {
	return se.grammarSym + se.stateSym
}

///////////////////////////////////////////////////////////////////

type termSym struct {
	tsym string
}

func (ts termSym) printTree() {
	ts.printTreeWork(0)
}

func (ts termSym) printTreeWork(indent int) {
	a := ""
	for i := 0; i < indent; i++ {
		a = a + "    "
	}
	fmt.Println(a + ts.tsym)
}

func construct(a string, ts *termSym) *termSym {
	ts.tsym = a
	return ts
}

func (ts termSym) toString() string {
	return ts.tsym
}

///////////////////////////////////////////////////////////////////

type treeStack struct {
	list []parseTree
}

func (tree *treeStack) pop() parseTree {
	n := len(tree.list) - 1
	a := tree.list[n]
	tree.list = tree.list[:n]
	return a
}

func (tree *treeStack) push(e parseTree) {
	tree.list = append(tree.list, e)
}

func (tree *treeStack) top() parseTree {
	return tree.list[len(tree.list)-1]
}

func (tree *treeStack) toString() string {
	a := ""
	for _, e := range tree.list {
		b := e.toString()
		a = b + " " + a
	}
	return a
}

///////////////////////////////////////////////////////////////////

type actionChoice int

const (
	accept        actionChoice = 0
	ungrammatical actionChoice = 1
	shift         actionChoice = 2
	reduce        actionChoice = 3
)

type parser struct {
	aTableIndex    map[string]int
	gTableIndex    map[string]int
	inputArray     []string
	outputTable    [20][11]string
	grammar        [6][]string
	choice         actionChoice
	actionValue    string
	newPush        *pstackEntry
	aTable         [12][6]string
	gTable         [12][3]string
	notGrammatical bool
	inputQueue     *inputQueue
	pStack         *parseStack
	tStack         *treeStack
}

func (p *parser) Initialize() {
	p.inputArray = []string{"id", "+", "id", "*", "id"}
	// p.inputArray = []string{"id"} // simplest possible grammatical input
	// p.inputArray = []string{ "id", "+", "id"};                       // simple grammatical input
	// p.inputArray = []string{ "id", "+", "id", "+", "id"};            // left assoc
	// p.inputArray = []string{ "(", "id", ")"};                        // parens1
	// p.inputArray = []string{ "id", "+", "(", "id", "+", "id", ")"};  // parens2
	// p.inputArray = []string{ "id", "+", "id", "*"};                  // ungrammatical input
	p.grammar = [6][]string{
		{"E", "->", "E", "+", "T"},
		{"E", "->", "T"},
		{"T", "->", "T", "*", "F"},
		{"T", "->", "F"},
		{"F", "->", "(", "E", ")"},
		{"F", "->", "id"},
	}
	p.choice = ungrammatical
	p.actionValue = ""
	p.newPush = nil
	p.aTable = [12][6]string{
		{"S5", "", "", "S4", "", ""},
		{"", "S6", "", "", "", "accept"},
		{"", "R2", "S7", "", "R2", "R2"},
		{"", "R4", "R4", "", "R4", "R4"},
		{"S5", "", "", "S4", "", ""},
		{"", "R6", "R6", "", "R6", "R6"},
		{"S5", "", "", "S4", "", ""},
		{"S5", "", "", "S4", "", ""},
		{"", "S6", "", "", "S11", ""},
		{"", "R1", "S7", "", "R1", "R1"},
		{"", "R3", "R3", "", "R3", "R3"},
		{"", "R5", "R5", "", "R5", "R5"},
	}
	p.gTable = [12][3]string{
		{"1", "2", "3"},
		{"", "", ""},
		{"", "", ""},
		{"", "", ""},
		{"8", "2", "3"},
		{"", "", ""},
		{"", "9", "3"},
		{"", "", "10"},
		{"", "", ""},
		{"", "", ""},
		{"", "", ""},
		{"", "", ""},
	}
	p.aTableIndex = make(map[string]int)
	p.gTableIndex = make(map[string]int)
	p.notGrammatical = false
	p.inputQueue = nil
	p.pStack = nil
	p.tStack = nil
	fmt.Print("\n\n")
	fmt.Print("Input: ", p.inputArray, "\n\n\n")
	p.outputTable = [20][11]string{}
	p.inputQueue = new(inputQueue)
	p.inputQueue.createQueue(p.inputArray)
	p.aTableIndex["id"] = 0
	p.aTableIndex["+"] = 1
	p.aTableIndex["*"] = 2
	p.aTableIndex["("] = 3
	p.aTableIndex[")"] = 4
	p.aTableIndex["$"] = 5
	p.gTableIndex["E"] = 0
	p.gTableIndex["T"] = 1
	p.gTableIndex["F"] = 2
	p.newPush = new(pstackEntry)
	p.newPush.setPstackEntry("0", "")
	p.pStack = new(parseStack)
}

func (p *parser) PrintParseTree() {
	fmt.Println("Parse Tree: ")
	p.tStack.top().printTree()
}

func (p *parser) Parse() {
	fmt.Println("                INPUT          ACTION    ACTION  VALUE   LENGTH  TEMP            GOTO      GOTO   STACK       ")
	fmt.Println("Stack           TOKENS         LOOKUP    VALUE   OF LHS  OF RHS  STACK           LOOKUP    VALUE  ACTION      PARSE TREE STACK")
	fmt.Println("______________________________________________________________________________________________________________________________")
	p.tStack = new(treeStack)
	p.Parse1step()
	for {
		if p.choice != reduce && p.choice != shift {
			break
		} else {
			p.Parse1step()
		}
	}
	if p.notGrammatical == false {
		fmt.Println("")
		p.PrintParseTree()
	}
}

func (p *parser) Parse1step() {
	valueofLHS := ""
	tempStack := ""
	gotoValueIndex1st := ""
	valueofLHSindex := -1
	lengthofRHS := 0
	lengthofRHSstr := ""
	gotoValue := ""
	initQue := p.inputQueue.toString()
	inputQueFront := p.inputQueue.first()
	p.pStack.push(*(p.newPush))
	initpStack := p.pStack.toString()
	newPushstr := ""
	gotoLookupstr := ""
	p.EvaluateActionChoice()
	actionIndex1st := p.pStack.top().stateSym
	switch p.choice {
	case accept:
		break
	case ungrammatical:
		p.notGrammatical = true
	case shift:
		a := []rune(p.actionValue)
		index := string(a[1])
		p.newPush = new(pstackEntry)
		p.newPush.setPstackEntry(index, p.inputQueue.first())
		if compareStrings(p.inputQueue.first(), "id") == true {
			b := new(termSym)
			b = construct("id", b)
			p.tStack.push(*b)
		}
		p.inputQueue.pop()
	case reduce:
		u := []rune(p.actionValue)
		h := string(u[1])
		q, err := strconv.Atoi(h)
		if err != nil {
			fmt.Println("err here")
		}
		valueofLHSindex = q
		valueofLHS = p.grammar[valueofLHSindex-1][0]
		lengthofRHS = len(p.grammar[valueofLHSindex-1]) - 2
		popped := new(parseStack)
		for i := 1; i <= lengthofRHS; i++ {
			popped.push(p.pStack.pop())
		}
		tempStack = p.pStack.toString()
		gotoValueIndex1st = p.pStack.top().stateSym
		y := []rune(gotoValueIndex1st)
		m := string(y[0])
		n, e := strconv.Atoi(m)
		if e != nil {
			fmt.Println("e here")
		}
		gotoValue = p.gTable[n][p.gTableIndex[valueofLHS]]
		p.newPush = new(pstackEntry)
		p.newPush.setPstackEntry(gotoValue, valueofLHS)
		if lengthofRHS == 1 {
			oldPT := p.tStack.pop()
			child := new(listOfTrees)
			child.constructor(oldPT)
			lhsSym := new(lhssym)
			lhsSym.constructor(valueofLHS)
			tr := new(nonLeafTree)
			tr = constructor(lhsSym, child, tr)
			p.tStack.push(tr)
		} else if lengthofRHS == 3 && compareStrings(popped.top().grammarSym, "(") {
			fmt.Println(popped.top().grammarSym)
			aarg2 := p.tStack.pop()
			aarg1 := new(termSym)
			aarg1 = construct("(", aarg1)
			aarg3 := new(termSym)
			aarg3 = construct(")", aarg3)
			llot := new(listOfTrees)
			llot.constructor(aarg1)
			llot.list = append(llot.list, aarg2)
			llot.list = append(llot.list, aarg3)
			pparent := new(lhssym)
			pparent.constructor(valueofLHS)
			ttr := new(nonLeafTree)
			ttr = constructor(pparent, llot, ttr)
			p.tStack.push(ttr)
		} else if lengthofRHS == 3 {
			arg3 := p.tStack.pop()
			arg1 := p.tStack.pop()
			popped.pop()
			operator := popped.top().grammarSym
			arg2 := new(termSym)
			arg2 = construct(operator, arg2)
			lot := new(listOfTrees)
			lot.constructor(arg1)
			lot.list = append(lot.list, arg2)
			lot.list = append(lot.list, arg3)
			parent := new(lhssym)
			parent.constructor(valueofLHS)
			tr := new(nonLeafTree)
			tr = constructor(parent, lot, tr)
			p.tStack.push(tr)
		}
	}
	if p.choice == shift || p.choice == reduce {
		newPushstr = "push " + p.newPush.toString()
	}
	if lengthofRHS > 0 {
		lengthofRHSstr = strconv.Itoa(lengthofRHS)
	}
	if p.choice == reduce {
		gotoLookupstr = "[" + gotoValueIndex1st + "," + valueofLHS + "]"
	}
	if p.notGrammatical == true {
		p.tStack = new(treeStack)
	}
	fmt.Printf("%-14s  %-14s [%2s,%2s]   %-6s  %-6s  %-7s %-13s   %-6s    %-5s  %-9s   %-20s", initpStack, initQue, actionIndex1st, inputQueFront, p.actionValue, valueofLHS, lengthofRHSstr, tempStack, gotoLookupstr, gotoValue, newPushstr, p.tStack.toString())
	fmt.Println()
}

func (p *parser) EvaluateActionChoice() {
	search := p.pStack.top().stateSym
	num, ero := strconv.ParseInt(search, 0, 64)
	if ero != nil {
		fmt.Println("ero here")
	}
	p.actionValue = p.aTable[num][p.aTableIndex[p.inputQueue.first()]]
	if compareStrings(p.actionValue, "accept") == true {
		p.choice = accept
	} else if compareStrings(p.actionValue, "") == true {
		p.choice = ungrammatical
		p.actionValue = "ungrammatical"
	} else if compareStrings("S", string(p.actionValue[0])) {
		p.choice = shift
	} else if compareStrings("R", string(p.actionValue[0])) {
		p.choice = reduce
	}
}

func compareStrings(a string, b string) bool {
	if a == b {
		return true
	}
	return false
}
