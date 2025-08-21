package docker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateDockerImageName(t *testing.T) {
	tests := []struct {
		name     string
		image    string
		expected bool
	}{
		// Valid image names (following Docker and OCI standards)
		{"Valid image name without registry or tag", "ubuntu", true},
		{"Valid image name with namespace", "library/ubuntu", true},
		{"Valid image name with registry", "docker.io/library/ubuntu", true},
		{"Valid image name with custom registry and port", "localhost:5000/myproject/ubuntu", true},
		{"Valid image name with tag", "myregistry/myproject/ubuntu:latest", true},
		{"Valid image name with version tag", "nginx:1.21.0", true},
		{"Valid image name with custom tag", "myapp:v1.0.0-beta", true},
		{"Valid image name with underscore", "my_project/app:latest", true},
		{"Valid image name with multiple slashes", "org/team/project/app:latest", true},
		{"Valid image name with IP address registry", "192.168.1.100:5000/app:latest", true},

		// Invalid image names
		{"Invalid image name with uppercase letters in repository", "MyRegistry/MyProject/Ubuntu", false},
		{"Invalid image name with invalid characters", "invalid!@#/image/name", false},
		{"Invalid image name with empty string", "", false},
		{"Invalid image name with only slashes", "///", false},
		{"Invalid image name with multiple @ symbols", "invalid@image@name", false},
		{"Invalid image name with invalid tag characters", "app:tag:with:colons", false},
		{"Invalid image name with leading slash", "/app:latest", false},
		{"Invalid image name with trailing slash", "app/:latest", false},
		{"Invalid image name with consecutive slashes", "app//name:latest", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DockerImageValidator(tt.image)
			assert.Equal(t, tt.expected, result, "Test case: %s", tt.name)
		})
	}
}
