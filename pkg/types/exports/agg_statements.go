package exports

import (
	"fmt"
	"math/big"
	"time"

	"github.com/TrueBlocks/trueblocks-approvals/pkg/types"
)

// AssetCharts series strategy constants
const (
	AddressOnly       = "address"
	SymbolOnly        = "symbol"
	AddressWithSymbol = "address+symbol"
)

// generateAssetIdentifier creates collision-safe asset identifiers for series naming
func generateAssetIdentifier(asset, symbol string, config types.FacetChartConfig) string {
	switch config.SeriesStrategy {
	case AddressOnly:
		chars := config.SeriesPrefixLen
		if chars < 8 {
			chars = 8
		} // Minimum (risky but functional)
		if chars > 15 {
			chars = 15
		} // Six sigma safety limit
		return asset[:2+chars] // "0x" + chars

	case SymbolOnly:
		// Risk: symbol conflicts exist in real data
		return symbol

	case AddressWithSymbol:
		chars := config.SeriesPrefixLen
		if chars < 8 {
			chars = 8
		}
		if chars > 15 {
			chars = 15
		}
		return fmt.Sprintf("%s_%s", asset[:2+chars], symbol)

	default:
		return asset[:14] // Default: 12 chars (practical safety)
	}
}

// groupStatementsByAsset organizes statements by their asset address
func groupStatementsByAsset(statements []*Statement) map[string][]*Statement {
	assetGroups := make(map[string][]*Statement)

	for _, statement := range statements {
		assetAddr := statement.Asset.Hex()
		assetGroups[assetAddr] = append(assetGroups[assetAddr], statement)
	}

	return assetGroups
}

// calculateFrequencyBuckets creates time-bucketed transaction frequency data
func calculateFrequencyBuckets(statements []*Statement) []types.Bucket {
	// Group statements by day
	dayGroups := make(map[string][]*Statement)
	for _, stmt := range statements {
		day := timestampToDailyBucket(int64(stmt.Timestamp))
		dayGroups[day] = append(dayGroups[day], stmt)
	}

	// Create buckets for each day with data (sparse bucketing)
	var buckets []types.Bucket
	for day, dayStatements := range dayGroups {
		bucket := types.Bucket{
			BucketKey:  day,
			Total:      float64(len(dayStatements)), // Count of transactions
			StartBlock: 0,                           // AssetCharts uses time-based buckets, not block-based
			EndBlock:   0,
			ColorValue: 0,
		}
		buckets = append(buckets, bucket)
	}

	return buckets
}

// calculateVolumeBuckets creates time-bucketed volume data (amountIn + amountOut)
func calculateVolumeBuckets(statements []*Statement) []types.Bucket {
	// Group statements by day
	dayGroups := make(map[string][]*Statement)
	for _, stmt := range statements {
		day := timestampToDailyBucket(int64(stmt.Timestamp))
		dayGroups[day] = append(dayGroups[day], stmt)
	}

	// Create buckets for each day with data
	var buckets []types.Bucket
	for day, dayStatements := range dayGroups {
		totalVolume := 0.0
		decimals := 18 // Default to Ethereum standard

		for _, stmt := range dayStatements {
			if stmt.Decimals > 0 {
				decimals = int(stmt.Decimals)
			}
			amountIn := statementValueToFloat64(stmt.AmountIn, decimals)
			amountOut := statementValueToFloat64(stmt.AmountOut, decimals)
			totalVolume += amountIn + amountOut
		}

		bucket := types.Bucket{
			BucketKey:  day,
			Total:      totalVolume,
			StartBlock: 0,
			EndBlock:   0,
			ColorValue: 0,
		}
		buckets = append(buckets, bucket)
	}

	return buckets
}

// calculateGasOutBuckets creates time-bucketed gas consumption data
func calculateGasOutBuckets(statements []*Statement) []types.Bucket {
	// Group statements by day
	dayGroups := make(map[string][]*Statement)
	for _, stmt := range statements {
		day := timestampToDailyBucket(int64(stmt.Timestamp))
		dayGroups[day] = append(dayGroups[day], stmt)
	}

	// Create buckets for each day with data
	var buckets []types.Bucket
	for day, dayStatements := range dayGroups {
		totalGas := 0.0

		for _, stmt := range dayStatements {
			gasOut := statementValueToFloat64(stmt.GasOut, 18) // Gas is always in wei (18 decimals)
			totalGas += gasOut
		}

		bucket := types.Bucket{
			BucketKey:  day,
			Total:      totalGas,
			StartBlock: 0,
			EndBlock:   0,
			ColorValue: 0,
		}
		buckets = append(buckets, bucket)
	}

	return buckets
}

// calculateEndBalBuckets creates time-bucketed end balance data (last value per day)
func calculateEndBalBuckets(statements []*Statement) []types.Bucket {
	// Group statements by day
	dayGroups := make(map[string][]*Statement)
	for _, stmt := range statements {
		day := timestampToDailyBucket(int64(stmt.Timestamp))
		dayGroups[day] = append(dayGroups[day], stmt)
	}

	// Create buckets for each day with data
	var buckets []types.Bucket
	for day, dayStatements := range dayGroups {
		if len(dayStatements) == 0 {
			continue
		}

		// Use the last statement's end balance for the day
		lastStmt := dayStatements[len(dayStatements)-1]
		decimals := 18
		if lastStmt.Decimals > 0 {
			decimals = int(lastStmt.Decimals)
		}
		endBal := statementValueToFloat64(lastStmt.EndBal, decimals)

		bucket := types.Bucket{
			BucketKey:  day,
			Total:      endBal,
			StartBlock: 0,
			EndBlock:   0,
			ColorValue: 0,
		}
		buckets = append(buckets, bucket)
	}

	return buckets
}

// calculateNetAmountBuckets creates time-bucketed net amount data (amountIn - amountOut)
func calculateNetAmountBuckets(statements []*Statement) []types.Bucket {
	// Group statements by day
	dayGroups := make(map[string][]*Statement)
	for _, stmt := range statements {
		day := timestampToDailyBucket(int64(stmt.Timestamp))
		dayGroups[day] = append(dayGroups[day], stmt)
	}

	// Create buckets for each day with data
	var buckets []types.Bucket
	for day, dayStatements := range dayGroups {
		totalNet := 0.0
		decimals := 18 // Default to Ethereum standard

		for _, stmt := range dayStatements {
			if stmt.Decimals > 0 {
				decimals = int(stmt.Decimals)
			}
			amountIn := statementValueToFloat64(stmt.AmountIn, decimals)
			amountOut := statementValueToFloat64(stmt.AmountOut, decimals)
			totalNet += amountIn - amountOut
		}

		bucket := types.Bucket{
			BucketKey:  day,
			Total:      totalNet,
			StartBlock: 0,
			EndBlock:   0,
			ColorValue: 0,
		}
		buckets = append(buckets, bucket)
	}

	return buckets
}

// calculateNeighborsBuckets creates time-bucketed unique address count data
func calculateNeighborsBuckets(statements []*Statement) []types.Bucket {
	// Group statements by day
	dayGroups := make(map[string][]*Statement)
	for _, stmt := range statements {
		day := timestampToDailyBucket(int64(stmt.Timestamp))
		dayGroups[day] = append(dayGroups[day], stmt)
	}

	// Create buckets for each day with data
	var buckets []types.Bucket
	for day, dayStatements := range dayGroups {
		uniqueAddresses := make(map[string]bool)

		for _, stmt := range dayStatements {
			if !stmt.Sender.IsZero() {
				uniqueAddresses[stmt.Sender.Hex()] = true
			}
			if !stmt.Recipient.IsZero() {
				uniqueAddresses[stmt.Recipient.Hex()] = true
			}
		}

		bucket := types.Bucket{
			BucketKey:  day,
			Total:      float64(len(uniqueAddresses)),
			StartBlock: 0,
			EndBlock:   0,
			ColorValue: 0,
		}
		buckets = append(buckets, bucket)
	}

	return buckets
}

// statementValueToFloat64 converts Statement.any fields to float64 with decimal adjustment
func statementValueToFloat64(value any, decimals int) float64 {
	if value == nil {
		return 0.0
	}

	// Convert to string first (Go's any to string conversion)
	str := fmt.Sprintf("%v", value)
	if str == "" || str == "0" {
		return 0.0
	}

	// Parse as big integer (wei/base units)
	bigInt, ok := new(big.Int).SetString(str, 10)
	if !ok {
		return 0.0
	}

	// Convert to float64 with decimal adjustment
	scale := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	bigFloat := new(big.Float).SetInt(bigInt)
	scaleFloat := new(big.Float).SetInt(scale)
	result := new(big.Float).Quo(bigFloat, scaleFloat)

	float64Result, _ := result.Float64()
	return float64Result
}

// timestampToDailyBucket converts Unix timestamp to daily bucket identifier (YYYYMMDD)
func timestampToDailyBucket(timestamp int64) string {
	t := time.Unix(timestamp, 0).UTC()
	return fmt.Sprintf("%04d%02d%02d", t.Year(), t.Month(), t.Day())
}

// updateAssetChartsBucket processes Statement data and populates asset chart buckets with dot notation series
func (c *ExportsCollection) updateAssetChartsBucket(statement *Statement) {
	if statement == nil {
		return
	}

	c.assetchartsFacet.UpdateBuckets(func(buckets *types.Buckets) {
		// Get all statements from the facet's store for complete aggregation
		statements := c.assetchartsFacet.GetStore().GetItems()

		// Clear existing series and recalculate from all data
		buckets.Series = make(map[string][]types.Bucket)

		// Get the actual facet configuration
		var config types.FacetChartConfig
		if viewConfig, err := c.GetConfig(); err == nil {
			if facetConfig, exists := viewConfig.Facets["assetcharts"]; exists && facetConfig.FacetChartConfig != nil {
				config = *facetConfig.FacetChartConfig
			} else {
				// Fallback to defaults
				config = types.FacetChartConfig{
					SeriesStrategy:  AddressWithSymbol,
					SeriesPrefixLen: 12,
				}
			}
		} else {
			// Fallback if config unavailable
			config = types.FacetChartConfig{
				SeriesStrategy:  AddressWithSymbol,
				SeriesPrefixLen: 12,
			}
		}

		// Group statements by asset
		assetGroups := groupStatementsByAsset(statements)

		// Process each asset group and create series with dot notation
		for asset, assetStatements := range assetGroups {
			if len(assetStatements) == 0 {
				continue
			}

			// Generate asset identifier using our collision-safe strategy
			assetIdentifier := generateAssetIdentifier(asset, assetStatements[0].Symbol, config)

			// Calculate all metrics for this asset
			metrics := calculateAssetMetrics(assetStatements)

			// Create series for each metric using dot notation
			for metricName, bucketData := range metrics {
				seriesName := fmt.Sprintf("%s.%s", assetIdentifier, metricName)
				buckets.SetSeries(seriesName, bucketData)
			}
		}
	})
}

// calculateAssetMetrics calculates all available metrics for an asset's statements
func calculateAssetMetrics(statements []*Statement) map[string][]types.Bucket {
	return map[string][]types.Bucket{
		"frequency": calculateFrequencyBuckets(statements),
		"volume":    calculateVolumeBuckets(statements),
		"gasOut":    calculateGasOutBuckets(statements),
		"endBal":    calculateEndBalBuckets(statements),
		"netAmount": calculateNetAmountBuckets(statements),
		"neighbors": calculateNeighborsBuckets(statements),
	}
}
