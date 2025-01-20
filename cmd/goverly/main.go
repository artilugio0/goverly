package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/artilugio0/goverly"
)

//go:embed wasm/*
var content embed.FS

func main() {
	configFile := flag.String("f", "config.json", "configuration file")
	flag.Parse()
	args := flag.Args()

	subcommand := ""
	if len(args) > 0 {
		subcommand = args[0]
	}

	switch subcommand {
	case "overlay":
		goverly.ServeOverlay(*configFile, content)
	case "config":
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "config subcommand not specified\n")
			os.Exit(1)
		}

		switch args[1] {
		case "set":
			if len(args) < 4 {
				fmt.Fprintf(os.Stderr, "config set values missing\n")
				os.Exit(1)
			}

			if err := configSet(*configFile, args[2], args[3]); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
		default:
			if err := widgetSpecificConfigSet(*configFile, args[1], args[2:]); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
		}
	}
}

func configSet(configFile, configPath, value string) error {
	config := goverly.Config{}

	configBytes, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}

	if json.Unmarshal(configBytes, &config); err != nil {
		return err
	}

	configFields := strings.Split(configPath, ".")
	if len(configFields) <= 1 {
		return fmt.Errorf("config path for widget '%s' does not specify an attribute", configFields[0])
	}

	widget, ok := config.Widgets[configFields[0]]
	if !ok {
		return fmt.Errorf("widget '%s' is not defined", configFields[0])
	}

	pathIndex := 1
	widgetValue := reflect.ValueOf(widget)
	st := reflect.TypeOf(widget)
OUTER:
	for {
		isPointer := st.Kind() == reflect.Pointer
		if isPointer {
			st = st.Elem()
		}

		if st.Kind() != reflect.Struct {
			return fmt.Errorf("invalid type for widget '%s'", configFields[0])
		}

		for i := range st.NumField() {
			field := st.Field(i)
			tag := field.Tag.Get("json")
			if tag != configFields[pathIndex] {
				continue
			}

			if isPointer {
				widgetValue = widgetValue.Elem()
			}
			f := widgetValue.Field(i)
			if !f.CanSet() {
				return fmt.Errorf("'%s' is not assignable", configPath)
			}

			switch field.Type.Kind() {
			case reflect.String:
				f.SetString(value)
			case reflect.Int, reflect.Int64:
				intValue, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return err
				}
				f.SetInt(intValue)
			case reflect.Bool:
				f.SetBool(value == "true")
			case reflect.Slice:
				if len(configFields) <= pathIndex+1 {
					return fmt.Errorf("missing array index in the specified config path")
				}

				if len(configFields) <= pathIndex+2 {
					return fmt.Errorf("missing field after array index in the specified config path")
				}

				sliceIndex, err := strconv.Atoi(configFields[pathIndex+1])
				if err != nil {
					return fmt.Errorf("invalid array index '%s'", configFields[pathIndex+1])
				}

				if f.Len() <= sliceIndex {
					return fmt.Errorf("array index '%d' out of bounds", sliceIndex)
				}

				widgetValue = f.Index(sliceIndex)
				st = widgetValue.Type()
				pathIndex += 2
				continue OUTER
			default:
				return fmt.Errorf("unsupported type for field '%s'", configPath)
			}

			break OUTER
		}

		return fmt.Errorf("invalid path for widget '%s'", configFields[0])
	}

	configBytes, err = json.MarshalIndent(config, "", "\t")
	if err != nil {
		return err
	}

	if err := os.WriteFile(configFile, configBytes, 0600); err != nil {
		return err
	}

	return nil
}

func widgetSpecificConfigSet(configFile, widgetName string, args []string) error {
	config := goverly.Config{}

	configBytes, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}

	if json.Unmarshal(configBytes, &config); err != nil {
		return err
	}

	widget, ok := config.Widgets[widgetName]
	if !ok {
		return fmt.Errorf("widget '%s' is not defined", widgetName)
	}

	ccWidget, ok := widget.(WidgetCustomConfig)
	if !ok {
		return fmt.Errorf("invalid operation for widget '%s'", widgetName)
	}

	if err := ccWidget.ApplyCustomConfig(args); err != nil {
		return err
	}

	configBytes, err = json.MarshalIndent(config, "", "\t")
	if err != nil {
		return err
	}

	if err := os.WriteFile(configFile, configBytes, 0600); err != nil {
		return err
	}

	return nil
}

type WidgetCustomConfig interface {
	ApplyCustomConfig([]string) error
}
