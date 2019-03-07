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

package modules

import (
	"github.com/xirtah/gopa-framework/core/module"
	"github.com/xirtah/gopa-framework/modules/persist"
	"github.com/xirtah/gopa-spider/modules/api"
	"github.com/xirtah/gopa-spider/modules/cluster"
	"github.com/xirtah/gopa-spider/modules/dispatch"
	"github.com/xirtah/gopa-spider/modules/filter"
	"github.com/xirtah/gopa-spider/modules/index"
	"github.com/xirtah/gopa-spider/modules/pipeline"
	"github.com/xirtah/gopa-spider/modules/queue"
	"github.com/xirtah/gopa-spider/modules/stats"
	"github.com/xirtah/gopa-framework/modules/storage"
	"github.com/xirtah/gopa-spider/modules/ui"
)

// Register is where modules are registered
func Register() {
	module.Register(module.Database, persist.DatabaseModule{}) //TODO: Update to use the framework version
	module.Register(module.Storage, storage.StorageModule{})
	module.Register(module.Filter, filter.FilterModule{})
	module.Register(module.Stats, stats.SimpleStatsModule{})
	module.Register(module.Queue, queue.DiskQueue{})
	module.Register(module.Index, index.IndexModule{}) //TODO: Comment out
	module.Register(module.System, pipeline.PipelineFrameworkModule{})
	module.Register(module.System, dispatch.DispatchModule{})
	module.Register(module.System, cluster.ClusterModule{})
	module.Register(module.System, ui.UIModule{}) //TODO: Comment out
	module.Register(module.System, api.APIModule{})
}
