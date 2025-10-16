import { useEffect, useMemo, useRef, useState } from 'react';

import { types } from '@models';

interface UsePlaceholderRowsProps {
  data: unknown[] | null | undefined;
  state: types.FacetState;
}

export function usePlaceholderRows({ data, state }: UsePlaceholderRowsProps) {
  const [cycleIndex, setCycleIndex] = useState(0);
  const timersRef = useRef<Set<ReturnType<typeof setTimeout>>>(new Set());

  // Determine if we should show placeholders based on state
  const shouldShowPlaceholders = useMemo(() => {
    switch (state) {
      case types.FacetState.FACET_FETCHING:
        return data == null || data.length === 0; // Only if no data yet
      case types.FacetState.FACET_STALE:
        return data == null || data.length === 0; // Initial load case
      case types.FacetState.FACET_LOADED:
      case types.FacetState.FACET_ERROR:
        return false; // Never show placeholders - we have final answer
      case types.FacetState.FACET_PARTIAL:
        return false; // We have some data, show it
      default:
        return false;
    }
  }, [state, data]);

  // Clear all timers helper
  const clearAllTimers = () => {
    timersRef.current.forEach((timer) => clearTimeout(timer));
    timersRef.current.clear();
  };

  useEffect(() => {
    // Clear timers if we shouldn't show placeholders
    if (!shouldShowPlaceholders) {
      clearAllTimers();
      setCycleIndex(0);
      return;
    }

    // Start after 1 second delay if we should show placeholders
    const delayTimer = setTimeout(() => {
      setCycleIndex(0); // Start cycling

      // Start interval timer for cycling through [3, 5, 7] counts
      const intervalTimer = setInterval(() => {
        setCycleIndex((prev) => (prev + 1) % 3);
      }, 1000);

      timersRef.current.add(intervalTimer);
    }, 1000); // 1 second delay

    timersRef.current.add(delayTimer);

    // Cleanup function
    return () => {
      clearAllTimers();
    };
  }, [shouldShowPlaceholders]);

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      clearAllTimers();
    };
  }, []);

  const placeholderCount = shouldShowPlaceholders ? [3, 5, 7][cycleIndex] : 0;

  return { placeholderCount };
}
