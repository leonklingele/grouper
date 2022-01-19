package analyzer

import (
	"flag"
)

const (
	FlagNameConstRequireSingleConst = "const-require-single-const"
	FlagNameConstRequireGrouping    = "const-require-grouping"

	FlagNameImportRequireSingleImport = "import-require-single-import"
	FlagNameImportRequireGrouping     = "import-require-grouping"
)

func Flags() flag.FlagSet {
	fs := flag.NewFlagSet(Name, flag.ExitOnError)

	fs.Bool(FlagNameConstRequireSingleConst, false, "require the use of a single global 'const' declaration only")
	fs.Bool(FlagNameConstRequireGrouping, false, "require the use of grouped global 'const' declarations")

	fs.Bool(FlagNameImportRequireSingleImport, false, "require the use of a single 'import' declaration only")
	fs.Bool(FlagNameImportRequireGrouping, false, "require the use of grouped 'import' declarations")

	return *fs
}
