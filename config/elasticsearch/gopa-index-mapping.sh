printf "Deleting existing gopa-index db\n"
curl --user elastic:changeme -XDELETE "http://localhost:9200/gopa-index"
printf "\n"
printf "\n"

printf "Creating new gopa-index db\n"
curl --user elastic:changeme -XPUT "http://localhost:9200/gopa-index" -H 'Content-Type: application/json' -d'
{
  "mappings": {
    "doc": {
      "properties": {
        "host": {
          "type": "keyword",
          "ignore_above": 256
        },
        "snapshot": {
          "properties": {
            "bold": {
              "type": "text"
            },
            "url": {
              "type": "keyword",
              "ignore_above": 256
            },
            "content_type": {
              "type": "keyword",
              "ignore_above": 256
            },
            "file": {
              "type": "keyword",
              "ignore_above": 256
            },
            "ext": {
              "type": "keyword",
              "ignore_above": 256
            },
            "h1": {
              "type": "text"
            },
            "h2": {
              "type": "text"
            },
            "h3": {
              "type": "text"
            },
            "h4": {
              "type": "text"
            },
            "hash": {
              "type": "keyword",
              "ignore_above": 256
            },
            "id": {
              "type": "keyword",
              "ignore_above": 256
            },
            "images": {
              "properties": {
                "external": {
                  "properties": {
                    "label": {
                      "type": "text"
                    },
                    "url": {
                      "type": "keyword",
                      "ignore_above": 256
                    }
                  }
                },
                "internal": {
                  "properties": {
                    "label": {
                      "type": "text"
                    },
                    "url": {
                      "type": "keyword",
                      "ignore_above": 256
                    }
                  }
                }
              }
            },
            "italic": {
              "type": "text"
            },
            "links": {
              "properties": {
                "external": {
                  "properties": {
                    "label": {
                      "type": "text"
                    },
                    "url": {
                      "type": "keyword",
                      "ignore_above": 256
                    }
                  }
                },
                "internal": {
                  "properties": {
                    "label": {
                      "type": "text"
                    },
                  "url": {
                    "type": "keyword",
                    "ignore_above": 256
                  }
                }
              }
            }
          },
          "path": {
            "type": "keyword",
            "ignore_above": 256
          },
          "sim_hash": {
            "type": "keyword",
            "ignore_above": 256
          },
          "lang": {
            "type": "keyword",
            "ignore_above": 256
          },
          "screenshot_id": {
            "type": "keyword",
            "ignore_above": 256
          },
          "size": {
            "type": "long"
          },
          "organisations": {
            "type": "keyword"
          },
          "persons": {
            "type": "keyword"
          },
          "text": {
            "type": "text",
            "analyzer": "english"
          },
          "title": {
            "type": "text",
            "analyzer": "english",
            "fields": {
              "keyword": {
                  "type": "keyword"
                }
              }
            },
            "version": {
              "type": "long"
            }
          }
        },
        "task": {
          "properties": {
            "breadth": {
              "type": "long"
            },
            "created": {
             "type": "date"
            },
            "depth": {
              "type": "long"
            },
            "id": {
              "type": "keyword",
              "ignore_above": 256
            },
            "host" : {
              "type": "keyword"
            },
            "original_url": {
              "type": "keyword",
              "ignore_above": 256
            },
            "reference_url": {
              "type": "keyword",
              "ignore_above": 256
            },
            "schema": {
              "type": "keyword",
              "ignore_above": 256
            },
            "status": {
              "type": "integer"
            },
            "updated": {
              "type": "date"
            },
            "url": {
              "type": "keyword"
            },
            "last_screenshot_id": {
              "type": "keyword",
              "ignore_above": 256
            }
          }
        }
      }
    }
  }
}'
printf "\n"
printf "\n"

printf "Closing gopa-index db\n"
curl -XPOST "http://localhost:9200/gopa-index/_close"
printf "\n"
printf "\n"

printf "Updating gopa-index db settings\n"
curl --user elastic:changeme -XPUT "http://localhost:9200/gopa-index/_settings" -H 'Content-Type: application/json' -d'
{
  "settings": {
    "analysis": {
      "filter": {
        "english_stop": {
          "type":       "stop",
          "stopwords":  "_english_" 
        },
        "english_keywords": {
          "type":       "keyword_marker",
          "keywords":   [] 
        },
        "english_stemmer": {
          "type":       "stemmer",
          "language":   "english"
        },
        "english_possessive_stemmer": {
          "type":       "stemmer",
          "language":   "possessive_english"
        }
      },
      "analyzer": {
        "english": {
          "tokenizer":  "standard",
          "filter": [
            "english_possessive_stemmer",
            "lowercase",
            "english_stop",
            "english_keywords",
            "english_stemmer"
          ]
        }
      }
    }
  }
}
'
printf "\n"
printf "\n"

printf "Opening gopa-index db\n"
curl -XPOST "http://localhost:9200/gopa-index/_open"
printf "\n"
printf "\n"

printf "Deleting existing gopa-task db\n"
curl --user elastic:changeme -XDELETE "http://localhost:9200/gopa-task"
printf "\n"
printf "\n"

printf "Creating new gopa-task db\n"
curl --user elastic:changeme -XPUT "http://localhost:9200/gopa-task" -H 'Content-Type: application/json' -d'
{
  "mappings": {
    "doc": {
      "properties": {
        "host": {
          "type": "keyword",
          "ignore_above": 256
        }
      }
    }
  }
}'
printf "\n"
printf "\n"