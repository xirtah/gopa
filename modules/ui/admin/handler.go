package admin

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/xirtah/gopa-framework/core/global"
	api "github.com/xirtah/gopa-framework/core/http"
	httprouter "github.com/xirtah/gopa-framework/core/http/router"
	"github.com/xirtah/gopa-framework/core/model"
	"github.com/xirtah/gopa-framework/core/persist"
	"github.com/xirtah/gopa-framework/core/util"
	"github.com/xirtah/gopa-spider/modules/config"
	"github.com/xirtah/gopa-spider/modules/ui/admin/console"
	"github.com/xirtah/gopa-spider/modules/ui/admin/dashboard"
	"github.com/xirtah/gopa-spider/modules/ui/admin/explore"
	"github.com/xirtah/gopa-spider/modules/ui/admin/setting"
	"github.com/xirtah/gopa-spider/modules/ui/admin/tasks"
	"gopkg.in/yaml.v2"
)

type AdminUI struct {
	api.Handler
}

func (h AdminUI) DashboardAction(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	dashboard.Index(w, r)
}

func (h AdminUI) TasksPageAction(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	var task []model.Task
	var count1, count2 int
	var host = h.GetParameterOrDefault(r, "host", "")
	var from = h.GetIntOrDefault(r, "from", 0)
	var size = h.GetIntOrDefault(r, "size", 20)
	var status = h.GetIntOrDefault(r, "status", -1)
	count1, task, _ = model.GetTaskList(from, size, host, status)

	err, hvs := model.GetHostStatus(status)
	if err != nil {
		panic(err)
	}

	err, kvs := model.GetTaskStatus(host)
	if err != nil {
		panic(err)
	}

	tasks.Index(w, r, host, status, from, size, count1, task, count2, hvs, kvs)
}

func (h AdminUI) TaskViewPageAction(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	if id == "" {
		panic(errors.New("id is nill"))
	}

	task, err := model.GetTask(id)
	if err != nil {
		panic(err)
	}

	total, snapshots, err := model.GetSnapshotList(0, 10, id)
	task.Snapshots = snapshots
	task.SnapshotCount = total

	tasks.View(w, r, task)
}

func (h AdminUI) ConsolePageAction(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	console.Index(w, r)
}

func (h AdminUI) ExplorePageAction(w http.ResponseWriter, r *http.Request) {

	explore.Index(w, r)
}

func (h AdminUI) GetScreenshotAction(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	bytes, err := persist.GetValue(config.ScreenshotBucketKey, []byte(id))
	if err != nil {
		h.Error(w, err)
		return
	}
	w.Write(bytes)
}

func (h AdminUI) SettingPageAction(w http.ResponseWriter, r *http.Request) {

	o, _ := yaml.Marshal(global.Env().RuntimeConfig)
	setting.Setting(w, r, string(o))
}

func (h AdminUI) UpdateSettingAction(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	body, _ := h.GetRawBody(r)
	yaml.Unmarshal(body, global.Env().RuntimeConfig) //TODO extract method, save to file

	o, _ := yaml.Marshal(global.Env().RuntimeConfig)

	setting.Setting(w, r, string(o))
}

//TODO: Clean this up, it is a hack - we have this code duplicated in the gopa-ui project
func (h AdminUI) RedirectHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	url := h.Get(r, "url", "")
	http.Redirect(w, r, util.UrlDecode(url), 302)
	return
}

//TODO: Clean this up, it is a hack - we have this code duplicated in the gopa-ui project
func (h *AdminUI) GetSnapshotPayloadAction(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	snapshot, err := model.GetSnapshot(id)
	if err != nil {
		h.Error(w, err)
		return
	}

	compressed := h.GetParameterOrDefault(req, "compressed", "true")
	var bytes []byte
	if compressed == "true" {
		bytes, err = persist.GetCompressedValue(config.SnapshotBucketKey, []byte(id))
	} else {
		bytes, err = persist.GetValue(config.SnapshotBucketKey, []byte(id))
	}

	if err != nil {
		h.Error(w, err)
		return
	}

	if len(bytes) > 0 {
		h.Write(w, bytes)

		//add link rewrite
		if util.ContainStr(snapshot.ContentType, "text/html") {
			h.Write(w, []byte("<script language='JavaScript' type='text/javascript'>"))
			h.Write(w, []byte(`var dom=document.createElement("div");dom.innerHTML='<div style="overflow: hidden;z-index: 99999999999999999;width:100%;height:18px;position: absolute top:1px;background:#ebebeb;font-size: 12px;text-align:center;">`))
			h.Write(w, []byte(fmt.Sprintf(`<a href="/"><img border=0 style="float:left;height:18px" src="%s"></a><span style="font-size: 12px;">Saved by Gopa, %v, <a title="%v" href="%v">View original</a></span></div>';var first=document.body.firstChild;  document.body.insertBefore(dom,first);`, nil, snapshot.Created, snapshot.Url, snapshot.Url)))
			//			h.Write(w, []byte(fmt.Sprintf(`<a href="/"><img border=0 style="float:left;height:18px" src="%s"></a><span style="font-size: 12px;">Saved by Gopa, %v, <a title="%v" href="%v">View original</a></span></div>';var first=document.body.firstChild;  document.body.insertBefore(dom,first);`, h.Config.SiteLogo, snapshot.Created, snapshot.Url, snapshot.Url)))
			h.Write(w, []byte("</script>"))
			h.Write(w, []byte("<script src=\"/static/assets/js/snapshot_footprint.js?v=1\"></script> "))
		}
		return
	}

	h.Error404(w)

}
