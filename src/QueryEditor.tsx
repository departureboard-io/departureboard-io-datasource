import defaults from 'lodash/defaults';

import React, { ChangeEvent, PureComponent } from 'react';
import { LegacyForms, Checkbox } from '@grafana/ui';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from './DataSource';
import { defaultQuery, DepartureBoardIOOptions, DepartureBoardIOQuery } from './types';

const { FormField } = LegacyForms;

type Props = QueryEditorProps<DataSource, DepartureBoardIOQuery, DepartureBoardIOOptions>;

export class QueryEditor extends PureComponent<Props> {
  constructor(props: Props) {
    super(props);
  }
  onCheckboxChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, [event.target.name]: event.target.checked });
    onRunQuery();
  };

  onStationCRSChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, stationCRS: event.target.value });
    onRunQuery();
  };

  render() {
    const query = defaults(this.props.query, defaultQuery);
    const { arrivals, departures, serviceDetails, stationCRS } = query;

    return (
      <div className="gf-form-group">
        <div className="gf-form">
          <FormField
            labelWidth={10}
            value={stationCRS || ''}
            onChange={this.onStationCRSChange}
            placeholder="PAD"
            label="CRS code for station"
            tooltip=""
          />
        </div>
        <div className="gf-form">
          <Checkbox name="serviceDetails" label="Include service details?" width={20} value={serviceDetails} onChange={this.onCheckboxChange} />
          <Checkbox name="departures" label="Include departures?" width={20} value={departures} onChange={this.onCheckboxChange} />
          <Checkbox name="arrivals" label="Include arrivals?" width={20} value={arrivals} onChange={this.onCheckboxChange} />
        </div>
      </div>
    );
  }
}
