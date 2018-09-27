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

package impl

import (
	log "github.com/xirtah/gopa-framework/core/logger/seelog"
	"github.com/xirtah/gopa-framework/core/util"
	. "github.com/zeebo/sbloom"
	"hash/fnv"
	"io/ioutil"
)

type BloomFilter struct {
	persistFileName string
	filter          *Filter
}

func (filter *BloomFilter) Open(fileName string) error {

	filter.persistFileName = fileName

	//loading or initializing bloom filter
	if util.FileExists(fileName) {
		log.Debug("found bloomFilter,start reload,", fileName)
		n, err := ioutil.ReadFile(fileName)
		if err != nil {
			log.Error("bloomFilter:", fileName, err)
		}
		filter.filter = new(Filter)
		if err := filter.filter.GobDecode(n); err != nil {
			log.Error("bloomFilter:", fileName, err)
		}
		log.Info("bloomFilter successfully reloaded:", fileName)
	} else {

		probItems := 100000 //config.GetIntConfig("BloomFilter", "ItemSize", 100000)
		log.Debug("initializing bloom-filter", fileName, ",virual size is,", probItems)
		filter.filter = NewFilter(fnv.New64(), probItems)
		log.Info("bloomFilter successfully initialized:", fileName)
	}

	return nil
}

func (filter *BloomFilter) Close() error {

	log.Debug("bloomFilter start persist,file:", filter.persistFileName)

	//save bloom-filter
	m, err := filter.filter.GobEncode()
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(filter.persistFileName, m, 0600)
	if err != nil {
		panic(err)
	}
	log.Info("bloomFilter safety persisted.")

	return nil
}

func (filter *BloomFilter) Exists(key []byte) bool {
	return filter.filter.Lookup(key)
}

func (filter *BloomFilter) Add(key []byte) error {
	filter.filter.Add(key)
	return nil
}

func (filter *BloomFilter) Delete(key []byte) error {

	return nil
}
