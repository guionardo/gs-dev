package console

import (
	"math/rand"
	"testing"
)

func randomInts() []interface{} {
	r := make([]interface{}, 1)
	r[0] = rand.Int()
	return r
}

func TestOutput(t *testing.T) {
	type args struct {
		printLevel OutputLevel
		format     string
		args       []interface{}
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "INFO",
			args: args{INFO, "Test %d", randomInts()},
		}, {
			name: "WARNING",
			args: args{WARN, "Warning %d", randomInts()},
		}, {
			name: "SUCCESS",
			args: args{SUCCESS, "Success %d", randomInts()},
		}, {
			name: "ERROR",
			args: args{ERROR, "Error %d", randomInts()},
		}, {
			name: "NONE",
			args: args{NONE, "None %d", randomInts()},
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Output(tt.args.printLevel, tt.args.format, tt.args.args...)
		})
	}
}
