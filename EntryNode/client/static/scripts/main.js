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
