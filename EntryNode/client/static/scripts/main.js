const outputElem = document.getElementById('output');

async function makeHealthCheck() {
    console.log('Making health check');

    try {
        const response = await fetch('/health', { method: 'GET' });
        await handleResponse(response);
    } catch (error) {
        handleFetchError(error);
    }
}

async function getData() {
    const keyInp = document.querySelector('#getPanel > input[name=key]');
    const key = keyInp.value;
    console.log('Getting data with key:', key);

    try {
        const response = await fetch(`/data?key=${key}`, { method: 'GET' });
        await handleResponse(response);
    } catch (error) {
        handleFetchError(error);
    }
}

async function addData() {
    console.log('Adding data');

    const keyInp = document.querySelector('#addPanel > input[name=key]');
    const valueInp = document.querySelector('#addPanel > input[name=value]');
    const data = {
        key: keyInp.value,
        value: valueInp.value,
    };

    console.log('Data:', data);

    try {
        const response = await fetch('/data', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data),
        });
        await handleResponse(response);
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
            console.log(body);

            const listContainer = document.getElementById('nodeList');

            body.nodes.forEach((node) => {
                const nodeElem = document.createElement('div');
                nodeElem.innerText = node;
                nodeElem.onclick = async () => {
                    // Get health of node
                    try {
                        const healthRes = await fetch(`/health?node=${node}`, {
                            method: 'GET',
                        });
                        handleResponse(healthRes);
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
    const healthBtn = document.createElement('button');
    const cycleHealthBtn = document.createElement('button');

    nodeText.innerText = node;
    healthBtn.innerText = 'Health';
    cycleHealthBtn.innerText = 'Cycle Health';

    nodeElem.appendChild(nodeText);
    nodeElem.appendChild(healthBtn);
    nodeElem.appendChild(cycleHealthBtn);

    healthBtn.onclick = makeHealthCheck;

    cycleHealthBtn.onclick = makeCycleHealthCheck;

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

// Helper function to handle fetch errors
function handleFetchError(error) {
    outputElem.innerText = error;
    console.error(error);
}
