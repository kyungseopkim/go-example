package main

import "encoding/json"

type ArxmlConfig struct {
	Filename   string     `json:"filename"`
}

func (arxml ArxmlConfig) String () string  {
	content, err := json.Marshal(arxml)
	if err != nil {
		return ""
	}
	return string(content)
}
