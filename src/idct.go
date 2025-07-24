package main

import (
	"math"
)

func idct(dct_mat []int16,
	blks_x int,
	blks_y int) []byte {
	w := blks_x * 8
	h := blks_y * 8
	pix := make([]byte, w*h)
	for blk_id_x := 0; blk_id_x < blks_x*blks_y; blk_id_x++ {
		for pix_y := 0; pix_y < 8; pix_y++ {
			for pix_x := 0; pix_x < 8; pix_x++ {
				sum := 0.0
				for v := 0; v < 8; v++ {
					alpha_v := 1.0
					if v == 0 {
						alpha_v = 1.0 / math.Sqrt2
					}
					cos_yv := math.Cos(math.Pi * float64(v) * (float64(pix_y) + 0.5) / 8.0)
					for u := 0; u < 8; u++ {
						alpha_u := 1.0
						if u == 0 {
							alpha_u = 1.0 / math.Sqrt2
						}
						cos_xu := math.Cos(math.Pi * float64(u) * (float64(pix_x) + 0.5) / 8.0)
						coeff := float64(dct_mat[blk_id_x*64+v*8+u])
						partial := coeff * alpha_u * alpha_v * cos_xu * cos_yv
						sum += partial
					}
				}
				sum *= 0.25
				val := int(math.Round(sum + 128)) // pixel alignment (0 - 255)
				if val < 0 {
					val = 0
				}
				if val > 255 {
					val = 255
				}
				x := (blk_id_x%blks_x)*8 + pix_x
				y := (blk_id_x/blks_x)*8 + pix_y
				pix[y*w+x] = byte(val)
			}
		}
	}
	return pix
}
