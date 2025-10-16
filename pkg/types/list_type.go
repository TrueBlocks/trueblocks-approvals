package types

type StoreState int

const (
	StoreStateStale    StoreState = iota // Needs refresh
	StoreStateFetching                   // Currently loading
	StoreStateLoaded                     // Complete data
	StoreStateError                      // Error occurred
	StoreStateCanceled                   // User canceled
)

var AllStoreStates = []struct {
	Value  StoreState `json:"value"`
	TSName string     `json:"tsname"`
}{
	{StoreStateStale, "STALE"},
	{StoreStateFetching, "FETCHING"},
	{StoreStateLoaded, "LOADED"},
	{StoreStateError, "ERROR"},
	{StoreStateCanceled, "CANCELED"},
}

type FacetState string

const (
	FacetStateStale    FacetState = "stale"
	FacetStateFetching FacetState = "fetching"
	FacetStatePartial  FacetState = "partial"
	FacetStateLoaded   FacetState = "loaded"
	FacetStateError    FacetState = "error"
)

var AllFacetStates = []struct {
	Value  FacetState `json:"value"`
	TSName string     `json:"tsname"`
}{
	{FacetStateStale, "FACET_STALE"},
	{FacetStateFetching, "FACET_FETCHING"},
	{FacetStatePartial, "FACET_PARTIAL"},
	{FacetStateLoaded, "FACET_LOADED"},
	{FacetStateError, "FACET_ERROR"},
}
