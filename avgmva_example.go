package main

import "C"

// avgmva is example of UDF function working with mva and returning float value
// You can load it into the daemon with
//  CREATE FUNCTION avgmva RETURNS FLOAT SONAME 'udfexample.so';

//export avgmva_init
func avgmva_init(init *SPH_UDF_INIT, args *SPH_UDF_ARGS, errmsg *ERR_MSG) int32 {

	// check argument count
	if args.arg_count != 1 ||
		(args.arg_type(0) != SPH_UDF_TYPE_UINT32SET && args.arg_type(0) != SPH_UDF_TYPE_INT64SET) {
		return errmsg.say("AVGMVA() requires 1 MVA argument")
	}

	// store our mva vs mva64 flag to func_data
	// (strictly speaking that is not necessary, since we always can fetch value from args,
	// but let it be for education)
	init.setuint32(uint32(args.arg_type(0)))
	return 0
}

// mva values stored in following form:
// args.lenval(idx) contains number of values in mva (despite size - 32 or 64 bits)
// args.valueptr contains pointer to actual values
// in order to iterate over mvas we cast raw c pointer into go slice, and then iterate
//export avgmva
func avgmva(init *SPH_UDF_INIT, args *SPH_UDF_ARGS, err *ERR_FLAG) float64 {

	result := float64(0.0)
	nvalues := args.lenval(0)

	if nvalues == 0 {
		return result
	}

	is64 := SPH_UDF_TYPE(init.getuint32())
	switch is64 {
	case SPH_UDF_TYPE_UINT32SET:
		{
			mvas := args.mva32(0)
			for _, value := range mvas {
				result += float64(value)
			}
		}
	case SPH_UDF_TYPE_INT64SET:
		{
			mvas := args.mva64(0)
			for _, value := range mvas {
				result += float64(value)
			}
		}
	}

	return result / float64(nvalues)
}
