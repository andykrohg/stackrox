import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { withRouter } from 'react-router-dom';
import Query from 'Components/ThrowingQuery';
import URLSearchInputWithAutocomplete from 'Components/URLSearchInputWithAutocomplete';

import SEARCH_AUTOCOMPLETE_QUERY from 'queries/searchAutocomplete';

// eslint-disable-next-line react/no-multi-comp
const URLSearchInput = props => {
    const [autoCompleteQuery, setAutoCompleteQuery] = useState('');

    function clearAutocomplete() {
        setAutoCompleteQuery('');
    }

    function fetchAutocomplete({ query }) {
        setAutoCompleteQuery(query);
    }

    if (!autoCompleteQuery) {
        return <URLSearchInputWithAutocomplete fetchAutocomplete={fetchAutocomplete} {...props} />;
    }
    return (
        <Query
            query={SEARCH_AUTOCOMPLETE_QUERY}
            variables={{
                query: autoCompleteQuery,
                categories: props.categories
            }}
        >
            {({ data }) => {
                const autoCompleteResults = data ? data.searchAutocomplete : [];
                return (
                    <URLSearchInputWithAutocomplete
                        autoCompleteResults={autoCompleteResults}
                        clearAutocomplete={clearAutocomplete}
                        fetchAutocomplete={fetchAutocomplete}
                        {...props}
                    />
                );
            }}
        </Query>
    );
};

URLSearchInput.propTypes = {
    categories: PropTypes.arrayOf(PropTypes.string)
};

URLSearchInput.defaultProps = {
    categories: []
};

export default withRouter(URLSearchInput);
