package fbxapi

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/Sirupsen/logrus"
)

func isIn(a string, b []string) bool {
	for _, v := range b {
		if a == v {
			return true
		}
	}
	return false
}

func getData(rawResult json.RawMessage) (data map[string]interface{}) {
	json.Unmarshal(rawResult, &data)
	return
}

func getFirst(rawResult json.RawMessage) map[string]interface{} {
	var data []map[string]interface{}
	json.Unmarshal(rawResult, &data)
	return data[0]
}

func checkOrphans(aStruct interface{}, result json.RawMessage) (bool, []string, []string) {
	t := reflect.TypeOf(aStruct).Elem()

	data := make(map[string]interface{})
	tehStruct := t

	if data != nil {
		switch t.Kind() {
		case reflect.Slice:
			data = getFirst(result)
			tehStruct = t.Elem()
		case reflect.Struct:
			data = getData(result)
		}
	}

	var structFields []string
	for i := 0; i < tehStruct.NumField(); i++ {
		tag := tehStruct.Field(i).Tag.Get("json")
		tagName := strings.SplitN(tag, ",", 2)
		if tagName[0] != "" {
			structFields = append(structFields, tagName[0])
		}
	}

	jsonKeys := make([]string, 0, len(data))
	for k := range data {
		jsonKeys = append(jsonKeys, k)
	}

	var newKeys []string
	var expiredKeys []string

	for _, k := range jsonKeys {
		if !isIn(k, structFields) {
			newKeys = append(newKeys, k)
		}
	}

	for _, k := range structFields {
		if !isIn(k, jsonKeys) {
			expiredKeys = append(expiredKeys, k)
		}
	}

	if len(newKeys) > 0 {
		logrus.Warnf("%s has new fields: %v", tehStruct.Name(), newKeys)
		// spew.Dump(data)
	}

	if len(expiredKeys) > 0 {
		logrus.Warnf("%s has expired fields: %v", tehStruct.Name(), expiredKeys)
	}

	return !(len(expiredKeys) == 0 && len(newKeys) == 0), newKeys, expiredKeys
}

func failOnError(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func EndpointTester(t *testing.T, ep *Endpoint, data interface{}, urlparams map[string]string, body interface{}) {
	resp := new(APIResponse)
	err := testClient.Query(ep).As(urlparams).WithBody(body).Inspect(resp).Do(&data)
	failOnError(t, err)

	checkOrphans(data, resp.Result)
}
