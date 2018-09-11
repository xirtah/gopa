/*
Copyright 2018 Sameer Saini
*/

package joint

import (
	"testing"

	"github.com/xirtah/gopa/core/env"
	"github.com/xirtah/gopa/core/global"
	"github.com/xirtah/gopa/core/model"
)

func TestProcessNamedEntityRecognition(t *testing.T) {

	global.RegisterEnv(env.EmptyEnv())

	context := model.Context{}
	context.Set(model.CONTEXT_TASK_URL, "http://microsoft.com/")
	context.Set(model.CONTEXT_TASK_Depth, 1)
	snapshot := model.Snapshot{}
	context.Set(model.CONTEXT_SNAPSHOT, &snapshot)
	snapshot.ContentType = "text/html"
	snapshot.Text = "Microsoft"

	NERParser := NamedEntityRecognitionJoint{}
	NERParser.Process(&context)

	//TODO: Add assertions

	//load file
	// b, e := ioutil.ReadFile("../../../test/samples/default.html")
	// if e != nil {
	// 	panic(e)
	// }
	// snapshot = model.Snapshot{}
	// context.Set(model.CONTEXT_SNAPSHOT, &snapshot)
	// snapshot.Payload = b
	// snapshot.ContentType = "text/html"

	// parse.Process(&context)

	// text = snapshot.Text
	// assert.Equal(t, "\nElastic中文社区\nlink\nHidden text, should not displayed!\nH1 title\nH2 title\n", text)

	//load file
	// b, e = ioutil.ReadFile("../../../test/samples/discuss.html")
	// if e != nil {
	// 	panic(e)
	// }
	// snapshot = model.Snapshot{}
	// context.Set(model.CONTEXT_SNAPSHOT, &snapshot)
	// snapshot.Payload = b
	// snapshot.ContentType = "text/html"

	// parse.Process(&context)

	// text = snapshot.Text
	// fmt.Println(text)

}
