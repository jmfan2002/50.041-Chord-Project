const outputElem = document.getElementById('output');

getNodes();

async function makeHealthCheck(node) {
    console.log('Making health check');

    try {
        const response = await fetch('/health?node=' + node, { method: 'GET' });
        await handleResponse(response);
    } catch (error) {
        handleFetchError(error);
    }
}

async function makeCycleHealthCheck() {
    // Cycle health of node
    console.log('Cycling health');
    try {
        const healthRes = await fetch(`/cycleHealth`, {
            method: 'GET',
        });

        handleResponse(healthRes);
    } catch (error) {
        handleFetchError(error);
    }
}

async function getData() {
    const keyInp = document.querySelector('input[name=getKey]');
    const key = keyInp.value;
    console.log('Getting data with key:', key);

    try {
        const res = await fetch(`/data?key=${key}`, { method: 'GET' });
        const outputElem = document.getElementById('dataOutput');
        if (res.ok) {
            try {
                const data = await res.json();
                if (!data?.value) {
                    outputElem.innerText = 'Value: No data found';
                } else {
                    outputElem.innerText = 'Value: ' + data.value;
                }
            } catch (error) {
                console.log(data);
            }
        } else {
            const errorMessage = `Error: ${res.statusText}`;
            outputElem.innerText = errorMessage;
            outputElem.innerText = 'Value: Something went wrong';
            console.error(errorMessage);
        }
    } catch (error) {
        handleFetchError(error);
    }
}

async function addData() {
    console.log('Adding data');

    const keyInp = document.querySelector('input[name=setKey]');
    const valueInp = document.querySelector('input[name=setValue]');
    const data = {
        key: keyInp.value,
        value: valueInp.value,
    };

    console.log('Data:', data);

    try {
        const res = await fetch('/data', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data),
        });

        if (!res.ok) {
            const errorMessage = `Error: ${res.statusText}`;
            outputElem.innerText = errorMessage;
            console.error(errorMessage);
        } else {
            const body = await res.json();
            const responseElem = document.querySelector('.response');
            responseElem.innerText = 'Response: ' + body?.message;
            console.log(body);
        }
    } catch (error) {
        handleFetchError(error);
    }
}

async function getNodes() {
    console.log('Retrieving nodes');

    try {
        const res = await fetch('/nodes', {
            method: 'GET',
        });

        if (!res.ok) {
            const errorMessage = `Error: ${res.statusText}`;
            outputElem.innerText = errorMessage;
            console.error(errorMessage);
        } else {
            const body = await res.json();

            const listContainer = document.getElementById('nodeList');
            listContainer.innerHTML = '';

            body.nodes?.forEach((node) => {
                listContainer.appendChild(createNodeComponent(node));
            });
        }
    } catch (error) {
        handleFetchError(error);
    }
}

async function getHashTable() {
    console.log('Retrieving hash table');

    try {
        const res = await fetch('/data/hashTable', {
            method: 'GET',
        });

        if (!res.ok) {
            const errorMessage = `Error: ${res.statusText}`;
            outputElem.innerText = errorMessage;
            console.error(errorMessage);
        } else {
            const body = await res.json();

            const listContainer = document.getElementById('hashTable');
            listContainer.innerHTML = '';

            body.hashTable?.forEach((kvPair) => {
                listContainer.appendChild(createKvPairComponent(kvPair));
            });
        }
    } catch (error) {
        handleFetchError(error);
    }
}

function createNodeComponent(node) {
    const nodeElem = document.createElement('div');
    const nodeText = document.createElement('span');

    nodeText.innerText = node;

    nodeElem.appendChild(nodeText);

    return nodeElem;
}

function createKvPairComponent(kvPair) {
    // kvPair = {node: "adress", key: 'key', value: 'value'}
    const kvPairElem = document.createElement('div');
    const kvPairText = document.createElement('span');

    kvPairText.innerText = `${kvPair.node} - ${kvPair.key} : ${kvPair.value}`;

    kvPairElem.appendChild(kvPairText);

    return kvPairElem;
}

async function handleResponse(res) {
    if (!res.ok) {
        const errorMessage = `Error: ${res.statusText}`;
        outputElem.innerText = errorMessage;
        console.error(errorMessage);
    } else {
        const body = await res.json();
        outputElem.innerText = JSON.stringify(body);
        console.log(body);
    }
}

setInterval(async () => {
    // Get node from the node container
    document.querySelector('.loader').classList.remove('hidden');
    const nodeList = document.getElementById('nodeList');
    await getNodes();

    // For each node call the health check on the node
    for (let i = 0; i < nodeList.children.length; i++) {
        const nodeTextElem = nodeList.children[i].querySelector('span');
        const nodeAdress = nodeTextElem.innerText;
        try {
            const res = await fetch('/health?node=' + nodeAdress, {
                method: 'GET',
            });
            if (!res.ok) {
                const errorMessage = `Error: ${res.statusText}`;
                console.error(errorMessage);
                nodeTextElem.classList.remove('nodeUp');
                nodeTextElem.classList.add('nodeDown');
            } else {
                const body = await res.json();
                if (body?.val === 'Node is healthy') {
                    nodeTextElem.classList.remove('nodeDown');
                    nodeTextElem.classList.add('nodeUp');
                } else {
                    nodeTextElem.classList.remove('nodeUp');
                    nodeTextElem.classList.add('nodeDown');
                }
            }
        } catch (error) {
            handleFetchError(error);
        }
    }
    document.querySelector('.loader').classList.add('hidden');
}, 10000);

// Helper function to handle fetch errors
function handleFetchError(error) {
    outputElem.innerText = error;
    console.error(error);
}
