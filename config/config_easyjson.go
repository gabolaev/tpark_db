// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package config

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

func easyjson6615c02eDecodeGithubComGabolaevTparkDbConfig(in *jlexer.Lexer, out *Database) {
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
		case "SchemaFile":
			out.SchemaFile = string(in.String())
		case "Timestamp":
			out.Timestamp = string(in.String())
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
func easyjson6615c02eEncodeGithubComGabolaevTparkDbConfig(out *jwriter.Writer, in Database) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"SchemaFile\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.SchemaFile))
	}
	{
		const prefix string = ",\"Timestamp\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Timestamp))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Database) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6615c02eEncodeGithubComGabolaevTparkDbConfig(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Database) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6615c02eEncodeGithubComGabolaevTparkDbConfig(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Database) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6615c02eDecodeGithubComGabolaevTparkDbConfig(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Database) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6615c02eDecodeGithubComGabolaevTparkDbConfig(l, v)
}
func easyjson6615c02eDecodeGithubComGabolaevTparkDbConfig1(in *jlexer.Lexer, out *Config) {
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
		case "Database":
			(out.Database).UnmarshalEasyJSON(in)
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
func easyjson6615c02eEncodeGithubComGabolaevTparkDbConfig1(out *jwriter.Writer, in Config) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Database\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(in.Database).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Config) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6615c02eEncodeGithubComGabolaevTparkDbConfig1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Config) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6615c02eEncodeGithubComGabolaevTparkDbConfig1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Config) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6615c02eDecodeGithubComGabolaevTparkDbConfig1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Config) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6615c02eDecodeGithubComGabolaevTparkDbConfig1(l, v)
}
