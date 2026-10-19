package cdio

import (
	"fmt"
	"unsafe"

	"github.com/ebitengine/purego"
)

type Device unsafe.Pointer
type Cdda unsafe.Pointer
type Paranoia unsafe.Pointer
type Lsn int32

type Track uint8
type Cdtext unsafe.Pointer
type DriverReturnCode int8 // 0 is success, negative values are error
type CdtextField uint8
type ParanoiaMode uint8

type LogLevel int

const (
	LogDebug  LogLevel = 1
	LogInfo   LogLevel = 2
	LogWarn   LogLevel = 3
	LogError  LogLevel = 4
	LogAssert LogLevel = 5
)

const CDIO_CD_FRAMESIZE_RAW int = 2352

const (
	ParanoiaModeOff       ParanoiaMode = 0
	ParanoiaModeOverlap   ParanoiaMode = 1
	ParanoiaModeScratch   ParanoiaMode = 2
	ParanoiaModeRepair    ParanoiaMode = 4
	ParanoiaModeFull      ParanoiaMode = 0xff
	ParanoiaModeNeverSkip ParanoiaMode = 0x20
)

const LeadoutTrack Track = 0xAA
const MessageForgetIt int32 = 0
const MessagePrintIt int32 = 1
const MessageLogIt int32 = 2
const SeekSet int32 = 0 // libc

var (
	// libcdio.so
	Open                func(string, unsafe.Pointer) Device
	OpenAm              func(string, int32, unsafe.Pointer) Device
	rawLibInit          func() int32
	Destroy             func(Device)
	rawGetDefaultDevice func(Device) *byte
	GetFirstTrackNum    func(Device) Track
	GetLastTrackNum     func(Device) Track
	GetNumTracks        func(Device) Track
	GetCdtext           func(Device) Cdtext
	GetCdtextRaw        func(Device) unsafe.Pointer
	rawCdtextGet        func(Cdtext, CdtextField, Track) *byte
	rawGetMcn           func(Device) *byte
	rawGetTrackIsrc     func(Device, Track) *byte
	GetTrackLba         func(Device, Track) int32
	GetTrackLsn         func(Device, Track) Lsn
	GetTrackSecCount    func(Device, Track) uint32
	GetTrackChannels    func(Device, Track) int32
	GetTrackCopyPermit  func(Device, Track) int32
	GetTrackFormat      func(Device, Track) int32
	GetTrackPreemphasis func(Device, Track) int32
	GetTrackLastLsn     func(Device, Track) Lsn
	GetDiscLastLsn      func(Device) Lsn
	GetDiscmode         func(Device) int32
	GetDriverId         func(Device) int32
	rawGetDriverName    func(Device) *byte
	GetHwinfo           func(Device) unsafe.Pointer
	rawGetMediaChanged  func(Device) int32
	EjectMedia          func(Device) int32
	EjectMediaDrive     func(string) int32
	CloseTray           func(string, unsafe.Pointer) int32
	SetSpeed            func(Device, int32) int32
	GetDevices          func(int32) **byte
	GetDevicesRet       func(int32, *uint32) **byte
	GetDevicesWithCap   func(unsafe.Pointer, int32) **byte
	FreeDeviceList      func(**byte)
	IsDevice            func(string) bool
	IsDiscmodeCdrom     func(int32) int32
	IsDiscmodeDvd       func(int32) int32
	ReadAudioSector     func(Device, unsafe.Pointer, Lsn) int32
	ReadAudioSectors    func(Device, unsafe.Pointer, Lsn, uint32) int32
	LogSetHandler       func(uintptr)

	// Mmc specific (often part of libcdio)
	rawMmcGetMcn        func(Device) *byte
	rawMmcGetTrackIsrc  func(Device, Track) *byte
	MmcEjectMedia       func(Device) int32
	MmcTestUnitReady    func(Device, uint32) DriverReturnCode // uint32 is timeout in ms
	rawMmcGetTrayStatus func(Device) int32

	// libcdio_cdda.so
	CddapOpen               func(Cdda) DriverReturnCode
	rawCddapClose           func(Cdda) int32
	rawCddapCloseNoFreeCdio func(Cdda) int32
	CddapIdentify           func(string, int32, unsafe.Pointer) Cdda
	CddapIdentifyCdio       func(Device, int32, unsafe.Pointer) Cdda
	CddapDiscFirstsector    func(Cdda) Lsn
	CddapDiscLastsector     func(Cdda) Lsn
	rawCddapMessages        func(Cdda) *byte
	rawCddapErrors          func(Cdda) *byte
	CddapFreeMessages       func(Cdda)
	CddapVerboseSet         func(Cdda, int32, int32)
	CddapSpeedSet           func(Cdda, int32) int32
	CddapTrackFirstsector   func(Cdda, Track) Lsn
	CddapTrackLastsector    func(Cdda, Track) Lsn
	CddapTrackAudiop        func(Cdda, Track) int32
	CddapTrackChannels      func(Cdda, Track) int32
	CddapTrackCopyp         func(Cdda, Track) int32
	CddapTrackPreemp        func(Cdda, Track) int32
	CddapTracks             func(Cdda) int32
	CddapSectorGettrack     func(Cdda, Lsn) Track
	CddapRead               func(Cdda, unsafe.Pointer) *byte
	CddapReadTimed          func(Cdda, unsafe.Pointer, *int32) *byte
	rawCddapVersion         func() *byte
	rawCddapFindACdrom      func() *byte

	// libcdio_paranoia.so
	ParanoiaInit           func(Cdda) Paranoia
	ParanoiaFree           func(Paranoia)
	ParanoiaModeset        func(Paranoia, ParanoiaMode)
	ParanoiaSeek           func(Paranoia, Lsn, int32) Lsn
	ParanoiaRead           func(Paranoia, unsafe.Pointer) *[CDIO_CD_FRAMESIZE_RAW]byte
	ParanoiaReadLimited    func(Paranoia, unsafe.Pointer, int32) *[CDIO_CD_FRAMESIZE_RAW]byte
	ParanoiaOverlapset     func(Paranoia, int32)
	ParanoiaSetRange       func(Paranoia, Track, Track)
	ParanoiaCachemodelSize func(Paranoia, int32)
	rawParanoiaVersion     func() *byte
)

func cStringToString(ptr *byte) string {
	if ptr == nil {
		return ""
	}
	var b []byte
	p := uintptr(unsafe.Pointer(ptr))
	for {
		c := *(*byte)(unsafe.Pointer(p))
		if c == 0 {
			break
		}
		b = append(b, c)
		p++
	}
	return string(b)
}

func LibInit() bool {
	return rawLibInit() != 0
}

func GetDefaultDevice(d Device) string {
	return cStringToString(rawGetDefaultDevice(d))
}

func CdtextGet(cdtext Cdtext, key CdtextField, track Track) string {
	return cStringToString(rawCdtextGet(cdtext, key, track))
}

func GetMcn(d Device) string {
	return cStringToString(rawGetMcn(d))
}

func GetTrackIsrc(d Device, track Track) string {
	return cStringToString(rawGetTrackIsrc(d, track))
}

func GetDriverName(d Device) string {
	return cStringToString(rawGetDriverName(d))
}

func GetMediaChanged(d Device) bool {
	return rawGetMediaChanged(d) != 0
}

func MmcGetMcn(d Device) string {
	return cStringToString(rawMmcGetMcn(d))
}

func MmcGetTrackIsrc(d Device, track Track) string {
	return cStringToString(rawMmcGetTrackIsrc(d, track))
}

func MmcGetTrayStatus(d Device) bool {
	return rawMmcGetTrayStatus(d) != 0
}

func CddapClose(cdda Cdda) bool {
	return rawCddapClose(cdda) != 0
}

func CddapCloseNoFreeCdio(cdda Cdda) bool {
	return rawCddapCloseNoFreeCdio(cdda) != 0
}

func CddapMessages(cdda Cdda) string {
	return cStringToString(rawCddapMessages(cdda))
}

func CddapErrors(cdda Cdda) string {
	return cStringToString(rawCddapErrors(cdda))
}

func CddapVersion() string {
	return cStringToString(rawCddapVersion())
}

func CddapFindACdrom() string {
	return cStringToString(rawCddapFindACdrom())
}

func ParanoiaVersion() string {
	return cStringToString(rawParanoiaVersion())
}

func Init() error {
	libcdio, err := purego.Dlopen("libcdio.so", purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		return fmt.Errorf("failed to load libcdio.so: %w", err)
	}

	libcdio_cdda, err := purego.Dlopen("libcdio_cdda.so", purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		return fmt.Errorf("failed to load libcdio_cdda.so: %w", err)
	}

	libcdio_paranoia, err := purego.Dlopen("libcdio_paranoia.so", purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		return fmt.Errorf("failed to load libcdio_paranoia.so: %w", err)
	}

	// libcdio.so
	purego.RegisterLibFunc(&Open, libcdio, "cdio_open")
	purego.RegisterLibFunc(&OpenAm, libcdio, "cdio_open_am")
	purego.RegisterLibFunc(&rawLibInit, libcdio, "cdio_init")
	purego.RegisterLibFunc(&Destroy, libcdio, "cdio_destroy")
	purego.RegisterLibFunc(&rawGetDefaultDevice, libcdio, "cdio_get_default_device")
	purego.RegisterLibFunc(&GetFirstTrackNum, libcdio, "cdio_get_first_track_num")
	purego.RegisterLibFunc(&GetLastTrackNum, libcdio, "cdio_get_last_track_num")
	purego.RegisterLibFunc(&GetNumTracks, libcdio, "cdio_get_num_tracks")
	purego.RegisterLibFunc(&GetCdtext, libcdio, "cdio_get_cdtext")
	purego.RegisterLibFunc(&GetCdtextRaw, libcdio, "cdio_get_cdtext_raw")
	purego.RegisterLibFunc(&rawCdtextGet, libcdio, "cdtext_get")
	purego.RegisterLibFunc(&rawGetMcn, libcdio, "cdio_get_mcn")
	purego.RegisterLibFunc(&rawGetTrackIsrc, libcdio, "cdio_get_track_isrc")
	purego.RegisterLibFunc(&GetTrackLba, libcdio, "cdio_get_track_lba")
	purego.RegisterLibFunc(&GetTrackLsn, libcdio, "cdio_get_track_lsn")
	purego.RegisterLibFunc(&GetTrackSecCount, libcdio, "cdio_get_track_sec_count")
	purego.RegisterLibFunc(&GetTrackChannels, libcdio, "cdio_get_track_channels")
	purego.RegisterLibFunc(&GetTrackCopyPermit, libcdio, "cdio_get_track_copy_permit")
	purego.RegisterLibFunc(&GetTrackFormat, libcdio, "cdio_get_track_format")
	purego.RegisterLibFunc(&GetTrackPreemphasis, libcdio, "cdio_get_track_preemphasis")
	purego.RegisterLibFunc(&GetTrackLastLsn, libcdio, "cdio_get_track_last_lsn")
	purego.RegisterLibFunc(&GetDiscLastLsn, libcdio, "cdio_get_disc_last_lsn")
	purego.RegisterLibFunc(&GetDiscmode, libcdio, "cdio_get_discmode")
	purego.RegisterLibFunc(&GetDriverId, libcdio, "cdio_get_driver_id")
	purego.RegisterLibFunc(&rawGetDriverName, libcdio, "cdio_get_driver_name")
	purego.RegisterLibFunc(&GetHwinfo, libcdio, "cdio_get_hwinfo")
	purego.RegisterLibFunc(&rawGetMediaChanged, libcdio, "cdio_get_media_changed")
	purego.RegisterLibFunc(&EjectMedia, libcdio, "cdio_eject_media")
	purego.RegisterLibFunc(&EjectMediaDrive, libcdio, "cdio_eject_media_drive")
	purego.RegisterLibFunc(&CloseTray, libcdio, "cdio_close_tray")
	purego.RegisterLibFunc(&SetSpeed, libcdio, "cdio_set_speed")
	purego.RegisterLibFunc(&GetDevices, libcdio, "cdio_get_devices")
	purego.RegisterLibFunc(&GetDevicesRet, libcdio, "cdio_get_devices_ret")
	purego.RegisterLibFunc(&GetDevicesWithCap, libcdio, "cdio_get_devices_with_cap")
	purego.RegisterLibFunc(&FreeDeviceList, libcdio, "cdio_free_device_list")
	purego.RegisterLibFunc(&IsDevice, libcdio, "cdio_is_device")
	purego.RegisterLibFunc(&IsDiscmodeCdrom, libcdio, "cdio_is_discmode_cdrom")
	purego.RegisterLibFunc(&IsDiscmodeDvd, libcdio, "cdio_is_discmode_dvd")
	purego.RegisterLibFunc(&ReadAudioSector, libcdio, "cdio_read_audio_sector")
	purego.RegisterLibFunc(&ReadAudioSectors, libcdio, "cdio_read_audio_sectors")
	purego.RegisterLibFunc(&LogSetHandler, libcdio, "cdio_log_set_handler")

	// MMC
	purego.RegisterLibFunc(&rawMmcGetMcn, libcdio, "mmc_get_mcn")
	purego.RegisterLibFunc(&rawMmcGetTrackIsrc, libcdio, "mmc_get_track_isrc")
	purego.RegisterLibFunc(&MmcEjectMedia, libcdio, "mmc_eject_media")
	purego.RegisterLibFunc(&MmcTestUnitReady, libcdio, "mmc_test_unit_ready")
	purego.RegisterLibFunc(&rawMmcGetTrayStatus, libcdio, "mmc_get_tray_status")

	// libcdio_cdda.so
	purego.RegisterLibFunc(&CddapOpen, libcdio_cdda, "cdio_cddap_open")
	purego.RegisterLibFunc(&rawCddapClose, libcdio_cdda, "cdio_cddap_close")
	purego.RegisterLibFunc(&rawCddapCloseNoFreeCdio, libcdio_cdda, "cdio_cddap_close_no_free_cdio")
	purego.RegisterLibFunc(&CddapIdentify, libcdio_cdda, "cdio_cddap_identify")
	purego.RegisterLibFunc(&CddapIdentifyCdio, libcdio_cdda, "cdio_cddap_identify_cdio")
	purego.RegisterLibFunc(&CddapDiscFirstsector, libcdio_cdda, "cdio_cddap_disc_firstsector")
	purego.RegisterLibFunc(&CddapDiscLastsector, libcdio_cdda, "cdio_cddap_disc_lastsector")
	purego.RegisterLibFunc(&rawCddapMessages, libcdio_cdda, "cdio_cddap_messages")
	purego.RegisterLibFunc(&rawCddapErrors, libcdio_cdda, "cdio_cddap_errors")
	purego.RegisterLibFunc(&CddapFreeMessages, libcdio_cdda, "cdio_cddap_free_messages")
	purego.RegisterLibFunc(&CddapVerboseSet, libcdio_cdda, "cdio_cddap_verbose_set")
	purego.RegisterLibFunc(&CddapSpeedSet, libcdio_cdda, "cdio_cddap_speed_set")
	purego.RegisterLibFunc(&CddapTrackFirstsector, libcdio_cdda, "cdio_cddap_track_firstsector")
	purego.RegisterLibFunc(&CddapTrackLastsector, libcdio_cdda, "cdio_cddap_track_lastsector")
	purego.RegisterLibFunc(&CddapTrackAudiop, libcdio_cdda, "cdio_cddap_track_audiop")
	purego.RegisterLibFunc(&CddapTrackChannels, libcdio_cdda, "cdio_cddap_track_channels")
	purego.RegisterLibFunc(&CddapTrackCopyp, libcdio_cdda, "cdio_cddap_track_copyp")
	purego.RegisterLibFunc(&CddapTrackPreemp, libcdio_cdda, "cdio_cddap_track_preemp")
	purego.RegisterLibFunc(&CddapTracks, libcdio_cdda, "cdio_cddap_tracks")
	purego.RegisterLibFunc(&CddapSectorGettrack, libcdio_cdda, "cdio_cddap_sector_gettrack")
	purego.RegisterLibFunc(&CddapRead, libcdio_cdda, "cdio_cddap_read")
	purego.RegisterLibFunc(&CddapReadTimed, libcdio_cdda, "cdio_cddap_read_timed")
	purego.RegisterLibFunc(&rawCddapVersion, libcdio_cdda, "cdio_cddap_version")
	purego.RegisterLibFunc(&rawCddapFindACdrom, libcdio_cdda, "cdio_cddap_find_a_cdrom")

	// libcdio_paranoia.so
	purego.RegisterLibFunc(&ParanoiaInit, libcdio_paranoia, "cdio_paranoia_init")
	purego.RegisterLibFunc(&ParanoiaFree, libcdio_paranoia, "cdio_paranoia_free")
	purego.RegisterLibFunc(&ParanoiaModeset, libcdio_paranoia, "cdio_paranoia_modeset")
	purego.RegisterLibFunc(&ParanoiaSeek, libcdio_paranoia, "cdio_paranoia_seek")
	purego.RegisterLibFunc(&ParanoiaRead, libcdio_paranoia, "cdio_paranoia_read")
	purego.RegisterLibFunc(&ParanoiaReadLimited, libcdio_paranoia, "cdio_paranoia_read_limited")
	purego.RegisterLibFunc(&ParanoiaOverlapset, libcdio_paranoia, "cdio_paranoia_overlapset")
	purego.RegisterLibFunc(&ParanoiaSetRange, libcdio_paranoia, "cdio_paranoia_set_range")
	purego.RegisterLibFunc(&ParanoiaCachemodelSize, libcdio_paranoia, "cdio_paranoia_cachemodel_size")
	purego.RegisterLibFunc(&rawParanoiaVersion, libcdio_paranoia, "cdio_paranoia_version")

	return nil
}
