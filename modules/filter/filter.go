package filter

import (
	log "github.com/cihub/seelog"
	. "github.com/xirtah/gopa/core/config"
	"github.com/xirtah/gopa/core/filter"
	. "github.com/xirtah/gopa/core/filter"
	"github.com/xirtah/gopa/modules/config"
	"github.com/xirtah/gopa/modules/filter/impl"
	"sync"
)

type FilterModule struct {
}

func (module FilterModule) Name() string {
	return "Filter"
}

func (module FilterModule) Exists(bucket Key, key []byte) bool {
	f := filters[bucket]
	return f.Exists(key)
}

func (module FilterModule) Add(bucket Key, key []byte) error {
	f := filters[bucket]
	return f.Add(key)
}

func (module FilterModule) Delete(bucket Key, key []byte) error {
	f := filters[bucket]
	return f.Delete(key)
}

var l sync.RWMutex

func (module FilterModule) CheckThenAdd(bucket Key, key []byte) (b bool, err error) {
	f := filters[bucket]
	l.Lock()
	defer l.Unlock()
	b = f.Exists(key)
	if !b {
		err = f.Add(key)
	}
	return b, err
}

func initFilter(key Key) {
	f := impl.PersistFilter{FilterBucket: string(key)}
	filters[key] = &f
}

var filters map[Key]*impl.PersistFilter

func (module FilterModule) Start(cfg *Config) {

	filters = map[Key]*impl.PersistFilter{}

	//TODO dynamic config
	initFilter(config.DispatchFilter)
	initFilter(config.FetchFilter)
	initFilter(config.CheckFilter)
	initFilter(config.ContentHashFilter)

	filter.Register(module)
}

func (module FilterModule) Stop() error {
	for _, v := range filters {
		err := (*v).Close()
		if err != nil {
			log.Error(err)
		}
	}
	return nil

}
