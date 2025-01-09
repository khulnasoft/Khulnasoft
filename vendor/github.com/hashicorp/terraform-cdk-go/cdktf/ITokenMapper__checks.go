// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !no_runtime_type_checking

package cdktf

import (
	"fmt"
)

func (i *jsiiProxy_ITokenMapper) validateMapTokenParameters(t IResolvable) error {
	if t == nil {
		return fmt.Errorf("parameter t is required, but nil was provided")
	}

	return nil
}

