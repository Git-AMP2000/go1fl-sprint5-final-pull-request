package trainings

import (
	"testing"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/personaldata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestActionInfoReturnsEmptyStringOnError(t *testing.T) {
	tests := []struct {
		name     string
		training Training
		wantErr  string
	}{
		{
			name: "неизвестный тип тренировки",
			training: Training{
				Steps:        6000,
				TrainingType: "Плавание",
				Duration:     time.Hour,
				Personal:     personaldata.Personal{Weight: 75, Height: 1.75},
			},
			wantErr: "неизвестный тип тренировки",
		},
		{
			name: "некорректный вес",
			training: Training{
				Steps:        6000,
				TrainingType: "Бег",
				Duration:     time.Hour,
				Personal:     personaldata.Personal{Weight: 0, Height: 1.75},
			},
			wantErr: "вес должен быть положительным конечным числом",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.training.ActionInfo()
			require.EqualError(t, err, tt.wantErr)
			assert.Empty(t, got)
		})
	}
}
