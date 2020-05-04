package config

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestEnvBinding(t *testing.T) {
	t.Parallel()

	expected := &Config{Name: "Bob", Age: 32, Banana: false, Personality: Personality{Funny: true}}
	os.Setenv("SPEC_FUNNY", "true")

	v := viper.New()
	f := flags()
	env(v)

	actual, _ := configure(v, f, nil)

	assert.Equal(t, expected, actual)
}

func TestConfigure(t *testing.T) {
	t.Parallel()

	tests := []struct{
		description string
		env [][]string
		flags []string
		expected *Config
	} {
		{
			description: "Defaults",
			expected: &Config{Name: "Bob", Age: 32, Banana: false, Personality: Personality{Funny: false}},
		},
		{
			description: "Name and banana changed by env",
			env: [][]string{{"name", "Tim"}, {"banana", "true"}},
			expected: &Config{Name: "Tim", Age: 32, Banana: true, Personality: Personality{Funny: false}},
		},
		{
			description: "Name and age changed by flags",
			flags: []string{"-nJoe", "-a13"},
			expected: &Config{Name: "Joe", Age: 13, Banana: false, Personality: Personality{Funny: false}},
		},
		{
			description: "Flags beat env",
			env: [][]string{{"name", "Tim"}, {"banana", "true"}},
			flags: []string{"-nJoe", "-a13"},
			expected: &Config{Name: "Joe", Age: 13, Banana: true, Personality: Personality{Funny: false}},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			v := viper.New()
			f := flags()

			prefix := strings.ToUpper(fmt.Sprintf("T%s", strings.ReplaceAll(uuid.New().String(), "-", "")))
			v.SetEnvPrefix(prefix)

			for _, env := range test.env {
				_ = os.Setenv(strings.ToUpper(fmt.Sprintf("%s_%s", prefix, env[0])), env[1])
			}

			_ = f.Parse(test.flags)

			actual, _ := configure(v, f, test.env)

			assert.Equal(t, test.expected, actual)
		})
	}
}
