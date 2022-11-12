package com.liuyue.painter.shape;

import android.graphics.Canvas;

import com.liuyue.painter.Constants;
import com.liuyue.painter.utils.InterSectUtil;

public class Line extends Shape {
    private Point mStartPoint;
    private Point mEndPoint;

    public Point getStartPoint() {
        return mStartPoint;
    }

    public Point getEndPoint() {
        return mEndPoint;
    }

    public void setStartPoint(Point startPoint) {
        this.mStartPoint = startPoint;
    }

    public void setEndPoint(Point endPoint) {
        this.mEndPoint = endPoint;
    }

    @Override
    public void draw(Canvas mCanvas) {
        if (mStartPoint != null && mEndPoint != null) {
            mCanvas.drawLine(
                    mStartPoint.getX(),
                    mStartPoint.getY(),
                    mEndPoint.getX(),
                    mEndPoint.getY(),
                    paint);
        }
    }

    @Override
    public void downAction(float x, float y) {
        // 设置初始点和终结点
        setStartPoint(new Point(x, y));
        setEndPoint(new Point(x, y));
    }

    @Override
    public void moveAction(float mx, float my, float x, float y) {
        // 修改终结点
        setEndPoint(new Point(x, y));
    }

    @Override
    public void upAction(float x, float y) {
        // 设置终结点
        setEndPoint(new Point(x, y));
    }

    @Override
    public int getKind() {
        return Constants.LINE;
    }

    @Override
    public void setOwnProperty() {
        // 获取关键点
        setStartPoint(pointList.get(0));
        setEndPoint(pointList.get(1));
    }

    @Override
    public boolean isInterSect(float lastx, float lasty, float x, float y) {
        // 直线从逻辑上来讲和曲线是一样的
        for (int i = 1; i < pointList.size(); i++) {
            if (new InterSectUtil(new Point(lastx, lasty), new Point(x, y), pointList.get(i - 1), pointList.get(i)).Segment_Intersect()) {
                return true;
            }
        }
        return false;
    }
}
