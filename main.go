package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Service represents a single service in the docker-compose file
type Service struct {
	Image       string   `yaml:"image"`
	Environment []string `yaml:"environment,omitempty"`
}

// ComposeFile represents the entire docker-compose file
type ComposeFile struct {
	Version  string              `yaml:"version"`
	Services map[string]*Service `yaml:"services"`
}

func run(reader io.Reader, file string) (string, error) {
	// Read docker-compose.yaml file
	data, err := os.ReadFile(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading docker-compose.yaml: %v\n", err)
		return "", err
	}

	// Parse the docker-compose.yaml file
	var compose ComposeFile
	if err := yaml.Unmarshal(data, &compose); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing docker-compose.yaml: %v\n", err)
		return "", err
	}

	// Read environment variables from stdin
	envVars, err := readEnvVars(reader)
	if err != nil {
		return "", err
	}

	// Replace process_env with the environment variables
	for _, service := range compose.Services {
		if len(service.Environment) > 0 {
			// Find the index of the process_env environment variable
			index := -1
			for i, envVar := range service.Environment {
				if strings.HasPrefix(envVar, "ALL_ENV_VARS") {
					index = i
					break
				}
			}
			if index == -1 {
				continue
			}

			or := service.Environment[0:index]
			or2 := make([]string, len(service.Environment)-(index+1))
			if len(service.Environment) > index+1 {
				t := service.Environment[index+1 : len(service.Environment)]
				copy(or2, t)
			}

			service.Environment = append(or, envVars...)
			service.Environment = append(service.Environment, or2...)
		}
	}

	// Unmarshal the YAML into a map
	var configMap map[string]interface{}
	err = yaml.Unmarshal(data, &configMap)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		return "", err
	}

	for key, val := range configMap {
		if key == "services" {
			services := val.(map[string]interface{})
			for name, service := range services {
				serviceMap := service.(map[string]interface{})
				envs := compose.Services[name].Environment
				serviceMap["environment"] = envs
			}
		}
	}

	// Marshal the modified compose file back to YAML
	output, err := yaml.Marshal(&configMap)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling modified docker-compose.yaml: %v\n", err)
		return "", err
	}

	// Write the modified YAML to stdout
	return string(output), nil
}

func main() {
	output, err := run(os.Stdin, "docker-compose.yaml")
	if err != nil {
		os.Exit(1)
	}
	fmt.Println(string(output))
}

// readEnvVars reads environment variables from stdin
func readEnvVars(reader io.Reader) ([]string, error) {
	var envVars []string
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		_, after, found := strings.Cut(line, "=")
		if found && after == "" {
			fmt.Fprintf(os.Stderr, "Skip: Environment variable %q has an empty value\n", line)
			continue
		}
		envVars = append(envVars, line)

	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
		return nil, err
	}
	return envVars, nil
}
