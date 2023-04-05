package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"unicode/utf16"
)

const (
	VolumeIDAndLocalBasePath               uint32 = 0x00000001
	CommonNetworkRelativeLink              uint32 = 0x00000002
	HasPathSuffix                          uint32 = 0x00000004
	CommonNetworkRelativeLinkAndPathSuffix uint32 = 0x00000002
	NameString                             uint32 = 1
	RelativePathString                     uint32 = 2
	WorkingDirString                       uint32 = 3
	CommandLineArgumentsString             uint32 = 4
	IconLocationString                     uint32 = 5
	HasLinkTargetIDList                    uint32 = 0x00000001
	HasLinkInfo                            uint32 = 0x00000002
	HasName                                uint32 = 0x00000004
	HasRelativePath                        uint32 = 0x00000008
	HasWorkingDir                          uint32 = 0x00000010
	HasArguments                           uint32 = 0x00000020
	HasIconLocation                        uint32 = 0x00000040
	IsUnicode                              uint32 = 0x00000080
	ForceNoLinkInfo                        uint32 = 0x00000100
	HasExpString                           uint32 = 0x00000200
	RunInSeparateProcess                   uint32 = 0x00000400
	Unused1                                uint32 = 0x00000800
	HasDarwinID                            uint32 = 0x00001000
	RunAsUser                              uint32 = 0x00002000
	HasExpIcon                             uint32 = 0x00004000
	NoPidlAlias                            uint32 = 0x00008000
	Unused2                                uint32 = 0x00010000
	RunWithShimLayer                       uint32 = 0x00020000
	ForceNoLinkTrack                       uint32 = 0x00040000
	EnableTargetMetadata                   uint32 = 0x00080000
	DisableLinkPathTracking                uint32 = 0x00100000
	DisableKnownFolderTracking             uint32 = 0x00200000
	DisableKnownFolderAlias                uint32 = 0x00400000
	AllowLinkToLink                        uint32 = 0x00800000
	UnaliasOnSave                          uint32 = 0x01000000
	PreferEnvironmentPath                  uint32 = 0x02000000
	KeepLocalIDListForUNCTarget            uint32 = 0x04000000
	Reserved                               uint32 = 0x08000000
	HTMLNoSubDirCreation                   uint32 = 0x10000000
	DisallowUserView                       uint32 = 0x20000000
	ForcePerceivedTypeSystem               uint32 = 0x40000000
	IncludeSlowInfo                        uint32 = 0x80000000
	ReservedForNYI                         uint32 = 0x00000000
	DisableShadowCopy                      uint32 = 0x00000000
	DisableKnownFolderAliasMigration       uint32 = 0x00000000
	DisableShellFolderVirtualization       uint32 = 0x00000000
)

type ShellLink struct {
	Header           ShellLinkHeader
	LinkTargetIDList *LinkTargetIDList
	LinkInfo         *LinkInfo
	StringData       StringData
	ExtraData        []*ExtraDataBlock
}

type LinkFlags struct {
	HasLinkTargetIDList              bool
	HasLinkInfo                      bool
	HasName                          bool
	HasRelativePath                  bool
	HasWorkingDir                    bool
	HasArguments                     bool
	HasIconLocation                  bool
	IsUnicode                        bool
	ForceNoLinkInfo                  bool
	HasExpString                     bool
	RunInSeparateProcess             bool
	Unused1                          bool
	HasDarwinID                      bool
	RunAsUser                        bool
	HasExpIcon                       bool
	NoPidlAlias                      bool
	Unused2                          bool
	RunWithShimLayer                 bool
	ForceNoLinkTrack                 bool
	EnableTargetMetadata             bool
	DisableLinkPathTracking          bool
	DisableKnownFolderTracking       bool
	DisableKnownFolderAlias          bool
	AllowLinkToLink                  bool
	UnaliasOnSave                    bool
	PreferEnvironmentPath            bool
	KeepLocalIDListForUNCTarget      bool
	HTMLNoSubDirCreation             bool
	DisallowUserView                 bool
	ForcePerceivedTypeSystem         bool
	IncludeSlowInfo                  bool
	ReservedForNYI                   bool
	DisableShadowCopy                bool
	DisableKnownFolderAliasMigration bool
	DisableShellFolderVirtualization bool
}

type FileAttributes struct {
	ReadOnly          bool
	Hidden            bool
	System            bool
	Archive           bool
	NTFS_EA           bool
	Temporary         bool
	SparseFile        bool
	ReparsePoint      bool
	Compressed        bool
	Offline           bool
	NotContentIndexed bool
	Encrypted         bool
}

type UUID [16]byte

type ShellLinkHeader struct {
	HeaderSize        uint32
	LinkCLSID         UUID
	LinkFlagsRaw      uint32
	FileAttributesRaw uint32
	CreationTime      uint64
	AccessTime        uint64
	WriteTime         uint64
	FileSize          uint32
	IconIndex         int32
	ShowCommand       uint32
	HotKey            uint16
	Reserved1         uint16
	Reserved2         uint32
	Reserved3         uint32
	LinkFlags         LinkFlags
	FileAttributes    FileAttributes
}

type LinkTargetIDList struct {
	IDListSize uint16
	IDList     []*ItemID
}

type LinkInfo struct {
	LinkInfoSize                    uint32
	LinkInfoFlagsRaw                uint32
	VolumeIDOffset                  uint32
	LocalBasePathOffset             uint32
	CommonNetworkRelativeLinkOffset uint32
	CommonPathSuffixOffset          uint32
	VolumeID                        []byte
	LocalBasePath                   string
	CommonNetworkRelativeLink       []byte
	CommonPathSuffix                string
}

type StringData struct {
	Name                 *string
	RelativePath         *string
	WorkingDir           *string
	CommandLineArguments *string
	IconLocation         *string
}

type StringDataItem struct {
	Name   uint32
	Size   uint16
	String string
}

type ItemID struct {
	ItemIDSize uint16
	Data       []byte
}

func parseIDList(idListData []byte) ([]*ItemID, error) {
	reader := bytes.NewReader(idListData)
	itemIDs := []*ItemID{}

	for {
		var itemIDSize uint16
		err := binary.Read(reader, binary.LittleEndian, &itemIDSize)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// Check for null-terminator
		if itemIDSize == 0 {
			break
		}

		dataSize := itemIDSize - 2 // Subtract the size of the ItemIDSize field itself
		data := make([]byte, dataSize)
		_, err = reader.Read(data)
		if err != nil {
			return nil, err
		}

		itemIDs = append(itemIDs, &ItemID{
			ItemIDSize: itemIDSize,
			Data:       data,
		})
	}

	return itemIDs, nil
}

type ExtraDataBlock struct {
	BlockSize      uint32
	BlockSignature uint32
	Data           []byte
}

func parseExtraData(reader *bytes.Reader) ([]*ExtraDataBlock, error) {
	extraDataBlocks := []*ExtraDataBlock{}

	for {
		var blockSize uint32
		err := binary.Read(reader, binary.LittleEndian, &blockSize)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// Check for null-terminator
		if blockSize == 0 {
			break
		}

		var blockSignature uint32
		err = binary.Read(reader, binary.LittleEndian, &blockSignature)
		if err != nil {
			return nil, err
		}

		dataSize := blockSize - 8 // Subtract the size of the BlockSize and BlockSignature fields
		data := make([]byte, dataSize)
		_, err = reader.Read(data)
		if err != nil {
			return nil, err
		}

		extraDataBlocks = append(extraDataBlocks, &ExtraDataBlock{
			BlockSize:      blockSize,
			BlockSignature: blockSignature,
			Data:           data,
		})
	}

	return extraDataBlocks, nil
}

func parseAllStringData(reader *bytes.Reader, isUnicode bool) (*StringData, error) {
	stringData := &StringData{}
	var err error

	for {
		item, err := parseStringData(reader, isUnicode)
		if err != nil {
			break
		}

		switch item.Name {
		case NameString:
			stringData.Name = &item.String
		case RelativePathString:
			stringData.RelativePath = &item.String
		case WorkingDirString:
			stringData.WorkingDir = &item.String
		case CommandLineArgumentsString:
			stringData.CommandLineArguments = &item.String
		case IconLocationString:
			stringData.IconLocation = &item.String
		}
	}

	if err != io.EOF {
		return nil, err
	}

	return stringData, nil
}

func parseShellLink(data []byte) (*ShellLink, error) {
	reader := bytes.NewReader(data)
	shellLink := &ShellLink{}

	// Parse ShellLinkHeader
	header, err := parseHeader(reader)
	if err != nil {
		return nil, err
	}
	shellLink.Header = *header

	// Parse LinkTargetIDList
	if header.LinkFlags.HasLinkTargetIDList {
		linkTargetIDList, err := parseLinkTargetIDList(reader)
		if err != nil {
			return nil, err
		}
		shellLink.LinkTargetIDList = linkTargetIDList
	}

	// Parse LinkInfo
	if header.LinkFlags.HasLinkInfo {
		linkInfo, err := parseLinkInfo(reader)
		if err != nil {
			return nil, fmt.Errorf("failed to parse link info: %v", err)
		}

		shellLink.LinkInfo = linkInfo
	}

	// Parse StringData
	if header.LinkFlags.HasName || header.LinkFlags.HasRelativePath || header.LinkFlags.HasWorkingDir || header.LinkFlags.HasArguments || header.LinkFlags.HasIconLocation {
		// Parse all StringData if present
		isUnicode := (header.LinkFlagsRaw & IsUnicode) != 0
		stringData, err := parseAllStringData(reader, isUnicode)
		if err != nil {
			return nil, fmt.Errorf("failed to parse string data: %v", err)
		}
		if stringData != nil {
			shellLink.StringData = *stringData
		}

	}

	// Parse ExtraData
	extraDataBlocks, err := parseExtraData(reader)
	if err != nil {
		return nil, err
	}
	shellLink.ExtraData = extraDataBlocks

	return shellLink, nil
}

func parseLinkFlags(flagsRaw uint32) LinkFlags {
	return LinkFlags{
		HasLinkTargetIDList:              (flagsRaw & HasLinkTargetIDList) != 0,
		HasLinkInfo:                      (flagsRaw & HasLinkInfo) != 0,
		HasName:                          (flagsRaw & HasName) != 0,
		HasRelativePath:                  (flagsRaw & HasRelativePath) != 0,
		HasWorkingDir:                    (flagsRaw & HasWorkingDir) != 0,
		HasArguments:                     (flagsRaw & HasArguments) != 0,
		HasIconLocation:                  (flagsRaw & HasIconLocation) != 0,
		IsUnicode:                        (flagsRaw & IsUnicode) != 0,
		ForceNoLinkInfo:                  (flagsRaw & ForceNoLinkInfo) != 0,
		HasExpString:                     (flagsRaw & HasExpString) != 0,
		RunInSeparateProcess:             (flagsRaw & RunInSeparateProcess) != 0,
		Unused1:                          (flagsRaw & Unused1) != 0,
		HasDarwinID:                      (flagsRaw & HasDarwinID) != 0,
		RunAsUser:                        (flagsRaw & RunAsUser) != 0,
		HasExpIcon:                       (flagsRaw & HasExpIcon) != 0,
		NoPidlAlias:                      (flagsRaw & NoPidlAlias) != 0,
		Unused2:                          (flagsRaw & Unused2) != 0,
		RunWithShimLayer:                 (flagsRaw & RunWithShimLayer) != 0,
		ForceNoLinkTrack:                 (flagsRaw & ForceNoLinkTrack) != 0,
		EnableTargetMetadata:             (flagsRaw & EnableTargetMetadata) != 0,
		DisableLinkPathTracking:          (flagsRaw & DisableLinkPathTracking) != 0,
		DisableKnownFolderTracking:       (flagsRaw & DisableKnownFolderTracking) != 0,
		DisableKnownFolderAlias:          (flagsRaw & DisableKnownFolderAlias) != 0,
		AllowLinkToLink:                  (flagsRaw & AllowLinkToLink) != 0,
		UnaliasOnSave:                    (flagsRaw & UnaliasOnSave) != 0,
		PreferEnvironmentPath:            (flagsRaw & PreferEnvironmentPath) != 0,
		KeepLocalIDListForUNCTarget:      (flagsRaw & KeepLocalIDListForUNCTarget) != 0,
		HTMLNoSubDirCreation:             (flagsRaw & HTMLNoSubDirCreation) != 0,
		DisallowUserView:                 (flagsRaw & DisallowUserView) != 0,
		ForcePerceivedTypeSystem:         (flagsRaw & ForcePerceivedTypeSystem) != 0,
		IncludeSlowInfo:                  (flagsRaw & IncludeSlowInfo) != 0,
		ReservedForNYI:                   (flagsRaw & ReservedForNYI) != 0,
		DisableShadowCopy:                (flagsRaw & DisableShadowCopy) != 0,
		DisableKnownFolderAliasMigration: (flagsRaw & DisableKnownFolderAliasMigration) != 0,
		DisableShellFolderVirtualization: (flagsRaw & DisableShellFolderVirtualization) != 0,
		// Reserved:                         (flagsRaw & Reserved) != 0,
	}
}

func parseHeader(reader *bytes.Reader) (*ShellLinkHeader, error) {
	header := &ShellLinkHeader{}
	err := binary.Read(reader, binary.LittleEndian, &header.HeaderSize)

	if err != nil {
		return nil, err
	}

	// Read LinkCLSID
	if _, err := reader.Read(header.LinkCLSID[:]); err != nil {
		return nil, err
	}

	err = binary.Read(reader, binary.LittleEndian, &header.LinkFlagsRaw)
	if err != nil {
		return nil, err
	}
	header.LinkFlags = parseLinkFlags(header.LinkFlagsRaw)

	err = binary.Read(reader, binary.LittleEndian, &header.FileAttributesRaw)
	if err != nil {
		return nil, err
	}

	err = binary.Read(reader, binary.LittleEndian, &header.CreationTime)
	if err != nil {
		return nil, err
	}

	err = binary.Read(reader, binary.LittleEndian, &header.AccessTime)
	if err != nil {
		return nil, err
	}

	err = binary.Read(reader, binary.LittleEndian, &header.WriteTime)
	if err != nil {
		return nil, err
	}

	err = binary.Read(reader, binary.LittleEndian, &header.FileSize)
	if err != nil {
		return nil, err
	}

	err = binary.Read(reader, binary.LittleEndian, &header.IconIndex)
	if err != nil {
		return nil, err
	}

	err = binary.Read(reader, binary.LittleEndian, &header.ShowCommand)
	if err != nil {
		return nil, err
	}

	err = binary.Read(reader, binary.LittleEndian, &header.HotKey)
	if err != nil {
		return nil, err
	}

	err = binary.Read(reader, binary.LittleEndian, &header.Reserved1)
	if err != nil {
		return nil, err
	}

	err = binary.Read(reader, binary.LittleEndian, &header.Reserved2)
	if err != nil {
		return nil, err
	}

	err = binary.Read(reader, binary.LittleEndian, &header.Reserved3)
	if err != nil {
		return nil, err
	}

	return header, nil
}

func parseLinkTargetIDList(reader *bytes.Reader) (*LinkTargetIDList, error) {
	var idListSize uint16
	err := binary.Read(reader, binary.LittleEndian, &idListSize)
	if err != nil {
		return nil, err
	}

	idListData := make([]byte, idListSize)
	_, err = reader.Read(idListData)
	if err != nil {
		return nil, err
	}

	itemIDs, err := parseIDList(idListData)
	if err != nil {
		return nil, err
	}

	return &LinkTargetIDList{
		IDListSize: idListSize,
		IDList:     itemIDs,
	}, nil
}

func readNullTerminatedString(reader io.Reader) (string, error) {
	var buf bytes.Buffer
	for {
		var b byte
		if err := binary.Read(reader, binary.LittleEndian, &b); err != nil {
			return "", err
		}
		if b == 0 {
			break
		}
		buf.WriteByte(b)
	}
	return buf.String(), nil
}

func parseLinkInfo(reader *bytes.Reader) (*LinkInfo, error) {
	var linkInfoSize uint32
	if err := binary.Read(reader, binary.LittleEndian, &linkInfoSize); err != nil {
		return nil, err
	}

	var linkInfoFlagsRaw uint32
	if err := binary.Read(reader, binary.LittleEndian, &linkInfoFlagsRaw); err != nil {
		return nil, err
	}

	linkInfo := &LinkInfo{
		LinkInfoSize:     linkInfoSize,
		LinkInfoFlagsRaw: linkInfoFlagsRaw,
	}

	if err := binary.Read(reader, binary.LittleEndian, &linkInfo.VolumeIDOffset); err != nil {
		return nil, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &linkInfo.LocalBasePathOffset); err != nil {
		return nil, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &linkInfo.CommonNetworkRelativeLinkOffset); err != nil {
		return nil, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &linkInfo.CommonPathSuffixOffset); err != nil {
		return nil, err
	}

	if linkInfo.LinkInfoFlagsRaw&VolumeIDAndLocalBasePath != 0 {
		var volumeIDSize uint32
		if _, err := reader.Seek(int64(linkInfo.VolumeIDOffset), io.SeekStart); err != nil {
			return nil, err
		}
		if err := binary.Read(reader, binary.LittleEndian, &volumeIDSize); err != nil {
			return nil, err
		}
		volumeIDData := make([]byte, volumeIDSize)
		if _, err := reader.Read(volumeIDData); err != nil {
			return nil, err
		}
		linkInfo.VolumeID = volumeIDData

		if _, err := reader.Seek(int64(linkInfo.LocalBasePathOffset), io.SeekStart); err != nil {
			return nil, err
		}
		localBasePath, err := readNullTerminatedString(reader)
		if err != nil {
			return nil, err
		}
		linkInfo.LocalBasePath = localBasePath
	}

	if linkInfo.LinkInfoFlagsRaw&CommonNetworkRelativeLinkAndPathSuffix != 0 {
		if _, err := reader.Seek(int64(linkInfo.CommonNetworkRelativeLinkOffset), io.SeekStart); err != nil {
			return nil, err
		}
		commonNetworkRelativeLinkSize := make([]byte, 4)
		if _, err := reader.Read(commonNetworkRelativeLinkSize); err != nil {
			return nil, err
		}
		linkInfo.CommonNetworkRelativeLink = commonNetworkRelativeLinkSize

		if _, err := reader.Seek(int64(linkInfo.CommonPathSuffixOffset), io.SeekStart); err != nil {
			return nil, err
		}
		commonPathSuffix, err := readNullTerminatedString(reader)
		if err != nil {
			return nil, err
		}
		linkInfo.CommonPathSuffix = commonPathSuffix
	}

	return linkInfo, nil
}

func parseStringData(reader *bytes.Reader, isUnicode bool) (*StringDataItem, error) {
	var strSize uint16
	err := binary.Read(reader, binary.LittleEndian, &strSize)
	if err != nil {
		return nil, err
	}

	if isUnicode {
		utf16Str := make([]uint16, strSize)
		err = binary.Read(reader, binary.LittleEndian, &utf16Str)
		if err != nil {
			return nil, err
		}
		return &StringDataItem{Size: strSize, String: string(utf16.Decode(utf16Str))}, nil
	}

	strData := make([]byte, strSize)
	_, err = reader.Read(strData)
	if err != nil {
		return nil, err
	}

	return &StringDataItem{Size: strSize, String: string(strData)}, nil
}

func formatGUID(guid string) string {
	// Assuming guid is a 16-byte binary string, convert it to the standard GUID format
	// For example, {12345678-1234-1234-1234-1234567890AB}
	b := []byte(guid)
	return fmt.Sprintf("{%08X-%04X-%04X-%04X-%02X%02X%02X%02X%02X%02X}", binary.LittleEndian.Uint32(b[0:4]), binary.LittleEndian.Uint16(b[4:6]), binary.LittleEndian.Uint16(b[6:8]), binary.LittleEndian.Uint16(b[8:10]), b[10], b[11], b[12], b[13], b[14], b[15])
}

func filetimeToUnix(ft int64) int64 {
	// Convert a Windows FILETIME (100-nanosecond intervals since January 1, 1601 UTC) to Unix time (seconds since January 1, 1970 UTC)
	return (ft - 116444736000000000) / 10000000
}

func utf16BytesToUTF8(u16s []uint16) string {
	return string(utf16.Decode(u16s))
}

func main() {
	//test file
	lnkPath := "Android Studio.lnk"

	// Read the contents of the file into a byte slice
	fileContent, err := ioutil.ReadFile(lnkPath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Parse the Shell Link file
	shellLink, err := parseShellLink(fileContent)
	if err != nil {
		log.Fatalf("Failed to parse Shell Link file: %v", err)
	}

	// Convert the Shell Link object to JSON
	shellLinkJSON, err := json.MarshalIndent(shellLink, "", "  ")
	if err != nil {
		log.Fatalf("Failed to convert Shell Link object to JSON: %v", err)
	}

	// Print the JSON representation of the Shell Link object
	fmt.Println(string(shellLinkJSON))
}
