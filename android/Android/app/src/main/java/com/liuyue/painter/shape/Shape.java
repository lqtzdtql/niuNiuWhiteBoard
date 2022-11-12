package com.liuyue.painter.shape;

import android.graphics.Canvas;
import android.graphics.Paint;
import android.graphics.RectF;

import java.util.ArrayList;
import java.util.List;

public abstract class Shape {
    /**
     * 颜色
     */
    protected String color;
    /**
     * 宽度
     */
    protected float width;
    /**
     * 笔迹上点集合
     */
    protected List<Point> pointList;
    /**
     * 画笔
     */
    protected final Paint paint;

    public Shape() {
        pointList = new ArrayList<>();
        paint = new Paint();
    }

    public String getColor() {
        return color;
    }

    public float getWidth() {
        return width;
    }

    public List<Point> getPointList() {
        return pointList;
    }

    public Paint getPaint() {
        return paint;
    }

    public void setColor(String color) {
        this.color = color;
    }

    public void setWidth(float width) {
        this.width = width;
    }

    public void setPointList(List<Point> pointList) {
        this.pointList = pointList;
    }

    public void setPaint(Paint paint) {
        this.paint.setColor(paint.getColor());
        this.paint.setStyle(Paint.Style.STROKE);
        this.paint.setStrokeWidth(paint.getStrokeWidth());
    }

    /**
     * 添加点函数
     */
    public void addPoint(float x, float y) {
        pointList.add(new Point(x, y));
    }

    /**
     * 绘制函数
     */
    public abstract void draw(Canvas mCanvas);

    /**
     * 按下操作对应的相关处理
     */
    public abstract void downAction(float x, float y);

    /**
     * 移动过程中相关操作
     */
    public abstract void moveAction(float mx, float my, float x, float y);

    /**
     * 抬起操作对应的相关处理
     */
    public abstract void upAction(float x, float y);

    /**
     * 返回自己对应种类
     */
    public abstract int getKind();

    /**
     * 设置自己特有属性
     */
    public abstract void setOwnProperty();

    /**
     * 找到笔迹的边缘矩形
     */
    public RectF findShapeEdge() {
        float minx = pointList.get(0).getX();
        float miny = pointList.get(0).getY();
        float maxx = pointList.get(0).getX();
        float maxy = pointList.get(0).getY();
        for (int i = 1; i < pointList.size(); i++) {
            if (maxx < pointList.get(i).getX()) {
                maxx = pointList.get(i).getX();
            }
            if (minx > pointList.get(i).getX()) {
                minx = pointList.get(i).getX();
            }
            if (maxy < pointList.get(i).getY()) {
                maxy = pointList.get(i).getY();
            }
            if (miny > pointList.get(i).getY()) {
                miny = pointList.get(i).getY();
            }
        }
        return new RectF(minx, maxy, maxx, miny);
    }

    /**
     * 判断是否相交并返回对应的list的位置
     */
    public abstract boolean isInterSect(float lastx, float lasty, float x, float y);

    /**
     * 判断是否进入边缘矩形
     */
    public boolean isEnterShapeEdge(float x, float y) {
        float minx = pointList.get(0).getX();
        float miny = pointList.get(0).getY();
        float maxx = pointList.get(0).getX();
        float maxy = pointList.get(0).getY();
        for (int i = 1; i < pointList.size(); i++) {
            if (maxx < pointList.get(i).getX()) {
                maxx = pointList.get(i).getX();
            }
            if (minx > pointList.get(i).getX()) {
                minx = pointList.get(i).getX();
            }
            if (maxy < pointList.get(i).getY()) {
                maxy = pointList.get(i).getY();
            }
            if (miny > pointList.get(i).getY()) {
                miny = pointList.get(i).getY();
            }
        }
        return (x >= minx && x <= maxx) && (y >= miny && y <= maxy);
    }

}
