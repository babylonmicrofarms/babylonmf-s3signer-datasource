import React, { ChangeEvent } from 'react';
import { Alert, InlineField, Input, SecretInput } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { MyDataSourceOptions, MySecureJsonData } from '../types';

interface Props extends DataSourcePluginOptionsEditorProps<MyDataSourceOptions> {}

export function ConfigEditor(props: Props) {
  const { onOptionsChange, options } = props;
  const onBucketChange = (event: ChangeEvent<HTMLInputElement>) => {
    const jsonData = {
      ...options.jsonData,
      bucket: event.target.value,
    };
    onOptionsChange({ ...options, jsonData });
  };

  // Secure field (only sent to the backend)
   const onKeyIdChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      secureJsonData: {
        aws_access_key_id: event.target.value,
      },
    });
  };

  const onSecretKeyChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      secureJsonData: {
        aws_secret_access_key: event.target.value,
      },
    });
  };

  const onResetKeyId = () => {
    onOptionsChange({
      ...options,
      secureJsonFields: {
        ...options.secureJsonFields,
        aws_access_key_id: false,
      },
      secureJsonData: {
        ...options.secureJsonData,
        aws_access_key_id: '',
      },
    });
  };

  const onResetSecretKey = () => {
    onOptionsChange({
      ...options,
      secureJsonFields: {
        ...options.secureJsonFields,
        aws_secret_access_key: false,
      },
      secureJsonData: {
        ...options.secureJsonData,
        aws_secret_access_key: '',
      },
    });
  };

  const { jsonData, secureJsonFields } = options;
  const secureJsonData = (options.secureJsonData || {}) as MySecureJsonData;

  return (
    <div className="gf-form-group">
      <Alert title='Credentials Note' severity='info'>
        Even though this plugin is using your aws credentials it does not communicate with AWS to ensure that your credentials are correct or that they have access to the bucket you are trying to use. Double check those things if something isn&apos;t working.
      </Alert>
      <InlineField label="bucket" labelWidth={12}>
        <Input
          onChange={onBucketChange}
          value={jsonData.bucket || ''}
          placeholder="s3 bucket that you want to sign objects from"
          width={40}
        />
      </InlineField>
      <InlineField label="AWS access key id" labelWidth={24}>
        <SecretInput
          isConfigured={(secureJsonFields && secureJsonFields.aws_access_key_id) as boolean}
          value={secureJsonData.aws_access_key_id || ''}
          placeholder="secure json field (backend only)"
          width={40}
          onReset={onResetKeyId}
          onChange={onKeyIdChange}
        />
      </InlineField>
      <InlineField label="AWS Secret Access Key" labelWidth={24}>
        <SecretInput
          isConfigured={(secureJsonFields && secureJsonFields.aws_secret_access_key) as boolean}
          value={secureJsonData.aws_secret_access_key || ''}
          placeholder="secure json field (backend only)"
          width={40}
          onReset={onResetSecretKey}
          onChange={onSecretKeyChange}
        />
      </InlineField>
    </div>
  );
}
