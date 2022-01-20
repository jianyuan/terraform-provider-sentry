package logging

import "testing"

func TestLogLevelDisplayed(t *testing.T) {
	Info("Hello world", "this is", "fun")
	Infof("%s %s %s", "Hello world", "this is", "fun")
	Debug("Hello world", "this is", "fun")
	Debugf("%s %s %s", "Hello world", "this is", "fun")
	Trace("Hello world", "this is", "fun")
	Tracef("%s %s %s", "Hello world", "this is", "fun")
	Warning("Hello world", "this is", "fun")
	Warningf("%s %s %s", "Hello world", "this is", "fun")
	Error("Hello world", "this is", "fun")
	Errorf("%s %s %s", "Hello world", "this is", "fun")
}
