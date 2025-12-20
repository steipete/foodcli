package foodora

import (
	"encoding/json"
	"strconv"
)

// FlexibleInt decodes ints that sometimes come back as strings (API drift).
type FlexibleInt int

func (i *FlexibleInt) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		*i = 0
		return nil
	}
	if len(b) > 0 && b[0] == '"' {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		if s == "" {
			*i = 0
			return nil
		}
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		*i = FlexibleInt(v)
		return nil
	}

	var n int64
	if err := json.Unmarshal(b, &n); err == nil {
		*i = FlexibleInt(n)
		return nil
	}

	var f float64
	if err := json.Unmarshal(b, &f); err != nil {
		return err
	}
	*i = FlexibleInt(int64(f))
	return nil
}
