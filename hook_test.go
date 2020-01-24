package datadog

import (
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		log.Printf("%s:%d: "+msg+"\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		log.Printf("%s:%d: unexpected error: %s\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		log.Printf("%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

func getLogger(t *testing.T, isJSON bool) (*Hook, *logrus.Logger) {
	host := os.Getenv("DATADOG_HOST")
	apiKey := os.Getenv("DATADOG_APIKEY")
	Debug = true

	if host == "" {
		host = DatadogUSHost
	}
	if apiKey == "" {
		t.Fatal("skipping test; DATADOG_APIKEY not set")
	}

	hostName, _ := os.Hostname()
	hook := NewHook(host, apiKey, isJSON, 3, 5*time.Second)
	hook.Hostname = hostName
	l := logrus.New()
	l.Hooks.Add(hook)
	return hook, l
}

func TestHook(t *testing.T) {
	hook, l := getLogger(t, true)

	for _, level := range hook.Levels() {
		if len(l.Hooks[level]) != 1 {
			t.Errorf("Hook was not added. The length of l.Hooks[%v]: %v", level, len(l.Hooks[level]))
		}
	}
}
func TestSendingJSON(t *testing.T) {
	_, l := getLogger(t, true)

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			l.WithField("from", "unitest").Infof("TestSendingJSON - %d", i)
		}()
		time.Sleep(1 * time.Second)
	}

	wg.Wait()
}

func TestSendingPlain(t *testing.T) {
	_, l := getLogger(t, false)

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			l.WithField("from", "unitest").Infof("TestSendingPlain - %d", i)
		}()
		time.Sleep(1 * time.Second)
	}

	wg.Wait()
}
