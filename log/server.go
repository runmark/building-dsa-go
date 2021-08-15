package log

import (
	"io/ioutil"
	stdlog "log"
	"net/http"
	"os"
)

var logger *stdlog.Logger

type fileLog string

func (fl fileLog) Write(data []byte) (n int, err error) {
	f, err := os.OpenFile(string(fl), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	defer f.Close()
	if err != nil {
		return 0, err
	}

	return f.Write(data)
}

func Run(destination string) {
	logger = stdlog.New(fileLog(destination), "", stdlog.LstdFlags)
}

func RegisterHandlers() {
	http.HandleFunc("/log", func(rw http.ResponseWriter, req *http.Request) {
		data, err := ioutil.ReadAll(req.Body)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		write(string(data))
	})
}

func write(message string) {
	logger.Printf("%v\n", message)
}