// Code generated by jsonenums -type=GamePhase; DO NOT EDIT.

package models

import (
	"encoding/json"
	"fmt"
)

var (
	_GamePhaseNameToValue = map[string]GamePhase{
		"WaitingForPlayers": WaitingForPlayers,
		"Starting":          Starting,
	}

	_GamePhaseValueToName = map[GamePhase]string{
		WaitingForPlayers: "WaitingForPlayers",
		Starting:          "Starting",
	}
)

func init() {
	var v GamePhase
	if _, ok := interface{}(v).(fmt.Stringer); ok {
		_GamePhaseNameToValue = map[string]GamePhase{
			interface{}(WaitingForPlayers).(fmt.Stringer).String(): WaitingForPlayers,
			interface{}(Starting).(fmt.Stringer).String():          Starting,
		}
	}
}

// MarshalJSON is generated so GamePhase satisfies json.Marshaler.
func (r GamePhase) MarshalJSON() ([]byte, error) {
	if s, ok := interface{}(r).(fmt.Stringer); ok {
		return json.Marshal(s.String())
	}
	s, ok := _GamePhaseValueToName[r]
	if !ok {
		return nil, fmt.Errorf("invalid GamePhase: %d", r)
	}
	return json.Marshal(s)
}

// UnmarshalJSON is generated so GamePhase satisfies json.Unmarshaler.
func (r *GamePhase) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("GamePhase should be a string, got %s", data)
	}
	v, ok := _GamePhaseNameToValue[s]
	if !ok {
		return fmt.Errorf("invalid GamePhase %q", s)
	}
	*r = v
	return nil
}
