// Copyright 2019-2022 Charles Korn.
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

package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

type serviceConfig struct {
	ServiceName     string
	ServiceVersion  string
	Port            string
	ProjectID       string
	HoneycombAPIKey string
}

func getConfig() (*serviceConfig, error) {
	port, err := getPort()

	if err != nil {
		return nil, fmt.Errorf("could not get port for service to listen to: %w", err)
	}

	projectID, err := getProjectID()

	if err != nil {
		return nil, fmt.Errorf("could not get project ID: %w", err)
	}

	honeycombAPIKey, err := getHoneycombAPIKey()

	if err != nil {
		return nil, fmt.Errorf("could not get Honeycomb API key: %w", err)
	}

	return &serviceConfig{
		ServiceName:     getServiceName(),
		ServiceVersion:  getServiceVersion(),
		Port:            port,
		ProjectID:       projectID,
		HoneycombAPIKey: honeycombAPIKey,
	}, nil
}

func getServiceName() string {
	return getEnvOrDefault("K_SERVICE", "abacus")
}

func getServiceVersion() string {
	return getEnvOrDefault("K_REVISION", "local")
}

func getEnvOrDefault(name string, fallback string) string {
	if value, ok := os.LookupEnv(name); ok {
		return value
	}

	return fallback
}

func getPort() (string, error) {
	return getEnv("PORT")
}

func getProjectID() (string, error) {
	return getEnv("GOOGLE_PROJECT")
}

func getHoneycombAPIKey() (string, error) {
	return getEnv("HONEYCOMB_API_KEY")
}

func getCredentialsFilePath() string {
	variableName := "GOOGLE_APPLICATION_CREDENTIALS"
	value := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")

	if value == "" {
		logrus.WithField("variable", variableName).Info("Credentials file environment variable is not set, will fallback to default credential sources for GCP connections.")
	}

	return value
}

func getEnv(name string) (string, error) {
	value := os.Getenv(name)

	if value == "" {
		return "", fmt.Errorf("environment variable '%v' is not set", name)
	}

	return value, nil
}
