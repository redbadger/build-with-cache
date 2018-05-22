package root

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

func Test_parseDockerfile(t *testing.T) {
	type args struct {
		reader  io.Reader
		imgName string
		cache   string
	}
	tests := []struct {
		name       string
		args       args
		wantStages []string
		wantImages map[string]string
		wantErr    bool
	}{
		{
			"named stages",
			args{strings.NewReader(`
FROM golang:alpine as builder

RUN echo xx > /hello

FROM alpine as release

COPY --from=builder /hello hello
`),
				"host.docker.internal:5000/redbadger/deploy:latest", ""},
			[]string{"builder", "release"},
			map[string]string{
				"builder": "host.docker.internal:5000/redbadger/deploy-builder",
				"release": "host.docker.internal:5000/redbadger/deploy-release",
			},
			false,
		},
		{
			"unnamed stages",
			args{strings.NewReader(`
FROM golang:alpine

RUN echo xx > /hello

FROM alpine

COPY --from=0 /hello hello
`),
				"host.docker.internal:5000/redbadger/deploy", ""},
			[]string{"0", "1"},
			map[string]string{
				"0": "host.docker.internal:5000/redbadger/deploy-0",
				"1": "host.docker.internal:5000/redbadger/deploy-1",
			},
			false,
		},
		{
			"mixed naming with cache",
			args{strings.NewReader(`
FROM golang:alpine as builder

RUN echo xx > /hello

FROM alpine

COPY --from=builder /hello hello
`),
				"host.docker.internal:5000/redbadger/deploy:latest", "localhost:5001"},
			[]string{"builder", "1"},
			map[string]string{
				"builder": "localhost:5001/redbadger/deploy-builder",
				"1":       "localhost:5001/redbadger/deploy-1",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStages, gotImages, err := parseDockerfile(tt.args.reader, tt.args.imgName, tt.args.cache)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDockerfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotStages, tt.wantStages) {
				t.Errorf("parseDockerfile() = %v, want %v", gotStages, tt.wantStages)
			}
			if !reflect.DeepEqual(gotImages, tt.wantImages) {
				t.Errorf("parseDockerfile() = %v, want %v", gotImages, tt.wantImages)
			}
		})
	}
}
