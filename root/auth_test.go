package root

import "testing"

func Test_getRegistryHostnameFromName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name         string
		args         args
		wantHostname string
		wantErr      bool
	}{
		{
			"full canonical name",
			args{"myregistry.com:5000/namespace/name:tag"},
			"myregistry.com:5000",
			false,
		},
		{
			"full canonical name without port",
			args{"myregistry.com/namespace/name:tag"},
			"myregistry.com",
			false,
		},
		{
			"canonical name without tag",
			args{"myregistry.com/name"},
			"myregistry.com",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHostname, err := getRegistryHostnameFromName(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("getRegistryHostnameFromName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotHostname != tt.wantHostname {
				t.Errorf("getRegistryHostnameFromName() = %v, want %v", gotHostname, tt.wantHostname)
			}
		})
	}
}
