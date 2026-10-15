package tui

import (
	"github.com/RakhaYandra/shiro/internal/cleanup"
	"github.com/RakhaYandra/shiro/internal/output"
)

type planView struct {
	plan cleanup.CleanupPlan
}

func newPlanView(plan cleanup.CleanupPlan) planView { return planView{plan} }

func (p planView) View() string {
	return output.RenderPlanText(p.plan)
}
