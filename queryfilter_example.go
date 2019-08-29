package main

import "C"
import (
	"fmt"
	"reflect"
	"unsafe"
)

// This is query tokenizer filter example.
// plugin is loaded per-query, you have to invoke it directly as option.
// For example for query:
// select * from ru where match ('сталь');
// with filter you have to write something as
// select * from ru where match ('сталь') OPTION token_filter='udfexample.so:queryshow:bla';

// Plugin  is intentionally noisy - it will report into daemons' debug about everything it sees.
// That is not good for production, but really helpful for this very education purposes

// this function gets called once per index prior to parsing query with parameters -
// max token length and string set by token_filter option
//export queryshow_init
func queryshow_init(ppuserdata *uintptr, max_len int32, options *C.char, errmsg *ERR_MSG) int32 {

	sphWarning(fmt.Sprintf("Called queryshow_init: %X, max_len %d, options %s, err %p",
		*ppuserdata, max_len, C.GoString(options), errmsg))
	return 0
}

// this function gets called once for token right before it got passed to morphology processor
// with reference to token and stopword flag. It might set stopword flag to mark token as stopword.
//export queryshow_pre_morph
func queryshow_pre_morph(puserdata uintptr, token *C.char, stopword *int32) {

	// we got C string, but want to operate with Go string. So, wrap it first
	var rawtoken string
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&rawtoken))
	hdr.Data = uintptr(unsafe.Pointer(token))
	hdr.Len = strlen(token)

	// let's for example return stopword for russian word "сталь"
	if rawtoken == "сталь" {
		*stopword = 1
	}

	sphWarning(fmt.Sprintf("%X Called queryshow_pre_morph: token '%s', res %d",
		puserdata, rawtoken, *stopword))
}

// this function gets called once for token after it processed by morphology processor
// with reference to token and stopword flag. It might set stopword flag to mark token as stopword.
// It must return flag non-zero value of which means to use token prior to morphology processing.
//export queryshow_post_morph
func queryshow_post_morph(puserdata uintptr, token *C.char, stopword *int32) int32 {

	// we got C string, but want to operate with Go string. So, wrap it first
	var rawtoken string
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&rawtoken))
	hdr.Data = uintptr(unsafe.Pointer(token))
	hdr.Len = strlen(token)

	ires := int32(0)
	// let's for example return stopword for russian word "сталь"
	if rawtoken == "сталь" {
		*stopword = 1
	}

	// let's signal to use non-morphed for word "стальная"
	if rawtoken == "стальная" {
		ires = 1
	}

	sphWarning(fmt.Sprintf("%X Called queryshow_post_morph: token '%s', stop %d, result %d",
		puserdata, rawtoken, *stopword, ires))
	return ires
}

// this function gets called once for each new token produced by base tokenizer with parameters:
// token produced by base tokenizer, pointer to raw token at source query string and raw token length.
// It must return token and delta position for token.
//export queryshow_push_token
func queryshow_push_token(puserdata uintptr, token *C.char, delta *int32, rawtoken *C.char, rawtokenlen int32) *C.char {

	// we got C string, but want to operate with Go string. So, wrap it first
	var crawtoken string
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&crawtoken))
	hdr.Data = uintptr(unsafe.Pointer(rawtoken))
	hdr.Len = int(rawtokenlen)

	var ctoken string
	hdr1 := (*reflect.StringHeader)(unsafe.Pointer(&ctoken))
	hdr1.Data = uintptr(unsafe.Pointer(token))
	hdr1.Len = strlen(token)

	sphWarning(fmt.Sprintf("%X Called queryshow_push_token: token '%s', rawtoken '%s'",
		puserdata, ctoken, crawtoken))

	return token
}
