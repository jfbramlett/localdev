package runtest

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/splice/platform/localdev/v2/localdev/internal/command"
	"github.com/splice/platform/localdev/v2/localdev/internal/files"
	"os"
	"strings"
)

const (
	TestTypeUnit        = "unit"
	TestTypeIntegration = "integration"
	TestTypeAll         = "all"
)

func NewRunTest(testType string) *RunTest {
	return &RunTest{
		testType: testType,
	}
}

type RunTest struct {
	testType string
}

func (r *RunTest) Run(ctx context.Context) error {
	switch r.testType {
	case TestTypeUnit:
		return r.runUnitTest(ctx)
	case TestTypeIntegration:
		return r.runIntegrationTest(ctx)
	case TestTypeAll:
		if err := r.runUnitTest(ctx); err != nil {
			return err
		}
		if err := r.runIntegrationTest(ctx); err != nil {
			return err
		}
		return nil
	}

	return errors.New("invalid test type specified")
}

func (r *RunTest) runUnitTest(ctx context.Context) error {
	fmt.Println("running unit tests")
	return r.runTest(ctx, "unit")
}

func (r *RunTest) runIntegrationTest(ctx context.Context) error {
	fmt.Println("running integration tests")
	return r.runTest(ctx, "integration")
}

func (r *RunTest) runTest(ctx context.Context, buildTag string) error {
	tags := fmt.Sprintf("-tags=%s", buildTag)
	coverProfile := fmt.Sprintf("-coverprofile=./%s_coverage.out", buildTag)
	envFile := fmt.Sprintf("%s_test.env", buildTag)
	dockerFile := fmt.Sprintf("docker-compose-%s.yml", buildTag)
	dockerSeedFile := fmt.Sprintf("docker-compose-%s-seed.yml", buildTag)

	envPath := os.Getenv("SPLICE_PLATFORM")
	dockerEnvFile := fmt.Sprintf("%s/private.env", envPath)

	if files.Exists(dockerFile) {
		dockerCompose := command.NewDockerComposeCommand("--env-file", dockerEnvFile, "-f", dockerFile, "up", "-V")
		if err := dockerCompose.Start(ctx); err != nil {
			return errors.Wrapf(err, "failed starting %s", dockerFile)
		}
		defer func() {
			_ = dockerCompose.Stop(ctx)
		}()

		if files.Exists(dockerSeedFile) {
			dockerComposeSeed := command.NewDockerComposeCommand("--env-file", dockerEnvFile, "-f", dockerSeedFile, "up")
			if err := dockerComposeSeed.Run(ctx); err != nil {
				return errors.Wrapf(err, "failed running %s", dockerSeedFile)
			}
		}

	}

	testCommand := command.NewShellCommand("gotestsum", "--format", "pkgname", "--", tags, coverProfile, "./...")

	if files.Exists(envFile) {
		testEnv, err := files.ReadAll(envFile)
		if err != nil {
			return errors.Wrapf(err, "failed loading our environment file %s", envFile)
		}
		envVars := strings.Split(testEnv, "\n")
		for _, envVar := range envVars {
			envVar = strings.TrimSpace(envVar)
			if len(envVar) == 0 {
				continue
			}
			varParts := strings.Split(envVar, "=")
			if len(varParts) != 2 {
				continue
			}
			testCommand.AddEnvironment(strings.TrimSpace(varParts[0]), strings.TrimSpace(varParts[1]))
		}
	}

	err := testCommand.Run(ctx)
	if err != nil {
		return errors.Wrap(err, "failed executing tests")
	}

	return nil

}
