import { CSSProperties } from 'react';

import { Badge, Flex, Text } from '@mantine/core';
import { types } from '@models';

interface StateDisplayProps {
  state: types.LoadState;
  facetName: string;
  totalItems?: number;
  style?: CSSProperties;
}

const getStateColor = (state: types.LoadState) => {
  switch (state) {
    case types.LoadState.STALE:
      return 'gray';
    case types.LoadState.FETCHING:
      return 'blue';
    case types.LoadState.PARTIAL:
      return 'yellow';
    case types.LoadState.LOADED:
      return 'green';
    case types.LoadState.ERROR:
      return 'red';
    default:
      return 'gray';
  }
};

const getStateLabel = (state: types.LoadState) => {
  switch (state) {
    case types.LoadState.STALE:
      return 'Stale';
    case types.LoadState.FETCHING:
      return 'Fetching...';
    case types.LoadState.PARTIAL:
      return 'Partial';
    case types.LoadState.LOADED:
      return 'Loaded';
    case types.LoadState.ERROR:
      return 'Error';
    default:
      return 'Unknown';
  }
};

export const StateDisplay = ({
  state,
  facetName,
  totalItems,
  style,
}: StateDisplayProps) => {
  return (
    <Flex
      gap="sm"
      align="center"
      style={{
        padding: '8px 16px',
        borderBottom: '1px solid #e0e0e0',
        ...style,
      }}
    >
      <Text size="sm" c="dimmed">
        Facet: <strong>{facetName}</strong>
      </Text>
      <Badge color={getStateColor(state)} variant="light" size="sm">
        {getStateLabel(state)}
      </Badge>
      {totalItems !== undefined && (
        <Text size="sm" c="dimmed">
          Items: <strong>{totalItems}</strong>
        </Text>
      )}
    </Flex>
  );
};
