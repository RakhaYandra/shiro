package cleanup

import (
	"os"

	"github.com/RakhaYandra/shiro/internal/policy"
)

type ItemResult struct {
	Path           string `json:"path"`
	Status         string `json:"status"`
	BytesReclaimed int64  `json:"bytes_reclaimed"`
	Reason         string `json:"reason"`
}

type Report struct {
	PlanID         string       `json:"plan_id"`
	Success        int          `json:"success"`
	Skipped        int          `json:"skipped"`
	Denied         int          `json:"denied"`
	Failed         int          `json:"failed"`
	BytesReclaimed int64        `json:"bytes_reclaimed"`
	Results        []ItemResult `json:"results"`
}

func Execute(pol policy.Policy, plan CleanupPlan) Report {
	rep := Report{PlanID: plan.ID}
	if rep.Results == nil {
		rep.Results = []ItemResult{}
	}
	for _, it := range plan.Items {
		v := Revalidate(pol, it, plan.Root)
		switch v.Status {
		case StatusStale:
			rep.Skipped++
			rep.Results = append(rep.Results, ItemResult{it.Path, "skipped", 0, v.Reason})
			continue
		case StatusDenied:
			rep.Denied++
			rep.Results = append(rep.Results, ItemResult{it.Path, "denied", 0, v.Reason})
			continue
		}
		if err := os.RemoveAll(it.Path); err != nil {
			rep.Failed++
			rep.Results = append(rep.Results, ItemResult{it.Path, "failed", 0, err.Error()})
			continue
		}
		if _, err := os.Lstat(it.Path); !os.IsNotExist(err) {
			rep.Failed++
			rep.Results = append(rep.Results, ItemResult{it.Path, "unknown", 0, "target still exists after removal"})
			continue
		}
		rep.Success++
		rep.BytesReclaimed += it.ExpectedSize
		rep.Results = append(rep.Results, ItemResult{it.Path, "success", it.ExpectedSize, "removed and verified gone"})
	}
	return rep
}
