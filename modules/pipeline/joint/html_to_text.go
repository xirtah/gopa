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
	"regexp"
	"strings"
	"sync"

	log "github.com/xirtah/gopa-framework/core/logger/seelog"
	"github.com/xirtah/gopa-framework/core/global"
	"github.com/xirtah/gopa-framework/core/model"
	"github.com/xirtah/gopa-framework/core/util"
)

type HtmlToTextJoint struct {
	model.Parameters
}

//merge whitespace and \n
const mergeWhitespace model.ParaKey = "merge_whitespace"
const removeNonScript model.ParaKey = "remove_nonscript"

func (joint HtmlToTextJoint) Name() string {
	return "html2text"
}

type cleanRule struct {
	l                   sync.RWMutex
	replaceRules        []*regexp.Regexp
	inited              bool
	removeTagsRule      *regexp.Regexp
	removeBreaksRule    *regexp.Regexp
	removeNonScriptRule *regexp.Regexp
	lowercase           bool
}

var rules = cleanRule{replaceRules: []*regexp.Regexp{}}

func getRule(str string) *regexp.Regexp {
	re, _ := regexp.Compile(str)
	return re
}

func initRules() {
	rules.l.Lock()
	defer rules.l.Unlock()
	if rules.inited {
		return
	}

	log.Trace("init html2text rule")

	//remove STYLE
	rules.replaceRules = append(rules.replaceRules, getRule(`<style[\S\s]+?\</style\>`))

	//remove META
	rules.replaceRules = append(rules.replaceRules, getRule(`\<meta[\S\s]+?\</meta\>`))

	//remove comments
	rules.replaceRules = append(rules.replaceRules, getRule(`<!--[\S\s]*?-->`))

	//remove SCRIPT
	rules.replaceRules = append(rules.replaceRules, getRule(`\<script[\S\s]+?.*?\</script\>`))

	//remove NOSCRIPT
	rules.removeNonScriptRule = getRule(`\<noscript[\S\s]+?\</noscript\>`)

	//remove iframe,frame
	rules.replaceRules = append(rules.replaceRules, getRule(`\<iframe[\S\s]+?\</iframe\>`))
	rules.replaceRules = append(rules.replaceRules, getRule(`\<frame[\S\s]+?\</frame\>`))
	rules.replaceRules = append(rules.replaceRules, getRule(`\<frameset[\S\s]+?\</frameset\>`))
	rules.replaceRules = append(rules.replaceRules, getRule(`\<noframes[\S\s]+?\</noframes\>`))

	//remove embed objects
	rules.replaceRules = append(rules.replaceRules, getRule(`\<noembed[\S\s]+?\</noembed\>`))
	rules.replaceRules = append(rules.replaceRules, getRule(`\<embed[\S\s]+?\</embed\>`))
	rules.replaceRules = append(rules.replaceRules, getRule(`\<applet[\S\s]+?\</applet\>`))
	rules.replaceRules = append(rules.replaceRules, getRule(`\<object[\S\s]+?\</object\>`))
	rules.replaceRules = append(rules.replaceRules, getRule(`\<base[\S\s]+?\</base\>`))

	//remove code blocks
	rules.replaceRules = append(rules.replaceRules, getRule(`\<pre[\S\s]+?\</pre\>`))
	rules.replaceRules = append(rules.replaceRules, getRule(`\<code[\S\s]+?\</code\>`))

	//remove all HTML tags and replaced with \n
	rules.removeTagsRule, _ = regexp.Compile("\\<[\\S\\s]+?\\>")

	//remove continued break lines
	rules.removeBreaksRule, _ = regexp.Compile("\\s{2,}")

	//lowercase all the text
	rules.lowercase = true

	rules.inited = true

}

// should equal to regex("\\<[\\S\\s]+?\\>").ReplaceAllStringFunc(str, strings.ToLower)
func lowercaseTag(str []byte) {

	startLowercase := false
	startLowercaseIndex := -1
	endLowercase := false
	endLowercaseIndex := -1

	for i, s := range str {
		if s == 60 {
			startLowercase = true
			startLowercaseIndex = i
		}
		if s == 62 {
			endLowercase = true
			endLowercaseIndex = i
		}
		if startLowercase && endLowercase {
			for j := startLowercaseIndex; j < endLowercaseIndex; j++ {
				x := str[j]
				if x > 64 && x < 91 {
					str[j] = x + 32
				}
			}
			startLowercase = false
			endLowercase = false
			startLowercaseIndex = -1
			endLowercaseIndex = -1
		}
	}
}

var empty = []byte(" ")

func replaceAll(src []byte) []byte {

	if rules.lowercase {
		lowercaseTag(src)
	}

	if rules.replaceRules != nil {
		for _, rule := range rules.replaceRules {
			src = rule.ReplaceAll(src, empty)
		}
	}

	if rules.removeTagsRule != nil {
		src = rules.removeTagsRule.ReplaceAll(src, []byte("\n"))
	}

	if rules.removeBreaksRule != nil {
		src = rules.removeBreaksRule.ReplaceAll(src, []byte("\n"))
	}

	return src
}

func (joint HtmlToTextJoint) Process(context *model.Context) error {
	initRules()

	snapshot := context.MustGet(model.CONTEXT_SNAPSHOT).(*model.Snapshot)

	if !util.PrefixStr(snapshot.ContentType, "text/") {
		log.Debugf("snapshot is not text, %s, %s , %s", snapshot.ID, snapshot.Url, snapshot.ContentType)
		return nil
	}

	body := snapshot.Payload
	if joint.GetBool(removeNonScript, true) {
		body = rules.removeNonScriptRule.ReplaceAll(body, []byte(""))
	}

	body = replaceAll(body)

	src := string(body)

	if joint.GetBool(mergeWhitespace, false) {
		src = util.MergeSpace(src)
	}

	src = strings.Replace(src, "&#8216;", "'", -1)
	src = strings.Replace(src, "&#8217;", "'", -1)
	src = strings.Replace(src, "&#8220;", "\"", -1)
	src = strings.Replace(src, "&#8221;", "\"", -1)
	src = strings.Replace(src, "&nbsp;", " ", -1)
	src = strings.Replace(src, "&quot;", "\"", -1)
	src = strings.Replace(src, "&apos;", "'", -1)
	src = strings.Replace(src, "&#34;", "\"", -1)
	src = strings.Replace(src, "&#39;", "'", -1)
	src = strings.Replace(src, "&amp; ", "& ", -1)
	src = strings.Replace(src, "&amp;amp; ", "& ", -1)

	snapshot.Text = util.XSSHandle(src)
	//probably set summary to text if it is empty (maybe first 300 chars or something)
	//snapshot.Summary = "KOALAS ARE COOL!" //TODO: Probably wrong place to set summary
	// if snapshot.Summary == "" {

	// 	snapshot.Summary = "SAMEER:" + snapshot.Text[0:300]
	// }

	if global.Env().IsDebug {
		log.Trace("get text: ", src)
	}

	//TODO: should probably move this into own separate module.
	//Parse text to CORENLP

	// //url := "http://localhost:9000/?properties={'annotators': 'ner', 'outputFormat':'json'}"
	// url := "http://localhost:9000/?properties=%7B%22annotators%22:%20%22ner%22,%20%22outputFormat%22:%22json%22%7D"
	// req := util.NewPostRequest(url, []byte(src))
	// //req := util.NewPostRequest(url, []byte("Microsoft"))

	// //req.SetBasicAuth(c.Config.Username, c.Config.Password)
	// response, err := util.ExecuteRequest(req)
	// log.Info("URL:", url)
	// if err != nil {
	// 	//Do nothing
	// 	//return nil, err"
	// 	log.Info("Failed to get response from coreNLP")
	// } else {
	// 	log.Info("1")
	// 	// 	testJson := `{
	// 	// 	"sentences": [
	// 	// 		{
	// 	// 			"index": 0,
	// 	// 			"entitymentions": [],
	// 	// 			"tokens": [
	// 	// 				{
	// 	// 					"index": 1,
	// 	// 					"word": "Koalas",
	// 	// 					"originalText": "Koalas",
	// 	// 					"lemma": "koala",
	// 	// 					"characterOffsetBegin": 0,
	// 	// 					"characterOffsetEnd": 6,
	// 	// 					"pos": "NNS",
	// 	// 					"ner": "O",
	// 	// 					"before": "",
	// 	// 					"after": " "
	// 	// 				}]
	// 	// 		}
	// 	// 	]
	// 	// }`
	// 	//var result map[string]interface{}
	// 	var result coreNLPResult

	// 	log.Info("2")
	// 	//json.Unmarshal([]byte(testJson), &result)
	// 	//err = json.Unmarshal(response.Body, &result)
	// 	json.Unmarshal(response.Body, &result)
	// 	log.Info("3")
	// 	//log.Info(result.Sentences[0].Tokens[0].Word)

	// 	for _, sentence := range result.Sentences {
	// 		for _, entitymentions := range sentence.Entitymentions {
	// 			log.Info("W:", entitymentions.Text, " N:", entitymentions.Ner)
	// 			if entitymentions.Ner == "ORGANIZATION" {
	// 				snapshot.Organisations = append(snapshot.Organisations, entitymentions.Text)
	// 			}
	// 			if entitymentions.Ner == "PERSON" {
	// 				snapshot.Persons = append(snapshot.Persons, entitymentions.Text)
	// 			}
	// 		}
	// 	}
	// 	// nlp := result["sentences"].(map[string]interface{})
	// 	// log.Info("4")
	// 	// for key, value := range nlp {
	// 	// 	log.Info("5")
	// 	// 	log.Info(key, value.(string))
	// 	// }
	// 	// log.Info("6")

	// 	//log.Info("response from coreNLP: ", string(response.Body))
	// }

	return nil
}

// type coreNLPResult struct {
// 	Sentences []struct {
// 		Index          int `json:"index"`
// 		Entitymentions []struct {
// 			DocTokenBegin        int    `json:"docTokenBegin"`
// 			DocTokenEnd          int    `json:"docTokenEnd"`
// 			TokenBegin           int    `json:"tokenBegin"`
// 			TokenEnd             int    `json:"tokenEnd"`
// 			Text                 string `json:"text"`
// 			CharacterOffsetBegin int    `json:"characterOffsetBegin"`
// 			CharacterOffsetEnd   int    `json:"characterOffsetEnd"`
// 			Ner                  string `json:"ner"`
// 		} `json:"entitymentions"`
// 		// Tokens []struct {
// 		// 	Index                int    `json:"index"`
// 		// 	Word                 string `json:"word"`
// 		// 	OriginalText         string `json:"originalText"`
// 		// 	Lemma                string `json:"lemma"`
// 		// 	CharacterOffsetBegin int    `json:"characterOffsetBegin"`
// 		// 	CharacterOffsetEnd   int    `json:"characterOffsetEnd"`
// 		// 	Pos                  string `json:"pos"`
// 		// 	Ner                  string `json:"ner"`
// 		// 	Before               string `json:"before"`
// 		// 	After                string `json:"after"`
// 		// } `json:"tokens"`
// 	} `json:"sentences"`
// }

// type coreNLPResult struct {
// 	Sentences []struct {
// 		Index          int           `json:"index"`
// 		Entitymentions []interface{} `json:"entitymentions"`
// 		Tokens         []struct {
// 			Index                int    `json:"index"`
// 			Word                 string `json:"word"`
// 			OriginalText         string `json:"originalText"`
// 			Lemma                string `json:"lemma"`
// 			CharacterOffsetBegin int    `json:"characterOffsetBegin"`
// 			CharacterOffsetEnd   int    `json:"characterOffsetEnd"`
// 			Pos                  string `json:"pos"`
// 			Ner                  string `json:"ner"`
// 			Before               string `json:"before"`
// 			After                string `json:"after"`
// 		} `json:"tokens"`
// 	} `json:"sentences"`
// }
