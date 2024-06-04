package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	input := "ENVIRONMENT1=hello\nENVIRONMENT2=world\nENVIRONMENT3=\n"
	reader := strings.NewReader(input)

	output, err := run(reader, "docker-compose.yaml")
	assert.NoError(t, err)

	b, err := os.ReadFile("docker-compose.yaml.expected")
	assert.NoError(t, err)

	assert.Equal(t, string(b), output)
}

func TestRunNoFile(t *testing.T) {
	reader := strings.NewReader("")

	_, err := run(reader, "nofile")
	assert.Error(t, err)
	assert.ErrorContains(t, err, "open nofile: no such file or directory")
}

func TestRunErrorStdin(t *testing.T) {
	reader := iotest.ErrReader(fmt.Errorf("error reading string from stdin"))

	_, err := run(reader, "docker-compose.yaml.expected")
	assert.Error(t, err)
	assert.ErrorContains(t, err, "error reading string from stdin")
}

func TestReadEnvVars(t *testing.T) {
	input := "ENVIRONMENT1=hello\nENVIRONMENT2=world\nENVIRONMENT3=\n\n"
	reader := strings.NewReader(input)

	envVars, err := readEnvVars(reader)
	assert.NoError(t, err)

	assert.Equal(t, 2, len(envVars))
	assert.Equal(t, "ENVIRONMENT1=hello", envVars[0])
	assert.Equal(t, "ENVIRONMENT2=world", envVars[1])
}

func TestReadEnvVarsErrors(t *testing.T) {
	reader := iotest.ErrReader(fmt.Errorf("error reading string from stdin"))

	envVars, err := readEnvVars(reader)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "error reading string from stdin")
	assert.Nil(t, envVars)
}
