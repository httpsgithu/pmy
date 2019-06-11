package pmy

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	shellBufferLeftVariableName  = "__pmy_out_buffer_left"
	shellBufferRightVariableName = "__pmy_out_buffer_right"
	shellCommandVariableName     = "__pmy_out_command"
	shellSourcesVariableName     = "__pmy_out_sources"
)

// pmyOut represents Output of pmy against zsh routine.
// This struct has strings exported to shell, whose embedded
// variables are all expanded.
type pmyOut struct {
	BufferLeft  string `json:"bufferLeft"`
	BufferRight string `json:"bufferRight"`
	Sources     string `json:"source"`
}

// newPmyOutFromRule create new pmyOut from rule
// which matches query and already has paramMap
func newPmyOutFromRule(rule *pmyRule) pmyOut {
	out := pmyOut{}
	out.Sources, _ = rule.CmdGroups.GetSources()
	out.BufferLeft = rule.BufferLeft
	out.BufferRight = rule.BufferRight
	out.expandAll(rule.paramMap)
	return out
}

// toShellVariables create zsh statement where pmyOut's attributes are
// passed into shell variables
func (out *pmyOut) toShellVariables() string {
	res := ""
	res += fmt.Sprintf("%v='%v';", shellBufferLeftVariableName, out.BufferLeft)
	res += fmt.Sprintf("%v='%v';", shellBufferRightVariableName, out.BufferRight)
	res += fmt.Sprintf("%v='%v';", shellSourcesVariableName, out.Sources)
	return res
}

func expand(org string, paramMap map[string]string) string {
	res := org
	for name, value := range paramMap {
		res = strings.Replace(
			res,
			fmt.Sprintf("<%v>", name),
			value,
			-1,
		)
	}
	return res
}

func (out *pmyOut) expandAll(paramMap map[string]string) {
	out.BufferLeft = expand(out.BufferLeft, paramMap)
	out.BufferRight = expand(out.BufferRight, paramMap)
	return
}

func (out *pmyOut) toJSON() string {
	bytes, _ := json.Marshal(out)
	str := string(bytes)
	return str
}

// func (out *pmyOut) serialize() string {
//     return out.BufferLeft + delimiter + out.BufferRight + delimiter + out.Command
// }
