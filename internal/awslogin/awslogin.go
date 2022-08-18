package awslogin

import (
	"context"
	"github.com/pkg/errors"
	"github.com/splice/platform/localdev/v2/localdev/internal/command"
)

func NewAWSLogin() *AWSLogin {
	return &AWSLogin{}
}

type AWSLogin struct {
}

func (a *AWSLogin) ECRLogin(ctx context.Context) error {
	awsVault := `aws-vault exec splice -- aws ecr get-login-password --region us-west-1 | docker login --username AWS --password-stdin "118139069697.dkr.ecr.us-west-1.amazonaws.com"`

	testCommand := command.NewShellCommand("bash", "-c", awsVault)
	err := testCommand.Run(ctx)
	if err != nil {
		return errors.Wrap(err, "failed executing aws vault")
	}
	return nil
}
