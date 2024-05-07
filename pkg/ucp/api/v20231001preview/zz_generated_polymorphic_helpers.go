//go:build go1.18
// +build go1.18

// Licensed under the Apache License, Version 2.0 . See LICENSE in the repository root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package v20231001preview

import "encoding/json"

func unmarshalAwsCredentialPropertiesClassification(rawMsg json.RawMessage) (AwsCredentialPropertiesClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var m map[string]any
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b AwsCredentialPropertiesClassification
	switch m["kind"] {
	case string(AWSCredentialKindAccessKey):
		b = &AwsAccessKeyCredentialProperties{}
	default:
		b = &AwsCredentialProperties{}
	}
	if err := json.Unmarshal(rawMsg, b); err != nil {
		return nil, err
	}
	return b, nil
}

func unmarshalAzureCredentialPropertiesClassification(rawMsg json.RawMessage) (AzureCredentialPropertiesClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var m map[string]any
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b AzureCredentialPropertiesClassification
	switch m["kind"] {
	case string(AzureCredentialKindServicePrincipal):
		b = &AzureServicePrincipalProperties{}
	case string(AzureCredentialKindWorkloadIdentity):
		b = &AzureWorkloadIdentityProperties{}
	default:
		b = &AzureCredentialProperties{}
	}
	if err := json.Unmarshal(rawMsg, b); err != nil {
		return nil, err
	}
	return b, nil
}

func unmarshalCredentialStoragePropertiesClassification(rawMsg json.RawMessage) (CredentialStoragePropertiesClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var m map[string]any
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b CredentialStoragePropertiesClassification
	switch m["kind"] {
	case string(CredentialStorageKindInternal):
		b = &InternalCredentialStorageProperties{}
	default:
		b = &CredentialStorageProperties{}
	}
	if err := json.Unmarshal(rawMsg, b); err != nil {
		return nil, err
	}
	return b, nil
}

