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
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	protogen2 "github.com/vhive-serverless/stellar/src/setup/deployment/raw-code/functions/producer-consumer/proto_gen"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	log "github.com/sirupsen/logrus"
	"io"
	"strconv"
	"strings"
	"time"
)

var minioClientSingleton *minio.Client

func getMinioClient() *minio.Client {
	if minioClientSingleton != nil {
		return minioClientSingleton
	}

	const ( // vHive
		serverAddress = "10.96.0.46:9000"
		accessKey     = "minio"
		secretKey     = "minio123"
	)

	minioClient, err := minio.New(serverAddress, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("Could not create minio client: %s", err.Error())
	}

	minioClientSingleton = minioClient
	return minioClientSingleton
}

func isUsingStorage(requestGRPC *protogen2.InvokeChainRequest, requestHTTP *events.APIGatewayProxyRequest) bool {
	if requestHTTP != nil {
		value, hasStorageTransferField := requestHTTP.QueryStringParameters["StorageTransfer"]

		if !hasStorageTransferField {
			return false
		}

		storageTransfer, err := strconv.ParseBool(value)
		if err != nil {
			log.Errorf("Could not parse StorageTransfer: %s", err.Error())
		}

		return storageTransfer
	}

	// gRPC
	return requestGRPC.StorageTransfer == true
}

func loadObjectFromStorage(requestHTTP *events.APIGatewayProxyRequest, requestGRPC *protogen2.InvokeChainRequest) string {
	var objectKey, objectBucket string
	if requestHTTP != nil {
		objectKey = requestHTTP.QueryStringParameters["Key"]
		objectBucket = requestHTTP.QueryStringParameters["Bucket"]
	} else {
		objectKey = requestGRPC.Key
		objectBucket = requestGRPC.Bucket
	}

	if requestHTTP != nil {
		s3svc, _ := authenticateS3Client()
		object, err := s3svc.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(objectBucket),
			Key:    aws.String(objectKey),
		})
		if err != nil {
			log.Infof("Object %q not found in S3 bucket %q: %s", objectKey, objectBucket, err.Error())
		}

		payload, err := io.ReadAll(object.Body)
		if err != nil {
			log.Infof("Error reading object body: %v", err)
			return ""
		}

		return string(payload)
	}

	// when using anything but HTTP, at the moment, automatically resort to minio
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	storageClient := getMinioClient()
	object, err := storageClient.GetObject(
		ctx,
		objectBucket,
		objectKey,
		minio.GetObjectOptions{},
	)
	if err != nil {
		log.Infof("Object %q not found in bucket %q: %s", objectKey, objectBucket, err.Error())
	}

	var buf bytes.Buffer
	if _, err = io.Copy(&buf, object); err != nil {
		log.Infof("Error reading object body: %v", err)
		return ""
	}

	return buf.String()
}

func saveObjectToStorage(requestHTTP *events.APIGatewayProxyRequest, stringPayload string, requestGRPC *protogen2.InvokeChainRequest) {
	if requestHTTP != nil {
		key := saveObject(stringPayload, requestHTTP.QueryStringParameters["Bucket"], false)
		requestHTTP.QueryStringParameters["Key"] = key
	} else {
		// when using anything but HTTP, at the moment, automatically resort to minio
		key := saveObject(stringPayload, requestGRPC.Bucket, true)
		requestGRPC.Key = key
	}
}

func saveObject(payload string, bucket string, useMinio bool) string {
	key := fmt.Sprintf("transfer-payload-%dbytes-%s", int64(len(payload)), generateTrulyRandomBytes(20))
	log.Infof(`Using storage, saving transfer payload (~%d bytes) as %q to %q.`, len(payload), key, bucket)

	var uploadResult string
	if useMinio {
		storageClient := getMinioClient()

		uploadOutput, err := storageClient.PutObject(
			context.Background(),
			bucket,
			key,
			strings.NewReader(payload),
			int64(len(payload)),
			minio.PutObjectOptions{},
		)
		if err != nil {
			log.Fatalf("Unable to upload %q to %q, %v", key, bucket, err.Error())
		}
		uploadResult = uploadOutput.Location
	} else {
		_, s3uploader := authenticateS3Client()

		uploadOutput, err := s3uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
			Body:   strings.NewReader(payload),
		})
		if err != nil {
			log.Fatalf("Unable to upload %q to %q, %v", key, bucket, err.Error())
		}
		uploadResult = uploadOutput.Location
	}

	log.Infof("Successfully uploaded %q to bucket %q (%s)", key, bucket, uploadResult)
	return key
}
