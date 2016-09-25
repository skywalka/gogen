package timer

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/coccyx/gogen/internal"
	"github.com/stretchr/testify/assert"
)

func TestTimer(t *testing.T) {
	os.Setenv("GOGEN_HOME", "..")
	c := config.NewConfig()
	s := c.FindSampleByName("translog")
	gq := make(chan *config.GenQueueItem)
	oq := make(chan *config.OutQueueItem)

	timer := &Timer{S: s, GQ: gq, OQ: oq}
	go timer.NewTimer()

	item := <-gq

	// Test that we get a GenQueueItem
	var gqi *config.GenQueueItem
	assert.Equal(t, reflect.TypeOf(gqi), reflect.ValueOf(item).Type())

	// Test that we're about the same interval
	n := time.Now()
	timer = &Timer{S: s, GQ: gq, OQ: oq}
	go timer.NewTimer()
	item = <-gq
	cur := time.Now()

	gt := cur.Sub(n) > (time.Duration(s.Interval) * time.Second)
	assert.Equal(t, true, gt)
}

func TestBackfill(t *testing.T) {
	os.Setenv("GOGEN_HOME", "..")
	os.Setenv("GOGEN_ALWAYS_REFRESH", "1")
	home := ".."
	os.Setenv("GOGEN_SAMPLES_DIR", filepath.Join(home, "config", "tests", "tokens.yml"))

	s := config.FindSampleInFile(home, "backfill")

	gq := make(chan *config.GenQueueItem, 1000)
	oq := make(chan *config.OutQueueItem)
	done := make(chan int)
	gqs := make([]*config.GenQueueItem, 0, 10)

	timer := &Timer{S: s, GQ: gq, OQ: oq, Done: done}
	go timer.NewTimer()
	<-done
Loop:
	for {
		select {
		case i := <-gq:
			gqs = append(gqs, i)
		default:
			break Loop
		}
	}
	assert.Equal(t, 6, len(gqs))
}

func TestBackfillRealtime(t *testing.T) {
	os.Setenv("GOGEN_HOME", "..")
	os.Setenv("GOGEN_ALWAYS_REFRESH", "1")
	home := ".."
	os.Setenv("GOGEN_SAMPLES_DIR", filepath.Join(home, "config", "tests", "tokens.yml"))

	s := config.FindSampleInFile(home, "backfillrealtime")

	gq := make(chan *config.GenQueueItem, 1000)
	oq := make(chan *config.OutQueueItem)
	done := make(chan int)
	gqs := make([]*config.GenQueueItem, 0, 10)

	timer := &Timer{S: s, GQ: gq, OQ: oq, Done: done}
	go timer.NewTimer()

	time.Sleep(2 * time.Second)
Loop:
	for {
		select {
		case i := <-gq:
			gqs = append(gqs, i)
		default:
			break Loop
		}
	}
	inrange := len(gqs) >= 31 && len(gqs) <= 33
	assert.Equal(t, inrange, true)
}

func TestBackfillFutureEnd(t *testing.T) {
	os.Setenv("GOGEN_HOME", "..")
	os.Setenv("GOGEN_ALWAYS_REFRESH", "1")
	home := ".."
	os.Setenv("GOGEN_SAMPLES_DIR", filepath.Join(home, "config", "tests", "tokens.yml"))

	s := config.FindSampleInFile(home, "backfillfutureend")

	gq := make(chan *config.GenQueueItem, 1000)
	oq := make(chan *config.OutQueueItem)
	done := make(chan int)
	gqs := make([]*config.GenQueueItem, 0, 10)

	timer := &Timer{S: s, GQ: gq, OQ: oq, Done: done}
	go timer.NewTimer()
	<-done
Loop:
	for {
		select {
		case i := <-gq:
			gqs = append(gqs, i)
		default:
			break Loop
		}
	}
	assert.Equal(t, 10, len(gqs))
}
