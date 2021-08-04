package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func RegisterService(reg Registration) error {

	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(reg)
	if err != nil {
		return err
	}

	res, err := http.Post(ServiceURL, "application/json", buf)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to register service. Register service respond with code %v", res.StatusCode)
	}

	return nil
}

func ShutdownService(serviceURL string) error {
	http.Delete(ServiceURL,)
}
