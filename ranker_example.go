package main

import "C"
import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

func (ranker *SPH_RANKER_INIT) String() string {

	var weights []int32
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&weights))
	sliceHeader.Cap = int(ranker.num_field_weights)
	sliceHeader.Len = int(ranker.num_field_weights)
	sliceHeader.Data = uintptr(unsafe.Pointer(ranker.field_weights))

	var sweights []string
	for _, num := range weights {
		sweights = append(sweights, fmt.Sprintf("%d", num))
	}
	line := "weights: [" + strings.Join(sweights, ",") + "], "
	line += fmt.Sprintf("options: %s, payload_mask: %d, num_query_words: %d, max_qpos: %d",
		GoString(ranker.options), ranker.payload_mask, ranker.num_query_words, ranker.max_qpos)
	return line
}

// This is ranker plugin example.
// plugin is loaded per-query, you have to invoke it directly as option.
// For example for query:
// select * from ru where match ('сталь');
// with ranker you have to write something as
// CREATE PLUGIN myrank TYPE 'ranker' SONAME 'udfexample.so';
// SELECT * from ru WHERE match ('сталь') OPTION ranker=myrank('option1=1');

// Plugin  is intentionally noisy - it will report into daemons' debug about everything it sees.
// That is not good for production, but really helpful for this very education purposes

// this function gets called once per query per index, in the very beginning.
// A few query-wide options are passed to it through a SPH_RANKER_INIT structure,
/*
type SPH_RANKER_INIT struct {
        Num_field_weights       int32
        Field_weights           *int32
        Options                 *int8
        Payload_mask            uint32
        Num_query_words         int32
        Max_qpos                int32
}
*/
// including the user options strings (in the example just above, “option1=1” is that string).
//export myrank_init
func myrank_init(ppuserdata *uintptr, ranker *SPH_RANKER_INIT, errmsg *ERR_MSG) int32 {

	sphWarning(fmt.Sprintf("Called myrank_init with %v", ranker))
	return 0
}

func (hit *SPH_RANKER_HIT) String() string {
	return fmt.Sprintf("doc_id=%d, hit_pos=%d, query_pos=%d, node_pos=%d, span_length=%d, match_length=%d, weight=%d, query_pos_mask=%d",
		hit.doc_id, hit.hit_pos, hit.query_pos, hit.node_pos, hit.span_length, hit.match_length,
		hit.weight, hit.query_pos_mask)
}

// this function gets called multiple times per matched document, with every matched keyword
// occurrence passed as its parameter, a SPH_RANKER_HIT structure.
/*
type SPH_RANKER_HIT struct {
        Doc_id          uint64
        Hit_pos         uint32
        Query_pos       uint16
        Node_pos        uint16
        Span_length     uint16
        Match_length    uint16
        Weight          uint32
        Query_pos_mask  uint32
}
*/
// The occurrences within each document are guaranteed to be passed in the order
// of ascending hit->hit_pos values.
//export myrank_update
func myrank_update(puserdata uintptr, hit *SPH_RANKER_HIT) {
	sphWarning(fmt.Sprintf("Called myrank_update with %v", hit))
}

// this function gets called once per matched document, once there are no more keyword occurrences.
// It must return the WEIGHT() value. This is the only mandatory function.
//export myrank_finalize
func myrank_finalize(puserdata uintptr, match_weight int32) uint32 {

	sphWarning(fmt.Sprintf("myrank_finalize called with %d and will return 1000", match_weight))
	return 1000
}

// this function gets called once per query, in the very end.
//export myrank_deinit
func myrank_deinit(puserdata uintptr) int32 {
	sphWarning("myrank_deinit called and return 0")
	return 0
}
