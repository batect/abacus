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
// +build journeyTests

package main_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("The application", func() {
	Context("when pinged", func() {
		var resp *http.Response

		BeforeEach(func() {
			var err error
			resp, err = http.Get("http://app:8080/ping")
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns a HTTP 200 response", func() {
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})
	})

	Context("when sent a session", func() {
		var resp *http.Response
		var eventReceivedByHoneycomb string

		BeforeEach(func() {
			event := `{
				"sessionId": "11112222-3333-4444-5555-666677778888", 
				"userId": "99990000-3333-4444-5555-666677778888", 
				"sessionStartTime": "2019-01-02T03:04:05.678Z", 
				"sessionEndTime": "2019-01-02T09:04:05.678Z", 
				"applicationId": "my-app", 
				"applicationVersion": "1.0.0",
				"metadata": {
					"operatingSystem": "Mac",
					"dockerVersion": "19.3.5"
				}
			}`

			resp = httpPut("http://app:8080/v1/sessions", event, "application/json")

			eventReceivedByHoneycomb = getFirstEventSentToHoneycomb()
		})

		It("returns a HTTP 201 response", func() {
			Expect(resp.StatusCode).To(Equal(http.StatusCreated))
		})

		It("returns an empty response body", func() {
			Expect(resp.Body).To(WithTransform(readAllBytes, BeEmpty()))
		})

		It("sends the event to Honeycomb", func() {
			Expect(eventReceivedByHoneycomb).To(MatchJSON(`{
				"time": "2019-01-02T03:04:05.678Z",
				"data": {
			        "sessionId": "11112222-3333-4444-5555-666677778888",
				    "userId": "99990000-3333-4444-5555-666677778888",
				    "sessionStartTime": "2019-01-02T03:04:05.678Z",
				    "sessionEndTime": "2019-01-02T09:04:05.678Z",
				    "applicationId": "my-app",
				    "applicationVersion": "1.0.0",
				    "metadata": {
					    "dockerVersion": "19.3.5",
					    "operatingSystem": "Mac"
				    }
				}
			}`))
		})
	})
})

func httpPut(url string, body string, contentType string) *http.Response {
	req, err := http.NewRequest("PUT", url, bytes.NewReader([]byte(body)))
	Expect(err).ToNot(HaveOccurred())

	req.Header.Set("Content-Type", contentType)

	resp, err := http.DefaultClient.Do(req)
	Expect(err).ToNot(HaveOccurred())

	return resp
}

func getFirstEventSentToHoneycomb() string {
	resp, err := http.Get("http://honeycomb-fake:3000/fake/events/batect-abacus-test/0")
	Expect(err).NotTo(HaveOccurred())
	defer resp.Body.Close()

	Expect(resp.StatusCode).To(Equal(http.StatusOK))

	return readAllBytes(resp.Body)
}

func readAllBytes(r io.ReadCloser) string {
	bytes, err := ioutil.ReadAll(r)
	Expect(err).NotTo(HaveOccurred())

	return string(bytes)
}
