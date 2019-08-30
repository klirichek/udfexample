package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// curl is example of stateful UDF function, just for fun.
// It takes 1 string param with url, and returns str contents of resource
// You can load it into the daemon with
//  CREATE FUNCTION curl RETURNS STRING SONAME 'udfexample.so';
// SELECT curl('https://yandex.ru/robots.txt');

// in this function we check only num of args (1) and the type (must be string)
// Also we set zero to init value
//export curl_init
func curl_init(_ *SPH_UDF_INIT, args *SPH_UDF_ARGS, errmsg *ERR_MSG) int32 {
	// check argument count
	if args.arg_count != 1 || args.arg_type(0) != SPH_UDF_TYPE_STRING {
		return errmsg.say("WEB() requires 1 string argument")
	}
	return 0
}

// here we execute provided action: extract arguments, make necessary calculations and return result back
//export curl
func curl(_ *SPH_UDF_INIT, args *SPH_UDF_ARGS, errf *ERR_FLAG) uintptr {

	url := args.stringval(0)

	// Get the data
	resp, _ := http.Get(url)
	defer func() { _ = resp.Body.Close() }()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return args.return_string(fmt.Sprintf("Bad status: %s", resp.Status))
	}

	// let's check content-type and avoid anything, but text
	contentType := resp.Header.Get("Content-Type")
	parts := strings.Split(contentType, "/")
	if len(parts) >= 1 && parts[0] == "text" {
		// retrieve whole body
		text, _ := ioutil.ReadAll(resp.Body)
		return args.return_string(string(text))
	}

	return args.return_string("Content type: " + contentType + ", will NOT download")
}
