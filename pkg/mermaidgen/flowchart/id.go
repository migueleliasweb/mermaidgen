package flowchart

import (
	"fmt"
	"strconv"
	"sync/atomic"
)

// nodesCount helps generate easy unique IDs for nodes when necessary
var nodesCount = uint64(0)

func generateID() string {
	id := fmt.Sprintf("id%s", strconv.Itoa(int(nodesCount)))
	atomic.AddUint64(&nodesCount, 1)

	return id
}
