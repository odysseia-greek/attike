<script>
    import mermaid from 'mermaid';


    mermaid.initialize({
        startOnLoad: true,
        theme: 'dark',
        htmlLabels: true,
        gantt: {
            titleTopMargin: 50,
            barHeight: 30,
            barGap: 20,
            topPadding: 75,
            rightPadding: 75,
            leftPadding: 75,
            fontSize: 15,
            sectionFontSize: 25,
            gridLineStartPadding: 20,
        },
        xyChart: {
            width: 2000,
            height: 600,
        }
    });


    export let mermaidDiagram; // Accepting the JSON string as a prop
    let mermaidContainer;

    async function renderDiagram() {
        const x = Math.floor(Math.random()*5000);
        const { svg } = await mermaid.render("diagram"+x, mermaidDiagram);
        mermaidContainer.innerHTML = svg;
    }

    $: if (mermaidDiagram) {
        try {
            renderDiagram();
        } catch (error) {
            console.error('Invalid JSON:', error);
            // Instead of setting mermaidDiagram to '', consider handling the error differently
        }
    }


</script>

<pre bind:this={mermaidContainer} class="mermaid">
</pre>

<style>
    .mermaid {
        width: 90%;
        margin: auto;
        overflow: auto;
    }
</style>