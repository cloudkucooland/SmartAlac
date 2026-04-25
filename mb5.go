package sa

import (
	"fmt"
	"unsafe"

	"github.com/ebitengine/purego"
)

// C-compatible types
type mb5_release_list unsafe.Pointer
type mb5_release unsafe.Pointer
type mb5_query unsafe.Pointer
type mb5_metadata unsafe.Pointer
type mb5_discid string
type mb5_tQueryResult int // enum eQuery_Success=0, >0 == err
type mb5_disc unsafe.Pointer
type mb5_artist_credit unsafe.Pointer
type mb5_media_list unsafe.Pointer
type mb5_medium unsafe.Pointer
type mb5_track_list unsafe.Pointer
type mb5_track unsafe.Pointer
type mb5_recording unsafe.Pointer
type mb5_namecreditlist unsafe.Pointer
type mb5_namecredit unsafe.Pointer
type mb5_artist unsafe.Pointer
type mb5_label_info_list unsafe.Pointer
type mb5_label_info unsafe.Pointer
type mb5_label unsafe.Pointer
type mb5_release_group unsafe.Pointer
type mb5_relation_list_list unsafe.Pointer
type mb5_relation_list unsafe.Pointer
type mb5_relation unsafe.Pointer
type mb5_work unsafe.Pointer
type mb5_iswc_list unsafe.Pointer
type mb5_iswc unsafe.Pointer
type mb5_isrc_list unsafe.Pointer
type mb5_isrc unsafe.Pointer

var (
	mb5_query_new                    func(string, string, int) mb5_query
	mb5_query_lookup_discid          func(mb5_query, mb5_discid) mb5_release_list
	mb5_release_list_size            func(mb5_release_list) int
	mb5_release_list_item            func(mb5_release_list, int) mb5_release
	mb5_release_get_id               func(unsafe.Pointer, *byte, int)
	mb5_query_query                  func(mb5_query, string, string, string, int, unsafe.Pointer, unsafe.Pointer) mb5_metadata
	mb5_metadata_get_release         func(mb5_metadata) mb5_release
	mb5_metadata_get_releaselist     func(mb5_metadata) mb5_release_list
	mb5_query_get_lasterrormessage   func(mb5_query, *byte, int)
	mb5_query_get_lastresult         func(mb5_query) mb5_tQueryResult
	mb5_query_get_lasthttpcode       func(mb5_query) int
	mb5_metadata_get_disc            func(mb5_metadata) mb5_disc
	mb5_disc_get_id                  func(mb5_disc, *byte, int)
	mb5_disc_get_releaselist         func(mb5_disc) mb5_release_list
	mb5_release_get_title            func(unsafe.Pointer, *byte, int)
	mb5_release_get_artistcredit     func(mb5_release) mb5_artist_credit
	mb5_metadata_delete              func(mb5_metadata)
	mb5_release_media_matching_discid func(mb5_release, string) mb5_media_list
	mb5_medium_get_tracklist         func(mb5_medium) mb5_track_list
	mb5_medium_list_item             func(mb5_media_list, int) mb5_medium
	mb5_medium_list_size             func(mb5_media_list) int
	mb5_medium_get_position          func(mb5_medium) int
	mb5_medium_get_format            func(unsafe.Pointer, *byte, int)
	mb5_disc_clone                   func(mb5_disc) mb5_disc
	mb5_disc_delete                  func(mb5_disc)
	mb5_release_list_clone           func(mb5_release_list) mb5_release_list
	mb5_release_list_delete          func(mb5_release_list)
	mb5_release_clone                func(mb5_release) mb5_release
	mb5_release_delete               func(mb5_release)
	mb5_medium_list_get_trackcount   func(mb5_media_list) int
	mb5_metadata_clone               func(mb5_metadata) mb5_metadata
	mb5_release_get_mediumlist       func(mb5_release) mb5_media_list
	mb5_track_list_clone             func(mb5_track_list) mb5_track_list
	mb5_track_list_delete            func(mb5_track_list)
	mb5_track_list_item              func(mb5_track_list, int) mb5_track
	mb5_track_get_title              func(unsafe.Pointer, *byte, int)
	mb5_track_get_artistcredit       func(mb5_track) mb5_artist_credit
	mb5_track_get_recording          func(mb5_track) mb5_recording
	mb5_recording_get_id             func(unsafe.Pointer, *byte, int)
	mb5_artistcredit_get_namecreditlist func(mb5_artist_credit) mb5_namecreditlist
	mb5_namecredit_list_get_count    func(mb5_namecreditlist) int
	mb5_namecredit_list_item         func(mb5_namecreditlist, int) mb5_namecredit
	mb5_namecredit_get_name          func(unsafe.Pointer, *byte, int)
	mb5_namecredit_get_joinphrase    func(unsafe.Pointer, *byte, int)
	mb5_track_list_get_count         func(mb5_track_list) int
	mb5_track_get_position           func(mb5_track) int
	mb5_recording_get_title          func(unsafe.Pointer, *byte, int)
	mb5_recording_get_artistcredit   func(mb5_recording) mb5_artist_credit
	mb5_namecredit_get_artist        func(mb5_namecredit) mb5_artist
	mb5_artist_get_name              func(unsafe.Pointer, *byte, int)
	mb5_artist_get_id                func(unsafe.Pointer, *byte, int)
	mb5_artist_get_sortname          func(unsafe.Pointer, *byte, int)

	// Additional fields for better tagging
	mb5_release_get_date             func(unsafe.Pointer, *byte, int)
	mb5_release_get_barcode          func(unsafe.Pointer, *byte, int)
	mb5_release_get_asin             func(unsafe.Pointer, *byte, int)
	mb5_release_get_country          func(unsafe.Pointer, *byte, int)
	mb5_release_get_disambiguation   func(unsafe.Pointer, *byte, int)
	mb5_release_get_labelinfolist    func(mb5_release) mb5_label_info_list
	mb5_labelinfo_list_item          func(mb5_label_info_list, int) mb5_label_info
	mb5_labelinfo_list_size          func(mb5_label_info_list) int
	mb5_labelinfo_get_label          func(mb5_label_info) mb5_label
	mb5_labelinfo_get_catalognumber  func(unsafe.Pointer, *byte, int)
	mb5_label_get_name               func(unsafe.Pointer, *byte, int)
	mb5_release_get_releasegroup     func(mb5_release) mb5_release_group
	mb5_releasegroup_get_id          func(unsafe.Pointer, *byte, int)
	mb5_releasegroup_get_title       func(unsafe.Pointer, *byte, int)
	mb5_releasegroup_get_firstreleasedate func(unsafe.Pointer, *byte, int)

	// Relations and Works
	mb5_recording_get_relationlistlist func(mb5_recording) mb5_relation_list_list
	mb5_relationlist_list_size         func(mb5_relation_list_list) int
	mb5_relationlist_list_item         func(mb5_relation_list_list, int) mb5_relation_list
	mb5_relation_list_size             func(mb5_relation_list) int
	mb5_relation_list_item             func(mb5_relation_list, int) mb5_relation
	mb5_relation_get_type              func(unsafe.Pointer, *byte, int)
	mb5_relation_get_target            func(unsafe.Pointer, *byte, int)
	mb5_relation_get_artist            func(mb5_relation) mb5_artist
	mb5_relation_get_work              func(mb5_relation) mb5_work
	mb5_work_get_id                    func(unsafe.Pointer, *byte, int)
	mb5_work_get_title                 func(unsafe.Pointer, *byte, int)
	mb5_work_get_iswclist              func(mb5_work) mb5_iswc_list
	mb5_iswc_list_size                 func(mb5_iswc_list) int
	mb5_iswc_list_item                 func(mb5_iswc_list, int) mb5_iswc
	mb5_iswc_get_iswc                  func(unsafe.Pointer, *byte, int)
	mb5_recording_get_isrclist         func(mb5_recording) mb5_isrc_list
	mb5_isrc_list_size                 func(mb5_isrc_list) int
	mb5_isrc_list_item                 func(mb5_isrc_list, int) mb5_isrc
	mb5_isrc_get_id                    func(unsafe.Pointer, *byte, int)
)

func (c *Curator) initMB5() error {
	libmusicbrainz5, err := purego.Dlopen("libmusicbrainz5.so", purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		return fmt.Errorf("failed to load libmusicbrainz5.so: %w", err)
	}

	purego.RegisterLibFunc(&mb5_query_new, libmusicbrainz5, "mb5_query_new")
	purego.RegisterLibFunc(&mb5_query_lookup_discid, libmusicbrainz5, "mb5_query_lookup_discid")
	purego.RegisterLibFunc(&mb5_release_list_size, libmusicbrainz5, "mb5_release_list_size")
	purego.RegisterLibFunc(&mb5_release_list_item, libmusicbrainz5, "mb5_release_list_item")
	purego.RegisterLibFunc(&mb5_release_get_id, libmusicbrainz5, "mb5_release_get_id")
	purego.RegisterLibFunc(&mb5_query_query, libmusicbrainz5, "mb5_query_query")
	purego.RegisterLibFunc(&mb5_metadata_get_release, libmusicbrainz5, "mb5_metadata_get_release")
	purego.RegisterLibFunc(&mb5_metadata_get_releaselist, libmusicbrainz5, "mb5_metadata_get_releaselist")
	purego.RegisterLibFunc(&mb5_query_get_lasterrormessage, libmusicbrainz5, "mb5_query_get_lasterrormessage")
	purego.RegisterLibFunc(&mb5_query_get_lastresult, libmusicbrainz5, "mb5_query_get_lastresult")
	purego.RegisterLibFunc(&mb5_query_get_lasthttpcode, libmusicbrainz5, "mb5_query_get_lasthttpcode")
	purego.RegisterLibFunc(&mb5_metadata_get_disc, libmusicbrainz5, "mb5_metadata_get_disc")
	purego.RegisterLibFunc(&mb5_disc_get_id, libmusicbrainz5, "mb5_disc_get_id")
	purego.RegisterLibFunc(&mb5_disc_get_releaselist, libmusicbrainz5, "mb5_disc_get_releaselist")
	purego.RegisterLibFunc(&mb5_release_get_title, libmusicbrainz5, "mb5_release_get_title")
	purego.RegisterLibFunc(&mb5_release_get_artistcredit, libmusicbrainz5, "mb5_release_get_artistcredit")
	purego.RegisterLibFunc(&mb5_release_media_matching_discid, libmusicbrainz5, "mb5_release_media_matching_discid")
	purego.RegisterLibFunc(&mb5_medium_get_tracklist, libmusicbrainz5, "mb5_medium_get_tracklist")
	purego.RegisterLibFunc(&mb5_medium_list_item, libmusicbrainz5, "mb5_medium_list_item")
	purego.RegisterLibFunc(&mb5_medium_list_size, libmusicbrainz5, "mb5_medium_list_size")
	purego.RegisterLibFunc(&mb5_medium_get_position, libmusicbrainz5, "mb5_medium_get_position")
	purego.RegisterLibFunc(&mb5_medium_get_format, libmusicbrainz5, "mb5_medium_get_format")
	purego.RegisterLibFunc(&mb5_disc_clone, libmusicbrainz5, "mb5_disc_clone")
	purego.RegisterLibFunc(&mb5_disc_delete, libmusicbrainz5, "mb5_disc_delete")
	purego.RegisterLibFunc(&mb5_release_list_clone, libmusicbrainz5, "mb5_release_list_clone")
	purego.RegisterLibFunc(&mb5_release_list_delete, libmusicbrainz5, "mb5_release_list_delete")
	purego.RegisterLibFunc(&mb5_release_clone, libmusicbrainz5, "mb5_release_clone")
	purego.RegisterLibFunc(&mb5_release_delete, libmusicbrainz5, "mb5_release_delete")
	purego.RegisterLibFunc(&mb5_medium_list_get_trackcount, libmusicbrainz5, "mb5_medium_list_get_trackcount")
	purego.RegisterLibFunc(&mb5_metadata_clone, libmusicbrainz5, "mb5_metadata_clone")
	purego.RegisterLibFunc(&mb5_metadata_delete, libmusicbrainz5, "mb5_metadata_delete")
	purego.RegisterLibFunc(&mb5_release_get_mediumlist, libmusicbrainz5, "mb5_release_get_mediumlist")
	purego.RegisterLibFunc(&mb5_track_list_clone, libmusicbrainz5, "mb5_track_list_clone")
	purego.RegisterLibFunc(&mb5_track_list_delete, libmusicbrainz5, "mb5_track_list_delete")
	purego.RegisterLibFunc(&mb5_track_list_item, libmusicbrainz5, "mb5_track_list_item")
	purego.RegisterLibFunc(&mb5_track_get_title, libmusicbrainz5, "mb5_track_get_title")
	purego.RegisterLibFunc(&mb5_track_get_artistcredit, libmusicbrainz5, "mb5_track_get_artistcredit")
	purego.RegisterLibFunc(&mb5_track_get_recording, libmusicbrainz5, "mb5_track_get_recording")
	purego.RegisterLibFunc(&mb5_recording_get_id, libmusicbrainz5, "mb5_recording_get_id")
	purego.RegisterLibFunc(&mb5_artistcredit_get_namecreditlist, libmusicbrainz5, "mb5_artistcredit_get_namecreditlist")
	purego.RegisterLibFunc(&mb5_namecredit_list_get_count, libmusicbrainz5, "mb5_namecredit_list_get_count")
	purego.RegisterLibFunc(&mb5_namecredit_list_item, libmusicbrainz5, "mb5_namecredit_list_item")
	purego.RegisterLibFunc(&mb5_namecredit_get_name, libmusicbrainz5, "mb5_namecredit_get_name")
	purego.RegisterLibFunc(&mb5_namecredit_get_joinphrase, libmusicbrainz5, "mb5_namecredit_get_joinphrase")
	purego.RegisterLibFunc(&mb5_track_list_get_count, libmusicbrainz5, "mb5_track_list_get_count")
	purego.RegisterLibFunc(&mb5_track_get_position, libmusicbrainz5, "mb5_track_get_position")
	purego.RegisterLibFunc(&mb5_recording_get_title, libmusicbrainz5, "mb5_recording_get_title")
	purego.RegisterLibFunc(&mb5_recording_get_artistcredit, libmusicbrainz5, "mb5_recording_get_artistcredit")
	purego.RegisterLibFunc(&mb5_namecredit_get_artist, libmusicbrainz5, "mb5_namecredit_get_artist")
	purego.RegisterLibFunc(&mb5_artist_get_name, libmusicbrainz5, "mb5_artist_get_name")
	purego.RegisterLibFunc(&mb5_artist_get_id, libmusicbrainz5, "mb5_artist_get_id")
	purego.RegisterLibFunc(&mb5_artist_get_sortname, libmusicbrainz5, "mb5_artist_get_sortname")

	purego.RegisterLibFunc(&mb5_release_get_date, libmusicbrainz5, "mb5_release_get_date")
	purego.RegisterLibFunc(&mb5_release_get_barcode, libmusicbrainz5, "mb5_release_get_barcode")
	purego.RegisterLibFunc(&mb5_release_get_asin, libmusicbrainz5, "mb5_release_get_asin")
	purego.RegisterLibFunc(&mb5_release_get_country, libmusicbrainz5, "mb5_release_get_country")
	purego.RegisterLibFunc(&mb5_release_get_disambiguation, libmusicbrainz5, "mb5_release_get_disambiguation")
	purego.RegisterLibFunc(&mb5_release_get_labelinfolist, libmusicbrainz5, "mb5_release_get_labelinfolist")
	purego.RegisterLibFunc(&mb5_labelinfo_list_item, libmusicbrainz5, "mb5_labelinfo_list_item")
	purego.RegisterLibFunc(&mb5_labelinfo_list_size, libmusicbrainz5, "mb5_labelinfo_list_size")
	purego.RegisterLibFunc(&mb5_labelinfo_get_label, libmusicbrainz5, "mb5_labelinfo_get_label")
	purego.RegisterLibFunc(&mb5_labelinfo_get_catalognumber, libmusicbrainz5, "mb5_labelinfo_get_catalognumber")
	purego.RegisterLibFunc(&mb5_label_get_name, libmusicbrainz5, "mb5_label_get_name")
	purego.RegisterLibFunc(&mb5_release_get_releasegroup, libmusicbrainz5, "mb5_release_get_releasegroup")
	purego.RegisterLibFunc(&mb5_releasegroup_get_id, libmusicbrainz5, "mb5_releasegroup_get_id")
	purego.RegisterLibFunc(&mb5_releasegroup_get_title, libmusicbrainz5, "mb5_releasegroup_get_title")
	purego.RegisterLibFunc(&mb5_releasegroup_get_firstreleasedate, libmusicbrainz5, "mb5_releasegroup_get_firstreleasedate")

	purego.RegisterLibFunc(&mb5_recording_get_relationlistlist, libmusicbrainz5, "mb5_recording_get_relationlistlist")
	purego.RegisterLibFunc(&mb5_relationlist_list_size, libmusicbrainz5, "mb5_relationlist_list_size")
	purego.RegisterLibFunc(&mb5_relationlist_list_item, libmusicbrainz5, "mb5_relationlist_list_item")
	purego.RegisterLibFunc(&mb5_relation_list_size, libmusicbrainz5, "mb5_relation_list_size")
	purego.RegisterLibFunc(&mb5_relation_list_item, libmusicbrainz5, "mb5_relation_list_item")
	purego.RegisterLibFunc(&mb5_relation_get_type, libmusicbrainz5, "mb5_relation_get_type")
	purego.RegisterLibFunc(&mb5_relation_get_target, libmusicbrainz5, "mb5_relation_get_target")
	purego.RegisterLibFunc(&mb5_relation_get_artist, libmusicbrainz5, "mb5_relation_get_artist")
	purego.RegisterLibFunc(&mb5_relation_get_work, libmusicbrainz5, "mb5_relation_get_work")
	purego.RegisterLibFunc(&mb5_work_get_id, libmusicbrainz5, "mb5_work_get_id")
	purego.RegisterLibFunc(&mb5_work_get_title, libmusicbrainz5, "mb5_work_get_title")
	purego.RegisterLibFunc(&mb5_work_get_iswclist, libmusicbrainz5, "mb5_work_get_iswclist")
	purego.RegisterLibFunc(&mb5_iswc_list_size, libmusicbrainz5, "mb5_iswc_list_size")
	purego.RegisterLibFunc(&mb5_iswc_list_item, libmusicbrainz5, "mb5_iswc_list_item")
	purego.RegisterLibFunc(&mb5_iswc_get_iswc, libmusicbrainz5, "mb5_iswc_get_iswc")
	purego.RegisterLibFunc(&mb5_recording_get_isrclist, libmusicbrainz5, "mb5_recording_get_isrclist")
	purego.RegisterLibFunc(&mb5_isrc_list_size, libmusicbrainz5, "mb5_isrc_list_size")
	purego.RegisterLibFunc(&mb5_isrc_list_item, libmusicbrainz5, "mb5_isrc_list_item")
	purego.RegisterLibFunc(&mb5_isrc_get_id, libmusicbrainz5, "mb5_isrc_get_id")

	return nil
}

// Go Helper Functions for MB5

func mb5String(f func(p1 unsafe.Pointer, buf *byte, size int), p unsafe.Pointer) string {
	if p == nil {
		return ""
	}
	buf := make([]byte, 256)
	f(p, &buf[0], len(buf))
	return string(buf[:cStringLen(buf)])
}

func cStringLen(buf []byte) int {
	for i, b := range buf {
		if b == 0 {
			return i
		}
	}
	return len(buf)
}
