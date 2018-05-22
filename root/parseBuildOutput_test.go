package root

import (
	"reflect"
	"testing"
)

func Test_parseBuildOutput(t *testing.T) {
	type args struct {
		buildOutput string
		imgTag      string
		cache       string
		stages      []string
	}
	tests := []struct {
		name     string
		args     args
		wantRefs map[string]string
		wantErr  bool
	}{
		{
			"simple output",
			args{
				`
Sending build context to Docker daemon  7.222MB
Step 1/4 : FROM golang:alpine as builder
 ---> 05fe62871090
Step 2/4 : RUN echo xx > /hello
 ---> Using cache
 ---> cf811a451aed
Step 3/4 : FROM alpine as release
 ---> 3fd9065eaf02
Step 4/4 : COPY --from=builder /hello hello
 ---> Using cache
 ---> d7c5da99a676
Successfully built d7c5da99a676
Successfully tagged host.docker.internal:5000/redbadger/deploy:latest
`,
				"host.docker.internal:5000/redbadger/deploy",
				"",
				[]string{
					"builder",
					"release",
				},
			},
			map[string]string{
				"cf811a451aed": "host.docker.internal:5000/redbadger/deploy-builder",
				"d7c5da99a676": "host.docker.internal:5000/redbadger/deploy-release",
			},
			false,
		},
		{
			"final stage not named",
			args{
				`
Sending build context to Docker daemon  7.222MB
Step 1/4 : FROM golang:alpine as builder
 ---> 05fe62871090
Step 2/4 : RUN echo xx > /hello
 ---> Using cache
 ---> cf811a451aed
Step 3/4 : FROM alpine
 ---> 3fd9065eaf02
Step 4/4 : COPY --from=builder /hello hello
 ---> Using cache
 ---> d7c5da99a676
Successfully built d7c5da99a676
Successfully tagged host.docker.internal:5000/redbadger/deploy:latest
`,
				"host.docker.internal:5000/redbadger/deploy",
				"",
				[]string{
					"builder",
					"1",
				},
			},
			map[string]string{
				"cf811a451aed": "host.docker.internal:5000/redbadger/deploy-builder",
				"d7c5da99a676": "host.docker.internal:5000/redbadger/deploy-1",
			},
			false,
		},
		{
			"no stages are named",
			args{
				`
Sending build context to Docker daemon  7.222MB
Step 1/4 : FROM golang:alpine
 ---> 05fe62871090
Step 2/4 : RUN echo xx > /hello
 ---> Using cache
 ---> cf811a451aed
Step 3/4 : FROM alpine
 ---> 3fd9065eaf02
Step 4/4 : COPY --from=0 /hello hello
 ---> Using cache
 ---> d7c5da99a676
Successfully built d7c5da99a676
Successfully tagged host.docker.internal:5000/redbadger/deploy:latest
`,
				"host.docker.internal:5000/redbadger/deploy",
				"",
				[]string{
					"0",
					"1",
				},
			},
			map[string]string{
				"cf811a451aed": "host.docker.internal:5000/redbadger/deploy-0",
				"d7c5da99a676": "host.docker.internal:5000/redbadger/deploy-1",
			},
			false,
		},
		{
			"cache registry supplied",
			args{
				`
Sending build context to Docker daemon  7.222MB
Step 1/4 : FROM golang:alpine
 ---> 05fe62871090
Step 2/4 : RUN echo xx > /hello
 ---> Using cache
 ---> cf811a451aed
Step 3/4 : FROM alpine
 ---> 3fd9065eaf02
Step 4/4 : COPY --from=0 /hello hello
 ---> Using cache
 ---> d7c5da99a676
Successfully built d7c5da99a676
Successfully tagged host.docker.internal:5000/redbadger/deploy:latest
`,
				"host.docker.internal:5000/redbadger/deploy",
				"localhost:5000",
				[]string{
					"0",
					"1",
				},
			},
			map[string]string{
				"cf811a451aed": "localhost:5000/redbadger/deploy-0",
				"d7c5da99a676": "localhost:5000/redbadger/deploy-1",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRefs, err := parseBuildOutput(tt.args.buildOutput, tt.args.imgTag, tt.args.cache, tt.args.stages)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseBuildOutput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRefs, tt.wantRefs) {
				t.Errorf("parseBuildOutput() = %v, want %v", gotRefs, tt.wantRefs)
			}
		})
	}
}
