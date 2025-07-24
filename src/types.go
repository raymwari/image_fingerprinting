package main

// q tables:
type q_tab struct {
	data []byte
}

// huffman:
type tab_ref struct {
	c1 []int // [component 1, dc/ac]
	c2 []int // [component 2, dc/ac]
	c3 []int // [component 3, dc/ac]
}

type pair struct {
	len byte
	sym byte
}

type h_tab struct {
	t_type   string
	sym_t    int      // total no of symbols
	sym_lens []byte   // array of symbol lengths
	syms     []byte   // actual symbol bytes
	codes    []string // codes for table
}

// headers:
var mk byte = 0xff
var sos byte = 0xda // start of scan
var dri byte = 0xdd // restart interval
var dqt byte = 0xdb // quant tables
