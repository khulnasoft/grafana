package v1alpha1

import (
	"github.com/grafana/grafana/pkg/apis/dashboard"
	conversion "k8s.io/apimachinery/pkg/conversion"
)

func Convert_v1alpha1_Dashboard_To_dashboard_Dashboard(in *Dashboard, out *dashboard.Dashboard, s conversion.Scope) error {
	panic("TODO: implement v1alpha1 -> internal conversion")
}
