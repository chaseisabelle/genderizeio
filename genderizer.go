package genderizer

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"strings"
)

const ENDPOINT = "https://api.genderize.io"

type Genderization struct {
	Name        string  `json:"name"`
	Gender      string  `json:"gender"`
	Probability float64 `json:"probability"`
	Count       uint64  `json:"count"`
}

type Genderizer struct {
	Client   *http.Client
	Endpoint string
	Key      string
}

func New() *Genderizer {
	return &Genderizer{
		Client:   &http.Client{},
		Endpoint: ENDPOINT,
		Key:      "",
	}
}

func (genderizer *Genderizer) Genderize(names ...string) ([]*Genderization, error) {
	// do we have any input?
	if len(names) == 0 {
		return nil, errors.New("Must provide at least one name.")
	}

	// check/prep input
	for index, name := range names {
		name = strings.TrimSpace(name)

		if name == "" {
			return nil, fmt.Errorf("Empty name detected at index %+v.", index)
		}

		names[index] = name
	}

	// build the request
	request, err := http.NewRequest("GET", genderizer.Endpoint, nil)

	if err != nil {
		return nil, err
	}

	// build the query params
	query := request.URL.Query()

	if genderizer.Key != "" {
		query.Add("apikey", genderizer.Key)
	}

	for index, name := range names {
		query.Add(fmt.Sprintf("name[%d]", index), name)
	}

	request.URL.RawQuery = query.Encode()

	// execute the request
	response, err := genderizer.Client.Do(request)

	if response != nil {
		defer response.Body.Close()
	}

	if err != nil {
		return nil, err
	}

	// load the response body
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	// parse the response body
	var payload interface{}

	err = json.Unmarshal(body, &payload)

	// check the response payload for error before checking status code
	// the payload will have a more detailed error message
	if err == nil {
		failure, ok := payload.(map[string]string)

		if ok {
			error, ok := failure["error"]

			if ok {
				return nil, errors.New(error)
			}
		}
	}

	// if the payload does not have error message then check status code
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%+v", http.StatusText(response.StatusCode))
	}

	// the payload failed to parse, return error
	if err != nil {
		return nil, err
	}

	// parse results from the payload
	results, ok := payload.([]interface{})

	if !ok {
		return nil, errors.Errorf("Malformed payload: %+v", payload)
	}

	var genderizations []*Genderization

	for _, result := range results {
		mapped := result.(map[string]interface{})
		genderization := &Genderization{}

		value, ok := mapped["name"]

		if ok {
			name, ok := value.(string)

			if ok {
				genderization.Name = name
			}
		}

		value, ok = mapped["gender"]

		if ok {
			gender, ok := value.(string)

			if ok {
				genderization.Gender = gender
			}
		}

		value, ok = mapped["probability"]

		if ok {
			probability, ok := value.(float64)

			if ok {
				genderization.Probability = probability
			}
		}

		value, ok = mapped["count"]

		if ok {
			// why cant i just do value.(uint64) ?
			count, ok := value.(float64)

			if ok {
				genderization.Count = uint64(count)
			}
		}

		genderizations = append(genderizations, genderization)
	}

	// and thats how you do the hokey pokey
	return genderizations, nil
}
