package genderizer

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"strings"
)

var endpoint string
var client = &http.Client{}

type Genderization struct {
	Name        string  `json:"name"`
	Gender      string  `json:"gender"`
	Probability float64 `json:"probability"`
	Count       uint64  `json:"count"`
}

func Genderize(names ...string) ([]*Genderization, error) {
	if len(names) == 0 {
		return nil, errors.New("Must provide at least one name.")
	}

	for index, name := range names {
		name = strings.TrimSpace(name)

		if name == "" {
			return nil, fmt.Errorf("Empty name detected at index %+v.", index)
		}

		names[index] = name
	}

	if endpoint == "" {
		endpoint = "https://api.genderize.io"
	}

	request, err := http.NewRequest("GET", endpoint, nil)

	if err != nil {
		return nil, err
	}

	query := request.URL.Query()

	for index, name := range names {
		query.Add(fmt.Sprintf("name[%d]", index), name)
	}

	request.URL.RawQuery = query.Encode()

	if client == nil {
		client = &http.Client{}
	}

	response, err := client.Do(request)

	if response != nil {
		defer response.Body.Close()
	}

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	var failure struct{
		error string `json:"error,omitempty"`
	}

	err = json.Unmarshal(body, &failure)

	if err == nil && failure.error != "" {
		return nil, errors.New(failure.error)
	}

	var results []*Genderization

	err = json.Unmarshal(body, &results)

	if err != nil {
		return nil, err
	}

	return results, nil
}
