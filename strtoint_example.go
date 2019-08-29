package main

import "fmt"

// strtoint is example of stateless UDF function. It takes 1 or 2 params and returns integer from hexadecimal string
// or given Go formatter
// You can load it into the daemon with
//  CREATE FUNCTION strtoint RETURNS INT SONAME 'udfexample.so';
// SELECT strtoint('10000'); // 65536, hexadecimal default
// SELECT strtoint('10000', '%b'); // 16, binary
// SELECT strtoint('10000', '%o'); // 4096, octal

// in this function we check only num of args (1 or 2) and the type of them (must be string)
// Nothing from init section is necessary and used.
//export strtoint_init
func strtoint_init(init *SPH_UDF_INIT, args *SPH_UDF_ARGS, errmsg *ERR_MSG) int32 {
	switch args.arg_count {
	case 1:
		if args.arg_type(0)!=SPH_UDF_TYPE_STRING {
		return errmsg.say("STRTOINT() requires 1st string argument")
	}
	case 2: if args.arg_type(1)!=SPH_UDF_TYPE_STRING {
		return errmsg.say("STRTOINT() requires 2nd string argument")
	}
	default: return errmsg.say("STRTOINT() requires 1 or 2 string arguments")
	}
	return 0
}

// here we execute provided action: extract arguments, make necessary calculations and return result back
//export strtoint
func strtoint(init *SPH_UDF_INIT, args *SPH_UDF_ARGS, err *ERR_FLAG) int64 {
	fmtstr := "%X"
	if args.arg_count==2 {
		fmtstr = args.stringval(1)
	}
	var a int64
	_, _ = fmt.Sscanf(args.stringval(0), fmtstr, &a)
	return a
}
