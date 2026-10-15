package cleanup

import "testing"

func FuzzDecodePlan(f *testing.F) {
	for _, seed := range []string{
		`{}`,
		`{"schema_version":1,"id":"a","root":"/x"}`,
		`{"schema_version":99,"id":"a","root":"/x"}`,
		`{not json`,
		"",
		`{"schema_version":1,"id":"a","root":"/x","items":[{"path":"/x","action":"remove"}]}`,
	} {
		f.Add([]byte(seed))
	}
	f.Fuzz(func(t *testing.T, raw []byte) {
		p, err := decodePlan(raw)
		if err == nil && (p.SchemaVersion != 1 || p.ID == "" || p.Root == "") {
			t.Fatalf("accepted invalid plan: %+v", p)
		}
	})
}
