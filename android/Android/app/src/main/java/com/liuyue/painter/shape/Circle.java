package com.liuyue.painter.shape;

import android.graphics.Canvas;

import com.liuyue.painter.Constants;

public class Circle extends Shape {
    Point startPoint;
    Point endPoint;

    public static float getDistance(float x1, float y1, float x2, float y2) {
        return (float) Math.sqrt((x2 - x1) * (x2 - x1) + (y2 - y1) * (y2 - y1));
    }

    public Point getStartPoint() {
        return startPoint;
    }

    public void setStartPoint(Point startPoint) {
        this.startPoint = startPoint;
    }

    public Point getEndPoint() {
        return endPoint;
    }

    public void setEndPoint(Point endPoint) {
        this.endPoint = endPoint;
    }

    @Override
    public void draw(Canvas mCanvas) {
        if (startPoint != null && endPoint != null) {
            float radius = Math.abs(endPoint.getY() - startPoint.getY()) >= Math.abs(endPoint.getX() - startPoint.getX()) ? Math.abs(endPoint.getX() - startPoint.getX()) / 2 : Math.abs(endPoint.getY() - startPoint.getY()) / 2;
            mCanvas.drawCircle(
                    (endPoint.getX() + startPoint.getX()) / 2,
                    (endPoint.getY() + startPoint.getY()) / 2,
                    radius,
                    paint);
        }
    }

    @Override
    public void DownAction(float x, float y) {
        // 设置初始点和终止点
        setStartPoint(new Point(x, y));
        setEndPoint(new Point(x, y));
    }

    @Override
    public void MoveAction(float mx, float my, float x, float y) {
        // 修改终止点
        setEndPoint(new Point(x, y));
    }

    @Override
    public void UpAction(float x, float y) {
        // 设置终止点
        setEndPoint(new Point(x, y));
    }

    @Override
    public int GetKind() {
        return Constants.CIRCLE;
    }

    @Override
    public void setOwnProperty() {
        // 获取关键点
        setStartPoint(pointList.get(0));
        setEndPoint(pointList.get(1));
    }

    @Override
    public boolean IsInterSect(float lastx, float lasty, float x, float y) {
        Point center = new Point((endPoint.getX() + startPoint.getX()) / 2, (endPoint.getY() + startPoint.getY()) / 2);
        float radius = Math.abs(
                endPoint.getY() - startPoint.getY()) >= Math.abs(endPoint.getX() - startPoint.getX())
                ? Math.abs(endPoint.getX() - startPoint.getX()) / 2
                : Math.abs(endPoint.getY() - startPoint.getY()) / 2;
        // TODO 目前采取这种方式
        return getDistance(x, y, center.getX(), center.getY()) < radius;
    }
}
