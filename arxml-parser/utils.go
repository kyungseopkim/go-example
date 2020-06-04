package arxml

import "encoding/json"

func ToJson(data interface{}) string {
    bstr, _ := json.Marshal(data)
    return string(bstr)
}