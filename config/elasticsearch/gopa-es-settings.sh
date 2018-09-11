curl -XPOST "http://localhost:9200/gopa-index/_close"

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

curl -XPOST "http://localhost:9200/gopa-index/_open"