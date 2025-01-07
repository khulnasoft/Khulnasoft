package gomockgengazelle

import (
	"log"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/language"
	golang "github.com/bazelbuild/bazel-gazelle/language/go"
	"github.com/bazelbuild/bazel-gazelle/repo"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"
	"golang.org/x/exp/maps"

	mockgenc "github.com/khulnasoft/khulnasoft/dev/go-mockgen-transformer/config"
)

type gomockgen struct {
	language.BaseLang
	language.BaseLifecycleManager
}

var (
	_ (language.Language)            = (*gomockgen)(nil)
	_ (language.LifecycleManager)    = (*gomockgen)(nil)
	_ (language.ModuleAwareLanguage) = (*gomockgen)(nil)
)

var (
	logger = log.New(log.Writer(), log.Prefix(), log.Flags())

	yamlPayload mockgenc.YamlPayload

	// Keep track of all output directories mentioned in go-mockgen config, so that
	// we can notify if any where not visited & therefore no mocks were generated/updated
	// for them. Unlike with go-mockgen, we currently don't (or can't?) create new directories
	// as part of sg bazel configure, so we give the user a bash command to run to prepare the
	// missing directories (+ shell Go files) before re-running sg bazel configure.
	unvisitedDirs = make(map[string]bool)

	allOutputDirs = make(map[string]mockgenc.YamlMock)

	// Haven't discovered a way to get the workspace root dir yet (the cwd is deep in the Bazel output base somewhere),
	// so to read the go-mockgen files we set this value once we start traversing the workspace root.
	rootDir string

	// Currently, the go_mockgen rule extracts the required go-mockgen config at build time, so their labels
	// need to be passed as inputs.
	manifests []string
	// We want to load the mockgen.yaml (and related) config file(s) only once, so we use a sync.OnceValue.
	loadConfig = sync.OnceValue[error](func() (err error) {
		yamlPayload, err = mockgenc.ReadManifest(filepath.Join(rootDir, "mockgen.yaml"))
		if err != nil {
			return err
		}

		for _, mock := range yamlPayload.Mocks {
			allOutputDirs[filepath.Dir(mock.Filename)] = mock
			unvisitedDirs[filepath.Dir(mock.Filename)] = true
		}

		manifests = []string{"//:mockgen.yaml"}
		for _, manifest := range yamlPayload.IncludeConfigPaths {
			manifests = append(manifests, "//:"+manifest)
		}
		return nil
	})
)

func NewLanguage() language.Language {
	return &gomockgen{}
}

func (*gomockgen) Name() string { return "gomockgen" }

func (*gomockgen) Kinds() map[string]rule.KindInfo {
	return map[string]rule.KindInfo{
		"go_mockgen": {
			MatchAny: true,
			// I cant tell if these work or not...
			MergeableAttrs: map[string]bool{
				"deps":      true,
				"out":       true,
				"manifests": true,
			},
		},
		"write_source_files": {
			MatchAttrs: []string{"name"},
			MergeableAttrs: map[string]bool{
				"additional_update_targets": true,
			},
		},
	}
}

// From ModuleAwareLanguage.ApparentLoads:
//
// ApparentLoads returns .bzl files and symbols they define. Every rule
// generated by GenerateRules, now or in the past, should be loadable from
// one of these files.
//
// https://sourcegraph.com/github.com/bazelbuild/bazel-gazelle@ec1591cb193591b9544b52b98b8ce52833b34c58/-/blob/language/lang.go?L137:2-137:15
func (*gomockgen) ApparentLoads(moduleToApparentName func(string) string) []rule.LoadInfo {
	return []rule.LoadInfo{
		{
			Name:    "//dev:go_mockgen.bzl",
			Symbols: []string{"go_mockgen"},
		},
		{
			Name:    "@aspect_bazel_lib//lib:write_source_files.bzl",
			Symbols: []string{"write_source_files"},
		},
	}
}

func (g *gomockgen) GenerateRules(args language.GenerateArgs) language.GenerateResult {
	rootDir = args.Config.RepoRoot
	if err := loadConfig(); err != nil {
		log.Fatalf("failed to load go-mockgen config: %v", err)
	}

	// if we're in the ./dev folder, we want to generate an "all" target.
	if args.Rel == "dev" {
		var targets []string
		for _, mock := range yamlPayload.Mocks {
			targets = append(targets, "//"+filepath.Dir(mock.Filename)+":generate_mocks")
		}

		slices.Sort(targets)
		targets = slices.Compact(targets)

		catchallRule := rule.NewRule("write_source_files", "go_mockgen")
		catchallRule.SetAttr("additional_update_targets", targets)

		return language.GenerateResult{
			Gen:     []*rule.Rule{catchallRule},
			Imports: []interface{}{nil},
		}
	}

	// is this a directory we care about? if not, we don't want to generate any rules.
	mock, ok := allOutputDirs[args.Rel]
	if !ok {
		return language.GenerateResult{}
	}

	delete(unvisitedDirs, args.Rel)

	outputFilename := filepath.Base(mock.Filename)

	r := rule.NewRule("go_mockgen", "generate_mocks")
	r.SetAttr("out", outputFilename)
	r.SetAttr("manifests", manifests)

	// we want to add the generated file to either the go_library rule or the go_test rule, depending
	// on whether the file is a _test.go file or not.
	goRuleIndex := slices.IndexFunc(args.OtherGen, func(r *rule.Rule) bool {
		if strings.HasSuffix(outputFilename, "_test.go") {
			return r.Kind() == "go_test"
		} else {
			return r.Kind() == "go_library"
		}
	})
	if goRuleIndex == -1 {
		// We can revisit this output if it's something we hit in practice.
		log.Fatalf("couldn't find a go_{library,test} rule in \"%s/BUILD.bazel\"", args.Rel)
	}

	goRule := args.OtherGen[goRuleIndex]

	goRule.SetAttr("srcs", append(goRule.AttrStrings("srcs"), filepath.Base(mock.Filename)))

	imports := gatherDependencies(mock)

	return language.GenerateResult{
		Gen: []*rule.Rule{r, goRule},
		// Gen and Imports correspond per-index aka 'r' above is associated with 'imports' below.
		// This value gets passed as 'rawImports' in 'Resolve' below.
		Imports: []interface{}{imports, nil},
	}
}

func (g *gomockgen) DoneGeneratingRules() {
	if len(unvisitedDirs) > 0 {
		var b strings.Builder
		for _, dir := range maps.Keys(unvisitedDirs) {
			b.WriteString("mkdir -p ")
			b.WriteString(dir)
			b.WriteString(" && ")
			b.WriteString("echo 'package ")
			if allOutputDirs[dir].Package != "" {
				b.WriteString(allOutputDirs[dir].Package)
			} else {
				b.WriteString(filepath.Base(dir))
			}
			b.WriteString("' > ")
			b.WriteString(allOutputDirs[dir].Filename)
			b.WriteString(" \\\n && ")
		}
		b.WriteString("echo 'Done preparing! Re-running `sg bazel configure`' && sg bazel configure")
		logger.Fatalf("Some declared go-mockgen output files were not created due to their output directory missing. Please run the following to resolve this:\n%s", b.String())
	}
}

// Here we translate the Go import paths into Bazel labels for the go_mockgen rules.
func (g *gomockgen) Resolve(c *config.Config, ix *resolve.RuleIndex, rc *repo.RemoteCache, r *rule.Rule, rawImports interface{}, from label.Label) {
	if r.Kind() != "go_mockgen" {
		return
	}
	imports := rawImports.([]string)

	r.DelAttr("deps")

	labels := make([]string, 0, len(imports))
	for _, i := range imports {
		result, err := golang.ResolveGo(c, ix, rc, i, from)
		if err != nil {
			log.Fatalf("failed to resolve Go import path (%s) to Bazel label: %v", i, err)
		}
		labels = append(labels, result.Rel(from.Repo, from.Pkg).String())
	}
	r.SetAttr("deps", labels)
}

func gatherDependencies(mock mockgenc.YamlMock) (deps []string) {
	if mock.Path != "" {
		deps = append(deps, mock.Path)
	}
	for _, source := range mock.Sources {
		deps = append(deps, source.Path)
	}
	return
}
