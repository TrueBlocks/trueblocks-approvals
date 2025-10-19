package chunks

import (
	"github.com/TrueBlocks/trueblocks-approvals/pkg/types"
)

func (c *ChunksCollection) updateIndexBucket(index *Index) {
	if index == nil {
		return
	}

	c.indexFacet.UpdateBuckets(func(bucket *types.Buckets) {
		// Parse the range string to get block numbers
		firstBlock, lastBlock, err := parseRangeString(index.Range)
		if err != nil {
			return
		}

		size := bucket.GridInfo.Size
		lastBucketIndex := int(lastBlock / size)

		// Define index metrics and their values
		metrics := map[string]float64{
			"nAddresses":   float64(index.NAddresses),
			"nAppearances": float64(index.NAppearances),
			"fileSize":     float64(index.FileSize),
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

		// Maintain backwards compatibility with legacy fields
		ensureBucketsExist(&bucket.Series0, lastBucketIndex, size)
		ensureBucketsExist(&bucket.Series1, lastBucketIndex, size)
		ensureBucketsExist(&bucket.Series2, lastBucketIndex, size)

		distributeToBuckets(&bucket.Series0, firstBlock, lastBlock, float64(index.NAddresses), size)
		distributeToBuckets(&bucket.Series1, firstBlock, lastBlock, float64(index.NAppearances), size)
		distributeToBuckets(&bucket.Series2, firstBlock, lastBlock, float64(index.FileSize), size)

		// Update grid info
		legacyMaxBuckets := len(bucket.Series0)
		if len(bucket.Series1) > legacyMaxBuckets {
			legacyMaxBuckets = len(bucket.Series1)
		}
		if len(bucket.Series2) > legacyMaxBuckets {
			legacyMaxBuckets = len(bucket.Series2)
		}

		finalMaxBuckets := maxBuckets
		if legacyMaxBuckets > finalMaxBuckets {
			finalMaxBuckets = legacyMaxBuckets
		}

		updateGridInfo(&bucket.GridInfo, finalMaxBuckets, lastBlock)
	})
}
