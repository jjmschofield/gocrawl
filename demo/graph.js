(async function () {
    const getData = async () => {
        const response = await fetch('data/pages.jsonl');

        if (!response.ok) throw new Error(`failed to fetch pages with status ${response.status}`);

        const text = await response.text();

        const jsonStrs = text.split('\n');

        const nodes = [];
        const edges = [];

        jsonStrs.forEach(str => {
            if (str.length > 1) {
                const page = JSON.parse(str);

                nodes.push({ id: page.id, label: page.url });

                for (const outId in page.outPages.internal) {
                    edges.push({ source: page.id, target: outId });
                }
            }
        });

        return {
            nodes,
            edges
        };
    };

    const data = await getData();

    cytoscape({
        container: document.getElementById('graph'),

        elements: {
            nodes: data.nodes.map(node => {
                return { data: node }
            }),
            edges: data.edges.map(edge => {
                return { data: { id: edge.source + edge.target, source: edge.source, target: edge.target } }
            })
        },

        layout: {
            name: 'concentric',
            concentric: function( node ){
                return node.degree();
            },
            levelWidth: function( nodes ){
                return 25;
            }
        },

        style: [
            {
                selector: 'node',
                style: {
                    'content': 'data(label)',
                    'height': 20,
                    'width': 20,
                    'background-color': '#30c9bc'
                }
            },

            {
                selector: 'edge',
                style: {
                    'curve-style': 'haystack',
                    'haystack-radius': 0,
                    'width': 5,
                    'opacity': 0.5,
                    'line-color': '#a8eae5'
                }
            }
        ],
    });
})();