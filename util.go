package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func parseInt(string string) (int64, error) {
	integer, err := strconv.ParseInt(string, 10, 64)
	if err != nil {
		return 0, err
	}

	return integer, nil
}

func parseFloat(string string) (float64, error) {
	number, err := strconv.ParseFloat(string, 64)
	if err != nil {
		return 0, err
	}

	return number, nil
}

func jsonSet(json []byte, path string, value interface{}) ([]byte, error) {
	keys := strings.Split(path, ".")
	path = ""
	for i := 0; i < len(keys); i++ {
		key := keys[i]
		if key == "#" {
			key = strconv.Itoa(len(gjson.GetBytes(json, path).Array()) - 1)
		}

		if len(path) > 0 {
			path += "."
		}
		path += key
	}
	return sjson.SetBytes(json, path, value)
}

func serializeQuery(params []byte) url.Values {
  serialized := url.Values{}

  var serialize func(value gjson.Result, path string)
  serialize = func(res gjson.Result, path string) {
    if res.IsObject() {
      for key, value := range res.Map() {
        newPath := path
        if len(newPath) == 0 {
          newPath += key
        } else {
          newPath = "[" + key + "]"
        }

        serialize(value, newPath)
      }
    } else if res.IsArray() {
      for _, value := range res.Array() {
        serialize(value, path)
      }
    } else {
      serialized.Add(path, res.String())
    }
  }
  serialize(gjson.GetBytes(params, "@this"), "")

  for key, values := range serialized {
    serialized.Set(key, strings.Join(values, ","))
  }

  return serialized
}

func serializeHeader(params []byte) http.Header {
  serialized := http.Header{}

  var serialize func(value gjson.Result, path string)
  serialize = func(res gjson.Result, path string) {
    if res.IsObject() {
      for key, value := range res.Map() {
        newPath := path
        if len(newPath) > 0 {
          newPath += "."
        }
        newPath += key

        serialize(value, newPath)
      }
    } else if res.IsArray() {
      for _, value := range res.Array() {
        serialize(value, path)
      }
    } else {
      serialized.Add(path, res.String())
    }
  }
  serialize(gjson.GetBytes(params, "@this"), "")

  return serialized
}

func getStdInput() []byte {
	if !isInputPiped() {
		return nil
	}
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return data
}

func isInputPiped() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}
