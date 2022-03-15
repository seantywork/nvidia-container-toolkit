/**
# Copyright (c) 2022, NVIDIA CORPORATION.  All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
**/

package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGerRuntimeConfigWithCustomConfig(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	// By default debug is disabled
	contents := []byte("[nvidia-container-runtime]\ndebug = \"/nvidia-container-toolkit.log\"")
	testDir := filepath.Join(wd, "test")
	filename := filepath.Join(testDir, configFilePath)

	os.Setenv(configOverride, testDir)

	require.NoError(t, os.MkdirAll(filepath.Dir(filename), 0766))
	require.NoError(t, ioutil.WriteFile(filename, contents, 0766))

	defer func() { require.NoError(t, os.RemoveAll(testDir)) }()

	cfg, err := GetRuntimeConfig()
	require.NoError(t, err)
	require.Equal(t, cfg.DebugFilePath, "/nvidia-container-toolkit.log")
}

func TestGerRuntimeConfig(t *testing.T) {
	testCases := []struct {
		description    string
		contents       []string
		expectedError  error
		expectedConfig *RuntimeConfig
	}{
		{
			description: "empty config is default",
			expectedConfig: &RuntimeConfig{
				DebugFilePath: "/dev/null",
			},
		},
		{
			description: "config options set inline",
			contents: []string{
				"nvidia-container-runtime.debug = \"/foo/bar\"",
			},
			expectedConfig: &RuntimeConfig{
				DebugFilePath: "/foo/bar",
			},
		},
		{
			description: "config options set in section",
			contents: []string{
				"[nvidia-container-runtime]",
				"debug = \"/foo/bar\"",
			},
			expectedConfig: &RuntimeConfig{
				DebugFilePath: "/foo/bar",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			reader := strings.NewReader(strings.Join(tc.contents, "\n"))

			cfg, err := getRuntimeConfigFrom(reader)
			if tc.expectedError != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.EqualValues(t, tc.expectedConfig, cfg)
		})
	}
}
