package chunks

import "github.com/TrueBlocks/trueblocks-approvals/pkg/types"

func (c *ChunksCollection) updateStatsBucket(stats *Stats) {
	if stats == nil {
		return
	}

	c.statsFacet.UpdateBuckets(func(bucket *types.Buckets) {
		// For stats, use time-based daily buckets if range dates are available
		if stats.RangeDates != nil && stats.RangeDates.FirstDate != "" && stats.RangeDates.LastDate != "" {
			c.updateStatsBucketTimeBase(stats, bucket)
		} else {
			c.updateStatsBucketBlockBase(stats, bucket)
		}
	})
}

// updateStatsBucketTimeBase handles time-based daily bucketing for stats
func (c *ChunksCollection) updateStatsBucketTimeBase(stats *Stats, bucket *types.Buckets) {
	startBucket, err := parseDateToDailyBucket(stats.RangeDates.FirstDate)
	if err != nil {
		return
	}

	endBucket, err := parseDateToDailyBucket(stats.RangeDates.LastDate)
	if err != nil {
		return
	}

	// Find or create buckets for the date range
	c.ensureTimeBucketsExist(bucket, startBucket, endBucket)

	// Distribute the values across the time buckets
	c.distributeToTimeBuckets(bucket, startBucket, endBucket, stats)
}

// updateStatsBucketBlockBase handles traditional block-based bucketing (fallback)
func (c *ChunksCollection) updateStatsBucketBlockBase(stats *Stats, bucket *types.Buckets) {
	// Parse the range string to get block numbers
	firstBlock, lastBlock, err := parseRangeString(stats.Range)
	if err != nil {
		return
	}

	size := bucket.GridInfo.Size
	lastBucketIndex := int(lastBlock / size)

	// Define metrics and their values
	metrics := map[string]float64{
		"ratio":         float64(stats.Ratio),
		"appsPerBlock":  float64(stats.AppsPerBlock),
		"addrsPerBlock": float64(stats.AddrsPerBlock),
		"appsPerAddr":   float64(stats.AppsPerAddr),
	}

	// Process each metric using the flexible series structure
	maxBuckets := 0
	for seriesName, value := range metrics {
		bucket.EnsureSeriesExists(seriesName)
		series := bucket.GetSeries(seriesName)
		ensureBucketsExist(&series, lastBucketIndex, size)
		distributeToBuckets(&series, firstBlock, lastBlock, value, size)
		bucket.SetSeries(seriesName, series)

		if len(series) > maxBuckets {
			maxBuckets = len(series)
		}
	}

	// Maintain backwards compatibility by also updating legacy fields
	ensureBucketsExist(&bucket.Series0, lastBucketIndex, size)
	ensureBucketsExist(&bucket.Series1, lastBucketIndex, size)
	ensureBucketsExist(&bucket.Series2, lastBucketIndex, size)
	ensureBucketsExist(&bucket.Series3, lastBucketIndex, size)

	distributeToBuckets(&bucket.Series0, firstBlock, lastBlock, float64(stats.Ratio), size)
	distributeToBuckets(&bucket.Series1, firstBlock, lastBlock, float64(stats.AppsPerBlock), size)
	distributeToBuckets(&bucket.Series2, firstBlock, lastBlock, float64(stats.AddrsPerBlock), size)
	distributeToBuckets(&bucket.Series3, firstBlock, lastBlock, float64(stats.AppsPerAddr), size)

	// Update grid info
	legacyMaxBuckets := len(bucket.Series0)
	if len(bucket.Series1) > legacyMaxBuckets {
		legacyMaxBuckets = len(bucket.Series1)
	}
	if len(bucket.Series2) > legacyMaxBuckets {
		legacyMaxBuckets = len(bucket.Series2)
	}
	if len(bucket.Series3) > legacyMaxBuckets {
		legacyMaxBuckets = len(bucket.Series3)
	}

	finalMaxBuckets := maxBuckets
	if legacyMaxBuckets > finalMaxBuckets {
		finalMaxBuckets = legacyMaxBuckets
	}

	updateGridInfo(&bucket.GridInfo, finalMaxBuckets, lastBlock)
}
