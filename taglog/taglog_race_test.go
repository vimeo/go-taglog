// +build race

package taglog

import (
	"bytes"
	"fmt"
	"sync"
	"testing"
)

func TestRace(t *testing.T) {
	n := 100
	buf := &bytes.Buffer{}
	buf.Grow(1 << 20)
	var wg sync.WaitGroup

	lg := New(buf, "", 0)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < n; i++ {
			lg.Println("line")
		}
	}()

	for i := 0; i < 10; i++ {
		wg.Add(1)
		l := lg.Copy()
		l.SetPrefix(fmt.Sprintf("Copy%d ", i+1))
		go func(l *Logger) {
			defer wg.Done()
			for i := 0; i < n; i++ {
				l.Println("line")
			}
		}(l)
	}

	wg.Wait()

	t.Logf("%db written", len(buf.Bytes()))
}
