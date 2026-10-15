package output

import (
	"strings"
	"testing"

	"github.com/RakhaYandra/shiro/internal/scanner"
)

func deepFixture() scanner.Summary {
	return scanner.Summary{
		Root:          "/x",
		TotalApparent: 300,
		Files:         3,
		Dirs:          3,
		Tree: &scanner.DirNode{
			Path:     "/x",
			Apparent: 300,
			Children: []*scanner.DirNode{
				{
					Path:     "/x/sub",
					Apparent: 200,
					Children: []*scanner.DirNode{{Path: "/x/sub/leaf", Apparent: 150}},
				},
			},
		},
	}
}

func TestRenderText_Depth(t *testing.T) {
	out1 := RenderText(deepFixture(), DisplayOpts{Depth: 1})
	if strings.Contains(out1, "/x/sub/leaf") {
		t.Fatalf("depth 1 must hide grandchildren:\n%s", out1)
	}
	out2 := RenderText(deepFixture(), DisplayOpts{Depth: 2})
	if !strings.Contains(out2, "/x/sub/leaf") {
		t.Fatalf("depth 2 must show grandchildren:\n%s", out2)
	}
}
