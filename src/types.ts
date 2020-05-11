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
  serviceDetails?: boolean;
  filterCRS?: string;
}

export const defaultQuery: Partial<DepartureBoardIOQuery> = {
  departures: true,
  arrivals: false,
  serviceDetails: false,
};
