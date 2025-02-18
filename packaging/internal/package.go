/*
 * Copyright ©1998-2022 by Richard A. Wilkes. All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, version 2.0. If a copy of the MPL was not distributed with
 * this file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 * This Source Code Form is "Incompatible With Secondary Licenses", as
 * defined by the Mozilla Public License, version 2.0.
 */

package internal

import (
	"github.com/richardwilkes/gcs/v5/early"
	"github.com/richardwilkes/gcs/v5/ux"
)

// Package performs the platform-specific packaging for GCS.
func Package() error {
	early.Configure()
	if err := loadBaseImages(); err != nil {
		return err
	}
	ux.RegisterExternalFileTypes()
	ux.RegisterGCSFileTypes()
	return platformPackage()
}
