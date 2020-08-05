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
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	cloudstorage "cloud.google.com/go/storage"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

type cloudStorageSessionStore struct {
	client *cloudstorage.Client
	bucket *cloudstorage.BucketHandle
}

func NewCloudStorageSessionStore(bucketName string, opts ...option.ClientOption) (SessionStore, error) {
	client, err := cloudstorage.NewClient(context.Background(), opts...)

	if err != nil {
		return nil, fmt.Errorf("could not create Cloud Storage client: %w", err)
	}

	store := cloudStorageSessionStore{
		client: client,
		bucket: client.Bucket(bucketName),
	}

	return &store, nil
}

func (c *cloudStorageSessionStore) Store(ctx context.Context, session *Session) error {
	w := c.bucket.
		Object(fmt.Sprintf("v1/%v.json", session.SessionID)).
		If(cloudstorage.Conditions{DoesNotExist: true}).
		NewWriter(ctx)
	w.ContentType = "application/json"

	bytes, err := json.Marshal(session)

	if err != nil {
		return fmt.Errorf("converting session to JSON failed: %w", err)
	}

	if _, err := w.Write(bytes); err != nil {
		return fmt.Errorf("writing to Cloud Storage failed: %w", err)
	}

	if err := w.Close(); err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == http.StatusPreconditionFailed {
			return AlreadyExistsError
		}

		return fmt.Errorf("storing session in Cloud Storage failed: %w", err)
	}

	return nil
}

func (c *cloudStorageSessionStore) CheckIfExists(ctx context.Context, session *Session) (bool, error) {
	panic("implement me")
}
