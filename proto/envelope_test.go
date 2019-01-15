package proto

import (
	"reflect"
	"testing"
)

func TestNewMessage(t *testing.T) {
	type args struct {
		cmd     string
		to      string
		content []byte
	}
	tests := []struct {
		name string
		args args
		want Envelope
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEnvelope(tt.args.cmd, tt.args.content); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEnvelope() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessage_Serialize(t *testing.T) {
	type fields struct {
		id      []byte
		cmd     string
		content []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name: "serialize",
			fields: fields{
				id:      []byte("1111111111111111"),
				cmd:     "COMMAND",
				content: []byte("44"),
			},
			want: []byte("44"),
		},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewEnvelope(tt.fields.cmd, tt.fields.content)
			if got := m.Serialize(); !reflect.DeepEqual(got[len(got)-2:], tt.want) {
				t.Errorf("Envelope.Serialize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnSerialize(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		args args
		want *Envelope
	}{
		{
			name: "unserialize",
			args: args{b: NewEnvelope("2222", []byte("44")).Serialize()},
			want: NewEnvelope("2222", []byte("44")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UnSerialize(tt.args.b); !reflect.DeepEqual(got.Content, tt.want.Content) {
				t.Errorf("UnSerialize() = %v, want %v", got, tt.want)
			}
		})
	}
}
