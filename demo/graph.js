(async function () {
    const getNodes = async () => {
        const response = await fetch('data/www.magicleap.com/pages.jsonl');

        if (!response.ok) throw new Error(`failed to fetch pages with status ${response.status}`);

        const pages = new Map();

        const text = await response.text();

        const jsonStrs = text.split('\n');

        jsonStrs.forEach(str => {
            if (str.length > 1) {
                const page = JSON.parse(str);
                pages.set(page.id, page);
            }
        });

        const nodes = Array.from(pages.values()).map(page => {
            // const label = page.url.replace("https://www.magicleap.com", "");

            return {id: page.id, label: page.url}
        });

        return new vis.DataSet(nodes);
    };

    const getEdges = async () => {
        const response = await fetch('data/www.magicleap.com/edges.csv');

        if (!response.ok) throw new Error(`failed to fetch pages with status ${response.status}`);

        const connections = new Map();

        const text = await response.text();

        const lines = text.split('\n');

        lines.forEach(line => {
            if (line.length > 1) {
                const edge = line.split(',');
                connections.set(`${edge[0]}:${edge[1]}`, edge);
            }
        });

        const edges = Array.from(connections.values()).map(edge => {
            return {from: edge[0], to: edge[1]}
        });

        return new vis.DataSet(edges);
    };

    const nodes = await getNodes();
    const edges = await getEdges();

    // create a network

    var data = {
        nodes: nodes,
        edges: edges
    };

    console.log("data", data);

    var options = {
        layout: {
            improvedLayout: true,
            hierarchical:{
                enabled: true,
                levelSeparation: 100,
                nodeSpacing: 1
            }
        },
        interaction: {dragNodes :true},
        physics: {
            enabled: true,
            timestep: 0.3,
            hierarchicalRepulsion:{
                springLength:1000,
                damping: 0.19
            },
            solver: "hierarchicalRepulsion",
            stabilization: {
                enabled : true,
            }
        },
        configure: {
            showButton:false
        },
        nodes: {
            shape: "box",
            widthConstraint: {
                maximum: 200
            },
            scaling:{
                min: 15,
                max: 20,
                label:{
                    enabled: true,
                    drawThreshold: 0,
                }
            }
        },
        edges: {
            arrows:{
              to:{
                  enabled: true,
                  scaleFactor: 0.1
              }
            },
            font: {
                size: 12
            },
            widthConstraint: {
                maximum: 90
            }
        },
    };

    var container = document.getElementById('graph');
    var network = new vis.Network(container, data, options);
})();