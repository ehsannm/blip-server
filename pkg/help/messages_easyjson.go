// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package help

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

func easyjson66c1e240DecodeGitRonaksoftwareComBlipServerPkgHelp(in *jlexer.Lexer, out *UnsetDefaultConfig) {
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
		case "key":
			out.Key = string(in.String())
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
func easyjson66c1e240EncodeGitRonaksoftwareComBlipServerPkgHelp(out *jwriter.Writer, in UnsetDefaultConfig) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"key\":"
		out.RawString(prefix[1:])
		out.String(string(in.Key))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v UnsetDefaultConfig) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson66c1e240EncodeGitRonaksoftwareComBlipServerPkgHelp(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UnsetDefaultConfig) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson66c1e240EncodeGitRonaksoftwareComBlipServerPkgHelp(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UnsetDefaultConfig) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson66c1e240DecodeGitRonaksoftwareComBlipServerPkgHelp(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UnsetDefaultConfig) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson66c1e240DecodeGitRonaksoftwareComBlipServerPkgHelp(l, v)
}
func easyjson66c1e240DecodeGitRonaksoftwareComBlipServerPkgHelp1(in *jlexer.Lexer, out *SetDefaultConfig) {
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
		case "key":
			out.Key = string(in.String())
		case "value":
			out.Value = string(in.String())
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
func easyjson66c1e240EncodeGitRonaksoftwareComBlipServerPkgHelp1(out *jwriter.Writer, in SetDefaultConfig) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"key\":"
		out.RawString(prefix[1:])
		out.String(string(in.Key))
	}
	{
		const prefix string = ",\"value\":"
		out.RawString(prefix)
		out.String(string(in.Value))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v SetDefaultConfig) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson66c1e240EncodeGitRonaksoftwareComBlipServerPkgHelp1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v SetDefaultConfig) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson66c1e240EncodeGitRonaksoftwareComBlipServerPkgHelp1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *SetDefaultConfig) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson66c1e240DecodeGitRonaksoftwareComBlipServerPkgHelp1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *SetDefaultConfig) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson66c1e240DecodeGitRonaksoftwareComBlipServerPkgHelp1(l, v)
}
func easyjson66c1e240DecodeGitRonaksoftwareComBlipServerPkgHelp2(in *jlexer.Lexer, out *Config) {
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
		case "update_available":
			out.UpdateAvailable = bool(in.Bool())
		case "update_force":
			out.UpdateForce = bool(in.Bool())
		case "store_link":
			out.StoreLink = string(in.String())
		case "show_blip_link":
			out.ShowBlipLink = bool(in.Bool())
		case "show_share_link":
			out.ShowShareLink = bool(in.Bool())
		case "authorized":
			out.Authorized = bool(in.Bool())
		case "vas_enabled":
			out.VasEnabled = bool(in.Bool())
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
func easyjson66c1e240EncodeGitRonaksoftwareComBlipServerPkgHelp2(out *jwriter.Writer, in Config) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"update_available\":"
		out.RawString(prefix[1:])
		out.Bool(bool(in.UpdateAvailable))
	}
	{
		const prefix string = ",\"update_force\":"
		out.RawString(prefix)
		out.Bool(bool(in.UpdateForce))
	}
	{
		const prefix string = ",\"store_link\":"
		out.RawString(prefix)
		out.String(string(in.StoreLink))
	}
	{
		const prefix string = ",\"show_blip_link\":"
		out.RawString(prefix)
		out.Bool(bool(in.ShowBlipLink))
	}
	{
		const prefix string = ",\"show_share_link\":"
		out.RawString(prefix)
		out.Bool(bool(in.ShowShareLink))
	}
	{
		const prefix string = ",\"authorized\":"
		out.RawString(prefix)
		out.Bool(bool(in.Authorized))
	}
	{
		const prefix string = ",\"vas_enabled\":"
		out.RawString(prefix)
		out.Bool(bool(in.VasEnabled))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Config) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson66c1e240EncodeGitRonaksoftwareComBlipServerPkgHelp2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Config) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson66c1e240EncodeGitRonaksoftwareComBlipServerPkgHelp2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Config) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson66c1e240DecodeGitRonaksoftwareComBlipServerPkgHelp2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Config) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson66c1e240DecodeGitRonaksoftwareComBlipServerPkgHelp2(l, v)
}
