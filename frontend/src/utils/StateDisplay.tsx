import { CSSProperties } from 'react';

import { Badge, Flex, Text } from '@mantine/core';
import { types } from '@models';

interface StateDisplayProps {
  state: types.StoreState;
  facetName: string;
  totalItems?: number;
  style?: CSSProperties;
}

const getStateColor = (state: types.StoreState) => {
  switch (state) {
    case types.StoreState.STORE_STALE:
      return 'gray';
    case types.StoreState.STORE_FETCHING:
      return 'blue';
    case types.StoreState.STORE_LOADED:
      return 'green';
    default:
      return 'gray';
  }
};

const getStateLabel = (state: types.StoreState) => {
  switch (state) {
    case types.StoreState.STORE_STALE:
      return 'Stale';
    case types.StoreState.STORE_FETCHING:
      return 'Fetching...';
    case types.StoreState.STORE_LOADED:
      return 'Loaded';
    default:
      return 'Unknown' + String(state);
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
