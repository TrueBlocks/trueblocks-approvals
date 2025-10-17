// Copyright 2016, 2026 The Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were auto generated. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */

package contracts

import "github.com/TrueBlocks/trueblocks-approvals/pkg/types"

func (c *ContractsCollection) GetBuckets(payload *types.Payload) (*types.Buckets, error) {
	var facet types.BucketInterface

	switch payload.DataFacet {
	case ContractsDashboard:
		facet = c.dashboardFacet
	case ContractsExecute:
		facet = c.executeFacet
	case ContractsEvents:
		facet = c.eventsFacet
	default:
		return &types.Buckets{
			Series0:      []types.Bucket{},
			Series0Stats: types.BucketStats{},
			Series1:      []types.Bucket{},
			Series1Stats: types.BucketStats{},
			Series2:      []types.Bucket{},
			Series2Stats: types.BucketStats{},
			Series3:      []types.Bucket{},
			Series3Stats: types.BucketStats{},
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
