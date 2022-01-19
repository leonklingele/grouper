package analyzer

import (
	"flag"
)

const (
	FlagNameImportRequireSingleImport = "import-require-single-import"
	FlagNameImportRequireGrouping     = "import-require-grouping"
)

func Flags() flag.FlagSet {
	fs := flag.NewFlagSet(Name, flag.ExitOnError)

	fs.Bool(FlagNameImportRequireSingleImport, false, "require the use of a single 'import' declaration only")
	fs.Bool(FlagNameImportRequireGrouping, false, "require the use of grouped 'import' declarations")

	return *fs
}
