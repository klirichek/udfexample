package main

// Common constants for daemon and client.
const (
	/// current udf version
	SPH_UDF_VERSION=9

	/// error buffer size
	SPH_UDF_ERROR_LEN=256
)

/// UDF argument and result value types
type SPH_UDF_TYPE int32

const (
	SPH_UDF_TYPE_UINT32 SPH_UDF_TYPE = iota+1	///< unsigned 32-bit integer
	SPH_UDF_TYPE_UINT32SET		///< sorted set of unsigned 32-bit integers
	SPH_UDF_TYPE_INT64			///< signed 64-bit integer
	SPH_UDF_TYPE_FLOAT			///< single-precision IEEE 754 float
	SPH_UDF_TYPE_STRING			///< non-ASCIIZ string, with a separately stored length
	SPH_UDF_TYPE_INT64SET		///< sorted set of signed 64-bit integers
	SPH_UDF_TYPE_FACTORS		///< packed ranking factors
	SPH_UDF_TYPE_JSON			///< whole json or particular field as a string
)



