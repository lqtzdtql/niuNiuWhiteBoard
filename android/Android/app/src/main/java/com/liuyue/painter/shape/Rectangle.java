package com.liuyue.painter.shape;

import android.graphics.Canvas;

import com.liuyue.painter.Constants;
import com.liuyue.painter.utils.InterSectUtil;

import java.util.ArrayList;
import java.util.List;

public class Rectangle extends Shape {
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
            mCanvas.drawRect(mStartPoint.getX(), mStartPoint.getY(), mEndPoint.getX(), mEndPoint.getY(), paint);
        }
    }

    @Override
    public void downAction(float x, float y) {
        // 设置初始点和终止点
        setStartPoint(new Point(x, y));
        setEndPoint(new Point(x, y));
    }

    @Override
    public void moveAction(float mx, float my, float x, float y) {
        // 修改终止点
        setEndPoint(new Point(x, y));
    }

    @Override
    public void upAction(float x, float y) {
        // 设置终止点
        setEndPoint(new Point(x, y));
    }

    @Override
    public int getKind() {
        return Constants.RECT;
    }

    @Override
    public void setOwnProperty() {
        // 获取关键点
        setStartPoint(pointList.get(0));
        setEndPoint(pointList.get(1));
    }

    @Override
    public boolean isInterSect(float lastx, float lasty, float x, float y) {
        // 矩形逻辑和曲线也是一样的，一共四个点
        List<Point> judgePointList = new ArrayList<>();
        judgePointList.add(pointList.get(0));
        judgePointList.add(new Point(pointList.get(1).getX(), pointList.get(0).getY()));
        judgePointList.add(new Point(pointList.get(1).getX(), pointList.get(1).getY()));
        judgePointList.add(new Point(pointList.get(0).getX(), pointList.get(1).getY()));
        judgePointList.add(pointList.get(0));
        for (int i = 1; i < judgePointList.size(); i++) {
            if (new InterSectUtil(new Point(lastx, lasty), new Point(x, y), judgePointList.get(i - 1), judgePointList.get(i)).Segment_Intersect()) {
                return true;
            }
        }
        return false;
    }
}
