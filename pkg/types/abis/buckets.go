// Copyright 2016, 2026 The Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were auto generated. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */

package abis

import "github.com/TrueBlocks/trueblocks-approvals/pkg/types"

func (c *AbisCollection) GetBuckets(payload *types.Payload) (*types.Buckets, error) {
	var facet types.BucketInterface

	switch payload.DataFacet {
	case AbisDownloaded:
		facet = c.downloadedFacet
	case AbisKnown:
		facet = c.knownFacet
	case AbisFunctions:
		facet = c.functionsFacet
	case AbisEvents:
		facet = c.eventsFacet
	default:
		return &types.Buckets{
			Series0: []types.Bucket{},
			Series1: []types.Bucket{},
			Series2: []types.Bucket{},
			Series3: []types.Bucket{},
			GridInfo: types.GridInfo{
				Size:        100000,
				Rows:        0,
				Columns:     20,
				BucketCount: 0,
				MaxBlock:    0,
			},
		}, nil
	}

	buckets := facet.GetBuckets()
	return buckets, nil
}
