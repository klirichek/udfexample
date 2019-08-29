package main

import "C"
import (
	"fmt"
	"reflect"
	"unsafe"
)

// very simple email hider with exception
// symbols @ and . should be in charset_table
// it works on indexing by processing tokens as:
// - keep emails to domain 'space.io' and transform them to 'mailto:any@space.io' then searching for query with 'mailto:*' should return all documents with emails
// - deletes all other emails, ie returns NULL for 'test@gmail.com'

//export hideemail_init
func hideemail_init ( ppuserdata **byte, num_fields int32, field_names **C.char, options *C.char, errmsg *ERR_MSG ) int32 {
sphWarning ( "Called hideemail_init" )
	*(*uintptr) (unsafe.Pointer(*ppuserdata)) = uintptr( malloc ( 256 ))
return 0
}

func hideemail_push_token ( puserdata *byte, token *C.char, extra *int32, delta *int32) *C.char {
	var s string
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&s))
	hdr.Data = uintptr(unsafe.Pointer(token))
	hdr.Len = strlen(token)
	sphWarning (fmt.Sprintf("Got token %s", s))
	return nil
}

func hideemail_get_extra_token ( puserdata *byte, delta *int32 ) * C.char {
	*delta = 0
	return nil
}

func hideemail_begin_document ( puserdata *byte, delta *int32 ) * C.char {
	return nil
}

func hideemail_deinit ( puserdata *byte ) {
	free ( unsafe.Pointer(puserdata) )
}