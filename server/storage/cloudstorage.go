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

package storage

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	cloudstorage "cloud.google.com/go/storage"
	"github.com/batect/abacus/server/types"
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

func (c *cloudStorageSessionStore) Store(ctx context.Context, session *types.Session) error {
	w := c.bucket.
		Object(fmt.Sprintf("v1/%v/%v/%v.json", session.ApplicationID, session.ApplicationVersion, session.SessionID)).
		If(cloudstorage.Conditions{DoesNotExist: true}).
		NewWriter(ctx)

	w.ContentType = "application/json"
	w.ContentEncoding = "gzip"
	gzipper := gzip.NewWriter(w)

	bytes, err := json.Marshal(session)

	if err != nil {
		return fmt.Errorf("converting session to JSON failed: %w", err)
	}

	if _, err := gzipper.Write(bytes); err != nil {
		return fmt.Errorf("writing to Cloud Storage failed: %w", err)
	}

	if err := gzipper.Close(); err != nil {
		return fmt.Errorf("closing gzip stream failed: %w", err)
	}

	if err := w.Close(); err != nil {
		var gerr *googleapi.Error

		if errors.As(err, &gerr) && gerr.Code == http.StatusPreconditionFailed {
			return ErrAlreadyExists
		}

		return fmt.Errorf("storing session in Cloud Storage failed: %w", err)
	}

	return nil
}
