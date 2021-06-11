import React, { ReactElement } from 'react';
import { TableComposable, Tbody, Td, Thead, Th, Tr } from '@patternfly/react-table';

import { availableAuthProviders } from 'constants/accessControl';
import { AuthProvider } from 'services/AuthService';

import { AccessControlEntityLink } from '../AccessControlLinks';

// TODO import from where?
const unselectedRowStyle = {};
const selectedRowStyle = {
    borderLeft: '3px solid var(--pf-global--primary-color--100)',
};

function getAuthProviderTypeLabel(type: string): string {
    return availableAuthProviders.find(({ value }) => value === type)?.label ?? '';
}

const entityType = 'AUTH_PROVIDER';

export type AuthProvidersListProps = {
    entityId?: string;
    authProviders: AuthProvider[];
};

function AuthProvidersList({ entityId, authProviders }: AuthProvidersListProps): ReactElement {
    return (
        <TableComposable variant="compact">
            <Thead>
                <Tr>
                    <Th>Name</Th>
                    <Th>Type</Th>
                    <Th>Minimum access role</Th>
                    <Th>Rules</Th>
                </Tr>
            </Thead>
            <Tbody>
                {authProviders.map(({ id, name, type }) => {
                    const typeLabel = getAuthProviderTypeLabel(type);
                    // TODO for minimumAccessRoleName see getDefaultRoleByAuthProviderId in classic code

                    return (
                        <Tr
                            key={id}
                            style={id === entityId ? selectedRowStyle : unselectedRowStyle}
                        >
                            <Td dataLabel="Name">
                                <AccessControlEntityLink
                                    entityType={entityType}
                                    entityId={id}
                                    entityName={name}
                                />
                            </Td>
                            <Td dataLabel="Type">{typeLabel}</Td>
                            <Td dataLabel="Minimum access role">
                                {/* <AccessControlEntityLink
                                    entityType="ROLE"
                                    entityId={minimumAccessRoleName}
                                    entityName={minimumAccessRoleName}
                                /> */}
                            </Td>
                            <Td dataLabel="Rules">{/* TODO */}</Td>
                        </Tr>
                    );
                })}
            </Tbody>
        </TableComposable>
    );
}

export default AuthProvidersList;
