/*
Copyright 2018 Sameer Saini

*/

package joint

import (
	"encoding/json"

	"github.com/xirtah/gopa-framework/core/global"
	log "github.com/xirtah/gopa-framework/core/logger/seelog"
	"github.com/xirtah/gopa-framework/core/model"
	"github.com/xirtah/gopa-framework/core/util"
)

type NamedEntityRecognitionJoint struct {
	model.Parameters
}

const URLcoreNLP model.ParaKey = "URL_CoreNLP"

func (joint NamedEntityRecognitionJoint) Name() string {
	return "ner"
}

func (joint NamedEntityRecognitionJoint) Process(context *model.Context) error {

	snapshot := context.MustGet(model.CONTEXT_SNAPSHOT).(*model.Snapshot)

	if !util.PrefixStr(snapshot.ContentType, "text/") {
		log.Debugf("snapshot is not text, %s, %s , %s", snapshot.ID, snapshot.Url, snapshot.ContentType)
		return nil
	}

	if global.Env().IsDebug {
		log.Trace("text to be parsed for named entity recognition: ", snapshot.Text)
	}

	url := joint.GetStringOrDefault(URLcoreNLP, "http://localhost:9000/") + "?properties=%7B%22annotators%22:%20%22ner%22,%20%22outputFormat%22:%22json%22%7D"
	req := util.NewPostRequest(url, []byte(snapshot.Text))
	//req.SetBasicAuth(c.Config.Username, c.Config.Password)

	response, err := util.ExecuteRequest(req)

	if err != nil {
		//Do nothing
		log.Warn("Failed to get response from coreNLP server, URL:", url)
	} else {
		log.Debug("Response from coreNLP: ", string(response.Body))
		var result NERResultCoreNLP

		json.Unmarshal(response.Body, &result)

		for _, sentence := range result.Sentences {
			for _, entitymentions := range sentence.Entitymentions {
				//log.Debug("Parsed Entity - Word:", entitymentions.Text, "| Type:", entitymentions.Ner)
				if entitymentions.Ner == "ORGANIZATION" {
					snapshot.Organisations = append(snapshot.Organisations, entitymentions.Text)
				}
				if entitymentions.Ner == "PERSON" {
					snapshot.Persons = append(snapshot.Persons, entitymentions.Text)
				}
			}
		}
	}

	return nil
}

type NERResultCoreNLP struct {
	Sentences []struct {
		Index          int `json:"index"`
		Entitymentions []struct {
			DocTokenBegin        int    `json:"docTokenBegin"`
			DocTokenEnd          int    `json:"docTokenEnd"`
			TokenBegin           int    `json:"tokenBegin"`
			TokenEnd             int    `json:"tokenEnd"`
			Text                 string `json:"text"`
			CharacterOffsetBegin int    `json:"characterOffsetBegin"`
			CharacterOffsetEnd   int    `json:"characterOffsetEnd"`
			Ner                  string `json:"ner"`
		} `json:"entitymentions"`
		// Tokens []struct {
		// 	Index                int    `json:"index"`
		// 	Word                 string `json:"word"`
		// 	OriginalText         string `json:"originalText"`
		// 	Lemma                string `json:"lemma"`
		// 	CharacterOffsetBegin int    `json:"characterOffsetBegin"`
		// 	CharacterOffsetEnd   int    `json:"characterOffsetEnd"`
		// 	Pos                  string `json:"pos"`
		// 	Ner                  string `json:"ner"`
		// 	Before               string `json:"before"`
		// 	After                string `json:"after"`
		// } `json:"tokens"`
	} `json:"sentences"`
}
