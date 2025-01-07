package utils

import (
	"encoding/base64"
	acbitypes "github.com/cometbft/cometbft/abci/types"
	"regexp"
)

func ResolveBase64Events(events []acbitypes.Event) []acbitypes.Event {
	base64Pattern := regexp.MustCompile(`^([A-Za-z\d+/]{4})*([A-Za-z\d+/]{3}=|[A-Za-z\d+/]{2}==)?$`)

	isBase64 := func(str string) bool {
		return base64Pattern.MatchString(str)
	}

	out := make([]acbitypes.Event, 0, len(events))
	for _, ie := range events {
		oe := acbitypes.Event{
			Type:       ie.Type,
			Attributes: make([]acbitypes.EventAttribute, 0, len(ie.Attributes)),
		}

		for _, ia := range ie.Attributes {
			oa := acbitypes.EventAttribute{
				Key:   ia.Key,
				Value: ia.Value,
			}

			sourceIsBase64 := (ia.Key == "" || isBase64(ia.Key)) && (ia.Value == "" || isBase64(ia.Value))

			if sourceIsBase64 {
				var key, value []byte
				var errK, errV error
				if ia.Key != "" {
					key, errK = base64.StdEncoding.DecodeString(ia.Key)
				}
				if ia.Value != "" {
					value, errV = base64.StdEncoding.DecodeString(ia.Value)
				}

				if errK == nil && errV == nil {
					oa.Key = string(key)
					oa.Value = string(value)
				}
			}

			oe.Attributes = append(oe.Attributes, oa)
		}

		out = append(out, oe)
	}

	return out
}
