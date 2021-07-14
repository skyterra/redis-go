package rgo

import "time"

func NowMs() uint64 {
	return uint64(time.Now().UnixNano() / 1e6)
}
