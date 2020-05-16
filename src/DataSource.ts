import { DataSourceInstanceSettings } from '@grafana/data';
import { DataSourceWithBackend, getTemplateSrv } from '@grafana/runtime';

import { DepartureBoardIOQuery, DepartureBoardIOOptions } from './types';

export class DataSource extends DataSourceWithBackend<DepartureBoardIOQuery, DepartureBoardIOOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<DepartureBoardIOOptions>) {
    super(instanceSettings);
  }

  // Support template variables for stationCRS.
  applyTemplateVariables(query: DepartureBoardIOQuery) {
    const templateSrv = getTemplateSrv();
    return {
      ...query,
      stationCRS: query.stationCRS ? templateSrv.replace(query.stationCRS) : '',
    };
  }
}
