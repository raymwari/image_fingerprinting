package main

import (
	"log"
)

func dequantize(dct_mat []int16,
	q_tabs []q_tab) []int16 {
	var lum_tab,
		chrom_tab [64]uint8
	for _, qt := range q_tabs {
		if len(qt.data) < 65 {
			log.Fatalln("encountered quantization table length!")
		}
		switch qt.data[0] {
		case 0:
			copy(lum_tab[:],
				qt.data[1:65])
		case 1:
			copy(chrom_tab[:],
				qt.data[1:65])
		default:
			log.Fatalln("encountered unknown quantization table ID")
		}
	}
	out := make([]int16,
		len(dct_mat))
	num_blks := len(dct_mat) / 64
	for blk_num := 0; blk_num < num_blks; blk_num++ {
		var qt [64]uint8
		if blk_num%6 < 4 {
			qt = lum_tab
		} else {
			qt = chrom_tab
		}
		strt_idx := blk_num * 64
		for i := 0; i < 64; i++ {
			out[strt_idx+i] = dct_mat[strt_idx+i] * int16(qt[i])
		}
	}
	return out
}

func extract_qtabs(byts []byte,
	i int) {
	if byts[i+1] == dqt {
		if byts[i+3] != 132 {
			log.Fatalln("re-encode image(s) and try again (view assumptions section)")
		}
		off_set := i + 3
		off_set++ // padding
		q_tabs = make([]q_tab,
			2) // luminance | chrominance
		q_tabs[0].data = byts[off_set:(off_set + 65)]
		q_tabs[1].data = byts[(off_set + 65):(off_set + 130)]
	}
}
