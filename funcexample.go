package main

/*
#include "sphinxudf.h"
#include <string.h>
static void malloc_fn ( sphinx_malloc_fn f, int iSize)
{
	f(iSize);
}
size_t _GoStringLen(_GoString_ s);
const char *_GoStringPtr(_GoString_ s);
static void errmsg (char * sDst, _GoString_ s)
{
	size_t n = _GoStringLen(s);
	if ( n>SPH_UDF_ERROR_LEN )
		n = SPH_UDF_ERROR_LEN;
	strncpy ( sDst, (const char*) _GoStringPtr(s), n);
	sDst[n] = '\0';
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

var stored_malloc *C.sphinx_malloc_fn;
func malloc ( size int32) {
	C.malloc_fn (stored_malloc, C.int(size))
}

type SPH_UDF_INIT C.struct_st_sphinx_udf_init
type SPH_UDF_ARGS C.struct_st_sphinx_udf_args


//export udfexample_ver
func udfexample_ver() int {
	return SPH_UDF_VERSION
}

//export udfexample_reinit
func udfexample_reinit() {
	fmt.Println("udfreinit called")
}

//export strtoint_init
func strtoint_init(init *SPH_UDF_INIT, args *SPH_UDF_ARGS, errmsg *C.char) int32 {
	if args.arg_count!=1 || (*args.arg_types)!=uint32(SPH_UDF_TYPE_STRING) {
		C.errmsg(errmsg,fmt.Sprintln ("STRTOINT() requires 1 string argument"))
		return 1
	}
	return 0
}

//export strtoint
func strtoint(init *SPH_UDF_INIT, args *SPH_UDF_ARGS, errmsg *C.char) int64 {
	bts:= C.GoStringN(*args.arg_values, *args.str_lengths)
	var a int64
	_, _ = fmt.Sscanf(bts, "%X", &a)
	return a
}

func get(arr unsafe.Pointer, idx int) int32 {
	ptr := uintptr(arr) + uintptr ( 4*idx )
	return *(*int32)(unsafe.Pointer(ptr))
}

/// UDF initialization
/// gets called on every query, when query begins
/// args are filled with values for a particular query
//export sequence_init
func sequence_init(init *SPH_UDF_INIT, args *SPH_UDF_ARGS, errmsg *C.char) int32 {

	// check argument count
	if args.arg_count>1 {
		C.errmsg(errmsg,fmt.Sprintln ("SEQUENCE() takes either 0 or 1 arguments"))
		return 1
	}

	// check argument type
	if args.arg_count==1 && (*args.arg_types)!=uint32(SPH_UDF_TYPE_UINT32) {
		C.errmsg(errmsg,fmt.Sprintln ("SEQUENCE() requires 1st argument to be uint"))
		return 1
	}

	*(*int32)(unsafe.Pointer(&init.func_data)) = 1

	return 0
}

/// UDF deinitialization
/// gets called on every query, when query ends
//export sequence_deinit
func sequence_deinit(init *SPH_UDF_INIT ) {
	*(*int32)(unsafe.Pointer(&init.func_data)) = 0
}

/// UDF implementation
/// gets called for every row, unless optimized away
//export sequence
func sequence(init *SPH_UDF_INIT, args *SPH_UDF_ARGS, errmsg *C.char) int64 {
	pres := (*int32)(unsafe.Pointer(&init.func_data))
	*pres++
	if args.arg_count>0 {
		*pres = *pres + get(unsafe.Pointer(*args.arg_values),0)
	}
	return int64(*pres)
}

//export avgmva_init
func avgmva_init(init *SPH_UDF_INIT, args *SPH_UDF_ARGS, errmsg *C.char) int32 {

	// check argument count
	if args.arg_count!=1 ||
		((*args.arg_types)!=uint32(SPH_UDF_TYPE_UINT32SET) && (*args.arg_types)!=uint32(SPH_UDF_TYPE_INT64SET)) {
		C.errmsg(errmsg,fmt.Sprintln ("AVGMVA() requires 1 MVA argument"))
		return 1
	}

	// check argument type
	if args.arg_count==1 && (*args.arg_types)!=uint32(SPH_UDF_TYPE_UINT32) {
		C.errmsg(errmsg,fmt.Sprintln ("SEQUENCE() requires 1st argument to be uint"))
		return 1
	}
	// store our mva vs mva64 flag to func_data
	if (*args.arg_types)==uint32(SPH_UDF_TYPE_INT64SET) {
		*(*int32)(unsafe.Pointer(&init.func_data)) = 1
	} else {
		*(*int32)(unsafe.Pointer(&init.func_data)) = 0
	}
	return 0
}

func main() {
	// We need the main function to make possible
	// CGO compiler to compile the package as C shared library
}