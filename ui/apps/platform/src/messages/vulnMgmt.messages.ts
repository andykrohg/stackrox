/* eslint-disable import/prefer-default-export */

export type ScanMessages = {
    header?: string;
    body?: string;
};

export const imageScanMessages = {
    missingMetadata: {
        header: 'Failed to retrieve metadata from the registry.',
        body: 'Couldn’t retrieve metadata from the registry, check registry connection.',
    },
    missingScanData: {
        header: 'Failed to get the base OS information.',
        body:
            'Failed to get the base OS information. Either the integrated scanner can’t find the OS or the base OS is unidentifiable.',
    },
    osUnavailable: {
        header: 'The scanner doesn’t provide OS information.',
        body:
            'Failed to get the base OS information. Either the integrated scanner can’t find the OS or the base OS is unidentifiable.',
    },
    languageCvesUnavailable: {
        header: 'Unable to retrieve the Language CVE data, only OS CVE data is available.',
        body:
            'Only showing information about the OS CVEs. Turn on the Language CVE feature in a scanner to view additional details.',
    },
    osCvesUnavailable: {
        header: 'Unable to retrieve the OS CVE data, only Language CVE data is available.',
        body: 'Only showing information about the Language CVEs.',
    },
    osCvesStale: {
        header: 'Stale OS CVE data..',
        body: 'The source no longer provides data updates.',
        extra: '',
    },
};
