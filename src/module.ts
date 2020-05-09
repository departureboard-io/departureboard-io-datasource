import { DataSourcePlugin } from '@grafana/data';
import { DataSource } from './DataSource';
import { ConfigEditor } from './ConfigEditor';
import { QueryEditor } from './QueryEditor';
import { DepartureBoardIOQuery, DepartureBoardIOOptions } from './types';

export const plugin = new DataSourcePlugin<DataSource, DepartureBoardIOQuery, DepartureBoardIOOptions>(DataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor);
