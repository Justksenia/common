package cron

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type idleJob struct{}

func (*idleJob) Run() {}

func TestParseSpecification(t *testing.T) {
	testCases := []struct {
		name    string
		spec    string
		isError bool
	}{
		{
			name:    "valid spec by rfc",
			spec:    "*/10 * * *",
			isError: false,
		},
		{
			name:    "valid spec, every tenth second expected",
			spec:    "0/10",
			isError: false,
		},
		{
			name:    "valid spec by description",
			spec:    "@midnight",
			isError: false,
		},
		{
			name:    "invalid spec, exceed max value",
			spec:    "*/10000 * * * *",
			isError: true,
		},
		{
			name:    "invalid spec, random set of symbols",
			spec:    "fwfnwifn",
			isError: true,
		},
	}
	ij := &idleJob{}
	sch := NewCronScheduler()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := sch.AddJob(ij, tc.spec)
			assert.Equal(t, tc.isError, err != nil)
		})
	}
}
