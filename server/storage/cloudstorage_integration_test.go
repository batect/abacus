// Copyright 2019-2023 Charles Korn.
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

package storage_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	cloudstorage "cloud.google.com/go/storage"
	"github.com/batect/abacus/server/storage"
	"github.com/batect/abacus/server/types"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gomega_types "github.com/onsi/gomega/types"
	"google.golang.org/api/option"
)

var _ = Describe("Saving sessions to Cloud Storage", func() {
	var bucket *cloudstorage.BucketHandle
	var store storage.SessionStore

	session := &types.Session{
		SessionID:          "11112222-3333-4444-5555-666677778888",
		UserID:             "99990000-3333-4444-5555-666677778888",
		SessionStartTime:   time.Date(2019, 1, 2, 3, 4, 5, 678000000, time.UTC),
		SessionEndTime:     time.Date(2019, 1, 2, 9, 4, 5, 678000000, time.UTC),
		IngestionTime:      time.Date(2019, 1, 2, 20, 4, 5, 678000000, time.UTC),
		ApplicationID:      "my-app",
		ApplicationVersion: "1.0.0",
		Attributes: map[string]interface{}{
			"operatingSystem": "Mac",
			"dockerVersion":   "19.3.5",
			"counter":         json.Number("123"),
			"duration":        json.Number("1.3"),
			"isEnabled":       true,
			"nullValue":       nil,
		},
		Events: []types.Event{
			{
				Type: "ThingHappened",
				Time: time.Date(2019, 1, 2, 3, 4, 6, 678000000, time.UTC),
				Attributes: map[string]interface{}{
					"operatingSystem": "Mac",
					"counter":         json.Number("123"),
					"duration":        json.Number("1.3"),
					"isEnabled":       true,
					"nullValue":       nil,
				},
			},
		},
		Spans: []types.Span{
			{
				Type:      "LoadingThings",
				StartTime: time.Date(2019, 1, 2, 3, 4, 7, 678000000, time.UTC),
				EndTime:   time.Date(2019, 1, 2, 3, 4, 8, 678000000, time.UTC),
				Attributes: map[string]interface{}{
					"operatingSystem": "Mac",
					"counter":         json.Number("123"),
					"duration":        json.Number("1.3"),
					"isEnabled":       true,
					"nullValue":       nil,
				},
			},
		},
	}

	expectedJSON := `{
		"sessionId": "11112222-3333-4444-5555-666677778888",
		"userId": "99990000-3333-4444-5555-666677778888",
		"sessionStartTime": "2019-01-02T03:04:05.678Z",
		"sessionEndTime": "2019-01-02T09:04:05.678Z",
		"ingestionTime": "2019-01-02T20:04:05.678Z",
		"applicationId": "my-app",
		"applicationVersion": "1.0.0",
		"attributes": {
			"operatingSystem": "Mac",
			"dockerVersion": "19.3.5",
			"counter": 123,
			"duration": 1.3,
			"isEnabled": true,
			"nullValue": null
		},
		"events": [
			{ 
				"type": "ThingHappened", 
				"time": "2019-01-02T03:04:06.678Z", 
				"attributes": { 
					"operatingSystem": "Mac",
					"counter": 123,
					"duration": 1.3,
					"isEnabled": true,
					"nullValue": null
				} 
			}
		],
		"spans": [
			{ 
				"type": "LoadingThings", 
				"startTime": "2019-01-02T03:04:07.678Z", 
				"endTime": "2019-01-02T03:04:08.678Z", 
				"attributes": { 
					"operatingSystem": "Mac",
					"counter": 123,
					"duration": 1.3,
					"isEnabled": true,
					"nullValue": null
				} 
			}
		]
	}`

	BeforeEach(func() {
		project := "my-project"
		bucketName := "test-bucket-" + uuid.New().String()

		// Note that we also have to set the STORAGE_EMULATOR_HOST environment variable so that object downloads
		// are done from the correct host and over HTTP (rather than HTTPS).
		opts := []option.ClientOption{
			option.WithEndpoint("http://cloud-storage/storage/v1/"),
		}

		client, err := cloudstorage.NewClient(context.Background(), opts...)
		Expect(err).ToNot(HaveOccurred())

		bucket = client.Bucket(bucketName)
		err = bucket.Create(context.Background(), project, nil)
		Expect(err).ToNot(HaveOccurred())

		store, err = storage.NewCloudStorageSessionStore(bucketName, opts...)
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("given the session does not already exist", func() {
		var err error

		BeforeEach(func() {
			err = store.Store(context.Background(), session)
		})

		It("does not return an error", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("stores the session in the bucket at the expected path", func() {
			Expect(bucket.Object("v1/my-app/1.0.0/11112222-3333-4444-5555-666677778888.json")).To(HaveContent(MatchJSON(expectedJSON)))
		})

		It("stores the session in the bucket with the JSON media type", func() {
			Expect(bucket.Object("v1/my-app/1.0.0/11112222-3333-4444-5555-666677778888.json")).To(HaveContentType("application/json"))
		})

		It("stores the session in the bucket compressed", func() {
			Expect(bucket.Object("v1/my-app/1.0.0/11112222-3333-4444-5555-666677778888.json")).To(HaveContentEncoding("gzip"))
		})
	})

	Describe("given the session already exists", func() {
		var err error

		BeforeEach(func() {
			err = store.Store(context.Background(), session)
			Expect(err).ToNot(HaveOccurred())

			updatedSession := &types.Session{
				SessionID:          session.SessionID,
				UserID:             session.UserID,
				SessionStartTime:   session.SessionStartTime,
				SessionEndTime:     session.SessionEndTime,
				IngestionTime:      session.IngestionTime,
				ApplicationID:      session.ApplicationID,
				ApplicationVersion: session.ApplicationVersion,
				Attributes: map[string]interface{}{
					"some-new-attribute": "some value",
				},
			}

			err = store.Store(context.Background(), updatedSession)
		})

		It("returns an error that indicates the session already exists", func() {
			Expect(err).To(MatchError(storage.ErrAlreadyExists))
		})

		It("does not overwrite the existing session", func() {
			Expect(bucket.Object("v1/my-app/1.0.0/11112222-3333-4444-5555-666677778888.json")).To(HaveContent(MatchJSON(expectedJSON)))
		})
	})
})

type haveContentMatcher struct {
	expectedContentMatcher gomega_types.GomegaMatcher
	actualContent          string
}

func HaveContent(expectedContentMatcher gomega_types.GomegaMatcher) gomega_types.GomegaMatcher {
	return &haveContentMatcher{expectedContentMatcher, ""}
}

func (c *haveContentMatcher) Match(actual interface{}) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	reader, err := actual.(*cloudstorage.ObjectHandle).NewReader(ctx)

	if err != nil {
		return false, fmt.Errorf("could not get content of object: %w", err)
	}

	defer reader.Close()

	actualBytes, err := io.ReadAll(reader)

	if err != nil {
		return false, fmt.Errorf("could not read content of object: %w", err)
	}

	c.actualContent = string(actualBytes)

	return c.expectedContentMatcher.Match(c.actualContent)
}

func (c *haveContentMatcher) FailureMessage(_ interface{}) string {
	return c.expectedContentMatcher.FailureMessage(c.actualContent)
}

func (c *haveContentMatcher) NegatedFailureMessage(_ interface{}) string {
	return c.expectedContentMatcher.NegatedFailureMessage(c.actualContent)
}

type haveContentTypeMatcher struct {
	expectedContentType string
	actualContentType   string
}

func HaveContentType(expectedContentType string) gomega_types.GomegaMatcher {
	return &haveContentTypeMatcher{expectedContentType, ""}
}

func (c *haveContentTypeMatcher) Match(actual interface{}) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	attrs, err := actual.(*cloudstorage.ObjectHandle).Attrs(ctx)

	if err != nil {
		return false, fmt.Errorf("could not get attributes of object: %w", err)
	}

	c.actualContentType = attrs.ContentType

	return c.expectedContentType == c.actualContentType, nil
}

func (c *haveContentTypeMatcher) FailureMessage(actual interface{}) string {
	//nolint:forcetypeassert
	return fmt.Sprintf("Expected object '%v' to have content type '%v', but it was '%v'", actual.(*cloudstorage.ObjectHandle).ObjectName(), c.expectedContentType, c.actualContentType)
}

func (c *haveContentTypeMatcher) NegatedFailureMessage(actual interface{}) string {
	//nolint:forcetypeassert
	return fmt.Sprintf(
		"Expected object '%v' to not have content type '%v', but it was '%v'",
		actual.(*cloudstorage.ObjectHandle).ObjectName(),
		c.expectedContentType,
		c.actualContentType,
	)
}

type haveContentEncodingMatcher struct {
	expectedContentEncoding string
	actualContentEncoding   string
}

func HaveContentEncoding(expectedContentEncoding string) gomega_types.GomegaMatcher {
	return &haveContentEncodingMatcher{expectedContentEncoding, ""}
}

func (c *haveContentEncodingMatcher) Match(actual interface{}) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	attrs, err := actual.(*cloudstorage.ObjectHandle).Attrs(ctx)

	if err != nil {
		return false, fmt.Errorf("could not get attributes of object: %w", err)
	}

	c.actualContentEncoding = attrs.ContentEncoding

	return c.expectedContentEncoding == c.actualContentEncoding, nil
}

func (c *haveContentEncodingMatcher) FailureMessage(actual interface{}) string {
	//nolint:forcetypeassert
	return fmt.Sprintf(
		"Expected object '%v' to have content encoding '%v', but it was '%v'",
		actual.(*cloudstorage.ObjectHandle).ObjectName(),
		c.expectedContentEncoding,
		c.actualContentEncoding,
	)
}

func (c *haveContentEncodingMatcher) NegatedFailureMessage(actual interface{}) string {
	//nolint:forcetypeassert
	return fmt.Sprintf(
		"Expected object '%v' to not have content encoding '%v', but it was '%v'",
		actual.(*cloudstorage.ObjectHandle).ObjectName(),
		c.expectedContentEncoding,
		c.actualContentEncoding,
	)
}
