import { DataQuery, DataSourceJsonData } from '@grafana/data';

export interface MyQuery extends DataQuery {
  image_keys?: string;
}

export const DEFAULT_QUERY: Partial<MyQuery> = {
  image_keys: "33884862/packs/c0de1446-0000-feed-f00d-5a1ad2c0ffee/zone1_20230223-175257_DEF.png,33884862/packs/c0de1446-0000-feed-f00d-5a1ad2c0ffee/zone2_20230223-175348_DEF.png",
};

/**
 * These are options configured for each DataSource instance
 */
export interface MyDataSourceOptions extends DataSourceJsonData {
  bucket?: string;
}

/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */
export interface MySecureJsonData {
  aws_access_key_id?: string;
  aws_secret_access_key?: string;
}
