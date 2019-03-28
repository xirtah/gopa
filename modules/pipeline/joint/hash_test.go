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

package joint

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xirtah/gopa-framework/core/model"
)

// func TestProcessHash(t *testing.T) {
// 	body := "Just some test content,你好"

// 	context := model.Context{}
// 	task := model.Task{}
// 	task.Url = "http://elasticsearch.cn/"
// 	task.Depth = 1

// 	context.Set(model.CONTEXT_TASK_URL, "http://elasticsearch.cn/")
// 	context.Set(model.CONTEXT_TASK_Depth, 1)
// 	parse := HashJoint{}

// 	snapshot := model.Snapshot{}
// 	snapshot.Payload = []byte(body)
// 	context.Set(model.CONTEXT_SNAPSHOT, &snapshot)

// 	parse.Process(&context)

// 	hash := snapshot.Hash
// 	assert.Equal(t, "b96aa2d91a6b69648250c8d4938b19c0750d7c50", hash)

// 	//TODO: work out what simhash is for - lzrbear
// 	hash1 := task.SnapshotSimHash
// 	//assert.Equal(t, "13442536247490772857", hash1)

// 	body = `Just some test content,你好啊,!! <a href="https://www.w3schools.com">Visit W3Schools.com!</a>`

// 	context = model.Context{}
// 	context.Set(model.CONTEXT_TASK_URL, "http://elasticsearch.cn/")
// 	context.Set(model.CONTEXT_TASK_Depth, 1)
// 	parse = HashJoint{}

// 	snapshot = model.Snapshot{}
// 	snapshot.Payload = []byte(body)
// 	context.Set(model.CONTEXT_SNAPSHOT, &snapshot)

// 	parse.Process(&context)
// 	hash2 := snapshot.SimHash
// 	fmt.Println(hash1)
// 	fmt.Println(hash2)
// 	assert.Equal(t, hash2, hash1)
// }

func TestProcessHashWithLinks(t *testing.T) {
	body := "Just some test content,你好啊,!!"

	context := model.Context{}
	task := model.Task{}
	task.Url = "http://elasticsearch.cn/"
	task.Depth = 1

	context.Set(model.CONTEXT_TASK_URL, "http://elasticsearch.cn/")
	context.Set(model.CONTEXT_TASK_Depth, 1)
	parse := HashJoint{}

	snapshot := model.Snapshot{}
	snapshot.Payload = []byte(body)
	context.Set(model.CONTEXT_SNAPSHOT, &snapshot)

	parse.Process(&context)

	hash1 := snapshot.Hash
	assert.Equal(t, "23fa8d5da158a827ccd12329935cc72ae6642109", hash1)

	body = `Just some test content,你好啊,!!<a href="https://www.w3schools.com">Visit W3Schools.com!</a>`

	context = model.Context{}
	context.Set(model.CONTEXT_TASK_URL, "http://elasticsearch.cn/")
	context.Set(model.CONTEXT_TASK_Depth, 1)
	parse = HashJoint{}

	snapshot = model.Snapshot{}
	snapshot.Payload = []byte(body)
	context.Set(model.CONTEXT_SNAPSHOT, &snapshot)

	parse.Process(&context)
	hash2 := snapshot.Hash
	fmt.Println(hash1)
	fmt.Println(hash2)
	assert.Equal(t, hash2, hash1)
}
