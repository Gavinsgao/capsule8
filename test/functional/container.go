// Copyright 2017 Capsule8, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package functional

import (
	"os/exec"
	"strings"
	"testing"

	"golang.org/x/net/context"
)

// Container represents a running container
type Container struct {
	t       *testing.T
	Path    string
	ImageID string
	command *exec.Cmd
}

func (c *Container) dockerBuildArgs(quiet bool, buildargs []string) []string {
	args := []string{"build"}
	if quiet {
		args = append(args, "-q")
	}
	args = append(args, buildargs...)
	args = append(args, c.Path)
	return args
}

// Build executes docker build
func (c *Container) Build(buildargs ...string) error {
	dockerArgs := c.dockerBuildArgs(false, buildargs)
	docker := exec.Command("docker", dockerArgs...)
	err := docker.Run()
	if err != nil {
		return err
	}

	dockerArgs = c.dockerBuildArgs(true, buildargs)
	docker = exec.Command("docker", dockerArgs...)
	dockerOutput, err := docker.Output()
	if err != nil {
		return err
	}

	trimmed := strings.TrimSpace(string(dockerOutput))
	c.ImageID = strings.TrimPrefix(trimmed, "sha256:")

	return nil
}

func (c *Container) dockerRunArgs(runargs []string) []string {
	args := append([]string{"run", "--rm"}, runargs...)
	args = append(args, c.ImageID)
	return args
}

// Start executes docker start
func (c *Container) Start(runargs ...string) error {
	dockerArgs := c.dockerRunArgs(runargs)
	c.command = exec.Command("docker", dockerArgs...)
	return c.command.Start()
}

// StartContext executes docker start with a context object
func (c *Container) StartContext(ctx context.Context, runargs ...string) error {
	dockerArgs := c.dockerRunArgs(runargs)
	c.command = exec.CommandContext(ctx, "docker", dockerArgs...)
	return c.command.Start()
}

// Wait executes docker wait
func (c *Container) Wait() error {
	return c.command.Wait()
}

// Run executes docker run
func (c *Container) Run(runargs ...string) error {
	dockerArgs := c.dockerRunArgs(runargs)
	c.command = exec.Command("docker", dockerArgs...)
	return c.command.Run()
}

// RunContext executes docker run with a context object
func (c *Container) RunContext(ctx context.Context, runargs ...string) error {
	dockerArgs := c.dockerRunArgs(runargs)
	c.command = exec.CommandContext(ctx, "docker", dockerArgs...)
	return c.command.Run()
}

// NewContainer returns a new container testing object
func NewContainer(t *testing.T, path string) *Container {
	return &Container{
		t:    t,
		Path: path,
	}
}
