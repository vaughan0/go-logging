package logging

import (
	"strings"
	"testing"
)

type msgSlice []*Message

func (m *msgSlice) Output(msg *Message) {
	*m = append(*m, msg)
}

func (m *msgSlice) CreateOutputter(options map[string]string) (Outputter, error) {
	return m, nil
}

func mockSetup() {
	config := `
  [loggers]
  root = INFO, mock

  [mock]
  type = mock
  `
	if err := SetupReader(strings.NewReader(config)); err != nil {
		panic(err)
	}
}

func TestMultipleSetup(t *testing.T) {
	var logs msgSlice
	checkLogs := func(msgs ...string) {
		t.Logf("checking logs %v against strings %v", logs, msgs)
		if len(logs) != len(msgs) {
			t.FailNow()
		}
		for i, logged := range logs {
			if logged.Msg != msgs[i] {
				t.FailNow()
			}
		}
		logs = nil
	}
	RegisterOutputPlugin("mock", &logs)
	logger := Get("test")

	mockSetup()
	logger.Info("first")
	checkLogs("first")

	mockSetup()
	logger.Info("second")
	logger.Info("third")
	checkLogs("second", "third")

	mockSetup()
	checkLogs()
}
