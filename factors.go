package main

import (
	"fmt"
	"unsafe"
)

type dword_ptr uintptr

func (in dword_ptr) atp(n int) unsafe.Pointer {
	return unsafe.Pointer(in + dword_ptr(n*4))
}

func (in dword_ptr) at(n int) uint32 {
	return *(*uint32)(in.atp(n))
}

func (in dword_ptr) ati(n int) int {
	return int(in.at(n))
}

func (in dword_ptr) atf(n int) float32 {
	return *(*float32)(in.atp(n))
}

func (in *dword_ptr) add(n int) {
	*in = *in + dword_ptr(n*4)
}

func (in *dword_ptr) bit(n int) int {
	if *(*uint32)(in.atp(n/32)) & ( 1 << ( uint32(n % 32) ) ) == 0 {
		return 0
	}
	return 1
}

func skip_fields(in dword_ptr, n int) dword_ptr {
	in.add(6 + ((in.ati(5)+31)/32)*2) // skip heading document factors and 2 exact masks
	for i := 0; i < n; i++ {
		if in.ati(0) > 0 { // skip 15 ints per matched field, or 1 per unmatched
			in.add(15)
		} else {
			in.add(1)
		}
	}
	return in
}

func skip_terms(in dword_ptr, n int) dword_ptr {
	in.add(1) // skip max_uniq_qpos
	for i := 0; i < n; i++ {
		if in.ati(0) > 0 { // skip 4 ints per matched term, or 1 per unmatched
			in.add(4)
		} else {
			in.add(1)
		}
	}
	return in
}

/// returns a pointer to the field factors, or NULL for a non-matched field index
func sphinx_get_field_factors(in dword_ptr, field int) dword_ptr {
	if in == 0 || field < 0 || field > in.ati(5) {
		return 0 // blob[5] is num_fields, do a sanity check
	}
	in = skip_fields(in, field)
	if in.ati(0) == 0 {
		return 0 // no hits, no fun
	}
	if in.ati(1) != field {
		return 0 // field[1] is field_id, do a sanity check
	}
	return in
}

/// returns a pointer to the term factors, or NULL for a non-matched field index
func sphinx_get_term_factors(in dword_ptr, term int) dword_ptr {

	if in == 0 || term < 0 {
		return 0
	}
	in = skip_fields(in, in.ati(5)) // skip all fields
	if term > in.ati(0) {
		return 0 // sanity check vs max_uniq_qpos ( qpos and terms range - [1, max_uniq_qpos]
	}
	in = skip_terms(in, term-1)
	if in.ati(0) == 0 {
		return 0 // unmatched term
	}
	if in.ati(1) != term {
		return 0 // term[1] is keyword_id, sanity check failed
	}
	return in
}

type sphinx_doc_factor int

const (
	SPH_DOCF_BM25             sphinx_doc_factor = iota + 1 ///< int
	SPH_DOCF_BM25A                                         ///< float
	SPH_DOCF_MATCHED_FIELDS                                ///< unsigned int
	SPH_DOCF_DOC_WORD_COUNT                                ///< int
	SPH_DOCF_NUM_FIELDS                                    ///< int
	SPH_DOCF_MAX_UNIQ_QPOS                                 ///< int
	SPH_DOCF_EXACT_HIT_MASK                                ///< unsigned int
	SPH_DOCF_EXACT_ORDER_MASK                              ///< v.4, unsigned int
)

/// returns a document factor value, interpreted as integer
func sphinx_get_doc_factor_int(in dword_ptr, f sphinx_doc_factor) int {
	switch f {
	case SPH_DOCF_BM25:
		return in.ati(1)
	case SPH_DOCF_BM25A:
		return in.ati(2)
	case SPH_DOCF_MATCHED_FIELDS:
		return in.ati(3)
	case SPH_DOCF_DOC_WORD_COUNT:
		return in.ati(4)
	case SPH_DOCF_NUM_FIELDS:
		return in.ati(5)
	case SPH_DOCF_MAX_UNIQ_QPOS:
		in = skip_fields(in, in.ati(5))
		return in.ati(0)
	case SPH_DOCF_EXACT_HIT_MASK:
		return in.ati(6)
	case SPH_DOCF_EXACT_ORDER_MASK:
		fields_size := (in.ati(5) + 31) / 32
		return in.ati(6 + fields_size)
	}
	return 0
}

/// returns a document factor value, interpreted as float
func sphinx_get_doc_factor_float(in dword_ptr, f sphinx_doc_factor) float32 {
	switch f {
	case SPH_DOCF_BM25A:
		return in.atf(2)
	}
	return 0.0
}

/// returns a pointer to document factor value, interpreted as vector of integers
func sphinx_get_doc_factor_ptr(in dword_ptr, f sphinx_doc_factor) dword_ptr {
	switch f {
	case SPH_DOCF_EXACT_HIT_MASK:
		in.add(6)
		return in
	case SPH_DOCF_EXACT_ORDER_MASK:
		fields_size := (in.ati(5) + 31) / 32
		in.add(6 + fields_size)
		return in
	}
	return 0
}

type sphinx_field_factor int

const (
	SPH_FIELDF_HIT_COUNT         sphinx_field_factor = iota + 1 ///< int
	SPH_FIELDF_LCS                                              ///< unsigned int
	SPH_FIELDF_WORD_COUNT                                       ///< unsigned int
	SPH_FIELDF_TF_IDF                                           ///< float
	SPH_FIELDF_MIN_IDF                                          ///< float
	SPH_FIELDF_MAX_IDF                                          ///< float
	SPH_FIELDF_SUM_IDF                                          ///< float
	SPH_FIELDF_MIN_HIT_POS                                      ///< int
	SPH_FIELDF_MIN_BEST_SPAN_POS                                ///< int
	SPH_FIELDF_MAX_WINDOW_HITS                                  ///< int
	SPH_FIELDF_MIN_GAPS                                         ///< v.3, int
	SPH_FIELDF_ATC                                              ///< v.4, float
	SPH_FIELDF_LCCS                                             ///< v.5, int
	SPH_FIELDF_WLCCS                                            ///< v.5, float
)

/// returns a field factor value, interpreted as integer
func sphinx_get_field_factor_int(in dword_ptr, f sphinx_field_factor) int {
	if in == 0 {
		return 0
	}
	switch f {
	case SPH_FIELDF_HIT_COUNT:
		return in.ati(0)
	case SPH_FIELDF_LCS:
		return in.ati(2)
	case SPH_FIELDF_WORD_COUNT:
		return in.ati(3)
	case SPH_FIELDF_TF_IDF:
		return in.ati(4)
	case SPH_FIELDF_MIN_IDF:
		return in.ati(5)
	case SPH_FIELDF_MAX_IDF:
		return in.ati(6)
	case SPH_FIELDF_SUM_IDF:
		return in.ati(7)
	case SPH_FIELDF_MIN_HIT_POS:
		return in.ati(8)
	case SPH_FIELDF_MIN_BEST_SPAN_POS:
		return in.ati(9)
	case SPH_FIELDF_MAX_WINDOW_HITS:
		return in.ati(10)
	case SPH_FIELDF_MIN_GAPS:
		return in.ati(11)
	case SPH_FIELDF_ATC:
		return in.ati(12)
	case SPH_FIELDF_LCCS:
		return in.ati(13)
	case SPH_FIELDF_WLCCS:
		return in.ati(14)
	}
	return 0
}

/// returns a field factor value, interpreted as float
func sphinx_get_field_factor_float(in dword_ptr, f sphinx_field_factor) float32 {
	r := sphinx_get_field_factor_int(in, f)
	return *(*float32)(unsafe.Pointer(&r))
}

type sphinx_term_factor int

const (
	SPH_TERMF_KEYWORD_MASK sphinx_term_factor = iota + 1 ///< unsigned int
	SPH_TERMF_TF                                         ///< int
	SPH_TERMF_IDF                                        ///< float
)

/// returns a term factor value, interpreted as integer
func sphinx_get_term_factor_int(in dword_ptr, f sphinx_term_factor) int {
	if in == 0 {
		return 0
	}
	switch f {
	case SPH_TERMF_KEYWORD_MASK:
		return in.ati(0)
	case SPH_TERMF_TF:
		return in.ati(2)
	case SPH_TERMF_IDF:
		return in.ati(3)
	}
	return 0
}

/// returns a term factor value, interpreted as float
func sphinx_get_term_factor_float(in dword_ptr, f sphinx_term_factor) float32 {
	r := sphinx_get_term_factor_int(in, f)
	return *(*float32)(unsafe.Pointer(&r))
}


// format provided factors as string similar to packedfactors() in original output
func format_factors(ptr unsafe.Pointer) string {
	in := dword_ptr(ptr)
	if in==0 {
		return "null factors"
	}
	sBmFmt := "bm25=%d, bm25a=%f, field_mask=%d, doc_word_count=%d"
	sFieldFmt := "field%d=(lcs=%d, hit_count=%d, word_count=%d, tf_idf=%f, min_idf=%f, max_idf=%f, sum_idf=%f, min_hit_pos=%d, min_best_span_pos=%d, exact_hit=%d, max_window_hits=%d, min_gaps=%d, exact_order=%d, lccs=%d, wlccs=%f, atc=%f)"
	sWordFmt := "word%d=(tf=%d, idf=%f)"
	var res string
	res = fmt.Sprintf (sBmFmt, sphinx_get_doc_factor_int (in, SPH_DOCF_BM25),
		sphinx_get_doc_factor_float (in, SPH_DOCF_BM25A),
		sphinx_get_doc_factor_int (in, SPH_DOCF_MATCHED_FIELDS),
		sphinx_get_doc_factor_int (in, SPH_DOCF_DOC_WORD_COUNT))
	pExactHit := sphinx_get_doc_factor_ptr ( in, SPH_DOCF_EXACT_HIT_MASK )
	pExactOrder := sphinx_get_doc_factor_ptr ( in, SPH_DOCF_EXACT_ORDER_MASK )
	iFields := sphinx_get_doc_factor_int (in, SPH_DOCF_NUM_FIELDS)
	for i:=0; i<iFields; i++ {
		pField := sphinx_get_field_factors(in, i)
		if sphinx_get_field_factor_int ( pField, SPH_FIELDF_HIT_COUNT)==0 {
			continue
		}
		if res!="" {
			res += ", "
		}
		res += fmt.Sprintf ( sFieldFmt,
			i, sphinx_get_field_factor_int ( pField, SPH_FIELDF_LCS ),
			sphinx_get_field_factor_int ( pField, SPH_FIELDF_HIT_COUNT ),
			sphinx_get_field_factor_int ( pField, SPH_FIELDF_WORD_COUNT ),
			sphinx_get_field_factor_float ( pField, SPH_FIELDF_TF_IDF ),
			sphinx_get_field_factor_float ( pField, SPH_FIELDF_MIN_IDF ),
			sphinx_get_field_factor_float ( pField, SPH_FIELDF_MAX_IDF ),
			sphinx_get_field_factor_float ( pField, SPH_FIELDF_SUM_IDF ),
			sphinx_get_field_factor_int ( pField, SPH_FIELDF_MIN_HIT_POS ),
			sphinx_get_field_factor_int ( pField, SPH_FIELDF_MIN_BEST_SPAN_POS ),
			pExactHit.bit(i),
			sphinx_get_field_factor_int ( pField, SPH_FIELDF_MAX_WINDOW_HITS ),
			sphinx_get_field_factor_int ( pField, SPH_FIELDF_MIN_GAPS ),
			pExactOrder.bit(i),
			sphinx_get_field_factor_int ( pField, SPH_FIELDF_LCCS ),
			sphinx_get_field_factor_float ( pField, SPH_FIELDF_WLCCS ),
			sphinx_get_field_factor_float ( pField, SPH_FIELDF_ATC ))

	}

	iUniqQpos := sphinx_get_doc_factor_int (in, SPH_DOCF_MAX_UNIQ_QPOS )
	for i:=0; i<iUniqQpos; i++ {
		pTerm := sphinx_get_term_factors ( in, i+1 )
		if sphinx_get_term_factor_int ( pTerm, SPH_TERMF_KEYWORD_MASK) == 0 {
			continue
		}
		if res!="" {
			res += ", "
		}
		res += fmt.Sprintf ( sWordFmt, i, sphinx_get_term_factor_int ( pTerm, SPH_TERMF_TF ),
			sphinx_get_term_factor_float ( pTerm, SPH_TERMF_IDF ) )
	}
	return res
}