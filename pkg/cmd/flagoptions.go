package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"github.com/stainless-api/stainless-api-cli/internal/apiform"
	"github.com/stainless-api/stainless-api-cli/internal/apiquery"
	"github.com/stainless-api/stainless-api-cli/internal/debugmiddleware"
	"github.com/stainless-api/stainless-api-cli/internal/requestflag"
	"github.com/stainless-api/stainless-api-go/option"

	"github.com/goccy/go-yaml"
	"github.com/urfave/cli/v3"
)

type BodyContentType int

const (
	MultipartFormEncoded BodyContentType = iota
	ApplicationJSON
)

func flagOptions(
	cmd *cli.Command,
	nestedFormat apiquery.NestedQueryFormat,
	arrayFormat apiquery.ArrayQueryFormat,
	bodyType BodyContentType,
) ([]option.RequestOption, error) {
	var options []option.RequestOption
	if cmd.Bool("debug") {
		options = append(options, option.WithMiddleware(debugmiddleware.NewRequestLogger().Middleware()))
	}

	flagContents := requestflag.ExtractRequestContents(cmd)

	var bodyData any
	if isInputPiped() {
		var err error
		pipeData, err := io.ReadAll(os.Stdin)
		if err != nil {
			return nil, err
		}

		if err := yaml.Unmarshal(pipeData, &bodyData); err == nil {
			if bodyMap, ok := bodyData.(map[string]any); ok {
				if flagMap, ok := flagContents.Body.(map[string]any); ok {
					for k, v := range flagMap {
						bodyMap[k] = v
					}
				} else {
					bodyData = flagContents.Body
				}
			} else if flagMap, ok := flagContents.Body.(map[string]any); ok && len(flagMap) > 0 {
				return nil, fmt.Errorf("Cannot merge flags with a body that is not a map: %v", bodyData)
			}
		}
	} else {
		// No piped input, just use body flag values as a map
		bodyData = flagContents.Body
	}

	querySettings := apiquery.QuerySettings{
		NestedFormat: nestedFormat,
		ArrayFormat:  arrayFormat,
	}

	// Add query parameters:
	if values, err := apiquery.MarshalWithSettings(flagContents.Queries, querySettings); err != nil {
		return nil, err
	} else {
		for k, vs := range values {
			if len(vs) == 0 {
				options = append(options, option.WithQueryDel(k))
			} else {
				options = append(options, option.WithQuery(k, vs[0]))
				for _, v := range vs[1:] {
					options = append(options, option.WithQueryAdd(k, v))
				}
			}
		}
	}

	// Add header parameters
	if values, err := apiquery.MarshalWithSettings(flagContents.Headers, querySettings); err != nil {
		return nil, err
	} else {
		for k, vs := range values {
			if len(vs) == 0 {
				options = append(options, option.WithHeaderDel(k))
			} else {
				options = append(options, option.WithHeader(k, vs[0]))
				for _, v := range vs[1:] {
					options = append(options, option.WithHeaderAdd(k, v))
				}
			}
		}
	}

	switch bodyType {
	case MultipartFormEncoded:
		buf := new(bytes.Buffer)
		writer := multipart.NewWriter(buf)

		// For multipart/form-encoded, we need a map structure
		bodyMap, ok := bodyData.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("Cannot send a non-map value to a form-encoded endpoint: %v\n", bodyData)
		}
		if err := apiform.MarshalWithSettings(bodyMap, writer, apiform.FormatComma); err != nil {
			return nil, err
		}
		if err := writer.Close(); err != nil {
			return nil, err
		}
		options = append(options, option.WithRequestBody(writer.FormDataContentType(), buf))
	case ApplicationJSON:
		bodyBytes, err := json.Marshal(bodyData)
		if err != nil {
			return nil, err
		}
		options = append(options, option.WithRequestBody("application/json", bodyBytes))
	default:
		panic("Invalid body content type!")
	}

	return options, nil
}
