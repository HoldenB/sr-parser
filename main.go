package main

func main() {
	// Test data
	input := []string{"id", "+", "id", "*", "id"}
	// input := []string{"id"} // simplest possible grammatical input
	// input := []string{ "id", "+", "id"};                       // simple grammatical input
	// input := []string{ "id", "+", "id", "+", "id"};            // left assoc
	// input := []string{ "(", "id", ")"};                        // parens1
	// input := []string{"id", "+", "(", "id", "+", "id", ")"} // parens2
	// input := []string{"id", "+", "id", "*"} // ungrammatical input

	gammer := [][]string{
		{"E", "->", "E", "+", "T"},
		{"E", "->", "T"},
		{"T", "->", "T", "*", "F"},
		{"T", "->", "F"},
		{"F", "->", "(", "E", ")"},
		{"F", "->", "id"},
	}

	actionTable := [][]string{
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

	gammarTable := [][]string{
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

	actionMap := map[string]int{
		"id": 0,
		"+":  1,
		"*":  2,
		"(":  3,
		")":  4,
		"$":  5,
	}

	grammarMap := map[string]int{
		"E": 0,
		"T": 1,
		"F": 2,
	}

	p := new(parser)
	p.Initialize(input, gammer, actionTable, gammarTable, actionMap, grammarMap)
	p.Parse()
}
