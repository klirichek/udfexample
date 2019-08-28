package main

import "C"
import (
	"fmt"
	"unsafe"
)

//export avgmva_init
func avgmva_init(init *SPH_UDF_INIT, args *SPH_UDF_ARGS, errmsg *ERR_MSG) int32 {

	// check argument count
	if args.arg_count != 1 ||
		((*args.arg_types) != SPH_UDF_TYPE_UINT32SET && (*args.arg_types) != SPH_UDF_TYPE_INT64SET) {
		errmsg.say(fmt.Sprintln("AVGMVA() requires 1 MVA argument"))
		return 1
	}

	// check argument type
	if args.arg_count == 1 && (*args.arg_types) != SPH_UDF_TYPE_UINT32 {
		errmsg.say(fmt.Sprintln("SEQUENCE() requires 1st argument to be uint"))
		return 1
	}
	// store our mva vs mva64 flag to func_data
	if (*args.arg_types) == SPH_UDF_TYPE_INT64SET {
		*(*int32)(unsafe.Pointer(&init.func_data)) = 1
	} else {
		*(*int32)(unsafe.Pointer(&init.func_data)) = 0
	}
	return 0
}
