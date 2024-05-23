package utils

import (
	"strings"
	"sync"
)

func BoolToIntString(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func BoolToStr(b bool, trueStr, falseStr string) string {
	if b {
		return trueStr
	}
	return falseStr
}

type TextPacket struct {
	Data sync.Map
}

func (T *TextPacket) Parse(text string) {
	lines := strings.Split(text, "\n")
	wg := sync.WaitGroup{}
	for _, line := range lines {
		wg.Add(1)
		go func(line string) {
			defer wg.Done()
			parts := strings.Split(line, "|")
			if len(parts) != 2 {
				return
			}
			T.Data.Store(parts[0], parts[1])
		}(line)
	}
	wg.Wait()
}

func (T *TextPacket) HasKey(key string) bool {
	_, ok := T.Data.Load(key)
	return ok
}

func (T *TextPacket) GetFromKey(key string) string {
	val, ok := T.Data.Load(key)
	if !ok {
		return ""
	}
	return val.(string)
}
