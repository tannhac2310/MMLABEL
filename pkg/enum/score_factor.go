package enum

import (
	"bytes"
	"fmt"
)

type ScoreFactor uint8

const (
	ScoreFactorNone ScoreFactor = iota
	ScoreFactor1
	ScoreFactor2
)

var ScoreFactorName = map[ScoreFactor]string{
	ScoreFactorNone: "none",
	ScoreFactor1:    "factor1",
	ScoreFactor2:    "factor2",
}

var ScoreFactorValue = func() map[string]ScoreFactor {
	value := map[string]ScoreFactor{}
	for k, v := range ScoreFactorName {
		value[v] = k
		value[fmt.Sprintf("%v", k)] = k
	}

	return value
}()

func (e ScoreFactor) MarshalJSON() ([]byte, error) {
	v, ok := ScoreFactorName[e]
	if !ok {
		return nil, fmt.Errorf("invalid values of ScoreFactor")
	}

	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(v)
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (e *ScoreFactor) UnmarshalJSON(data []byte) error {
	data = bytes.Trim(data, "\"")
	v, ok := ScoreFactorValue[string(data)]
	if !ok {
		return fmt.Errorf("enum '%s' is not register, must be one of: %v", data, e.EnumDescriptions())
	}

	*e = v

	return nil
}

func (*ScoreFactor) EnumDescriptions() []string {
	vals := []string{}

	for _, name := range ScoreFactorName {
		vals = append(vals, name)
	}

	return vals
}
