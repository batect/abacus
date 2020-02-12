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

package main

import (
	"net/url"
	"os"

	"github.com/sirupsen/logrus"
)

func getPort() string {
	return getEnvOrExit("PORT")
}

func getProjectID() string {
	return getEnvOrExit("GOOGLE_PROJECT")
}

func getHoneycombBaseURL() url.URL {
	return getURLFromEnvOrExit("HONEYCOMB_BASE_URL")
}

func getHoneycombDatasetName() string {
	return getEnvOrExit("HONEYCOMB_DATASET_NAME")
}

func getHoneycombAPIKey() string {
	return getEnvOrExit("HONEYCOMB_API_KEY")
}

func getURLFromEnvOrExit(name string) url.URL {
	value := getEnvOrExit(name)

	var parsed *url.URL
	var err error

	parsed, err = url.Parse(value)

	if err != nil {
		logrus.WithField("variable", name).WithField("value", value).WithError(err).Fatal("Environment variable is not a valid URL.")
	}

	return *parsed
}

func getEnvOrExit(name string) string {
	value := os.Getenv(name)

	if value == "" {
		logrus.WithField("variable", name).Fatal("Environment variable is not set.")
	}

	return value
}
