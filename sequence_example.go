package main

import "C"
import (
	"fmt"
	"unsafe"
)

/// UDF initialization
/// gets called on every query, when query begins
/// args are filled with values for a particular query
//export sequence_init
func sequence_init(init *SPH_UDF_INIT, args *SPH_UDF_ARGS, errmsg *ERR_MSG) int32 {

	// check argument count
	if args.arg_count > 1 {
		errmsg.say(fmt.Sprintln("SEQUENCE() takes either 0 or 1 arguments"))
		return 1
	}

	// check argument type
	if args.arg_count == 1 && (*args.arg_types) != SPH_UDF_TYPE_UINT32 {

		errmsg.say(fmt.Sprintln("SEQUENCE() requires 1st argument to be uint"))
		return 1
	}

	*(*int32)(unsafe.Pointer(&init.func_data)) = 1

	return 0
}

/// UDF deinitialization
/// gets called on every query, when query ends
//export sequence_deinit
func sequence_deinit(init *SPH_UDF_INIT) {
	*(*int32)(unsafe.Pointer(&init.func_data)) = 0
}

/// UDF implementation
/// gets called for every row, unless optimized away
//export sequence
func sequence(init *SPH_UDF_INIT, args *SPH_UDF_ARGS, err *ERR_FLAG) int64 {
	pres := (*int32)(unsafe.Pointer(&init.func_data))
	*pres++
	if args.arg_count > 0 {
		*pres = *pres + get(unsafe.Pointer(*args.arg_values), 0)
	}
	return int64(*pres)
}
