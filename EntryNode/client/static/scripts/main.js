const outputElem = document.getElementById('output');

async function makeHealthCheck() {
    console.log('Making health check');

    // Make the check
    const response = await fetch('/health', {
        method: 'GET',
    });

    const body = await response.json();

    outputElem.innerText = JSON.stringify(body);
    console.log(body);
}

async function getData() {
    console.log('Getting data');

    // Make the check
    const response = await fetch('/data', {
        method: 'GET',
    });

    const body = await response.json();

    outputElem.innerText = JSON.stringify(body);
    console.log(body);
}

async function addData() {
    console.log('Adding data');

    // Make the check
    const response = await fetch('/data', {
        method: 'POST',
    });

    const body = await response.json();

    outputElem.innerText = JSON.stringify(body);
    console.log(body);
}
