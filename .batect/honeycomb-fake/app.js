const express = require('express');
const app = express();
const port = 3000;

const datasets = {};

app.use(express.json());

app.get('/ping', (req, res) => res.send('Honeycomb fake is running'));

app.post('/1/events/:datasetName', (req, res) => {
    if (req.headers['content-type'] !== 'application/json') {
        return res.status(400).send('missing or incorrect Content-Type header');
    }

    const datasetName = req.params.datasetName;

    if (req.headers['x-honeycomb-team'] === undefined) {
        return res.status(400).send('missing X-Honeycomb-Team header');
    }

    if (datasets[datasetName] === undefined) {
        datasets[datasetName] = [];
    }

    let time = req.headers['x-honeycomb-event-time'];

    if (time === undefined) {
        time = (new Date()).toISOString();
    }

    console.log('Received', req.body);

    const event = {
        time,
        data: req.body
    };

    datasets[datasetName].push(event);

    return res.status(200).send();
});

app.get('/fake/events', (req, res) => {
    return res.json(datasets);
});

app.get('/fake/events/:datasetName', (req, res) => {
    const datasetName = req.params.datasetName;
    const events = datasets[datasetName];

    if (events === undefined) {
        return res.status(404).send(`no events for dataset '${datasetName}'`);
    }

    return res.json(events);
});

app.get('/fake/events/:datasetName/:index', (req, res) => {
    const datasetName = req.params.datasetName;
    const index = req.params.index;
    const events = datasets[datasetName];

    if (events === undefined) {
        return res.status(404).send(`no events for dataset '${datasetName}'`);
    }

    if (index >= events.length) {
        return res.status(404).send(`no event ${index} for dataset '${datasetName}'`);
    }

    return res.json(events[index]);
});

app.listen(port, () => console.log(`Honeycomb fake listening on port ${port}.`));
