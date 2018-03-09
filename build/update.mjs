import fs from 'fs';
import util from 'util';
import path from 'path';
import pkg from '../package.json';

const basePath = path.dirname(process.argv[1]);

const pkgFile = path.resolve(basePath, '../package.json');

const writeFile = util.promisify(fs.writeFile);
const verions = ['major', 'minor', 'patch'];

const version = process.argv[2] || 'patch';

const foundIndex = verions.indexOf(version);
if (foundIndex === -1) {
  console.error('the update version type is wrong');
} else {
  const versionList = pkg.version.split('.');
  versionList[foundIndex] = Number.parseInt(versionList[foundIndex], 10) + 1;
  const newVersion = versionList.join('.');
  pkg.version = newVersion;
  writeFile(pkgFile, JSON.stringify(pkg, null, 2))
    .then(() => {
      console.info(`update version success, new version is:${newVersion}`);
    })
    .catch(err => {
      console.error(`update version fail, ${err.message}`);
    });
}
