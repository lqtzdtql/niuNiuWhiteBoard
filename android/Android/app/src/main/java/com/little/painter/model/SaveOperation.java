package com.little.painter.model;

import com.little.painter.view.ArtBoard;

public abstract class SaveOperation {
    String filepath;
    String filename = null;

    public String getFilepath() {
        return filepath;
    }

    public void setFilepath(String filepath) {
        this.filepath = filepath;
    }

    public String getFilename() {
        return filename;
    }

    public void setFilename(String filename) {
        this.filename = filename;
    }

    public abstract String GetAbusoluteFileName();

    public abstract void SavePainting();

    public abstract void GetContent(ArtBoard artBoard);
}
