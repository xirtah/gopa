/*
Copyright 2016 Medcl (m AT medcl.net)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package storage

import (
	"github.com/xirtah/gopa-framework/core/config"
	"github.com/xirtah/gopa-framework/core/errors"
	"github.com/xirtah/gopa-framework/core/index"
	"github.com/xirtah/gopa-framework/core/persist"
	"github.com/xirtah/gopa-spider/modules/storage/boltdb"
	"github.com/xirtah/gopa-spider/modules/storage/elastic"
)

var impl boltdb.BoltdbStore

func (this StorageModule) Name() string {
	return "Storage"
}

var storeConfig *StorageConfig

type BoltdbConfig struct {
}

type StorageConfig struct {
	//Driver only `boltdb` and `elasticsearch` are available
	Driver  string                     `config:"driver"`
	Boltdb  *BoltdbConfig              `config:"boltdb"`
	Elastic *index.ElasticsearchConfig `config:"elasticsearch"`
}

var (
	defaultConfig = StorageConfig{
		Driver: "boltdb",
		Boltdb: &BoltdbConfig{},
		Elastic: &index.ElasticsearchConfig{
			Endpoint:    "http://localhost:9200",
			IndexPrefix: "gopa-", //TODO: Add support for elasticsearch credentials
		},
	}
)

func getDefaultConfig() StorageConfig {
	return defaultConfig
}

func (module StorageModule) Start(cfg *config.Config) {
	//Sameer - Looks like the storage module is used to store the snapshots of each url hit

	//init config
	config := getDefaultConfig()
	cfg.Unpack(&config)
	storeConfig = &config

	switch config.Driver {
	case "elasticsearch":
		client := index.ElasticsearchClient{Config: config.Elastic}
		handler := elastic.ElasticsearchStore{Client: &client}
		persist.RegisterKVHandler(handler)
	//TODO: Consider removing boltdb as a storage module driver as it does not support concurrent connections
	// case "boltdb":
	// 	folder := path.Join(global.Env().SystemConfig.GetWorkingDir(), "blob")
	// 	os.MkdirAll(folder, 0777)
	// 	impl = boltdb.BoltdbStore{FileName: path.Join(folder, "/bolt.db")}
	// 	err := impl.Open()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	persist.RegisterKVHandler(impl)
	default:
		panic(errors.Errorf("invalid driver, %s", config.Driver))
	}
}

func (module StorageModule) Stop() error {
	// if storeConfig.Driver == "boltdb" {
	// 	return impl.Close()
	// }
	return nil
}

type StorageModule struct {
}
