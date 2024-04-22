<!-- SearchFilter.svelte -->
<script>
    import JSONViewWrapper from './JSONViewWrapper.svelte'; // Import the JSONViewWrapper component
    import builderQuery from '../graphql/queries/builderTrace.graphql';
    import gql from 'graphql-tag';
    import { onMount } from 'svelte';
    import client from "../graphql/client.js";

    const formatDateTime = (dateTime) => {
        if (dateTime) {
            const date = new Date(dateTime);
            return date.toISOString().slice(0, 23); // Format as "YYYY-MM-DDTHH:MM:SS.sss"
        }
        return undefined;
    };

    let exampleQuery = {
        ids: [""],
        podName: "",
        operation: "",
        statusCode: "",
        totalTimeHigherThan: "",
        beginTime: "",
        endTime: "",
    };

    let traceIDs = [];
    let podName = ''; // Store the podName filter
    let operation = ''; // Store the operation filter
    let statusCode = ''; // Store the statusCode filter
    let totalTimeHigherThan = ''; // Store the totalTimeHigherThan filter

    // Initialize variables for the time range
    let beginTime = ''; // Store the beginTime filter
    let endTime = ''; // Store the endTime filter

    // Function to trigger the GraphQL query with the search criteria
    function search() {
        // Create an object to store the filter criteria
        const filters = {};

        // Add non-empty values to the filters object
        if (traceIDs.length > 0) {
            filters.ids = traceIDs;
        }

        if (podName !== '') {
            filters.podName = podName;
        }
        if (operation !== '') {
            filters.operation = operation;
        }
        if (statusCode !== '') {
            filters.statusCode = statusCode;
        }
        if (totalTimeHigherThan !== '') {
            filters.totalTimeHigherThan = totalTimeHigherThan;
        }
        if (beginTime !== '') {
            beginTime = formatDateTime(beginTime)
            filters.beginTime = beginTime;
        }
        if (endTime !== '') {
            endTime = formatDateTime(endTime)
            filters.endTime = endTime;
        }

        // Trigger the custom search event and pass the filter criteria
        const customEvent = new CustomEvent('searchfiltersubmit', { detail: filters });
        document.dispatchEvent(customEvent);
    }

    function clear() {
        traceIDs = [];
        podName = '';
        operation = '';
        statusCode = '';
        totalTimeHigherThan = '';
        beginTime = '';
        endTime = '';

        exampleQuery = {
            ids: [""],
            podName: "",
            operation: "",
            statusCode: "",
            totalTimeHigherThan: "",
            beginTime: "",
            endTime: "",
        };

        search()
    }

    // Function to update the example query based on filter changes
    function updateExampleQuery() {
        exampleQuery = {
            ids: traceIDs,
            podName: podName.trim() !== '' ? podName : undefined,
            operation: operation.trim() !== '' ? operation : undefined,
            statusCode: statusCode.trim() !== '' ? statusCode : undefined,
            totalTimeHigherThan: totalTimeHigherThan.trim() !== '' ? totalTimeHigherThan : undefined,
            beginTime: formatDateTime(beginTime),
            endTime: formatDateTime(endTime),
        };

        search()
    }

    // Define your GraphQL query for the initial data retrieval
    const BUILDER_TRACES = gql`${builderQuery}`;

    let ids = [];
    let responseCodes = [];
    let operations = [];
    let podNames = [];

    // Perform the initial query and populate the options
    onMount(async () => {
        search()
        const variables = {
            input: {},
        };
        // Trigger the GraphQL query with the updated variables
        const result = await client.query({
            query: BUILDER_TRACES,
            variables: variables
        });

        if (result.data && result.data.traces) {
            const traces = result.data.traces;

            traces.forEach(trace => {
                if ('traceID' in trace) {
                    ids.push(trace.traceID)
                }
                if ('responseCode' in trace) {
                    const responseCode = trace.responseCode;
                    if (responseCode !== '' && !responseCodes.includes(responseCode)) {
                        responseCodes.push(responseCode);
                    }
                }
                trace.items.forEach(item => {
                    if ('operation' in item) {
                        const operation = item.operation;
                        if (operation !== '' && !operations.includes(operation)) {
                            operations.push(operation);
                        }
                    }
                    if ('podName' in item) {
                        const podName = item.podName;
                        if (podName !== '' && !podNames.includes(podName)) {
                            podNames.push(podName);
                        }
                    }
                });
            });
        }
    });

    // Variables to control the visibility of options
    let showIdsOptions = false;
    let showPodNameOptions = false;
    let showOperationOptions = false;
    let showStatusCodeOptions = false;

    // Function to toggle visibility of options
    function showOptions(field) {
        switch (field) {
            case 'ids':
                showIdsOptions = !showIdsOptions;
                break;
            case 'podName':
                showPodNameOptions = !showPodNameOptions;
                break;
            case 'operation':
                showOperationOptions = !showOperationOptions;
                break;
            case 'statusCode':
                showStatusCodeOptions = !showStatusCodeOptions;
                break;
            // Add more cases for other input fields if needed
        }
    }

    // Function to set the search term and hide options
    function setSearchTerm(inputField, term) {
        switch (inputField) {
            case 'ids':
                traceIDs = [term];
                showIdsOptions = false;
                break;
            case 'podName':
                podName = term;
                showPodNameOptions = false;
                break;
            case 'operation':
                operation = term;
                showOperationOptions = false;
                break;
            case 'statusCode':
                statusCode = term;
                showStatusCodeOptions = false;
                break;
            // Add more cases for other input fields if needed
        }
        updateExampleQuery()
    }

</script>

<style>
    .search-container {
        display: flex;
        justify-content: space-between; /* Spread the elements horizontally */
        align-items: flex-start;
        max-width: 100%; /* Occupy 100% of the main container */
        margin: 0 auto;
    }

    /* Adjust width of search options and json view */
    .search-options,
    .json-view {
        flex: 1; /* Occupy equal space */
        padding: 0 10px; /* Add padding to separate the elements */
    }

    .search-options {
        flex: 1;
        display: flex;
        flex-direction: column; /* Stack elements vertically within .search-options */
        align-items: center;
    }

    .search-options h1 {
        color: #ff7100;
        font-weight: 100;
        letter-spacing: 0.01em;
        margin-left: 15px;
        margin-bottom: 35px;
        text-transform: uppercase;
    }

    .search-options input[type="text"],
    .search-options input[type="datetime-local"],
    .search-options button {
        margin-top: 35px;
        background-color: #94999f;
        border: 1px solid #ff7100;
        line-height: 1;
        font-size: 17px;
        display: inline-block;
        box-sizing: border-box;
        padding: 20px 15px;
        border-radius: 60px;
        color: #444649;
        font-weight: 100;
        letter-spacing: 0.01em;
        position: relative;
        z-index: 1;
        width: 15em;
    }

    /* Change button colors on hover */
    .search-options button:hover {
        background-color: #ff7100;
        color: #444649;
    }

    .search-options input[type="text"]:focus,
    .search-options input[type="datetime-local"]:focus {
        outline: none;
        background: #ff7100;
        color: #444649;
        margin-top: 1px;
    }

    .search-options input[type="text"]:valid {
        margin-top: 30px;
    }

    /* Input styles with placeholder color */
    .search-options input[type="text"]::placeholder {
        color: #444649;
    }

    .search-options input[type="text"]:focus ~ label,
    .search-options input[type="datetime-local"]:focus ~ label {
        transform: translate(0, -35px);
    }

    .search-options input[type="text"]:valid ~ label,
    .search-options input[type="datetime-local"]:valid ~ label {
        text-transform: uppercase;
        font-style: italic;
        transform: translate(5px, -35px) scale(0.6);
    }

    .search-options label {
        transform-origin: left center;
        color: #ff7100;
        font-weight: 100;
        letter-spacing: 0.01em;
        font-size: 17px;
        box-sizing: border-box;
        padding: 10px 15px;
        display: block;
        position: absolute;
        margin-top: -40px;
        z-index: 2;
        pointer-events: none;
    }

</style>


<div class="search-container">
    <!-- ... (other parts of your Svelte component) ... -->

    <div class="search-options">
        <!-- ... (other input fields) ... -->

        <input
                type="text"
                placeholder="PodName"
                bind:value={podName}
                on:input={updateExampleQuery}
                on:click={() => showOptions('podName')}
        />
        {#if showPodNameOptions}
            <div class="options-dropdown">
                <select on:change={(event) => setSearchTerm('podName', event.target.value)}>
                    <option value="">Select a PodName</option>
                    {#each podNames as pod }
                        <option value={pod}>{pod}</option>
                    {/each}
                </select>
            </div>
        {/if}

        <input
                type="text"
                placeholder="Ids"
                bind:value={ids}
                on:input={updateExampleQuery}
                on:click={() => showOptions('ids')}
        />
        {#if showIdsOptions}
            <div class="options-dropdown">
                <select on:change={(event) => setSearchTerm('ids', event.target.value)}>
                    <option value="">Select an ID</option>
                    {#each ids as id }
                        <option value={id}>{id}</option>
                    {/each}
                </select>
            </div>
        {/if}

        <input
                type="text"
                placeholder="Operation"
                bind:value={operation}
                on:input={updateExampleQuery}
                on:click={() => showOptions('operation')}
        />
        {#if showOperationOptions}
            <div class="options-dropdown">
                <select on:change={(event) => setSearchTerm('operation', event.target.value)}>
                    <option value="">Select an Operation</option>
                    {#each operations as op }
                        <option value={op}>{op}</option>
                    {/each}
                </select>
            </div>
        {/if}

        <input
                type="text"
                placeholder="StatusCode"
                bind:value={statusCode}
                on:input={updateExampleQuery}
                on:click={() => showOptions('statusCode')}
        />
        {#if showStatusCodeOptions}
            <div class="options-dropdown">
                <select on:change={(event) => setSearchTerm('statusCode', event.target.value)}>
                    <option value="">Select a StatusCode</option>
                    {#each responseCodes as code }
                        <option value={code}>{code}</option>
                    {/each}
                </select>
            </div>
        {/if}

        <input type="text" placeholder="QueryTime" bind:value={totalTimeHigherThan} on:input={updateExampleQuery} />
        <input type="datetime-local" placeholder="Begin Time" bind:value={beginTime} on:input={updateExampleQuery} />
        <input type="datetime-local" placeholder="End Time" bind:value={endTime} on:input={updateExampleQuery} />
        <button on:click={search}>Search</button>
        <button on:click={clear}>Clear</button>
    </div>
    <!-- Display the dynamic example query to the user -->
    <div class="json-view">
        <h2>Query Parameters</h2>
        <pre>
            <JSONViewWrapper jsonString={JSON.stringify(exampleQuery, null, 2)} />
        </pre>
    </div>

</div>

