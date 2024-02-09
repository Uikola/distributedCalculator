package discovery

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/rs/zerolog/log"
	"net/http"
)

// RegistryService регистрирует сервис на стороне оркестратора.
func RegistryService(name string) error {
	client := http.Client{}

	values := map[string]string{"name": name}
	jsonData, err := json.Marshal(values)
	if err != nil {
		log.Error().Err(err).Msg("can't marshal the data")
		return err
	}

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/registry", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Error().Err(err).Msg("can't create the request")
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("can't do the request")
		return err
	}

	if resp.StatusCode != 201 {
		return errors.New("error")
	}

	return nil
}
