package internal

import "sync"

var InfoPool = sync.Pool{
	New: func() any {
		return map[string]any{}
	},
}
