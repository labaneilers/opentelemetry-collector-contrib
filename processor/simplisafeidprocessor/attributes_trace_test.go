// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package simplisafeidprocessor

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/filter/filterset"
)

func createConfig(matchType filterset.MatchType) *filterset.Config {
	return &filterset.Config{
		MatchType: matchType,
	}
}
