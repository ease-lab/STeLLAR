// MIT License
//
// Copyright (c) 2020 Theodor Amariucai
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package setup

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

//Configuration is the schema for all experiment configurations.
type Configuration struct {
	Sequential     bool            `json:"Sequential"`
	Provider       string          `json:"Provider"`
	Runtime        string          `json:"Runtime"`
	SubExperiments []SubExperiment `json:"SubExperiments"`
}

//SubExperiment is the schema for sub-experiment configurations.
type SubExperiment struct {
	Title                   string   `json:"Title"`
	Bursts                  int      `json:"Bursts"`
	BurstSizes              []int    `json:"BurstSizes"`
	PayloadLengthBytes      int      `json:"PayloadLengthBytes"`
	IATSeconds              float64  `json:"IATSeconds"`
	FunctionIncrementLimits []int64  `json:"FunctionIncrementLimits"`
	DesiredServiceTimes     []string `json:"DesiredServiceTimes"`
	IATType                 string   `json:"IATType"`
	PackageType             string   `json:"PackageType"`
	GatewaysNumber          int      `json:"GatewaysNumber"`
	Visualization           string   `json:"Visualization"`
	FunctionMemoryMB        int64    `json:"FunctionMemoryMB"`
	FunctionImageSizeMB     float64  `json:"FunctionImageSizeMB"`
	DataTransferChainLength int      `json:"DataTransferChainLength"`
	S3Transfer              bool     `json:"S3Transfer"`
	GatewayEndpoints        []GatewayEndpoint
	ID                      int
}

const (
	defaultVisualization             = "cdf"
	defaultIATType                   = "stochastic"
	defaultProvider                  = "aws"
	defaultRuntime                   = "go1.x"
	defaultPackageType               = "Zip"
	defaultGatewaysNumber            = 1
	defaultDataTransferChainLength   = 1
	defaultFunctionMemoryMB          = 128
	manyRequestsInBurstWarnThreshold = 2000
	manyFilesWarnThreshold           = 500
)

func extractConfiguration(configFile *os.File) Configuration {
	configByteValue, _ := ioutil.ReadAll(configFile)

	var parsedConfig Configuration
	if err := json.Unmarshal(configByteValue, &parsedConfig); err != nil {
		log.Fatalf("Could not extract experiment configuration from file: %s", err.Error())
	}

	if parsedConfig.Provider == "" {
		parsedConfig.Provider = defaultProvider
	}
	if parsedConfig.Runtime == "" {
		parsedConfig.Runtime = defaultRuntime
	}

	for index := range parsedConfig.SubExperiments {
		if parsedConfig.SubExperiments[index].Visualization == "" {
			parsedConfig.SubExperiments[index].Visualization = defaultVisualization
		}
		if parsedConfig.SubExperiments[index].PackageType == "" {
			parsedConfig.SubExperiments[index].PackageType = defaultPackageType
		}
		if parsedConfig.SubExperiments[index].IATType == "" {
			parsedConfig.SubExperiments[index].IATType = defaultIATType
		}
		if parsedConfig.SubExperiments[index].DataTransferChainLength == 0 {
			parsedConfig.SubExperiments[index].DataTransferChainLength = defaultDataTransferChainLength
		}
		if parsedConfig.SubExperiments[index].FunctionMemoryMB == 0 {
			parsedConfig.SubExperiments[index].FunctionMemoryMB = defaultFunctionMemoryMB
		}
		if parsedConfig.SubExperiments[index].GatewaysNumber == 0 {
			parsedConfig.SubExperiments[index].GatewaysNumber = defaultGatewaysNumber
		}
	}

	log.Debugf("Extracted %d sub-experiments from given configuration file.", len(parsedConfig.SubExperiments))
	return parsedConfig
}