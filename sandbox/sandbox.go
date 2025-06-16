package sandbox

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/gin-gonic/gin"
)

type Sandbox struct {
	client *client.Client
}

type ContainerRequest struct {
	Name        string            `json:"name"`
	Image       string            `json:"image"`
	Environment map[string]string `json:"environment"`
	Ports       map[string]string `json:"ports"`
}

type ContainerInfo struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Image  string   `json:"image"`
	Status string   `json:"status"`
	Ports  []string `json:"ports"`
	URL    string   `json:"url,omitempty"`
}

func NewSandbox() (*Sandbox, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return &Sandbox{client: cli}, nil
}

func (s *Sandbox) CreateContainer(c *gin.Context) {
	var req ContainerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Default to webtop if no image specified
	if req.Image == "" {
		req.Image = "lscr.io/linuxserver/webtop:latest"
	}

	// Pull image first
	reader, err := s.client.ImagePull(context.Background(), req.Image, image.PullOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer func(reader io.ReadCloser) {
		_ = reader.Close()
	}(reader)

	if _, err := io.Copy(io.Discard, reader); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Setup environment variables
	env := []string{
		"PUID=1000",
		"PGID=1000",
		"TZ=UTC",
		"SUBFOLDER=/",
		"KEYBOARD=en-us-qwerty",
	}
	for key, value := range req.Environment {
		env = append(env, fmt.Sprintf("%s=%s", key, value))
	}

	// Setup port bindings
	exposedPorts := nat.PortSet{}
	portBindings := nat.PortMap{}

	// Default webtop port
	containerPort := "3000/tcp"
	exposedPorts[nat.Port(containerPort)] = struct{}{}

	if req.Ports["3000"] != "" {
		hostPort := req.Ports["3000"]
		portBindings[nat.Port(containerPort)] = []nat.PortBinding{
			{HostPort: hostPort},
		}
	} else {
		// Auto-assign port
		portBindings[nat.Port(containerPort)] = []nat.PortBinding{
			{HostPort: "0"},
		}
	}

	// Container configuration
	config := &container.Config{
		Image:        req.Image,
		Env:          env,
		ExposedPorts: exposedPorts,
	}

	hostConfig := &container.HostConfig{
		PortBindings: portBindings,
		SecurityOpt:  []string{"seccomp=unconfined"},
		ShmSize:      1073741824, // 1GB
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: fmt.Sprintf("/tmp/%s-config", req.Name),
				Target: "/config",
			},
		},
		RestartPolicy: container.RestartPolicy{Name: "unless-stopped"},
	}

	// Create container
	resp, err := s.client.ContainerCreate(
		context.Background(),
		config,
		hostConfig,
		nil,
		nil,
		req.Name,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      resp.ID,
		"message": "Container created successfully",
	})
}

func (s *Sandbox) ListContainers(c *gin.Context) {
	containers, err := s.client.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var containerInfos []ContainerInfo
	for _, item := range containers {
		info := ContainerInfo{
			ID:     item.ID[:12],
			Name:   item.Names[0][1:], // Remove leading slash
			Image:  item.Image,
			Status: item.Status,
			Ports:  []string{},
		}
		// Extract port information and generate URL
		for _, port := range item.Ports {
			if port.PublicPort > 0 {
				portStr := fmt.Sprintf("%d:%d", port.PublicPort, port.PrivatePort)
				info.Ports = append(info.Ports, portStr)
				// Generate access URL for webtop containers
				if port.PrivatePort == 3000 {
					info.URL = fmt.Sprintf("http://localhost:%d", port.PublicPort)
				}
			}
		}
		containerInfos = append(containerInfos, info)
	}

	c.JSON(http.StatusOK, containerInfos)
}

func (s *Sandbox) StartContainer(c *gin.Context) {
	containerID := c.Param("id")

	err := s.client.ContainerStart(context.Background(), containerID, container.StartOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Container started successfully"})
}

func (s *Sandbox) StopContainer(c *gin.Context) {
	containerID := c.Param("id")
	timeout := 30

	err := s.client.ContainerStop(context.Background(), containerID, container.StopOptions{Timeout: &timeout})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Container stopped successfully"})
}

func (s *Sandbox) RemoveContainer(c *gin.Context) {
	containerID := c.Param("id")

	err := s.client.ContainerRemove(context.Background(), containerID, container.RemoveOptions{Force: true})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Container removed successfully"})
}

func (s *Sandbox) GetLogs(c *gin.Context) {
	containerID := c.Param("id")

	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       c.DefaultQuery("lines", "100"),
	}

	logs, err := s.client.ContainerLogs(context.Background(), containerID, options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer func(logs io.ReadCloser) {
		_ = logs.Close()
	}(logs)

	c.Header("Content-Type", "text/plain")

	if _, err := io.Copy(c.Writer, logs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}
