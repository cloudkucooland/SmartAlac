package mb5

import (
	"fmt"
	"unsafe"

	"github.com/ebitengine/purego"
)

// C-compatible types
type Alias unsafe.Pointer
type AliasList unsafe.Pointer
type Annotation unsafe.Pointer
type AnnotationList unsafe.Pointer
type Artist unsafe.Pointer
type ArtistList unsafe.Pointer
type ArtistCredit unsafe.Pointer
type Attribute unsafe.Pointer
type AttributeList unsafe.Pointer
type CDStub unsafe.Pointer
type CDStubList unsafe.Pointer
type Collection unsafe.Pointer
type CollectionList unsafe.Pointer
type Disc unsafe.Pointer
type DiscList unsafe.Pointer
type FreeDBDisc unsafe.Pointer
type FreeDBDiscList unsafe.Pointer
type IPI unsafe.Pointer
type IPIList unsafe.Pointer
type ISRC unsafe.Pointer
type ISRCList unsafe.Pointer
type ISWC unsafe.Pointer
type ISWCList unsafe.Pointer
type Label unsafe.Pointer
type LabelList unsafe.Pointer
type LabelInfo unsafe.Pointer
type LabelInfoList unsafe.Pointer
type Lifespan unsafe.Pointer
type Medium unsafe.Pointer
type MediumList unsafe.Pointer
type Message unsafe.Pointer
type Metadata unsafe.Pointer
type NameCredit unsafe.Pointer
type NameCreditList unsafe.Pointer
type NonMBTrack unsafe.Pointer
type NonMBTrackList unsafe.Pointer
type Offset unsafe.Pointer
type OffsetList unsafe.Pointer
type PUID unsafe.Pointer
type PUIDList unsafe.Pointer
type Query unsafe.Pointer
type Rating unsafe.Pointer
type Recording unsafe.Pointer
type RecordingList unsafe.Pointer
type Relation unsafe.Pointer
type RelationList unsafe.Pointer
type RelationListList unsafe.Pointer
type Release unsafe.Pointer
type ReleaseList unsafe.Pointer
type ReleaseGroup unsafe.Pointer
type ReleaseGroupList unsafe.Pointer
type SecondaryType unsafe.Pointer
type SecondaryTypeList unsafe.Pointer
type Tag unsafe.Pointer
type TagList unsafe.Pointer
type TextRepresentation unsafe.Pointer
type Track unsafe.Pointer
type TrackList unsafe.Pointer
type UserRating unsafe.Pointer
type UserTag unsafe.Pointer
type UserTagList unsafe.Pointer
type Work unsafe.Pointer
type WorkList unsafe.Pointer

type QueryResult int // enum eQuery_Success=0, >0 == err
type DiscID string

var (
	// Query
	QueryNew                     func(string, string, int) Query
	QueryDelete                  func(Query)
	QuerySetUsername             func(Query, string)
	QuerySetPassword             func(Query, string)
	QuerySetProxyhost           func(Query, string)
	QuerySetProxyport           func(Query, int)
	QuerySetProxyusername       func(Query, string)
	QuerySetProxypassword       func(Query, string)
	QueryLookupDiscid           func(Query, string) ReleaseList
	QueryLookupRelease          func(Query, string) Release
	QueryQuery                   func(Query, string, string, string, int, unsafe.Pointer, unsafe.Pointer) Metadata
	QueryAddCollectionEntries    func(Query, Collection, int, int) bool
	QueryDeleteCollectionEntries func(Query, Collection, int, int) bool
	QueryGetLastresult           func(Query) QueryResult
	QueryGetLasthttpcode         func(Query) int
	QueryGetLasterrormessage     func(Query, *byte, int) int
	QueryGetVersion              func(Query, *byte, int) int

	// Entity (Generic Ext Attributes/Elements)
	EntityExtAttributesSize func(unsafe.Pointer) int
	EntityExtAttributeName  func(unsafe.Pointer, int, *byte, int) int
	EntityExtAttributeValue func(unsafe.Pointer, int, *byte, int) int
	EntityExtElementsSize   func(unsafe.Pointer) int
	EntityExtElementName    func(unsafe.Pointer, int, *byte, int) int
	EntityExtElementValue   func(unsafe.Pointer, int, *byte, int) int

	// Alias
	AliasClone        func(Alias) Alias
	AliasDelete       func(Alias)
	AliasGetLocale    func(unsafe.Pointer, *byte, int) int
	AliasGetText      func(unsafe.Pointer, *byte, int) int
	AliasGetSortname  func(unsafe.Pointer, *byte, int) int
	AliasGetType      func(unsafe.Pointer, *byte, int) int
	AliasGetPrimary   func(unsafe.Pointer, *byte, int) int
	AliasGetBegindate func(unsafe.Pointer, *byte, int) int
	AliasGetEnddate   func(unsafe.Pointer, *byte, int) int
	AliasListSize     func(AliasList) int
	AliasListItem     func(AliasList, int) Alias
	AliasListClone    func(AliasList) AliasList
	AliasListDelete   func(AliasList)
	AliasListGetCount func(AliasList) int
	AliasListGetOffset func(AliasList) int

	// Annotation
	AnnotationClone     func(Annotation) Annotation
	AnnotationDelete    func(Annotation)
	AnnotationGetType   func(unsafe.Pointer, *byte, int) int
	AnnotationGetEntity func(unsafe.Pointer, *byte, int) int
	AnnotationGetName   func(unsafe.Pointer, *byte, int) int
	AnnotationGetText   func(unsafe.Pointer, *byte, int) int
	AnnotationListSize  func(AnnotationList) int
	AnnotationListItem  func(AnnotationList, int) Annotation
	AnnotationListClone func(AnnotationList) AnnotationList
	AnnotationListDelete func(AnnotationList)
	AnnotationListGetCount func(AnnotationList) int
	AnnotationListGetOffset func(AnnotationList) int

	// Artist
	ArtistClone               func(Artist) Artist
	ArtistDelete              func(Artist)
	ArtistGetID               func(unsafe.Pointer, *byte, int) int
	ArtistGetType             func(unsafe.Pointer, *byte, int) int
	ArtistGetName             func(unsafe.Pointer, *byte, int) int
	ArtistGetSortname         func(unsafe.Pointer, *byte, int) int
	ArtistGetGender           func(unsafe.Pointer, *byte, int) int
	ArtistGetCountry          func(unsafe.Pointer, *byte, int) int
	ArtistGetDisambiguation   func(unsafe.Pointer, *byte, int) int
	ArtistGetIpilist          func(Artist) IPIList
	ArtistGetLifespan         func(Artist) Lifespan
	ArtistGetAliaslist        func(Artist) AliasList
	ArtistGetRecordinglist    func(Artist) RecordingList
	ArtistGetReleaselist      func(Artist) ReleaseList
	ArtistGetReleasegrouplist func(Artist) ReleaseGroupList
	ArtistGetLabellist        func(Artist) LabelList
	ArtistGetWorklist         func(Artist) WorkList
	ArtistGetRelationlistlist func(Artist) RelationListList
	ArtistGetTaglist          func(Artist) TagList
	ArtistGetUsertaglist      func(Artist) UserTagList
	ArtistGetRating           func(Artist) Rating
	ArtistGetUserrating       func(Artist) UserRating
	ArtistListSize            func(ArtistList) int
	ArtistListItem            func(ArtistList, int) Artist
	ArtistListClone           func(ArtistList) ArtistList
	ArtistListDelete          func(ArtistList)
	ArtistListGetCount        func(ArtistList) int
	ArtistListGetOffset       func(ArtistList) int

	// ArtistCredit
	ArtistcreditClone             func(ArtistCredit) ArtistCredit
	ArtistcreditDelete            func(ArtistCredit)
	ArtistcreditGetNamecreditlist func(ArtistCredit) NameCreditList

	// Attribute
	AttributeClone         func(Attribute) Attribute
	AttributeDelete        func(Attribute)
	AttributeGetText       func(unsafe.Pointer, *byte, int) int
	AttributeListSize      func(AttributeList) int
	AttributeListItem      func(AttributeList, int) Attribute
	AttributeListClone     func(AttributeList) AttributeList
	AttributeListDelete    func(AttributeList)
	AttributeListGetCount  func(AttributeList) int
	AttributeListGetOffset func(AttributeList) int

	// CDStub
	CdstubClone        func(CDStub) CDStub
	CdstubDelete       func(CDStub)
	CdstubGetID        func(unsafe.Pointer, *byte, int) int
	CdstubGetTitle     func(unsafe.Pointer, *byte, int) int
	CdstubGetArtist    func(unsafe.Pointer, *byte, int) int
	CdstubGetBarcode   func(unsafe.Pointer, *byte, int) int
	CdstubGetComment   func(unsafe.Pointer, *byte, int) int
	CdstubGetNonmbtracklist func(CDStub) NonMBTrackList
	CdstubListSize     func(CDStubList) int
	CdstubListItem     func(CDStubList, int) CDStub
	CdstubListClone    func(CDStubList) CDStubList
	CdstubListDelete   func(CDStubList)
	CdstubListGetCount func(CDStubList) int
	CdstubListGetOffset func(CDStubList) int

	// Collection
	CollectionClone        func(Collection) Collection
	CollectionDelete       func(Collection)
	CollectionGetID        func(unsafe.Pointer, *byte, int) int
	CollectionGetName      func(unsafe.Pointer, *byte, int) int
	CollectionGetEditor    func(unsafe.Pointer, *byte, int) int
	CollectionGetReleaselist func(Collection) ReleaseList
	CollectionListSize     func(CollectionList) int
	CollectionListItem     func(CollectionList, int) Collection
	CollectionListClone    func(CollectionList) CollectionList
	CollectionListDelete   func(CollectionList)
	CollectionListGetCount func(CollectionList) int
	CollectionListGetOffset func(CollectionList) int

	// Disc
	DiscClone        func(Disc) Disc
	DiscDelete       func(Disc)
	DiscGetID        func(unsafe.Pointer, *byte, int) int
	DiscGetSectors   func(Disc) int
	DiscGetOffsetlist func(Disc) OffsetList
	DiscGetReleaselist func(Disc) ReleaseList
	DiscListSize     func(DiscList) int
	DiscListItem     func(DiscList, int) Disc
	DiscListClone    func(DiscList) DiscList
	DiscListDelete   func(DiscList)
	DiscListGetCount func(DiscList) int
	DiscListGetOffset func(DiscList) int

	// FreeDBDisc
	FreedbdiscClone        func(FreeDBDisc) FreeDBDisc
	FreedbdiscDelete       func(FreeDBDisc)
	FreedbdiscGetID        func(unsafe.Pointer, *byte, int) int
	FreedbdiscGetTitle     func(unsafe.Pointer, *byte, int) int
	FreedbdiscGetArtist    func(unsafe.Pointer, *byte, int) int
	FreedbdiscGetCategory  func(unsafe.Pointer, *byte, int) int
	FreedbdiscGetYear      func(unsafe.Pointer, *byte, int) int
	FreedbdiscGetNonmbtracklist func(FreeDBDisc) NonMBTrackList
	FreedbdiscListSize     func(FreeDBDiscList) int
	FreedbdiscListItem     func(FreeDBDiscList, int) FreeDBDisc
	FreedbdiscListClone    func(FreeDBDiscList) FreeDBDiscList
	FreedbdiscListDelete   func(FreeDBDiscList)
	FreedbdiscListGetCount func(FreeDBDiscList) int
	FreedbdiscListGetOffset func(FreeDBDiscList) int

	// IPI
	IpiClone        func(IPI) IPI
	IpiDelete       func(IPI)
	IpiGetIpi       func(unsafe.Pointer, *byte, int) int
	IpiListSize     func(IPIList) int
	IpiListItem     func(IPIList, int) IPI
	IpiListClone    func(IPIList) IPIList
	IpiListDelete   func(IPIList)
	IpiListGetCount func(IPIList) int
	IpiListGetOffset func(IPIList) int

	// ISRC
	ISRCClone        func(ISRC) ISRC
	ISRCDelete       func(ISRC)
	ISRCGetID        func(unsafe.Pointer, *byte, int) int
	ISRCListSize     func(ISRCList) int
	ISRCListItem     func(ISRCList, int) ISRC
	ISRCListClone    func(ISRCList) ISRCList
	ISRCListDelete   func(ISRCList)
	ISRCListGetCount func(ISRCList) int
	ISRCListGetOffset func(ISRCList) int

	// ISWC
	ISWCClone        func(ISWC) ISWC
	ISWCDelete       func(ISWC)
	ISWCGetIswc      func(unsafe.Pointer, *byte, int) int
	ISWCListSize     func(ISWCList) int
	ISWCListItem     func(ISWCList, int) ISWC
	ISWCListClone    func(ISWCList) ISWCList
	ISWCListDelete   func(ISWCList)
	ISWCListGetCount func(ISWCList) int
	ISWCListGetOffset func(ISWCList) int

	// Label
	LabelClone               func(Label) Label
	LabelDelete              func(Label)
	LabelGetID               func(unsafe.Pointer, *byte, int) int
	LabelGetType             func(unsafe.Pointer, *byte, int) int
	LabelGetName             func(unsafe.Pointer, *byte, int) int
	LabelGetSortname         func(unsafe.Pointer, *byte, int) int
	LabelGetLabelcode        func(unsafe.Pointer, *byte, int) int
	LabelGetIpilist          func(Label) IPIList
	LabelGetDisambiguation   func(unsafe.Pointer, *byte, int) int
	LabelGetCountry          func(unsafe.Pointer, *byte, int) int
	LabelGetLifespan         func(Label) Lifespan
	LabelGetAliaslist        func(Label) AliasList
	LabelGetReleaselist      func(Label) ReleaseList
	LabelGetRelationlistlist func(Label) RelationListList
	LabelGetTaglist          func(Label) TagList
	LabelGetUsertaglist      func(Label) UserTagList
	LabelGetRating           func(Label) Rating
	LabelGetUserrating       func(Label) UserRating
	LabelListSize            func(LabelList) int
	LabelListItem            func(LabelList, int) Label
	LabelListClone           func(LabelList) LabelList
	LabelListDelete          func(LabelList)
	LabelListGetCount        func(LabelList) int
	LabelListGetOffset       func(LabelList) int

	// LabelInfo
	LabelinfoClone           func(LabelInfo) LabelInfo
	LabelinfoDelete          func(LabelInfo)
	LabelinfoGetCatalognumber func(unsafe.Pointer, *byte, int) int
	LabelinfoGetLabel         func(LabelInfo) Label
	LabelinfoListSize        func(LabelInfoList) int
	LabelinfoListItem        func(LabelInfoList, int) LabelInfo
	LabelinfoListClone       func(LabelInfoList) LabelInfoList
	LabelinfoListDelete      func(LabelInfoList)
	LabelinfoListGetCount    func(LabelInfoList) int
	LabelinfoListGetOffset    func(LabelInfoList) int

	// Lifespan
	LifespanClone     func(Lifespan) Lifespan
	LifespanDelete    func(Lifespan)
	LifespanGetBegin  func(unsafe.Pointer, *byte, int) int
	LifespanGetEnd    func(unsafe.Pointer, *byte, int) int
	LifespanGetEnded  func(unsafe.Pointer, *byte, int) int

	// Medium
	MediumClone        func(Medium) Medium
	MediumDelete       func(Medium)
	MediumGetTitle     func(unsafe.Pointer, *byte, int) int
	MediumGetPosition  func(Medium) int
	MediumGetFormat    func(unsafe.Pointer, *byte, int) int
	MediumGetDisclist  func(Medium) DiscList
	MediumGetTracklist func(Medium) TrackList
	MediumContainsDiscid func(Medium, string) bool
	MediumListSize     func(MediumList) int
	MediumListItem     func(MediumList, int) Medium
	MediumListClone    func(MediumList) MediumList
	MediumListDelete   func(MediumList)
	MediumListGetCount func(MediumList) int
	MediumListGetOffset func(MediumList) int
	MediumListGetTrackcount func(MediumList) int

	// Message
	MessageClone  func(Message) Message
	MessageDelete func(Message)
	MessageGetText func(unsafe.Pointer, *byte, int) int

	// Metadata
	MetadataClone           func(Metadata) Metadata
	MetadataDelete          func(Metadata)
	MetadataGetArtist       func(Metadata) Artist
	MetadataGetRelease      func(Metadata) Release
	MetadataGetReleaseGroup func(Metadata) ReleaseGroup
	MetadataGetRecording    func(Metadata) Recording
	MetadataGetWork         func(Metadata) Work
	MetadataGetLabel        func(Metadata) Label
	MetadataGetDisc         func(Metadata) Disc
	MetadataGetPUID         func(Metadata) PUID
	MetadataGetISRC         func(Metadata) ISRC
	MetadataGetLabelinfolist func(Metadata) LabelInfoList
	MetadataGetRating       func(Metadata) Rating
	MetadataGetUserrating   func(Metadata) UserRating
	MetadataGetCollection    func(Metadata) Collection
	MetadataGetArtistlist    func(Metadata) ArtistList
	MetadataGetReleaselist   func(Metadata) ReleaseList
	MetadataGetReleaseGroupList func(Metadata) ReleaseGroupList
	MetadataGetRecordinglist func(Metadata) RecordingList
	MetadataGetLabellist     func(Metadata) LabelList
	MetadataGetWorklist      func(Metadata) WorkList
	MetadataGetISRCList      func(Metadata) ISRCList
	MetadataGetAnnotationlist func(Metadata) AnnotationList
	MetadataGetCdstublist    func(Metadata) CDStubList
	MetadataGetFreedbdisclist func(Metadata) FreeDBDiscList
	MetadataGetTaglist       func(Metadata) TagList
	MetadataGetUsertaglist   func(Metadata) UserTagList
	MetadataGetCollectionlist func(Metadata) CollectionList
	MetadataGetCdstub        func(Metadata) CDStub
	MetadataGetMessage       func(Metadata) Message

	// NonMBTrack
	NonmbtrackClone        func(NonMBTrack) NonMBTrack
	NonmbtrackDelete       func(NonMBTrack)
	NonmbtrackGetTitle     func(unsafe.Pointer, *byte, int) int
	NonmbtrackGetArtist    func(unsafe.Pointer, *byte, int) int
	NonmbtrackGetLength    func(NonMBTrack) int
	NonmbtrackListSize     func(NonMBTrackList) int
	NonmbtrackListItem     func(NonMBTrackList, int) NonMBTrack
	NonmbtrackListClone    func(NonMBTrackList) NonMBTrackList
	NonmbtrackListDelete   func(NonMBTrackList)
	NonmbtrackListGetCount func(NonMBTrackList) int
	NonmbtrackListGetOffset func(NonMBTrackList) int

	// NameCredit
	NamecreditClone        func(NameCredit) NameCredit
	NamecreditDelete       func(NameCredit)
	NamecreditGetJoinphrase func(unsafe.Pointer, *byte, int) int
	NamecreditGetName       func(unsafe.Pointer, *byte, int) int
	NamecreditGetArtist     func(NameCredit) Artist
	NamecreditListSize     func(NameCreditList) int
	NamecreditListItem     func(NameCreditList, int) NameCredit
	NamecreditListClone    func(NameCreditList) NameCreditList
	NamecreditListDelete   func(NameCreditList)
	NamecreditListGetCount func(NameCreditList) int
	NamecreditListGetOffset func(NameCreditList) int

	// Offset
	OffsetClone        func(Offset) Offset
	OffsetDelete       func(Offset)
	OffsetGetPosition  func(Offset) int
	OffsetGetOffset    func(Offset) int
	OffsetListSize     func(OffsetList) int
	OffsetListItem     func(OffsetList, int) Offset
	OffsetListClone    func(OffsetList) OffsetList
	OffsetListDelete   func(OffsetList)
	OffsetListGetCount func(OffsetList) int
	OffsetListGetOffset func(OffsetList) int

	// PUID
	PuidClone           func(PUID) PUID
	PuidDelete          func(PUID)
	PuidGetID           func(unsafe.Pointer, *byte, int) int
	PuidGetRecordinglist func(PUID) RecordingList
	PuidListSize        func(PUIDList) int
	PuidListItem        func(PUIDList, int) PUID
	PuidListClone       func(PUIDList) PUIDList
	PuidListDelete      func(PUIDList)
	PuidListGetCount    func(PUIDList) int
	PuidListGetOffset    func(PUIDList) int

	// Rating
	RatingClone        func(Rating) Rating
	RatingDelete       func(Rating)
	RatingGetVotescount func(Rating) int
	RatingGetRating    func(Rating) float64

	// Recording
	RecordingClone               func(Recording) Recording
	RecordingDelete              func(Recording)
	RecordingGetID               func(unsafe.Pointer, *byte, int) int
	RecordingGetTitle            func(unsafe.Pointer, *byte, int) int
	RecordingGetLength           func(Recording) int
	RecordingGetDisambiguation   func(unsafe.Pointer, *byte, int) int
	RecordingGetArtistcredit     func(Recording) ArtistCredit
	RecordingGetReleaselist      func(Recording) ReleaseList
	RecordingGetPuidlist         func(Recording) PUIDList
	RecordingGetISRCList      func(Recording) ISRCList
	RecordingGetRelationlistlist func(Recording) RelationListList
	RecordingGetTaglist          func(Recording) TagList
	RecordingGetUsertaglist      func(Recording) UserTagList
	RecordingGetRating           func(Recording) Rating
	RecordingGetUserrating       func(Recording) UserRating
	RecordingListSize            func(RecordingList) int
	RecordingListItem            func(RecordingList, int) Recording
	RecordingListClone           func(RecordingList) RecordingList
	RecordingListDelete          func(RecordingList)
	RecordingListGetCount        func(RecordingList) int
	RecordingListGetOffset        func(RecordingList) int

	// Relation
	RelationClone           func(Relation) Relation
	RelationDelete          func(Relation)
	RelationGetTarget       func(unsafe.Pointer, *byte, int) int
	RelationGetType         func(unsafe.Pointer, *byte, int) int
	RelationGetDirection    func(unsafe.Pointer, *byte, int) int
	RelationGetAttributelist func(Relation) AttributeList
	RelationGetBegin        func(unsafe.Pointer, *byte, int) int
	RelationGetEnd          func(unsafe.Pointer, *byte, int) int
	RelationGetEnded        func(unsafe.Pointer, *byte, int) int
	RelationGetArtist       func(Relation) Artist
	RelationGetRelease      func(Relation) Release
	RelationGetReleasegroup func(Relation) ReleaseGroup
	RelationGetRecording    func(Relation) Recording
	RelationGetLabel        func(Relation) Label
	RelationGetWork         func(Relation) Work

	// RelationList
	RelationListSize        func(RelationList) int
	RelationListItem        func(RelationList, int) Relation
	RelationListClone       func(RelationList) RelationList
	RelationListDelete      func(RelationList)
	RelationListGetTargettype func(unsafe.Pointer, *byte, int) int
	RelationListGetCount    func(RelationList) int
	RelationListGetOffset    func(RelationList) int

	// RelationListList
	RelationlistListSize     func(RelationListList) int
	RelationlistListItem     func(RelationListList, int) RelationList
	RelationlistListClone    func(RelationListList) RelationListList
	RelationlistListDelete   func(RelationListList)
	RelationlistListGetCount func(RelationListList) int
	RelationlistListGetOffset func(RelationListList) int

	// Release
	ReleaseClone                 func(Release) Release
	ReleaseDelete                func(Release)
	ReleaseGetID                 func(unsafe.Pointer, *byte, int) int
	ReleaseGetTitle              func(unsafe.Pointer, *byte, int) int
	ReleaseGetStatus             func(unsafe.Pointer, *byte, int) int
	ReleaseGetQuality            func(unsafe.Pointer, *byte, int) int
	ReleaseGetDisambiguation     func(unsafe.Pointer, *byte, int) int
	ReleaseGetPackaging          func(unsafe.Pointer, *byte, int) int
	ReleaseGetTextrepresentation func(Release) TextRepresentation
	ReleaseGetArtistcredit       func(Release) ArtistCredit
	ReleaseGetReleasegroup       func(Release) ReleaseGroup
	ReleaseGetDate               func(unsafe.Pointer, *byte, int) int
	ReleaseGetCountry            func(unsafe.Pointer, *byte, int) int
	ReleaseGetBarcode            func(unsafe.Pointer, *byte, int) int
	ReleaseGetAsin               func(unsafe.Pointer, *byte, int) int
	ReleaseGetLabelinfolist      func(Release) LabelInfoList
	ReleaseGetMediumlist         func(Release) MediumList
	ReleaseGetRelationlistlist   func(Release) RelationListList
	ReleaseGetCollectionlist     func(Release) CollectionList
	ReleaseMediaMatchingDiscid   func(Release, string) MediumList
	ReleaseListSize              func(ReleaseList) int
	ReleaseListItem              func(ReleaseList, int) Release
	ReleaseListClone             func(ReleaseList) ReleaseList
	ReleaseListDelete            func(ReleaseList)
	ReleaseListGetCount          func(ReleaseList) int
	ReleaseListGetOffset          func(ReleaseList) int

	// ReleaseGroup
	ReleasegroupClone               func(ReleaseGroup) ReleaseGroup
	ReleasegroupDelete              func(ReleaseGroup)
	ReleasegroupGetID               func(unsafe.Pointer, *byte, int) int
	ReleasegroupGetPrimarytype      func(unsafe.Pointer, *byte, int) int
	ReleasegroupGetFirstreleasedate func(unsafe.Pointer, *byte, int) int
	ReleasegroupGetTitle            func(unsafe.Pointer, *byte, int) int
	ReleasegroupGetDisambiguation   func(unsafe.Pointer, *byte, int) int
	ReleasegroupGetArtistcredit     func(ReleaseGroup) ArtistCredit
	ReleasegroupGetReleaselist      func(ReleaseGroup) ReleaseList
	ReleasegroupGetRelationlistlist func(ReleaseGroup) RelationListList
	ReleasegroupGetTaglist          func(ReleaseGroup) TagList
	ReleasegroupGetUsertaglist      func(ReleaseGroup) UserTagList
	ReleasegroupGetRating           func(ReleaseGroup) Rating
	ReleasegroupGetUserrating       func(ReleaseGroup) UserRating
	ReleasegroupGetSecondarytypelist func(ReleaseGroup) SecondaryTypeList
	ReleasegroupListSize            func(ReleaseGroupList) int
	ReleasegroupListItem            func(ReleaseGroupList, int) ReleaseGroup
	ReleasegroupListClone           func(ReleaseGroupList) ReleaseGroupList
	ReleasegroupListDelete          func(ReleaseGroupList)
	ReleasegroupListGetCount        func(ReleaseGroupList) int
	ReleasegroupListGetOffset        func(ReleaseGroupList) int

	// SecondaryType
	SecondarytypeListSize     func(SecondaryTypeList) int
	SecondarytypeListItem     func(SecondaryTypeList, int) SecondaryType
	SecondarytypeListClone    func(SecondaryTypeList) SecondaryTypeList
	SecondarytypeListDelete   func(SecondaryTypeList)
	SecondarytypeListGetCount func(SecondaryTypeList) int
	SecondarytypeListGetOffset func(SecondaryTypeList) int

	// Tag
	TagClone        func(Tag) Tag
	TagDelete       func(Tag)
	TagGetCount     func(Tag) int
	TagGetName      func(unsafe.Pointer, *byte, int) int
	TagListSize     func(TagList) int
	TagListItem     func(TagList, int) Tag
	TagListClone    func(TagList) TagList
	TagListDelete   func(TagList)
	TagListGetCount func(TagList) int
	TagListGetOffset func(TagList) int

	// TextRepresentation
	TextrepresentationClone         func(TextRepresentation) TextRepresentation
	TextrepresentationDelete        func(TextRepresentation)
	TextrepresentationGetLanguage   func(unsafe.Pointer, *byte, int) int
	TextrepresentationGetScript     func(unsafe.Pointer, *byte, int) int

	// Track
	TrackClone           func(Track) Track
	TrackDelete          func(Track)
	TrackGetPosition     func(Track) int
	TrackGetTitle        func(unsafe.Pointer, *byte, int) int
	TrackGetRecording    func(Track) Recording
	TrackGetLength       func(Track) int
	TrackGetArtistcredit func(Track) ArtistCredit
	TrackGetNumber       func(unsafe.Pointer, *byte, int) int
	TrackListSize        func(TrackList) int
	TrackListItem        func(TrackList, int) Track
	TrackListClone       func(TrackList) TrackList
	TrackListDelete      func(TrackList)
	TrackListGetCount    func(TrackList) int
	TrackListGetOffset    func(TrackList) int

	// UserRating
	UserratingClone         func(UserRating) UserRating
	UserratingDelete        func(UserRating)
	UserratingGetUserrating func(UserRating) int

	// UserTag
	UsertagClone        func(UserTag) UserTag
	UsertagDelete       func(UserTag)
	UsertagGetName      func(unsafe.Pointer, *byte, int) int
	UsertagListSize     func(UserTagList) int
	UsertagListItem     func(UserTagList, int) UserTag
	UsertagListClone    func(UserTagList) UserTagList
	UsertagListDelete   func(UserTagList)
	UsertagListGetCount func(UserTagList) int
	UsertagListGetOffset func(UserTagList) int

	// Work
	WorkClone               func(Work) Work
	WorkDelete              func(Work)
	WorkGetID               func(unsafe.Pointer, *byte, int) int
	WorkGetType             func(unsafe.Pointer, *byte, int) int
	WorkGetTitle            func(unsafe.Pointer, *byte, int) int
	WorkGetArtistcredit     func(Work) ArtistCredit
	WorkGetISWCList      func(Work) ISWCList
	WorkGetDisambiguation   func(unsafe.Pointer, *byte, int) int
	WorkGetAliaslist        func(Work) AliasList
	WorkGetRelationlistlist func(Work) RelationListList
	WorkGetTaglist          func(Work) TagList
	WorkGetUsertaglist      func(Work) UserTagList
	WorkGetRating           func(Work) Rating
	WorkGetUserrating       func(Work) UserRating
	WorkGetLanguage         func(unsafe.Pointer, *byte, int) int
	WorkListSize            func(WorkList) int
	WorkListItem            func(WorkList, int) Work
	WorkListClone           func(WorkList) WorkList
	WorkListDelete          func(WorkList)
	WorkListGetCount        func(WorkList) int
	WorkListGetOffset        func(WorkList) int
)

func Init() error {
	libmusicbrainz5, err := purego.Dlopen("libmusicbrainz5.so", purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		return fmt.Errorf("failed to load libmusicbrainz5.so: %w", err)
	}

	purego.RegisterLibFunc(&QueryNew, libmusicbrainz5, "mb5_query_new")
	purego.RegisterLibFunc(&QueryDelete, libmusicbrainz5, "mb5_query_delete")
	purego.RegisterLibFunc(&QuerySetUsername, libmusicbrainz5, "mb5_query_set_username")
	purego.RegisterLibFunc(&QuerySetPassword, libmusicbrainz5, "mb5_query_set_password")
	purego.RegisterLibFunc(&QuerySetProxyhost, libmusicbrainz5, "mb5_query_set_proxyhost")
	purego.RegisterLibFunc(&QuerySetProxyport, libmusicbrainz5, "mb5_query_set_proxyport")
	purego.RegisterLibFunc(&QuerySetProxyusername, libmusicbrainz5, "mb5_query_set_proxyusername")
	purego.RegisterLibFunc(&QuerySetProxypassword, libmusicbrainz5, "mb5_query_set_proxypassword")
	purego.RegisterLibFunc(&QueryLookupDiscid, libmusicbrainz5, "mb5_query_lookup_discid")
	purego.RegisterLibFunc(&QueryLookupRelease, libmusicbrainz5, "mb5_query_lookup_release")
	purego.RegisterLibFunc(&QueryQuery, libmusicbrainz5, "mb5_query_query")
	purego.RegisterLibFunc(&QueryAddCollectionEntries, libmusicbrainz5, "mb5_query_add_collection_entries")
	purego.RegisterLibFunc(&QueryDeleteCollectionEntries, libmusicbrainz5, "mb5_query_delete_collection_entries")
	purego.RegisterLibFunc(&QueryGetLastresult, libmusicbrainz5, "mb5_query_get_lastresult")
	purego.RegisterLibFunc(&QueryGetLasthttpcode, libmusicbrainz5, "mb5_query_get_lasthttpcode")
	purego.RegisterLibFunc(&QueryGetLasterrormessage, libmusicbrainz5, "mb5_query_get_lasterrormessage")
	purego.RegisterLibFunc(&QueryGetVersion, libmusicbrainz5, "mb5_query_get_version")

	// Entity
	purego.RegisterLibFunc(&EntityExtAttributesSize, libmusicbrainz5, "mb5_entity_ext_attributes_size")
	purego.RegisterLibFunc(&EntityExtAttributeName, libmusicbrainz5, "mb5_entity_ext_attribute_name")
	purego.RegisterLibFunc(&EntityExtAttributeValue, libmusicbrainz5, "mb5_entity_ext_attribute_value")
	purego.RegisterLibFunc(&EntityExtElementsSize, libmusicbrainz5, "mb5_entity_ext_elements_size")
	purego.RegisterLibFunc(&EntityExtElementName, libmusicbrainz5, "mb5_entity_ext_element_name")
	purego.RegisterLibFunc(&EntityExtElementValue, libmusicbrainz5, "mb5_entity_ext_element_value")

	// Alias
	purego.RegisterLibFunc(&AliasClone, libmusicbrainz5, "mb5_alias_clone")
	purego.RegisterLibFunc(&AliasDelete, libmusicbrainz5, "mb5_alias_delete")
	purego.RegisterLibFunc(&AliasGetLocale, libmusicbrainz5, "mb5_alias_get_locale")
	purego.RegisterLibFunc(&AliasGetText, libmusicbrainz5, "mb5_alias_get_text")
	purego.RegisterLibFunc(&AliasGetSortname, libmusicbrainz5, "mb5_alias_get_sortname")
	purego.RegisterLibFunc(&AliasGetType, libmusicbrainz5, "mb5_alias_get_type")
	purego.RegisterLibFunc(&AliasGetPrimary, libmusicbrainz5, "mb5_alias_get_primary")
	purego.RegisterLibFunc(&AliasGetBegindate, libmusicbrainz5, "mb5_alias_get_begindate")
	purego.RegisterLibFunc(&AliasGetEnddate, libmusicbrainz5, "mb5_alias_get_enddate")
	purego.RegisterLibFunc(&AliasListSize, libmusicbrainz5, "mb5_alias_list_size")
	purego.RegisterLibFunc(&AliasListItem, libmusicbrainz5, "mb5_alias_list_item")
	purego.RegisterLibFunc(&AliasListClone, libmusicbrainz5, "mb5_alias_list_clone")
	purego.RegisterLibFunc(&AliasListDelete, libmusicbrainz5, "mb5_alias_list_delete")
	purego.RegisterLibFunc(&AliasListGetCount, libmusicbrainz5, "mb5_alias_list_get_count")
	purego.RegisterLibFunc(&AliasListGetOffset, libmusicbrainz5, "mb5_alias_list_get_offset")

	// Annotation
	purego.RegisterLibFunc(&AnnotationClone, libmusicbrainz5, "mb5_annotation_clone")
	purego.RegisterLibFunc(&AnnotationDelete, libmusicbrainz5, "mb5_annotation_delete")
	purego.RegisterLibFunc(&AnnotationGetType, libmusicbrainz5, "mb5_annotation_get_type")
	purego.RegisterLibFunc(&AnnotationGetEntity, libmusicbrainz5, "mb5_annotation_get_entity")
	purego.RegisterLibFunc(&AnnotationGetName, libmusicbrainz5, "mb5_annotation_get_name")
	purego.RegisterLibFunc(&AnnotationGetText, libmusicbrainz5, "mb5_annotation_get_text")
	purego.RegisterLibFunc(&AnnotationListSize, libmusicbrainz5, "mb5_annotation_list_size")
	purego.RegisterLibFunc(&AnnotationListItem, libmusicbrainz5, "mb5_annotation_list_item")
	purego.RegisterLibFunc(&AnnotationListClone, libmusicbrainz5, "mb5_annotation_list_clone")
	purego.RegisterLibFunc(&AnnotationListDelete, libmusicbrainz5, "mb5_annotation_list_delete")
	purego.RegisterLibFunc(&AnnotationListGetCount, libmusicbrainz5, "mb5_annotation_list_get_count")
	purego.RegisterLibFunc(&AnnotationListGetOffset, libmusicbrainz5, "mb5_annotation_list_get_offset")

	// Artist
	purego.RegisterLibFunc(&ArtistClone, libmusicbrainz5, "mb5_artist_clone")
	purego.RegisterLibFunc(&ArtistDelete, libmusicbrainz5, "mb5_artist_delete")
	purego.RegisterLibFunc(&ArtistGetID, libmusicbrainz5, "mb5_artist_get_id")
	purego.RegisterLibFunc(&ArtistGetType, libmusicbrainz5, "mb5_artist_get_type")
	purego.RegisterLibFunc(&ArtistGetName, libmusicbrainz5, "mb5_artist_get_name")
	purego.RegisterLibFunc(&ArtistGetSortname, libmusicbrainz5, "mb5_artist_get_sortname")
	purego.RegisterLibFunc(&ArtistGetGender, libmusicbrainz5, "mb5_artist_get_gender")
	purego.RegisterLibFunc(&ArtistGetCountry, libmusicbrainz5, "mb5_artist_get_country")
	purego.RegisterLibFunc(&ArtistGetDisambiguation, libmusicbrainz5, "mb5_artist_get_disambiguation")
	purego.RegisterLibFunc(&ArtistGetIpilist, libmusicbrainz5, "mb5_artist_get_ipilist")
	purego.RegisterLibFunc(&ArtistGetLifespan, libmusicbrainz5, "mb5_artist_get_lifespan")
	purego.RegisterLibFunc(&ArtistGetAliaslist, libmusicbrainz5, "mb5_artist_get_aliaslist")
	purego.RegisterLibFunc(&ArtistGetRecordinglist, libmusicbrainz5, "mb5_artist_get_recordinglist")
	purego.RegisterLibFunc(&ArtistGetReleaselist, libmusicbrainz5, "mb5_artist_get_releaselist")
	purego.RegisterLibFunc(&ArtistGetReleasegrouplist, libmusicbrainz5, "mb5_artist_get_releasegrouplist")
	purego.RegisterLibFunc(&ArtistGetLabellist, libmusicbrainz5, "mb5_artist_get_labellist")
	purego.RegisterLibFunc(&ArtistGetWorklist, libmusicbrainz5, "mb5_artist_get_worklist")
	purego.RegisterLibFunc(&ArtistGetRelationlistlist, libmusicbrainz5, "mb5_artist_get_relationlistlist")
	purego.RegisterLibFunc(&ArtistGetTaglist, libmusicbrainz5, "mb5_artist_get_taglist")
	purego.RegisterLibFunc(&ArtistGetUsertaglist, libmusicbrainz5, "mb5_artist_get_usertaglist")
	purego.RegisterLibFunc(&ArtistGetRating, libmusicbrainz5, "mb5_artist_get_rating")
	purego.RegisterLibFunc(&ArtistGetUserrating, libmusicbrainz5, "mb5_artist_get_userrating")
	purego.RegisterLibFunc(&ArtistListSize, libmusicbrainz5, "mb5_artist_list_size")
	purego.RegisterLibFunc(&ArtistListItem, libmusicbrainz5, "mb5_artist_list_item")
	purego.RegisterLibFunc(&ArtistListClone, libmusicbrainz5, "mb5_artist_list_clone")
	purego.RegisterLibFunc(&ArtistListDelete, libmusicbrainz5, "mb5_artist_list_delete")
	purego.RegisterLibFunc(&ArtistListGetCount, libmusicbrainz5, "mb5_artist_list_get_count")
	purego.RegisterLibFunc(&ArtistListGetOffset, libmusicbrainz5, "mb5_artist_list_get_offset")

	// ArtistCredit
	purego.RegisterLibFunc(&ArtistcreditClone, libmusicbrainz5, "mb5_artistcredit_clone")
	purego.RegisterLibFunc(&ArtistcreditDelete, libmusicbrainz5, "mb5_artistcredit_delete")
	purego.RegisterLibFunc(&ArtistcreditGetNamecreditlist, libmusicbrainz5, "mb5_artistcredit_get_namecreditlist")

	// Attribute
	purego.RegisterLibFunc(&AttributeClone, libmusicbrainz5, "mb5_attribute_clone")
	purego.RegisterLibFunc(&AttributeDelete, libmusicbrainz5, "mb5_attribute_delete")
	purego.RegisterLibFunc(&AttributeGetText, libmusicbrainz5, "mb5_attribute_get_text")
	purego.RegisterLibFunc(&AttributeListSize, libmusicbrainz5, "mb5_attribute_list_size")
	purego.RegisterLibFunc(&AttributeListItem, libmusicbrainz5, "mb5_attribute_list_item")
	purego.RegisterLibFunc(&AttributeListClone, libmusicbrainz5, "mb5_attribute_list_clone")
	purego.RegisterLibFunc(&AttributeListDelete, libmusicbrainz5, "mb5_attribute_list_delete")
	purego.RegisterLibFunc(&AttributeListGetCount, libmusicbrainz5, "mb5_attribute_list_get_count")
	purego.RegisterLibFunc(&AttributeListGetOffset, libmusicbrainz5, "mb5_attribute_list_get_offset")

	// CDStub
	purego.RegisterLibFunc(&CdstubClone, libmusicbrainz5, "mb5_cdstub_clone")
	purego.RegisterLibFunc(&CdstubDelete, libmusicbrainz5, "mb5_cdstub_delete")
	purego.RegisterLibFunc(&CdstubGetID, libmusicbrainz5, "mb5_cdstub_get_id")
	purego.RegisterLibFunc(&CdstubGetTitle, libmusicbrainz5, "mb5_cdstub_get_title")
	purego.RegisterLibFunc(&CdstubGetArtist, libmusicbrainz5, "mb5_cdstub_get_artist")
	purego.RegisterLibFunc(&CdstubGetBarcode, libmusicbrainz5, "mb5_cdstub_get_barcode")
	purego.RegisterLibFunc(&CdstubGetComment, libmusicbrainz5, "mb5_cdstub_get_comment")
	purego.RegisterLibFunc(&CdstubGetNonmbtracklist, libmusicbrainz5, "mb5_cdstub_get_nonmbtracklist")
	purego.RegisterLibFunc(&CdstubListSize, libmusicbrainz5, "mb5_cdstub_list_size")
	purego.RegisterLibFunc(&CdstubListItem, libmusicbrainz5, "mb5_cdstub_list_item")
	purego.RegisterLibFunc(&CdstubListClone, libmusicbrainz5, "mb5_cdstub_list_clone")
	purego.RegisterLibFunc(&CdstubListDelete, libmusicbrainz5, "mb5_cdstub_list_delete")
	purego.RegisterLibFunc(&CdstubListGetCount, libmusicbrainz5, "mb5_cdstub_list_get_count")
	purego.RegisterLibFunc(&CdstubListGetOffset, libmusicbrainz5, "mb5_cdstub_list_get_offset")

	// Collection
	purego.RegisterLibFunc(&CollectionClone, libmusicbrainz5, "mb5_collection_clone")
	purego.RegisterLibFunc(&CollectionDelete, libmusicbrainz5, "mb5_collection_delete")
	purego.RegisterLibFunc(&CollectionGetID, libmusicbrainz5, "mb5_collection_get_id")
	purego.RegisterLibFunc(&CollectionGetName, libmusicbrainz5, "mb5_collection_get_name")
	purego.RegisterLibFunc(&CollectionGetEditor, libmusicbrainz5, "mb5_collection_get_editor")
	purego.RegisterLibFunc(&CollectionGetReleaselist, libmusicbrainz5, "mb5_collection_get_releaselist")
	purego.RegisterLibFunc(&CollectionListSize, libmusicbrainz5, "mb5_collection_list_size")
	purego.RegisterLibFunc(&CollectionListItem, libmusicbrainz5, "mb5_collection_list_item")
	purego.RegisterLibFunc(&CollectionListClone, libmusicbrainz5, "mb5_collection_list_clone")
	purego.RegisterLibFunc(&CollectionListDelete, libmusicbrainz5, "mb5_collection_list_delete")
	purego.RegisterLibFunc(&CollectionListGetCount, libmusicbrainz5, "mb5_collection_list_get_count")
	purego.RegisterLibFunc(&CollectionListGetOffset, libmusicbrainz5, "mb5_collection_list_get_offset")

	// Disc
	purego.RegisterLibFunc(&DiscClone, libmusicbrainz5, "mb5_disc_clone")
	purego.RegisterLibFunc(&DiscDelete, libmusicbrainz5, "mb5_disc_delete")
	purego.RegisterLibFunc(&DiscGetID, libmusicbrainz5, "mb5_disc_get_id")
	purego.RegisterLibFunc(&DiscGetSectors, libmusicbrainz5, "mb5_disc_get_sectors")
	purego.RegisterLibFunc(&DiscGetOffsetlist, libmusicbrainz5, "mb5_disc_get_offsetlist")
	purego.RegisterLibFunc(&DiscGetReleaselist, libmusicbrainz5, "mb5_disc_get_releaselist")
	purego.RegisterLibFunc(&DiscListSize, libmusicbrainz5, "mb5_disc_list_size")
	purego.RegisterLibFunc(&DiscListItem, libmusicbrainz5, "mb5_disc_list_item")
	purego.RegisterLibFunc(&DiscListClone, libmusicbrainz5, "mb5_disc_list_clone")
	purego.RegisterLibFunc(&DiscListDelete, libmusicbrainz5, "mb5_disc_list_delete")
	purego.RegisterLibFunc(&DiscListGetCount, libmusicbrainz5, "mb5_disc_list_get_count")
	purego.RegisterLibFunc(&DiscListGetOffset, libmusicbrainz5, "mb5_disc_list_get_offset")

	// FreeDBDisc
	purego.RegisterLibFunc(&FreedbdiscClone, libmusicbrainz5, "mb5_freedbdisc_clone")
	purego.RegisterLibFunc(&FreedbdiscDelete, libmusicbrainz5, "mb5_freedbdisc_delete")
	purego.RegisterLibFunc(&FreedbdiscGetID, libmusicbrainz5, "mb5_freedbdisc_get_id")
	purego.RegisterLibFunc(&FreedbdiscGetTitle, libmusicbrainz5, "mb5_freedbdisc_get_title")
	purego.RegisterLibFunc(&FreedbdiscGetArtist, libmusicbrainz5, "mb5_freedbdisc_get_artist")
	purego.RegisterLibFunc(&FreedbdiscGetCategory, libmusicbrainz5, "mb5_freedbdisc_get_category")
	purego.RegisterLibFunc(&FreedbdiscGetYear, libmusicbrainz5, "mb5_freedbdisc_get_year")
	purego.RegisterLibFunc(&FreedbdiscGetNonmbtracklist, libmusicbrainz5, "mb5_freedbdisc_get_nonmbtracklist")
	purego.RegisterLibFunc(&FreedbdiscListSize, libmusicbrainz5, "mb5_freedbdisc_list_size")
	purego.RegisterLibFunc(&FreedbdiscListItem, libmusicbrainz5, "mb5_freedbdisc_list_item")
	purego.RegisterLibFunc(&FreedbdiscListClone, libmusicbrainz5, "mb5_freedbdisc_list_clone")
	purego.RegisterLibFunc(&FreedbdiscListDelete, libmusicbrainz5, "mb5_freedbdisc_list_delete")
	purego.RegisterLibFunc(&FreedbdiscListGetCount, libmusicbrainz5, "mb5_freedbdisc_list_get_count")
	purego.RegisterLibFunc(&FreedbdiscListGetOffset, libmusicbrainz5, "mb5_freedbdisc_list_get_offset")

	// IPI
	purego.RegisterLibFunc(&IpiClone, libmusicbrainz5, "mb5_ipi_clone")
	purego.RegisterLibFunc(&IpiDelete, libmusicbrainz5, "mb5_ipi_delete")
	purego.RegisterLibFunc(&IpiGetIpi, libmusicbrainz5, "mb5_ipi_get_ipi")
	purego.RegisterLibFunc(&IpiListSize, libmusicbrainz5, "mb5_ipi_list_size")
	purego.RegisterLibFunc(&IpiListItem, libmusicbrainz5, "mb5_ipi_list_item")
	purego.RegisterLibFunc(&IpiListClone, libmusicbrainz5, "mb5_ipi_list_clone")
	purego.RegisterLibFunc(&IpiListDelete, libmusicbrainz5, "mb5_ipi_list_delete")
	purego.RegisterLibFunc(&IpiListGetCount, libmusicbrainz5, "mb5_ipi_list_get_count")
	purego.RegisterLibFunc(&IpiListGetOffset, libmusicbrainz5, "mb5_ipi_list_get_offset")

	// ISRC
	purego.RegisterLibFunc(&ISRCClone, libmusicbrainz5, "mb5_isrc_clone")
	purego.RegisterLibFunc(&ISRCDelete, libmusicbrainz5, "mb5_isrc_delete")
	purego.RegisterLibFunc(&ISRCGetID, libmusicbrainz5, "mb5_isrc_get_id")
	purego.RegisterLibFunc(&ISRCListSize, libmusicbrainz5, "mb5_isrc_list_size")
	purego.RegisterLibFunc(&ISRCListItem, libmusicbrainz5, "mb5_isrc_list_item")
	purego.RegisterLibFunc(&ISRCListClone, libmusicbrainz5, "mb5_isrc_list_clone")
	purego.RegisterLibFunc(&ISRCListDelete, libmusicbrainz5, "mb5_isrc_list_delete")
	purego.RegisterLibFunc(&ISRCListGetCount, libmusicbrainz5, "mb5_isrc_list_get_count")
	purego.RegisterLibFunc(&ISRCListGetOffset, libmusicbrainz5, "mb5_isrc_list_get_offset")

	// ISWC
	purego.RegisterLibFunc(&ISWCClone, libmusicbrainz5, "mb5_iswc_clone")
	purego.RegisterLibFunc(&ISWCDelete, libmusicbrainz5, "mb5_iswc_delete")
	purego.RegisterLibFunc(&ISWCGetIswc, libmusicbrainz5, "mb5_iswc_get_iswc")
	purego.RegisterLibFunc(&ISWCListSize, libmusicbrainz5, "mb5_iswc_list_size")
	purego.RegisterLibFunc(&ISWCListItem, libmusicbrainz5, "mb5_iswc_list_item")
	purego.RegisterLibFunc(&ISWCListClone, libmusicbrainz5, "mb5_iswc_list_clone")
	purego.RegisterLibFunc(&ISWCListDelete, libmusicbrainz5, "mb5_iswc_list_delete")
	purego.RegisterLibFunc(&ISWCListGetCount, libmusicbrainz5, "mb5_iswc_list_get_count")
	purego.RegisterLibFunc(&ISWCListGetOffset, libmusicbrainz5, "mb5_iswc_list_get_offset")

	// Label
	purego.RegisterLibFunc(&LabelClone, libmusicbrainz5, "mb5_label_clone")
	purego.RegisterLibFunc(&LabelDelete, libmusicbrainz5, "mb5_label_delete")
	purego.RegisterLibFunc(&LabelGetID, libmusicbrainz5, "mb5_label_get_id")
	purego.RegisterLibFunc(&LabelGetType, libmusicbrainz5, "mb5_label_get_type")
	purego.RegisterLibFunc(&LabelGetName, libmusicbrainz5, "mb5_label_get_name")
	purego.RegisterLibFunc(&LabelGetSortname, libmusicbrainz5, "mb5_label_get_sortname")
	purego.RegisterLibFunc(&LabelGetLabelcode, libmusicbrainz5, "mb5_label_get_labelcode")
	purego.RegisterLibFunc(&LabelGetIpilist, libmusicbrainz5, "mb5_label_get_ipilist")
	purego.RegisterLibFunc(&LabelGetDisambiguation, libmusicbrainz5, "mb5_label_get_disambiguation")
	purego.RegisterLibFunc(&LabelGetCountry, libmusicbrainz5, "mb5_label_get_country")
	purego.RegisterLibFunc(&LabelGetLifespan, libmusicbrainz5, "mb5_label_get_lifespan")
	purego.RegisterLibFunc(&LabelGetAliaslist, libmusicbrainz5, "mb5_label_get_aliaslist")
	purego.RegisterLibFunc(&LabelGetReleaselist, libmusicbrainz5, "mb5_label_get_releaselist")
	purego.RegisterLibFunc(&LabelGetRelationlistlist, libmusicbrainz5, "mb5_label_get_relationlistlist")
	purego.RegisterLibFunc(&LabelGetTaglist, libmusicbrainz5, "mb5_label_get_taglist")
	purego.RegisterLibFunc(&LabelGetUsertaglist, libmusicbrainz5, "mb5_label_get_usertaglist")
	purego.RegisterLibFunc(&LabelGetRating, libmusicbrainz5, "mb5_label_get_rating")
	purego.RegisterLibFunc(&LabelGetUserrating, libmusicbrainz5, "mb5_label_get_userrating")
	purego.RegisterLibFunc(&LabelListSize, libmusicbrainz5, "mb5_label_list_size")
	purego.RegisterLibFunc(&LabelListItem, libmusicbrainz5, "mb5_label_list_item")
	purego.RegisterLibFunc(&LabelListClone, libmusicbrainz5, "mb5_label_list_clone")
	purego.RegisterLibFunc(&LabelListDelete, libmusicbrainz5, "mb5_label_list_delete")
	purego.RegisterLibFunc(&LabelListGetCount, libmusicbrainz5, "mb5_label_list_get_count")
	purego.RegisterLibFunc(&LabelListGetOffset, libmusicbrainz5, "mb5_label_list_get_offset")

	// LabelInfo
	purego.RegisterLibFunc(&LabelinfoClone, libmusicbrainz5, "mb5_labelinfo_clone")
	purego.RegisterLibFunc(&LabelinfoDelete, libmusicbrainz5, "mb5_labelinfo_delete")
	purego.RegisterLibFunc(&LabelinfoGetCatalognumber, libmusicbrainz5, "mb5_labelinfo_get_catalognumber")
	purego.RegisterLibFunc(&LabelinfoGetLabel, libmusicbrainz5, "mb5_labelinfo_get_label")
	purego.RegisterLibFunc(&LabelinfoListSize, libmusicbrainz5, "mb5_labelinfo_list_size")
	purego.RegisterLibFunc(&LabelinfoListItem, libmusicbrainz5, "mb5_labelinfo_list_item")
	purego.RegisterLibFunc(&LabelinfoListClone, libmusicbrainz5, "mb5_labelinfo_list_clone")
	purego.RegisterLibFunc(&LabelinfoListDelete, libmusicbrainz5, "mb5_labelinfo_list_delete")
	purego.RegisterLibFunc(&LabelinfoListGetCount, libmusicbrainz5, "mb5_labelinfo_list_get_count")
	purego.RegisterLibFunc(&LabelinfoListGetOffset, libmusicbrainz5, "mb5_labelinfo_list_get_offset")

	// Lifespan
	purego.RegisterLibFunc(&LifespanClone, libmusicbrainz5, "mb5_lifespan_clone")
	purego.RegisterLibFunc(&LifespanDelete, libmusicbrainz5, "mb5_lifespan_delete")
	purego.RegisterLibFunc(&LifespanGetBegin, libmusicbrainz5, "mb5_lifespan_get_begin")
	purego.RegisterLibFunc(&LifespanGetEnd, libmusicbrainz5, "mb5_lifespan_get_end")
	purego.RegisterLibFunc(&LifespanGetEnded, libmusicbrainz5, "mb5_lifespan_get_ended")

	// Medium
	purego.RegisterLibFunc(&MediumClone, libmusicbrainz5, "mb5_medium_clone")
	purego.RegisterLibFunc(&MediumDelete, libmusicbrainz5, "mb5_medium_delete")
	purego.RegisterLibFunc(&MediumGetTitle, libmusicbrainz5, "mb5_medium_get_title")
	purego.RegisterLibFunc(&MediumGetPosition, libmusicbrainz5, "mb5_medium_get_position")
	purego.RegisterLibFunc(&MediumGetFormat, libmusicbrainz5, "mb5_medium_get_format")
	purego.RegisterLibFunc(&MediumGetDisclist, libmusicbrainz5, "mb5_medium_get_disclist")
	purego.RegisterLibFunc(&MediumGetTracklist, libmusicbrainz5, "mb5_medium_get_tracklist")
	purego.RegisterLibFunc(&MediumContainsDiscid, libmusicbrainz5, "mb5_medium_contains_discid")
	purego.RegisterLibFunc(&MediumListSize, libmusicbrainz5, "mb5_medium_list_size")
	purego.RegisterLibFunc(&MediumListItem, libmusicbrainz5, "mb5_medium_list_item")
	purego.RegisterLibFunc(&MediumListClone, libmusicbrainz5, "mb5_medium_list_clone")
	purego.RegisterLibFunc(&MediumListDelete, libmusicbrainz5, "mb5_medium_list_delete")
	purego.RegisterLibFunc(&MediumListGetCount, libmusicbrainz5, "mb5_medium_list_get_count")
	purego.RegisterLibFunc(&MediumListGetOffset, libmusicbrainz5, "mb5_medium_list_get_offset")
	purego.RegisterLibFunc(&MediumListGetTrackcount, libmusicbrainz5, "mb5_medium_list_get_trackcount")

	// Message
	purego.RegisterLibFunc(&MessageClone, libmusicbrainz5, "mb5_message_clone")
	purego.RegisterLibFunc(&MessageDelete, libmusicbrainz5, "mb5_message_delete")
	purego.RegisterLibFunc(&MessageGetText, libmusicbrainz5, "mb5_message_get_text")

	// Metadata
	purego.RegisterLibFunc(&MetadataClone, libmusicbrainz5, "mb5_metadata_clone")
	purego.RegisterLibFunc(&MetadataDelete, libmusicbrainz5, "mb5_metadata_delete")
	purego.RegisterLibFunc(&MetadataGetArtist, libmusicbrainz5, "mb5_metadata_get_artist")
	purego.RegisterLibFunc(&MetadataGetRelease, libmusicbrainz5, "mb5_metadata_get_release")
	purego.RegisterLibFunc(&MetadataGetReleaseGroup, libmusicbrainz5, "mb5_metadata_get_releasegroup")
	purego.RegisterLibFunc(&MetadataGetRecording, libmusicbrainz5, "mb5_metadata_get_recording")
	purego.RegisterLibFunc(&MetadataGetWork, libmusicbrainz5, "mb5_metadata_get_work")
	purego.RegisterLibFunc(&MetadataGetLabel, libmusicbrainz5, "mb5_metadata_get_label")
	purego.RegisterLibFunc(&MetadataGetDisc, libmusicbrainz5, "mb5_metadata_get_disc")
	purego.RegisterLibFunc(&MetadataGetPUID, libmusicbrainz5, "mb5_metadata_get_puid")
	purego.RegisterLibFunc(&MetadataGetISRC, libmusicbrainz5, "mb5_metadata_get_isrc")
	purego.RegisterLibFunc(&MetadataGetLabelinfolist, libmusicbrainz5, "mb5_metadata_get_labelinfolist")
	purego.RegisterLibFunc(&MetadataGetRating, libmusicbrainz5, "mb5_metadata_get_rating")
	purego.RegisterLibFunc(&MetadataGetUserrating, libmusicbrainz5, "mb5_metadata_get_userrating")
	purego.RegisterLibFunc(&MetadataGetCollection, libmusicbrainz5, "mb5_metadata_get_collection")
	purego.RegisterLibFunc(&MetadataGetArtistlist, libmusicbrainz5, "mb5_metadata_get_artistlist")
	purego.RegisterLibFunc(&MetadataGetReleaselist, libmusicbrainz5, "mb5_metadata_get_releaselist")
	purego.RegisterLibFunc(&MetadataGetReleaseGroupList, libmusicbrainz5, "mb5_metadata_get_releasegrouplist")
	purego.RegisterLibFunc(&MetadataGetRecordinglist, libmusicbrainz5, "mb5_metadata_get_recordinglist")
	purego.RegisterLibFunc(&MetadataGetLabellist, libmusicbrainz5, "mb5_metadata_get_labellist")
	purego.RegisterLibFunc(&MetadataGetWorklist, libmusicbrainz5, "mb5_metadata_get_worklist")
	purego.RegisterLibFunc(&MetadataGetISRCList, libmusicbrainz5, "mb5_metadata_get_isrclist")
	purego.RegisterLibFunc(&MetadataGetAnnotationlist, libmusicbrainz5, "mb5_metadata_get_annotationlist")
	purego.RegisterLibFunc(&MetadataGetCdstublist, libmusicbrainz5, "mb5_metadata_get_cdstublist")
	purego.RegisterLibFunc(&MetadataGetFreedbdisclist, libmusicbrainz5, "mb5_metadata_get_freedbdisclist")
	purego.RegisterLibFunc(&MetadataGetTaglist, libmusicbrainz5, "mb5_metadata_get_taglist")
	purego.RegisterLibFunc(&MetadataGetUsertaglist, libmusicbrainz5, "mb5_metadata_get_usertaglist")
	purego.RegisterLibFunc(&MetadataGetCollectionlist, libmusicbrainz5, "mb5_metadata_get_collectionlist")
	purego.RegisterLibFunc(&MetadataGetCdstub, libmusicbrainz5, "mb5_metadata_get_cdstub")
	purego.RegisterLibFunc(&MetadataGetMessage, libmusicbrainz5, "mb5_metadata_get_message")

	// NonMBTrack
	purego.RegisterLibFunc(&NonmbtrackClone, libmusicbrainz5, "mb5_nonmbtrack_clone")
	purego.RegisterLibFunc(&NonmbtrackDelete, libmusicbrainz5, "mb5_nonmbtrack_delete")
	purego.RegisterLibFunc(&NonmbtrackGetTitle, libmusicbrainz5, "mb5_nonmbtrack_get_title")
	purego.RegisterLibFunc(&NonmbtrackGetArtist, libmusicbrainz5, "mb5_nonmbtrack_get_artist")
	purego.RegisterLibFunc(&NonmbtrackGetLength, libmusicbrainz5, "mb5_nonmbtrack_get_length")
	purego.RegisterLibFunc(&NonmbtrackListSize, libmusicbrainz5, "mb5_nonmbtrack_list_size")
	purego.RegisterLibFunc(&NonmbtrackListItem, libmusicbrainz5, "mb5_nonmbtrack_list_item")
	purego.RegisterLibFunc(&NonmbtrackListClone, libmusicbrainz5, "mb5_nonmbtrack_list_clone")
	purego.RegisterLibFunc(&NonmbtrackListDelete, libmusicbrainz5, "mb5_nonmbtrack_list_delete")
	purego.RegisterLibFunc(&NonmbtrackListGetCount, libmusicbrainz5, "mb5_nonmbtrack_list_get_count")
	purego.RegisterLibFunc(&NonmbtrackListGetOffset, libmusicbrainz5, "mb5_nonmbtrack_list_get_offset")

	// NameCredit
	purego.RegisterLibFunc(&NamecreditClone, libmusicbrainz5, "mb5_namecredit_clone")
	purego.RegisterLibFunc(&NamecreditDelete, libmusicbrainz5, "mb5_namecredit_delete")
	purego.RegisterLibFunc(&NamecreditGetJoinphrase, libmusicbrainz5, "mb5_namecredit_get_joinphrase")
	purego.RegisterLibFunc(&NamecreditGetName, libmusicbrainz5, "mb5_namecredit_get_name")
	purego.RegisterLibFunc(&NamecreditGetArtist, libmusicbrainz5, "mb5_namecredit_get_artist")
	purego.RegisterLibFunc(&NamecreditListSize, libmusicbrainz5, "mb5_namecredit_list_size")
	purego.RegisterLibFunc(&NamecreditListItem, libmusicbrainz5, "mb5_namecredit_list_item")
	purego.RegisterLibFunc(&NamecreditListClone, libmusicbrainz5, "mb5_namecredit_list_clone")
	purego.RegisterLibFunc(&NamecreditListDelete, libmusicbrainz5, "mb5_namecredit_list_delete")
	purego.RegisterLibFunc(&NamecreditListGetCount, libmusicbrainz5, "mb5_namecredit_list_get_count")
	purego.RegisterLibFunc(&NamecreditListGetOffset, libmusicbrainz5, "mb5_namecredit_list_get_offset")

	// Offset
	purego.RegisterLibFunc(&OffsetClone, libmusicbrainz5, "mb5_offset_clone")
	purego.RegisterLibFunc(&OffsetDelete, libmusicbrainz5, "mb5_offset_delete")
	purego.RegisterLibFunc(&OffsetGetPosition, libmusicbrainz5, "mb5_offset_get_position")
	purego.RegisterLibFunc(&OffsetGetOffset, libmusicbrainz5, "mb5_offset_get_offset")
	purego.RegisterLibFunc(&OffsetListSize, libmusicbrainz5, "mb5_offset_list_size")
	purego.RegisterLibFunc(&OffsetListItem, libmusicbrainz5, "mb5_offset_list_item")
	purego.RegisterLibFunc(&OffsetListClone, libmusicbrainz5, "mb5_offset_list_clone")
	purego.RegisterLibFunc(&OffsetListDelete, libmusicbrainz5, "mb5_offset_list_delete")
	purego.RegisterLibFunc(&OffsetListGetCount, libmusicbrainz5, "mb5_offset_list_get_count")
	purego.RegisterLibFunc(&OffsetListGetOffset, libmusicbrainz5, "mb5_offset_list_get_offset")

	// PUID
	purego.RegisterLibFunc(&PuidClone, libmusicbrainz5, "mb5_puid_clone")
	purego.RegisterLibFunc(&PuidDelete, libmusicbrainz5, "mb5_puid_delete")
	purego.RegisterLibFunc(&PuidGetID, libmusicbrainz5, "mb5_puid_get_id")
	purego.RegisterLibFunc(&PuidGetRecordinglist, libmusicbrainz5, "mb5_puid_get_recordinglist")
	purego.RegisterLibFunc(&PuidListSize, libmusicbrainz5, "mb5_puid_list_size")
	purego.RegisterLibFunc(&PuidListItem, libmusicbrainz5, "mb5_puid_list_item")
	purego.RegisterLibFunc(&PuidListClone, libmusicbrainz5, "mb5_puid_list_clone")
	purego.RegisterLibFunc(&PuidListDelete, libmusicbrainz5, "mb5_puid_list_delete")
	purego.RegisterLibFunc(&PuidListGetCount, libmusicbrainz5, "mb5_puid_list_get_count")
	purego.RegisterLibFunc(&PuidListGetOffset, libmusicbrainz5, "mb5_puid_list_get_offset")

	// Rating
	purego.RegisterLibFunc(&RatingClone, libmusicbrainz5, "mb5_rating_clone")
	purego.RegisterLibFunc(&RatingDelete, libmusicbrainz5, "mb5_rating_delete")
	purego.RegisterLibFunc(&RatingGetVotescount, libmusicbrainz5, "mb5_rating_get_votescount")
	purego.RegisterLibFunc(&RatingGetRating, libmusicbrainz5, "mb5_rating_get_rating")

	// Recording
	purego.RegisterLibFunc(&RecordingClone, libmusicbrainz5, "mb5_recording_clone")
	purego.RegisterLibFunc(&RecordingDelete, libmusicbrainz5, "mb5_recording_delete")
	purego.RegisterLibFunc(&RecordingGetID, libmusicbrainz5, "mb5_recording_get_id")
	purego.RegisterLibFunc(&RecordingGetTitle, libmusicbrainz5, "mb5_recording_get_title")
	purego.RegisterLibFunc(&RecordingGetLength, libmusicbrainz5, "mb5_recording_get_length")
	purego.RegisterLibFunc(&RecordingGetDisambiguation, libmusicbrainz5, "mb5_recording_get_disambiguation")
	purego.RegisterLibFunc(&RecordingGetArtistcredit, libmusicbrainz5, "mb5_recording_get_artistcredit")
	purego.RegisterLibFunc(&RecordingGetReleaselist, libmusicbrainz5, "mb5_recording_get_releaselist")
	purego.RegisterLibFunc(&RecordingGetPuidlist, libmusicbrainz5, "mb5_recording_get_puidlist")
	purego.RegisterLibFunc(&RecordingGetISRCList, libmusicbrainz5, "mb5_recording_get_isrclist")
	purego.RegisterLibFunc(&RecordingGetRelationlistlist, libmusicbrainz5, "mb5_recording_get_relationlistlist")
	purego.RegisterLibFunc(&RecordingGetTaglist, libmusicbrainz5, "mb5_recording_get_taglist")
	purego.RegisterLibFunc(&RecordingGetUsertaglist, libmusicbrainz5, "mb5_recording_get_usertaglist")
	purego.RegisterLibFunc(&RecordingGetRating, libmusicbrainz5, "mb5_recording_get_rating")
	purego.RegisterLibFunc(&RecordingGetUserrating, libmusicbrainz5, "mb5_recording_get_userrating")
	purego.RegisterLibFunc(&RecordingListSize, libmusicbrainz5, "mb5_recording_list_size")
	purego.RegisterLibFunc(&RecordingListItem, libmusicbrainz5, "mb5_recording_list_item")
	purego.RegisterLibFunc(&RecordingListClone, libmusicbrainz5, "mb5_recording_list_clone")
	purego.RegisterLibFunc(&RecordingListDelete, libmusicbrainz5, "mb5_recording_list_delete")
	purego.RegisterLibFunc(&RecordingListGetCount, libmusicbrainz5, "mb5_recording_list_get_count")
	purego.RegisterLibFunc(&RecordingListGetOffset, libmusicbrainz5, "mb5_recording_list_get_offset")

	// Relation
	purego.RegisterLibFunc(&RelationClone, libmusicbrainz5, "mb5_relation_clone")
	purego.RegisterLibFunc(&RelationDelete, libmusicbrainz5, "mb5_relation_delete")
	purego.RegisterLibFunc(&RelationGetTarget, libmusicbrainz5, "mb5_relation_get_target")
	purego.RegisterLibFunc(&RelationGetType, libmusicbrainz5, "mb5_relation_get_type")
	purego.RegisterLibFunc(&RelationGetDirection, libmusicbrainz5, "mb5_relation_get_direction")
	purego.RegisterLibFunc(&RelationGetAttributelist, libmusicbrainz5, "mb5_relation_get_attributelist")
	purego.RegisterLibFunc(&RelationGetBegin, libmusicbrainz5, "mb5_relation_get_begin")
	purego.RegisterLibFunc(&RelationGetEnd, libmusicbrainz5, "mb5_relation_get_end")
	purego.RegisterLibFunc(&RelationGetEnded, libmusicbrainz5, "mb5_relation_get_ended")
	purego.RegisterLibFunc(&RelationGetArtist, libmusicbrainz5, "mb5_relation_get_artist")
	purego.RegisterLibFunc(&RelationGetRelease, libmusicbrainz5, "mb5_relation_get_release")
	purego.RegisterLibFunc(&RelationGetReleasegroup, libmusicbrainz5, "mb5_relation_get_releasegroup")
	purego.RegisterLibFunc(&RelationGetRecording, libmusicbrainz5, "mb5_relation_get_recording")
	purego.RegisterLibFunc(&RelationGetLabel, libmusicbrainz5, "mb5_relation_get_label")
	purego.RegisterLibFunc(&RelationGetWork, libmusicbrainz5, "mb5_relation_get_work")

	// RelationList
	purego.RegisterLibFunc(&RelationListSize, libmusicbrainz5, "mb5_relation_list_size")
	purego.RegisterLibFunc(&RelationListItem, libmusicbrainz5, "mb5_relation_list_item")
	purego.RegisterLibFunc(&RelationListClone, libmusicbrainz5, "mb5_relation_list_clone")
	purego.RegisterLibFunc(&RelationListDelete, libmusicbrainz5, "mb5_relation_list_delete")
	purego.RegisterLibFunc(&RelationListGetTargettype, libmusicbrainz5, "mb5_relation_list_get_targettype")
	purego.RegisterLibFunc(&RelationListGetCount, libmusicbrainz5, "mb5_relation_list_get_count")
	purego.RegisterLibFunc(&RelationListGetOffset, libmusicbrainz5, "mb5_relation_list_get_offset")

	// RelationListList
	purego.RegisterLibFunc(&RelationlistListSize, libmusicbrainz5, "mb5_relationlist_list_size")
	purego.RegisterLibFunc(&RelationlistListItem, libmusicbrainz5, "mb5_relationlist_list_item")
	purego.RegisterLibFunc(&RelationlistListClone, libmusicbrainz5, "mb5_relationlist_list_clone")
	purego.RegisterLibFunc(&RelationlistListDelete, libmusicbrainz5, "mb5_relationlist_list_delete")
	purego.RegisterLibFunc(&RelationlistListGetCount, libmusicbrainz5, "mb5_relationlist_list_get_count")
	purego.RegisterLibFunc(&RelationlistListGetOffset, libmusicbrainz5, "mb5_relationlist_list_get_offset")

	// Release
	purego.RegisterLibFunc(&ReleaseClone, libmusicbrainz5, "mb5_release_clone")
	purego.RegisterLibFunc(&ReleaseDelete, libmusicbrainz5, "mb5_release_delete")
	purego.RegisterLibFunc(&ReleaseGetID, libmusicbrainz5, "mb5_release_get_id")
	purego.RegisterLibFunc(&ReleaseGetTitle, libmusicbrainz5, "mb5_release_get_title")
	purego.RegisterLibFunc(&ReleaseGetStatus, libmusicbrainz5, "mb5_release_get_status")
	purego.RegisterLibFunc(&ReleaseGetQuality, libmusicbrainz5, "mb5_release_get_quality")
	purego.RegisterLibFunc(&ReleaseGetDisambiguation, libmusicbrainz5, "mb5_release_get_disambiguation")
	purego.RegisterLibFunc(&ReleaseGetPackaging, libmusicbrainz5, "mb5_release_get_packaging")
	purego.RegisterLibFunc(&ReleaseGetTextrepresentation, libmusicbrainz5, "mb5_release_get_textrepresentation")
	purego.RegisterLibFunc(&ReleaseGetArtistcredit, libmusicbrainz5, "mb5_release_get_artistcredit")
	purego.RegisterLibFunc(&ReleaseGetReleasegroup, libmusicbrainz5, "mb5_release_get_releasegroup")
	purego.RegisterLibFunc(&ReleaseGetDate, libmusicbrainz5, "mb5_release_get_date")
	purego.RegisterLibFunc(&ReleaseGetCountry, libmusicbrainz5, "mb5_release_get_country")
	purego.RegisterLibFunc(&ReleaseGetBarcode, libmusicbrainz5, "mb5_release_get_barcode")
	purego.RegisterLibFunc(&ReleaseGetAsin, libmusicbrainz5, "mb5_release_get_asin")
	purego.RegisterLibFunc(&ReleaseGetLabelinfolist, libmusicbrainz5, "mb5_release_get_labelinfolist")
	purego.RegisterLibFunc(&ReleaseGetMediumlist, libmusicbrainz5, "mb5_release_get_mediumlist")
	purego.RegisterLibFunc(&ReleaseGetRelationlistlist, libmusicbrainz5, "mb5_release_get_relationlistlist")
	purego.RegisterLibFunc(&ReleaseGetCollectionlist, libmusicbrainz5, "mb5_release_get_collectionlist")
	purego.RegisterLibFunc(&ReleaseMediaMatchingDiscid, libmusicbrainz5, "mb5_release_media_matching_discid")
	purego.RegisterLibFunc(&ReleaseListSize, libmusicbrainz5, "mb5_release_list_size")
	purego.RegisterLibFunc(&ReleaseListItem, libmusicbrainz5, "mb5_release_list_item")
	purego.RegisterLibFunc(&ReleaseListClone, libmusicbrainz5, "mb5_release_list_clone")
	purego.RegisterLibFunc(&ReleaseListDelete, libmusicbrainz5, "mb5_release_list_delete")
	purego.RegisterLibFunc(&ReleaseListGetCount, libmusicbrainz5, "mb5_release_list_get_count")
	purego.RegisterLibFunc(&ReleaseListGetOffset, libmusicbrainz5, "mb5_release_list_get_offset")

	// ReleaseGroup
	purego.RegisterLibFunc(&ReleasegroupClone, libmusicbrainz5, "mb5_releasegroup_clone")
	purego.RegisterLibFunc(&ReleasegroupDelete, libmusicbrainz5, "mb5_releasegroup_delete")
	purego.RegisterLibFunc(&ReleasegroupGetID, libmusicbrainz5, "mb5_releasegroup_get_id")
	purego.RegisterLibFunc(&ReleasegroupGetPrimarytype, libmusicbrainz5, "mb5_releasegroup_get_primarytype")
	purego.RegisterLibFunc(&ReleasegroupGetFirstreleasedate, libmusicbrainz5, "mb5_releasegroup_get_firstreleasedate")
	purego.RegisterLibFunc(&ReleasegroupGetTitle, libmusicbrainz5, "mb5_releasegroup_get_title")
	purego.RegisterLibFunc(&ReleasegroupGetDisambiguation, libmusicbrainz5, "mb5_releasegroup_get_disambiguation")
	purego.RegisterLibFunc(&ReleasegroupGetArtistcredit, libmusicbrainz5, "mb5_releasegroup_get_artistcredit")
	purego.RegisterLibFunc(&ReleasegroupGetReleaselist, libmusicbrainz5, "mb5_releasegroup_get_releaselist")
	purego.RegisterLibFunc(&ReleasegroupGetRelationlistlist, libmusicbrainz5, "mb5_releasegroup_get_relationlistlist")
	purego.RegisterLibFunc(&ReleasegroupGetTaglist, libmusicbrainz5, "mb5_releasegroup_get_taglist")
	purego.RegisterLibFunc(&ReleasegroupGetUsertaglist, libmusicbrainz5, "mb5_releasegroup_get_usertaglist")
	purego.RegisterLibFunc(&ReleasegroupGetRating, libmusicbrainz5, "mb5_releasegroup_get_rating")
	purego.RegisterLibFunc(&ReleasegroupGetUserrating, libmusicbrainz5, "mb5_releasegroup_get_userrating")
	purego.RegisterLibFunc(&ReleasegroupGetSecondarytypelist, libmusicbrainz5, "mb5_releasegroup_get_secondarytypelist")
	purego.RegisterLibFunc(&ReleasegroupListSize, libmusicbrainz5, "mb5_releasegroup_list_size")
	purego.RegisterLibFunc(&ReleasegroupListItem, libmusicbrainz5, "mb5_releasegroup_list_item")
	purego.RegisterLibFunc(&ReleasegroupListClone, libmusicbrainz5, "mb5_releasegroup_list_clone")
	purego.RegisterLibFunc(&ReleasegroupListDelete, libmusicbrainz5, "mb5_releasegroup_list_delete")
	purego.RegisterLibFunc(&ReleasegroupListGetCount, libmusicbrainz5, "mb5_releasegroup_list_get_count")
	purego.RegisterLibFunc(&ReleasegroupListGetOffset, libmusicbrainz5, "mb5_releasegroup_list_get_offset")

	// SecondaryType
	purego.RegisterLibFunc(&SecondarytypeListSize, libmusicbrainz5, "mb5_secondarytype_list_size")
	purego.RegisterLibFunc(&SecondarytypeListItem, libmusicbrainz5, "mb5_secondarytype_list_item")
	purego.RegisterLibFunc(&SecondarytypeListClone, libmusicbrainz5, "mb5_secondarytype_list_clone")
	purego.RegisterLibFunc(&SecondarytypeListDelete, libmusicbrainz5, "mb5_secondarytype_list_delete")
	purego.RegisterLibFunc(&SecondarytypeListGetCount, libmusicbrainz5, "mb5_secondarytype_list_get_count")
	purego.RegisterLibFunc(&SecondarytypeListGetOffset, libmusicbrainz5, "mb5_secondarytype_list_get_offset")

	// Tag
	purego.RegisterLibFunc(&TagClone, libmusicbrainz5, "mb5_tag_clone")
	purego.RegisterLibFunc(&TagDelete, libmusicbrainz5, "mb5_tag_delete")
	purego.RegisterLibFunc(&TagGetCount, libmusicbrainz5, "mb5_tag_get_count")
	purego.RegisterLibFunc(&TagGetName, libmusicbrainz5, "mb5_tag_get_name")
	purego.RegisterLibFunc(&TagListSize, libmusicbrainz5, "mb5_tag_list_size")
	purego.RegisterLibFunc(&TagListItem, libmusicbrainz5, "mb5_tag_list_item")
	purego.RegisterLibFunc(&TagListClone, libmusicbrainz5, "mb5_tag_list_clone")
	purego.RegisterLibFunc(&TagListDelete, libmusicbrainz5, "mb5_tag_list_delete")
	purego.RegisterLibFunc(&TagListGetCount, libmusicbrainz5, "mb5_tag_list_get_count")
	purego.RegisterLibFunc(&TagListGetOffset, libmusicbrainz5, "mb5_tag_list_get_offset")

	// TextRepresentation
	purego.RegisterLibFunc(&TextrepresentationClone, libmusicbrainz5, "mb5_textrepresentation_clone")
	purego.RegisterLibFunc(&TextrepresentationDelete, libmusicbrainz5, "mb5_textrepresentation_delete")
	purego.RegisterLibFunc(&TextrepresentationGetLanguage, libmusicbrainz5, "mb5_textrepresentation_get_language")
	purego.RegisterLibFunc(&TextrepresentationGetScript, libmusicbrainz5, "mb5_textrepresentation_get_script")

	// Track
	purego.RegisterLibFunc(&TrackClone, libmusicbrainz5, "mb5_track_clone")
	purego.RegisterLibFunc(&TrackDelete, libmusicbrainz5, "mb5_track_delete")
	purego.RegisterLibFunc(&TrackGetPosition, libmusicbrainz5, "mb5_track_get_position")
	purego.RegisterLibFunc(&TrackGetTitle, libmusicbrainz5, "mb5_track_get_title")
	purego.RegisterLibFunc(&TrackGetRecording, libmusicbrainz5, "mb5_track_get_recording")
	purego.RegisterLibFunc(&TrackGetLength, libmusicbrainz5, "mb5_track_get_length")
	purego.RegisterLibFunc(&TrackGetArtistcredit, libmusicbrainz5, "mb5_track_get_artistcredit")
	purego.RegisterLibFunc(&TrackGetNumber, libmusicbrainz5, "mb5_track_get_number")
	purego.RegisterLibFunc(&TrackListSize, libmusicbrainz5, "mb5_track_list_size")
	purego.RegisterLibFunc(&TrackListItem, libmusicbrainz5, "mb5_track_list_item")
	purego.RegisterLibFunc(&TrackListClone, libmusicbrainz5, "mb5_track_list_clone")
	purego.RegisterLibFunc(&TrackListDelete, libmusicbrainz5, "mb5_track_list_delete")
	purego.RegisterLibFunc(&TrackListGetCount, libmusicbrainz5, "mb5_track_list_get_count")
	purego.RegisterLibFunc(&TrackListGetOffset, libmusicbrainz5, "mb5_track_list_get_offset")

	// UserRating
	purego.RegisterLibFunc(&UserratingClone, libmusicbrainz5, "mb5_userrating_clone")
	purego.RegisterLibFunc(&UserratingDelete, libmusicbrainz5, "mb5_userrating_delete")
	purego.RegisterLibFunc(&UserratingGetUserrating, libmusicbrainz5, "mb5_userrating_get_userrating")

	// UserTag
	purego.RegisterLibFunc(&UsertagClone, libmusicbrainz5, "mb5_usertag_clone")
	purego.RegisterLibFunc(&UsertagDelete, libmusicbrainz5, "mb5_usertag_delete")
	purego.RegisterLibFunc(&UsertagGetName, libmusicbrainz5, "mb5_usertag_get_name")
	purego.RegisterLibFunc(&UsertagListSize, libmusicbrainz5, "mb5_usertag_list_size")
	purego.RegisterLibFunc(&UsertagListItem, libmusicbrainz5, "mb5_usertag_list_item")
	purego.RegisterLibFunc(&UsertagListClone, libmusicbrainz5, "mb5_usertag_list_clone")
	purego.RegisterLibFunc(&UsertagListDelete, libmusicbrainz5, "mb5_usertag_list_delete")
	purego.RegisterLibFunc(&UsertagListGetCount, libmusicbrainz5, "mb5_usertag_list_get_count")
	purego.RegisterLibFunc(&UsertagListGetOffset, libmusicbrainz5, "mb5_usertag_list_get_offset")

	// Work
	purego.RegisterLibFunc(&WorkClone, libmusicbrainz5, "mb5_work_clone")
	purego.RegisterLibFunc(&WorkDelete, libmusicbrainz5, "mb5_work_delete")
	purego.RegisterLibFunc(&WorkGetID, libmusicbrainz5, "mb5_work_get_id")
	purego.RegisterLibFunc(&WorkGetType, libmusicbrainz5, "mb5_work_get_type")
	purego.RegisterLibFunc(&WorkGetTitle, libmusicbrainz5, "mb5_work_get_title")
	purego.RegisterLibFunc(&WorkGetArtistcredit, libmusicbrainz5, "mb5_work_get_artistcredit")
	purego.RegisterLibFunc(&WorkGetISWCList, libmusicbrainz5, "mb5_work_get_iswclist")
	purego.RegisterLibFunc(&WorkGetDisambiguation, libmusicbrainz5, "mb5_work_get_disambiguation")
	purego.RegisterLibFunc(&WorkGetAliaslist, libmusicbrainz5, "mb5_work_get_aliaslist")
	purego.RegisterLibFunc(&WorkGetRelationlistlist, libmusicbrainz5, "mb5_work_get_relationlistlist")
	purego.RegisterLibFunc(&WorkGetTaglist, libmusicbrainz5, "mb5_work_get_taglist")
	purego.RegisterLibFunc(&WorkGetUsertaglist, libmusicbrainz5, "mb5_work_get_usertaglist")
	purego.RegisterLibFunc(&WorkGetRating, libmusicbrainz5, "mb5_work_get_rating")
	purego.RegisterLibFunc(&WorkGetUserrating, libmusicbrainz5, "mb5_work_get_userrating")
	purego.RegisterLibFunc(&WorkGetLanguage, libmusicbrainz5, "mb5_work_get_language")
	purego.RegisterLibFunc(&WorkListSize, libmusicbrainz5, "mb5_work_list_size")
	purego.RegisterLibFunc(&WorkListItem, libmusicbrainz5, "mb5_work_list_item")
	purego.RegisterLibFunc(&WorkListClone, libmusicbrainz5, "mb5_work_list_clone")
	purego.RegisterLibFunc(&WorkListDelete, libmusicbrainz5, "mb5_work_list_delete")
	purego.RegisterLibFunc(&WorkListGetCount, libmusicbrainz5, "mb5_work_list_get_count")
	purego.RegisterLibFunc(&WorkListGetOffset, libmusicbrainz5, "mb5_work_list_get_offset")

	return nil
}

// Go Helper Functions for MB5

func String(f func(p1 unsafe.Pointer, buf *byte, size int) int, p unsafe.Pointer) string {
	if p == nil {
		return ""
	}
	buf := make([]byte, 256)
	f(p, &buf[0], len(buf))
	return string(buf[:CStringLen(buf)])
}

func CStringLen(buf []byte) int {
	for i, b := range buf {
		if b == 0 {
			return i
		}
	}
	return len(buf)
}
