package com.little.painter.shape;

import android.graphics.Bitmap;

import java.util.ArrayList;
import java.util.List;

public class Gallery {
    /**
     * 画册名
     */
    private String mName;
    private int mNum;
    private List<List<Shape>> mPaintingList;
    /**
     * 缩略图集合
     */
    private List<Bitmap> mBitmapList;

    public Gallery() {
        mPaintingList = new ArrayList<>();
        mBitmapList = new ArrayList<>();
        mNum = 0;
    }

    public String getName() {
        return mName;
    }

    public int getNum() {
        return mNum;
    }

    public List<List<Shape>> getPaintingList() {
        return mPaintingList;
    }

    public void setPaintingList(List<List<Shape>> paintingList) {
        mPaintingList = paintingList;
    }

    public void setName(String name) {
        this.mName = name;
    }

    public void setNum(int num) {
        this.mNum = num;
    }

    public List<Bitmap> getBitmapList() {
        return mBitmapList;
    }

    public void setBitmapList(List<Bitmap> bitmapList) {
        mBitmapList = bitmapList;
    }

    public void AddPainting(List<Shape> painting, Bitmap bitmap) {
        List<Shape> shapes = new ArrayList<>();
        shapes.addAll(painting);
        mPaintingList.add(shapes);
        Bitmap bitmapobject = Bitmap.createBitmap(bitmap);
        mBitmapList.add(bitmapobject);
        mNum++;
    }

    public void CoverPainting(List<Shape> painting, Bitmap bitmap, int position) {
        List<Shape> shapes = new ArrayList<>();
        shapes.addAll(painting);
        mPaintingList.set(position, shapes);
        mBitmapList.set(position, bitmap);
    }
}
