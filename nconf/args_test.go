package nconf

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestArgsLoad(t *testing.T) {
	cfg := &struct {
		Something  string
		Other      int
		Overridden string
	}{
		Something:  "default",
		Overridden: "this should change",
	}

	tmp, err := ioutil.TempFile("", "something*.env")
	require.NoError(t, err)
	cfgStr := `
PF_OTHER=10
PF_OVERRIDDEN=not-that
PF_LOG_LEVEL=debug
PF_LOG_QUOTE_EMPTY_FIELDS=true
`
	require.NoError(t, ioutil.WriteFile(tmp.Name(), []byte(cfgStr), 0644))

	args := RootArgs{
		Prefix:     "pf",
		ConfigFile: tmp.Name(),
	}

	log, err := args.Setup(cfg, "", "")
	require.NoError(t, err)

	// check that we did call configure the logger
	assert.NotNil(t, log)
	entry := log.(*logrus.Entry)
	assert.Equal(t, logrus.DebugLevel, entry.Logger.Level)
	assert.True(t, entry.Logger.Formatter.(*logrus.TextFormatter).QuoteEmptyFields)

	assert.Equal(t, "default", cfg.Something)
	assert.Equal(t, 10, cfg.Other)
	assert.Equal(t, "not-that", cfg.Overridden)
}

func TestArgsAddToCmd(t *testing.T) {
	args := new(RootArgs)
	var called int
	cmd := cobra.Command{
		Run: func(_ *cobra.Command, _ []string) {
			assert.Equal(t, "PF", args.Prefix)
			assert.Equal(t, "file.env", args.ConfigFile)
			called++
		},
	}
	cmd.PersistentFlags().AddFlag(args.ConfigFlag())
	cmd.PersistentFlags().AddFlag(args.PrefixFlag())
	cmd.SetArgs([]string{"--config", "file.env", "--prefix", "PF"})
	require.NoError(t, cmd.Execute())
	assert.Equal(t, 1, called)
}

func TestArgsLoadDefault(t *testing.T) {
	configVals := map[string]interface{}{
		"log": map[string]interface{}{
			"level": "debug",
			"fields": map[string]interface{}{
				"something": 1,
			},
		},
		"bugsnag": map[string]interface{}{
			"api_key":         "secrets",
			"project_package": "package",
		},
		"metrics": map[string]interface{}{
			"enabled": true,
			"port":    8125,
			"tags": map[string]string{
				"env": "prod",
			},
		},
		"tracing": map[string]interface{}{
			"enabled":      true,
			"port":         "8125",
			"enable_debug": true,
		},
		"featureflag": map[string]interface{}{
			"key":             "magicalkey",
			"request_timeout": "10s",
			"enabled":         true,
		},
	}

	scenes := []struct {
		ext string
		enc func(v interface{}) ([]byte, error)
	}{
		{"json", json.Marshal},
		{"yaml", yaml.Marshal},
	}
	for _, s := range scenes {
		t.Run(s.ext, func(t *testing.T) {
			f, err := ioutil.TempFile("", "test-config-*."+s.ext)
			require.NoError(t, err)
			defer os.Remove(f.Name())

			b, err := s.enc(&configVals)
			require.NoError(t, err)
			_, err = f.Write(b)
			require.NoError(t, err)

			args := RootArgs{
				ConfigFile: f.Name(),
			}
			cfg, err := args.loadDefaultConfig()
			require.NoError(t, err)

			// logging
			assert.Equal(t, "debug", cfg.Log.Level)
			assert.Equal(t, true, cfg.Log.QuoteEmptyFields)
			assert.Equal(t, "", cfg.Log.File)
			assert.Equal(t, false, cfg.Log.DisableColors)
			assert.Equal(t, "", cfg.Log.TSFormat)

			assert.Len(t, cfg.Log.Fields, 1)
			assert.EqualValues(t, 1, cfg.Log.Fields["something"])
			assert.Equal(t, false, cfg.Log.UseNewLogger)

			// bugsnag
			assert.Equal(t, "", cfg.BugSnag.Environment)
			assert.Equal(t, "secrets", cfg.BugSnag.APIKey)
			assert.Equal(t, false, cfg.BugSnag.LogHook)
			assert.Equal(t, "package", cfg.BugSnag.ProjectPackage)

			// metrics
			assert.Equal(t, true, cfg.Metrics.Enabled)
			assert.Equal(t, "", cfg.Metrics.Host)
			assert.Equal(t, 8125, cfg.Metrics.Port)
			assert.Equal(t, map[string]string{"env": "prod"}, cfg.Metrics.Tags)

			// tracing
			assert.Equal(t, true, cfg.Tracing.Enabled)
			assert.Equal(t, "", cfg.Tracing.Host)
			assert.Equal(t, "8125", cfg.Tracing.Port)
			assert.Empty(t, cfg.Tracing.Tags)
			assert.Equal(t, true, cfg.Tracing.EnableDebug)

			// featureflag
			assert.Equal(t, "magicalkey", cfg.FeatureFlag.Key)
			assert.Equal(t, 10*time.Second, cfg.FeatureFlag.RequestTimeout.Duration)
			assert.Equal(t, true, cfg.FeatureFlag.Enabled)
			assert.Equal(t, false, cfg.FeatureFlag.DisableEvents)
			assert.Equal(t, "", cfg.FeatureFlag.RelayHost)
		})
	}
}

func TestArgsLoadFromYAML(t *testing.T) {
	f, err := ioutil.TempFile("", "test-config-*.yaml")
	require.NoError(t, err)
	defer os.Remove(f.Name())

	args := RootArgs{
		ConfigFile: f.Name(),
	}

	t.Run("empty-file", func(t *testing.T) {
		cfg := RootConfig{
			Log: DefaultLoggingConfig(),
		}
		require.NoError(t, args.load(&cfg))

		assert.True(t, cfg.Log.QuoteEmptyFields)
		assert.Equal(t, DefaultLoggingConfig(), cfg.Log)
		assert.Nil(t, cfg.BugSnag)
	})

	_, err = f.WriteString(`
log:
  level: debug
  fields:
    string: value
    int: 4
`)
	require.NoError(t, err)

	t.Run("set log field", func(t *testing.T) {
		cfg := RootConfig{Log: DefaultLoggingConfig()}
		require.NoError(t, args.load(&cfg))

		// retains original value
		assert.True(t, cfg.Log.QuoteEmptyFields)
		assert.Equal(t, "debug", cfg.Log.Level)
		require.Len(t, cfg.Log.Fields, 2)
		assert.Equal(t, "value", cfg.Log.Fields["string"])
		assert.Equal(t, 4, cfg.Log.Fields["int"])
	})
}
