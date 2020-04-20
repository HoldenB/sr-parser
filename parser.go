package main

import (
	"fmt"
	"strconv"
	"strings"
)

type inputBuffer struct {
	list []string
}

func (iq *inputBuffer) createEmptyInputBuffer() {
	iq.list = nil
}

func (iq *inputBuffer) createBuffer(list []string) {
	iq.list = list
	iq.list = append(iq.list, "$")
}

func (iq *inputBuffer) pop() string {
	a := iq.list[0]
	iq.list = iq.list[1:]
	return a
}

func (iq *inputBuffer) first() string {
	return iq.list[0]
}

func (iq *inputBuffer) repr() string {
	var sb strings.Builder
	for i := 0; i < len(iq.list); i++ {
		sb.WriteString(iq.list[i])
	}
	return sb.String()
}

type LHS_tree struct {
	sym string
}

func (lh *LHS_tree) initialize(a string) {
	lh.sym = a
}

func (lh LHS_tree) treeRepr(indent int) {
	a := ""
	for i := 0; i < indent; i++ {
		a = a + "    "
	}
	fmt.Println(a + lh.sym)
}

func (lh LHS_tree) repr() string {
	return lh.sym
}

type treeHandler struct {
	list []parseTreeInterface
}

func (handler *treeHandler) initialize(tree parseTreeInterface) {
	handler.list = append(handler.list, tree)
}

func (handler treeHandler) treeRepr(indent int) {
	for i := 0; i < len(handler.list); i++ {
		pt := handler.list[i]
		pt.treeRepr(indent)
	}
}

func (handler *treeHandler) repr() string {
	a := ""
	for _, e := range handler.list {
		b := e.repr()
		a += " " + b
	}
	return a
}

type NLT struct {
	parent   *LHS_tree
	children *treeHandler
}

func initialize(par *LHS_tree, handler *treeHandler, nlt *NLT) *NLT {
	nlt.parent = par
	nlt.children = handler
	return nlt
}

func (nlt NLT) printTree() {
	nlt.treeRepr(0)
}

func (nlt NLT) treeRepr(indent int) {
	nlt.parent.treeRepr(indent)
	nlt.children.treeRepr(indent + 1)
}

func (nlt NLT) repr() string {
	return "[" + nlt.parent.repr() + " " + nlt.children.repr() + "]"
}

type parseStack struct {
	list []entry
}

func (ps *parseStack) pop() entry {
	n := len(ps.list) - 1
	a := ps.list[n]

	ps.list[n] = entry{"", ""}
	ps.list = ps.list[:n]

	return a
}

func (ps *parseStack) push(a entry) {
	ps.list = append(ps.list, a)
}

func (ps *parseStack) isEmpty() bool {
	if len(ps.list) == 0 {
		return true
	}

	return false
}

func (ps *parseStack) top() entry {
	return ps.list[len(ps.list)-1]
}

func (ps *parseStack) repr() string {
	a := ""
	for i := 0; i < len(ps.list); i++ {
		a = a + ps.list[i].repr()
	}
	return a
}

type parseTreeInterface interface {
	printTree()
	treeRepr(indent int)
	repr() string
}

type entry struct {
	grammer string
	state   string
}

func (se *entry) createEntry(a string, b string) {
	se.grammer = b
	se.state = a
}

func (se *entry) repr() string {
	return se.grammer + se.state
}

type termSymbol struct {
	symbol string
}

func (ts termSymbol) printTree() {
	ts.treeRepr(0)
}

func (ts termSymbol) treeRepr(indent int) {
	a := ""
	// Indent value spaces
	for i := 0; i < indent; i++ {
		a = a + "    "
	}
	fmt.Println(a + ts.symbol)
}

func createtermSymbolbol(a string, ts *termSymbol) *termSymbol {
	ts.symbol = a
	return ts
}

func (ts termSymbol) repr() string {
	return ts.symbol
}

type actionChoice int

const (
	accept        actionChoice = 0
	ungrammatical actionChoice = 1
	shift         actionChoice = 2
	reduce        actionChoice = 3
)

type treeStack struct {
	list []parseTreeInterface
}

func (tree *treeStack) pop() parseTreeInterface {
	n := len(tree.list) - 1
	a := tree.list[n]
	tree.list = tree.list[:n]
	return a
}

func (tree *treeStack) push(e parseTreeInterface) {
	tree.list = append(tree.list, e)
}

func (tree *treeStack) top() parseTreeInterface {
	return tree.list[len(tree.list)-1]
}

func (tree *treeStack) repr() string {
	a := ""
	// Tree empty spaces
	for _, e := range tree.list {
		b := e.repr()
		a = b + " " + a
	}
	return a
}

func printHeader() {
	fmt.Println("                input          action    action  value   length  temp            goto      goto   stack       ")
	fmt.Println("stack           tokens         lookup    value   of LHS  of rhs  stack           lookup    value  action      parse tree stack")
	fmt.Println("______________________________________________________________________________________________________________________________")
}

func cmpStr(a string, b string) bool {
	if a == b {
		return true
	}
	return false
}

type parser struct {
	choice           actionChoice
	ungrammatical    bool
	actionValue      string
	parserStack      *parseStack
	inputBuffer      *inputBuffer
	newlyPlacedEntry *entry
	ts               *treeStack
	input            []string
	outputTable      [][]string
	grammar          [][]string
	grammarTable     [][]string
	actionTable      [][]string
	actionMap        map[string]int
	grammarMap       map[string]int
}

func (p *parser) Initialize(
	input []string,
	grammar [][]string,
	actionTable [][]string,
	gammarTable [][]string,
	actionMap map[string]int,
	grammarMap map[string]int) {

	// Output the input value
	p.input = input
	fmt.Print("\n\n")
	fmt.Print("Input: ", p.input, "\n\n\n")

	// Initialization
	p.grammar = grammar
	p.actionTable = actionTable
	p.grammarTable = gammarTable
	p.actionMap = actionMap
	p.grammarMap = grammarMap

	// Initialize nil values
	p.inputBuffer = nil
	p.parserStack = nil
	p.newlyPlacedEntry = nil
	p.ts = nil

	p.ungrammatical = false
	p.choice = ungrammatical
	p.actionValue = ""

	// Input buffer initialization
	p.outputTable = [][]string{}
	p.inputBuffer = new(inputBuffer)
	p.inputBuffer.createBuffer(p.input)

	p.newlyPlacedEntry = new(entry)
	p.newlyPlacedEntry.createEntry("0", "")
	p.parserStack = new(parseStack)
}

func (p *parser) PrintTree() {
	fmt.Println("Parse Tree: ")
	p.ts.top().printTree()
}

func (p *parser) EvaluateActionChoice() {
	search := p.parserStack.top().state

	num, err := strconv.ParseInt(search, 0, 64)
	if err != nil {
		fmt.Println("ERROR")
	}

	// Get the action value from the action table with the first in the input buffer
	p.actionValue = p.actionTable[num][p.actionMap[p.inputBuffer.first()]]

	// Compare the action value to see if we should accept/or if the value
	// is ungrammatical
	if cmpStr(p.actionValue, "accept") == true {
		p.choice = accept
	} else if cmpStr(p.actionValue, "") == true {
		p.choice = ungrammatical
		p.actionValue = "ungrammatical"
	} else if cmpStr("S", string(p.actionValue[0])) {
		p.choice = shift
	} else if cmpStr("R", string(p.actionValue[0])) {
		p.choice = reduce
	}
}

func (p *parser) Parse() {
	printHeader()

	p.ts = new(treeStack)
	p.Parse1step()

	// While loop to parse as long as we can either
	// reduce or shift
	for {
		if p.choice != reduce && p.choice != shift {
			break
		} else {
			p.Parse1step()
		}
	}

	if p.ungrammatical == false {
		fmt.Println()
		p.PrintTree()
	}
}

func (p *parser) Parse1step() {
	// LHS/RHS temp variables
	LHSvalue := ""
	LHSindexValue := -1

	RHSlength := 0
	RHSlengthstr := ""
	//////////////////////////

	// Default str initializers
	gotoValue := ""
	temparserStack := ""
	gotoValueIndex1st := ""
	newlyPlacedEntrystr := ""
	gotoLookupstr := ""
	//////////////////////////

	initQue := p.inputBuffer.repr()
	inputQueFront := p.inputBuffer.first()
	p.parserStack.push(*(p.newlyPlacedEntry))
	initparserStack := p.parserStack.repr()

	p.EvaluateActionChoice()
	actionIndexFirst := p.parserStack.top().state

	// Switch on our choice and evaluate the case
	switch p.choice {
	// Choice is accepted
	case accept:
		break

	// Choice is ungrammatical
	case ungrammatical:
		p.ungrammatical = true

	// We need to shift to the next input buffer value
	case shift:
		a := []rune(p.actionValue)

		index := string(a[1])
		p.newlyPlacedEntry = new(entry)
		p.newlyPlacedEntry.createEntry(index, p.inputBuffer.first())

		if cmpStr(p.inputBuffer.first(), "id") == true {
			b := new(termSymbol)
			b = createtermSymbolbol("id", b)
			p.ts.push(*b)
		}
		p.inputBuffer.pop()

	// We need to reduce our input action and pop off the queue
	// if we have a successful reduction
	case reduce:
		u := []rune(p.actionValue)
		h := string(u[1])
		q, err := strconv.Atoi(h)
		if err != nil {
			fmt.Println("ERROR")
		}

		// LHS
		LHSindexValue = q
		LHSvalue = p.grammar[LHSindexValue-1][0]
		RHSlength = len(p.grammar[LHSindexValue-1]) - 2
		popped := new(parseStack)

		// Iterate RHS length
		for i := 1; i <= RHSlength; i++ {
			popped.push(p.parserStack.pop())
		}

		temparserStack = p.parserStack.repr()
		gotoValueIndex1st = p.parserStack.top().state
		y := []rune(gotoValueIndex1st)
		m := string(y[0])

		n, e := strconv.Atoi(m)
		if e != nil {
			fmt.Println("ERROR")
		}

		// New entry
		gotoValue = p.grammarTable[n][p.grammarMap[LHSvalue]]
		p.newlyPlacedEntry = new(entry)
		p.newlyPlacedEntry.createEntry(gotoValue, LHSvalue)

		if RHSlength == 1 {
			old := p.ts.pop()
			child := new(treeHandler)
			child.initialize(old)

			LHSsymbol := new(LHS_tree)
			LHSsymbol.initialize(LHSvalue)

			tr := new(NLT)
			tr = initialize(LHSsymbol, child, tr)

			p.ts.push(tr)

			// Evaluate the case where we have a left paren
		} else if RHSlength == 3 && cmpStr(popped.top().grammer, "(") {
			fmt.Println(popped.top().grammer)

			a := p.ts.pop()
			b := new(termSymbol)
			b = createtermSymbolbol("(", b)
			c := new(termSymbol)
			c = createtermSymbolbol(")", c)

			handler := new(treeHandler)
			handler.initialize(b)
			handler.list = append(handler.list, a)
			handler.list = append(handler.list, c)

			parent := new(LHS_tree)
			parent.initialize(LHSvalue)

			nlt := new(NLT)
			nlt = initialize(parent, handler, nlt)
			p.ts.push(nlt)

			// No left paren
		} else if RHSlength == 3 {
			c := p.ts.pop()
			a := p.ts.pop()

			popped.pop()
			operator := popped.top().grammer

			b := new(termSymbol)
			b = createtermSymbolbol(operator, b)

			handler := new(treeHandler)
			handler.initialize(a)
			handler.list = append(handler.list, b)
			handler.list = append(handler.list, c)

			parent := new(LHS_tree)
			parent.initialize(LHSvalue)

			tr := new(NLT)
			tr = initialize(parent, handler, tr)
			p.ts.push(tr)
		}
	}

	if p.choice == shift || p.choice == reduce {
		newlyPlacedEntrystr = "push " +
			p.newlyPlacedEntry.repr()
	}

	if RHSlength > 0 {
		RHSlengthstr = strconv.Itoa(RHSlength)
	}

	// If our choice is to reduce we need to add the lookup str
	// and formatted output
	if p.choice == reduce {
		gotoLookupstr = "[" +
			gotoValueIndex1st +
			"," +
			LHSvalue +
			"]"
	}

	// Is our string ungrammatical?
	if p.ungrammatical == true {
		p.ts = new(treeStack)
	}

	// Format output after we're done with this iteration of parsing
	fmt.Printf("%-14s  %-14s [%2s,%2s]   %-6s  %-6s  %-7s %-13s   %-6s    %-5s  %-9s   %-20s",
		initparserStack,
		initQue,
		actionIndexFirst,
		inputQueFront,
		p.actionValue,
		LHSvalue,
		RHSlengthstr,
		temparserStack,
		gotoLookupstr,
		gotoValue,
		newlyPlacedEntrystr,
		p.ts.repr())

	fmt.Println()
}
