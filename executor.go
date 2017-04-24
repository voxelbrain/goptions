package goptions

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
)

var (
	//ErrNoCommand This should never happen.
	ErrNoCommand = errors.New("No command specified")
)

//Executer used to allow auto execution of commands.
type Executer interface {
	Execute(args []string) error
}

// ParseAndExecute ...
func ParseAndExecute(v Executer) error {
	globalFlagSet = NewFlagSet(filepath.Base(os.Args[0]), v)
	if err := globalFlagSet.Parse(os.Args[1:]); err != nil {
		return err
	}
	exe, args, err := getExecuterArgs(v)
	if err != nil {
		return err
	}
	return exe.Execute(args)
}

// ParseExecuteAndFail ...
func ParseExecuteAndFail(v Executer) {
	err := ParseAndExecute(v)
	if err != nil {
		errCode := 0
		if err != ErrHelpRequest {
			errCode = 1
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		}
		PrintHelp()
		os.Exit(errCode)
	}
}

//GetCommands returns a list of all commands avaliable for the executer.
func GetCommands(v Executer) []string {
	if globalFlagSet == nil {
		panic("You must call parse before you can call get commands")
	}
	commands := []string{}
	for k := range globalFlagSet.Verbs {
		commands = append(commands, k)
	}
	return commands
}

func getExecuterArgs(v Executer) (Executer, []string, error) {
	var cmd string
	var args []string
	structValue := reflect.ValueOf(v)
	if structValue.Kind() == reflect.Ptr {
		structValue = structValue.Elem()
	}

	if structValue.Kind() != reflect.Struct {
		panic("Value type is not a pointer to a struct")
	}

	var i int
	for i = 0; i < structValue.Type().NumField(); i++ {
		if StartsWithLowercase(structValue.Type().Field(i).Name) {
			continue
		}
		fieldValue := structValue.Field(i)
		if fieldValue.Type().Name() == "Remainder" {
			switch v := fieldValue.Interface().(type) {
			case Remainder:
				args = v
			default:
				fmt.Printf("%##v", fieldValue.Interface())
				panic("Remainder was not a []string")

			}
		}
		if fieldValue.Type().Name() == "Verbs" {
			cmd = fmt.Sprint(fieldValue)
			break
		}
	}
	if cmd == "" {
		return v, args, nil
	}

	for i++; i < structValue.Type().NumField(); i++ {
		// fieldValue := structValue.Field(i)
		tag := structValue.Type().Field(i).Tag.Get("goptions")
		if tag == cmd {
			x := structValue.Field(i).Interface()

			switch y := x.(type) {
			case Executer:
				return getExecuterArgs(y)
			default:
				return nil, nil, fmt.Errorf("CMD: %s is not an Executer", structValue.Field(i).Type().Name())
			}
		}
	}
	return nil, nil, ErrNoCommand
}
