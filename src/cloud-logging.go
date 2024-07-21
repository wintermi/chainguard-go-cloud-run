// Copyright 2024, Matthew Winter
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"log"
)

//---------------------------------------------------------------------------------------

// Cloud Logging Log Entry
type LogEntry struct {
	Message   string `json:"message"`
	Severity  string `json:"severity,omitempty"`
	Trace     string `json:"logging.googleapis.com/trace,omitempty"`
	Component string `json:"component,omitempty"`
}

// Cloud Logging Severity Type
type SeverityType int32

const (
	DEFAULT   SeverityType = 0
	DEBUG     SeverityType = 100
	INFO      SeverityType = 200
	NOTICE    SeverityType = 300
	WARNING   SeverityType = 400
	ERROR     SeverityType = 500
	CRITICAL  SeverityType = 600
	ALERT     SeverityType = 700
	EMERGENCY SeverityType = 800
)

// Cloud Logging Severity Message
var SeverityMessage = map[SeverityType]string{
	DEFAULT:   "DEFAULT",
	DEBUG:     "DEBUG",
	INFO:      "INFO",
	NOTICE:    "NOTICE",
	WARNING:   "WARNING",
	ERROR:     "ERROR",
	CRITICAL:  "CRITICAL",
	ALERT:     "ALERT",
	EMERGENCY: "EMERGENCY",
}

//---------------------------------------------------------------------------------------

// Print a Log Entry in a JSON format compatible with Cloud Logging.
func PrintLogEntry(severity SeverityType, message string) {

	entry, err := json.Marshal(LogEntry{
		Severity:  SeverityMessage[severity],
		Message:   message,
		Component: "export-api-data",
	})
	if err != nil {
		log.Printf("Failed to Execute JSON Marshal: %v", err)
		return
	}
	log.Println(string(entry))

}
