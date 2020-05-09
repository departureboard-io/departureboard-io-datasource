import defaults from 'lodash/defaults';

import React, { ChangeEvent, PureComponent } from 'react';
import { LegacyForms } from '@grafana/ui';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from './DataSource';
import { defaultQuery, DepartureBoardIOOptions, DepartureBoardIOQuery } from './types';

const { FormField } = LegacyForms;

type Props = QueryEditorProps<DataSource, DepartureBoardIOQuery, DepartureBoardIOOptions>;

export class QueryEditor extends PureComponent<Props> {
  onQueryTextChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, stationCRS: event.target.value });
    onRunQuery();
  };

  render() {
    const query = defaults(this.props.query, defaultQuery);
    const { stationCRS } = query;

    return (
      <div className="gf-form">
        <FormField
          labelWidth={10}
          value={stationCRS || ''}
          onChange={this.onQueryTextChange}
          placeholder="PAD"
          label="CRS code for station"
          tooltip=""
        />
      </div>
    );
  }
}
