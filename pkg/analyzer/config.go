package analyzer

import (
	"github.com/leonklingele/grouper/pkg/analyzer/consts"
	"github.com/leonklingele/grouper/pkg/analyzer/imports"
)

type Config struct {
	ConstsConfig  *consts.Config
	ImportsConfig *imports.Config
}
