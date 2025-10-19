import { types } from '@models';

import { AssetChartsFacet } from './facets';
import {
  renderApprovalDetailPanel,
  renderStatementDetailPanel,
} from './panels';

export * from './panels';
export * from './facets';

export const renderers = {
  panels: {
    [types.DataFacet.OPENAPPROVALS]: renderApprovalDetailPanel,
    [types.DataFacet.STATEMENTS]: renderStatementDetailPanel,
  },
  facets: {
    [types.DataFacet.ASSETCHARTS]: AssetChartsFacet,
  },
};
