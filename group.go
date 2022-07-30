package main

import (
	"encoding/json"
)

type RunResult struct {
	SubmissionId string  `json:"submission_id"`
	TestIndex    uint64  `json:"test_index"`
	Status       string  `json:"status"`
	TimeUsage    float64 `json:"time_usage"`
	MemoryUsage  uint64  `json:"memory_usage"`
	Score        float64 `json:"score"`
	Message      string  `json:"message"`
}

type Group struct {
	Score        float64     `json:"score"`
	FullScore    float64     `json:"full_score"`
	SubmissionId string      `json:"submission_id"`
	GroupIndex   uint64      `json:"group_index"`
	RunResult    []RunResult `json:"run_result"`
}

func GroupToJSONString(group Group) string {
	bytes, _ := json.Marshal(group)

	return string(bytes)
}
