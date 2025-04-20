// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package simplisafeidprocessor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/processor/processortest"

	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal/attraction"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/simplisafeidprocessor/internal/metadata"
)

func TestFactory_Type(t *testing.T) {
	factory := NewFactory()
	assert.Equal(t, factory.Type(), metadata.Type)
}

func TestFactory_CreateDefaultConfig(t *testing.T) {
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()
	assert.Equal(t, &Config{}, cfg)
	assert.NoError(t, componenttest.CheckConfigStruct(cfg))
}

// func TestValidateConfig(t *testing.T) {
// 	factory := NewFactory()
// 	cfg := factory.CreateDefaultConfig()
// 	assert.Error(t, xconfmap.Validate(cfg))
// }

func TestFactoryCreateTraces_NotImplemented(t *testing.T) {
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()

	ap, err := factory.CreateTraces(context.Background(), processortest.NewNopSettings(metadata.Type), cfg, consumertest.NewNop())
	assert.Error(t, err)
	assert.Nil(t, ap)

}

func TestFactoryCreateMetrics_NotImplemented(t *testing.T) {
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()

	mp, err := factory.CreateMetrics(context.Background(), processortest.NewNopSettings(metadata.Type), cfg, consumertest.NewNop())
	require.Error(t, err)
	assert.Nil(t, mp)

}

func TestFactoryCreateLogs(t *testing.T) {
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()
	oCfg := cfg.(*Config)
	oCfg.Actions = []attraction.ActionKeyValue{
		{Key: "a key", Action: attraction.DELETE},
	}

	tp, err := factory.CreateLogs(
		context.Background(), processortest.NewNopSettings(metadata.Type), cfg, consumertest.NewNop())
	assert.NotNil(t, tp)
	assert.NoError(t, err)

	oCfg.Actions = []attraction.ActionKeyValue{
		{Action: attraction.DELETE},
	}
	tp, err = factory.CreateLogs(
		context.Background(), processortest.NewNopSettings(metadata.Type), cfg, consumertest.NewNop())
	assert.Nil(t, tp)
	assert.Error(t, err)
}
