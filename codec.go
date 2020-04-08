package tapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

// Codec 编解码器接口
type Codec interface {
	Encode(v interface{}) ([]byte, error)
	Decode(data []byte, v interface{}) error
}

// json 编解码器
type jsonc struct {
}

// JSONC json 编解码器
var JSONC jsonc

func (jsonc) Encode(v interface{}) ([]byte, error) {
	return json.MarshalIndent(v, "", "\t")
}

func (jsonc) Decode(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// toml 编解码器
type tomlc struct {
}

// TOMLC toml 编解码器
var TOMLC tomlc

func (tomlc) Encode(v interface{}) ([]byte, error) {
	var b bytes.Buffer
	enc := toml.NewEncoder(&b)
	enc.Indent = "\t"
	err := enc.Encode(v)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (tomlc) Decode(data []byte, v interface{}) error {
	_, err := toml.Decode(string(data), v)
	return err
}

func loadFromFile(c Codec, filename string, v interface{}) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}
	if err = c.Decode(data, v); err != nil {
		return fmt.Errorf("decode: %w", err)
	}
	return nil
}

func writeToFile(c Codec, filename string, v interface{}) error {
	data, err := c.Encode(v)
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}
	err = ioutil.WriteFile(filename, data, 0666)
	if err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	return nil
}
