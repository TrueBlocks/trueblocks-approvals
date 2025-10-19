package types

type Buckets struct {
	// New flexible series structure using named metrics
	Series map[string][]Bucket `json:"series"`

	// Legacy fields for backwards compatibility - will be deprecated
	Series0 []Bucket `json:"series0"`
	Series1 []Bucket `json:"series1"`
	Series2 []Bucket `json:"series2"`
	Series3 []Bucket `json:"series3"`

	GridInfo GridInfo `json:"gridInfo"`
}

// NewBuckets creates a new Buckets struct with proper initialization
func NewBuckets() *Buckets {
	return &Buckets{
		Series:   make(map[string][]Bucket),
		Series0:  []Bucket{},
		Series1:  []Bucket{},
		Series2:  []Bucket{},
		Series3:  []Bucket{},
		GridInfo: NewGridInfo(),
	}
}

// NewBuckets creates a new Buckets struct with proper initialization
func NewBucketsWithGridInfo(gi *GridInfo) *Buckets {
	ret := NewBuckets()
	if gi != nil {
		ret.GridInfo = *gi
	}
	return ret
}

// GetSeries returns a series by name, creating it if it doesn't exist
func (b *Buckets) GetSeries(name string) []Bucket {
	if b.Series == nil {
		b.Series = make(map[string][]Bucket)
	}
	if series, exists := b.Series[name]; exists {
		return series
	}
	b.Series[name] = []Bucket{}
	return b.Series[name]
}

// SetSeries sets a series by name
func (b *Buckets) SetSeries(name string, buckets []Bucket) {
	if b.Series == nil {
		b.Series = make(map[string][]Bucket)
	}
	b.Series[name] = buckets
}

// EnsureSeriesExists ensures a series exists, creating it if necessary
func (b *Buckets) EnsureSeriesExists(name string) {
	if b.Series == nil {
		b.Series = make(map[string][]Bucket)
	}
	if _, exists := b.Series[name]; !exists {
		b.Series[name] = []Bucket{}
	}
}

// BucketInterface defines bucket operations that facets must implement
type BucketInterface interface {
	GetBuckets() *Buckets
	ClearBuckets()
	SetBuckets(buckets *Buckets)
	UpdateBuckets(updateFunc func(*Buckets))
}

type Bucket struct {
	BucketKey  string  `json:"bucketIndex"`
	StartBlock uint64  `json:"startBlock"`
	EndBlock   uint64  `json:"endBlock"`
	Total      float64 `json:"total"`
	ColorValue float64 `json:"colorValue"`
}

// NewBucket creates a new Bucket with the specified parameters
func NewBucket(bucketKey string, startBlock, endBlock uint64) Bucket {
	return Bucket{
		BucketKey:  bucketKey,
		StartBlock: startBlock,
		EndBlock:   endBlock,
		Total:      0,
		ColorValue: 0,
	}
}

type BucketStats struct {
	Total   float64 `json:"total"`
	Average float64 `json:"average"`
	Min     float64 `json:"min"`
	Max     float64 `json:"max"`
	Count   int     `json:"count"`
}

// NewBucketStats creates a new BucketStats with zero values
func NewBucketStats() BucketStats {
	return BucketStats{
		Total:   0,
		Average: 0,
		Min:     0,
		Max:     0,
		Count:   0,
	}
}

type GridInfo struct {
	Rows        int    `json:"rows"`
	Columns     int    `json:"columns"`
	MaxBlock    uint64 `json:"maxBlock"`
	Size        uint64 `json:"size"`
	BucketCount int    `json:"bucketCount"`
}

// NewGridInfo creates a new GridInfo with proper default values
func NewGridInfo() GridInfo {
	return GridInfo{
		Size:        100000,
		Rows:        0,
		Columns:     20,
		BucketCount: 0,
		MaxBlock:    0,
	}
}
