/*
 * Copyright (c) 1998-2020 by Richard A. Wilkes. All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, version 2.0. If a copy of the MPL was not distributed with
 * this file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 * This Source Code Form is "Incompatible With Secondary Licenses", as
 * defined by the Mozilla Public License, version 2.0.
 */

package com.trollworks.gcs.notes;

import com.trollworks.gcs.character.GURPSCharacter;
import com.trollworks.gcs.datafile.DataFile;
import com.trollworks.gcs.datafile.LoadState;
import com.trollworks.gcs.io.xml.XMLReader;
import com.trollworks.gcs.io.xml.XMLWriter;
import com.trollworks.gcs.ui.RetinaIcon;
import com.trollworks.gcs.ui.image.Images;
import com.trollworks.gcs.ui.widget.outline.Column;
import com.trollworks.gcs.ui.widget.outline.ListRow;
import com.trollworks.gcs.ui.widget.outline.RowEditor;
import com.trollworks.gcs.utility.I18n;
import com.trollworks.gcs.utility.text.Text;

import java.io.IOException;
import java.util.Map;
import java.util.Set;

/** A note. */
public class Note extends ListRow {
    private static final int    CURRENT_VERSION    = 1;
    /** The XML tag used for items. */
    public static final  String TAG_NOTE           = "note";
    /** The XML tag used for containers. */
    public static final  String TAG_NOTE_CONTAINER = "note_container";
    private static final String TAG_TEXT           = "text";
    /** The prefix used in front of all IDs for the notes. */
    public static final  String PREFIX             = GURPSCharacter.CHARACTER_PREFIX + "note.";
    /** The field ID for text changes. */
    public static final  String ID_TEXT            = PREFIX + "Text";
    /** The field ID for when the row hierarchy changes. */
    public static final  String ID_LIST_CHANGED    = PREFIX + "ListChanged";
    private              String mText;

    /**
     * Creates a new note.
     *
     * @param dataFile    The data file to associate it with.
     * @param isContainer Whether or not this row allows children.
     */
    public Note(DataFile dataFile, boolean isContainer) {
        super(dataFile, isContainer);
        mText = "";
    }

    /**
     * Creates a clone of an existing note and associates it with the specified data file.
     *
     * @param dataFile The data file to associate it with.
     * @param note     The note to clone.
     * @param deep     Whether or not to clone the children, grandchildren, etc.
     */
    public Note(DataFile dataFile, Note note, boolean deep) {
        super(dataFile, note);
        mText = note.mText;
        if (deep) {
            int count = note.getChildCount();
            for (int i = 0; i < count; i++) {
                addChild(new Note(dataFile, (Note) note.getChild(i), true));
            }
        }
    }

    /**
     * Loads a note and associates it with the specified data file.
     *
     * @param dataFile The data file to associate it with.
     * @param reader   The XML reader to load from.
     * @param state    The {@link LoadState} to use.
     */
    public Note(DataFile dataFile, XMLReader reader, LoadState state) throws IOException {
        this(dataFile, TAG_NOTE_CONTAINER.equals(reader.getName()));
        load(reader, state);
    }

    @Override
    public boolean isEquivalentTo(Object obj) {
        if (obj == this) {
            return true;
        }
        if (obj instanceof Note && super.isEquivalentTo(obj)) {
            return mText.equals(((Note) obj).mText);
        }
        return false;
    }

    @Override
    public String getLocalizedName() {
        return I18n.Text("Note");
    }

    @Override
    public String getListChangedID() {
        return ID_LIST_CHANGED;
    }

    @Override
    public String getXMLTagName() {
        return canHaveChildren() ? TAG_NOTE_CONTAINER : TAG_NOTE;
    }

    @Override
    public int getXMLTagVersion() {
        return CURRENT_VERSION;
    }

    @Override
    public String getRowType() {
        return I18n.Text("Note");
    }

    @Override
    protected void prepareForLoad(LoadState state) {
        super.prepareForLoad(state);
        mText = "";
    }

    @Override
    protected void loadSubElement(XMLReader reader, LoadState state) throws IOException {
        String name = reader.getName();
        if (TAG_TEXT.equals(name)) {
            mText = Text.standardizeLineEndings(reader.readText());
        } else if (!state.mForUndo && (TAG_NOTE.equals(name) || TAG_NOTE_CONTAINER.equals(name))) {
            addChild(new Note(mDataFile, reader, state));
        } else {
            super.loadSubElement(reader, state);
        }
    }

    @Override
    protected void saveSelf(XMLWriter out, boolean forUndo) {
        out.simpleTagNotEmpty(TAG_TEXT, mText);
    }

    /** @return The description. */
    public String getDescription() {
        return mText;
    }

    /**
     * @param description The description to set.
     * @return Whether it was modified.
     */
    public boolean setDescription(String description) {
        if (!mText.equals(description)) {
            mText = description;
            notifySingle(ID_TEXT);
            return true;
        }
        return false;
    }

    @Override
    public boolean contains(String text, boolean lowerCaseOnly) {
        if (getDescription().toLowerCase().contains(text)) {
            return true;
        }
        return super.contains(text, lowerCaseOnly);
    }

    @Override
    public Object getData(Column column) {
        return NoteColumn.values()[column.getID()].getData(this);
    }

    @Override
    public String getDataAsText(Column column) {
        return NoteColumn.values()[column.getID()].getDataAsText(this);
    }

    @Override
    public String toString() {
        return getDescription();
    }

    @Override
    public RetinaIcon getIcon(boolean marker) {
        return marker ? Images.NOT_MARKER : Images.NOT_FILE;
    }

    @Override
    public RowEditor<? extends ListRow> createEditor() {
        return new NoteEditor(this);
    }

    @Override
    public void fillWithNameableKeys(Set<String> set) {
        // No nameables
    }

    @Override
    public void applyNameableKeys(Map<String, String> map) {
        // No nameables
    }
}
