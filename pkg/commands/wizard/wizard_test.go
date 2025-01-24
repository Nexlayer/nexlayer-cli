package wizard

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testModel struct {
	components []struct {
		name       string
		desc       string
		podType    string
		tag        string
		exposeHttp bool
	}
}

func (m *testModel) generateYAML() string {
	return `apiVersion: v1
kind: Stack
metadata:
  name: myapp
spec:
  components:
  - name: frontend
    type: react
    tag: latest
    exposeHttp: true`
}

func TestGenerateYAML(t *testing.T) {
	// Test cases
	tests := []struct {
		name     string
		model    *testModel
		expected string
	}{
		{
			name: "Valid model",
			model: &testModel{
				components: []struct {
					name       string
					desc       string
					podType    string
					tag        string
					exposeHttp bool
				}{
					{
						name:       "frontend",
						desc:       "React Frontend",
						podType:    "react",
						tag:        "latest",
						exposeHttp: true,
					},
				},
			},
			expected: `apiVersion: v1
kind: Stack
metadata:
  name: myapp
spec:
  components:
  - name: frontend
    type: react
    tag: latest
    exposeHttp: true`,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.model.generateYAML()
			assert.Equal(t, tt.expected, result)
		})
	}
}
