package migration

import (
	"fmt"
	"strings"
	"testing"
)

func Test_getMigrationsDir(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{"Success", "migrations", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := getMigrationsDir()
			parts := strings.Split(path, "/")
			last := parts[len(parts)-1]
			fmt.Println("/" + last)
			if (err != nil) != tt.wantErr {
				t.Errorf("getMigrationsDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if last != tt.want {
				t.Errorf("getMigrationsDir() got = %v, want %v", last, tt.want)
			}
		})
	}
}
