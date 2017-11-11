// Package init bootstraps an Apex project.
package init

import (
	"errors"

	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/tj/cobra"

	"github.com/apex/apex/boot"
	"github.com/apex/apex/cmd/apex/root"
)

var credentialsError = `

  AWS region missing, are your credentials set up? Try:

  $ export AWS_PROFILE=myapp-stage
  $ apex init

  Visit http://apex.run/#aws-credentials for more details on
  setting up AWS credentials and specifying which profile to
  use.

`

var name string
var description string
var noPrompt bool
var skipSkeleton bool

// Command config.
var Command = &cobra.Command{
	Use:              "init",
	Short:            "Initialize a project",
	PersistentPreRun: root.PreRunNoop,
	RunE:             run,
}

// Initialize.
func init() {
	root.Register(Command)

	f := Command.PersistentFlags()

	f.StringVarP(&name, "name", "n", "", "Project name")
	f.StringVarP(&description, "description", "d", "", "Project description")
	f.BoolVarP(&noPrompt, "no-prompt", "q", false, "Non interactive mode")
	f.BoolVarP(&skipSkeleton, "skip-skeleton", "s", false, "Dont't create hello function")
}

// Run command.
func run(c *cobra.Command, args []string) error {
	if err := root.Prepare(c, args); err != nil {
		return err
	}

	region := root.Config.Region
	if region == nil {
		return errors.New(credentialsError)
	}

	b := boot.Bootstrapper{
		IAM:    iam.New(root.Session),
		Region: *region,
		Name: name,
		Description: description,
		SkipSkeleton: skipSkeleton,
	}

	return b.Boot(noPrompt)
}
