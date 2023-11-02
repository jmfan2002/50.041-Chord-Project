const outputElem = document.getElementById('output');

async function makeHealthCheck() {
    console.log('Making health check');

    // Make the check
    try {
        const res = await fetch('/health', {
            method: 'GET',
        });

        if (!res.ok) {
            outputElem.innerText = `Error: ${res.statusText}`;
            console.error(res);
            return;
        }

        const body = await res.json();

        outputElem.innerText = JSON.stringify(body);
        console.log(body);
    } catch (error) {
        outputElem.innerText = error;
        console.error(error);
    }
}

async function getData() {
    const keyInp = document.querySelector('#getPanel > input[name=key]');
    const key = keyInp.value;
    console.log('Getting data with key:', key);

    // Make the get call
    try {
        const res = await fetch(`/data?key=${key}`, {
            method: 'GET',
        });

        if (!res.ok) {
            outputElem.innerText = `Error: ${res.statusText}`;
            console.error(res);
            return;
        }

        const body = await res.json();
        outputElem.innerText = JSON.stringify(body);
        console.log(body);
    } catch (error) {
        outputElem.innerText = error;
        console.error(error);
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
        const res = await fetch('/data', {
            method: 'POST',
            body: data,
        });

        if (!res.ok) {
            outputElem.innerText = `Error: ${res.statusText}`;
            console.error(res);
            return;
        }

        const body = await res.json();

        outputElem.innerText = JSON.stringify(body);

        console.log(body);
    } catch (error) {
        outputElem.innerText = error;
        console.error(error);
    }
}
