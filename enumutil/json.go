package main

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
)

type KeyValue struct {
	Key   string
	Value interface{}
}

type OrderedMap []KeyValue

func (order OrderedMap) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString("{")
	for i, elem := range order {
		if i != 0 {
			buf.WriteString(",")
		}

		// marshal key
		key, err := json.Marshal(elem.Key)
		if err != nil {
			return nil, err
		}

		buf.Write(key)
		buf.WriteString(":")

		// marshal value
		v, err := json.Marshal(elem.Value)
		if err != nil {
			return nil, err
		}
		buf.Write(v)
	}

	buf.WriteString("}")
	return buf.Bytes(), nil
}

func GenerateJSON(filename string, values []string) (OrderedMap, error) {
	list := make(OrderedMap, len(values))

	for i, v := range values {
		list[i].Key = v
		list[i].Value = []string{}
	}

	if _, err := os.Stat(filename); err == nil {
		file, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			return nil, err
		}

		var result map[string]interface{}
		json.Unmarshal([]byte(content), &result)

		for i, elem := range list {
			v, ok := result[elem.Key]
			if !ok {
				continue
			}

			list[i].Value = v
			delete(result, elem.Key)
		}

		for key, value := range result {
			list = append(list, KeyValue{
				key,
				value,
			})
		}
	}

	return list, nil
}

func PrintOutJSON(filename string, list OrderedMap) error {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	return enc.Encode(list)
}
