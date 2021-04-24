package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	nodeconfig "github.com/harmony-one/harmony/internal/configs/node"
	"github.com/stretchr/testify/assert"
)

var (
	trueBool = true
)

type testCfgOpt func(config *HarmonyConfig)

func makeTestConfig(nt nodeconfig.NetworkType, opt testCfgOpt) HarmonyConfig {
	cfg := GetDefaultHmyConfigCopy(nt)
	if opt != nil {
		opt(&cfg)
	}
	return cfg
}

var testBaseDir = ".testdata"

func init() {
	if _, err := os.Stat(testBaseDir); os.IsNotExist(err) {
		os.MkdirAll(testBaseDir, 0777)
	}
}

func TestV1_0_0Config(t *testing.T) {
	testConfig := `Version = "1.0.4"

[BLSKeys]
  KMSConfigFile = ""
  KMSConfigSrcType = "shared"
  KMSEnabled = false
  KeyDir = "./.hmy/blskeys"
  KeyFiles = []
  MaxKeys = 10
  PassEnabled = true
  PassFile = ""
  PassSrcType = "auto"
  SavePassphrase = false

[General]
  DataDir = "./"
  IsArchival = false
  NoStaking = false
  NodeType = "validator"
  ShardID = -1

[HTTP]
  Enabled = true
  IP = "127.0.0.1"
  Port = 9500

[Log]
  FileName = "harmony.log"
  Folder = "./latest"
  RotateSize = 100
  Verbosity = 3

[Network]
  BootNodes = ["/dnsaddr/bootstrap.t.hmny.io"]
  DNSPort = 9000
  DNSZone = "t.hmny.io"
  LegacySyncing = false
  NetworkType = "mainnet"

[P2P]
  KeyFile = "./.hmykey"
  Port = 9000

[Pprof]
  Enabled = false
  ListenAddr = "127.0.0.1:6060"

[TxPool]
  BlacklistFile = "./.hmy/blacklist.txt"

[Sync]
  Downloader = false
  Concurrency = 6
  DiscBatch = 8
  DiscHardLowCap = 6
  DiscHighCap = 128
  DiscSoftLowCap = 8
  InitStreams = 8
  LegacyClient = true
  LegacyServer = true
  MinPeers = 6

[WS]
  Enabled = true
  IP = "127.0.0.1"
  Port = 9800`
	testDir := filepath.Join(testBaseDir, t.Name())
	os.RemoveAll(testDir)
	os.MkdirAll(testDir, 0777)
	file := filepath.Join(testDir, "test.config")
	err := ioutil.WriteFile(file, []byte(testConfig), 0644)
	if err != nil {
		t.Fatal(err)
	}
	config, err := LoadHarmonyConfig(file)
	if err != nil {
		t.Fatal(err)
	}
	defConf := GetDefaultHmyConfigCopy(nodeconfig.Mainnet)
	if config.HTTP.RosettaEnabled {
		t.Errorf("Expected rosetta http server to be disabled when loading old config")
	}
	if config.General.IsOffline {
		t.Errorf("Expect node to de online when loading old config")
	}
	if config.P2P.IP != defConf.P2P.IP {
		t.Errorf("Expect default p2p IP if old config is provided")
	}
	if config.Version != "1.0.4" {
		t.Errorf("Expected config version: 1.0.4, not %v", config.Version)
	}
	config.Version = defConf.Version // Shortcut for testing, value checked above
	assert.Equal(t, defConf, config)
}

func TestPersistConfig(t *testing.T) {
	testDir := filepath.Join(testBaseDir, t.Name())
	os.RemoveAll(testDir)
	os.MkdirAll(testDir, 0777)

	tests := []struct {
		config HarmonyConfig
	}{
		{
			config: makeTestConfig("mainnet", nil),
		},
		{
			config: makeTestConfig("devnet", nil),
		},
		{
			config: makeTestConfig("mainnet", func(cfg *HarmonyConfig) {
				consensus := GetDefaultConsensusConfigCopy()
				cfg.Consensus = &consensus

				devnet := GetDefaultDevnetConfigCopy()
				cfg.Devnet = &devnet

				revert := GetDefaultRevertConfigCopy()
				cfg.Revert = &revert

				webHook := "web hook"
				cfg.Legacy = &LegacyConfig{
					WebHookConfig:         &webHook,
					TPBroadcastInvalidTxn: &trueBool,
				}

				logCtx := GetDefaultLogContextCopy()
				cfg.Log.Context = &logCtx
			}),
		},
	}
	for i, test := range tests {
		file := filepath.Join(testDir, fmt.Sprintf("%d.conf", i))

		if err := WriteHarmonyConfigToFile(test.config, file); err != nil {
			t.Fatal(err)
		}
		config, err := LoadHarmonyConfig(file)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, test.config, config, "test %d: configs should match", i)
	}
}