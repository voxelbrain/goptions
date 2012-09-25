package goptions

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	_LONG_FLAG_REGEXP     = `--[[:word:]-]+`
	_SHORT_FLAG_REGEXP    = `-[[:alnum:]]`
	_BOOL_OPTION_REGEXP   = `[[:word:]-]+`
	_QUOTED_STRING_REGEXP = `'((?:\\'|[^\\'])+)'`
	_VALUE_OPTION_REGEXP  = `[[:word:]-]+=` + _QUOTED_STRING_REGEXP
)

var (
	optionRegexp = regexp.MustCompile(`^(` + strings.Join([]string{_SHORT_FLAG_REGEXP, _LONG_FLAG_REGEXP, _BOOL_OPTION_REGEXP, _VALUE_OPTION_REGEXP}, "|") + `)(?:,|$)`)
)

func parseTag(tag string) (*Flag, error) {
	f := &Flag{
		Short: make([]string, 0),
		Long:  make([]string, 0),
	}
	for {
		tag = strings.TrimSpace(tag)
		if len(tag) == 0 {
			break
		}
		idx := optionRegexp.FindStringSubmatchIndex(tag)
		if idx == nil {
			return nil, fmt.Errorf("Could not find a valid flag definition at the beginning of \"%s\"", tag)
		}
		option := tag[idx[2]:idx[3]]
		tag = tag[idx[1]:]

		if strings.HasPrefix(option, "--") {
			f.Long = append(f.Long, option[2:])
		} else if strings.HasPrefix(option, "-") {
			f.Short = append(f.Short, option[1:])
		} else if strings.HasPrefix(option, "description=") {
			f.Description = strings.Replace(option[idx[4]:idx[5]], `\`, ``, -1)
		} else if strings.HasPrefix(option, "mutexgroup=") {
			f.MutexGroup = option[idx[4]:idx[5]]
		} else {
			switch option {
			case "accumulate":
				f.Accumulate = true
			case "obligatory":
				f.Obligatory = true
			default:
				return nil, fmt.Errorf("Unknown option %s", option)
			}
		}
	}
	return f, nil
}
