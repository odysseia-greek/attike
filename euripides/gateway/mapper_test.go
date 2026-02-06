package gateway

import (
	"encoding/json"
	"testing"

	"github.com/odysseia-greek/attike/euripides/models"
	"github.com/stretchr/testify/assert"
)

func TestToModelTrace(t *testing.T) {
	// Raw Elasticsearch response JSON
	esResponseJSON := `{
  "_index": "tracing-2026.01.16",
  "_id": "0e0136fa-5f9b-469b-895e-07dac8e714df",
  "_version": 1,
  "_seq_no": 23,
  "_primary_term": 1,
  "found": true,
  "_source": {
    "items": [
      {
        "timestamp": "2026-01-16T10:13:48.403",
        "item_type": "TRACE_START",
        "span_id": "ae6bf4f023a1deaa",
        "parent_span_id": "ae6bf4f023a1deaa",
        "pod_name": "homeros-7499c7d86-bfjjs",
        "namespace": "olympia",
        "payload": {
          "common": {
            "traceId": "0e0136fa-5f9b-469b-895e-07dac8e714df",
            "spanId": "ae6bf4f023a1deaa",
            "parentSpanId": "ae6bf4f023a1deaa",
            "timestamp": "2026-01-16T10:13:48.403",
            "podName": "homeros-7499c7d86-bfjjs",
            "namespace": "olympia",
            "itemType": "TRACE_START"
          },
          "traceStart": {
            "method": "POST",
            "url": "/graphql",
            "host": "byzantium.odysseia-greek:8080",
            "remoteAddress": "::1",
            "rootQuery": "query mediaAnswer($input: MediaAnswerInput!) {\n  mediaAnswer(input: $input) {\n    correct\n    quizWord\n    finished\n    similarWords {\n      greek\n      english\n      __typename\n    }\n    foundInText {\n      rootword\n      conjugations {\n        word\n        rule\n        __typename\n      }\n      texts {\n        author\n        reference\n        referenceLink\n        book\n        text {\n          translations\n          greek\n          __typename\n        }\n        __typename\n      }\n      __typename\n    }\n    progress {\n      greek\n      translation\n      playCount\n      correctCount\n      incorrectCount\n      lastPlayed\n      __typename\n    }\n    __typename\n  }\n}",
            "operation": "mediaAnswer"
          }
        }
      },
      {
        "timestamp": "2026-01-16T10:13:48.403",
        "item_type": "GRAPHQL",
        "span_id": "c6244bb3bd606c29",
        "parent_span_id": "ae6bf4f023a1deaa",
        "pod_name": "sokrates-84b4f6d5c7-z85tx",
        "namespace": "apologia",
        "payload": {
          "common": {
            "traceId": "0e0136fa-5f9b-469b-895e-07dac8e714df",
            "spanId": "c6244bb3bd606c29",
            "parentSpanId": "ae6bf4f023a1deaa",
            "timestamp": "2026-01-16T10:13:48.403",
            "podName": "sokrates-84b4f6d5c7-z85tx",
            "namespace": "apologia",
            "itemType": "GRAPHQL"
          },
          "graphql": {
            "operation": "mediaAnswer",
            "rootQuery": "\nquery mediaAnswer($input: MediaAnswerInput!) {\n  mediaAnswer(input: $input) {\n    correct\n    quizWord\n    finished\n    similarWords {\n      greek\n      english\n      __typename\n    }\n    foundInText {\n      rootword\n      conjugations {\n        word\n        rule\n        __typename\n      }\n      texts {\n        author\n        reference\n        referenceLink\n        book\n        text {\n          translations\n          greek\n          __typename\n        }\n        __typename\n      }\n      __typename\n    }\n    progress {\n      greek\n      translation\n      playCount\n      correctCount\n      incorrectCount\n      lastPlayed\n      __typename\n    }\n    __typename\n  }\n}"
          }
        }
      },
      {
        "timestamp": "2026-01-16T10:13:48.404",
        "item_type": "TRACE_HOP",
        "span_id": "3f8741dd7a3fa07e",
        "parent_span_id": "c6244bb3bd606c29",
        "pod_name": "aristippos-5b4677c854-d2gvk",
        "namespace": "apologia",
        "payload": {
          "common": {
            "traceId": "0e0136fa-5f9b-469b-895e-07dac8e714df",
            "spanId": "3f8741dd7a3fa07e",
            "parentSpanId": "c6244bb3bd606c29",
            "timestamp": "2026-01-16T10:13:48.404",
            "podName": "aristippos-5b4677c854-d2gvk",
            "namespace": "apologia",
            "itemType": "TRACE_HOP"
          },
          "traceHop": {
            "method": "/aristippos.v1.Aristippos/Answer",
            "url": "/aristippos.v1.Aristippos/Answer",
            "host": "10.244.0.42:33962"
          }
        }
      },
      {
        "timestamp": "2026-01-16T10:13:48.404",
        "item_type": "TRACE_HOP_STOP",
        "span_id": "3f8741dd7a3fa07e",
        "parent_span_id": "c6244bb3bd606c29",
        "pod_name": "aristippos-5b4677c854-d2gvk",
        "namespace": "apologia",
        "payload": {
          "common": {
            "traceId": "0e0136fa-5f9b-469b-895e-07dac8e714df",
            "spanId": "3f8741dd7a3fa07e",
            "parentSpanId": "c6244bb3bd606c29",
            "timestamp": "2026-01-16T10:13:48.404",
            "podName": "aristippos-5b4677c854-d2gvk",
            "namespace": "apologia",
            "itemType": "TRACE_HOP_STOP"
          },
          "traceHopStop": {}
        }
      },
      {
        "timestamp": "2026-01-16T10:13:48.404",
        "item_type": "ACTION",
        "span_id": "ef770fa7c43fde94",
        "parent_span_id": "3f8741dd7a3fa07e",
        "pod_name": "aristippos-5b4677c854-d2gvk",
        "namespace": "apologia",
        "payload": {
          "common": {
            "traceId": "0e0136fa-5f9b-469b-895e-07dac8e714df",
            "spanId": "ef770fa7c43fde94",
            "parentSpanId": "3f8741dd7a3fa07e",
            "timestamp": "2026-01-16T10:13:48.404",
            "podName": "aristippos-5b4677c854-d2gvk",
            "namespace": "apologia",
            "itemType": "ACTION"
          },
          "action": {
            "action": "taken from cache with key: Basic+1+First Words 1",
            "status": "{\"content\":[{\"greek\":\"ὁ λόγος\",\"imageURL\":\"logos.webp\",\"translation\":\"word\"},{\"greek\":\"ὁ χρόνος\",\"imageURL\":\"time.webp\",\"translation\":\"time\"},{\"greek\":\"ὁ ποταμός\",\"imageURL\":\"river.webp\",\"translation\":\"river\"},{\"greek\":\"ὁ πόλεμος\",\"imageURL\":\"war.webp\",\"translation\":\"war\"},{\"greek\":\"ή ναυτικός\",\"imageURL\":\"naval.webp\",\"translation\":\"seafaring, naval\"},{\"greek\":\"ή στρατιά\",\"imageURL\":\"army.webp\",\"translation\":\"army\"},{\"greek\":\"τό δένδρον\",\"imageURL\":\"tree.webp\",\"translation\":\"tree\"},{\"greek\":\"ή ἀγαθός\",\"imageURL\":\"good.webp\",\"translation\":\"good\"},{\"greek\":\"ή κακός\",\"imageURL\":\"bad.webp\",\"translation\":\"bad\"},{\"greek\":\"ὁ ἵππος\",\"imageURL\":\"horse.webp\",\"translation\":\"horse\"},{\"greek\":\"ὁ θεός\",\"imageURL\":\"god.webp\",\"translation\":\"god\"},{\"greek\":\"ή θεά\",\"imageURL\":\"goddess.webp\",\"translation\":\"goddess\"},{\"greek\":\"ὁ οἶνος\",\"imageURL\":\"wine.webp\",\"translation\":\"wine\"},{\"greek\":\"ή θάλασσα\",\"imageURL\":\"sea.webp\",\"translation\":\"sea\"},{\"greek\":\"ή γέφυρα\",\"imageURL\":\"bridge.webp\",\"translation\":\"bridge\"},{\"greek\":\"ή σελήνη\",\"imageURL\":\"moon.webp\",\"translation\":\"moon\"},{\"greek\":\"ὁ ἥλιος\",\"imageURL\":\"sun.webp\",\"translation\":\"sun\"},{\"greek\":\"ή ἀσπίς\",\"imageURL\":\"shield.webp\",\"translation\":\"shield\"},{\"greek\":\"ή γυνή\",\"imageURL\":\"woman.webp\",\"translation\":\"woman\"},{\"greek\":\"ὁ ἀνήρ\",\"imageURL\":\"man.webp\",\"translation\":\"man\"}],\"quizMetadata\":{\"language\":\"English\"},\"segment\":\"First Words 1\",\"set\":1,\"theme\":\"Basic\"}"
          }
        }
      },
      {
        "timestamp": "2026-01-16T10:13:48.416",
        "item_type": "DB_SPAN",
        "span_id": "9e8e6f570d686dc8",
        "parent_span_id": "521d55de7ad8fe5d",
        "pod_name": "antigonos-f858797f4-6sljh",
        "namespace": "makedonia",
        "payload": {
          "common": {
            "traceId": "0e0136fa-5f9b-469b-895e-07dac8e714df",
            "spanId": "9e8e6f570d686dc8",
            "parentSpanId": "521d55de7ad8fe5d",
            "timestamp": "2026-01-16T10:13:48.416",
            "podName": "antigonos-f858797f4-6sljh",
            "namespace": "makedonia",
            "itemType": "DB_SPAN"
          },
          "dbSpan": {
            "action": "search",
            "query": "{\"query\":{\"bool\":{\"minimum_should_match\":1,\"should\":[{\"fuzzy\":{\"greek\":{\"fuzziness\":2,\"value\":\"ὁ πόλεμος\"}}},{\"fuzzy\":{\"normalized\":{\"fuzziness\":2,\"value\":\"πολεμος\"}}}]}},\"size\":20}",
            "hits": "5",
            "tookMs": "6"
          }
        }
      },
      {
        "timestamp": "2026-01-16T10:13:48.441",
        "item_type": "TRACE_STOP",
        "span_id": "ae6bf4f023a1deaa",
        "parent_span_id": "ae6bf4f023a1deaa",
        "pod_name": "homeros-7499c7d86-bfjjs",
        "namespace": "olympia",
        "payload": {
          "common": {
            "traceId": "0e0136fa-5f9b-469b-895e-07dac8e714df",
            "spanId": "ae6bf4f023a1deaa",
            "parentSpanId": "ae6bf4f023a1deaa",
            "timestamp": "2026-01-16T10:13:48.441",
            "podName": "homeros-7499c7d86-bfjjs",
            "namespace": "olympia",
            "itemType": "TRACE_STOP"
          },
          "traceStop": {
            "responseBody": "{\"data\":{\"mediaAnswer\":{\"correct\":true,\"quizWord\":\"ὁ πόλεμος\",\"finished\":false,\"similarWords\":[{\"greek\":\"πολεμέω\",\"english\":\"make war\",\"__typename\":\"Hit\"},{\"greek\":\"ποταμός\",\"english\":\"river\",\"__typename\":\"Hit\"},{\"greek\":\"πολέμιος -α -ον\",\"english\":\"hostile (m.pl.: the enemy)\",\"__typename\":\"Hit\"},{\"greek\":\"πότερος -α -ον\",\"english\":\"which of the two?\",\"__typename\":\"Hit\"}],\"foundInText\":{\"rootword\":\"ὁ πόλεμος\",\"conjugations\":null,\"texts\":[{\"author\":\"Thucydides\",\"reference\":\"2.65\",\"referenceLink\":\"/texts?author=Thucydides&book=History of the Peloponnesian War&reference=2.65\",\"book\":\"History of the Peloponnesian War\",\"text\":{\"translations\":[\"For as long as he was in authority in the city in time of peace, he governed the same with moderation and was a faithful watchman of it; and in his time it was at the greatest.\"],\"greek\":\"ὅσον τε γὰρ χρόνον προύστη τῆς πόλεως ἐν τῇ εἰρήνῃ, μετρίως ἐξηγεῖτο καὶ ἀσφαλῶς διεφύλαξεν αὐτήν, καὶ ἐγένετο ἐπ' ἐκείνου μεγίστη, ἐπειδή τε &&ὁ πόλεμος&& κατέστη, ὁ δὲ φαίνεται καὶ ἐν τούτῳ προγνοὺς τὴν δύναμιν.\",\"__typename\":\"Rhema\"},\"__typename\":\"AnalyzeResult\"}],\"__typename\":\"AnalyzeTextResponse\"},\"progress\":[{\"greek\":\"ή γυνή\",\"translation\":\"\",\"playCount\":0,\"correctCount\":0,\"incorrectCount\":0,\"lastPlayed\":\"0001-01-01T00:00:00Z\",\"__typename\":\"ProgressEntry\"},{\"greek\":\"ὁ ἵππος\",\"translation\":\"\",\"playCount\":0,\"correctCount\":0,\"incorrectCount\":0,\"lastPlayed\":\"0001-01-01T00:00:00Z\",\"__typename\":\"ProgressEntry\"},{\"greek\":\"ή θεά\",\"translation\":\"\",\"playCount\":0,\"correctCount\":0,\"incorrectCount\":0,\"lastPlayed\":\"0001-01-01T00:00:00Z\",\"__typename\":\"ProgressEntry\"},{\"greek\":\"ή σελήνη\",\"translation\":\"\",\"playCount\":0,\"correctCount\":0,\"incorrectCount\":0,\"lastPlayed\":\"0001-01-01T00:00:00Z\",\"__typename\":\"ProgressEntry\"},{\"greek\":\"ὁ ἥλιος\",\"translation\":\"\",\"playCount\":0,\"correctCount\":0,\"incorrectCount\":0,\"lastPlayed\":\"0001-01-01T00:00:00Z\",\"__typename\":\"ProgressEntry\"},{\"greek\":\"ή ναυτικός\",\"translation\":\"\",\"playCount\":0,\"correctCount\":0,\"incorrectCount\":0,\"lastPlayed\":\"0001-01-01T00:00:00Z\",\"__typename\":\"ProgressEntry\"},{\"greek\":\"ή ἀγαθός\",\"translation\":\"\",\"playCount\":0,\"correctCount\":0,\"incorrectCount\":0,\"lastPlayed\":\"0001-01-01T00:00:00Z\",\"__typename\":\"ProgressEntry\"},{\"greek\":\"ὁ θεός\",\"translation\":\"\",\"playCount\":0,\"correctCount\":0,\"incorrectCount\":0,\"lastPlayed\":\"0001-01-01T00:00:00Z\",\"__typename\":\"ProgressEntry\"},{\"greek\":\"ὁ οἶνος\",\"translation\":\"\",\"playCount\":0,\"correctCount\":0,\"incorrectCount\":0,\"lastPlayed\":\"0001-01-01T00:00:00Z\",\"__typename\":\"ProgressEntry\"},{\"greek\":\"ή θάλασσα\",\"translation\":\"\",\"playCount\":0,\"correctCount\":0,\"incorrectCount\":0,\"lastPlayed\":\"0001-01-01T00:00:00Z\",\"__typename\":\"ProgressEntry\"},{\"greek\":\"ὁ ἀνήρ\",\"translation\":\"\",\"playCount\":0,\"correctCount\":0,\"incorrectCount\":0,\"lastPlayed\":\"0001-01-01T00:00:00Z\",\"__typename\":\"ProgressEntry\"},{\"greek\":\"ὁ λόγος\",\"translation\":\"\",\"playCount\":0,\"correctCount\":0,\"incorrectCount\":0,\"lastPlayed\":\"0001-01-01T00:00:00Z\",\"__typename\":\"ProgressEntry\"},{\"greek\":\"ὁ χρόνος\",\"translation\":\"\",\"playCount\":0,\"correctCount\":0,\"incorrectCount\":0,\"lastPlayed\":\"0001-01-01T00:00:00Z\",\"__typename\":\"ProgressEntry\"},{\"greek\":\"ή στρατιά\",\"translation\":\"\",\"playCount\":0,\"correctCount\":0,\"incorrectCount\":0,\"lastPlayed\":\"0001-01-01T00:00:00Z\",\"__typename\":\"ProgressEntry\"},{\"greek\":\"τό δένδρον\",\"translation\":\"\",\"playCount\":0,\"correctCount\":0,\"incorrectCount\":0,\"lastPlayed\":\"0001-01-01T00:00:00Z\",\"__typename\":\"ProgressEntry\"},{\"greek\":\"ή γέφυρα\",\"translation\":\"\",\"playCount\":0,\"correctCount\":0,\"incorrectCount\":0,\"lastPlayed\":\"0001-01-01T00:00:00Z\",\"__typename\":\"ProgressEntry\"},{\"greek\":\"ὁ ποταμός\",\"translation\":\"\",\"playCount\":0,\"correctCount\":0,\"incorrectCount\":0,\"lastPlayed\":\"0001-01-01T00:00:00Z\",\"__typename\":\"ProgressEntry\"},{\"greek\":\"ὁ πόλεμος\",\"translation\":\"war\",\"playCount\":1,\"correctCount\":1,\"incorrectCount\":0,\"lastPlayed\":\"2026-01-16T10:13:46Z\",\"__typename\":\"ProgressEntry\"},{\"greek\":\"ή κακός\",\"translation\":\"\",\"playCount\":0,\"correctCount\":0,\"incorrectCount\":0,\"lastPlayed\":\"0001-01-01T00:00:00Z\",\"__typename\":\"ProgressEntry\"},{\"greek\":\"ή ἀσπίς\",\"translation\":\"shield\",\"playCount\":1,\"correctCount\":1,\"incorrectCount\":0,\"lastPlayed\":\"2026-01-16T10:13:41Z\",\"__typename\":\"ProgressEntry\"}],\"__typename\":\"ComprehensiveResponse\"}}}",
            "responseCode": 200,
            "timeStarted": "2026-01-16T10:13:48.403",
            "timeEnded": "2026-01-16T10:13:48.441",
            "totalTimeMs": "38",
            "isClosed": true
          }
        }
      }
    ],
    "isActive": false,
    "timeStarted": "2026-01-16T10:13:48.403",
    "timeEnded": "2026-01-16T10:13:48.441",
    "totalTime": 38,
    "responseCode": 200
  }
}`

	// Parse the ES response structure
	var esResponse struct {
		Source json.RawMessage `json:"_source"`
	}
	err := json.Unmarshal([]byte(esResponseJSON), &esResponse)
	assert.NoError(t, err)

	// Parse the _source into esTraceDoc
	var esModel models.EsTraceDoc
	err = json.Unmarshal(esResponse.Source, &esModel)
	assert.NoError(t, err)

	// Call the mapper function
	id := "0e0136fa-5f9b-469b-895e-07dac8e714df"
	traceModel, err := toModelTrace(id, esModel)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, traceModel)
	assert.Equal(t, id, traceModel.ID)
	assert.False(t, traceModel.IsActive)
	assert.Equal(t, int32(38), traceModel.TotalTimeMs)
	assert.Equal(t, int32(200), traceModel.ResponseCode)

	// Check top-level metadata
	assert.NotNil(t, traceModel.Namespace)
	assert.Equal(t, "olympia", *traceModel.Namespace)
	assert.NotNil(t, traceModel.PodName)
	assert.Equal(t, "homeros-7499c7d86-bfjjs", *traceModel.PodName)
	assert.NotNil(t, traceModel.Operation)
	assert.Equal(t, "mediaAnswer", *traceModel.Operation)

	// Check flags
	assert.True(t, traceModel.HasDbSpan)
	assert.True(t, traceModel.HasAction)

	// Check items count
	assert.Equal(t, 7, len(traceModel.Items))

	// Verify first item is TRACE_START
	firstItem := traceModel.Items[0]
	assert.Equal(t, "TRACE_START", string(firstItem.ItemType))
	assert.NotNil(t, firstItem.SpanID)
	assert.Equal(t, "ae6bf4f023a1deaa", *firstItem.SpanID)

	// Verify DB_SPAN item
	dbSpanItem := traceModel.Items[5]
	assert.Equal(t, "DB_SPAN", string(dbSpanItem.ItemType))
	assert.NotNil(t, dbSpanItem.SpanID)
	assert.Equal(t, "9e8e6f570d686dc8", *dbSpanItem.SpanID)

	// Verify ACTION item
	actionItem := traceModel.Items[4]
	assert.Equal(t, "ACTION", string(actionItem.ItemType))
	assert.NotNil(t, actionItem.SpanID)
	assert.Equal(t, "ef770fa7c43fde94", *actionItem.SpanID)

	// Verify last item is TRACE_STOP
	lastItem := traceModel.Items[6]
	assert.Equal(t, "TRACE_STOP", string(lastItem.ItemType))
	assert.NotNil(t, lastItem.SpanID)
	assert.Equal(t, "ae6bf4f023a1deaa", *lastItem.SpanID)
}
