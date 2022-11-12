package com.liuyue.painter.shape;

public class Point {
    private float mX;
    private float mY;

    public Point() {
    }

    public Point(float x, float y) {
        this.mX = x;
        this.mY = y;
    }

    public float getX() {
        return mX;
    }

    public float getY() {
        return mY;
    }

    public void setX(float x) {
        this.mX = x;
    }

    public void setY(float y) {
        this.mY = y;
    }
}
