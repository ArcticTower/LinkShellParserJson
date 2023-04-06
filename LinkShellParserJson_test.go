package main

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

func Test_parseIDList(t *testing.T) {
	type args struct {
		idListData []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []*ItemID
		wantErr bool
	}{
		{
			name: "Test Case 1",
			args: args{
				idListData: []byte{0x14, 0x00, 0x1F, 0x50, 0xE0, 0x4F, 0xD0, 0x20, 0xEA, 0x3A, 0x69, 0x10, 0xA2, 0xD8, 0x08, 0x00, 0x2B, 0x30, 0x30, 0x9D},
			},
			want: []*ItemID{
				{
					ItemIDSize: 0x14,
					Data:       []byte{0x1F, 0x50, 0xE0, 0x4F, 0xD0, 0x20, 0xEA, 0x3A, 0x69, 0x10, 0xA2, 0xD8, 0x08, 0x00, 0x2B, 0x30, 0x30, 0x9D},
				},
			},
			wantErr: false,
		},
		// Add more test cases here
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseIDList(tt.args.idListData)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseIDList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseIDList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseExtraData(t *testing.T) {
	type args struct {
		reader *bytes.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    []*ExtraDataBlock
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseExtraData(tt.args.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseExtraData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseExtraData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseAllStringData(t *testing.T) {
	type args struct {
		reader    *bytes.Reader
		isUnicode bool
	}
	tests := []struct {
		name    string
		args    args
		want    *StringData
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseAllStringData(tt.args.reader, tt.args.isUnicode)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseAllStringData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseAllStringData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseShellLink(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *ShellLink
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseShellLink(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseShellLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseShellLink() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseLinkFlags(t *testing.T) {
	type args struct {
		flagsRaw uint32
	}
	tests := []struct {
		name string
		args args
		want LinkFlags
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseLinkFlags(tt.args.flagsRaw); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseLinkFlags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseHeader(t *testing.T) {
	type args struct {
		reader *bytes.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    *ShellLinkHeader
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseHeader(tt.args.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseHeader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseLinkTargetIDList(t *testing.T) {
	type args struct {
		reader *bytes.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    *LinkTargetIDList
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseLinkTargetIDList(tt.args.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseLinkTargetIDList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseLinkTargetIDList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readNullTerminatedString(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readNullTerminatedString(tt.args.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("readNullTerminatedString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("readNullTerminatedString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseLinkInfo(t *testing.T) {
	type args struct {
		reader *bytes.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    *LinkInfo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseLinkInfo(tt.args.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseLinkInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseLinkInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseStringData(t *testing.T) {
	type args struct {
		reader    *bytes.Reader
		isUnicode bool
	}
	tests := []struct {
		name    string
		args    args
		want    *StringDataItem
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseStringData(tt.args.reader, tt.args.isUnicode)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseStringData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseStringData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_formatGUID(t *testing.T) {
	type args struct {
		guid string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatGUID(tt.args.guid); got != tt.want {
				t.Errorf("formatGUID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filetimeToUnix(t *testing.T) {
	type args struct {
		ft int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filetimeToUnix(tt.args.ft); got != tt.want {
				t.Errorf("filetimeToUnix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_utf16BytesToUTF8(t *testing.T) {
	type args struct {
		u16s []uint16
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utf16BytesToUTF8(tt.args.u16s); got != tt.want {
				t.Errorf("utf16BytesToUTF8() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}
