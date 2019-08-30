package main

import "C"
import (
	"fmt"
	"strings"
)

// inspect is example of UDF function to inspect different arguments
// it is stateless, and just output info about provided args as string
// You can load it into the daemon with
//  CREATE FUNCTION inspect RETURNS STRING SONAME 'udfexample.so';
// most comprehensive usage example:
// SELECT inspect(id*1000,tags,1.145,tagl,'abc',WEIGHT(),jsn.data,jsn,packedfactors())
// 	FROM testall where match(' what | foo | what ') option ranker=expr('1')

// in this function we check only num of args (1 or 2) and the type of them (must be string)
// Nothing from init section is necessary and used.
//export inspect_init
func inspect_init(init *SPH_UDF_INIT, args *SPH_UDF_ARGS, errmsg *ERR_MSG) int32 {
	return 0
}

func (tp SPH_UDF_TYPE) String() string {
	switch tp {

	case SPH_UDF_TYPE_UINT32:
		return "UINT32"
	case SPH_UDF_TYPE_UINT32SET:
		return "UINT32SET"
	case SPH_UDF_TYPE_INT64:
		return "INT64"
	case SPH_UDF_TYPE_FLOAT:
		return "FLOAT"
	case SPH_UDF_TYPE_STRING:
		return "STRING"
	case SPH_UDF_TYPE_INT64SET:
		return "INT64SET"
	case SPH_UDF_TYPE_FACTORS:
		return "FACTORS"
	case SPH_UDF_TYPE_JSON:
		return "JSON"
	default:
		return fmt.Sprintf("unknown(%d)", tp)
	}
}

func format_mva32(values []uint32) string {
	var mvas []string
	for _, num := range values {
		mvas = append(mvas, fmt.Sprintf("%d", num))
	}
	return "(" + strings.Join(mvas, ",") + ")"
}

func format_mva64(values []int64) string {
	var mvas []string
	for _, num := range values {
		mvas = append(mvas, fmt.Sprintf("%d", num))
	}
	return "(" + strings.Join(mvas, ",") + ")"
}

// here we execute provided action: extract arguments,
// and make print details about each of them
//export inspect
func inspect(init *SPH_UDF_INIT, args *SPH_UDF_ARGS, err *ERR_FLAG) uintptr {
	var res string
	res += fmt.Sprintf("passed %d args\n", args.arg_count)
	for i := 0; i < int(args.arg_count); i++ {
		tp := args.arg_type(i)
		res += fmt.Sprintf("arg %d: type %v", i, tp)

		// print len, if necessary
		switch tp {
		case SPH_UDF_TYPE_UINT32SET, SPH_UDF_TYPE_STRING, SPH_UDF_TYPE_INT64SET, SPH_UDF_TYPE_JSON:
			res += fmt.Sprintf(", len %d", args.lenval(i))
		default:
		}
		res += ": "

		// output the value
		// here you can see how different types may be extracted from args
		switch tp {
		case SPH_UDF_TYPE_UINT32:
			res += fmt.Sprintf("%v", *(*uint32)(args.valueptr(i)))
		case SPH_UDF_TYPE_UINT32SET:
			res += format_mva32(args.mva32(i))
		case SPH_UDF_TYPE_INT64:
			res += fmt.Sprintf("%v", *(*int64)(args.valueptr(i)))
		case SPH_UDF_TYPE_FLOAT:
			res += fmt.Sprintf("%v", *(*float32)(args.valueptr(i)))
		case SPH_UDF_TYPE_STRING, SPH_UDF_TYPE_JSON:
			res += "'" + args.stringval(i) + "'"
		case SPH_UDF_TYPE_INT64SET:
			res += format_mva64(args.mva64(i))
		case SPH_UDF_TYPE_FACTORS:
			res += format_factors(args.valueptr(i))
		default:
			res += fmt.Sprintf("other (%d) todo!", tp)
		}

		res += "\n" // trailing cr
	}

	return args.return_string(res)
}
