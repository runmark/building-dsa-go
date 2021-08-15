package log

import (
	"bytes"
	"fmt"
	"github.com/runmark/distribute-app-go/registry"
	stdlog "log"
	"net/http"
)

func SetClientLog(logURL string, serviceName registry.ServiceName) {
	stdlog.SetPrefix(fmt.Sprintf("[%v] - ", serviceName))
	stdlog.SetFlags(0)
	stdlog.SetOutput(&clientLogger{})
}

type clientLogger struct {
	url string
}

func (cl clientLogger) Write(data []byte) (int, error) {

	res, err := http.Post(cl.url+"/log", "text/plain", bytes.NewReader(data))
	if err != nil {
		return 0, err
	}

	if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("Failed to send log message. Service respond with: %v", res.StatusCode)
	}

	return len(data), nil
}
