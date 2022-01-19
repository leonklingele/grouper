package analyzer_test

import (
	"flag"
	"path/filepath"
	"testing"

	"github.com/leonklingele/grouper/pkg/analyzer"

	"golang.org/x/tools/go/analysis/analysistest"
)

// TODO(leon): Add fuzzing

func TestConst(t *testing.T) {
	t.Parallel()

	fixtures := []struct {
		name  string
		flags flag.FlagSet
	}{
		{
			name: "single-grouped",
			flags: flags().
				withConstRequireGrouping().
				build(),
		},
		{
			name: "single-ungrouped",
			flags: flags().
				withConstRequireGrouping().
				build(),
		},

		{
			name: "multi-grouped",
			flags: flags().
				withConstRequireSingleConst().
				withConstRequireGrouping().
				build(),
		},
		{
			name: "multi-ungrouped",
			flags: flags().
				withConstRequireSingleConst().
				withConstRequireGrouping().
				build(),
		},

		{
			name: "mixed-require-single-const",
			flags: flags().
				withConstRequireSingleConst().
				build(),
		},
		{
			name: "mixed-require-grouping",
			flags: flags().
				withConstRequireGrouping().
				build(),
		},

		{
			name: "mixed-named-with-vars",
			flags: flags().
				withConstRequireSingleConst().
				withConstRequireGrouping().
				build(),
		},
	}

	for _, f := range fixtures {
		f := f

		t.Run(f.name, func(t *testing.T) {
			t.Parallel()

			a := analyzer.New()
			a.Flags = f.flags

			testdata := filepath.Join(analysistest.TestData(), "const")
			_ = analysistest.Run(t, testdata, a, f.name)
		})
	}
}

func TestImport(t *testing.T) {
	t.Parallel()

	fixtures := []struct {
		name  string
		flags flag.FlagSet
	}{
		{
			name: "single-grouped",
			flags: flags().
				withImportRequireGrouping().
				build(),
		},
		{
			name: "single-ungrouped",
			flags: flags().
				withImportRequireGrouping().
				build(),
		},

		{
			name: "multi-grouped",
			flags: flags().
				withImportRequireSingleImport().
				withImportRequireGrouping().
				build(),
		},
		{
			name: "multi-ungrouped",
			flags: flags().
				withImportRequireSingleImport().
				withImportRequireGrouping().
				build(),
		},

		{
			name: "mixed-require-single-import",
			flags: flags().
				withImportRequireSingleImport().
				build(),
		},
		{
			name: "mixed-require-grouping",
			flags: flags().
				withImportRequireGrouping().
				build(),
		},

		{
			name: "mixed-named-require-single-import",
			flags: flags().
				withImportRequireSingleImport().
				build(),
		},
		{
			name: "mixed-named-require-grouping",
			flags: flags().
				withImportRequireGrouping().
				build(),
		},
	}

	for _, f := range fixtures {
		f := f

		t.Run(f.name, func(t *testing.T) {
			t.Parallel()

			a := analyzer.New()
			a.Flags = f.flags

			testdata := filepath.Join(analysistest.TestData(), "import")
			_ = analysistest.Run(t, testdata, a, f.name)
		})
	}
}

type flagger struct {
	fs *flag.FlagSet
}

func (f *flagger) withConstRequireSingleConst() *flagger {
	if err := f.fs.Lookup(analyzer.FlagNameConstRequireSingleConst).Value.Set("true"); err != nil {
		panic(err)
	}

	return f
}

func (f *flagger) withConstRequireGrouping() *flagger {
	if err := f.fs.Lookup(analyzer.FlagNameConstRequireGrouping).Value.Set("true"); err != nil {
		panic(err)
	}

	return f
}

func (f *flagger) withImportRequireSingleImport() *flagger {
	if err := f.fs.Lookup(analyzer.FlagNameImportRequireSingleImport).Value.Set("true"); err != nil {
		panic(err)
	}

	return f
}

func (f *flagger) withImportRequireGrouping() *flagger {
	if err := f.fs.Lookup(analyzer.FlagNameImportRequireGrouping).Value.Set("true"); err != nil {
		panic(err)
	}

	return f
}

func (f *flagger) build() flag.FlagSet {
	return *f.fs
}

func flags() *flagger {
	fs := analyzer.Flags()

	return &flagger{
		fs: &fs,
	}
}
