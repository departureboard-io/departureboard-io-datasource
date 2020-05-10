import defaults from 'lodash/defaults';

import React, { ChangeEvent, PureComponent } from 'react';
import { LegacyForms, Checkbox } from '@grafana/ui';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { DataSource } from './DataSource';
import { defaultQuery, DepartureBoardIOOptions, DepartureBoardIOQuery } from './types';

const { FormField } = LegacyForms;

type Props = QueryEditorProps<DataSource, DepartureBoardIOQuery, DepartureBoardIOOptions>;

export class QueryEditor extends PureComponent<Props> {
  constructor(props: Props) {
    super(props);
  }
  onArrivalsChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, arrivals: event.target.checked });
    onRunQuery();
  };

  onDeparturesChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, departures: event.target.checked });
    onRunQuery();
  };

  onStationCRSChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, stationCRS: event.target.value });
    onRunQuery();
  };

  onQueryTypeChange = (selected: SelectableValue<string>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, queryType: selected.value });
  };

  render() {
    const query = defaults(this.props.query, defaultQuery);
    const { arrivals, departures, stationCRS } = query;

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
          <Checkbox label="Include departures?" width={20} value={departures} onChange={this.onDeparturesChange} />
        </div>
        <div className="gf-form">
          <Checkbox label="Include arrivals?" width={20} value={arrivals} onChange={this.onArrivalsChange} />
        </div>
      </div>
    );
  }
}
