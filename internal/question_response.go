package internal

import "time"

type QuestionResponse struct {
	Question                string        `json:"question" yaml:"question"`
	AnnotatedQuestion       string        `json:"annotated_question" yaml:"annotated_question"`
	Command                 string        `json:"command" yaml:"command"`
	Answer                  string        `json:"answer" yaml:"answer"`
	Explanation             string        `json:"explanation" yaml:"explanation"`
	QuestionResponseTime    time.Duration `json:"question_response_time" yaml:"question_response_time"`
	ExplanationResponseTime time.Duration `json:"explanation_response_time" yaml:"explanation_response_time"`
}
