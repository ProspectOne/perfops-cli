package cmd

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"testing"
)

func TestEnvPerfOpsAPIKey(t *testing.T) {
	os.Unsetenv("PERFOPS_API_KEY")

	rootCmd.ResetFlags()
	initRootCmd()

	if err := rootCmd.ParseFlags([]string{}); err != nil {
		t.Fatalf("exepected nil; got %v", err)
	}
	if got := apiKey; got != "" {
		t.Fatalf("expected no key; got %v", got)
	}

	const envAPIKey = "Meep"
	os.Setenv("PERFOPS_API_KEY", envAPIKey)
	defer os.Unsetenv("PERFOPS_API_KEY")

	rootCmd.ResetFlags()
	initRootCmd()

	if got := apiKey; got != envAPIKey {
		t.Fatalf("expected %v; got %v", envAPIKey, got)
	}

	const apiKey2 = "Moo"
	if err := rootCmd.ParseFlags([]string{"-K", apiKey2}); err != nil {
		t.Fatalf("exepected nil; got %v", err)
	}
	if got := apiKey; got != apiKey2 {
		t.Fatalf("expected %v; got %v", apiKey2, got)
	}
}

func TestVersionFlag(t *testing.T) {
	versionOutput := fmt.Sprintf(versionTmpl, version, buildDate, commitHash,
		runtime.Version(), runtime.Compiler, runtime.GOOS, runtime.GOARCH)

	var b bytes.Buffer
	rootCmd.ResetFlags()
	rootCmd.SetOutput(&b)
	initRootCmd()

	if err := rootCmd.ParseFlags([]string{"-v"}); err != nil {
		t.Fatalf("exepected nil; got %v", err)
	}
	if got := showVersion; !got {
		t.Fatalf("expected true; got %v", got)
	}
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("exepected nil; got %v", err)
	}
	if got := b.String(); got != versionOutput {
		t.Fatalf("expected %q; got %q", versionOutput, got)
	}
}
