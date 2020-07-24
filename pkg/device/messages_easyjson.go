// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package device

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson66c1e240DecodeGitRonaksoftwareComBlipServerPkgDevice(in *jlexer.Lexer, out *RegisterDeviceReq) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "token_type":
			out.TokenType = string(in.String())
		case "token":
			out.Token = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson66c1e240EncodeGitRonaksoftwareComBlipServerPkgDevice(out *jwriter.Writer, in RegisterDeviceReq) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"token_type\":"
		out.RawString(prefix[1:])
		out.String(string(in.TokenType))
	}
	{
		const prefix string = ",\"token\":"
		out.RawString(prefix)
		out.String(string(in.Token))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v RegisterDeviceReq) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson66c1e240EncodeGitRonaksoftwareComBlipServerPkgDevice(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v RegisterDeviceReq) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson66c1e240EncodeGitRonaksoftwareComBlipServerPkgDevice(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *RegisterDeviceReq) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson66c1e240DecodeGitRonaksoftwareComBlipServerPkgDevice(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *RegisterDeviceReq) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson66c1e240DecodeGitRonaksoftwareComBlipServerPkgDevice(l, v)
}
