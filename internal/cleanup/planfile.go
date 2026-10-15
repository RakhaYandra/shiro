package cleanup

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

func SavePlan(path string, p CleanupPlan) error {
	raw, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(raw, '\n'), 0o644)
}

func LoadPlan(path string) (CleanupPlan, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return CleanupPlan{}, err
	}
	return decodePlan(raw)
}

func decodePlan(raw []byte) (CleanupPlan, error) {
	var p CleanupPlan
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&p); err != nil {
		return CleanupPlan{}, fmt.Errorf("invalid plan file: %w", err)
	}
	if p.SchemaVersion != 1 {
		return CleanupPlan{}, fmt.Errorf("unsupported plan schema_version %d", p.SchemaVersion)
	}
	if p.ID == "" || p.Root == "" {
		return CleanupPlan{}, fmt.Errorf("plan missing id or root")
	}
	return p, nil
}
