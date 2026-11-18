package spentenergy

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDistanceAndMeanSpeedInvalidInputs(t *testing.T) {
	invalidHeights := []float64{0, -1.75, math.NaN(), math.Inf(1), math.Inf(-1)}
	for _, height := range invalidHeights {
		assert.Equal(t, 0.0, Distance(1000, height))
		assert.Equal(t, 0.0, MeanSpeed(1000, height, time.Hour))
	}

	assert.Equal(t, 0.0, Distance(-1000, 1.75))
	assert.Equal(t, 0.0, MeanSpeed(-1000, 1.75, time.Hour))
}

func TestCaloriesRejectInvalidAndNonFiniteInputs(t *testing.T) {
	tests := []struct {
		name           string
		steps          int
		weight, height float64
		duration       time.Duration
	}{
		{name: "нулевой рост при беге", steps: 1000, weight: 75, height: 0, duration: time.Hour},
		{name: "NaN в весе", steps: 1000, weight: math.NaN(), height: 1.75, duration: time.Hour},
		{name: "бесконечный вес", steps: 1000, weight: math.Inf(1), height: 1.75, duration: time.Hour},
		{name: "NaN в росте", steps: 1000, weight: 75, height: math.NaN(), duration: time.Hour},
		{name: "бесконечный рост", steps: 1000, weight: 75, height: math.Inf(1), duration: time.Hour},
		{name: "переполнение результата", steps: 6000, weight: math.MaxFloat64, height: 1.75, duration: time.Hour},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			running, err := RunningSpentCalories(tt.steps, tt.weight, tt.height, tt.duration)
			require.Error(t, err)
			assert.Equal(t, 0.0, running)

			walking, err := WalkingSpentCalories(tt.steps, tt.weight, tt.height, tt.duration)
			require.Error(t, err)
			assert.Equal(t, 0.0, walking)
		})
	}
}

func TestWalkingUsesHalfOfRunningCalories(t *testing.T) {
	running, err := RunningSpentCalories(6000, 75, 1.75, time.Hour)
	require.NoError(t, err)

	walking, err := WalkingSpentCalories(6000, 75, 1.75, time.Hour)
	require.NoError(t, err)

	assert.Equal(t, running*walkingCaloriesCoefficient, walking)
}
