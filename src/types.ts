import { DataQuery, DataSourceJsonData } from '@grafana/data';

export interface DepartureBoardIOOptions extends DataSourceJsonData {
  apiEndpoint: string;
}

export const defaultOptions: DepartureBoardIOOptions = {
  apiEndpoint: 'https://api.departureboard.io/api/v2.0',
};

export interface DepartureBoardIOSecureJSONData {
  apiKey: string;
}

export interface DepartureBoardIOQuery extends DataQuery {
  stationCRS?: string;
  departures?: boolean;
  arrivals?: boolean;
}

export const defaultQuery: Partial<DepartureBoardIOQuery> = {
  queryType: 'advanced',
  departures: true,
  arrivals: false,
};
