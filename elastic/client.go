// Copyright 2017 The go2 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
// Setup env JSON value:
// go2_elastic={
//   "urls": [],
//   "healthcheck": {
//      "enable": false
//   },
//   "sniff": {
//       "enable": false
//       "scheme": "http"
//    }
// }
package elastic

import (
	"github.com/qiangli/go2/config"
	"github.com/qiangli/go2/logging"
	es "gopkg.in/olivere/elastic.v3"
)

var settings = config.AppSettings()
var log = logging.Logger()

//TODO add more options?
type ElasticEnv struct {
	Urls              []string      `env:"go2_elastic.urls"`
	HealthcheckEnable bool          `env:"go2_elastic.healthcheck.enable"`
	SniffEnable       bool          `env:"go2_elastic.sniff.enable"`
	SniffScheme       string        `env:"go2_elastic.sniff.scheme"`
}

var client *es.Client

func init() {
	env := ElasticEnv{}

	err := settings.Parse(&env)
	if err != nil {
		log.Errorf("Elastic init error: %v", err)
		return
	}

	log.Debugf("Elastic env: %v", env)

	client = initES(env)
}

func initES(env ElasticEnv) *es.Client {
	var options []es.ClientOptionFunc

	options = append(options, es.SetURL(env.Urls...))
	options = append(options, es.SetHealthcheck(env.HealthcheckEnable))
	options = append(options, es.SetSniff(env.SniffEnable))
	options = append(options, es.SetScheme(env.SniffScheme))

	c, err := es.NewClient(options...)
	if err != nil {
		panic(err)
	}

	return c
}

func Client() *es.Client {
	return client
}