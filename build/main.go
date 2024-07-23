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
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

//---------------------------------------------------------------------------------------

// Cloud Run Service Configuration
type ServiceConfig struct {
	Port           int
	TimeoutSeconds time.Duration
}

//---------------------------------------------------------------------------------------

// Initialise the Service
func init() {
	// Disable log prefixes such as the default timestamp.
	// Prefix text prevents the message from being parsed as JSON.
	// A timestamp is added when shipping logs to Cloud Logging.
	log.SetFlags(0)
}

//---------------------------------------------------------------------------------------

// Main Service
func main() {

	config := NewServiceConfig()
	PrintLogEntry(INFO, "Starting Service...")
	PrintLogEntry(INFO, fmt.Sprintf("listening on port %d", config.Port))

	// Initialise the engine and define router API Endpoints
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	apiv1 := router.Group("/api/v1")
	{
		apiv1.POST("/joshua", config.v1Joshua)
	}

	// Start listening and serve the API responses
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", config.Port),
		Handler:        router,
		ReadTimeout:    config.TimeoutSeconds,
		WriteTimeout:   config.TimeoutSeconds,
		MaxHeaderBytes: 1 << 20,
	}
	_ = s.ListenAndServe()

	PrintLogEntry(INFO, "Ending Service...")
}

//---------------------------------------------------------------------------------------

// Get a New Service Configuration from Environment Variables
func NewServiceConfig() ServiceConfig {

	// Default the Port to 8080 if not provided
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil || port == 0 {
		port = 8080
	}

	// Default the Timeout to 10 seconds if not provided
	timeout, err := strconv.Atoi(os.Getenv("CLOUD_RUN_TIMEOUT_SECONDS"))
	if err != nil || timeout == 0 {
		timeout = 10
	}

	return ServiceConfig{
		Port:           port,
		TimeoutSeconds: time.Duration(timeout) * time.Second,
	}

}

//---------------------------------------------------------------------------------------

// Abort the request with an error message
func AbortWithError(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}

//---------------------------------------------------------------------------------------

// JSON Structure defining the Request Message
type RequestMessage struct {
	Name string `json:"name"`
}

//---------------------------------------------------------------------------------------

// Accept a JSON Request Message and return a JSON Response Message
func (config *ServiceConfig) v1Joshua(c *gin.Context) {

	// Attempt to Bind the JSON Payload to a PubSubMessage instance
	var jsonData RequestMessage
	if err := c.ShouldBindBodyWith(&jsonData, binding.JSON); err != nil {
		msg := fmt.Sprintf("Unable to Bind RequestMessage JSON: %v", err)
		PrintLogEntry(DEBUG, msg)
		AbortWithError(c, http.StatusBadRequest, msg)
		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("Hello %s, would you like to play a game?", jsonData.Name))

}
