package tools_generator

import (
	log "github.com/cihub/seelog"
	. "github.com/xirtah/gopa/core/config"
	"github.com/xirtah/gopa/core/model"
	"github.com/xirtah/gopa/core/queue"
	"github.com/xirtah/gopa/core/util"
	"github.com/xirtah/gopa/modules/config"
	"time"
)

type GeneratorPlugin struct {
}

func (plugin GeneratorPlugin) Name() string {
	return "Generator"
}

func (plugin GeneratorPlugin) Start(cfg *Config) {

	generatorConfig := struct {
		TaskID  string `config:"task_id"`
		TaskUrl string `config:"task_url"`
	}{}

	cfg.Unpack(&generatorConfig)

	go func() {
		for {
			if generatorConfig.TaskUrl != "" {
				context := model.Context{IgnoreBroken: true}
				context.Set(model.CONTEXT_TASK_URL, generatorConfig.TaskUrl)
				err := queue.Push(config.CheckChannel, util.ToJSONBytes(context))
				if err != nil {
					log.Error(err)
				}
			}

			if generatorConfig.TaskID != "" {

				context := model.Context{}
				context.Set(model.CONTEXT_TASK_ID, generatorConfig.TaskID)
				err := queue.Push(config.FetchChannel, util.ToJSONBytes(context))
				if err != nil {
					log.Error(err)
				}
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()
}

func (plugin GeneratorPlugin) Stop() error {
	return nil
}
