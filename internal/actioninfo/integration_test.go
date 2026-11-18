package actioninfo

import (
	"bytes"
	"log"
	"os"
	"testing"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/daysteps"
	"github.com/Yandex-Practicum/tracker/internal/personaldata"
	"github.com/Yandex-Practicum/tracker/internal/trainings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Оба прикладных типа должны удовлетворять общему интерфейсу.
var (
	_ DataParser = (*daysteps.DaySteps)(nil)
	_ DataParser = (*trainings.Training)(nil)
)

func TestInfoWithRealParsers(t *testing.T) {
	person := personaldata.Personal{Name: "Тест", Weight: 75, Height: 1.75}

	tests := []struct {
		name        string
		dataset     []string
		makeParser  func() DataParser
		wantOutput  string
		wantLogPart string
	}{
		{
			name:    "прогулки: ошибочные записи пропускаются",
			dataset: []string{"6000,1h", "не число,1h", "0,1h", "3000,30m"},
			makeParser: func() DataParser {
				return &daysteps.DaySteps{Personal: person}
			},
			wantOutput: "Количество шагов: 6000.\nДистанция составила 4.72 км.\nВы сожгли 177.19 ккал.\n" +
				"Количество шагов: 3000.\nДистанция составила 2.36 км.\nВы сожгли 88.59 ккал.\n",
			wantLogPart: `ошибка разбора строки "0,1h"`,
		},
		{
			name:    "тренировки: неизвестный тип не прерывает журнал",
			dataset: []string{"6000,Ходьба,1h", "6000,Плавание,1h", "неверная запись", "6000,Бег,1h"},
			makeParser: func() DataParser {
				return &trainings.Training{Personal: person}
			},
			wantOutput: "Тип тренировки: Ходьба\nДлительность: 1.00 ч.\nДистанция: 4.72 км.\nСкорость: 4.72 км/ч\nСожгли калорий: 177.19\n" +
				"Тип тренировки: Бег\nДлительность: 1.00 ч.\nДистанция: 4.72 км.\nСкорость: 4.72 км/ч\nСожгли калорий: 354.38\n",
			wantLogPart: "неизвестный тип тренировки",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var logBuf bytes.Buffer
			log.SetOutput(&logBuf)
			defer log.SetOutput(os.Stderr)

			output := captureStdout(t, func() {
				Info(tt.dataset, tt.makeParser())
			})

			assert.Equal(t, tt.wantOutput, output)
			assert.Contains(t, logBuf.String(), tt.wantLogPart)
		})
	}
}

func TestParseKeepsPersonalDataAndPreviousStateOnError(t *testing.T) {
	person := personaldata.Personal{Name: "Тест", Weight: 75, Height: 1.75}

	t.Run("прогулка", func(t *testing.T) {
		ds := daysteps.DaySteps{Personal: person, Steps: 100, Duration: time.Minute}

		require.Error(t, ds.Parse("200,неверная продолжительность"))
		assert.Equal(t, 100, ds.Steps)
		assert.Equal(t, time.Minute, ds.Duration)
		assert.Equal(t, person, ds.Personal)

		require.NoError(t, ds.Parse("200,2h"))
		assert.Equal(t, 200, ds.Steps)
		assert.Equal(t, 2*time.Hour, ds.Duration)
		assert.Equal(t, person, ds.Personal)
	})

	t.Run("тренировка", func(t *testing.T) {
		tr := trainings.Training{Personal: person, Steps: 100, TrainingType: "Бег", Duration: time.Minute}

		require.Error(t, tr.Parse("200,Ходьба,неверная продолжительность"))
		assert.Equal(t, 100, tr.Steps)
		assert.Equal(t, "Бег", tr.TrainingType)
		assert.Equal(t, time.Minute, tr.Duration)
		assert.Equal(t, person, tr.Personal)

		require.NoError(t, tr.Parse("200,Ходьба,2h"))
		assert.Equal(t, 200, tr.Steps)
		assert.Equal(t, "Ходьба", tr.TrainingType)
		assert.Equal(t, 2*time.Hour, tr.Duration)
		assert.Equal(t, person, tr.Personal)
	})
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	old := os.Stdout
	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = w
	defer func() {
		os.Stdout = old
		_ = w.Close()
		_ = r.Close()
	}()

	fn()

	require.NoError(t, w.Close())
	var buf bytes.Buffer
	_, err = buf.ReadFrom(r)
	require.NoError(t, err)

	return buf.String()
}
