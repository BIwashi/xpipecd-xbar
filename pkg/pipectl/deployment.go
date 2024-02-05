package pipectl

import (
	"context"
	_ "embed"
	"fmt"
	"strings"

	api "github.com/pipe-cd/pipecd/pkg/app/server/service/apiservice"
	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pkg/errors"

	"github.com/BIwashi/xpipecd-xbar/pkg/xbar"
)

func (c *pipectl) listDeployments(ctx context.Context) error {
	labels := map[string]string{}
	for _, label := range c.labels {
		sp := strings.SplitN(label, ":", 2)
		if len(sp) == 2 {
			labels[sp[0]] = sp[1]
		}
	}

	listReq := &api.ListDeploymentsRequest{
		Statuses:        c.statuses,
		Kinds:           c.appKinds,
		ApplicationIds:  c.appIds,
		ApplicationName: c.appName,
		Limit:           c.limit,
		Cursor:          c.cursor,
		Labels:          labels,
	}

	resp, err := c.a.ListDeployments(ctx, listReq)
	if err != nil {
		return errors.Wrap(err, "failed to list deployments")
	}

	var (
		xbars            = make([]xbar.Xbar, 0, len(resp.Deployments))
		pending, running int
	)

	for _, d := range resp.Deployments {
		l := c.makeDeploymentLink(d)

		xbars = append(xbars,
			xbar.Xbar{
				Line: xbar.Line{
					Title:    fmt.Sprintf("%s %s", makeStatusIcon(d), d.ApplicationName),
					Href:     &l,
					Dropdown: &tr,
				},
				SubLine: []xbar.Xbar{
					{
						Line: xbar.Line{
							Title: fmt.Sprintf("Status: %s", d.Status),
						},
					},
				},
			},
			xbar.Xbar{
				Line: xbar.Line{
					Title: makeStageStatus(d.Stages),
				},
			},
			xbar.SeparateLine,
		)

		switch d.Status {
		case model.DeploymentStatus_DEPLOYMENT_PENDING:
			pending++
		case model.DeploymentStatus_DEPLOYMENT_RUNNING:
			running++
		}
	}

	xbars = append([]xbar.Xbar{
		{
			Line: xbar.Line{
				Title:         fmt.Sprintf("(%d/%d)", pending, running),
				TemplateImage: &pipecdIconBase64,
			},
		},
	}, xbars...)

	for _, x := range xbars {
		x.Print()
	}

	return nil
}

func (c *pipectl) makeDeploymentLink(deployment *model.Deployment) string {
	return fmt.Sprintf("https://%s/deployments/%s/", c.host, deployment.Id)
}

func makeStatusIcon(deployment *model.Deployment) string {
	switch deployment.Status {
	case model.DeploymentStatus_DEPLOYMENT_PENDING:
		return ":heavy_exclamation_mark:"
	case model.DeploymentStatus_DEPLOYMENT_RUNNING:
		return ":arrows_counterclockwise:"
	case model.DeploymentStatus_DEPLOYMENT_ROLLING_BACK:
		return ":rewind:"
	case model.DeploymentStatus_DEPLOYMENT_SUCCESS:
		return ":white_check_mark:"
	case model.DeploymentStatus_DEPLOYMENT_FAILURE:
		return ":x:"
	case model.DeploymentStatus_DEPLOYMENT_CANCELLED:
		return ":no_entry:"
	default:
		return ""
	}
}

func makeStageStatus(stages []*model.PipelineStage) string {
	var runningStage string
	for _, stage := range stages {
		if stage.Status == model.StageStatus_STAGE_RUNNING {
			runningStage = stage.Name
		}
	}

	return fmt.Sprintf("%s %s", makeStageStatusIcon(runningStage), runningStage)
}

// TODO: Add icon for each stage.
func makeStageStatusIcon(runningStage string) string {
	switch runningStage {
	case string(model.StageWait):
		return ""
	case string(model.StageWaitApproval):
		return ":hand:"
	case string(model.StageAnalysis):
		return ""
	case string(model.StageTerraformSync):
		return ":cloud:"
	case string(model.StageTerraformPlan):
		return ":cloud:"
	case string(model.StageTerraformApply):
		return ":cloud:"
	case string(model.StageECSSync):
		return ":cloud:"
	case string(model.StageECSCanaryRollout):
		return ":cloud:"
	case string(model.StageECSTrafficRouting):
		return ":cloud:"
	case string(model.StageECSCanaryClean):
		return ":cloud:"
	case string(model.StageCustomSync):
		return ":cloud:"
	case string(model.StageRollback):
		return ":rewind:"
	case string(model.StageCustomSyncRollback):
		return ":rewind:"
	default:
		return ""
	}
}

//go:embed pipecd_base64.txt
var pipecdIconBase64 string

var (
	tr = true  //nolint:unused
	fl = false //nolint:unused
)
