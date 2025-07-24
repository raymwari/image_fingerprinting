package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

func h_decoder(stream string,
	tabs []h_tab,
	int_mk int,
	refs tab_ref) []int16 {
	var pred_y,
		pred_cb,
		pred_cr int
	var out []int16
	mcu_c := 0
	pos := 0
	nxt_bit := func() (byte,
		bool) {
		if pos >= len(stream) {
			return 0, false
		}
		b := stream[pos]
		pos++
		return b,
			true
	}
	nxt_nbits := func(n int) (int,
		bool) {
		val := 0
		for i := 0; i < n; i++ {
			b, ok := nxt_bit()
			if !ok {
				return 0, false
			}
			val = (val << 1) | int(b-'0')
		}
		return val,
			true
	}
	decode_signed := func(bits int,
		size int) int {
		if size == 0 {
			return 0
		}
		half := 1 << (size - 1)
		if bits < half {
			return bits - ((1 << size) - 1)
		}
		return bits
	}
	// actual decoding:
	huff_decode := func(table h_tab) (byte,
		bool) {
		code := ""
		for {
			b, ok := nxt_bit()
			if !ok {
				return 0,
					false
			}
			code += string(b)
			for i, c := range table.codes {
				if c == code {
					return table.syms[i], true
				}
			}
		}
	}
	// mapping tables:
	resolve := func(ref_val int,
		lum_dc,
		lum_ac,
		chr_dc,
		chr_ac h_tab) (h_tab,
		h_tab) {
		dc_id := (ref_val >> 4) & 0xF
		ac_id := ref_val & 0xF
		if dc_id == 0 {
			if ac_id == 0 {
				return lum_dc, lum_ac
			}
			return lum_dc,
				chr_ac // unused
		}
		if ac_id == 0 {
			return chr_dc,
				lum_ac // unused
		}
		return chr_dc,
			chr_ac
	}
	// fetch tables:
	var lum_dc,
		lum_ac,
		chr_dc,
		chr_ac h_tab
	for _, t := range tabs {
		switch t.t_type {
		case "dl":
			lum_dc = t
		case "al":
			lum_ac = t
		case "dc":
			chr_dc = t
		case "ac":
			chr_ac = t
		}
	}
	// decode single block:
	decode_block := func(dc_tab h_tab,
		ac_tab h_tab,
		predictor *int) ([]int16,
		bool) {
		block := make([]int16, 64)
		// dc:
		dc_sym, ok := huff_decode(dc_tab)
		if !ok {
			return nil, false
		}
		size := int(dc_sym)
		bits, _ := nxt_nbits(size)
		dc_val := decode_signed(bits, size) + *predictor
		*predictor = dc_val
		block[0] = int16(dc_val)
		// ac:
		idx := 1
		for idx < 64 {
			ac_sym, ok := huff_decode(ac_tab)
			if !ok {
				break
			}
			if ac_sym == 0x00 { // eob
				break
			}
			if ac_sym == 0xF0 { // zrl
				idx += 16
				continue
			}
			run := int(ac_sym >> 4)
			size := int(ac_sym & 0x0F)
			idx += run
			if idx >= 64 {
				break
			}
			val_bits, _ := nxt_nbits(size)
			val := decode_signed(val_bits,
				size)
			block[idx] = int16(val)
			idx++
		}
		return block,
			true
	}
	// component tables:
	dcy, acy := resolve(refs.c1[1],
		lum_dc,
		lum_ac,
		chr_dc,
		chr_ac)
	dc_cb, ac_cb := resolve(refs.c2[1],
		lum_dc,
		lum_ac,
		chr_dc,
		chr_ac)
	dc_cr, ac_cr := resolve(refs.c3[1],
		lum_dc,
		lum_ac,
		chr_dc,
		chr_ac)
	// mcu loop:
	for {
		// 4 y blocks:
		for i := 0; i < 4; i++ {
			block, ok := decode_block(dcy,
				acy,
				&pred_y)
			if !ok {
				return out
			}
			out = append(out,
				block...)
		}
		// cb block:
		block, ok := decode_block(dc_cb,
			ac_cb,
			&pred_cb)
		if !ok {
			return out
		}
		out = append(out, block...)
		// cr block:
		block, ok = decode_block(dc_cr,
			ac_cr, &pred_cr)
		if !ok {
			return out
		}
		out = append(out, block...)
		// restart interval:
		mcu_c++
		if int_mk > 0 && mcu_c%int_mk == 0 {
			pred_y,
				pred_cb,
				pred_cr = 0,
				0,
				0
		}
		// stream ended:
		if pos >= len(stream) {
			break
		}
	}
	return out
}

func code_fitter(syms []byte,
	lens []byte) []string {
	var pairs []pair
	sym_ind := 0
	for len := 1; len <= 16; len++ {
		count := int(lens[len-1])
		for i := 0; i < count; i++ {
			pairs = append(pairs,
				pair{byte(len),
					syms[sym_ind]})
			sym_ind++
		}
	}
	code := 0
	var codes []string
	prev_len := byte(0)
	for _, p := range pairs {
		if p.len != prev_len {
			code <<= (p.len - prev_len)
			prev_len = p.len
		}
		c := fmt.Sprintf("%0*b",
			p.len,
			code)
		codes = append(codes,
			c)
		code++
	}
	return codes
}

func h_codes(h_tabs []h_tab) {
	for i := 0; i < len(h_tabs); i++ {
		if h_tabs[i].t_type == "in" {
			log.Fatalln("encountered invalid huffman table!")
		}
		codes := code_fitter(h_tabs[i].syms,
			h_tabs[i].sym_lens)
		h_tabs[i].codes = codes
	}
}

func h_dec(byts []byte,
	h_tabs []h_tab) {
	h_codes(h_tabs)
	mlen_off := 3
	cmp_off := 4
	rest_off := 5
	var cmp int   // number of components
	var m_len int // length of current marker (sos)
	rest_int := 0 // restart interval
	var data_off int
	var ref tab_ref
	for i := 0; i < len(byts); i++ {
		if byts[i] == mk {
			if byts[i+1] == sos {
				data_off = i + 1
				cmp = int(byts[i+cmp_off])
				ref.c1 = append(ref.c1,
					int(byts[i+cmp_off+1]))
				ref.c1 = append(ref.c1,
					int(byts[i+cmp_off+2]))
				ref.c2 = append(ref.c2,
					int(byts[i+cmp_off+3]))
				ref.c2 = append(ref.c2,
					int(byts[i+cmp_off+4]))
				ref.c3 = append(ref.c3,
					int(byts[i+cmp_off+5]))
				ref.c3 = append(ref.c3,
					int(byts[i+cmp_off+6]))
				m_len = int(byts[i+mlen_off])
				data_off = data_off + m_len
				if cmp > 3 {
					log.Fatalln("encountered specialized encoding! (code locked for baseline encoded jpegs)")
				}
			}
			if byts[i+1] == dri {
				rest_int = int(byts[i+rest_off])
			}
			go extract_qtabs(byts,
				i)
		}
	}
	data_off++ // padding
	byts = byts[data_off:]
	var bin []string
	var bin_j string
	for i := 0; i < len(byts); i++ {
		b := strconv.FormatInt(int64(byts[i]),
			2)
		bin = append(bin,
			b)
		bin_j = strings.Join(bin,
			"")
	}
	dct_mat = h_decoder(bin_j,
		h_tabs,
		rest_int,
		ref)
}

func huffman(byts []byte,
	h_tabs []h_tab) {
	var ht byte = 0xc4 // define huffman table
	id_offset := 4
	lens_offset := 5
	syms_offset := 21
	max_syms := 162 // 16 (max_lens) * 10 (max_values)
	t_cn := 0
	for i := 0; i < len(byts); i++ {
		if byts[i] == mk {
			if byts[i+1] == ht {
				id := strconv.FormatInt(int64(byts[i+id_offset]),
					16)
				switch id {
				case "0": // 00
					h_tabs[t_cn].t_type = "dl"
				case "1": // 01
					h_tabs[t_cn].t_type = "dc"
				case "10":
					h_tabs[t_cn].t_type = "al"
				case "11":
					h_tabs[t_cn].t_type = "ac"
				default:
					h_tabs[t_cn].t_type = "in"
				}
				lens := make([]byte,
					16)
				for j := 0; j < 16; j++ {
					lens[j] = byts[i+lens_offset+j]
				}
				h_tabs[t_cn].sym_lens = lens
				sym_lens := int(lens[0])
				for k := 1; k < len(lens); k++ {
					sym_lens = sym_lens + int(lens[k])
				}
				syms := make([]byte,
					max_syms)
				for j := 0; j < sym_lens; j++ {
					syms[j] = byts[i+syms_offset+j]
				}
				h_tabs[t_cn].sym_t = sym_lens
				h_tabs[t_cn].syms = syms
				t_cn++
			}
		}
	}
	h_dec(byts,
		h_tabs)
}
