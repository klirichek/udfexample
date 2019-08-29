package main

import "C"
import (
	"reflect"
	"strings"
	"unsafe"
)

// This is index tokenizer filter example.
// It has to be declared in config as:
//  index_token_filter = udfexample.so:hideemail

// very simple email hider with exception
// symbols @ and . should be in charset_table
// it works on indexing by processing tokens as:
// - keep emails to domain 'space.io' and transform them to 'mailto:any@space.io' then searching for query with 'mailto:*' should return all documents with emails
// - deletes all other emails, ie returns NULL for 'test@gmail.com'

// this function called after indexer created token filter and also get the schema. So, we have actual field list here
// (however it is not necessary at all for our simple plugin)
//export hideemail_init
func hideemail_init(ppuserdata **byte, num_fields int32, field_names **C.char, options *C.char, errmsg *ERR_MSG) int32 {

	sphWarning("Called hideemail_init")

	// initialize storage in C memory. That is necessary, because functions from viewpoint of Go lang are standalone,
	// so if we just return a pointer to Go allocated structure, gc may accidentally kill it.
	// and also, if we make global var - don't forget, that functions must be thread-safe, so we need kind of another
	// storage (map), protected by mutexes. To avoid all this stuff we just use plain C 'malloc' and keep allocated
	// pointer in user data
	*(*uintptr)(unsafe.Pointer(*ppuserdata)) = uintptr(malloc(256))
	return 0
}

// gets called once for each new token produced by base tokenizer with source token as its parameter.
// It must return token, count of extra tokens made by token filter and delta position for token.
// in case of hideemail we set extra to 0 (i.e. no more tokens), and delta to 1 (our only token in 1-st position)
//export hideemail_push_token
func hideemail_push_token(puserdata *byte, token *C.char, extra *int32, delta *int32) *C.char {

	if token == nil {
		return token // sanity bypass
	}

	// we got C string, but want to operate with Go string. So, wrap it first
	var rawtoken string
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&rawtoken))
	hdr.Data = uintptr(unsafe.Pointer(token))
	hdr.Len = strlen(token)

	parts := strings.Split(rawtoken, "@")
	if len(parts) != 2 {
		return token // not email case
	}

	if parts[1] != "space.io" {
		return nil // domain name does not match - hide email
	}
	result := (*C.char)(unsafe.Pointer(puserdata))

	// prefix and return
	putstr(result, "mailto:"+rawtoken)
	return result
}

// gets called multiple times in case hideemail_push_token reports extra tokens.
// It must return token and delta position for that extra token.
// For hideemail there are no extra tokens, so we just return null.
//export hideemail_get_extra_token
func hideemail_get_extra_token(puserdata *byte, delta *int32) *C.char {
	*delta = 0
	return nil
}

// this function is mandatory. It is called once for every RT insert/replace document.
// options are populated from the clause, like:
// INSERT INTO rt (id, title) VALUES (1, 'some text corp@space.io') OPTION token_filter_options='.io'
// - here we would get this '.io' in options string.
//export hideemail_begin_document
func hideemail_begin_document(puserdata *byte, options *C.char, errmsg *ERR_MSG) int32 {
	return 0
}

// here we have to free allocated resources
//export hideemail_deinit
func hideemail_deinit(puserdata *byte) {
	free(unsafe.Pointer(puserdata))
}
