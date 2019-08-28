package main

import "C"
import "fmt"

// strtoint is very simple UDF function. It takes 1 or 2 params and returns integer from hexadecimal string
// or given Go formatter
// You can load it into the daemon with
//  CREATE FUNCTION strtoint RETURNS INT SONAME 'udfexample.so';

// in this function we check only num of args (1 or 2) and the type of them (must be string)
// Nothing from init section is necessary and used.
//export strtoint_init
func strtoint_init(init *SPH_UDF_INIT, args *SPH_UDF_ARGS, errmsg *ERR_MSG) int32 {
	if args.arg_count != 1 || (*args.arg_types) != SPH_UDF_TYPE_STRING {
		return errmsg.say(fmt.Sprintln("STRTOINT() requires 1 string argument"))
	}
	return 0
}

// here we execute provided action: extract arguments, make necessary calculations and return result back
//export strtoint
func strtoint(init *SPH_UDF_INIT, args *SPH_UDF_ARGS, err *ERR_FLAG) int64 {
	bts := C.GoStringN(*args.arg_values, *args.str_lengths)
	var a int64
	_, _ = fmt.Sscanf(bts, "%X", &a)
	err.fail()
	return a
}
