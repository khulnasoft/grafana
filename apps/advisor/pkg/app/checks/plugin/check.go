package app

import (
	"context"
	"fmt"
	sysruntime "runtime"

	advisor "github.com/grafana/grafana/apps/advisor/pkg/apis/advisor/v0alpha1"
	"github.com/grafana/grafana/apps/advisor/pkg/app/common"
	"github.com/grafana/grafana/pkg/plugins/repo"
	"github.com/grafana/grafana/pkg/services/pluginsintegration/pluginstore"
)

func init() {
	common.RegisterCheck(&pluginCheckRegisterer{})
}

type pluginCheckRegisterer struct{}

func (p *pluginCheckRegisterer) New(cfg *common.AdvisorConfig) common.Check {
	return &PluginCheckImpl{
		pluginStore: cfg.PluginStore,
		pluginRepo:  cfg.PluginRepo,
	}
}

func (p *pluginCheckRegisterer) Type() string {
	return "plugin"
}

type PluginCheckImpl struct {
	pluginStore pluginstore.Store
	pluginRepo  repo.Service
}

func (c *PluginCheckImpl) Run(ctx context.Context, obj *advisor.CheckSpec) (*advisor.CheckV0alpha1StatusReport, error) {
	ps := c.pluginStore.Plugins(ctx)

	dsErrs := []advisor.CheckV0alpha1StatusReportErrors{}
	for _, p := range ps {
		// Check if plugin is deprecated
		i, err := c.pluginRepo.PluginInfo(ctx, p.ID)
		if err != nil {
			continue
		}
		if i.Status == "deprecated" {
			dsErrs = append(dsErrs, advisor.CheckV0alpha1StatusReportErrors{
				Type:   advisor.CheckStatusTypeInvestigation,
				Reason: fmt.Sprintf("Plugin deprecated: %s", p.ID),
				Action: "Look for alternatives",
			})
		}

		// Check if plugin has a newer version
		info, err := c.pluginRepo.GetPluginArchiveInfo(ctx, p.ID, "", repo.NewCompatOpts("", sysruntime.GOOS, sysruntime.GOARCH))
		if err != nil {
			continue
		}
		if info.Version != p.Info.Version { // TODO: Improve check for newer version
			dsErrs = append(dsErrs, advisor.CheckV0alpha1StatusReportErrors{
				Type:   advisor.CheckStatusTypeAction,
				Reason: fmt.Sprintf("Newer version available: %s", p.ID),
				Action: "Update plugin",
			})
		}
	}

	return &advisor.CheckV0alpha1StatusReport{
		Count:  int64(len(ps)),
		Errors: dsErrs,
	}, nil
}
