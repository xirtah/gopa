package joint

import (
	"fmt"

	"github.com/xirtah/gopa-framework/core/errors"
	"github.com/xirtah/gopa-framework/core/filter"
	log "github.com/xirtah/gopa-framework/core/logger/seelog"
	"github.com/xirtah/gopa-framework/core/model"
	"github.com/xirtah/gopa-spider/modules/config"
)

// ContentDeduplicationJoint used to check the hash of page body, if duplicated hash already exists, will break the pipeline
type ContentDeduplicationJoint struct {
	model.Parameters
}

// Name return: content_deduplication
func (joint ContentDeduplicationJoint) Name() string {
	return "content_deduplication"
}

// Process the content hash Deduplication
func (joint ContentDeduplicationJoint) Process(c *model.Context) error {
	snapshot := c.MustGet(model.CONTEXT_SNAPSHOT).(*model.Snapshot)
	url := c.MustGetString(model.CONTEXT_TASK_URL)
	if snapshot.Hash != "" || c.MustGetInt(model.CONTEXT_TASK_Status) != model.TaskRedirected {

		log.Trace("check content deduplication, ", url)

		snapshot.Url = url
		taskID := c.MustGetString(model.CONTEXT_TASK_ID)
		snapshot.TaskID = taskID

		exist, depTaskID, depSnapshotId, depUrl := checkBySimHash(snapshot, c)

		msg := fmt.Sprintf("same content hash found, %s, %s, %s, duplicated with task: %s, snapshotID: %s, url: %s", taskID, url, snapshot.Hash, depTaskID, depSnapshotId, depUrl)

		if exist {
			c.Set(model.CONTEXT_TASK_Status, model.TaskDuplicated)
			c.End(msg)
			return errors.New(msg)
		}

	}

	return nil
}

func checkBySimHash(snapshot *model.Snapshot, c *model.Context) (bool, string, string, string) {

	simhash := snapshot.SimHash

	//Check local hash first
	if c.GetBool("check_filter", false) {
		exist, _ := filter.CheckThenAdd(config.ContentHashFilter, []byte(simhash))
		if exist {
			return true, "", "local_filter_cache", ""
		}
	}

	//Check hash from db
	items, err := model.GetSnapshotByField("sim_hash", simhash)

	if err != nil {
		panic(err)
	}

	if len(items) > 0 {
		for _, v := range items {
			log.Tracef("%s vs  %s , %s vs %s", v.Url, snapshot.Url, snapshot.TaskID, v.TaskID)
			if v.Url != snapshot.Url && v.TaskID != snapshot.TaskID {
				return true, v.TaskID, v.ID, v.Url
			}
		}
	}

	return false, "", "", ""
}
