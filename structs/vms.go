package structs

// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    VMAction, err := UnmarshalWelcome(bytes)
//    bytes, err = welcome.Marshal()
import "encoding/json"

func UnmarshalWelcome(data []byte) (VMAction, error) {
	var r VMAction
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *VMAction) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type VMAction struct {
	Action string `json:"action"`
	VM     string `json:"vm"`
}
