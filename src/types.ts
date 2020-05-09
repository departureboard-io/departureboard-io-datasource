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
  endpoint?: string;
  stationCRS?: string;
}

export const defaultQuery: Partial<DepartureBoardIOQuery> = {
  endpoint: 'getDeparturesByCRS',
};
