package policy

import "github.com/RakhaYandra/shiro/internal/detector"

type Decision struct {
	Allowed      bool
	Reason       string
	NeedsConfirm bool
}

func Decide(p Policy, f detector.Finding, root string) Decision {
	if err := p.CheckIn(f.Path, root); err != nil {
		return Decision{Reason: err.Error()}
	}
	return Decision{
		Allowed:      true,
		Reason:       "artifact match " + f.RuleID + " (" + f.Confidence + " confidence); requires explicit confirmation",
		NeedsConfirm: true,
	}
}
