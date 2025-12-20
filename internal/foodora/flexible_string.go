package foodora

import (
	"encoding/json"
	"fmt"
)

// FlexibleString decodes strings that sometimes come back as numbers/bools (API drift).
type FlexibleString string

func (s *FlexibleString) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		*s = ""
		return nil
	}
	if len(b) > 0 && b[0] == '"' {
		var v string
		if err := json.Unmarshal(b, &v); err != nil {
			return err
		}
		*s = FlexibleString(v)
		return nil
	}

	var n json.Number
	if err := json.Unmarshal(b, &n); err == nil {
		*s = FlexibleString(n.String())
		return nil
	}

	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	*s = FlexibleString(fmt.Sprint(v))
	return nil
}
