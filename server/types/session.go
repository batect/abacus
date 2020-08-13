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

package types

import "time"

type Session struct {
	SessionID          string                 `json:"sessionId" validate:"required,uuid4"`
	UserID             string                 `json:"userId" validate:"required,uuid4"`
	SessionStartTime   time.Time              `json:"sessionStartTime" validate:"required"`
	SessionEndTime     time.Time              `json:"sessionEndTime" validate:"required,gtefield=SessionStartTime"`
	IngestionTime      time.Time              `json:"ingestionTime"`
	ApplicationID      string                 `json:"applicationId" validate:"required,applicationId"`
	ApplicationVersion string                 `json:"applicationVersion" validate:"required,version"`
	Attributes         map[string]interface{} `json:"attributes" validate:"dive,keys,required,attributeName,endkeys,attributeValue"`
	Events             []Event                `json:"events" validate:"dive"`
	Spans              []Span                 `json:"spans" validate:"dive"`
}

type Event struct {
	Type       string                 `json:"type" validate:"required"`
	Time       time.Time              `json:"time" validate:"required"`
	Attributes map[string]interface{} `json:"attributes" validate:"dive,keys,required,attributeName,endkeys,attributeValue"`
}

type Span struct {
	Type       string                 `json:"type" validate:"required"`
	StartTime  time.Time              `json:"startTime" validate:"required"`
	EndTime    time.Time              `json:"endTime" validate:"required,gtefield=StartTime"`
	Attributes map[string]interface{} `json:"attributes" validate:"dive,keys,required,attributeName,endkeys,attributeValue"`
}
