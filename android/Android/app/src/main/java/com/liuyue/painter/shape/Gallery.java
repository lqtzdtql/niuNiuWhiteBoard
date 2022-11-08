package com.liuyue.painter.shape;

import android.graphics.Bitmap;

import java.util.ArrayList;
import java.util.List;

public class Gallery {
    private String name;// 画册名
    private int num;
    private List<List<Shape>> PaintingList;

    private List<Bitmap> BitmapList; // 缩略图集合

    public Gallery() {
        PaintingList = new ArrayList<>();
        BitmapList = new ArrayList<>();
        num = 0;
    }


    public String getName() {
        return name;
    }

    public int getNum() {
        return num;
    }

    public List<List<Shape>> getPaintingList() {
        return PaintingList;
    }

    public void setPaintingList(List<List<Shape>> paintingList) {
        PaintingList = paintingList;
    }

    public void setName(String name) {
        this.name = name;
    }

    public void setNum(int num) {
        this.num = num;
    }

    public List<Bitmap> getBitmapList() {
        return BitmapList;
    }

    public void setBitmapList(List<Bitmap> bitmapList) {
        BitmapList = bitmapList;
    }

    public void AddPainting(List<Shape> painting, Bitmap bitmap) {
        List<Shape> shapes = new ArrayList<>();
        shapes.addAll(painting);
        PaintingList.add(shapes);
        Bitmap bitmapobject = Bitmap.createBitmap(bitmap);
        BitmapList.add(bitmapobject);
        num++;
    }

    public void CoverPainting(List<Shape> painting, Bitmap bitmap, int position) {
        List<Shape> shapes = new ArrayList<>();
        shapes.addAll(painting);
        PaintingList.set(position, shapes);
        BitmapList.set(position, bitmap);
    }
}
