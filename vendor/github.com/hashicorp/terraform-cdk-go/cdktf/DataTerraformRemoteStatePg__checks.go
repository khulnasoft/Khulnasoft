// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !no_runtime_type_checking

package cdktf

import (
	"fmt"

	_jsii_ "github.com/aws/jsii-runtime-go/runtime"

	"github.com/aws/constructs-go/constructs/v10"
)

func (d *jsiiProxy_DataTerraformRemoteStatePg) validateAddOverrideParameters(path *string, value interface{}) error {
	if path == nil {
		return fmt.Errorf("parameter path is required, but nil was provided")
	}

	if value == nil {
		return fmt.Errorf("parameter value is required, but nil was provided")
	}

	return nil
}

func (d *jsiiProxy_DataTerraformRemoteStatePg) validateGetParameters(output *string) error {
	if output == nil {
		return fmt.Errorf("parameter output is required, but nil was provided")
	}

	return nil
}

func (d *jsiiProxy_DataTerraformRemoteStatePg) validateGetBooleanParameters(output *string) error {
	if output == nil {
		return fmt.Errorf("parameter output is required, but nil was provided")
	}

	return nil
}

func (d *jsiiProxy_DataTerraformRemoteStatePg) validateGetListParameters(output *string) error {
	if output == nil {
		return fmt.Errorf("parameter output is required, but nil was provided")
	}

	return nil
}

func (d *jsiiProxy_DataTerraformRemoteStatePg) validateGetNumberParameters(output *string) error {
	if output == nil {
		return fmt.Errorf("parameter output is required, but nil was provided")
	}

	return nil
}

func (d *jsiiProxy_DataTerraformRemoteStatePg) validateGetStringParameters(output *string) error {
	if output == nil {
		return fmt.Errorf("parameter output is required, but nil was provided")
	}

	return nil
}

func (d *jsiiProxy_DataTerraformRemoteStatePg) validateOverrideLogicalIdParameters(newLogicalId *string) error {
	if newLogicalId == nil {
		return fmt.Errorf("parameter newLogicalId is required, but nil was provided")
	}

	return nil
}

func validateDataTerraformRemoteStatePg_IsConstructParameters(x interface{}) error {
	if x == nil {
		return fmt.Errorf("parameter x is required, but nil was provided")
	}

	return nil
}

func validateDataTerraformRemoteStatePg_IsTerraformElementParameters(x interface{}) error {
	if x == nil {
		return fmt.Errorf("parameter x is required, but nil was provided")
	}

	return nil
}

func validateNewDataTerraformRemoteStatePgParameters(scope constructs.Construct, id *string, config *DataTerraformRemoteStatePgConfig) error {
	if scope == nil {
		return fmt.Errorf("parameter scope is required, but nil was provided")
	}

	if id == nil {
		return fmt.Errorf("parameter id is required, but nil was provided")
	}

	if config == nil {
		return fmt.Errorf("parameter config is required, but nil was provided")
	}
	if err := _jsii_.ValidateStruct(config, func() string { return "parameter config" }); err != nil {
		return err
	}

	return nil
}

