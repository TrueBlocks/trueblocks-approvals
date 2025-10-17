package types

type StoreState string

const (
	StoreStateStale    StoreState = "stale"    // Needs refresh
	StoreStateFetching StoreState = "fetching" // Currently loading
	StoreStateLoaded   StoreState = "loaded"   // Complete data
)

var AllStoreStates = []struct {
	Value  StoreState `json:"value"`
	TSName string     `json:"tsname"`
}{
	{StoreStateStale, "STORE_STALE"},
	{StoreStateFetching, "STORE_FETCHING"},
	{StoreStateLoaded, "STORE_LOADED"},
}
