package utils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/KushamiNeko/go_fun/utils/pretty"
)

func KeyValueInput(color pretty.HexColor, message string) (map[string]string, error) {

	//const (
	//rePattern = `([^;]*)=([^;]*)`
	//)

	inputs := pretty.ColorInput(color, message)
	inputs = strings.TrimSpace(inputs)

	if inputs == "" {
		return nil, fmt.Errorf("empty inputs")
	}

	pair, err := keyValuePair(inputs)
	if err != nil {
		return nil, err
	}

	return pair, nil

	//pair := make(map[string]string)

	//re := regexp.MustCompile(rePattern)

	//match := re.FindAllStringSubmatch(inputs, -1)

	//for _, m := range match {
	//k := strings.TrimSpace(m[1])
	//v := strings.TrimSpace(m[2])

	//if k == "" || v == "" {
	//return nil, fmt.Errorf("empty key or value: %s=%s", k, v)
	//}

	//pair[k] = v
	//}

	//if len(pair) == 0 {
	//return nil, fmt.Errorf("empty inputs")
	//}

	//return pair, nil
}

func keyValuePair(inputs string) (map[string]string, error) {

	const (
		rePattern = `([^;]*)=([^;]*)`
	)

	pair := make(map[string]string)

	re := regexp.MustCompile(rePattern)

	match := re.FindAllStringSubmatch(inputs, -1)

	for _, m := range match {
		k := strings.TrimSpace(m[1])
		v := strings.TrimSpace(m[2])

		if k == "" || v == "" {
			return nil, fmt.Errorf("empty key or value: %s=%s", k, v)
		}

		pair[k] = v
	}

	if len(pair) == 0 {
		return nil, fmt.Errorf("empty inputs")
	}

	return pair, nil
}

func InputsAbbreviation(inputs, abbreviation map[string]string) map[string]string {
	n := make(map[string]string)

	for k, v := range inputs {
		fk, ok := abbreviation[k]
		if ok {
			n[fk] = v
		} else {
			n[k] = v
		}
	}

	return n
}
