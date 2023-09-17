<!-- JSONViewWrapper.svelte -->
<script>
    import { JsonView } from '@zerodevx/svelte-json-view'
    export let jsonString;

    // Parse the JSON string into a JavaScript object
    let json = {};

    $: {
        try {
            json = JSON.parse(jsonString);
        } catch (error) {
            console.error('Invalid JSON:', error);
        }
    }
</script>

<div class="wrap">
    <JsonView json={json} depth="1"/>
</div>


<style>
    .wrap {
        color-scheme: dark;
        font-family: ui-monospace,SFMono-Regular,Menlo,Monaco,Consolas,Liberation Mono,Courier New,monospace;
        font-size: .875rem;
        --jsonPaddingLeft: 2rem;
        --jsonBorderLeft: 1px dotted #A6ADBA;
        --jsonValColor: #4b5563;
        --jsonValStringColor: #34d399;
        --jsonValNumberColor: #ff7100;
        --jsonValBooleanColor: #60a5fa;
        --jsonKeyColor: #A6ADBA;
        --jsonSeparatorColor: #A6ADBA;
        --jsonBracketColor: #ff7100;
        text-align: left;
        margin-left: 2em;
        background-color: #16082f;
        padding-left: 2em;
        padding-bottom: 1em;
        padding-top: 1em;
    }
</style>
