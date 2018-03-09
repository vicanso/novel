// eslint-disable-next-line
import swaggerJSDoc from 'swagger-jsdoc';
import fs from 'fs';

import pkg from '../package';

const definition = {
  info: {
    title: pkg.name,
    version: pkg.version,
    description: 'SPT文档',
  },
  host: 'aslant.site',
};
const options = {
  swaggerDefinition: definition,
  apis: ['*/*.mjs', '*/*/*.mjs'],
};

const swaggerSpec = swaggerJSDoc(options);
fs.writeFileSync('./assets/api.json', JSON.stringify(swaggerSpec, null, 2));
