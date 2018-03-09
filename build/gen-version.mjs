import fs from 'fs';

const file = './assets/version';
fs.writeFileSync(file, new Date().toISOString());
