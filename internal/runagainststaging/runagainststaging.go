package runagainststaging

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/splice/platform/localdev/v2/localdev/internal/command"
	"github.com/splice/platform/localdev/v2/localdev/internal/files"
	"github.com/splice/platform/localdev/v2/localdev/internal/signal"
	"os"
)

func NewRunner(rebuild bool) *Runner {
	return &Runner{
		rebuild: rebuild,
	}
}

type Runner struct {
	rebuild bool
}

func (r *Runner) Run(ctx context.Context) error {
	if !files.Exists("docker-compose-staging.yml") {
		return errors.New("no docker-compose-staging.yml")
	}

	tunnelPrimary := command.NewShellCommand("ssh", "-N", "-L", "57526:splicestaging-primary-2020820.ckfmzkktjzob.us-west-1.rds.amazonaws.com:3306", "ned")
	tunnelReplica := command.NewShellCommand("ssh", "-N", "-L", "57527:splicestaging-replica-2020820.ckfmzkktjzob.us-west-1.rds.amazonaws.com:3306", "ned")

	fmt.Println("Starting database tunnel to primary")
	if err := tunnelPrimary.Start(ctx); err != nil {
		return errors.New("failed to start database tunnel to primary")
	}

	fmt.Println("Starting database tunnel to replica")
	if err := tunnelReplica.Start(ctx); err != nil {
		return errors.New("failed to start database tunnel to replica")
	}

	envPath := os.Getenv("SPLICE_PLATFORM")
	envFile := fmt.Sprintf("%s/private.env", envPath)

	if r.rebuild {
		dockerComposeBuild := command.NewDockerComposeCommand("--env-file", envFile, "-f", "docker-compose-staging.yml", "build")
		if err := dockerComposeBuild.Run(ctx); err != nil {
			return errors.New("failed to build docker compose environment using docker-compose-staging.yml")
		}
	}

	dockerCompose := command.NewDockerComposeCommand("--env-file", envFile, "-f", "docker-compose-staging.yml", "up")
	if err := dockerCompose.Start(ctx); err != nil {
		return errors.New("failed to start docker compose environment using docker-compose-staging.yml")
	}

	signal.WaitForSignal(ctx, dockerCompose, tunnelReplica, tunnelPrimary)

	return nil
}
