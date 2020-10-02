package main

import (
	"regexp"
	"time"
)

var (
	clauses        []string
	volatileMemory map[string]string
	rgx            = regexp.MustCompile(`^Run\((.*?)\)`)
)

func InitializeKB() {
	volatileMemory = make(map[string]string)
}

func AddClause(clause string) bool {
	clauses = append(clauses, clause)
	return true
}

func RemoveClause(targetClause string) bool {
	for i, clause := range clauses {
		if clause == targetClause {
			clauses = append(clauses[:i], clauses[i+1:]...)
			return true
		}
	}
	return false
}

func PrintClauses() []string {
	return clauses
}

func PrintMemory() map[string]string {
	return volatileMemory
}

func ClearMemory(subject string) {
	delete(volatileMemory, subject)
}

func Ask() (plugins []string) {
	for i := 0; i < len(clauses); i++ {
		rs := rgx.FindStringSubmatch(clauses[i])
		if len(rs) >= 1 {
			plugins = append(plugins, rs[1])
		}
	}
	return
}

// RunKnowledgebase handles KB related requests via channels
func RunKnowledgebase(fromMeasure chan RMQMessage, toScheduler chan string) {
	for {
		select {
		case message := <-fromMeasure:
			addMeasure(message)
			// Let scheduler know if there is any
			// interesting measure occurred
			if reasoning(message.Topic, memory[message.Topic]) {
				toScheduler <- "knowledgebase"
			}
		}
	}
}

func addMeasure(message RMQMessage) {
	measure := Measure{
		Timestamp: time.Unix(0, message.Timestamp),
		Value:     message.Value,
	}
	if _, ok := memory[message.Topic]; !ok {
		memory[message.Topic] = CreateTimeseriesMeasure(&measure)
	} else {
		timeSeriesMeasures := memory[message.Topic]
		timeSeriesMeasures.Add(&measure)
	}
}

// reasoning reasons about the received measure using rules
// to see if it logically entails a clause that matches
// any registered conditions
func reasoning(topic string, measures *TimeseriesMeasures) bool {
	// TODO: this will be removed as we bring KB here
	if topic == "env.rain_gauge" {
		_, value := measures.Get()
		if value > 3. {
			AddClause("Rain(Now)")
			AddClause("Run(Cloud)")
		} else {
			RemoveClause("Rain(Now)")
			RemoveClause("Run(Cloud)")
		}
	}
	// TODO: return true only when there is an interesting change occurred
	return true
}
