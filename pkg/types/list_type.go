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

type LoadState string

const (
	FacetStateStale    LoadState = "stale"
	FacetStateFetching LoadState = "fetching"
	FacetStatePartial  LoadState = "partial"
	FacetStateLoaded   LoadState = "loaded"
	FacetStateError    LoadState = "error"
)

var AllFacetStates = []struct {
	Value  LoadState `json:"value"`
	TSName string    `json:"tsname"`
}{
	{FacetStateStale, "STALE"},
	{FacetStateFetching, "FETCHING"},
	{FacetStatePartial, "PARTIAL"},
	{FacetStateLoaded, "LOADED"},
	{FacetStateError, "ERROR"},
}
