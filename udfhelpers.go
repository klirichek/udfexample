package main

/*
#include "sphinxudf.h"
#include <string.h>
#include <stdlib.h>
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
static void logmsg ( _GoString_ msg, sphinx_log_fn f )
{
	if (f)
		f ( _GoStringPtr(msg), _GoStringLen(msg) );
}
*/
import "C"
import (
	"reflect"
	"unsafe"
)

// Common constants for daemon and client.
const (
	/// current udf version
	SPH_UDF_VERSION = C.SPH_UDF_VERSION
)

/// UDF argument and result value types
const (
	SPH_UDF_TYPE_UINT32    = uint32(C.SPH_UDF_TYPE_UINT32)    ///< unsigned 32-bit integer
	SPH_UDF_TYPE_UINT32SET = uint32(C.SPH_UDF_TYPE_UINT32SET) ///< sorted set of unsigned 32-bit integers
	SPH_UDF_TYPE_INT64     = uint32(C.SPH_UDF_TYPE_INT64)     ///< signed 64-bit integer
	SPH_UDF_TYPE_FLOAT     = uint32(C.SPH_UDF_TYPE_FLOAT)     ///< single-precision IEEE 754 float
	SPH_UDF_TYPE_STRING    = uint32(C.SPH_UDF_TYPE_STRING)    ///< non-ASCIIZ string, with a separately stored length
	SPH_UDF_TYPE_INT64SET  = uint32(C.SPH_UDF_TYPE_INT64SET)  ///< sorted set of signed 64-bit integers
	SPH_UDF_TYPE_FACTORS   = uint32(C.SPH_UDF_TYPE_FACTORS)   ///< packed ranking factors
	SPH_UDF_TYPE_JSON      = uint32(C.SPH_UDF_TYPE_JSON)      ///< whole json or particular field as a string
)

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

// SPH_UDF_ARGS contain arguments passed to the function.
/*  -godefs shows this structure here:
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

func (args *SPH_UDF_ARGS) Arg_count() int32 {
	return int32(args.arg_count)
}

// internal: returns pointer go arg_type
func (args *SPH_UDF_ARGS) typeptr(idx int) *uint32 {
	return (*uint32)(unsafe.Pointer(uintptr(unsafe.Pointer(args.arg_types)) + unsafe.Sizeof(*args.arg_types)*uintptr(idx)))
}

// internal: returns unsafe pointer go arg_value
func (args *SPH_UDF_ARGS) valueptr(idx int) unsafe.Pointer {
	ptr:= unsafe.Pointer(uintptr(unsafe.Pointer(args.arg_values)) + unsafe.Sizeof(*args.arg_values)*uintptr(idx))
	return *(*unsafe.Pointer)(ptr)
}

// internal: returns unsafe pointer go arg_name
func (args *SPH_UDF_ARGS) nameptr(idx int) unsafe.Pointer {
	base:= uintptr(unsafe.Pointer(args.arg_names))
	if base==0 {
		return nil
	}

	ptr:= base + unsafe.Sizeof(*args.arg_names)*uintptr(idx)
	return *(*unsafe.Pointer)(unsafe.Pointer(ptr))
}

// internal: returns len of arg_value
func (args *SPH_UDF_ARGS) lenval(idx int) int {
	return int(*(*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(args.str_lengths)) + unsafe.Sizeof(*args.str_lengths)*uintptr(idx))))
}

// return name of the arg by idx
func (args *SPH_UDF_ARGS) arg_name(idx int) string {
	return C.GoString((*C.char)(args.nameptr(idx)))
}

// return type of the arg by idx
func (args *SPH_UDF_ARGS) arg_type(idx int) uint32 {
	return *args.typeptr(idx)
}

// return string value by idx
// it not copies value, but use backend C string instead
func (args *SPH_UDF_ARGS) stringval(idx int) string {
	var s string
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&s))
	hdr.Data = uintptr(args.valueptr(idx))
	hdr.Len = args.lenval(idx)
	return s
}

// return slice value by idx
// it not copies value, but use backend C string instead
func (args *SPH_UDF_ARGS) mva32(idx int) []uint32 {
	var mvas []uint32
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&mvas))
	sliceHeader.Cap = args.lenval(idx)
	sliceHeader.Len = args.lenval(idx)
	sliceHeader.Data = uintptr(args.valueptr(idx))
	return mvas
}

// return slice value by idx
// it not copies value, but use backend C string instead
func (args *SPH_UDF_ARGS) mva64(idx int) []int64 {
	var mvas []int64
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&mvas))
	sliceHeader.Cap = args.lenval(idx)
	sliceHeader.Len = args.lenval(idx)
	sliceHeader.Data = uintptr(args.valueptr(idx))
	return mvas
}

// convert Go string into C string result and return it
func (args *SPH_UDF_ARGS) return_string (result string) uintptr {
	return uintptr(unsafe.Pointer(C.retmsg(result,args.fn_malloc)))
}


/* UDF initialization
// -godefs shows this structure here:
type SPH_UDF_INIT struct {
        Func_data       *byte
        Is_const        int8
        Pad_cgo_0       [7]byte
}
 */
type SPH_UDF_INIT C.SPH_UDF_INIT

// set func_data to given value
// note that according CGO spec you can't provide any pointer, but only integer value,
// since gc manages all go pointers
func (init *SPH_UDF_INIT) setvalue(value uintptr) {
	*(*uintptr)(unsafe.Pointer(&init.func_data)) = value
}

func (init *SPH_UDF_INIT) setuint32(value uint32) {
	*(*uint32)(unsafe.Pointer(&init.func_data)) = value
}

// get func_data
func (init *SPH_UDF_INIT) getvalue() uintptr {
	return *(*uintptr)(unsafe.Pointer(&init.func_data))
}

func (init *SPH_UDF_INIT) getuint32() uint32 {
	return *(*uint32)(unsafe.Pointer(&init.func_data))
}

var cblog *C.sphinx_log_fn
func sphWarning (msg string) {
	C.logmsg(msg,cblog)
}

func strlen ( param *C.char ) int {
	return int(C.strlen(param))
}

func malloc ( param int ) unsafe.Pointer {
	return C.malloc((C.ulong)(param))
}

func free ( param unsafe.Pointer ) {
	C.free(param)
}
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

/// Reinit. Was called in workers=prefork, now it is just necessary stub.
//export udfexample_setlogcb
func udfexample_setlogcb(logfn *C.sphinx_log_fn) {
	cblog = logfn
}

func main() {
	// We need the main function to make possible
	// CGO compiler to compile the package as C shared library
}
