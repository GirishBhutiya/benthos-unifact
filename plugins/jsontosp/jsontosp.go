package jsontosp

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/benthosdev/benthos/v4/public/bloblang"
)

var JSONObjectSpec = bloblang.NewPluginSpec().
	Param(bloblang.NewStringParam("jsonstr"))

func init() {

	err := bloblang.RegisterFunctionV2(
		"jsontosp", JSONObjectSpec,
		func(args *bloblang.ParsedParams) (bloblang.Function, error) {
			jsonstr, err := args.GetString("jsonstr")
			if err != nil {
				return nil, err
			}
			return func() (interface{}, error) {
				var obj map[string]any
				err := json.Unmarshal([]byte(jsonstr), &obj)
				if err != nil {
					return nil, err
				}

				var sb strings.Builder
				for key, value := range obj {
					_, err := strconv.ParseFloat(value.(string), 64)
					if err == nil {
						sb.WriteString(fmt.Sprintf("%s%s=%s,", "@", key, value))
						continue
					}
					_, err = strconv.ParseBool(value.(string))
					if err == nil {
						sb.WriteString(fmt.Sprintf("%s%s=%s,", "@", key, value))
						continue
					}
					sb.WriteString(fmt.Sprintf("%s%s='%s',", "@", key, value))
				}
				return sb.String(), nil
			}, nil
		})
	if err != nil {
		panic(err)
	}
}
