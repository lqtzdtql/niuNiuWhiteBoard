package com.little.painter.shape;

import android.graphics.Canvas;
import android.graphics.Path;

import com.little.painter.Constants;
import com.little.painter.utils.InterSectUtil;

import java.util.List;

public class Ink extends Shape {
    private Path path;

    public Ink() {
        path = new Path();
    }

    public Path getPath() {
        return path;
    }

    public void setPath(Path path) {
        this.path = path;
    }

    public void createPath(List<Point> points) {
        Path newpath = new Path();
        // 起点
        newpath.moveTo(points.get(0).getX(), points.get(0).getY());
        for (int j = 1; j < points.size() - 1; j++) {
            float mx = points.get(j - 1).getX();
            float my = points.get(j - 1).getY();
            float x = points.get(j).getX();
            float y = points.get(j).getY();
            newpath.quadTo(mx, my, (x + mx) / 2, (y + my) / 2);
        }
        // 终点
        newpath.lineTo(points.get(points.size() - 1).getX(), points.get(points.size() - 1).getY());
        this.setPath(newpath);
    }

    @Override
    public void draw(Canvas mCanvas) {
        // 不为空则绘制
        if (this.path != null) {
            mCanvas.drawPath(path, paint);
        }
    }

    @Override
    public void downAction(float x, float y) {
        path.moveTo(x, y);
    }

    @Override
    public void moveAction(float mx, float my, float x, float y) {
        path.quadTo(mx, my, (x + mx) / 2, (y + my) / 2);
        // 保存点
        addPoint(x, y);
    }

    @Override
    public void upAction(float x, float y) {
        path.lineTo(x, y);
    }

    @Override
    public int getKind() {
        return Constants.INK;
    }

    @Override
    public void setOwnProperty() {
        Path newPath = new Path();
        newPath.moveTo(pointList.get(0).getX(), pointList.get(0).getY());
        int j;
        for (j = 1; j < pointList.size() - 1; j++) {
            newPath.quadTo(
                    pointList.get(j - 1).getX(),
                    pointList.get(j - 1).getY(),
                    (pointList.get(j).getX() + pointList.get(j - 1).getX()) / 2,
                    (pointList.get(j).getY() + pointList.get(j - 1).getY()) / 2);
        }
        newPath.lineTo(pointList.get(j).getX(), pointList.get(j).getY());
        setPath(newPath);
    }

    @Override
    public boolean isInterSect(float lastx, float lasty, float x, float y) {
        for (int i = 1; i < pointList.size(); i++) {
            if (new InterSectUtil(new Point(lastx, lasty), new Point(x, y), pointList.get(i - 1), pointList.get(i)).Segment_Intersect()) {
                return true;
            }
        }
        return false;
    }
}
