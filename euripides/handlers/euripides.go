package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/odysseia-greek/agora/aristoteles"
	elasticModels "github.com/odysseia-greek/agora/aristoteles/models"
	"github.com/odysseia-greek/agora/plato/logging"
)

type EuripidesHandler struct {
	Elastic aristoteles.Client
	Index   string
}

func (e *EuripidesHandler) Tracing(input map[string]interface{}) (*elasticModels.Hits, error) {
	var query map[string]interface{}
	var err error

	if len(input) == 0 {
		query = e.Elastic.Builder().MatchAll()
	} else {
		query, err = e.filterInputToQuery(input)
		if err != nil {
			return nil, err
		}
	}

	jsonQuery, _ := json.Marshal(query)
	logging.Debug(string(jsonQuery))
	response, err := e.Elastic.Query().MatchWithScroll(e.Index, query)
	if err != nil {
		return nil, err
	}

	return &response.Hits, nil
}

func (e *EuripidesHandler) filterInputToQuery(input map[string]interface{}) (map[string]interface{}, error) {
	// Create a query map
	query := make(map[string]interface{})

	// Create a filter map
	filter := make(map[string]interface{})

	// Create a list to accumulate filter conditions
	filterConditions := make([]map[string]interface{}, 0)

	// Extract and validate the "ids" filter
	if ids, ok := input["ids"].([]interface{}); ok {
		var validIDs []string
		for _, id := range ids {
			if idStr, ok := id.(string); ok {
				validIDs = append(validIDs, idStr)
			}
		}

		// Build a terms query for "_id" field
		query["terms"] = map[string]interface{}{
			"_id": validIDs,
		}
	} else {
		// Extract and validate the "statusCode" filter
		if statusCode, ok := input["statusCode"].(int); ok {
			// Add a filter for "statusCode"
			filterConditions = append(filterConditions, map[string]interface{}{
				"term": map[string]interface{}{
					"responseCode": statusCode,
				},
			})
		}

		// Extract and validate the "beginTime" and "endTime" filters
		if beginTime, ok := input["beginTime"].(string); ok {
			rangeFilter := map[string]interface{}{
				"timeEnded": map[string]interface{}{
					"gte": beginTime,
				},
			}

			// If "endTime" is also provided, add the "lte" condition
			if endTime, ok := input["endTime"].(string); ok {
				rangeFilter["timeEnded"].(map[string]interface{})["lte"] = endTime
			}

			// Add the range filter to filterConditions
			filterConditions = append(filterConditions, map[string]interface{}{
				"range": rangeFilter,
			})
		}

		// Extract and validate the "totalTimeHigherThan" filter
		if totalTimeHigherThan, ok := input["totalTimeHigherThan"].(int); ok {
			rangeFilter := map[string]interface{}{
				"totalTime": map[string]interface{}{
					"gte": totalTimeHigherThan,
				},
			}

			// Add the range filter to filterConditions
			filterConditions = append(filterConditions, map[string]interface{}{
				"range": rangeFilter,
			})
		}
		// Extract and validate the "podName" filter
		var podNameFilter map[string]interface{}
		if podName, ok := input["podName"].(string); ok {
			// Create a wildcard filter for "items.pod_name"
			podNameFilter = map[string]interface{}{
				"wildcard": map[string]interface{}{
					"items.common.pod_name.keyword": fmt.Sprintf("%s*", podName),
				},
			}
		}

		// Extract and validate the "operation" filter
		var operationFilter map[string]interface{}
		if operation, ok := input["operation"].(string); ok {
			// Create a wildcard filter for "items.operation"
			operationFilter = map[string]interface{}{
				"wildcard": map[string]interface{}{
					"items.operation": operation,
				},
			}
		}

		// Create a nested query for "items.podName" and "items.operation" if either exists
		if podNameFilter != nil || operationFilter != nil {
			nestedQuery := map[string]interface{}{
				"path": "items",
				"query": map[string]interface{}{
					"bool": map[string]interface{}{
						"should": []map[string]interface{}{},
					},
				},
			}

			// Add the "podName" wildcard filter to the nested query
			if podNameFilter != nil {
				nestedQuery["query"].(map[string]interface{})["bool"].(map[string]interface{})["should"] = append(
					nestedQuery["query"].(map[string]interface{})["bool"].(map[string]interface{})["should"].([]map[string]interface{}),
					podNameFilter,
				)
			}

			// Add the "operation" wildcard filter to the nested query
			if operationFilter != nil {
				nestedQuery["query"].(map[string]interface{})["bool"].(map[string]interface{})["should"] = append(
					nestedQuery["query"].(map[string]interface{})["bool"].(map[string]interface{})["should"].([]map[string]interface{}),
					operationFilter,
				)
			}

			// Add the nested query to the filterConditions
			filterConditions = append(filterConditions, map[string]interface{}{
				"nested": nestedQuery,
			})
		}

		// Add filterConditions to a bool query with "must"
		if len(filterConditions) > 0 {
			filter["bool"] = map[string]interface{}{
				"must": filterConditions,
			}
		}
	}

	if len(filter) > 0 {
		// Add the filter to the query
		query["bool"] = map[string]interface{}{
			"filter": []map[string]interface{}{
				filter,
			},
		}
	}

	return map[string]interface{}{"query": query}, nil
}
