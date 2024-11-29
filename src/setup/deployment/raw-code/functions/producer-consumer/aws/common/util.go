// MIT License
//
// Copyright (c) 2021 Theodor Amariucai and EASE Lab
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

package common

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	allowedChars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

//InitializeGlobalRandomPayload creates the initial transfer payload to be used for quicker random payload generation
func InitializeGlobalRandomPayload() {
	const (
		uniquePayloadBytesLength = 1024 * 1024 // 1 MB, unique and randomized
	)

	uniquePayload := generateTrulyRandomBytes(uniquePayloadBytesLength)
	GlobalRandomPayload = string(uniquePayload)

	// Double the random payload 3 times by concatenation, fastest way to get from 1MB to 8MB
	length := len(GlobalRandomPayload)
	for i := 0; i < 3; i++ {
		length *= 2
		GlobalRandomPayload = GeneratePayloadFromGlobalRandom(length)
	}
}

func generateTrulyRandomBytes(uniquePayloadBytesLength int) []byte {
	rand.Seed(time.Now().UnixNano())
	uniquePayload := make([]byte, uniquePayloadBytesLength)
	for i := range uniquePayload {
		uniquePayload[i] = allowedChars[rand.Intn(len(allowedChars))]
	}
	return uniquePayload
}

//GeneratePayloadFromGlobalRandom creates a transfer payload for the producer-consumer chain
func GeneratePayloadFromGlobalRandom(payloadLengthBytes int) string {
	var repeatedRandomPayload strings.Builder
	repeatedRandomPayload.WriteString(GlobalRandomPayload)

	// Doubling the payload every time is faster than concatenating with the static "GlobalRandomPayload"
	for repeatedRandomPayload.Len() < payloadLengthBytes {
		repeatedRandomPayload.WriteString(repeatedRandomPayload.String())
	}

	return repeatedRandomPayload.String()[:payloadLengthBytes]
}

//extractJSONTimestampChain will process raw bytes into a string array of timestamps
func extractJSONTimestampChain(responsePayload []byte) []string {
	var reply map[string]interface{}
	err := json.Unmarshal(responsePayload, &reply)
	if err != nil {
		log.Fatalf("Could not unmarshal lambda response into map[string]interface{}: %s", err)
	}

	var parsedReply ProducerConsumerResponse
	err = json.Unmarshal([]byte(reply["body"].(string)), &parsedReply)
	if err != nil {
		log.Fatalf("Could not unmarshal lambda response body into producerConsumerResponse: %s", err)
	}

	return parsedReply.TimestampChain
}

//AppendTimestampToChain will add a new timestamp to the chain
func AppendTimestampToChain(timestampChain []string) []string {
	timestampMilliString := strconv.FormatInt(time.Now().UnixNano()/(int64(time.Millisecond)/int64(time.Nanosecond)), 10)
	return append(timestampChain, timestampMilliString)
}

//StringArrayToArrayOfString will process, e.g., "[14 35 8]" into []string{14, 35, 8}
func StringArrayToArrayOfString(str string) []string {
	str = strings.Split(str, "]")[0]
	str = strings.Split(str, "[")[1]
	return strings.Split(str, " ")
}
