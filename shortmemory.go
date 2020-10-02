package main

import (
	"errors"
	"time"
)

const (
	memoryLength = 10
)

var (
	memory map[string]*TimeseriesMeasures = make(map[string]*TimeseriesMeasures)
)

type Measure struct {
	Older     *Measure
	Later     *Measure
	Timestamp time.Time
	Value     float64
}

// TimeSeriesMeasures is a linked list holding measures.
type TimeseriesMeasures struct {
	Oldest *Measure
	Latest *Measure
}

func (m *TimeseriesMeasures) Add(measure *Measure) {
	if m.Length() >= memoryLength {
		oldest := m.Oldest
		m.Oldest = oldest.Later
		oldest = nil
	}
	latest := m.Latest
	latest.Later = measure
	m.Latest = measure
}

func (m *TimeseriesMeasures) Get() (time.Time, float64) {
	return m.Latest.Timestamp, m.Latest.Value
}

func (m *TimeseriesMeasures) Length() (length int32) {
	measure := m.Oldest
	for {
		length += 1
		if measure.Later == nil {
			break
		}
		measure = measure.Later
	}
	return
}

func (m *TimeseriesMeasures) Print() []float64 {
	var result []float64
	var measure *Measure = m.Oldest
	for {
		result = append(result, measure.Value)
		if measure.Later == nil {
			break
		}
		measure = measure.Later
	}
	return result
}

func CreateTimeseriesMeasure(measure *Measure) *TimeseriesMeasures {
	timeseriesMeasure := TimeseriesMeasures{}
	timeseriesMeasure.Oldest = measure
	timeseriesMeasure.Latest = measure
	return &timeseriesMeasure
}

// Sum adds up the values stored in the topic
func (m *TimeseriesMeasures) Sum() (sum float64) {
	measure := m.Oldest
	for {
		sum += measure.Value
		if measure.Later == nil {
			break
		}
		measure = measure.Later
	}
	return
}

// Avg averages the values stored in the topic
func (m *TimeseriesMeasures) Avg() float64 {
	return m.Sum() / float64(m.Length())
}

// AddMeasure takes topic, timestamp, and value and stores them
// in a time series database. Timestamp is a nano-second precision epoch value.
func AddMeasure(topic string, timestampNs int64, value float64) {
	if _, ok := memory[topic]; !ok {
		memory[topic] = CreateTimeseriesMeasure(&Measure{
			Timestamp: time.Unix(0, timestampNs),
			Value:     value,
		})
	} else {
		timeSeriesMeasures := memory[topic]
		timeSeriesMeasures.Add(&Measure{
			Timestamp: time.Unix(0, timestampNs),
			Value:     value,
		})
	}
}

func GetTopic(topic string) (*TimeseriesMeasures, error) {
	if _, ok := memory[topic]; ok {
		return memory[topic], nil
	} else {
		return nil, errors.New("No such topic exists")
	}
}
