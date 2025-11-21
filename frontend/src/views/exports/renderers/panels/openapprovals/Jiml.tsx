// Copyright 2016, 2026 The Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
import React, { useCallback, useMemo, useState } from 'react';

import {
  BorderedSection,
  DetailPanelContainer,
  InfoAddressRenderer,
  PanelRow,
  PanelTable,
  StyledButton,
  appToAddressInfo,
  txToDetailsInfo,
} from '@components';
import { useWalletGatedAction } from '@hooks';
import { Group, Stack, Text } from '@mantine/core';
import { types } from '@models';
import {
  Log,
  LogError,
  addressToHex,
  buildTransaction,
  displayHash,
  formatNumericValue,
} from '@utils';
import {
  PreparedTransaction,
  TransactionData,
  useWalletConnection,
} from '@utils';

import '../../../../../components/detail/DetailTable.css';
import { TransactionReviewModal } from '../../../../contracts/renderers/components/execute/TransactionReviewModal';

/*
// Helper functions
const formatTimestamp = (timestamp: number | string): string => {
  const numTimestamp = Number(timestamp);
  if (isNaN(numTimestamp) || numTimestamp <= 0) {
    return 'No timestamp';
  }
  return new Date(numTimestamp * 1000).toLocaleString(undefined, {
    hour12: false,
  });
};

const truncateAddress = (address: unknown): string => {
  if (!address) return 'N/A';
  const hex = addressToHex(address);
  if (!hex || hex.length < 10) return hex;
  return `${hex.slice(0, 6)}...${hex.slice(-4)}`;
};
*/

export const OpenApprovalsPanel = (rowData: Record<string, unknown> | null) => {
  // Collapse state management
  const [collapsed, setCollapsed] = useState<Set<string>>(new Set());

  const [transactionModal, setTransactionModal] = useState<{
    opened: boolean;
    transactionData: TransactionData | null;
  }>({ opened: false, transactionData: null });

  const { createWalletGatedAction } = useWalletGatedAction();
  const { sendTransaction } = useWalletConnection({
    onTransactionSigned: (txHash: string) => {
      Log('✅ Revoke transaction signed:', txHash);
      Log('🔍 View on Etherscan: https://etherscan.io/tx/' + txHash);
      setTransactionModal({ opened: false, transactionData: null });
    },
    onError: (error: string) => {
      LogError('Revoke transaction error:', error);
    },
  });

  // Handle section toggle
  const handleToggle = (sectionName: string) => {
    const isCollapsed = collapsed.has(sectionName);
    if (isCollapsed) {
      setCollapsed((prev) => {
        const next = new Set(prev);
        next.delete(sectionName);
        return next;
      });
    } else {
      setCollapsed((prev) => new Set([...prev, sectionName]));
    }
  };

  // Memoize approval conversion to avoid dependency warnings
  const approval = useMemo(
    () => (rowData as unknown as types.Approval) || ({} as types.Approval),
    [rowData],
  );

  // Memoized converter functions (called before early return to maintain hook order)
  const detailsInfo = useMemo(() => {
    if (!rowData) return null;
    // Create a pseudo-transaction structure for the InfoDetailsRenderer
    const pseudoTransaction = {
      hash: approval.token, // Use token address as hash for display
      blockNumber: approval.blockNumber,
      blockHash: approval.token, // Use token as block hash placeholder
      transactionIndex: 0,
      timestamp: approval.timestamp,
      nonce: 0,
      type: 'approval',
      value: approval.allowance,
      from: approval.owner,
      fromName: approval.ownerName,
      to: approval.spender,
      toName: approval.spenderName,
      // Add missing required fields with default values
      gas: 0,
      gasPrice: 0,
      gasUsed: 0,
      hasToken: false,
      isError: false,
      receipt: null,
      traces: [],
    };
    return txToDetailsInfo(pseudoTransaction as unknown as types.Transaction);
  }, [rowData, approval]);

  const addressInfo = useMemo(() => {
    if (!rowData) return null;
    return appToAddressInfo(
      approval.owner,
      approval.ownerName,
      approval.spender,
      approval.spenderName,
    );
  }, [rowData, approval]);

  const tokenInfo = useMemo(() => {
    if (!rowData) return null;
    return [
      {
        label: 'Token',
        value: addressToHex(approval.token),
        name: approval.tokenName,
      },
      {
        label: 'Name',
        value: approval.tokenName || 'Unknown Token',
      },
    ];
  }, [rowData, approval]);

  const allowanceInfo = useMemo(() => {
    if (!rowData) return null;
    return [
      {
        label: 'Allowance',
        value: formatNumericValue(approval.allowance || 0),
        isHighlight: Number(approval.allowance) > 0,
      },
      {
        label: 'Block',
        value: approval.lastAppBlock?.toString() || 'N/A',
      },
      {
        label: 'Timestamp',
        value: approval.lastAppTs
          ? new Date(approval.lastAppTs * 1000).toLocaleString()
          : 'N/A',
      },
    ];
  }, [rowData, approval]);

  // Create revoke transaction
  const createRevokeTransaction = useCallback(() => {
    try {
      Log('Before approveFunction');

      const approveFunction: types.Function = {
        name: 'approve',
        type: 'function',
        inputs: [
          types.Parameter.createFrom({ name: 'spender', type: 'address' }),
          types.Parameter.createFrom({ name: 'amount', type: 'uint256' }),
        ],
        outputs: [types.Parameter.createFrom({ name: '', type: 'bool' })],
        stateMutability: 'nonpayable',
        encoding: '0x095ea7b3', // ERC20 approve function selector
        convertValues: () => {},
      };

      Log('Before transactionInputs');

      const transactionInputs = [
        {
          name: 'spender',
          type: 'address',
          value: addressToHex(approval.spender),
        },
        { name: 'amount', type: 'uint256', value: '0' }, // Set to zero to revoke
      ];

      Log('Before buildTransaction');

      const transactionData = buildTransaction(
        addressToHex(approval.token), // Contract address is the token
        approveFunction,
        transactionInputs,
      );

      Log('Revoke transaction data created:', JSON.stringify(transactionData));

      // Open the transaction modal
      setTransactionModal({
        opened: true,
        transactionData,
      });

      Log('After setTransactionModel');
    } catch (error) {
      LogError('Creating revoke transaction:', String(error));
    }
  }, [approval]);

  // Handle revoke action with wallet gating
  const handleRevoke = createWalletGatedAction(() => {
    createRevokeTransaction();
  }, 'Revoke');

  // Handle transaction confirmation from modal
  const handleConfirmTransaction = useCallback(
    async (preparedTx: PreparedTransaction) => {
      try {
        await sendTransaction(preparedTx);
      } catch (error) {
        LogError('Failed to send revoke transaction:', String(error));
      }
    },
    [sendTransaction],
  );

  // Handle modal close
  const handleModalClose = useCallback(() => {
    setTransactionModal({ opened: false, transactionData: null });
  }, []);

  // Show loading state if no data is provided
  if (!rowData) {
    return <div className="no-selection">Loading...</div>;
  }

  // Early return after all hooks if computed data is invalid
  if (!detailsInfo || !addressInfo || !tokenInfo || !allowanceInfo) {
    return null;
  }

  // Title component with key identifying info
  const titleComponent = () => (
    <Group justify="space-between" align="flex-start">
      <Text variant="primary" size="md" fw={600}>
        Approval {displayHash(approval.token)}
      </Text>
      <Text variant="primary" size="md" fw={600}>
        {approval.tokenName || 'Token'}
      </Text>
    </Group>
  );

  return (
    <Stack gap={0} className="fixed-prompt-width">
      <Group
        justify="flex-end"
        style={{ padding: '8px 0', marginBottom: '8px' }}
      >
        <StyledButton
          onClick={handleRevoke}
          variant="warning"
          size="sm"
          disabled={!approval.token || !approval.spender}
        >
          Revoke
        </StyledButton>
      </Group>

      <DetailPanelContainer title={titleComponent()}>
        <BorderedSection>
          <div
            onClick={() => handleToggle('Address Information')}
            style={{ cursor: 'pointer' }}
          >
            <Text variant="primary" size="sm">
              <div className="detail-section-header">
                {collapsed.has('Address Information') ? '▶ ' : '▼ '}Address
                Information
              </div>
            </Text>
          </div>
          {!collapsed.has('Address Information') && (
            <InfoAddressRenderer addressInfo={addressInfo} />
          )}
        </BorderedSection>

        <BorderedSection>
          <div
            onClick={() => handleToggle('Token Information')}
            style={{ cursor: 'pointer' }}
          >
            <Text variant="primary" size="sm">
              <div className="detail-section-header">
                {collapsed.has('Token Information') ? '▶ ' : '▼ '}Token
                Information
              </div>
            </Text>
          </div>
          {!collapsed.has('Token Information') && (
            <PanelTable>
              {tokenInfo.map((item, index) => (
                <PanelRow
                  label={item.label}
                  key={index}
                  value={<div className="panel-nested-name">{item.value}</div>}
                />
              ))}
            </PanelTable>
          )}
          {/* <PanelRow
                  key={index}
                  label={item.label}
                  value={
                    <span
                      style={{
                        fontFamily: item.label.includes('Address')
                          ? 'monospace'
                          : 'inherit',
                        fontSize: '14px',
                      }}
                      title={String(item.value)}
                    >
                      {item.name && item.label.includes('Address')
                        ? `${item.value.slice(0, 6)}...${item.value.slice(-4)}`
                        : item.value}
                      {item.name && item.label.includes('Address') && (
                        <div
                          style={{
                            fontSize: '12px',
                            color: 'var(--mantine-color-dimmed)',
                            marginTop: '2px',
                            overflow: 'hidden',
                            textOverflow: 'ellipsis',
                            whiteSpace: 'nowrap',
                          }}
                        >
                          {item.name}
                        </div>
                      )}
                    </span>
                  }
                /> */}
        </BorderedSection>

        {/* <BorderedSection>
          <div
            onClick={() => handleToggle('Allowance Details')}
            style={{ cursor: 'pointer' }}
          >
            <Text variant="primary" size="sm">
              <div className="detail-section-header">
                {collapsed.has('Allowance Details') ? '▶ ' : '▼ '}Allowance
                Details
              </div>
            </Text>
          </div>
          {!collapsed.has('Allowance Details') && (
            <PanelTable>
              {allowanceInfo.map((item, index) => (
                <PanelRow
                  key={index}
                  label={item.label}
                  value={
                    <span
                      style={{
                        fontFamily: item.label.includes('Allowance')
                          ? 'monospace'
                          : 'inherit',
                        fontSize: '14px',
                        fontWeight: item.isHighlight ? 600 : 'normal',
                        color: item.isHighlight
                          ? 'var(--mantine-color-red-6)'
                          : 'inherit',
                      }}
                      title={String(item.value)}
                    >
                      {item.value}
                    </span>
                  }
                />
              ))}
            </PanelTable>
          )}
        </BorderedSection> */}

        {/* <BorderedSection>
          <div
            onClick={() => handleToggle('Approval Details')}
            style={{ cursor: 'pointer' }}
          >
            <Text variant="primary" size="sm">
              <div className="detail-section-header">
                {collapsed.has('Approval Details') ? '▶ ' : '▼ '}Approval
                Details
              </div>
            </Text>
          </div>
          {!collapsed.has('Approval Details') && (
            <InfoDetailsRenderer detailsInfo={detailsInfo} />
          )}
        </BorderedSection> */}
      </DetailPanelContainer>

      {/* Transaction Review Modal */}
      <TransactionReviewModal
        opened={transactionModal.opened}
        onClose={handleModalClose}
        transactionData={transactionModal.transactionData}
        onConfirm={handleConfirmTransaction}
      />
    </Stack>
  );
};
