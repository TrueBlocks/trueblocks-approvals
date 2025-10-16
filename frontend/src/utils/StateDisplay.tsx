import { CSSProperties } from 'react';

import { Badge, Flex, Text } from '@mantine/core';
import { types } from '@models';

interface StateDisplayProps {
  state: types.FacetState;
  facetName: string;
  totalItems?: number;
  style?: CSSProperties;
}

const getStateColor = (state: types.FacetState) => {
  switch (state) {
    case types.FacetState.FACET_STALE:
      return 'gray';
    case types.FacetState.FACET_FETCHING:
      return 'blue';
    case types.FacetState.FACET_PARTIAL:
      return 'yellow';
    case types.FacetState.FACET_LOADED:
      return 'green';
    case types.FacetState.FACET_ERROR:
      return 'red';
    default:
      return 'gray';
  }
};

const getStateLabel = (state: types.FacetState) => {
  switch (state) {
    case types.FacetState.FACET_STALE:
      return 'Stale';
    case types.FacetState.FACET_FETCHING:
      return 'Fetching...';
    case types.FacetState.FACET_PARTIAL:
      return 'Partial';
    case types.FacetState.FACET_LOADED:
      return 'Loaded';
    case types.FacetState.FACET_ERROR:
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
