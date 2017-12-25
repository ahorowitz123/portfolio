import fs from 'fs';
import { printSchema } from 'graphql/utilities';
import path from 'path';

import Schema from './schema';

fs.writeFileSync(path.join(__dirname, './schema.graphql'), printSchema(Schema));
