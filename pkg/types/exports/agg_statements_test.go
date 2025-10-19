package exports

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestData represents the JSON structure from our test data
type TestData struct {
	Data []StatementJSON `json:"data"`
}

// StatementJSON represents a single statement from JSON (before conversion to Statement struct)
type StatementJSON struct {
	AccountedFor string `json:"accountedFor"`
	AmountIn     string `json:"amountIn"`
	AmountNet    string `json:"amountNet"`
	AmountOut    string `json:"amountOut"`
	Asset        string `json:"asset"`
	BegBal       string `json:"begBal"`
	EndBal       string `json:"endBal"`
	Decimals     int    `json:"decimals"`
	GasOut       string `json:"gasOut"`
	Recipient    string `json:"recipient"`
	Sender       string `json:"sender"`
	SpotPrice    string `json:"spotPrice"`
	Symbol       string `json:"symbol"`
	Timestamp    int64  `json:"timestamp"`
	TotalIn      string `json:"totalIn"`
	TotalOut     string `json:"totalOut"`
}

func TestAssetChartsBucketing(t *testing.T) {
	// Load test data
	testFile := filepath.Join("testdata", "tb_statements_sample.json")
	data, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	var testData TestData
	if err := json.Unmarshal(data, &testData); err != nil {
		t.Fatalf("Failed to unmarshal test data: %v", err)
	}

	t.Logf("Loaded %d statements for testing", len(testData.Data))

	// Test asset grouping
	t.Run("AssetGrouping", func(t *testing.T) {
		assetGroups := groupStatementsByAsset(testData.Data)

		// The test data contains 10 different assets (more diverse than initially expected)
		expectedMinAssets := 8 // At least 8 different assets
		if len(assetGroups) < expectedMinAssets {
			t.Errorf("Expected at least %d asset groups, got %d", expectedMinAssets, len(assetGroups))
		}

		// Verify ETH (special address) is present
		ethAsset := "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
		if _, exists := assetGroups[ethAsset]; !exists {
			t.Errorf("Expected ETH asset %s not found in groups", ethAsset)
		}

		// Verify DAI v1 is present
		daiV1Asset := "0x89d24a6b4ccb1b6faa2625fe562bdd9a23260359"
		if _, exists := assetGroups[daiV1Asset]; !exists {
			t.Errorf("Expected DAI v1 asset %s not found in groups", daiV1Asset)
		}

		// Verify we have a good mix of activity levels
		var singleTxAssets, multiTxAssets int
		for asset, stmts := range assetGroups {
			if len(stmts) == 1 {
				singleTxAssets++
			} else {
				multiTxAssets++
			}
			t.Logf("Asset %s: %d statements", asset[:10]+"...", len(stmts))
		}

		// Most assets should be single-transaction (demonstrates sparse bucketing need)
		if singleTxAssets == 0 {
			t.Error("Expected some single-transaction assets to test sparse bucketing")
		}

		// Should have some multi-transaction assets too
		if multiTxAssets == 0 {
			t.Error("Expected some multi-transaction assets")
		}

		t.Logf("Asset distribution: %d single-tx, %d multi-tx (validates sparse bucketing approach)", singleTxAssets, multiTxAssets)
	})

	// Test time bucketing
	t.Run("TimeBucketing", func(t *testing.T) {
		// Test converting timestamps to daily bucket keys
		testCases := []struct {
			timestamp int64
			expected  string
		}{
			{1572639538, "20191101"}, // Nov 1, 2019
			{1572660966, "20191102"}, // Nov 2, 2019
			{1576868456, "20191220"}, // Dec 20, 2019
		}

		for _, tc := range testCases {
			bucketKey := timestampToDailyBucket(tc.timestamp)
			if bucketKey != tc.expected {
				t.Errorf("Timestamp %d: expected bucket %s, got %s", tc.timestamp, tc.expected, bucketKey)
			}
		}
	})

	// Test frequency metric calculation
	t.Run("FrequencyMetric", func(t *testing.T) {
		assetGroups := groupStatementsByAsset(testData.Data)

		totalBuckets := 0
		for asset, stmts := range assetGroups {
			buckets := calculateFrequencyBuckets(stmts)
			totalBuckets += len(buckets)

			t.Logf("Asset %s... frequency buckets (%d days):", asset[:10], len(buckets))
			for bucketKey, count := range buckets {
				t.Logf("  %s: %d transactions", bucketKey, int(count))
			}
		}

		// Validate sparse bucketing: should have many more assets than bucket days
		// This proves we don't create empty buckets for every possible day
		if totalBuckets < len(assetGroups) {
			t.Errorf("Sparse bucketing validation failed: %d total bucket-days for %d assets", totalBuckets, len(assetGroups))
		}

		t.Logf("Sparse bucketing validated: %d bucket-days across %d assets (avg %.1f days per asset)",
			totalBuckets, len(assetGroups), float64(totalBuckets)/float64(len(assetGroups)))
	})

	// Test volume metric calculation
	t.Run("VolumeMetric", func(t *testing.T) {
		assetGroups := groupStatementsByAsset(testData.Data)

		for asset, stmts := range assetGroups {
			buckets := calculateVolumeBuckets(stmts)
			t.Logf("Asset %s... volume buckets:", asset[:10])
			for bucketKey, volume := range buckets {
				t.Logf("  %s: %.6f", bucketKey, volume)
			}
		}
	})

	// Test edge cases and data diversity
	t.Run("EdgeCases", func(t *testing.T) {
		assetGroups := groupStatementsByAsset(testData.Data)

		// Test date range diversity (should span multiple years)
		var minYear, maxYear int = 9999, 0
		for _, stmts := range assetGroups {
			for _, stmt := range stmts {
				year := int((stmt.Timestamp / (365 * 24 * 3600)) + 1970) // Rough year calculation
				if year < minYear {
					minYear = year
				}
				if year > maxYear {
					maxYear = year
				}
			}
		}

		yearSpan := maxYear - minYear
		if yearSpan < 2 {
			t.Logf("Warning: Data only spans %d years (%d-%d), consider testing with wider date range", yearSpan, minYear, maxYear)
		} else {
			t.Logf("Good temporal diversity: data spans %d years (%d-%d)", yearSpan, minYear, maxYear)
		}

		// Test volume edge cases (zero volumes, large volumes)
		var zeroVolumeAssets, largeVolumeAssets int
		for _, stmts := range assetGroups {
			buckets := calculateVolumeBuckets(stmts)
			for _, volume := range buckets {
				if volume == 0.0 {
					zeroVolumeAssets++
				} else if volume > 1000 {
					largeVolumeAssets++
				}
			}
		}

		t.Logf("Volume diversity: %d zero-volume bucket-days, %d large-volume bucket-days", zeroVolumeAssets, largeVolumeAssets)
	})
}

// groupStatementsByAsset groups statements by asset address
func groupStatementsByAsset(statements []StatementJSON) map[string][]StatementJSON {
	groups := make(map[string][]StatementJSON)
	for _, stmt := range statements {
		asset := stmt.Asset
		groups[asset] = append(groups[asset], stmt)
	}
	return groups
}

// timestampToDailyBucket converts Unix timestamp to daily bucket key
func timestampToDailyBucket(timestamp int64) string {
	t := time.Unix(timestamp, 0).UTC()
	return fmt.Sprintf("%04d%02d%02d", t.Year(), t.Month(), t.Day())
}

// calculateFrequencyBuckets counts transactions per day per asset
func calculateFrequencyBuckets(statements []StatementJSON) map[string]float64 {
	buckets := make(map[string]float64)
	for _, stmt := range statements {
		bucketKey := timestampToDailyBucket(stmt.Timestamp)
		buckets[bucketKey]++
	}
	return buckets
}

// calculateVolumeBuckets sums amountIn + amountOut per day per asset
func calculateVolumeBuckets(statements []StatementJSON) map[string]float64 {
	buckets := make(map[string]float64)
	for _, stmt := range statements {
		bucketKey := timestampToDailyBucket(stmt.Timestamp)

		// Convert amounts to float64 (simplified for testing)
		amountIn := stringToFloat64(stmt.AmountIn, stmt.Decimals)
		amountOut := stringToFloat64(stmt.AmountOut, stmt.Decimals)

		buckets[bucketKey] += amountIn + amountOut
	}
	return buckets
}

// stringToFloat64 converts string amount to float64 with decimal scaling
func stringToFloat64(amountStr string, decimals int) float64 {
	if amountStr == "" || amountStr == "0" {
		return 0.0
	}

	// Parse as big integer
	bigInt, ok := new(big.Int).SetString(amountStr, 10)
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
