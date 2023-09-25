/*
 license x
*/

package repository

import (
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
	for _, tp := range tests {
		tt := tp
		t.Run(tt.name, func(t *testing.T) {
			if tt.newURL != "" {
				t.Setenv(env, tt.newURL)
			}
			if got := DBURL(); got != tt.want {
				t.Errorf("DBURL() = %v, want %v", got, tt.want)
			}
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
	for _, tp := range tests {
		tt := tp
		t.Run(tt.name, func(t *testing.T) {

			if tt.new != "" {
				t.Setenv(env, tt.new)
			}
			if got := DBDriver(); got != tt.want {
				t.Errorf("DBDriver() = %v, want %v", got, tt.want)
			}
		})
	}
}
