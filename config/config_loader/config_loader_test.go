package configloader_test

import (
	"os"
	"path/filepath"
	"testing"

	"frisboo-bank/pkg/config/config_loader/config"
	"frisboo-bank/pkg/config/config_loader/contracts"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/syserrors"

	configloader "frisboo-bank/pkg/config/config_loader"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type appConfig struct {
	Name        string `mapstructure:"name"`
	Description string `mapstructure:"description"`
}

type httpServerConfig struct {
	Instances map[string]struct {
		Type   string `mapstructure:"type"`
		Port   string `mapstructure:"port"`
		Logger string `mapstructure:"logger"`
	} `mapstructure:"instances"`
}

type sampleConfig struct {
	App         appConfig        `mapstructure:"app"`
	HTTPServers httpServerConfig `mapstructure:"httpServers"`
}

func TestConfigLoader_Load_Success(t *testing.T) {
	tmp := t.TempDir()

	writeFile(t, tmp, "application.development.yaml", `
app:
  name: my-service
  description: my-service description
httpServers:
  instances:
    main:
      type: gin
      port: "8000"
      logger: http-server
`)

	loader := newLoader(t, tmp, "application")

	var cfg sampleConfig
	require.NoError(t, loader.Load(environment.Environments.DEVELOPMENT, &cfg))

	assert.Equal(t, "my-service", cfg.App.Name)
	assert.Equal(t, "my-service description", cfg.App.Description)

	main, ok := cfg.HTTPServers.Instances["main"]
	assert.True(t, ok)
	assert.Equal(t, "gin", main.Type)
	assert.Equal(t, "8000", main.Port)
	assert.Equal(t, "http-server", main.Logger)
}

func TestConfigLoader_LoadKey_Success(t *testing.T) {
	tmp := t.TempDir()

	writeFile(t, tmp, "application.development.yaml", `
app:
  name: my-service
  description: my-service description
httpServers:
  instances:
    main:
      type: gin
      port: "8000"
      logger: http-server
`)

	loader := newLoader(t, tmp, "application")

	var cfg httpServerConfig
	require.NoError(t, loader.LoadKey(environment.Environments.DEVELOPMENT, &cfg, "httpServers"))

	main, ok := cfg.Instances["main"]
	assert.True(t, ok)
	assert.Equal(t, "gin", main.Type)
}

func TestConfigLoader_HasKey(t *testing.T) {
	t.Parallel()
	tmp := t.TempDir()
	writeFile(t, tmp, "application.development.yaml", `
app:
  name: a
httpServers: {}
`)
	loader := newLoader(t, tmp, "application")

	ok, err := loader.HasKey(environment.Environments.DEVELOPMENT, "app")
	require.NoError(t, err)
	assert.True(t, ok)

	ok, err = loader.HasKey(environment.Environments.DEVELOPMENT, "httpServers")
	require.NoError(t, err)
	assert.True(t, ok)

	ok, err = loader.HasKey(environment.Environments.DEVELOPMENT, "doesNotExist")
	require.NoError(t, err)
	assert.False(t, ok)
}

func TestConfigLoader_Load_TargetMustBePointer(t *testing.T) {
	tmpDir := t.TempDir()

	loader := newLoader(t, tmpDir, "application")

	var loadedNotPtr sampleConfig
	err := loader.Load(environment.Environments.DEVELOPMENT, loadedNotPtr)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "target must be a pointer")
}

func TestConfigLoader_LoadKey_TargetMustBePointer(t *testing.T) {
	tmpDir := t.TempDir()

	loader := newLoader(t, tmpDir, "application")

	var loadedNotPtr sampleConfig
	err := loader.LoadKey(environment.Environments.DEVELOPMENT, loadedNotPtr, "app")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "target must be a pointer")
}

func TestConfigLoader_LoadKey_MissingKey(t *testing.T) {
	tmp := t.TempDir()

	writeFile(t, tmp, "application.development.yaml", "app: { name: svc }")

	loader := newLoader(t, tmp, "application")

	var httpCfg httpServerConfig
	err := loader.LoadKey(environment.Environments.DEVELOPMENT, &httpCfg, "httpServers")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "required key httpServers not found")
}

func TestConfigLoader_MissingConfigName(t *testing.T) {
	_, err := configloader.New(
		config.ConfigName(""),
		config.ConfigPath(t.TempDir()),
	)

	require.Error(t, err)
	assert.Contains(t, syserrors.Cause(err).Error(), "ConfigName: cannot be blank.")
}

func TestConfigLoader_MissingConfigPath(t *testing.T) {
	_, err := configloader.New(
		config.ConfigName("app"),
		config.ConfigPath(""),
	)

	require.Error(t, err)
	assert.Contains(t, syserrors.Cause(err).Error(), "ConfigPath: cannot be blank.")
}

func writeFile(t *testing.T, dist string, filename string, content string) {
	t.Helper()
	fullDist := filepath.Join(dist, filename)
	require.NoError(t, os.WriteFile(fullDist, []byte(content), 0o600))
}

func newLoader(t *testing.T, dist string, fileName string, extraOpts ...config.Option) contracts.ConfigLoader {
	t.Helper()

	opts := []config.Option{
		config.Debug(false),
		config.ConfigPath(dist),
		config.ConfigName(fileName),
	}
	opts = append(opts, extraOpts...)

	cl, err := configloader.New(opts...)
	require.NoError(t, err)

	return cl
}
