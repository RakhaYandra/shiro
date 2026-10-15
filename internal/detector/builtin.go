package detector

// DefaultRegistry returns the built-in M3 rule set.
func DefaultRegistry() *Registry {
	reg := NewRegistry()
	for _, r := range []Rule{
		{
			ID: "DET-003", Category: "node", Basename: "node_modules",
			Markers:  []string{"package.json", "package-lock.json", "pnpm-lock.yaml", "yarn.lock"},
			Describe: "Node.js dependencies (regenerable via package manager)",
		},
		{
			ID: "DET-002a", Category: "go", Basename: "go-build",
			Describe: "Go build cache (regenerable via go build)",
		},
		{
			ID: "DET-002b", Category: "go", Basename: "mod", Parent: "pkg",
			Describe: "Go module cache (regenerable via go mod download)",
		},
		{
			ID: "DET-004a", Category: "python", Basename: "__pycache__",
			Markers:  []string{"__init__.py", "setup.py", "pyproject.toml"},
			Describe: "Python bytecode cache (regenerable on import)",
		},
		{
			ID: "DET-004b", Category: "python", Basename: ".venv",
			Markers:  []string{"pyvenv.cfg"},
			Describe: "Python virtualenv (regenerable via pip install)",
		},
		{
			ID: "DET-004c", Category: "python", Basename: "venv",
			Markers:  []string{"pyvenv.cfg"},
			Describe: "Python virtualenv (regenerable via pip install)",
		},
		{
			ID: "DET-GEN", Category: "build", Basename: "target",
			Markers:  []string{"Cargo.toml"},
			Describe: "Rust build output (regenerable via cargo build)",
		},
		{
			ID: "DET-GEN", Category: "build", Basename: "build",
			Markers:  []string{"build.gradle", "build.gradle.kts", "pom.xml"},
			Describe: "JVM build output (regenerable via build tool)",
		},
		{
			ID: "DET-GEN", Category: "build", Basename: "dist",
			Markers:  []string{"package.json"},
			Describe: "Frontend build output (regenerable via build tool)",
		},
		{
			ID: "DET-GEN", Category: "build", Basename: ".next",
			Markers:  []string{"package.json"},
			Describe: "Next.js build output (regenerable via next build)",
		},
	} {
		reg.Register(r)
	}
	return reg
}
