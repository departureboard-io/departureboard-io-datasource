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

  onFormFieldChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, [event.target.name]: event.target.value });
    onRunQuery();
  };

  render() {
    const query = defaults(this.props.query, defaultQuery);
    const { arrivals, departures, filterCRS, serviceDetails, stationCRS } = query;

    return (
      <div className="gf-form-group">
        <div className="gf-form">
          <FormField
            name="stationCRS"
            label="CRS code for station"
            labelWidth={10}
            value={stationCRS || ''}
            onChange={this.onFormFieldChange}
            placeholder="PAD"
          />
          <FormField
            name="filterCRS"
            label="CRS code for filter station"
            labelWidth={15}
            value={filterCRS || ''}
            onChange={this.onFormFieldChange}
            placeholder="HAY"
            tooltip="Will only show services that include the filter station as the origin or destination."
          />
        </div>
        <div className="gf-form">
          <Checkbox name="serviceDetails" label="Include service details?" width={20} value={serviceDetails} onChange={this.onCheckboxChange} />
        </div>
        <div className="gf-form">
          <Checkbox name="departures" label="Include departures?" width={20} value={departures} onChange={this.onCheckboxChange} />
        </div>
        <div className="gf-form">
          <Checkbox name="arrivals" label="Include arrivals?" width={20} value={arrivals} onChange={this.onCheckboxChange} />
        </div>
      </div>
    );
  }
}
