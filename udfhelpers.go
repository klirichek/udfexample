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
static void cmsg (char * sDst, _GoString_ s)
{
	size_t n = _GoStringLen(s);
	if ( n>SPH_UDF_ERROR_LEN-1 )
		n = SPH_UDF_ERROR_LEN-1;
	strncpy ( sDst, (const char*) _GoStringPtr(s), n);
	sDst[n] = '\0';
}
static char* retmsg ( _GoString_ msg, sphinx_malloc_fn f )
{
	size_t iLen = _GoStringLen(msg);
	char* sRes = f(iLen+1);
	strncpy ( sRes, (const char*) _GoStringPtr(msg), iLen);
	sRes[iLen] = '\0';
	return sRes;
}
*/
import "C"
import "unsafe"

// Common constants for daemon and client.
const (
	/// current udf version
	SPH_UDF_VERSION = C.SPH_UDF_VERSION
)

/// UDF argument and result value types
const (
	SPH_UDF_TYPE_UINT32    = C.SPH_UDF_TYPE_UINT32    ///< unsigned 32-bit integer
	SPH_UDF_TYPE_UINT32SET = C.SPH_UDF_TYPE_UINT32SET ///< sorted set of unsigned 32-bit integers
	SPH_UDF_TYPE_INT64     = C.SPH_UDF_TYPE_INT64     ///< signed 64-bit integer
	SPH_UDF_TYPE_FLOAT     = C.SPH_UDF_TYPE_FLOAT     ///< single-precision IEEE 754 float
	SPH_UDF_TYPE_STRING    = C.SPH_UDF_TYPE_STRING    ///< non-ASCIIZ string, with a separately stored length
	SPH_UDF_TYPE_INT64SET  = C.SPH_UDF_TYPE_INT64SET  ///< sorted set of signed 64-bit integers
	SPH_UDF_TYPE_FACTORS   = C.SPH_UDF_TYPE_FACTORS   ///< packed ranking factors
	SPH_UDF_TYPE_JSON      = C.SPH_UDF_TYPE_JSON      ///< whole json or particular field as a string
)

// SPH_UDF_ARGS contain arguments passed to the function. -godefs shows this structure here:
/*
type SPH_UDF_ARGS struct {
        Arg_count       int32
        Arg_types       *uint32
        Arg_values      **int8
        Arg_names       **int8
        Str_lengths     *int32
        Fn_malloc       *[0]byte
}
*/
type SPH_UDF_ARGS C.SPH_UDF_ARGS

// ERR_MSG is the buffer for returning error messages from _init functions.
type ERR_MSG C.char

// Report packs message from go string into C string buffer to be returned from UDF function.
// It returns 1 in order to be used as shortcut (i.e. 'return msg.Report(...)' instead of 'msg.Report(...); return 1'
func (errmsg *ERR_MSG) say(message string) int32 {
	C.cmsg((*C.char)(errmsg), message)
	return 1
}

// ERR_FLAG points to success flag and may be used to indicate critical errors
type ERR_FLAG C.char

// fail set ERR_FLAG to 1
func (errflag *ERR_FLAG) fail() {
	*errflag = 1
}

func get(arr unsafe.Pointer, idx int) int32 {
	ptr := uintptr(arr) + uintptr(4*idx)
	return *(*int32)(unsafe.Pointer(ptr))
}

type SPH_UDF_INIT C.SPH_UDF_INIT
type SPH_UDF_FIELD_FACTORS C.SPH_UDF_FIELD_FACTORS

// global functions that must be in any udf plugin library

/// UDF version control. Named as LIBRARYNAME_ver (i.e. udfexample_ver in the case)
/// gets called once when the library is loaded
//export udfexample_ver
func udfexample_ver() int32 {
	return SPH_UDF_VERSION
}

/// Reinit. Was called in workers=prefork, now it is just necessary stub.
//export udfexample_reinit
func udfexample_reinit() {
}

func main() {
	// We need the main function to make possible
	// CGO compiler to compile the package as C shared library
}
