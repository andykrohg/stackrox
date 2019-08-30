import React, { useEffect, useState } from 'react';
import ReactRouterPropTypes from 'react-router-prop-types';
import * as Icon from 'react-feather';
import Tooltip from 'rc-tooltip';
import { generatePath } from 'react-router-dom';

import PageHeader from 'Components/PageHeader';
import Panel from 'Components/Panel';
import StatusField from 'Components/StatusField';
import ToggleSwitch from 'Components/ToggleSwitch';
import CheckboxTable from 'Components/CheckboxTable';
import { defaultColumnClassName, wrapClassName, rtTrActionsClassName } from 'Components/Table';
import TableHeader from 'Components/TableHeader';
import {
    fetchClustersAsArray,
    getAutoUpgradeConfig,
    saveAutoUpgradeConfig
} from 'services/ClustersService';
import { toggleRow, toggleSelectAll } from 'utils/checkboxUtils';
import { clustersPath } from 'routePaths';

import ClustersSidePanel from './ClustersSidePanel';

// @TODO, refactor these helper utilities to this folder,
//        when retiring clusters in Integrations section
import {
    formatClusterType,
    formatCollectionMethod,
    formatEnabledDisabledField,
    formatLastCheckIn,
    formatSensorVersion,
    parseUpgradeStatus
} from './cluster.helpers';

const ClustersPage = ({
    history,
    location: { search },
    match: {
        params: { clusterId }
    }
}) => {
    const [currentClusters, setCurrentClusters] = useState([]);
    const [autoUpgradeConfig, setAutoUpgradeConfig] = useState({});
    const [selectedClusters, setSelectedClusters] = useState([]);
    const [tableRef, setTableRef] = useState(null);
    const [selectedClusterId, setSelectedClusterId] = useState(clusterId);

    // @TODO, implement actual delete logic into this stub function
    const onDeleteHandler = cluster => e => {
        e.stopPropagation();
        setSelectedClusters([cluster.id]);
    };

    function renderRowActionButtons(cluster) {
        return (
            <div className="border-2 border-r-2 border-base-400 bg-base-100">
                <Tooltip placement="top" overlay={<div>Delete cluster</div>} mouseLeaveDelay={0}>
                    <button
                        type="button"
                        className="p-1 px-4 hover:bg-primary-200 text-primary-600 hover:text-primary-700"
                        onClick={onDeleteHandler(cluster)}
                    >
                        <Icon.Trash2 className="mt-1 h-4 w-4" />
                    </button>
                </Tooltip>
            </div>
        );
    }

    function toggleCluster(id) {
        const selection = toggleRow(id, selectedClusters);
        setSelectedClusters(selection);
    }

    function toggleAllClusters() {
        const rowsLength = selectedClusters.length;
        const ref = tableRef.reactTable;
        const selection = toggleSelectAll(rowsLength, selectedClusters, ref);
        setSelectedClusters(selection);
    }

    // @TODO: Change table component to use href for accessibility and better UX here, instead of an onclick
    function handleRowClick(cluster) {
        const newClusterId = (cluster && cluster.id) || '';
        setSelectedClusterId(newClusterId);
    }

    useEffect(() => {
        fetchClustersAsArray().then(clusters => {
            setCurrentClusters(clusters);
        });
    }, []);

    function fetchConfig() {
        getAutoUpgradeConfig().then(config => {
            setAutoUpgradeConfig(config);
        });
    }

    useEffect(() => {
        fetchConfig();
    }, []);

    // When the selected cluster changes, update the URL.
    useEffect(
        () => {
            const newPath = selectedClusterId
                ? generatePath(clustersPath, { clusterId: selectedClusterId })
                : clustersPath.replace('/:clusterId?', '');
            history.push({
                pathname: newPath,
                search
            });
        },
        [history, search, selectedClusterId]
    );

    function toggleAutoUpgrade() {
        // @TODO, wrap this settings change in a confirmation prompt of some sort
        const previousValue = autoUpgradeConfig.enableAutoUpgrade;
        const newConfig = { ...autoUpgradeConfig, enableAutoUpgrade: !previousValue };

        setAutoUpgradeConfig(newConfig); // optimistically set value before API call

        saveAutoUpgradeConfig(newConfig).catch(() => {
            // reverse the optimistic update of the control in the UI
            const rollbackConfig = { ...autoUpgradeConfig, enableAutoUpgrade: previousValue };
            setAutoUpgradeConfig(rollbackConfig);

            // also, re-fetch the data from the server, just in case it did update but we didn't get the network response
            fetchConfig();
        });
    }

    function getUpgradeStatusField(original) {
        const status = parseUpgradeStatus(original);
        return <StatusField displayValue={status.displayValue} type={status.type} />;
    }

    // @TODO: flesh out the new Clusters page layout, placeholders for now
    const paginationComponent = <div>Buttons and Pagination here</div>;

    const headerComponent = (
        <TableHeader length={currentClusters.length} type="Cluster" isViewFiltered={false} />
    );

    const clusterColumns = [
        {
            accessor: 'name',
            Header: 'Name',
            className: `${wrapClassName} ${defaultColumnClassName}`
        },
        {
            Header: 'Type',
            Cell: ({ original }) => formatClusterType(original.type)
        },
        {
            Header: 'Runtime Support',
            Cell: ({ original }) => formatCollectionMethod(original.collectionMethod)
        },
        {
            Header: 'Admission Controller Webhook',
            Cell: ({ original }) => formatEnabledDisabledField(original.admissionController)
        },
        {
            Header: 'Last Check-In',
            Cell: ({ original }) => formatLastCheckIn(original.status)
        },
        {
            Header: 'Upgrade status',
            Cell: ({ original }) => getUpgradeStatusField(original)
        },
        {
            Header: 'Current Sensor Version',
            Cell: ({ original }) => formatSensorVersion(original.status),
            className: `${wrapClassName} ${defaultColumnClassName} word-break`
        },
        {
            Header: '',
            accessor: '',
            headerClassName: 'hidden',
            className: rtTrActionsClassName,
            Cell: ({ original }) => renderRowActionButtons(original)
        }
    ];

    const headerText = 'Clusters';
    const subHeaderText = 'Resource list';

    const pageHeader = (
        <PageHeader header={headerText} subHeader={subHeaderText}>
            <div className="flex flex-1 justify-end">
                <div className="flex items-center">
                    <ToggleSwitch
                        id="enableAutoUpgrade"
                        toggleHandler={toggleAutoUpgrade}
                        label="Automatically upgrade secured clusters"
                        enabled={autoUpgradeConfig.enableAutoUpgrade}
                    />
                </div>
            </div>
        </PageHeader>
    );

    return (
        <section className="flex flex-1 flex-col h-full">
            <div className="flex flex-1 flex-col">
                {pageHeader}
                <div className="flex flex-1 relative">
                    <div className="shadow border-primary-300 bg-base-100 w-full overflow-hidden">
                        <Panel
                            headerTextComponent={headerComponent}
                            headerComponents={paginationComponent}
                        >
                            <div className="w-full">
                                <CheckboxTable
                                    ref={table => {
                                        setTableRef(table);
                                    }}
                                    rows={currentClusters}
                                    columns={clusterColumns}
                                    onRowClick={handleRowClick}
                                    toggleRow={toggleCluster}
                                    toggleSelectAll={toggleAllClusters}
                                    selection={selectedClusters}
                                    selectedRowId={selectedClusterId}
                                    noDataText="No clusters to show."
                                    minRows={20}
                                />
                            </div>
                        </Panel>
                    </div>
                    <ClustersSidePanel
                        selectedClusterId={selectedClusterId}
                        setSelectedClusterId={setSelectedClusterId}
                    />
                </div>
            </div>
        </section>
    );
};

ClustersPage.propTypes = {
    history: ReactRouterPropTypes.history.isRequired,
    location: ReactRouterPropTypes.location.isRequired,
    match: ReactRouterPropTypes.match.isRequired
};

export default ClustersPage;
