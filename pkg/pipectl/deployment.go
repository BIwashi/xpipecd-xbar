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

//go:embed pipecd_base64.txt
var pipecdIconBase64 string

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
		var (
			l  = c.makeDeploymentLink(d)
			tr = true
		)

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
		)

		if stageStatus, ok := makeStageStatus(d.Stages); ok {
			xbars = append(xbars,
				xbar.Xbar{
					Line: xbar.Line{
						Title: stageStatus,
					},
				},
			)
		}

		xbars = append(xbars, xbar.SeparateLine)

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
		xbar.SeparateLine,
	}, xbars...)

	for _, x := range xbars {
		x.Print()
	}

	return nil
}

func (c *pipectl) makeDeploymentLink(deployment *model.Deployment) string {
	return fmt.Sprintf("https://%s/deployments/%s/", c.host, deployment.Id)
}

const (
	iconHeavyExclamationMark   = ":heavy_exclamation_mark:"  // ❗
	iconArrowsCounterclockwise = ":arrows_counterclockwise:" // 🔄
	iconArrowForward           = ":arrow_forward:"           // ▶️
	iconRewind                 = ":rewind:"                  // ⏪
	iconWhiteCheckMark         = ":white_check_mark:"        // ✅
	iconX                      = ":x:"                       // ❌
	iconNoEntry                = ":no_entry:"                // ⛔
	iconHourglassFlowingSand   = ":hourglass_flowing_sand:"  // ⏳
	iconCloud                  = ":cloud:"                   // ☁️
	iconMag                    = ":mag:"                     // 🔍
	iconBabyChick              = ":baby_chick:"              // 🐤
	iconHand                   = ":hand:"                    // ✋
)

func makeStatusIcon(deployment *model.Deployment) string {
	switch deployment.Status {
	case model.DeploymentStatus_DEPLOYMENT_PENDING:
		return iconHeavyExclamationMark
	case model.DeploymentStatus_DEPLOYMENT_RUNNING:
		return iconArrowsCounterclockwise
	case model.DeploymentStatus_DEPLOYMENT_ROLLING_BACK:
		return iconRewind
	case model.DeploymentStatus_DEPLOYMENT_SUCCESS:
		return iconWhiteCheckMark
	case model.DeploymentStatus_DEPLOYMENT_FAILURE:
		return iconX
	case model.DeploymentStatus_DEPLOYMENT_CANCELLED:
		return iconNoEntry
	default:
		return ""
	}
}

func makeStageStatus(stages []*model.PipelineStage) (string, bool) {
	var runningStage string
	for _, stage := range stages {
		if stage.Status == model.StageStatus_STAGE_RUNNING {
			runningStage = stage.Name
		}
	}

	if runningStage == "" {
		return "", false
	}

	return fmt.Sprintf("%s %s", makeStageStatusIcon(runningStage), runningStage), true
}

func makeStageStatusIcon(runningStage string) string {
	switch runningStage {
	case string(model.StageWait):
		return iconHourglassFlowingSand
	case string(model.StageWaitApproval):
		return iconHand
	case string(model.StageAnalysis):
		return iconMag
	case string(model.StageTerraformSync):
		return iconArrowsCounterclockwise
	case string(model.StageTerraformPlan):
		return iconArrowForward
	case string(model.StageTerraformApply):
		return iconArrowForward
	case string(model.StageECSSync):
		return iconArrowsCounterclockwise
	case string(model.StageECSCanaryRollout):
		return iconBabyChick
	case string(model.StageECSTrafficRouting):
		return iconBabyChick
	case string(model.StageECSCanaryClean):
		return iconBabyChick
	case string(model.StageCustomSync):
		return iconArrowsCounterclockwise
	case string(model.StageRollback):
		return iconRewind
	case string(model.StageCustomSyncRollback):
		return iconArrowsCounterclockwise
	default:
		return ""
	}
}
