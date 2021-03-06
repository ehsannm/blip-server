// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package store

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

func easyjson66c1e240DecodeGitRonaksoftwareComBlipServerPkgStore(in *jlexer.Lexer, out *Stores) {
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
		case "stores":
			if in.IsNull() {
				in.Skip()
				out.Stores = nil
			} else {
				in.Delim('[')
				if out.Stores == nil {
					if !in.IsDelim(']') {
						out.Stores = make([]*Store, 0, 8)
					} else {
						out.Stores = []*Store{}
					}
				} else {
					out.Stores = (out.Stores)[:0]
				}
				for !in.IsDelim(']') {
					var v1 *Store
					if in.IsNull() {
						in.Skip()
						v1 = nil
					} else {
						if v1 == nil {
							v1 = new(Store)
						}
						(*v1).UnmarshalEasyJSON(in)
					}
					out.Stores = append(out.Stores, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
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
func easyjson66c1e240EncodeGitRonaksoftwareComBlipServerPkgStore(out *jwriter.Writer, in Stores) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"stores\":"
		out.RawString(prefix[1:])
		if in.Stores == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Stores {
				if v2 > 0 {
					out.RawByte(',')
				}
				if v3 == nil {
					out.RawString("null")
				} else {
					(*v3).MarshalEasyJSON(out)
				}
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Stores) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson66c1e240EncodeGitRonaksoftwareComBlipServerPkgStore(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Stores) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson66c1e240EncodeGitRonaksoftwareComBlipServerPkgStore(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Stores) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson66c1e240DecodeGitRonaksoftwareComBlipServerPkgStore(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Stores) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson66c1e240DecodeGitRonaksoftwareComBlipServerPkgStore(l, v)
}
func easyjson66c1e240DecodeGitRonaksoftwareComBlipServerPkgStore1(in *jlexer.Lexer, out *SaveStoreReq) {
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
		case "store_id":
			out.StoreID = int64(in.Int64())
		case "dsn":
			out.Dsn = string(in.String())
		case "region":
			out.Region = string(in.String())
		case "capacity":
			out.Capacity = int(in.Int())
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
func easyjson66c1e240EncodeGitRonaksoftwareComBlipServerPkgStore1(out *jwriter.Writer, in SaveStoreReq) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"store_id\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.StoreID))
	}
	{
		const prefix string = ",\"dsn\":"
		out.RawString(prefix)
		out.String(string(in.Dsn))
	}
	{
		const prefix string = ",\"region\":"
		out.RawString(prefix)
		out.String(string(in.Region))
	}
	{
		const prefix string = ",\"capacity\":"
		out.RawString(prefix)
		out.Int(int(in.Capacity))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v SaveStoreReq) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson66c1e240EncodeGitRonaksoftwareComBlipServerPkgStore1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v SaveStoreReq) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson66c1e240EncodeGitRonaksoftwareComBlipServerPkgStore1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *SaveStoreReq) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson66c1e240DecodeGitRonaksoftwareComBlipServerPkgStore1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *SaveStoreReq) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson66c1e240DecodeGitRonaksoftwareComBlipServerPkgStore1(l, v)
}
func easyjson66c1e240DecodeGitRonaksoftwareComBlipServerPkgStore2(in *jlexer.Lexer, out *GetStoreReq) {
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
		case "store_ids":
			if in.IsNull() {
				in.Skip()
				out.StoreIDs = nil
			} else {
				in.Delim('[')
				if out.StoreIDs == nil {
					if !in.IsDelim(']') {
						out.StoreIDs = make([]int64, 0, 8)
					} else {
						out.StoreIDs = []int64{}
					}
				} else {
					out.StoreIDs = (out.StoreIDs)[:0]
				}
				for !in.IsDelim(']') {
					var v4 int64
					v4 = int64(in.Int64())
					out.StoreIDs = append(out.StoreIDs, v4)
					in.WantComma()
				}
				in.Delim(']')
			}
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
func easyjson66c1e240EncodeGitRonaksoftwareComBlipServerPkgStore2(out *jwriter.Writer, in GetStoreReq) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"store_ids\":"
		out.RawString(prefix[1:])
		if in.StoreIDs == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v5, v6 := range in.StoreIDs {
				if v5 > 0 {
					out.RawByte(',')
				}
				out.Int64(int64(v6))
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v GetStoreReq) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson66c1e240EncodeGitRonaksoftwareComBlipServerPkgStore2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v GetStoreReq) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson66c1e240EncodeGitRonaksoftwareComBlipServerPkgStore2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *GetStoreReq) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson66c1e240DecodeGitRonaksoftwareComBlipServerPkgStore2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *GetStoreReq) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson66c1e240DecodeGitRonaksoftwareComBlipServerPkgStore2(l, v)
}
