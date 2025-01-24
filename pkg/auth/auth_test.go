package auth

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetToken(t *testing.T) {
	token, err := GetToken()
	assert.NoError(t, err)
	assert.Equal(t, sessionID, token)
}

func TestGetAuthURL(t *testing.T) {
	tests := []struct {
		name     string
		testMode string
		want     string
	}{
		{
			name:     "Production mode",
			testMode: "",
			want:     "https://app.nexlayer.io/auth/github",
		},
		{
			name:     "Test mode",
			testMode: "true",
			want:     "https://app.staging.nexlayer.io/auth/github",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldTestMode := os.Getenv("NEXLAYER_TEST_MODE")
			os.Setenv("NEXLAYER_TEST_MODE", tt.testMode)
			defer os.Setenv("NEXLAYER_TEST_MODE", oldTestMode)

			got := GetAuthURL()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSaveToken(t *testing.T) {
	// Create temporary directory for test
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	// Test saving token
	testToken := "test-token"
	err := SaveToken(testToken)
	assert.NoError(t, err)

	// Verify token was saved correctly
	tokenPath := filepath.Join(tmpDir, configDir, tokenFile)
	data, err := os.ReadFile(tokenPath)
	assert.NoError(t, err)

	var config TokenConfig
	err = json.Unmarshal(data, &config)
	assert.NoError(t, err)
	assert.Equal(t, testToken, config.Token)

	// Test saving empty token
	err = SaveToken("")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token cannot be empty")
}

func TestLogin(t *testing.T) {
	// Test login in test mode
	oldTestMode := os.Getenv("NEXLAYER_TEST_MODE")
	os.Setenv("NEXLAYER_TEST_MODE", "true")
	defer os.Setenv("NEXLAYER_TEST_MODE", oldTestMode)

	err := Login()
	assert.NoError(t, err)
}
