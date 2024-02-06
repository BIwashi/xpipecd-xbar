package pipectl

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"os"
	"strings"

	api "github.com/pipe-cd/pipecd/pkg/app/server/service/apiservice"
	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/rpc/rpcauth"
	"github.com/pipe-cd/pipecd/pkg/rpc/rpcclient"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/credentials"

	"github.com/BIwashi/xpipecd-xbar/pkg/cli"
)

type pipectl struct {
	a      APIClient
	host   string
	apiKey string

	statuses []string
	appKinds []string
	appIds   []string
	appName  string
	labels   []string
	limit    int32

	cursor string
	stdout io.Writer
}

func NewCommand() *cobra.Command {
	c := pipectl{
		statuses: []string{
			// model.DeploymentStatus_DEPLOYMENT_PENDING.String(),
			model.DeploymentStatus_DEPLOYMENT_RUNNING.String(),
			// model.DeploymentStatus_DEPLOYMENT_SUCCESS.String(),
		},
		limit:  50,
		stdout: os.Stdout,
	}
	cmd := &cobra.Command{
		Use:   "pipectl",
		Short: "run pipectl command",
		RunE:  cli.WithContext(c.run),
	}

	cmd.Flags().StringVar(&c.host, "host", "pipecd.jp:433", "The address of the PipeCD Control Plane Server Address.")
	cmd.Flags().StringVar(&c.apiKey, "api-key", "", "API key for pipectl")
	cmd.Flags().StringSliceVar(&c.statuses, "status", c.statuses, fmt.Sprintf("The list of application statuses to filter. (%s)", strings.Join(model.DeploymentStatusStrings(), "|")))
	cmd.Flags().StringSliceVar(&c.appIds, "app-id", c.appIds, fmt.Sprintf("The list of application ids to filter. (%s)", strings.Join(model.ApplicationKindStrings(), "|")))
	cmd.Flags().StringSliceVar(&c.appKinds, "app-kind", c.appKinds, fmt.Sprintf("The list of application kinds to filter. (%s)", strings.Join(model.ApplicationKindStrings(), "|")))
	cmd.Flags().StringVar(&c.appName, "app-name", c.appName, "The application name to filter.")
	cmd.Flags().StringVar(&c.cursor, "cursor", c.cursor, "The cursor which returned by the previous request applications list.")
	cmd.Flags().Int32Var(&c.limit, "limit", c.limit, "Upper limit on the number of return values. Default value is 30.")
	cmd.Flags().StringSliceVar(&c.labels, "label", c.labels, "The list of labels to filter. Expect input in the form KEY:VALUE.")

	return cmd
}

func (c *pipectl) run(ctx context.Context, input cli.Input) error {
	// logger := input.Logger
	// logger.Info("hello world")

	creds := rpcclient.NewPerRPCCredentials(c.apiKey, rpcauth.APIKeyCredentials, true)
	tlsConfig := &tls.Config{}
	options := []rpcclient.DialOption{
		// rpcclient.WithBlock(),
		rpcclient.WithPerRPCCredentials(creds),
		rpcclient.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
	}

	client, err := api.NewClient(ctx, c.host, options...)
	if err != nil {
		return errors.Wrap(err, "failed to create api client")
	}
	c.a = client
	defer c.a.Close()

	if err := c.listDeployments(ctx); err != nil {
		return errors.Wrap(err, "failed to list deployments")
	}

	return nil
}
