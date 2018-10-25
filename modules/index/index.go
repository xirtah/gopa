package index

import (
	"encoding/json"
	"runtime"

	. "github.com/xirtah/gopa-framework/core/config"
	"github.com/xirtah/gopa-framework/core/global"
	core "github.com/xirtah/gopa-framework/core/index"
	log "github.com/xirtah/gopa-framework/core/logger/seelog"
	"github.com/xirtah/gopa-framework/core/model"
	"github.com/xirtah/gopa-framework/core/queue"
	"github.com/xirtah/gopa-spider/modules/config"
	common "github.com/xirtah/gopa-spider/modules/index/ui/common"
)

type IndexModule struct {
}

func (this IndexModule) Name() string {
	return "Index"
}

var signalChannel chan bool

var (
	defaultConfig = common.IndexConfig{
		Elasticsearch: &core.ElasticsearchConfig{
			Endpoint:    "http://localhost:9200",
			IndexPrefix: "gopa-",
		},
		UIConfig: &common.UIConfig{
			Enabled:     true,
			SiteName:    "GOPA",
			SiteFavicon: "/static/assets/img/favicon.ico",
			SiteLogo:    "/static/assets/img/logo.svg",
		},
	}
)

func (module IndexModule) Start(cfg *Config) {

	indexConfig := defaultConfig
	cfg.Unpack(&indexConfig)

	signalChannel = make(chan bool, 1)
	client := core.ElasticsearchClient{Config: indexConfig.Elasticsearch}

	go func() {
		defer func() {

			if !global.Env().IsDebug {
				if r := recover(); r != nil {

					if r == nil {
						return
					}
					var v string
					switch r.(type) {
					case error:
						v = r.(error).Error()
					case runtime.Error:
						v = r.(runtime.Error).Error()
					case string:
						v = r.(string)
					}
					log.Error("error in indexer,", v)
				}
			}
		}()

		for {
			select {
			case <-signalChannel:
				log.Trace("indexer exited")
				return
			default:
				log.Trace("waiting index signal")
				er, v := queue.Pop(config.IndexChannel)
				log.Trace("got index signal, ", string(v))
				if er != nil {
					log.Error(er)
					continue
				}
				//indexing to es or blevesearch
				doc := model.IndexDocument{}
				err := json.Unmarshal(v, &doc)
				if err != nil {
					panic(err)
				}

				client.Index(doc.Index, doc.ID, doc.Source)
			}

		}
	}()
}

func (module IndexModule) Stop() error {
	signalChannel <- true
	return nil
}
