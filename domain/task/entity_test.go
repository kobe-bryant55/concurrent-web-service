package taskdomain

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/types"
)

func TestTask_SetStatus(t *testing.T) {
	testCases := []struct {
		name          string
		initialStatus types.Status
		newStatus     types.Status
		expectedErrs  int
	}{
		{
			name:          "Valid Status",
			initialStatus: types.Active,
			newStatus:     types.Passive,
			expectedErrs:  0,
		},
		{
			name:          "Invalid Status",
			initialStatus: types.Active,
			newStatus:     "InvalidStatus",
			expectedErrs:  1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			task := &Task{status: tc.initialStatus, description: "description", title: "title"}
			errs := task.SetStatus(tc.newStatus)

			assert.Len(t, errs, tc.expectedErrs)
		})
	}
}

func TestTask_Validate(t *testing.T) {
	testCases := []struct {
		name         string
		title        string
		description  string
		status       types.Status
		expectedErrs int
	}{
		{
			name:         "Valid Task",
			title:        "Valid Title",
			description:  "Valid Description",
			status:       types.Active,
			expectedErrs: 0,
		},
		{
			name:         "Short Description",
			title:        "Valid Title",
			description:  "de",
			status:       types.Active,
			expectedErrs: 1,
		},
		{
			name:         "Short Title",
			title:        "ti",
			description:  "description",
			status:       types.Active,
			expectedErrs: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			task := &Task{
				title:       tc.title,
				description: tc.description,
				status:      tc.status,
			}
			errs := task.validate()

			assert.Len(t, errs, tc.expectedErrs)
		})
	}
}
