// Copyright 2017 The go2 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
// Setup optional env JSON value:
// go2_blobstore={
//   "name": "",
// }
package blobstore

import (
	"github.com/qiangli/go2/config"
	"github.com/qiangli/go2/logging"
	"github.com/gosemver/aws_aws-sdk-go_v1.4.3-1-g1f24fa1/aws"
	"github.com/gosemver/aws_aws-sdk-go_v1.4.3-1-g1f24fa1/aws/credentials"
	"github.com/gosemver/aws_aws-sdk-go_v1.4.3-1-g1f24fa1/aws/session"
	"github.com/gosemver/aws_aws-sdk-go_v1.4.3-1-g1f24fa1/service/s3"
	"github.com/cloudfoundry-community/go-cfenv"
)

var settings = config.AppSettings()
var log = logging.Logger()

type BlobstoreEnv struct {
	Name string     `env:"go2_blobstore.name"`
}

func init() {
	env := BlobstoreEnv{}

	err := settings.Parse(&env)
	if err != nil {
		log.Errorf("Blobstore init error: %v", err)
		return
	}
	log.Debugf("Blobstore env: %v", env)

	initStore(env)
}

var (
	Config aws.Config
	Session *session.Session
	S3 *s3.S3
	BucketName string
)

func initStore(env BlobstoreEnv) {
	// The SDK requires a region. However, the endpoint will override this region.
	region := "us-east-1"
	disableSSL := true
	logLevel := aws.LogDebugWithRequestErrors

	s := settings.GetService(env.Name).(cfenv.Service)
	accessKeyID := s.Credentials["access_key_id"].(string)
	secretAccessKey := s.Credentials["secret_access_key"].(string)
	endpoint := s.Credentials["host"].(string)

	BucketName = s.Credentials["bucket_name"].(string)

	Config = aws.Config{
		Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
		Region:      &region,
		Endpoint:    &endpoint,
		DisableSSL:  &disableSSL,
		LogLevel:    &logLevel,
	}

	Session = session.New(&Config)

	S3 = s3.New(Session)

	S3.Handlers.Sign.Clear()
	S3.Handlers.Sign.PushBack(SignV2)
}
