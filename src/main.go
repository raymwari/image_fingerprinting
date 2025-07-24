package main

import (
	"bufio"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"log"
	"os"
)

var dct_mat []int16 // dct matrix
var q_tabs []q_tab  // quantization tables (luminance | chrominance)

// shomali11: https://github.com/shomali11/util/blob/master/xhashes/xhashes.go#L61
func stringHasher(algorithm hash.Hash, text string) string {
	algorithm.Write([]byte(text))
	return hex.EncodeToString(algorithm.Sum(nil))
}

func SHA512(text string) string {
	algorithm := sha512.New()
	return stringHasher(algorithm, text)
}

func up(arg string) string {
	f, _ := os.Open(arg)
	defer f.Close()
	st, _ := f.Stat()
	sz := st.Size()
	byts := make([]byte,
		sz)
	rd := bufio.NewReader(f)
	rd.Read(byts)
	h_tabs := make([]h_tab,
		4)
	// decoding the jpeg:
	huffman(byts, // 1. huffman decoding
		h_tabs)
	dct_mat = dequantize(dct_mat, // 2. dct matrix dequantization
		q_tabs)
	blks_x := 2
	blks_y := 1
	out := idct(dct_mat, // 3. inverse dct (luminance only)
		blks_x,
		blks_y)
	return SHA512(string(out))
}

func main() {
	if len(os.Args) < 3 {
		log.Fatalln("2 args expected... try again!")
	}
	args := os.Args
	if _, err := os.Stat(args[1]); errors.Is(err, os.ErrNotExist) {
		log.Fatalln("file doesn't exist:",
			args[1])
	}
	if _, err := os.Stat(args[2]); errors.Is(err, os.ErrNotExist) {
		log.Fatalln("file doesn't exist:",
			args[2])
	}
	if up(args[1]) == up(args[2]) {
		fmt.Println("Match: visually similar")
	} else {
		fmt.Println("Match: not similar")
	}

}
