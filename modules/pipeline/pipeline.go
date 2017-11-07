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

package pipeline

import (
	log "github.com/cihub/seelog"
	. "github.com/infinitbyte/gopa/core/config"
	"github.com/infinitbyte/gopa/core/errors"
	"github.com/infinitbyte/gopa/core/global"
	"github.com/infinitbyte/gopa/core/model"
	"github.com/infinitbyte/gopa/core/queue"
	"github.com/infinitbyte/gopa/core/stats"
	"github.com/infinitbyte/gopa/core/util"
	"github.com/infinitbyte/gopa/modules/config"
	. "github.com/infinitbyte/gopa/modules/pipeline/config"
	. "github.com/infinitbyte/gopa/modules/pipeline/joint"
	"runtime"
	"sync"
	"time"
)

var frameworkStarted bool

type PipelineFrameworkModule struct {
	Tasks map[string]*Task `config:"tasks"`
}

type Task struct {
	config         *TaskConfig
	l              sync.Mutex
	signalChannels []*chan bool
}

func (task *Task) Start(config *TaskConfig) {
	task.l.Lock()
	defer task.l.Unlock()

	numGoRoutine := task.config.MaxGoRoutine

	task.signalChannels = make([]*chan bool, numGoRoutine)
	//start fetcher
	for i := 0; i < numGoRoutine; i++ {
		log.Trace("start pipeline instance:", i)
		signalC := make(chan bool, 1)
		task.signalChannels[i] = &signalC
		go task.runPipeline(&signalC, i)

	}
}

func (task *Task) Update(config *TaskConfig) {
	task.Stop()
	task.Start(config)
}

func (task *Task) Stop() {
	task.l.Lock()
	defer task.l.Unlock()

	for i, item := range task.signalChannels {
		if item != nil {
			*item <- true
		}
		log.Debug("send exit signal to fetch channel: ", i)
	}
}

func (task *Task) runPipeline(singal *chan bool, shard int) {

	var taskInfo []byte
	for {
		select {
		case <-*singal:
			log.Trace("pipeline exit, shard:", shard)
			return
		case taskInfo = <-queue.ReadChan(config.FetchChannel):
			stats.Increment("queue."+string(config.FetchChannel), "pop")

			taskId, pipelineConfigId := model.DecodePipelineTask(taskInfo)

			pipelineConfig := task.config.DefaultPipelineConfig
			if pipelineConfigId != "" {
				var err error
				pipelineConfig, err = model.GetPipelineConfig(pipelineConfigId)
				if err != nil {
					panic(err)
				}
			}

			log.Trace("shard:", shard, ",task received:", taskId)
			task.execute(taskId, pipelineConfig)
			log.Trace("shard:", shard, ",task finished:", taskId)
		}
	}
}

func (task *Task) execute(taskId string, pipelineConfig *model.PipelineConfig) {
	var pipeline *model.Pipeline
	defer func() {
		if !global.Env().IsDebug {
			if r := recover(); r != nil {
				if e, ok := r.(runtime.Error); ok {
					log.Error("pipeline: ", pipeline.GetID(), ", taskId: ", taskId, ", ", util.GetRuntimeErrorMessage(e))
				}
				log.Error("error in pipeline,", util.ToJson(r, true), util.ToJson(pipeline.GetContext(), true))
			}
		}
	}()

	context := &model.Context{Phrase: config.PhraseCrawler}
	context.Set(CONTEXT_TASK_ID, taskId)

	pipeline = model.NewPipelineFromConfig(pipelineConfig, context)
	pipeline.Run()

	if task.config.ThresholdInMs > 0 {
		log.Debug("sleep ", task.config.ThresholdInMs, "ms to control crawling speed")
		time.Sleep(time.Duration(task.config.ThresholdInMs) * time.Millisecond)
		log.Debug("wake up now,continue crawing")
	}

	log.Trace("end pipeline")
}

// getDefaultTaskConfig return a default TaskConfig
func getDefaultCrawlerTaskConfig() TaskConfig {

	config := model.PipelineConfig{}
	start := model.JointConfig{}
	start.Enabled = true
	start.JointName = "init_task"
	config.StartJoint = &start
	save := model.JointConfig{}
	save.Enabled = true
	save.JointName = "save_task"

	urlNormalization := model.JointConfig{}
	urlNormalization.Enabled = true
	urlNormalization.JointName = "url_normalization"
	urlNormalization.Parameters = util.MapStr{
		"follow_all_domain": false,
		"follow_sub_domain": true,
	}

	fetchJoint := model.JointConfig{}
	fetchJoint.Enabled = true
	fetchJoint.JointName = "fetch"

	parse := model.JointConfig{}
	parse.Enabled = true
	parse.JointName = "parse"

	html2text := model.JointConfig{}
	html2text.Enabled = true
	html2text.JointName = "html2text"

	hash := model.JointConfig{}
	hash.Enabled = true
	hash.JointName = "hash"

	updateCheckTime := model.JointConfig{}
	updateCheckTime.Enabled = true
	updateCheckTime.JointName = "update_check_time"

	contentDeduplication := model.JointConfig{}
	contentDeduplication.Enabled = true
	contentDeduplication.JointName = "content_deduplication"

	langDetect := model.JointConfig{}
	langDetect.Enabled = true
	langDetect.JointName = "lang_detect"

	index := model.JointConfig{}
	index.Enabled = true
	index.JointName = "index"

	saveSnapshot := model.JointConfig{}
	saveSnapshot.Enabled = true
	saveSnapshot.JointName = "save_snapshot_db"

	config.EndJoint = &save
	config.ProcessJoints = []*model.JointConfig{
		&urlNormalization,
		&fetchJoint,
		&parse,
		&html2text,
		&hash,
		&updateCheckTime,
		&contentDeduplication,
		&langDetect,
		&saveSnapshot,
		&index,
	}

	defaultCrawlerConfig := TaskConfig{
		Name:                  "crawler",
		MaxGoRoutine:          10,
		TimeoutInMs:           60000,
		ThresholdInMs:         0,
		DefaultPipelineConfig: &config,
	}

	return defaultCrawlerConfig
}

func (module PipelineFrameworkModule) Name() string {
	return "Pipeline"
}

func (module PipelineFrameworkModule) Start(cfg *Config) {

	if frameworkStarted {
		log.Error("pipeline framework already started, please stop it first.")
		return
	}

	//init joints
	InitJoints()

	//config := GetDefaultTaskConfig()
	//cfg.Unpack(&config)
	//module.config = &config

	module.Tasks = map[string]*Task{}
	c := getDefaultCrawlerTaskConfig()

	if c.DefaultPipelineConfig == nil {
		panic(errors.Errorf("default pipeline config can't be null, %v", c))
	}

	module.Tasks[c.Name] = &Task{config: &c}

	for k, v := range module.Tasks {
		log.Debugf("startting pipeline: %s", k)
		v.Start(v.config)
		log.Infof("pipeline: %s started", k)
	}

	frameworkStarted = true
}

func (module PipelineFrameworkModule) Stop() error {
	if frameworkStarted {
		frameworkStarted = false
		log.Debug("start shutting down pipeline framework")
		for k, v := range module.Tasks {
			log.Infof("stopping pipeline: %s", k)
			v.Stop()
			log.Infof("pipeline: %s stopped", k)
		}
	} else {
		log.Error("pipeline framework is not started")
	}

	return nil
}