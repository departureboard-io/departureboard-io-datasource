import { DataSourceInstanceSettings } from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';

import { DepartureBoardIOQuery, DepartureBoardIOOptions } from './types';

export class DataSource extends DataSourceWithBackend<DepartureBoardIOQuery, DepartureBoardIOOptions> {
  /** @ngInject */
  constructor(instanceSettings: DataSourceInstanceSettings<DepartureBoardIOOptions>, private templateSrv: any) {
    super(instanceSettings);
  }

  // Support template variables for stationCRS.
  applyTemplateVariables(query: DepartureBoardIOQuery) {
    return {
      ...query,
      stationCRS: this.templateSrv.replace(query.stationCRS),
    };
  }
}
