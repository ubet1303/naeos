package strutil

import "testing"

func TestSlugify(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"Hello World", "hello-world"},
		{"my_project/name", "my-project-name"},
		{"  spaces  ", "spaces"},
		{"UPPER CASE", "upper-case"},
		{"already-slug", "already-slug"},
		{"a/b_c d", "a-b-c-d"},
		{"", ""},
	}
	for _, tt := range tests {
		got := Slugify(tt.input)
		if got != tt.want {
			t.Errorf("Slugify(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
