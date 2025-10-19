import { useCallback, useEffect, useState } from 'react';

import { GetExportsBuckets } from '@app';
import { usePayload, useViewConfig } from '@hooks';
import {
  Accordion,
  Group,
  Paper,
  Stack,
  Table,
  Text,
  Title,
} from '@mantine/core';
import { exports, types } from '@models';
import { LogError, useErrorHandler } from '@utils';

export const AssetChartsFacet = ({
  pageData: _pageData,
}: {
  pageData: exports.ExportsPage;
}) => {
  const [bucketsData, setBucketsData] = useState<types.Buckets | null>(null);
  const [loading, setLoading] = useState(true);
  const createPayload = usePayload();
  const { error, handleError, clearError } = useErrorHandler();
  const { config: viewConfig } = useViewConfig({ viewName: 'exports' });

  // Fetch buckets data
  const fetchBucketsData = useCallback(async () => {
    setLoading(true);
    clearError();
    try {
      const payload = createPayload(types.DataFacet.ASSETCHARTS);
      const result = await GetExportsBuckets(payload);
      setBucketsData(result);
    } catch (err: unknown) {
      handleError(err, 'Failed to fetch asset charts data');
    } finally {
      setLoading(false);
    }
  }, [createPayload, handleError, clearError]);

  useEffect(() => {
    fetchBucketsData();
  }, [fetchBucketsData]);

  // Parse series data to group by asset and metric
  const parseSeriesData = () => {
    if (!bucketsData?.series) return {};

    const grouped: Record<string, Record<string, types.Bucket[]>> = {};

    Object.entries(bucketsData.series).forEach(([seriesName, buckets]) => {
      // Parse dot notation: "0x1234567890ab_ETH.frequency" -> asset="0x1234567890ab_ETH", metric="frequency"
      const dotIndex = seriesName.lastIndexOf('.');
      if (dotIndex === -1) {
        LogError(`Invalid series name format: ${seriesName}`);
        return;
      }

      const assetKey = seriesName.substring(0, dotIndex);
      const metric = seriesName.substring(dotIndex + 1);

      if (!grouped[assetKey]) {
        grouped[assetKey] = {};
      }

      grouped[assetKey][metric] = buckets;
    });

    return grouped;
  };

  // Get current configuration from viewConfig
  const getCurrentConfig = () => {
    if (!viewConfig?.facets?.assetcharts?.facetChartConfig) {
      return {
        seriesStrategy: 'address+symbol',
        seriesPrefixLen: 12,
      };
    }

    const config = viewConfig.facets.assetcharts.facetChartConfig;
    return {
      seriesStrategy: config.seriesStrategy || 'address+symbol',
      seriesPrefixLen: config.seriesPrefixLen || 12,
    };
  };

  if (loading) {
    return (
      <Stack gap="md" p="xl" align="center" justify="center" h={400}>
        <Text size="lg" c="dimmed">
          Loading asset charts data...
        </Text>
      </Stack>
    );
  }

  if (error) {
    return (
      <Stack gap="md" p="xl" align="center" justify="center" h={400}>
        <Text size="lg" c="red">
          Error loading data: {error.message}
        </Text>
      </Stack>
    );
  }

  const groupedData = parseSeriesData();
  const config = getCurrentConfig();
  const assetCount = Object.keys(groupedData).length;
  const totalSeries = Object.keys(bucketsData?.series || {}).length;

  if (assetCount === 0) {
    return (
      <Stack gap="md" p="xl" align="center" justify="center" h={400}>
        <Text size="lg" c="dimmed">
          No asset chart data available
        </Text>
      </Stack>
    );
  }

  return (
    <Stack gap="md" p="md">
      {/* Configuration Info */}
      <Paper p="md" withBorder>
        <Group gap="xl">
          <div>
            <Text size="sm" c="dimmed">
              Strategy
            </Text>
            <Text fw={500}>{config.seriesStrategy}</Text>
          </div>
          <div>
            <Text size="sm" c="dimmed">
              Prefix Length
            </Text>
            <Text fw={500}>{config.seriesPrefixLen}</Text>
          </div>
          <div>
            <Text size="sm" c="dimmed">
              Assets
            </Text>
            <Text fw={500}>{assetCount}</Text>
          </div>
          <div>
            <Text size="sm" c="dimmed">
              Total Series
            </Text>
            <Text fw={500}>{totalSeries}</Text>
          </div>
        </Group>
      </Paper>

      {/* Asset Data */}
      <Accordion variant="contained">
        {Object.entries(groupedData).map(([assetKey, metrics]) => {
          const metricCount = Object.keys(metrics).length;
          const totalBuckets = Object.values(metrics).reduce(
            (sum, buckets) => sum + buckets.length,
            0,
          );

          return (
            <Accordion.Item key={assetKey} value={assetKey}>
              <Accordion.Control>
                <Group gap="md">
                  <Text fw={500}>{assetKey}</Text>
                  <Text size="sm" c="dimmed">
                    {metricCount} metrics, {totalBuckets} data points
                  </Text>
                </Group>
              </Accordion.Control>
              <Accordion.Panel>
                <Stack gap="md">
                  {Object.entries(metrics).map(([metric, buckets]) => (
                    <div key={metric}>
                      <Title order={5} mb="xs">
                        {metric} ({buckets.length} buckets)
                      </Title>
                      <Table striped>
                        <Table.Thead>
                          <Table.Tr>
                            <Table.Th>Date</Table.Th>
                            <Table.Th>Value</Table.Th>
                            <Table.Th>Blocks</Table.Th>
                          </Table.Tr>
                        </Table.Thead>
                        <Table.Tbody>
                          {buckets.slice(0, 10).map((bucket, index) => (
                            <Table.Tr key={index}>
                              <Table.Td>{bucket.bucketIndex}</Table.Td>
                              <Table.Td>
                                {bucket.total.toLocaleString()}
                              </Table.Td>
                              <Table.Td>
                                {bucket.startBlock}-{bucket.endBlock}
                              </Table.Td>
                            </Table.Tr>
                          ))}
                          {buckets.length > 10 && (
                            <Table.Tr>
                              <Table.Td colSpan={3}>
                                <Text size="sm" c="dimmed" ta="center">
                                  ... and {buckets.length - 10} more data points
                                </Text>
                              </Table.Td>
                            </Table.Tr>
                          )}
                        </Table.Tbody>
                      </Table>
                    </div>
                  ))}
                </Stack>
              </Accordion.Panel>
            </Accordion.Item>
          );
        })}
      </Accordion>
    </Stack>
  );
};
