package main

import "C"

// sequence is example of UDF function with state
// You can load it into the daemon with
//  CREATE FUNCTION sequence RETURNS INT SONAME 'udfexample.so';

/// UDF initialization
/// gets called on every query, when query begins
/// args are filled with values for a particular query
/// here you have to fill the state. On C you may allocate resources, etc.
/// On Go we can't save any allocations since gc is in game and it will just lose the pointer
/// so, we set initial value just to POD value 1.
//export sequence_init
func sequence_init(init *SPH_UDF_INIT, args *SPH_UDF_ARGS, errmsg *ERR_MSG) int32 {

	// check argument count
	if args.arg_count > 1 {
		return errmsg.say("SEQUENCE() takes either 0 or 1 arguments")
	}

	// check argument type
	if args.arg_count == 1 && args.arg_type(0) != SPH_UDF_TYPE_UINT32 {
		return errmsg.say("SEQUENCE() requires 1st argument to be uint")
	}

	init.setuint32(1)
	return 0
}

/// UDF deinitialization
/// gets called on every query, when query ends
/// here you have to reset the state. We have nothing to do actually, but for educational purpose
/// we 'release' resource by setting value to 0.
//export sequence_deinit
func sequence_deinit(init *SPH_UDF_INIT) {
	init.setuint32(0)
}

/// UDF implementation
/// gets called for every row, unless optimized away.
/// here we take and modify stored state value, and use it in calculation of the result.
//export sequence
func sequence(init *SPH_UDF_INIT, args *SPH_UDF_ARGS, err *ERR_FLAG) int64 {
	res := init.getuint32()
	init.setuint32(res + 1)
	if args.arg_count > 0 {
		res = res + *(*uint32)(args.valueptr(0))
	}
	return int64(res)
}
