package main

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"unicode/utf8"
)

// ShellLinkParsed
type ShellLinkParsed struct {
	header               ShellLinkHeader
	linkFlagsParsed      LinkFlagsParsed
	fileAttributesParsed FileAttributesParsed
	linkTargetIDList     LinkTargetIDList
	linkInfo             LinkInfo
}

// ShellLinkHeader represents the header of a .lnk file.
type ShellLinkHeader struct {
	HeaderSize     uint32
	LinkCLSID      [16]byte
	LinkFlags      uint32
	FileAttributes uint32
	CreationTime   uint64
	AccessTime     uint64
	WriteTime      uint64
	FileSize       uint32
	IconIndex      int32
	ShowCommand    uint32
	HotKey         uint16
	Reserved1      uint16
	Reserved2      uint32
	Reserved3      uint32
}

// LinkFlags
const (
	HasLinkTargetIDList         uint32 = 0x00000001
	HasLinkInfo                 uint32 = 0x00000002
	HasName                     uint32 = 0x00000004
	HasRelativePath             uint32 = 0x00000008
	HasWorkingDir               uint32 = 0x00000010
	HasArguments                uint32 = 0x00000020
	HasIconLocation             uint32 = 0x00000040
	IsUnicode                   uint32 = 0x00000080
	ForceNoLinkInfo             uint32 = 0x00000100
	HasExpString                uint32 = 0x00000200
	RunInSeparateProcess        uint32 = 0x00000400
	Unused1                     uint32 = 0x00000800
	HasDarwinID                 uint32 = 0x00001000
	RunAsUser                   uint32 = 0x00002000
	HasExpIcon                  uint32 = 0x00004000
	NoPidlAlias                 uint32 = 0x00008000
	Unused2                     uint32 = 0x00010000
	RunWithShimLayer            uint32 = 0x00020000
	ForceNoLinkTrack            uint32 = 0x00040000
	EnableTargetMetadata        uint32 = 0x00080000
	DisableLinkPathTracking     uint32 = 0x00100000
	DisableKnownFolderTracking  uint32 = 0x00200000
	DisableKnownFolderAlias     uint32 = 0x00400000
	AllowLinkToLink             uint32 = 0x00800000
	UnaliasOnSave               uint32 = 0x01000000
	PreferEnvironmentPath       uint32 = 0x02000000
	KeepLocalIDListForUNCTarget uint32 = 0x04000000
	//not shure:
	Reserved uint32 = 0x08000000
	//???
	HTMLNoSubDirCreation uint32 = 0x10000000
	//with HasExpString
	DisallowUserView uint32 = 0x20000000
	//???
	ForcePerceivedTypeSystem uint32 = 0x40000000
	//indicates that the perceived type of the link target should be used as the actual type when launching the target, regardless of the actual type of the file
	IncludeSlowInfo uint32 = 0x80000000
	//slow properties should be included in the link serialization
	//NYI?:
	// ReservedForNYI                   uint32 = 0x00000000
	// DisableShadowCopy                uint32 = 0x00000000
	// DisableKnownFolderAliasMigration uint32 = 0x00000000
	// DisableShellFolderVirtualization uint32 = 0x00000000
)

// LinkFlagsParsed
type LinkFlagsParsed struct {
	HasLinkTargetIDList         bool //uint32 = 0x00000001
	HasLinkInfo                 bool //uint32 = 0x00000002
	HasName                     bool //uint32 = 0x00000004
	HasRelativePath             bool //uint32 = 0x00000008
	HasWorkingDir               bool //uint32 = 0x00000010
	HasArguments                bool //uint32 = 0x00000020
	HasIconLocation             bool //uint32 = 0x00000040
	IsUnicode                   bool //uint32 = 0x00000080
	ForceNoLinkInfo             bool //uint32 = 0x00000100
	HasExpString                bool //uint32 = 0x00000200
	RunInSeparateProcess        bool //uint32 = 0x00000400
	Unused1                     bool //uint32 = 0x00000800
	HasDarwinID                 bool //uint32 = 0x00001000
	RunAsUser                   bool //uint32 = 0x00002000
	HasExpIcon                  bool //uint32 = 0x00004000
	NoPidlAlias                 bool //uint32 = 0x00008000
	Unused2                     bool //uint32 = 0x00010000
	RunWithShimLayer            bool //uint32 = 0x00020000
	ForceNoLinkTrack            bool //uint32 = 0x00040000
	EnableTargetMetadata        bool //uint32 = 0x00080000
	DisableLinkPathTracking     bool //uint32 = 0x00100000
	DisableKnownFolderTracking  bool //uint32 = 0x00200000
	DisableKnownFolderAlias     bool //uint32 = 0x00400000
	AllowLinkToLink             bool //uint32 = 0x00800000
	UnaliasOnSave               bool //uint32 = 0x01000000
	PreferEnvironmentPath       bool //uint32 = 0x02000000
	KeepLocalIDListForUNCTarget bool //uint32 = 0x04000000
	Reserved                    bool //uint32 = 0x08000000
	HTMLNoSubDirCreation        bool //uint32 = 0x10000000
	DisallowUserView            bool //uint32 = 0x20000000
	ForcePerceivedTypeSystem    bool //uint32 = 0x40000000
	IncludeSlowInfo             bool //uint32 = 0x80000000
}

// FileAttributes
const (
	FileAttributeReadOnly          uint32 = 0x00000001
	FileAttributeHidden            uint32 = 0x00000002
	FileAttributeSystem            uint32 = 0x00000004
	FileAttributeVolumeLabel       uint32 = 0x00000008 // Not used in link files
	FileAttributeDirectory         uint32 = 0x00000010
	FileAttributeArchive           uint32 = 0x00000020
	FileAttributeNormal            uint32 = 0x00000080
	FileAttributeTemporary         uint32 = 0x00000100
	FileAttributeSparseFile        uint32 = 0x00000200
	FileAttributeReparsePoint      uint32 = 0x00000400
	FileAttributeCompressed        uint32 = 0x00000800
	FileAttributeOffline           uint32 = 0x00001000
	FileAttributeNotContentIndexed uint32 = 0x00002000
	FileAttributeEncrypted         uint32 = 0x00004000
)

// FileAttributesParsed
type FileAttributesParsed struct {
	FileAttributeReadOnly          bool //uint32 = 0x00000001
	FileAttributeHidden            bool //uint32 = 0x00000002
	FileAttributeSystem            bool //uint32 = 0x00000004
	FileAttributeVolumeLabel       bool //uint32 = 0x00000008 // Not used in link files
	FileAttributeDirectory         bool //uint32 = 0x00000010
	FileAttributeArchive           bool //uint32 = 0x00000020
	FileAttributeNormal            bool //uint32 = 0x00000080
	FileAttributeTemporary         bool //uint32 = 0x00000100
	FileAttributeSparseFile        bool //uint32 = 0x00000200
	FileAttributeReparsePoint      bool //uint32 = 0x00000400
	FileAttributeCompressed        bool //uint32 = 0x00000800
	FileAttributeOffline           bool //uint32 = 0x00001000
	FileAttributeNotContentIndexed bool //uint32 = 0x00002000
	FileAttributeEncrypted         bool //uint32 = 0x00004000
}

// LinkTargetIDList represents the item ID list of a .lnk file.
type LinkTargetIDList struct {
	IDListSize uint16
	IDListData IDList
}

// IDList represents the IDList structure, which is a sequence of ItemID structures.
type IDList struct {
	ItemIDs []ItemID
	//ends with uint16 \0
}

// ItemID represents an ItemID structure in the IDList.
type ItemID struct {
	ItemIDSize       uint16
	ItemIDDataBase64 string
	ItemIDData       []byte
}

// LinkInfo represents the link information of a .lnk file.
type LinkInfo struct {
	LinkInfoSize                    uint32
	LinkInfoHeaderSize              uint32
	LinkInfoFlags                   uint32
	VolumeIDOffset                  uint32
	LocalBasePathOffset             uint32
	CommonNetworkRelativeLinkOffset uint32
	CommonPathSuffixOffset          uint32
	LocalBasePathOffsetUnicode      uint32
	CommonPathSuffixOffsetUnicode   uint32
	VolumeID                        VolumeID
	LocalBasePath                   string
	LocalBasePathBase64             string
	CommonNetworkRelativeLink       CommonNetworkRelativeLink
	CommonPathSuffix                string
	CommonPathSuffixBase64          string
	LocalBasePathUnicode            string // Optional, present if LinkFlag 'IsUnicode' is set
	CommonPathSuffixUnicode         string // Optional, present if LinkFlag 'IsUnicode' is set
}

// LinkInfoHeaderSize
const (
	LinkInfoHeaderSizeOptionalFieldsNotSpecified  uint32 = 0x0000001C
	LinkInfoHeaderSizeOptionalFieldsSpecifiedFrom uint32 = 0x00000024
)

// LinkInfoFlags
const (
	VolumeIDAndLocalBasePathPresent               uint32 = 0x00000001
	CommonNetworkRelativeLinkAndPathSuffixPresent uint32 = 0x00000002
)

type VolumeID struct {
	VolumeIDSize             uint32
	DriveType                uint32
	DriveSerialNumber        uint32
	VolumeLabelOffset        uint32
	VolumeLabelOffsetUnicode uint32 // Optional, present if LinkFlag 'IsUnicode' is set
	VolumeLabel              string
	VolumeLabelUnicode       string // Optional, present if LinkFlag 'IsUnicode' is set
	VolumeLableBase64        string
}

// VolumeLableOffset
const (
	VolumeLabelOffsetUnicodePresent uint32 = 0x00000014
	VolumeIDSizeMin                 uint32 = 0x00000010
)

// DriveType
const (
	DriveUnknown   uint32 = 0x00000000
	DriveNoRootDir uint32 = 0x00000001
	DriveRemovable uint32 = 0x00000002
	DriveFixed     uint32 = 0x00000003
	DriveRemote    uint32 = 0x00000004
	DriveCDROM     uint32 = 0x00000005
	DriveRAMDisk   uint32 = 0x00000006
)

// CommonNetworkRelativeLink represents the CommonNetworkRelativeLink structure in the LinkInfo.
type CommonNetworkRelativeLink struct {
	CommonNetworkRelativeLinkSize  uint32
	CommonNetworkRelativeLinkFlags uint32
	NetNameOffset                  uint32
	DeviceNameOffset               uint32
	NetworkProviderType            uint32
	NetNameOffsetUnicode           uint32 // Optional, present if LinkFlag 'IsUnicode' is set
	DeviceNameOffsetUnicode        uint32 // Optional, present if LinkFlag 'IsUnicode' is set
	NetName                        string
	NetNameBase64                  string
	DeviceName                     string
	DeviceNameBase64               string
	NetNameUnicode                 string // Optional, present if LinkFlag 'IsUnicode' is set
	DeviceNameUnicode              string // Optional, present if LinkFlag 'IsUnicode' is set
}

// CommonNetworkRelativeLinkFlags
const (
	ValidDevice                                uint32 = 0x00000001
	ValidNetType                               uint32 = 0x00000002
	CommonNetworkRelativeLinkUnicodeMinOffsets uint32 = 0x00000014
)

// NetworkProviderType
const (
	WNNC_NET_AVID        uint32 = 0x001A0000
	WNNC_NET_DOCUSPACE   uint32 = 0x001B0000
	WNNC_NET_MANGOSOFT   uint32 = 0x001C0000
	WNNC_NET_SERNET      uint32 = 0x001D0000
	WNNC_NET_RIVERFRONT1 uint32 = 0x001E0000
	WNNC_NET_RIVERFRONT2 uint32 = 0x001F0000
	WNNC_NET_DECORB      uint32 = 0x00200000
	WNNC_NET_PROTSTOR    uint32 = 0x00210000
	WNNC_NET_FJ_REDIR    uint32 = 0x00220000
	WNNC_NET_DISTINCT    uint32 = 0x00230000
	WNNC_NET_TWINS       uint32 = 0x00240000
	WNNC_NET_RDR2SAMPLE  uint32 = 0x00250000
	WNNC_NET_CSC         uint32 = 0x00260000
	WNNC_NET_3IN1        uint32 = 0x00270000
	WNNC_NET_EXTENDNET   uint32 = 0x00290000
	WNNC_NET_STAC        uint32 = 0x002A0000
	WNNC_NET_FOXBAT      uint32 = 0x002B0000
	WNNC_NET_YAHOO       uint32 = 0x002C0000
	WNNC_NET_EXIFS       uint32 = 0x002D0000
	WNNC_NET_DAV         uint32 = 0x002E0000
	WNNC_NET_KNOWARE     uint32 = 0x002F0000
	WNNC_NET_OBJECT_DIRE uint32 = 0x00300000
	WNNC_NET_MASFAX      uint32 = 0x00310000
	WNNC_NET_HOB_NFS     uint32 = 0x00320000
	WNNC_NET_SHIVA       uint32 = 0x00330000
	WNNC_NET_IBMAL       uint32 = 0x00340000
	WNNC_NET_LOCK        uint32 = 0x00350000
	WNNC_NET_TERMSRV     uint32 = 0x00360000
	WNNC_NET_SRT         uint32 = 0x00370000
	WNNC_NET_QUINCY      uint32 = 0x00380000
	WNNC_NET_OPENAFS     uint32 = 0x00390000
	WNNC_NET_AVID1       uint32 = 0x003A0000
	WNNC_NET_DFS         uint32 = 0x003B0000
	WNNC_NET_KWNP        uint32 = 0x003C0000
	WNNC_NET_ZENWORKS    uint32 = 0x003D0000
	WNNC_NET_DRIVEONWEB  uint32 = 0x003E0000
	WNNC_NET_VMWARE      uint32 = 0x003F0000
	WNNC_NET_RSFX        uint32 = 0x00400000
	WNNC_NET_MFILES      uint32 = 0x00410000
	WNNC_NET_MS_NFS      uint32 = 0x00420000
	WNNC_NET_GOOGLE      uint32 = 0x00430000
)

// StringData represents the string data section of a .lnk file.
type StringData struct {
	NameString            string
	RelativePath          string
	WorkingDir            string
	CommandLineArgs       string
	IconLocation          string
	NameStringBase64      string
	RelativePathBase64    string
	WorkingDirBase64      string
	CommandLineArgsBase64 string
	IconLocationBase64    string
}

// ExtraData represents the extra data section of a .lnk file.
type ExtraData struct {
	BlockSignature uint32
	BlockSize      uint32
	BlockData      []byte
}

// BlockSignature
const (
	ConsoleDataBlockSignature             uint32 = 0xA0000002
	ConsoleFEDataBlockSignature           uint32 = 0xA0000004
	DarwinDataBlockSignature              uint32 = 0xA0000006
	EnviromentVariableDataBlockSignature  uint32 = 0xA0000001
	IconEnviromentDataBlockSignature      uint32 = 0xA0000007
	KnownFolderDataBlockSignature         uint32 = 0xA000000B
	PropertyStoreDataBlockSignature       uint32 = 0xA0000009
	ShimDataBlockSignature                uint32 = 0xA0000008
	SpecialFolderDataBlockSignature       uint32 = 0xA0000005
	TrackerDataBlockSignature             uint32 = 0xA0000003
	VistaAndAboveIDListDataBlockSignature uint32 = 0xA000000C
)

// ConsoleDataBlock represents the ConsoleDataBlock structure in the ExtraData section.
type ConsoleDataBlock struct {
	BlockSignature         uint32
	BlockSize              uint32
	FillAttributes         uint16
	PopupFillAttributes    uint16
	ScreenBufferSizeX      uint16
	ScreenBufferSizeY      uint16
	WindowSizeX            uint16
	WindowSizeY            uint16
	WindowOriginX          uint16
	WindowOriginY          uint16
	Unused1                uint32
	Unused2                uint32
	FontSize               uint32
	FontFamily             uint32
	FontWeight             uint32
	FaceName               [32]byte
	CursorSize             uint32
	FullSreen              uint32
	QuickEdit              uint32
	InsertMode             uint32
	AutoPosition           uint32
	HistoryBufferSize      uint32
	NumberOfHistoryBuffers uint32
	HistoryNoDup           uint32
	ColorTable             [16]uint32
}

// FillAttributes
const (
	FILLATTR_FOREGROUND_BLUE      uint16 = 0x0001
	FILLATTR_FOREGROUND_GREEN     uint16 = 0x0002
	FILLATTR_FOREGROUND_RED       uint16 = 0x0004
	FILLATTR_FOREGROUND_INTENSITY uint16 = 0x0008
	FILLATTR_BACKGROUND_BLUE      uint16 = 0x0010
	FILLATTR_BACKGROUND_GREEN     uint16 = 0x0020
	FILLATTR_BACKGROUND_RED       uint16 = 0x0040
	FILLATTR_BACKGROUND_INTENSITY uint16 = 0x0080
)

// FontFamily
const (
	FONT_FAMILY_DONTCARE   uint32 = 0x0000
	FONT_FAMILY_ROMAN      uint32 = 0x0010
	FONT_FAMILY_SWISS      uint32 = 0x0020
	FONT_FAMILY_MODERN     uint32 = 0x0030
	FONT_FAMILY_SCRIPT     uint32 = 0x0040
	FONT_FAMILY_DECORATIVE uint32 = 0x0050
)

// FontWeight
const (
	BoldMin uint32 = 700 // val > BoldMin -> bold,
)

// CursorSize
const (
	SmallMax  uint32 = 25
	MediumMax uint32 = 50
	LargeMax  uint32 = 100
)

// FullScreen
const (
	FullScreenOff uint32 = 0x00000000 // > FullScreenOff -> on
)

// QuickEdit
const (
	QuickEditOff uint32 = 0x00000000
)

// InsertMode
const (
	InsertModeDisabled uint32 = 0x00000000
)

// AutoPosition
const (
	AutoPosition uint32 = 0x00000000 //WindowsOrigin used
)

// HistoryNoDup
const (
	HistoryNoDupAllowed uint32 = 0x00000000
)

// ConsoleFEDataBlock represents the ConsoleFEDataBlock structure in the ExtraData section.
type ConsoleFEDataBlock struct {
	BlockSignature uint32
	BlockSize      uint32
	CodePage       uint32
}

// DarwinDataBlock represents the DarwinDataBlock structure in the ExtraData section.
type DarwinDataBlock struct {
	BlockSignature    uint32
	BlockSize         uint32
	DarwinDataAnsi    string
	DarwinDataUnicode string
}

// EnvironmentVariableDataBlock represents the EnvironmentVariableDataBlock structure in the ExtraData section.
type EnvironmentVariableDataBlock struct {
	BlockSignature uint32
	BlockSize      uint32
	TargetAnsi     string
	TargetUnicode  string
}

// IconEnvironmentDataBlock represents the IconEnvironmentDataBlock structure in the ExtraData section.
type IconEnvironmentDataBlock struct {
	BlockSignature uint32
	BlockSize      uint32
	TargetAnsi     string
	TargetUnicode  string
}

// KnownFolderDataBlock represents the KnownFolderDataBlock structure in the ExtraData section.
type KnownFolderDataBlock struct {
	BlockSignature uint32
	BlockSize      uint32
	KnownFolderID  [16]byte
	Offset         int32
}

// PropertyStoreDataBlock represents the PropertyStoreDataBlock structure in the ExtraData section.
type PropertyStoreDataBlock struct {
	BlockSignature uint32
	BlockSize      uint32
	PropertyStore  []byte
}

// ShimDataBlock represents the ShimDataBlock structure in the ExtraData section.
type ShimDataBlock struct {
	BlockSignature uint32
	BlockSize      uint32
	LayerName      string
}

// SpecialFolderDataBlock represents the SpecialFolderDataBlock structure in the ExtraData section.
type SpecialFolderDataBlock struct {
	BlockSignature  uint32
	BlockSize       uint32
	SpecialFolderID uint32
	Offset          int32
}

// TrackerDataBlock represents the TrackerDataBlock structure in the ExtraData section.
type TrackerDataBlock struct {
	BlockSignature   uint32
	BlockSize        uint32
	Length           uint32
	Version          uint32
	MachineID        string
	DroidVolume      [16]byte
	DroidFile        [16]byte
	BirthDroidVolume [16]byte
	BirthDroidFile   [16]byte
}

// VistaAndAboveIDListDataBlock represents the VistaAndAboveIDListDataBlock structure in the ExtraData section.
type VistaAndAboveIDListDataBlock struct {
	BlockSignature uint32
	BlockSize      uint32
	IDList         []byte
}

//TODO: ExtraDataBlock sizes as consts

// Expected const list
const (
	HeaderSizeExpected uint32 = 0x0000004C
)

var LinkCLSIDExpected = [16]byte{0x01, 0x14, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46}

type ConstMismatchError struct {
	At       string
	Is       string
	Expected string
}

func (e *ConstMismatchError) Error() string {
	return fmt.Sprintf("parse at %v: const mismatch: is: %v, expected: %v", e.At, e.Is, e.Expected)
}

// ReadLnkFile reads the content of a .lnk file and returns the data as a byte slice.
func ReadLnkFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}

	return data, nil
}

// Parse ShellLinkHeader
func ParseShellLinkHeader(r *bytes.Reader) (ShellLinkHeader, error) {
	var head ShellLinkHeader

	err := binary.Read(r, binary.LittleEndian, &head)

	if err != nil {
		return head, err
	}

	if head.HeaderSize != HeaderSizeExpected {
		return head, &ConstMismatchError{
			At:       "ShellLinkHeader",
			Is:       strconv.FormatUint(uint64(head.HeaderSize), 16),
			Expected: strconv.FormatUint(uint64(HeaderSizeExpected), 16),
		}
	}

	if head.LinkCLSID != LinkCLSIDExpected {
		is := ""
		exp := ""
		for _, b := range head.LinkCLSID {
			is += fmt.Sprintf("0x%02x ", b)
		}
		for _, b := range LinkCLSIDExpected {
			exp += fmt.Sprintf("0x%02x ", b)
		}
		return head, &ConstMismatchError{
			At:       "ShellLinkHeader",
			Is:       is,
			Expected: exp,
		}
	}

	return head, nil
}

func ParseLinkFlags(flagsRaw uint32) LinkFlagsParsed {
	return LinkFlagsParsed{
		HasLinkTargetIDList:         (flagsRaw & HasLinkTargetIDList) != 0,
		HasLinkInfo:                 (flagsRaw & HasLinkInfo) != 0,
		HasName:                     (flagsRaw & HasName) != 0,
		HasRelativePath:             (flagsRaw & HasRelativePath) != 0,
		HasWorkingDir:               (flagsRaw & HasWorkingDir) != 0,
		HasArguments:                (flagsRaw & HasArguments) != 0,
		HasIconLocation:             (flagsRaw & HasIconLocation) != 0,
		IsUnicode:                   (flagsRaw & IsUnicode) != 0,
		ForceNoLinkInfo:             (flagsRaw & ForceNoLinkInfo) != 0,
		HasExpString:                (flagsRaw & HasExpString) != 0,
		RunInSeparateProcess:        (flagsRaw & RunInSeparateProcess) != 0,
		Unused1:                     (flagsRaw & Unused1) != 0,
		HasDarwinID:                 (flagsRaw & HasDarwinID) != 0,
		RunAsUser:                   (flagsRaw & RunAsUser) != 0,
		HasExpIcon:                  (flagsRaw & HasExpIcon) != 0,
		NoPidlAlias:                 (flagsRaw & NoPidlAlias) != 0,
		Unused2:                     (flagsRaw & Unused2) != 0,
		RunWithShimLayer:            (flagsRaw & RunWithShimLayer) != 0,
		ForceNoLinkTrack:            (flagsRaw & ForceNoLinkTrack) != 0,
		EnableTargetMetadata:        (flagsRaw & EnableTargetMetadata) != 0,
		DisableLinkPathTracking:     (flagsRaw & DisableLinkPathTracking) != 0,
		DisableKnownFolderTracking:  (flagsRaw & DisableKnownFolderTracking) != 0,
		DisableKnownFolderAlias:     (flagsRaw & DisableKnownFolderAlias) != 0,
		AllowLinkToLink:             (flagsRaw & AllowLinkToLink) != 0,
		UnaliasOnSave:               (flagsRaw & UnaliasOnSave) != 0,
		PreferEnvironmentPath:       (flagsRaw & PreferEnvironmentPath) != 0,
		KeepLocalIDListForUNCTarget: (flagsRaw & KeepLocalIDListForUNCTarget) != 0,

		HTMLNoSubDirCreation:     (flagsRaw & HTMLNoSubDirCreation) != 0,
		DisallowUserView:         (flagsRaw & DisallowUserView) != 0,
		ForcePerceivedTypeSystem: (flagsRaw & ForcePerceivedTypeSystem) != 0,
		IncludeSlowInfo:          (flagsRaw & IncludeSlowInfo) != 0,

		// ReservedForNYI:                   (flagsRaw & ReservedForNYI) != 0,
		// DisableShadowCopy:                (flagsRaw & DisableShadowCopy) != 0,
		// DisableKnownFolderAliasMigration: (flagsRaw & DisableKnownFolderAliasMigration) != 0,
		// DisableShellFolderVirtualization: (flagsRaw & DisableShellFolderVirtualization) != 0,
	}
}

func ParseFileAttributes(flagsRaw uint32) FileAttributesParsed {
	return FileAttributesParsed{
		FileAttributeReadOnly:          (flagsRaw & FileAttributeReadOnly) != 0,
		FileAttributeHidden:            (flagsRaw & FileAttributeHidden) != 0,
		FileAttributeSystem:            (flagsRaw & FileAttributeSystem) != 0,
		FileAttributeVolumeLabel:       (flagsRaw & FileAttributeVolumeLabel) != 0,
		FileAttributeDirectory:         (flagsRaw & FileAttributeDirectory) != 0,
		FileAttributeArchive:           (flagsRaw & FileAttributeArchive) != 0,
		FileAttributeNormal:            (flagsRaw & FileAttributeNormal) != 0,
		FileAttributeTemporary:         (flagsRaw & FileAttributeTemporary) != 0,
		FileAttributeSparseFile:        (flagsRaw & FileAttributeSparseFile) != 0,
		FileAttributeReparsePoint:      (flagsRaw & FileAttributeReparsePoint) != 0,
		FileAttributeCompressed:        (flagsRaw & FileAttributeCompressed) != 0,
		FileAttributeOffline:           (flagsRaw & FileAttributeOffline) != 0,
		FileAttributeNotContentIndexed: (flagsRaw & FileAttributeNotContentIndexed) != 0,
		FileAttributeEncrypted:         (flagsRaw & FileAttributeEncrypted) != 0,
	}
}

func ParseIDList(idListData []byte) ([]ItemID, error) {
	var itemIDList = []ItemID{}
	reader := bytes.NewReader(idListData)
	for {
		var itemIDSize uint16
		err := binary.Read(reader, binary.LittleEndian, &itemIDSize)
		if err == io.EOF {
			break
		}
		if err != nil {
			return itemIDList, err
		}

		// Check for null-terminator
		if itemIDSize == 0 {
			break
		}
		//TODO utf
		dataSize := itemIDSize - 2
		data := make([]byte, dataSize)
		_, err = reader.Read(data)
		if err != nil {
			return itemIDList, err
		}

		itemIDList = append(itemIDList, ItemID{
			ItemIDSize:       itemIDSize,
			ItemIDDataBase64: base64.StdEncoding.EncodeToString(data),
			ItemIDData:       data,
		})
	}
	return itemIDList, nil

}

func ParseLinkTargetIDList(r *bytes.Reader) (LinkTargetIDList, error) {
	var idListSize uint16

	err := binary.Read(r, binary.LittleEndian, &idListSize)
	//TODO redo
	if err != nil {
		return LinkTargetIDList{
			IDListSize: 0,
			IDListData: IDList{
				ItemIDs: []ItemID{
					ItemID{
						ItemIDSize: 0,
						ItemIDData: []byte{},
					},
				},
			},
		}, err
	}

	if idListSize == 0 {
		return LinkTargetIDList{
			IDListSize: 0,
			IDListData: IDList{
				ItemIDs: []ItemID{
					ItemID{
						ItemIDSize: 0,
						ItemIDData: []byte{},
					},
				},
			},
		}, nil
	}

	idListData := make([]byte, idListSize)

	_, err = r.Read(idListData)
	if err != nil {
		return LinkTargetIDList{
			IDListSize: 0,
			IDListData: IDList{
				ItemIDs: []ItemID{
					ItemID{
						ItemIDSize: 0,
						ItemIDData: []byte{},
					},
				},
			},
		}, err
	}

	itemID, err := ParseIDList(idListData)

	if err != nil {
		return LinkTargetIDList{
			IDListSize: 0,
			IDListData: IDList{
				ItemIDs: []ItemID{
					ItemID{
						ItemIDSize: 0,
						ItemIDData: []byte{},
					},
				},
			},
		}, err
	}

	return LinkTargetIDList{
		IDListSize: idListSize,
		IDListData: IDList{
			ItemIDs: itemID,
		},
	}, nil

}

func readByteStringZeroTerminated(r *bytes.Reader) (str string, b64 string, err error) {
	counter := 0
	var byteString []byte
	for {
		b, err := r.ReadByte()
		if err != nil {
			byteString = append(byteString, 0x0)
			return string(byteString), base64.StdEncoding.EncodeToString(byteString), err
		}
		byteString = append(byteString, b)
		if b == 0x0 {
			break
		}
		counter++
		if counter > 512 {
			byteString = append(byteString, 0x0)
			break
		}
	}
	return string(byteString), base64.StdEncoding.EncodeToString(byteString), nil
}

func ParseLinkInfo(r *bytes.Reader) (LinkInfo, error) {
	var linkInfo LinkInfo
	var linkInfoSize uint32
	err := binary.Read(r, binary.LittleEndian, &linkInfoSize)
	if err != nil {
		return linkInfo, err
	}
	linkInfo.LinkInfoSize = linkInfoSize

	dataSize := linkInfoSize - 4
	linkInfoData := make([]byte, dataSize)

	linkInfoReader := bytes.NewReader(linkInfoData)

	err = binary.Read(linkInfoReader, binary.LittleEndian, &linkInfo.LinkInfoHeaderSize)
	if err != nil {
		return linkInfo, err
	}
	err = binary.Read(linkInfoReader, binary.LittleEndian, &linkInfo.LinkInfoFlags)
	if err != nil {
		return linkInfo, err
	}
	err = binary.Read(linkInfoReader, binary.LittleEndian, &linkInfo.VolumeIDOffset)
	if err != nil {
		return linkInfo, err
	}
	err = binary.Read(linkInfoReader, binary.LittleEndian, &linkInfo.LocalBasePathOffset)
	if err != nil {
		return linkInfo, err
	}
	err = binary.Read(linkInfoReader, binary.LittleEndian, &linkInfo.CommonNetworkRelativeLinkOffset)
	if err != nil {
		return linkInfo, err
	}
	err = binary.Read(linkInfoReader, binary.LittleEndian, &linkInfo.CommonPathSuffixOffset)
	if err != nil {
		return linkInfo, err
	}
	if linkInfo.LinkInfoHeaderSize >= LinkInfoHeaderSizeOptionalFieldsSpecifiedFrom {
		err = binary.Read(linkInfoReader, binary.LittleEndian, &linkInfo.LocalBasePathOffsetUnicode)
		if err != nil {
			return linkInfo, err
		}
		err = binary.Read(linkInfoReader, binary.LittleEndian, &linkInfo.CommonPathSuffixOffsetUnicode)
		if err != nil {
			return linkInfo, err
		}
	} else {
		linkInfo.LocalBasePathOffsetUnicode = 0
		linkInfo.CommonPathSuffixOffsetUnicode = 0
	}
	if (linkInfo.LinkInfoFlags & VolumeIDAndLocalBasePathPresent) != 0 {
		var volumeID VolumeID
		if linkInfo.VolumeIDOffset != 0 {
			_, err = linkInfoReader.Seek(int64(linkInfo.VolumeIDOffset-4), io.SeekStart)
			if err != nil {
				return linkInfo, err
			}

			err = binary.Read(linkInfoReader, binary.LittleEndian, &volumeID.VolumeIDSize)
			if err != nil {
				return linkInfo, err
			}
			//TODO throw if < VolumeIDSizeMin

			err = binary.Read(linkInfoReader, binary.LittleEndian, &volumeID.DriveSerialNumber)
			if err != nil {
				return linkInfo, err
			}
			err = binary.Read(linkInfoReader, binary.LittleEndian, &volumeID.VolumeLabelOffset)
			if err != nil {
				return linkInfo, err
			}
			volumeLabelSize := volumeID.VolumeIDSize - VolumeIDSizeMin
			volumeLabelData := make([]byte, volumeLabelSize)
			if volumeID.VolumeLabelOffset == VolumeLabelOffsetUnicodePresent {
				err = binary.Read(linkInfoReader, binary.LittleEndian, &volumeID.VolumeLabelOffsetUnicode)
				if err != nil {
					return linkInfo, err
				}
				volumeID.VolumeLabel = ""
				//TODO set by offset
				err = binary.Read(linkInfoReader, binary.LittleEndian, volumeLabelData)

				if err != nil {
					return linkInfo, err
				}
				//TODO utf-16BE, utf-16LE, utf-32, utf-32LE, utf-32BE
				volumeID.VolumeLableBase64 = base64.StdEncoding.EncodeToString(volumeLabelData)
				if utf8.Valid(volumeLabelData) {
					volumeID.VolumeLabelUnicode = string(volumeLabelData)
				} else {
					// utf16s := utf16.Decode(volumeLabelData)
					// if err != nil {
					// 	//TODO other
					// }
					volumeID.VolumeLabelUnicode = string(volumeLabelData)
				}

			} else {
				volumeID.VolumeLabelOffsetUnicode = 0
				volumeID.VolumeLabelUnicode = ""
				//TODO set by offset
				err = binary.Read(linkInfoReader, binary.LittleEndian, volumeLabelData)
				if err != nil {
					return linkInfo, err
				}
				volumeID.VolumeLabel = string(volumeLabelData)
				volumeID.VolumeLableBase64 = base64.StdEncoding.EncodeToString(volumeLabelData)
			}
		}
		linkInfo.VolumeID = volumeID

		if linkInfo.LocalBasePathOffset != 0 {
			_, err = linkInfoReader.Seek(int64(linkInfo.LocalBasePathOffset-4), io.SeekStart)
			if err != nil {
				return linkInfo, err
			}
			linkInfo.LocalBasePath, linkInfo.LocalBasePathBase64, err = readByteStringZeroTerminated(linkInfoReader)
			if err != nil {
				return linkInfo, err
			}
		}

		if linkInfo.LinkInfoHeaderSize > LinkInfoHeaderSizeOptionalFieldsSpecifiedFrom {
			if linkInfo.LocalBasePathOffsetUnicode != 0 {
				_, err = linkInfoReader.Seek(int64(linkInfo.LocalBasePathOffsetUnicode-4), io.SeekStart)
				if err != nil {
					return linkInfo, err
				}
				linkInfo.LocalBasePath, linkInfo.LocalBasePathBase64, err = readByteStringZeroTerminated(linkInfoReader)
				if err != nil {
					return linkInfo, err
				}
			}
		}
		linkInfo.VolumeID = volumeID
	}

	if (linkInfo.LinkInfoFlags & CommonNetworkRelativeLinkAndPathSuffixPresent) != 0 {
		if linkInfo.CommonNetworkRelativeLinkOffset != 0 {
			var commonNetworkRelativeLink CommonNetworkRelativeLink
			_, err = linkInfoReader.Seek(int64(linkInfo.CommonNetworkRelativeLinkOffset-4), io.SeekStart)
			if err != nil {
				linkInfo.CommonNetworkRelativeLink = commonNetworkRelativeLink
				return linkInfo, err
			}
			err = binary.Read(linkInfoReader, binary.LittleEndian, commonNetworkRelativeLink.CommonNetworkRelativeLinkSize)
			if err != nil {
				return linkInfo, err
			}
			commonNetworkRelativeLinkData := make([]byte, commonNetworkRelativeLink.CommonNetworkRelativeLinkSize-4)
			err = binary.Read(linkInfoReader, binary.LittleEndian, commonNetworkRelativeLinkData)
			if err != nil {
				return linkInfo, err
			}
			commonNetworkRelativeLinkReader := bytes.NewReader(commonNetworkRelativeLinkData)

			err = binary.Read(commonNetworkRelativeLinkReader, binary.LittleEndian, commonNetworkRelativeLink.CommonNetworkRelativeLinkFlags)
			if err != nil {
				return linkInfo, err
			}

			err = binary.Read(commonNetworkRelativeLinkReader, binary.LittleEndian, commonNetworkRelativeLink.NetNameOffset)
			if err != nil {
				return linkInfo, err
			}

			err = binary.Read(commonNetworkRelativeLinkReader, binary.LittleEndian, commonNetworkRelativeLink.DeviceNameOffset)
			if err != nil {
				return linkInfo, err
			}
			//TODO throw const mismatch flag

			err = binary.Read(commonNetworkRelativeLinkReader, binary.LittleEndian, commonNetworkRelativeLink.NetworkProviderType)
			if err != nil {
				return linkInfo, err
			}
			//TODO throw const mismatch flag

			err = binary.Read(commonNetworkRelativeLinkReader, binary.LittleEndian, commonNetworkRelativeLink.NetNameOffsetUnicode)
			if err != nil {
				return linkInfo, err
			}

			err = binary.Read(commonNetworkRelativeLinkReader, binary.LittleEndian, commonNetworkRelativeLink.DeviceNameOffsetUnicode)
			if err != nil {
				return linkInfo, err
			}

			if commonNetworkRelativeLink.NetNameOffset != 0 {
				_, err = commonNetworkRelativeLinkReader.Seek(int64(commonNetworkRelativeLink.NetNameOffset-4), io.SeekStart)
				if err != nil {
					commonNetworkRelativeLink.NetName = ""
				} else {
					commonNetworkRelativeLink.NetName, commonNetworkRelativeLink.NetNameBase64, err = readByteStringZeroTerminated(commonNetworkRelativeLinkReader)
					if err != nil {
						commonNetworkRelativeLink.NetName = ""
					}
				}
			}
			if commonNetworkRelativeLink.DeviceNameOffset != 0 {
				_, err = commonNetworkRelativeLinkReader.Seek(int64(commonNetworkRelativeLink.DeviceNameOffset-4), io.SeekStart)
				if err != nil {
					commonNetworkRelativeLink.DeviceName = ""
				} else {
					commonNetworkRelativeLink.DeviceName, commonNetworkRelativeLink.DeviceNameBase64, err = readByteStringZeroTerminated(commonNetworkRelativeLinkReader)
					if err != nil {
						commonNetworkRelativeLink.DeviceName = ""
					}
				}
			}
			if commonNetworkRelativeLink.NetNameOffset > CommonNetworkRelativeLinkUnicodeMinOffsets {
				_, err = commonNetworkRelativeLinkReader.Seek(int64(commonNetworkRelativeLink.NetNameOffsetUnicode-4), io.SeekStart)
				if err != nil {
					commonNetworkRelativeLink.NetNameUnicode = ""
				} else {
					if commonNetworkRelativeLink.NetNameOffsetUnicode > 0 {
						commonNetworkRelativeLink.NetNameUnicode, commonNetworkRelativeLink.NetNameBase64, err = readByteStringZeroTerminated(commonNetworkRelativeLinkReader)
						if err != nil {
							commonNetworkRelativeLink.NetNameUnicode = ""
						}
					}
				}
				_, err = commonNetworkRelativeLinkReader.Seek(int64(commonNetworkRelativeLink.DeviceNameOffsetUnicode-4), io.SeekStart)
				if err != nil {
					commonNetworkRelativeLink.DeviceNameUnicode = ""
				} else {
					if commonNetworkRelativeLink.DeviceNameOffsetUnicode > 0 {
						commonNetworkRelativeLink.DeviceNameUnicode, commonNetworkRelativeLink.DeviceNameBase64, err = readByteStringZeroTerminated(commonNetworkRelativeLinkReader)
						if err != nil {
							commonNetworkRelativeLink.DeviceNameUnicode = ""
						}
					}
				}

			}
			linkInfo.CommonNetworkRelativeLink = commonNetworkRelativeLink
		}
	}

	if linkInfo.CommonPathSuffixOffset != 0 {
		_, err = linkInfoReader.Seek(int64(linkInfo.CommonPathSuffixOffset-4), io.SeekStart)
		if err != nil {
			return linkInfo, err
		}
		linkInfo.CommonPathSuffix, linkInfo.CommonPathSuffixBase64, err = readByteStringZeroTerminated(linkInfoReader)
		if err != nil {
			linkInfo.CommonPathSuffix = ""
		}
	}

	if linkInfo.LinkInfoHeaderSize >= LinkInfoHeaderSizeOptionalFieldsSpecifiedFrom {
		if linkInfo.CommonPathSuffixOffsetUnicode >= 0 {
			_, err = linkInfoReader.Seek(int64(linkInfo.CommonPathSuffixOffsetUnicode-4), io.SeekStart)
			if err != nil {
				return linkInfo, err
			}
			linkInfo.CommonPathSuffixUnicode, linkInfo.CommonPathSuffixBase64, err = readByteStringZeroTerminated(linkInfoReader)
			if err != nil {
				linkInfo.CommonPathSuffixUnicode = ""
			}
			//TODO read
		}
	}
	return linkInfo, nil
}

func readByteStringSizeSpecified(r *bytes.Reader, size uint64) (str string, b64 string, err error) {

	var byteString []byte
	for counter := 0; counter < size; counter++ {
		b, err := r.ReadByte()
		if err != nil {
			break
		}
		byteString = append(byteString, b)
	}
	byteString = append(byteString, 0x0)
	return string(byteString), base64.StdEncoding.EncodeToString(byteString), nil
}

func ParseStringData(r *bytes.Reader, linkFlagsParsed LinkFlagsParsed) (StringData, error) {
	var stringData StringData
	var countCharacters int16

	if linkFlagsParsed.HasName {
		err = binary.Read(r, binary.LittleEndian, countCharacters)
		if err != nil {
			stringData.NameString = ""
			return stringData, err
		}
		stringData.NameString, stringData.NameStringBase64, err = readByteStringSizeSpecified(r, uint64(countCharacters))
		if err != nil {
			return stringData, err
		}
	}
	//TODO firtsly
	if linkFlagsParsed.HasRelativePath {

	}

	if linkFlagsParsed.HasWorkingDir {

	}

	if linkFlagsParsed.HasArguments {

	}

	if linkFlagsParsed.HasIconLocation {

	}

}

func ParseData(r *bytes.Reader) (ShellLinkParsed, error) {
	var shellLinkParsed ShellLinkParsed
	shellLinkHeader, err := ParseShellLinkHeader(r)
	if err != nil {
		switch t := err.(type) {
		default:
			return shellLinkParsed, err
		case *ConstMismatchError:
			fmt.Println("ConstMismatchError", t)
			err.Error()
		}

	}
	shellLinkParsed.header = shellLinkHeader

	linkFlagsParsed := ParseLinkFlags(shellLinkHeader.LinkFlags)
	shellLinkParsed.linkFlagsParsed = linkFlagsParsed

	fileAttributesParsed := ParseFileAttributes(shellLinkHeader.FileAttributes)
	shellLinkParsed.fileAttributesParsed = fileAttributesParsed

	//TODO: parse HotKeyFlags

	var linkTargetIDList LinkTargetIDList
	if linkFlagsParsed.HasLinkTargetIDList {
		linkTargetIDList, err = ParseLinkTargetIDList(r)
		if err != nil {
			return shellLinkParsed, err
		}
	}
	shellLinkParsed.linkTargetIDList = linkTargetIDList

	var linkInfo LinkInfo
	if linkFlagsParsed.HasLinkInfo {
		linkInfo, err = ParseLinkInfo(r)
		if err != nil {
			return shellLinkParsed, err
		}
	}
	shellLinkParsed.linkInfo = linkInfo

	var stringData StringData
	if linkFlagsParsed.HasName || linkFlagsParsed.HasRelativePath ||
		linkFlagsParsed.HasWorkingDir || linkFlagsParsed.HasArguments ||
		linkFlagsParsed.HasIconLocation {
		stringData, err = ParseStringData(r, linkFlagsParsed)
	}
}

func main() {
	// Replace "example.lnk" with the path to your .lnk file.
	data, err := ReadLnkFile("FineReaderPortable.lnk")
	if err != nil {
		fmt.Printf("Error reading .lnk file: %v\n", err)
		return
	}

	reader := bytes.NewReader(data)
	header, err := ParseShellLinkHeader(reader)

}
