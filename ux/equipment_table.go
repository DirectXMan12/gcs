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

package ux

import (
	"os"
	"path/filepath"

	"github.com/richardwilkes/gcs/v5/model"
	"github.com/richardwilkes/unison"
)

type equipmentListProvider struct {
	carried []*model.Equipment
	other   []*model.Equipment
}

func (p *equipmentListProvider) Entity() *model.Entity {
	return nil
}

func (p *equipmentListProvider) CarriedEquipmentList() []*model.Equipment {
	return p.carried
}

func (p *equipmentListProvider) SetCarriedEquipmentList(list []*model.Equipment) {
	p.carried = list
}

func (p *equipmentListProvider) OtherEquipmentList() []*model.Equipment {
	return p.other
}

func (p *equipmentListProvider) SetOtherEquipmentList(list []*model.Equipment) {
	p.other = list
}

// NewEquipmentTableDockableFromFile loads a list of equipment from a file and creates a new unison.Dockable for them.
func NewEquipmentTableDockableFromFile(filePath string) (unison.Dockable, error) {
	equipment, err := model.NewEquipmentFromFile(os.DirFS(filepath.Dir(filePath)), filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	d := NewEquipmentTableDockable(filePath, equipment)
	d.needsSaveAsPrompt = false
	return d, nil
}

// NewEquipmentTableDockable creates a new unison.Dockable for equipment list files.
func NewEquipmentTableDockable(filePath string, equipment []*model.Equipment) *TableDockable[*model.Equipment] {
	provider := &equipmentListProvider{other: equipment}
	d := NewTableDockable(filePath, model.EquipmentExt, NewEquipmentProvider(provider, false, false),
		func(path string) error { return model.SaveEquipment(provider.OtherEquipmentList(), path) },
		NewOtherEquipmentItemID, NewOtherEquipmentContainerItemID)
	d.InstallCmdHandlers(ConvertToContainerItemID,
		func(_ any) bool { return CanConvertToContainer(d.table) },
		func(_ any) { ConvertToContainer(d, d.table) })
	return d
}
