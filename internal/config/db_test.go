/*
 license x
*/

package config

import (
	"os"
	"testing"
)

func TestDbURL(t *testing.T) {
	const env = "DB_URL"
	tests := []struct {
		name   string
		want   string
		newURL string
	}{
		{
			name:   "default value",
			newURL: "",
			want:   dbURLDefault,
		},
		{
			name:   "new value",
			newURL: "new URL value",
			want:   "new URL value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.newURL != "" {
				_ = os.Setenv(env, tt.newURL)
			}
			if got := DbURL(); got != tt.want {
				t.Errorf("DbURL() = %v, want %v", got, tt.want)
			}
			_ = os.Unsetenv(env)
		})
	}
}

func TestDbDriver(t *testing.T) {
	const env = "DB_DRIVER"

	tests := []struct {
		name string
		new  string
		want string
	}{
		{
			name: "default value",
			new:  "",
			want: dbDriverDefault,
		},
		{
			name: "new value",
			new:  "new driver value",
			want: "new driver value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.new != "" {
				_ = os.Setenv(env, tt.new)
			}
			if got := DbDriver(); got != tt.want {
				t.Errorf("DbDriver() = %v, want %v", got, tt.want)
			}
			_ = os.Unsetenv(env)
		})
	}
}
