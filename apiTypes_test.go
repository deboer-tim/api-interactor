package interactor

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// Ensure all types adhere to required interfaces
var (
	_ ApiType = new(Contest)
	_ ApiType = new(Problem)
	_ ApiType = new(Submission)

	_ json.Unmarshaler = new(ApiTime)
	_ JSONUnmarshaller = new(ApiTime)
	_ fmt.Stringer     = new(ApiTime)

	_ json.Unmarshaler = new(ApiRelTime)
	_ JSONUnmarshaller = new(ApiRelTime)
	_ fmt.Stringer     = new(ApiRelTime)
)

func TestApiTime_UnmarshalJSON(t *testing.T) {
	ti := struct {
		T ApiTime
	}{}

	// Supported formats
	formats := []string{"2006-01-02T15:04:05Z07", time.RFC3339}
	now := time.Now()

	jsonFormat := `{"T": "%v"}`
	for _, f := range formats {
		t.Run("format-"+f, func(t *testing.T) {
			jsonString := fmt.Sprintf(jsonFormat, now.Format(f))
			assert.Nil(t, json.Unmarshal([]byte(jsonString), &ti))
			assert.EqualValues(t, now.Truncate(time.Second).UnixNano(), ti.T.Time().UnixNano())

			// Also test when the value is null
			jsonString = `{"T": null}`
			assert.Nil(t, json.Unmarshal([]byte(jsonString), &ti))
			assert.EqualValues(t, time.Time{}.UnixNano(), ti.T.Time().UnixNano())
		})
	}
}

func TestApiRelTime_UnmarshalJSON(t *testing.T) {
	ti := struct {
		T ApiRelTime
	}{}

	// Only one format is allowed, test a single value. TODO perhaps some form of fuzz testing here?
	jsonString := `{"T": "0:03:38.749"}`
	duration := time.Minute*3 + time.Second*38 + time.Millisecond*749
	assert.Nil(t, json.Unmarshal([]byte(jsonString), &ti))
	assert.EqualValues(t, duration, ti.T.Duration())

	// Test null
	jsonString = `{"T": null}`
	assert.Nil(t, json.Unmarshal([]byte(jsonString), &ti))
	assert.EqualValues(t, time.Duration(0), ti.T.Duration())
}