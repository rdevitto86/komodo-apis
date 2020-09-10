/** @see https://www.youtube.com/watch?v=L72fhGm1tfE */

// imports
const express = require('express');
// const path = require('path');

// express server config
const service = express();

service.get('/api/test', (req, res) => res.json({
    a: {
        aa: 1.1,
        ab: 1.2,
        },
    b: {
        ba: 2.1,
        bb: 2.2,
    },
}));

// service.use(express.static(path.join(__dirname, 'public')));

service.use((req, res, next) => {
    const apiPath = `${req.protocol}://${req.get('host')}${req.originalUrl}`;
    console.log(`[INFO] ${apiPath}`);
    next();
});

// server port
const PORT = process.env.PORT || 5000;

//set server port
service.listen(PORT, () => console.log(`server started using port ${PORT}`));
