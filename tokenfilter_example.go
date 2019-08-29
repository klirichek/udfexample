package main

import "C"
import (
	"fmt"
	"strings"
	"unsafe"
)

// This is index tokenizer filter example.
// It has to be declared in config as:
//  index_token_filter = udfexample.so:hideemail
// also to make it works, ensure your index config has min_prefix_len set (otherwise '*' search will not work),
// and also has appropriate characters (@, ., :) in charset_table. For example:
// ...
//  	charset_table = 0..9, english, _, @, ., :
//		min_prefix_len = 4
//  	rt_field = content
//		index_token_filter = udfexample.so:hideemail

// plugin is quite noisy - it will report into daemons' debug about everything it sees.
// That is not good for production, but really helpful for this very education purposes

// very simple email hider with exception
// symbols @ and . should be in charset_table
// it works on indexing by processing tokens as:
// - keep emails to domain 'space.io' and transform them to 'mailto:any@space.io' then searching for query with 'mailto:*' should return all documents with emails
// - deletes all other emails, ie returns NULL for 'test@gmail.com'

// this function called after indexer created token filter and also get the schema. So, we have actual field list here
// (however it is not necessary at all for our simple plugin)
//export hideemail_init
func hideemail_init(ppuserdata *uintptr, num_fields int32, field_names **C.char, options *C.char, errmsg *ERR_MSG) int32 {

	sphWarning(fmt.Sprintf("Called hideemail_init: %X, %d, %p, %p:%s, %p",
		*ppuserdata, num_fields, field_names, options, GoString(options), errmsg))

	// initialize storage in C memory. That is necessary, because functions from viewpoint of Go lang are standalone,
	// so if we just return a pointer to Go allocated structure, gc may accidentally kill it.
	// and also, if we make global var - don't forget, that functions must be thread-safe, so we need kind of another
	// storage (map), protected by mutexes. To avoid all this stuff we just use plain C 'malloc' and keep allocated
	// pointer in user data.
	if *ppuserdata == 0 {
		*ppuserdata = malloc(256)
	}
	return 0
}

// gets called once for each new token produced by base tokenizer with source token as its parameter.
// It must return token, count of extra tokens made by token filter and delta position for token.
// in case of hideemail we set extra to 0 (i.e. no more tokens), and delta to 1 (our only token in 1-st position)
//export hideemail_push_token
func hideemail_push_token(puserdata uintptr, token *C.char, extra *int32, delta *int32) *C.char {

	if token == nil {
		return token // sanity bypass
	}

	// we got C string, but want to operate with Go string. So, wrap it first
	rawtoken := GoString(token)
	sphWarning(fmt.Sprintf("Called hideemail_push_token with %s, %d, %d", rawtoken, *extra, *delta))

	parts := strings.Split(rawtoken, "@")
	if len(parts) != 2 {
		sphWarning("not email")
		return token // not email case
	}

	if parts[1] != "space.io" {
		sphWarning(fmt.Sprintf("not space.io, but '%s'", parts[1]))
		return nil // domain name does not match - hide email
	}
	result := (*C.char)(unsafe.Pointer(puserdata))

	// prefix and return
	sphWarning(fmt.Sprintf("returning 'mailto:%s' as result", rawtoken))
	putstr(result, "mailto:"+rawtoken)
	return result
}

// gets called multiple times in case hideemail_push_token reports extra tokens.
// It must return token and delta position for that extra token.
// For hideemail there are no extra tokens, so we just return null.
//export hideemail_get_extra_token
func hideemail_get_extra_token(puserdata uintptr, delta *int32) *C.char {
	sphWarning(fmt.Sprintf("Called hideemail_get_extra_token with %d", *delta))

	*delta = 0
	return nil
}

// this function is mandatory. It is called once for every RT insert/replace document.
// options are populated from the clause, like:
// INSERT INTO rt (id, title) VALUES (1, 'some text corp@space.io') OPTION token_filter_options='.io'
// - here we would get this '.io' in options string.
//export hideemail_begin_document
func hideemail_begin_document(puserdata uintptr, options *C.char, errmsg *ERR_MSG) int32 {
	return 0
}

// here we have to free allocated resources
//export hideemail_deinit
func hideemail_deinit(puserdata uintptr) {
	free(puserdata)
}
