import React, { FormEvent } from 'react';
import { AutoSizeInput, InlineField } from '@grafana/ui';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from '../datasource';
import { MyDataSourceOptions, MyQuery } from '../types';

type Props = QueryEditorProps<DataSource, MyQuery, MyDataSourceOptions>;

export function QueryEditor({ query, onChange, onRunQuery }: Props) {

  const onKeysChange = (event: FormEvent<HTMLInputElement>) => {
    console.log(query)
    console.log(event.currentTarget.value)
    onChange({ ...query, image_keys: event.currentTarget.value});
    // executes the query
    onRunQuery();
  };

  const { image_keys } = query;

  return (
    <div className="gf-form">
      <InlineField label="Image Keys">
        <AutoSizeInput onCommitChange={onKeysChange} value={image_keys} width={20} />
      </InlineField>
    </div>
  );
}
