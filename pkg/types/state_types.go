package types

type StoreState string

const (
	StoreStateStale    StoreState = "stale"    // Needs refresh
	StoreStateFetching StoreState = "fetching" // Currently loading
	StoreStateLoaded   StoreState = "loaded"   // Complete data
	StoreStateError    StoreState = "error"    // Error occurred
)

var AllStoreStates = []struct {
	Value  StoreState `json:"value"`
	TSName string     `json:"tsname"`
}{
	{StoreStateStale, "STORE_STALE"},
	{StoreStateFetching, "STORE_FETCHING"},
	{StoreStateLoaded, "STORE_LOADED"},
	{StoreStateError, "STORE_ERROR"},
}
