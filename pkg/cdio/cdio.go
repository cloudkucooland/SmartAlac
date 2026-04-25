package cdio

import (
	"fmt"
	"unsafe"

	"github.com/ebitengine/purego"
)

type Device unsafe.Pointer
type Cdda unsafe.Pointer
type Paranoia unsafe.Pointer
type Lsn int // 32

type Track uint8
type Cdtext unsafe.Pointer
type DriverReturnCode int8 // 0 is success, negative values are error
type CdtextField uint8
type ParanoiaMode uint8

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
const MessageForgetIt int = 0
const MessagePrintIt int = 1
const MessageLogIt int = 2
const SeekSet int = 0 // libc

var (
	Open                  func(string, unsafe.Pointer) Device
	GetDefaultDevice      func(Device) string
	GetFirstTrackNum      func(Device) Track
	GetNumTracks          func(Device) Track
	GetCdtext             func(Device) Cdtext
	CdtextGet             func(Cdtext, CdtextField, Track) string
	MmcGetMcn             func(Device) string
	MmcGetTrackIsrc       func(Device, Track) string
	Destroy               func(Device)
	CdtextDestroy         func(Cdtext)
	GetTrackLba           func(Device, Track) int
	MmcEjectMedia         func(Device) int
	MmcTestUnitReady      func(Device, uint32) DriverReturnCode // uint32 is timeout in ms
	CddapOpen             func(Cdda) DriverReturnCode
	CddapDiscFirstsector  func(Cdda) Lsn
	CddapIdentifyCdio     func(Device, int, unsafe.Pointer) Cdda
	CddapMessages         func(Cdda) string
	CddapErrors           func(Cdda) string
	CddapCloseNoFreeCdio  func(Cdda) bool
	CddapVerboseSet       func(Cdda, int, int)
	CddapTrackFirstsector func(Cdda, Track) Lsn
	CddapTrackLastsector  func(Cdda, Track) Lsn
	ParanoiaInit          func(Cdda) Paranoia
	ParanoiaModeset       func(Paranoia, ParanoiaMode)
	ParanoiaSeek          func(Paranoia, Lsn, int) Lsn
	ParanoiaRead          func(Paranoia, unsafe.Pointer) *[CDIO_CD_FRAMESIZE_RAW]byte
	ParanoiaReadLimited   func(Paranoia, unsafe.Pointer, int) *[CDIO_CD_FRAMESIZE_RAW]byte
	ParanoiaFree          func(Paranoia)
	GetMediaChanged       func(Device) bool
	MmcGetTrayStatus      func(Device) bool
)

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

	purego.RegisterLibFunc(&Open, libcdio, "cdio_open")
	purego.RegisterLibFunc(&GetDefaultDevice, libcdio, "cdio_get_default_device")
	purego.RegisterLibFunc(&GetFirstTrackNum, libcdio, "cdio_get_first_track_num")
	purego.RegisterLibFunc(&GetNumTracks, libcdio, "cdio_get_num_tracks")
	purego.RegisterLibFunc(&GetCdtext, libcdio, "cdio_get_cdtext")
	purego.RegisterLibFunc(&CdtextGet, libcdio, "cdtext_get")
	purego.RegisterLibFunc(&MmcGetMcn, libcdio, "mmc_get_mcn")
	purego.RegisterLibFunc(&MmcGetTrackIsrc, libcdio, "mmc_get_track_isrc")
	purego.RegisterLibFunc(&Destroy, libcdio, "cdio_destroy")
	purego.RegisterLibFunc(&GetTrackLba, libcdio, "cdio_get_track_lba")
	purego.RegisterLibFunc(&MmcEjectMedia, libcdio, "mmc_eject_media")
	purego.RegisterLibFunc(&MmcTestUnitReady, libcdio, "mmc_test_unit_ready")
	purego.RegisterLibFunc(&GetMediaChanged, libcdio, "cdio_get_media_changed")
	purego.RegisterLibFunc(&MmcGetTrayStatus, libcdio, "mmc_get_tray_status")

	purego.RegisterLibFunc(&CddapOpen, libcdio_cdda, "cdio_cddap_open")
	purego.RegisterLibFunc(&CddapDiscFirstsector, libcdio_cdda, "cdio_cddap_disc_firstsector")
	purego.RegisterLibFunc(&CddapIdentifyCdio, libcdio_cdda, "cdio_cddap_identify_cdio")
	purego.RegisterLibFunc(&CddapMessages, libcdio_cdda, "cdio_cddap_messages")
	purego.RegisterLibFunc(&CddapErrors, libcdio_cdda, "cdio_cddap_errors")
	purego.RegisterLibFunc(&CddapCloseNoFreeCdio, libcdio_cdda, "cdio_cddap_close_no_free_cdio")
	purego.RegisterLibFunc(&CddapVerboseSet, libcdio_cdda, "cdio_cddap_verbose_set")
	purego.RegisterLibFunc(&CddapTrackFirstsector, libcdio_cdda, "cdio_cddap_track_firstsector")
	purego.RegisterLibFunc(&CddapTrackLastsector, libcdio_cdda, "cdio_cddap_track_lastsector")

	purego.RegisterLibFunc(&ParanoiaInit, libcdio_paranoia, "cdio_paranoia_init")
	purego.RegisterLibFunc(&ParanoiaModeset, libcdio_paranoia, "cdio_paranoia_modeset")
	purego.RegisterLibFunc(&ParanoiaSeek, libcdio_paranoia, "cdio_paranoia_seek")
	purego.RegisterLibFunc(&ParanoiaRead, libcdio_paranoia, "cdio_paranoia_read")
	purego.RegisterLibFunc(&ParanoiaReadLimited, libcdio_paranoia, "cdio_paranoia_read_limited")
	purego.RegisterLibFunc(&ParanoiaFree, libcdio_paranoia, "cdio_paranoia_free")

	return nil
}
