// Copyright 2019-2020 Charles Korn.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// and the Commons Clause License Condition v1.0 (the "Condition");
// you may not use this file except in compliance with both the License and Condition.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// You may obtain a copy of the Condition at
//
//     https://commonsclause.com/
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License and the Condition is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See both the License and the Condition for the specific language governing permissions and
// limitations under the License and the Condition.

package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type honeycombSessionStore struct {
	eventsEndpoint string
	apiKey         string
	client         *http.Client
}

func NewHoneycombSessionStore(baseURL url.URL, datasetName string, apiKey string) SessionStore {
	eventsEndpoint := baseURL.String() + "/1/events/" + datasetName

	return &honeycombSessionStore{
		eventsEndpoint: eventsEndpoint,
		apiKey:         apiKey,
		client: &http.Client{
			Timeout: 2 * time.Second,
		},
	}
}

func (h *honeycombSessionStore) Store(ctx context.Context, session *Session) error {
	var body []byte
	var err error

	if body, err = json.Marshal(session); err != nil {
		return fmt.Errorf("could not encode request body: %w", err)
	}

	var req *http.Request

	if req, err = http.NewRequestWithContext(ctx, "POST", h.eventsEndpoint, bytes.NewReader(body)); err != nil {
		return fmt.Errorf("could not create HTTP request: %w", err)
	}

	req.Header.Set("X-Honeycomb-Team", h.apiKey)
	req.Header.Set("X-Honeycomb-Event-Time", session.SessionStartTime.UTC().Format(time.RFC3339Nano))
	req.Header.Set("Content-Type", "application/json")

	var resp *http.Response

	if resp, err = h.client.Do(req); err != nil {
		return fmt.Errorf("request to Honeycomb (%s %s) failed with error: %w", req.Method, req.URL, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorMessage, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			errorMessage = []byte{}
		}

		return fmt.Errorf("request to Honeycomb (%s %s) received error response: HTTP %d: %s", req.Method, req.URL, resp.StatusCode, errorMessage)
	}

	return nil
}
