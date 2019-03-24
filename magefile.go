// +build mage

package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	outDir       = "build"
	basePackage  = "github.com/kasvith/kache"
	noGitLdflags = "-X $PACKAGE/internal/cobra-cmds.BuildDate=$BUILD_DATE"
)

var ldflags = `-X $PACKAGE/internal/cobra-cmds.CommitHash=$COMMIT_HASH 
				-X $PACKAGE/internal/cobra-cmds.BuildDate=$BUILD_DATE`

// allow user to override go executable by running as GOEXE=xxx make ... on unix-like systems
var goexe = "go"

func init() {
	if exe := os.Getenv("GOEXE"); exe != "" {
		goexe = exe
	}

	os.Setenv("GO111MODULE", "on")
}

func flagEnv(osType, osArch string) map[string]string {
	hash, _ := sh.Output("git", "rev-parse", "--short", "HEAD")
	return map[string]string{
		"PACKAGE":     basePackage,
		"OUT_DIR":     outDir,
		"COMMIT_HASH": hash,
		"BUILD_DATE":  time.Now().Format("2006-01-02T15:04:05Z0700"),
		"GOOS":        osType,
		"GOARCH":      osArch,
	}
}

func addOSExecType(osType, osArch, str string) string {
	str = fmt.Sprintf("%s-%s-%s", str, osType, osArch)

	if osType == "windows" {
		return str + ".exe"
	}

	return str
}

// get go imports
func getGoImports() error {
	return sh.Run(goexe, "get", "-u", "golang.org/x/tools/cmd/goimports")
}

// get go lint
func getGoLint() error {
	return sh.Run(goexe, "get", "-u", "golang.org/x/lint/golint")
}

// Install dependencies to vendor
func Vendor() error {
	mg.Deps(getGoLint)
	mg.Deps(getGoImports)

	return sh.Run(goexe, "mod", "download")
}

// Build kache
func Kache() error {
	return sh.RunWith(flagEnv(runtime.GOOS, runtime.GOARCH), goexe, "build", "-ldflags", ldflags, "-o", addOSExecType(runtime.GOOS, runtime.GOARCH, "$OUT_DIR/kache"), "$PACKAGE/cmd/kache")
}

// Build kache without git info
func KacheNoGitInfo() error {
	ldflags = noGitLdflags
	return Kache()
}

// Build kache cross platform
func KacheCrossBuild() error {
	// os/arch types
	platforms := map[string][]string{
		"linux":   []string{"386", "amd64", "arm"},
		"darwin":  []string{"386", "amd64"},
		"windows": []string{"386", "amd64"},
		"openbsd": []string{"386", "amd64"},
	}

	for osType, archs := range platforms {
		for _, arch := range archs {
			fmt.Printf("Building for %s-%s\n", osType, arch)
			err := sh.RunWith(flagEnv(osType, arch), goexe, "build", "-ldflags", ldflags, "-o", addOSExecType(osType, arch, "$OUT_DIR/kache"), "$PACKAGE/cmd/kache")
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Build kache-cli
func KacheCli() error {
	return sh.RunWith(flagEnv(runtime.GOOS, runtime.GOARCH), goexe, "build", "-ldflags", ldflags, "-o", addOSExecType(runtime.GOOS, runtime.GOARCH, "$OUT_DIR/kache-cli"), "$PACKAGE/cmd/kache-cli")
}

// Build kache-cli without git info
func KacheCliNoGitInfo() error {
	ldflags = noGitLdflags
	return Kache()
}

// Run gofmt, vet, imports and tests also with race
func Check() {
	fmt.Println("Checking started")

	if strings.Contains(runtime.Version(), "1.8") {
		// Go 1.8 doesn't play along with go test ./... and /vendor.
		// We could fix that, but that would take time.
		fmt.Printf("Skip Check on %s\n", runtime.Version())
		return
	}

	// Do this because CI memory error can occur
	mg.Deps(Lint)
	mg.Deps(Fmt)
	mg.Deps(Vet)
	mg.Deps(Imports)
	mg.Deps(Test)
}

// Run tests
func Test() error {
	return sh.RunV(goexe, "test", "-v", "./...")
}

// Run tests with race detector
func TestRace() error {
	return sh.RunV(goexe, "test", "-race", "./...")
}

// get all packages of kache
func kachePackages() ([]string, error) {
	var pkgPrefixLen = len(basePackage)
	s, err := sh.Output(goexe, "list", "./...")
	if err != nil {
		return nil, err
	}
	pkgs := strings.Split(s, "\n")
	for i := range pkgs {
		pkgs[i] = "." + pkgs[i][pkgPrefixLen:]
	}
	return pkgs, nil
}

// check if go is latest
func isGoLatest() bool {
	return strings.Contains(runtime.Version(), "1.11")
}

// Run gofmt linter
func Fmt() error {
	if !isGoLatest() {
		return nil
	}
	pkgs, err := kachePackages()
	if err != nil {
		return err
	}
	failed := false
	first := true
	for _, pkg := range pkgs {
		files, err := filepath.Glob(filepath.Join(pkg, "*.go"))
		if err != nil {
			return nil
		}
		for _, f := range files {
			// gofmt doesn't exit with non-zero when it finds unformatted code
			// so we have to explicitly look for output, and if we find any, we
			// should fail this target.
			s, err := sh.Output("gofmt", "-l", f)
			if err != nil {
				fmt.Printf("ERROR: running gofmt on %q: %v\n", f, err)
				failed = true
			}
			if s != "" {
				if first {
					fmt.Println("The following files are not gofmt'ed:")
					first = false
				}
				failed = true
				fmt.Println(s)
			}
		}
	}
	if failed {
		return errors.New("improperly formatted go files")
	}
	return nil
}

// Run goimports
func Imports() error {
	pkgs, err := kachePackages()
	if err != nil {
		return err
	}
	failed := false
	first := true
	for _, pkg := range pkgs {
		files, err := filepath.Glob(filepath.Join(pkg, "*.go"))
		if err != nil {
			return nil
		}
		for _, f := range files {
			// gofmt doesn't exit with non-zero when it finds unformatted code
			// so we have to explicitly look for output, and if we find any, we
			// should fail this target.
			s, err := sh.Output("goimports", "-w", `-local='github.com/kasvith/kache'`, "-l", f)
			if err != nil {
				fmt.Printf("ERROR: running goimports on %q: %v\n", f, err)
				failed = true
			}
			if s != "" {
				if first {
					fmt.Println("The following files are not goimport'ed:")
					first = false
				}
				failed = true
				fmt.Println(s)
			}
		}
	}
	if failed {
		return errors.New("run goimports in following files")
	}
	return nil
}

// Run golint
func Lint() error {
	pkgs, err := kachePackages()
	if err != nil {
		return err
	}
	failed := false
	for _, pkg := range pkgs {
		// We don't actually want to fail this target if we find golint errors,
		// so we don't pass -set_exit_status, but we still print out any failures.
		if _, err := sh.Exec(nil, os.Stderr, nil, "golint", "-set_exit_status", pkg); err != nil {
			fmt.Printf("ERROR: running go lint on %q: %v\n", pkg, err)
			failed = true
		}
	}
	if failed {
		return errors.New("errors running golint")
	}
	return nil
}

//  Run go vet linter
func Vet() error {
	if err := sh.Run(goexe, "vet", "./..."); err != nil {
		return fmt.Errorf("error running go vet: %v", err)
	}
	return nil
}

// Generate test coverage report
func TestCover() error {
	return sh.RunV(goexe, "test", "-race", "-coverprofile=coverage.txt", "-covermode=atomic", "./...")
}

// Verify that vendored packages match git HEAD
func CheckVendor() error {
	if err := sh.RunV("git", "diff-index", "--quiet", "HEAD", "vendor/"); err != nil {
		// yes, ignore errors from this, not much we can do.
		sh.Exec(nil, os.Stdout, os.Stderr, "git", "diff", "vendor/")
		return errors.New("check-vendor target failed: vendored packages out of sync")
	}
	return nil
}
