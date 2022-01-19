package analyzer

import (
	"github.com/leonklingele/grouper/pkg/analyzer/consts"
	"github.com/leonklingele/grouper/pkg/analyzer/imports"
	"github.com/leonklingele/grouper/pkg/analyzer/types"
)

type Config struct {
	ConstsConfig  *consts.Config
	ImportsConfig *imports.Config
	TypesConfig   *types.Config
}
