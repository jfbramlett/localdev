package command

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Command is a wrapper around exec.Command providing some convenience methods for how we are using the commands.
type Command interface {
	// Start starts the command in the background
	Start(ctx context.Context) error

	// Run runs the command and waits for it to complete
	Run(ctx context.Context) error

	// Stop stops the command started by Start
	Stop(ctx context.Context) error

	// AddEnvironment adds a new key=value to the environment for this command
	AddEnvironment(key, value string)
}

// NewShellCommand creates a new shell command mapping stdout to our current stdout (same with stderr)
func NewShellCommand(name string, args ...string) Command {
	command := &shellCommand{cmd: exec.Command(name, args...)}
	command.cmd.Env = os.Environ()
	command.cmd.Stdout = os.Stdout
	command.cmd.Stderr = os.Stderr
	command.cmd.Stdin = os.Stdin
	return command
}

// shellCommand represents a basic shell or script command.
type shellCommand struct {
	cmd *exec.Cmd
}

func (c *shellCommand) Start(_ context.Context) error {
	return c.cmd.Start()
}

func (c *shellCommand) Run(_ context.Context) error {
	return c.cmd.Run()
}

func (c *shellCommand) Stop(_ context.Context) error {
	return c.cmd.Process.Kill()
}

func (c *shellCommand) AddEnvironment(key, value string) {
	c.cmd.Env = append(c.cmd.Env, fmt.Sprintf("%s=%s", key, value))
}

func (c *shellCommand) String() string {
	return c.cmd.String()
}

// NewDockerComposeCommand creates a new docker-compose command (name is left off and will be added automatically by
// this)
func NewDockerComposeCommand(args ...string) Command {
	composeFile := "docker-compose.yml"
	for idx, a := range args {
		if a == "-f" {
			composeFile = args[idx+1]
			break
		}
	}
	command := &dockerComposeCommand{cmd: exec.Command("docker-compose", args...), composeFile: composeFile}
	command.cmd.Env = os.Environ()
	command.cmd.Stdout = os.Stdout
	command.cmd.Stderr = os.Stderr
	return command
}

// dockerComposeCommand represents a new request of docker-compose. We separate this from the traditional shell command
// because to stop what is run by docker-compose we need to do `docker-compose down` rather than killing the process.
type dockerComposeCommand struct {
	cmd         *exec.Cmd
	composeFile string
}

func (dc *dockerComposeCommand) Start(ctx context.Context) error {
	err := dc.cmd.Start()
	if err != nil {
		return err
	}

	if err := dc.waitForStarted(ctx); err != nil {
		return err
	}
	return nil
}

func (dc *dockerComposeCommand) Run(_ context.Context) error {
	return dc.cmd.Run()
}

func (dc *dockerComposeCommand) Stop(ctx context.Context) error {
	dockerCompose := NewShellCommand("docker-compose", "-f", dc.composeFile, "down")
	return dockerCompose.Run(ctx)
}

func (dc *dockerComposeCommand) String() string {
	return dc.cmd.String()
}

func (dc *dockerComposeCommand) AddEnvironment(key, value string) {
	dc.cmd.Env = append(dc.cmd.Env, fmt.Sprintf("%s=%s", key, value))
}

func (dc *dockerComposeCommand) waitForStarted(_ context.Context) error {
	containerListCmd := exec.Command("bash", "-c", `docker ps | awk '{if(NR>1) print $NF}'`)
	containerListCmd.Env = os.Environ()
	output, err := containerListCmd.Output()
	if err != nil {
		return err
	}

	containers := strings.Split(string(output), "\n")
	for _, container := range containers {
		container = strings.TrimSpace(container)
		if container == "" {
			continue
		}
		for {
			statusCmd := exec.Command("docker", "inspect", "--format", `"{{.State.Health.Status}}"`, container)
			statusCmd.Env = os.Environ()
			status, err := statusCmd.Output()
			if err != nil {
				return err
			}
			statuses := strings.Split(string(status), "\n")
			currentStatus := strings.ReplaceAll(statuses[0], "\"", "")

			if currentStatus == "healthy" {
				break
			}
			time.Sleep(2 * time.Second)
		}
	}

	return nil
}
