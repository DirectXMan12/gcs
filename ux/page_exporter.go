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
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/richardwilkes/gcs/v5/model"
	"github.com/richardwilkes/toolbox"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/xmath"
	"github.com/richardwilkes/unison"
)

var _ Rebuildable = &pageExporter{}

const pageKey = "pageKey"

type pageExporter struct {
	unison.Panel
	entity      *model.Entity
	targetMgr   *TargetMgr
	pages       []*Page
	currentPage int
}

func newPageExporter(entity *model.Entity) *pageExporter {
	p := &pageExporter{entity: entity}
	p.targetMgr = NewTargetMgr(p)
	pageSize := p.PageSize()
	r := unison.Rect{Size: pageSize}
	page, _ := createPageTopBlock(entity, p.targetMgr)
	p.AddChild(page)
	p.pages = append(p.pages, page)
	for _, col := range entity.SheetSettings.BlockLayout.ByRow() {
		startAt := make(map[string]int)
		for {
			rowPanel := unison.NewPanel()
			rowPanel.SetLayoutData(&unison.FlexLayoutData{
				HAlign: unison.FillAlignment,
				HGrab:  true,
			})
			for _, c := range col {
				switch c {
				case model.BlockLayoutReactionsKey:
					addRowPanel(rowPanel, NewReactionsPageList(entity), model.BlockLayoutReactionsKey, startAt)
				case model.BlockLayoutConditionalModifiersKey:
					addRowPanel(rowPanel, NewConditionalModifiersPageList(entity), model.BlockLayoutConditionalModifiersKey, startAt)
				case model.BlockLayoutMeleeKey:
					addRowPanel(rowPanel, NewMeleeWeaponsPageList(entity), model.BlockLayoutMeleeKey, startAt)
				case model.BlockLayoutRangedKey:
					addRowPanel(rowPanel, NewRangedWeaponsPageList(entity), model.BlockLayoutRangedKey, startAt)
				case model.BlockLayoutTraitsKey:
					addRowPanel(rowPanel, NewTraitsPageList(p, entity), model.BlockLayoutTraitsKey, startAt)
				case model.BlockLayoutSkillsKey:
					addRowPanel(rowPanel, NewSkillsPageList(p, entity), model.BlockLayoutSkillsKey, startAt)
				case model.BlockLayoutSpellsKey:
					addRowPanel(rowPanel, NewSpellsPageList(p, entity), model.BlockLayoutSpellsKey, startAt)
				case model.BlockLayoutEquipmentKey:
					addRowPanel(rowPanel, NewCarriedEquipmentPageList(p, entity), model.BlockLayoutEquipmentKey, startAt)
				case model.BlockLayoutOtherEquipmentKey:
					addRowPanel(rowPanel, NewOtherEquipmentPageList(p, entity), model.BlockLayoutOtherEquipmentKey, startAt)
				case model.BlockLayoutNotesKey:
					addRowPanel(rowPanel, NewNotesPageList(p, entity), model.BlockLayoutNotesKey, startAt)
				}
			}
			children := rowPanel.Children()
			if len(children) == 0 {
				break
			}
			rowPanel.SetLayout(&unison.FlexLayout{
				Columns:      len(children),
				HSpacing:     1,
				HAlign:       unison.FillAlignment,
				EqualColumns: true,
			})
			page.AddChild(rowPanel)
			page.SetFrameRect(r)
			page.MarkForLayoutRecursively()
			page.ValidateLayout()
			_, pref, _ := page.Sizes(unison.Size{Width: r.Width})
			excess := pref.Height - pageSize.Height
			if excess <= 0 {
				break // Not extending off the page, so move to the next row
			}
			remaining := (pageSize.Height - page.insets().Bottom) - rowPanel.FrameRect().Y
			startNewPage := false
			data := make([]*pageState, len(children))
			for i, child := range children {
				data[i] = newPageState(child)
				if remaining < data[i].minimum {
					startNewPage = true
				}
			}
			if startNewPage {
				// At least one of the columns can't fit at least the header plus the next row, so start a new page
				page.RemoveChild(rowPanel)
				page = NewPage(entity)
				p.AddChild(page)
				p.pages = append(p.pages, page)
				page.AddChild(rowPanel)
				page.SetFrameRect(r)
				page.MarkForLayoutRecursively()
				page.ValidateLayout()
				_, pref, _ = page.Sizes(unison.Size{Width: r.Width})
				if excess = pref.Height - pageSize.Height; excess <= 0 {
					break // Not extending off the page, so move to the next row
				}
				remaining = (pageSize.Height - page.insets().Bottom) - rowPanel.FrameRect().Y
			}
			startNewPage = false
			for _, one := range data {
				allowed := remaining - one.overhead
				start, endBefore := one.helper.CurrentDrawRowRange()
				startAt[one.key()] = len(one.heights) // Assume all remaining fit
				for i := start; i < endBefore; i++ {
					allowed -= one.heights[i] + 1
					if allowed < 0 {
						// No more fit, mark it
						one.helper.SetDrawRowRange(start, xmath.Max(i, start+1))
						if i == start {
							// I have to guard against the case where a single row is so large it can't fit on a single
							// page on its own. In this case, I just let it flow off the end and drop that extra
							// content.
							//
							// TODO: In the future, see if I can do sub-row partitioning.
							i = start + 1
						}
						startAt[one.key()] = i
						startNewPage = true
						break
					}
				}
			}
			if startNewPage {
				// We've filled the page, so add another
				page = NewPage(entity)
				p.AddChild(page)
				p.pages = append(p.pages, page)
			}
		}
	}
	for _, page = range p.pages {
		page.Force = true
		page.SetFrameRect(r)
		page.MarkForLayoutRecursively()
		page.ValidateLayout()
	}
	return p
}

type pageHelper interface {
	OverheadHeight() float32
	RowHeights() []float32
	CurrentDrawRowRange() (start, endBefore int)
	SetDrawRowRange(start, endBefore int)
}

type pageState struct {
	child    *unison.Panel
	helper   pageHelper
	current  float32
	overhead float32
	minimum  float32
	heights  []float32
}

func newPageState(child unison.Paneler) *pageState {
	panel := child.AsPanel()
	helper := panel.Self.(pageHelper) //nolint:errcheck // The only things used with this are pageHelper-compliant
	state := &pageState{
		child:    panel,
		helper:   helper,
		current:  panel.FrameRect().Height,
		overhead: helper.OverheadHeight(),
		heights:  helper.RowHeights(),
	}
	state.minimum = state.overhead
	start, _ := state.helper.CurrentDrawRowRange()
	if len(state.heights) > start {
		state.minimum += state.heights[start] + 1
	}
	return state
}

func (s *pageState) key() string {
	return s.child.ClientData()[pageKey].(string)
}

func addRowPanel[T model.NodeTypes](rowPanel *unison.Panel, list *PageList[T], key string, startAtMap map[string]int) {
	list.ClientData()[pageKey] = key
	count := list.RowCount()
	startAt := startAtMap[key]
	if count > startAt {
		list.SetDrawRowRange(startAt, count)
		rowPanel.AddChild(list)
	}
}

func (p *pageExporter) exportAsPDFBytes() ([]byte, error) {
	stream := unison.NewMemoryStream()
	defer stream.Close()
	if err := p.exportAsPDF(stream); err != nil {
		return nil, err
	}
	return stream.Bytes(), nil
}

func (p *pageExporter) exportAsPDFFile(filePath string) error {
	if err := os.Remove(filePath); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return errs.Wrap(err)
	}
	stream, err := unison.NewFileStream(filePath)
	if err != nil {
		return err
	}
	defer stream.Close()
	return p.exportAsPDF(stream)
}

func (p *pageExporter) exportAsPDF(stream unison.Stream) error {
	savedColorMode := p.saveTheme()
	defer p.restoreTheme(savedColorMode)
	if err := unison.CreatePDF(stream, &unison.PDFMetaData{
		Title:           p.entity.Profile.Name,
		Author:          toolbox.CurrentUserName(),
		Subject:         p.entity.Profile.Name,
		Keywords:        "GCS Character Sheet",
		Creator:         "GCS",
		RasterDPI:       300,
		EncodingQuality: 101,
	}, p); err != nil {
		return err
	}
	return nil
}

func (p *pageExporter) exportAsPNGs(filePathBase string) error {
	return p.exportAsImages(filePathBase, ".png", func(img *unison.Image) ([]byte, error) {
		return img.ToPNG()
	})
}

func (p *pageExporter) exportAsWEBPs(filePathBase string) error {
	return p.exportAsImages(filePathBase, ".webp", func(img *unison.Image) ([]byte, error) {
		return img.ToWebp(75)
	})
}

func (p *pageExporter) exportAsJPEGs(filePathBase string) error {
	return p.exportAsImages(filePathBase, ".jpeg", func(img *unison.Image) ([]byte, error) {
		return img.ToJPEG(80)
	})
}

func (p *pageExporter) exportAsImages(filePathBase, extension string, f func(img *unison.Image) ([]byte, error)) error {
	filePathBase = strings.TrimSuffix(filePathBase, extension)
	savedColorMode := p.saveTheme()
	defer p.restoreTheme(savedColorMode)
	resolution := model.GlobalSettings().General.ImageResolution
	pageNumber := 1
	for p.HasPage(pageNumber) {
		size := p.PageSize()
		var drawErr error
		img, err := unison.NewImageFromDrawing(int(size.Width), int(size.Height), resolution, func(c *unison.Canvas) {
			drawErr = p.DrawPage(c, pageNumber)
		})
		if err != nil {
			return err
		}
		if drawErr != nil {
			return drawErr
		}
		var data []byte
		if data, err = f(img); err != nil {
			return err
		}
		if err = os.WriteFile(fmt.Sprintf("%s-%d%s", filePathBase, pageNumber, extension), data, 0o640); err != nil {
			return err
		}
		pageNumber++
	}
	return nil
}

func (p *pageExporter) saveTheme() unison.ColorMode {
	savedColorMode := unison.CurrentColorMode()
	unison.SetColorMode(unison.LightColorMode)
	unison.ThemeChanged()
	unison.RebuildDynamicColors()
	return savedColorMode
}

func (p *pageExporter) restoreTheme(colorMode unison.ColorMode) {
	unison.SetColorMode(colorMode)
	unison.ThemeChanged()
	unison.RebuildDynamicColors()
}

// HasPage implements unison.PageProvider.
func (p *pageExporter) HasPage(pageNumber int) bool {
	p.currentPage = pageNumber
	return pageNumber > 0 && pageNumber <= len(p.pages)
}

// PageSize implements unison.PageProvider.
func (p *pageExporter) PageSize() unison.Size {
	w, h := p.entity.SheetSettings.Page.Orientation.Dimensions(p.entity.SheetSettings.Page.Size.Dimensions())
	return unison.NewSize(w.Pixels(), h.Pixels())
}

// DrawPage implements unison.PageProvider.
func (p *pageExporter) DrawPage(canvas *unison.Canvas, pageNumber int) error {
	p.currentPage = pageNumber
	if pageNumber > 0 && pageNumber <= len(p.pages) {
		page := p.pages[pageNumber-1]
		page.Draw(canvas, page.ContentRect(true))
		return nil
	}
	return errs.New("invalid page number")
}

func (p *pageExporter) String() string {
	return ""
}

func (p *pageExporter) Rebuild(_ bool) {
}
