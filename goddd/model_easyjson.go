// AUTOGENERATED FILE: easyjson marshaler/unmarshalers.

package goddd

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

func easyjsonC80ae7adDecodeGithubComHyperpilotioSnapPluginCollectorGodddGoddd(in *jlexer.Lexer, out *Summary) {
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
		case "sampleCount":
			out.SampleCount = uint64(in.Uint64())
		case "sampleSum":
			out.SampleSum = float64(in.Float64())
		case "quantile050":
			out.Quantile050 = float64(in.Float64())
		case "quantile090":
			out.Quantile090 = float64(in.Float64())
		case "quantile099":
			out.Quantile099 = float64(in.Float64())
		case "label":
			if in.IsNull() {
				in.Skip()
				out.Label = nil
			} else {
				in.Delim('[')
				if out.Label == nil {
					if !in.IsDelim(']') {
						out.Label = make([]*LabelStruct, 0, 8)
					} else {
						out.Label = []*LabelStruct{}
					}
				} else {
					out.Label = (out.Label)[:0]
				}
				for !in.IsDelim(']') {
					var v1 *LabelStruct
					if in.IsNull() {
						in.Skip()
						v1 = nil
					} else {
						if v1 == nil {
							v1 = new(LabelStruct)
						}
						(*v1).UnmarshalEasyJSON(in)
					}
					out.Label = append(out.Label, v1)
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
func easyjsonC80ae7adEncodeGithubComHyperpilotioSnapPluginCollectorGodddGoddd(out *jwriter.Writer, in Summary) {
	out.RawByte('{')
	first := true
	_ = first
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"sampleCount\":")
	out.Uint64(uint64(in.SampleCount))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"sampleSum\":")
	out.Float64(float64(in.SampleSum))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"quantile050\":")
	out.Float64(float64(in.Quantile050))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"quantile090\":")
	out.Float64(float64(in.Quantile090))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"quantile099\":")
	out.Float64(float64(in.Quantile099))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"label\":")
	if in.Label == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v2, v3 := range in.Label {
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
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Summary) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC80ae7adEncodeGithubComHyperpilotioSnapPluginCollectorGodddGoddd(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Summary) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC80ae7adEncodeGithubComHyperpilotioSnapPluginCollectorGodddGoddd(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Summary) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC80ae7adDecodeGithubComHyperpilotioSnapPluginCollectorGodddGoddd(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Summary) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC80ae7adDecodeGithubComHyperpilotioSnapPluginCollectorGodddGoddd(l, v)
}
func easyjsonC80ae7adDecodeGithubComHyperpilotioSnapPluginCollectorGodddGoddd1(in *jlexer.Lexer, out *LabelStruct) {
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
		case "name":
			out.Name = string(in.String())
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
func easyjsonC80ae7adEncodeGithubComHyperpilotioSnapPluginCollectorGodddGoddd1(out *jwriter.Writer, in LabelStruct) {
	out.RawByte('{')
	first := true
	_ = first
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"name\":")
	out.String(string(in.Name))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"value\":")
	out.String(string(in.Value))
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v LabelStruct) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC80ae7adEncodeGithubComHyperpilotioSnapPluginCollectorGodddGoddd1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v LabelStruct) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC80ae7adEncodeGithubComHyperpilotioSnapPluginCollectorGodddGoddd1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *LabelStruct) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC80ae7adDecodeGithubComHyperpilotioSnapPluginCollectorGodddGoddd1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *LabelStruct) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC80ae7adDecodeGithubComHyperpilotioSnapPluginCollectorGodddGoddd1(l, v)
}
func easyjsonC80ae7adDecodeGithubComHyperpilotioSnapPluginCollectorGodddGoddd2(in *jlexer.Lexer, out *CounterCache) {
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
		case "pre":
			out.Pre = float64(in.Float64())
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
func easyjsonC80ae7adEncodeGithubComHyperpilotioSnapPluginCollectorGodddGoddd2(out *jwriter.Writer, in CounterCache) {
	out.RawByte('{')
	first := true
	_ = first
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"pre\":")
	out.Float64(float64(in.Pre))
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v CounterCache) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC80ae7adEncodeGithubComHyperpilotioSnapPluginCollectorGodddGoddd2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v CounterCache) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC80ae7adEncodeGithubComHyperpilotioSnapPluginCollectorGodddGoddd2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *CounterCache) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC80ae7adDecodeGithubComHyperpilotioSnapPluginCollectorGodddGoddd2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *CounterCache) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC80ae7adDecodeGithubComHyperpilotioSnapPluginCollectorGodddGoddd2(l, v)
}
func easyjsonC80ae7adDecodeGithubComHyperpilotioSnapPluginCollectorGodddGoddd3(in *jlexer.Lexer, out *CacheType) {
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
		case "counterType":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				if !in.IsDelim('}') {
					out.CounterType = make(map[string]CounterCache)
				} else {
					out.CounterType = nil
				}
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v4 CounterCache
					(v4).UnmarshalEasyJSON(in)
					(out.CounterType)[key] = v4
					in.WantComma()
				}
				in.Delim('}')
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
func easyjsonC80ae7adEncodeGithubComHyperpilotioSnapPluginCollectorGodddGoddd3(out *jwriter.Writer, in CacheType) {
	out.RawByte('{')
	first := true
	_ = first
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"counterType\":")
	if in.CounterType == nil && (out.Flags&jwriter.NilMapAsEmpty) == 0 {
		out.RawString(`null`)
	} else {
		out.RawByte('{')
		v5First := true
		for v5Name, v5Value := range in.CounterType {
			if !v5First {
				out.RawByte(',')
			}
			v5First = false
			out.String(string(v5Name))
			out.RawByte(':')
			(v5Value).MarshalEasyJSON(out)
		}
		out.RawByte('}')
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v CacheType) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC80ae7adEncodeGithubComHyperpilotioSnapPluginCollectorGodddGoddd3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v CacheType) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC80ae7adEncodeGithubComHyperpilotioSnapPluginCollectorGodddGoddd3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *CacheType) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC80ae7adDecodeGithubComHyperpilotioSnapPluginCollectorGodddGoddd3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *CacheType) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC80ae7adDecodeGithubComHyperpilotioSnapPluginCollectorGodddGoddd3(l, v)
}
