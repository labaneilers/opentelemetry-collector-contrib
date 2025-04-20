// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package simplisafeidprocessor // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/simplisafeidprocessor"

import (
	"context"
	"regexp"
	"sort"
	"strings"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal/attraction"
	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/filter/expr"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottllog"
)

var id32Regexp = regexp.MustCompile(`\b[a-zA-Z0-9]{32}\b`)

type logsimplisafeidprocessor struct {
	logger   *zap.Logger
	attrProc *attraction.AttrProc
	skipExpr expr.BoolExpr[ottllog.TransformContext]
}

// newLogsimplisafeidprocessor returns a processor that modifies attributes of a
// log record. To construct the attributes processors, the use of the factory
// methods are required in order to validate the inputs.
func newLogsimplisafeidprocessor(logger *zap.Logger, attrProc *attraction.AttrProc, skipExpr expr.BoolExpr[ottllog.TransformContext]) *logsimplisafeidprocessor {
	return &logsimplisafeidprocessor{
		logger:   logger,
		attrProc: attrProc,
		skipExpr: skipExpr,
	}
}

func (a *logsimplisafeidprocessor) processLogs(ctx context.Context, ld plog.Logs) (plog.Logs, error) {
	rls := ld.ResourceLogs()
	for i := 0; i < rls.Len(); i++ {
		rs := rls.At(i)
		ilss := rs.ScopeLogs()

		for j := 0; j < ilss.Len(); j++ {
			ils := ilss.At(j)
			logs := ils.LogRecords()
			for k := 0; k < logs.Len(); k++ {
				lr := logs.At(k)

				topAttrs := lr.Attributes()

				var collectedIDs []string

				// recurse through all attributes
				var processAttributes func(attrs pcommon.Map)
				processAttributes = func(attrs pcommon.Map) {
					attrs.Range(func(key string, value pcommon.Value) bool {
						if value.Type() == pcommon.ValueTypeMap {
							// If the value is a nested map, recurse into it
							processAttributes(value.Map())
						} else if value.Type() == pcommon.ValueTypeStr {
							// If the value is a string, extract alphanumeric IDs
							strValue := value.Str()
							ids := id32Regexp.FindAllString(strValue, -1)
							if len(ids) > 0 {
								collectedIDs = append(collectedIDs, ids...)
							}
						}
						return true // Continue iteration
					})
				}

				processAttributes(topAttrs)

				if len(collectedIDs) > 0 {
					sort.Strings(collectedIDs) // Sort IDs alphanumerically
					topAttrs.PutStr("ss.ids", strings.Join(collectedIDs, ","))
					a.logger.Debug("Added ss.ids attribute", zap.Strings("ss.ids", collectedIDs))
				}

				a.attrProc.Process(ctx, a.logger, lr.Attributes())
			}
		}
	}
	return ld, nil
}
