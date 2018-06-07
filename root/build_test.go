package root

import (
	"reflect"
	"testing"
)

func Test_createCommand(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"creates build command when using space",
			args{
				[]string{
					"--cache", "myregistry",
					"--file", "Dockerfile",
					"--tag", "mytag",
					"--build-arg", "x=1",
					".",
				},
			},
			[]string{"docker", "build", "--file", "Dockerfile", "--tag", "mytag", "--build-arg", "x=1", "."},
		},
		{
			"creates build command when using equals sign",
			args{
				[]string{
					".",
					"--cache=myregistry",
					"--file=Dockerfile",
					"--tag=mytag",
					"--build-arg=x=1",
				},
			},
			[]string{"docker", "build", ".", "--file=Dockerfile", "--tag=mytag", "--build-arg=x=1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createCommand(tt.args.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}
