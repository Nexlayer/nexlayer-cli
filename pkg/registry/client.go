package registry

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Client handles Docker registry operations
type Client struct {
	config *RegistryConfig
}

// NewClient creates a new registry client
func NewClient(config *RegistryConfig) *Client {
	return &Client{
		config: config,
	}
}

// Login authenticates with the container registry
func (c *Client) Login() error {
	cmd := exec.Command("docker", "login", c.config.Registry,
		"-u", c.config.Username,
		"--password-stdin")
	
	cmd.Stdin = strings.NewReader(c.config.Token)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	return cmd.Run()
}

// BuildImage builds a Docker image for a service
func (c *Client) BuildImage(cfg ImageConfig) error {
	for _, tag := range cfg.Tags {
		imageTag := fmt.Sprintf("%s/%s/%s:%s", 
			cfg.Namespace,
			filepath.Base(cfg.Namespace), // organization name
			cfg.ServiceName,
			tag)
		
		cmd := exec.Command("docker", "build",
			"-t", imageTag,
			cfg.Path)
		
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to build image %s: %w", imageTag, err)
		}
	}
	
	return nil
}

// PushImage pushes a Docker image to the registry
func (c *Client) PushImage(cfg ImageConfig) error {
	for _, tag := range cfg.Tags {
		imageTag := fmt.Sprintf("%s/%s/%s:%s",
			cfg.Namespace,
			filepath.Base(cfg.Namespace), // organization name
			cfg.ServiceName,
			tag)
		
		cmd := exec.Command("docker", "push", imageTag)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to push image %s: %w", imageTag, err)
		}
	}
	
	return nil
}

// BuildAndPushImages builds and pushes multiple images
func (c *Client) BuildAndPushImages(cfg BuildConfig) error {
	for _, img := range cfg.Images {
		// Apply global namespace and tags if not set specifically for the image
		if img.Namespace == "" {
			img.Namespace = cfg.Namespace
		}
		if len(img.Tags) == 0 {
			img.Tags = cfg.Tags
		}
		
		if err := c.BuildImage(img); err != nil {
			return err
		}
		
		if err := c.PushImage(img); err != nil {
			return err
		}
	}
	
	return nil
}
