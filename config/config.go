package config

import (
	"os"
	"reflect"
	"strings"
)

var Version string = "Develop"

// Config contains the main configuration values.
type Config struct {
	GitHub   github
	Workflow workflow
}

type github struct {
	Token      string
	Repository string
}

type workflow struct {
	Files string
}

// Init creates a new config based on parsed environment variables.
func Init() (*Config, error) {
	return envParseConfig(&Config{
		GitHub: github{
			Token: "",
		},
		Workflow: workflow{
			Files: "",
		},
	}), nil
}

func envParseConfig(in *Config) (*Config) {
	numSubStructs := reflect.ValueOf(in).Elem().NumField()
	for i := 0; i < numSubStructs; i++ {
		iter := reflect.ValueOf(in).Elem().Field(i)
		subStruct := strings.ToUpper(iter.Type().Name())

		structType := iter.Type()
		for j := 0; j < iter.NumField(); j++ {
			fieldVal := iter.Field(j).String()
			fieldName := structType.Field(j).Name
			for _, prefix := range []string{"CAFFEINATE", "INPUT"} {
				evName := prefix + "_" + subStruct + "_" + strings.ToUpper(fieldName)
				evVal, evExists := os.LookupEnv(evName)
				if evExists && evVal != fieldVal {
					iter.FieldByName(fieldName).SetString(evVal)
				}
			}
		}
	}
	return in
}
