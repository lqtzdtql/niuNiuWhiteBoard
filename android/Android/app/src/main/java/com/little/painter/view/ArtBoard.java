package com.little.painter.view;

import android.content.Context;
import android.graphics.Bitmap;
import android.graphics.Canvas;
import android.graphics.Color;
import android.graphics.DashPathEffect;
import android.graphics.Paint;
import android.graphics.Path;
import android.graphics.PathEffect;
import android.graphics.RectF;
import android.os.Handler;
import android.os.Looper;
import android.os.Message;
import android.util.AttributeSet;
import android.view.MotionEvent;
import android.view.View;

import com.little.painter.Constants;
import com.little.painter.shape.Circle;
import com.little.painter.shape.Gallery;
import com.little.painter.shape.Ink;
import com.little.painter.shape.Line;
import com.little.painter.shape.Point;
import com.little.painter.shape.Rectangle;
import com.little.painter.shape.Shape;
import com.little.painter.utils.XmlOperation;

import org.dom4j.DocumentException;

import java.io.File;
import java.util.ArrayList;
import java.util.Iterator;
import java.util.List;

public class ArtBoard extends View {
    /**
     * 最初 list 的长度
     */
    private final int mPreShapeListSize = 0;
    /**
     * 画布的宽
     */
    private int mCanvasWidth;
    /**
     * 画布的高
     */
    private int mCanvasHeight;
    /**
     * 当前画笔位置
     */
    private float mX;
    private float mY;
    /**
     * 已保存笔迹List
     */
    private List<Shape> mSaveShapeList;
    /**
     * 删除笔迹List
     */
    private final List<Shape> mDeleteShapeList;
    /**
     * 画册类
     */
    private Gallery mGallery;
    private Shape mCurrentShape;
    private Path mEraserPath;
    private Path mMovePath;
    private List<Integer> mSweepList;
    private List<Shape> mNeedHandleList;
    private List<Integer> mMoveList;
    private List<Shape> mNeedMoveList;
    private RectF mNeedMoveRect;
    /**
     * 默认字体大小
     */
    private float mCurrentWidth = 5;
    /**
     * 默认字体颜色
     */
    private String mCurrentColor = Constants.colors[0];
    /**
     * 默认绘图类型
     */
    private int mCurrentKind = Constants.INK;
    /**
     * 默认当前页数
     */
    private int mCurrentPageNum = 1;
    /**
     * 默认当前页序号
     */
    private int mCurrentPageIndex = 1;
    /**
     * 画笔
     */
    private Paint mPaint;
    /**
     * 画布的 bitmap
     */
    private Bitmap mBitmap;
    /**
     * 画布
     */
    private Canvas mCanvas;
    /**
     * 是否有编辑过标志
     */
    private boolean mIsFileExist = false;
    /**
     * 是否处于橡皮擦状态
     */
    private boolean mEraserState = false;
    /**
     * 是否处于笔迹操作状态
     */
    private boolean mMoveShapeState = false;
    private final Handler mHandler = new Handler(Looper.getMainLooper()) {
        @Override
        public void handleMessage(Message msg) {
            super.handleMessage(msg);
            if (msg.what == 0) {
                try {
                    redrawOnBitmap();
                } catch (Exception e) {
                    e.printStackTrace();
                }
            }
        }
    };
    /**
     * 已选中状态
     */
    private boolean mSelected = false;
    /**
     * 移动笔迹的标志
     */
    private boolean mIsMoving = false;
    private ShapeChangeListener mShapeChangeListener;

    public ArtBoard(Context context) {
        this(context, null);
    }

    public ArtBoard(Context context, AttributeSet attrs) {
        this(context, attrs, 0);
    }

    public ArtBoard(Context context, AttributeSet attrs, int defStyleAttr) {
        super(context, attrs, defStyleAttr);
        mSaveShapeList = new ArrayList<>();
        mDeleteShapeList = new ArrayList<>();
        mGallery = new Gallery();

    }

    @Override
    protected void onMeasure(int widthMeasureSpec, int heightMeasureSpec) {
        super.onMeasure(widthMeasureSpec, heightMeasureSpec);
        mCanvasWidth = getMeasuredWidth();
        mCanvasHeight = getMeasuredHeight();
        initCanvas();
    }

    /**
     * 初始化画布
     */
    private void initCanvas() {
        // 初始化画布
        mBitmap = Bitmap.createBitmap(mCanvasWidth, mCanvasHeight, Bitmap.Config.ARGB_8888);
        // 所有 mCanvas 画的东西都被保存在了 mBitmap中
        mCanvas = new Canvas(mBitmap);
        mCanvas.drawColor(Color.WHITE);
        // 初始化画笔
        mPaint = new Paint();
        mPaint.setStyle(Paint.Style.STROKE);
        mPaint.setStrokeWidth(mCurrentWidth);
        mPaint.setColor(Color.parseColor(mCurrentColor));
    }

    @Override
    public boolean onTouchEvent(MotionEvent event) {
        float x = event.getX();
        float y = event.getY();
        switch (event.getAction()) {
            case MotionEvent.ACTION_DOWN:
                touchStart(x, y);
                invalidate();
                break;
            case MotionEvent.ACTION_MOVE:
                touchMove(x, y);
                break;
            case MotionEvent.ACTION_UP:
                touchUp(x, y);
                break;
            default:
                break;
        }
        invalidate();
        return true;
    }

    /**
     * 按下操作对应处理
     */
    private void touchStart(float x, float y) {
        // 进入橡皮擦模式
        if (mEraserState) {
            // 准备画虚线需要的相关属性
            PathEffect effects = new DashPathEffect(new float[]{8, 8, 8, 8}, 1);
            mPaint.setPathEffect(effects);
            mPaint.setColor(Color.parseColor(Constants.colors[1]));
            mSweepList = new ArrayList<>();
            mNeedHandleList = new ArrayList<>();
            // 创建橡皮擦的 path
            mEraserPath = new Path();
            mEraserPath.moveTo(x, y);
        }
        // 进入选择线条模式
        else if (mMoveShapeState) {
            if (mSelected) {
                // 判断点击的地方是否是在 NeedRect 内部如果
                if (!IsNotInside(x, y)) {
                    // 不在范围内，相关参数清零从头开始
                    mNeedMoveRect = null;
                    mMoveList = null;
                    mNeedMoveList = null;
                } else {
                    // 在范围内采取相关措施
                    mIsMoving = true;
                    // TODO 移动和缩放相关
                    // 脏区为 NeedMoveRect
                }
            } else {
                // 准备画虚线需要的相关属性
                PathEffect effects = new DashPathEffect(new float[]{8, 8, 8, 8}, 1);
                mPaint.setPathEffect(effects);
                mPaint.setColor(Color.parseColor(Constants.colors[1]));

                mMoveList = new ArrayList<>();
                mNeedMoveList = new ArrayList<>();
                // 创建选中笔迹的 path
                mMovePath = new Path();
                mMovePath.moveTo(x, y);
            }
        } else {
            // 判断当前类型，根据类型选择构造函数
            switch (mCurrentKind) {
                case Constants.INK:
                    mCurrentShape = new Ink();
                    break;
                case Constants.LINE:
                    mCurrentShape = new Line();
                    break;
                case Constants.RECT:
                    mCurrentShape = new Rectangle();
                    break;
                case Constants.CIRCLE:
                    mCurrentShape = new Circle();
                    break;
            }
            // 执行对应操作
            mCurrentShape.downAction(x, y);
            // 设置画笔
            mCurrentShape.setPaint(mPaint);
            // 记录起始点
            mCurrentShape.addPoint(x, y);
            // 获得颜色和宽度数据
            mCurrentShape.setColor(mCurrentColor);
            mCurrentShape.setWidth(mCurrentWidth);
        }
        mX = x;
        mY = y;
    }

    /**
     * 移动操作对应处理
     */
    private void touchMove(float x, float y) {
        if (mEraserState) {
            mEraserPath.quadTo(mX, mY, x, y);
            // 遍历笔迹
            for (int i = 0; i < mSaveShapeList.size(); i++) {
                // 判断进入对应的矩形
                if (mSaveShapeList.get(i).isEnterShapeEdge(x, y)) {
                    // 判断是否发生相交
                    if (mSaveShapeList.get(i).isInterSect(mX, mY, x, y)) {
                        // 记录当前 shape 的 position
                        mSweepList.add(i);
                    }
                }
            }
        } else if (mMoveShapeState) {
            if (mSelected) {
                if (mIsMoving) {
                    for (int k = 0; k < mNeedMoveList.size(); k++) {
                        // 取出
                        Shape shape = mNeedMoveList.get(k);
                        for (int j = 0; j < shape.getPointList().size(); j++) {
                            // 修改坐标
                            float moveX = shape.getPointList().get(j).getX() + (x - mX);
                            float moveY = shape.getPointList().get(j).getY() + (y - mY);
                            // 保存坐标
                            shape.getPointList().set(j, new Point(moveX, moveY));
                        }
                        if (mShapeChangeListener != null) {
                            mShapeChangeListener.onMoveShape(shape);
                        }
                        // 替换 SaveList 中的对应 shape
                        mSaveShapeList.set(mMoveList.get(k), shape);
                    }
                }
            } else {
                mMovePath.quadTo(mX, mY, x, y);
                // 遍历笔迹
                for (int i = 0; i < mSaveShapeList.size(); i++) {
                    // 判断进入对应的矩形
                    if (mSaveShapeList.get(i).isEnterShapeEdge(x, y)) {
                        // 判断是否发生相交
                        if (mSaveShapeList.get(i).isInterSect(mX, mY, x, y)) {
                            // 记录当前 shape 的 position
                            mMoveList.add(i);
                        }
                    }
                }
            }
        } else {
            // 执行相关操作
            mCurrentShape.moveAction(mX, mY, x, y);
        }
        // 记录当前坐标点
        mX = x;
        mY = y;
    }

    /**
     * 点击区域不在对应范围内
     */
    private boolean IsNotInside(float x, float y) {
        if (mNeedMoveRect != null) {
            return x >= mNeedMoveRect.left && x <= mNeedMoveRect.right && y >= mNeedMoveRect.bottom && y <= mNeedMoveRect.top;
        }
        return false;
    }

    /**
     * 抬起操作对应处理
     */
    private void touchUp(float x, float y) {
        if (mEraserState) {
            mEraserPath.lineTo(x, y);
            mEraserPath = null;
            if (!mSweepList.isEmpty()) {
                // 删除选中的 shape
                for (int i = 0; i < mSweepList.size(); i++) {
                    // 根据下标取出对象
                    mNeedHandleList.add(mSaveShapeList.get(mSweepList.get(i)));
                }
                // 遍历对象依次删除
                for (int j = 0; j < mNeedHandleList.size(); j++) {
                    Shape deleteObject = mNeedHandleList.get(j);
                    Iterator<Shape> it = mSaveShapeList.iterator();
                    while (it.hasNext()) {
                        Shape shape = it.next();
                        if (shape == deleteObject) {
                            // 删除的笔迹放入DeleteList
                            mDeleteShapeList.add(shape);
                            if (mShapeChangeListener != null) {
                                mShapeChangeListener.onDeleteShape(shape);
                            }
                            it.remove();
                        }
                    }
                }
            }
            // 相关参数清空
            mSweepList = null;
            mNeedHandleList = null;
            // 通知系统重绘
            Message msg = new Message();
            msg.what = 0;
            mHandler.sendMessageDelayed(msg, 100);
        } else if (mMoveShapeState) {
            if (mSelected) {
                if (mNeedMoveList == null) {
                    mSelected = false;
                }
                try {
                    redrawOnBitmap();
                } catch (Exception e) {
                    e.printStackTrace();
                }
            } else {
                mMovePath.lineTo(x, y);
                if (!mMoveList.isEmpty()) {
                    // 删除选中的shape
                    for (int i = 0; i < mMoveList.size(); i++) {
                        // 根据下标取出对象
                        mNeedMoveList.add(mSaveShapeList.get(mMoveList.get(i)));
                    }
                    // 遍历找到笔迹最大的Rect区域
                    mNeedMoveRect = findBiggestRect(mNeedMoveList);
                }
                // 相关参数清空
                mMovePath = null;
                // 设置为已选中状态
                mSelected = true;
                // 现在 NeedMoveList 中有保存对应笔迹 NeedMoveRect 不为空 MoveList 也保存有笔迹对应下标
            }
            // 通知系统重绘
            Message msg = new Message();
            msg.what = 0;
            mHandler.sendMessageDelayed(msg, 100);
        } else {
            // 执行相关操作
            mCurrentShape.upAction(x, y);
            // 绘制到 Bitmap 上去
            mCurrentShape.draw(mCanvas);
            // 保存终结点
            mCurrentShape.addPoint(x, y);
            // 将笔迹添加到栈中
            mSaveShapeList.add(mCurrentShape);
            if (mShapeChangeListener != null) {
                mShapeChangeListener.onAddShape(mCurrentShape);
            }
            // 对象置空
            mCurrentShape = null;
        }
    }

    /**
     * 找到要移动的区域
     */
    private RectF findBiggestRect(List<Shape> needMoveList) {
        float minx = needMoveList.get(0).getPointList().get(0).getX();
        float miny = needMoveList.get(0).getPointList().get(0).getY();
        float maxx = needMoveList.get(0).getPointList().get(0).getX();
        float maxy = needMoveList.get(0).getPointList().get(0).getY();
        for (int k = 0; k < needMoveList.size(); k++) {
            List<Point> pointList = needMoveList.get(k).getPointList();
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
        }
        return new RectF(minx, maxy, maxx, miny);
    }

    @Override
    protected void onDraw(Canvas canvas) {
        super.onDraw(canvas);
        canvas.drawBitmap(mBitmap, 0, 0, null);
        if (mCurrentShape != null) {
            mCurrentShape.draw(canvas);
        }
        if (mEraserState && mEraserPath != null) {
            canvas.drawPath(mEraserPath, mPaint);
        }
        if (mMoveShapeState && mMovePath != null) {
            canvas.drawPath(mMovePath, mPaint);
        }
    }

    /**
     * 撤销操作
     */
    public void undo() {
        if (mSaveShapeList != null && mSaveShapeList.size() >= 1) {
            Shape shape = mSaveShapeList.get(mSaveShapeList.size() - 1);
            mDeleteShapeList.add(shape);
            mSaveShapeList.remove(shape);
            if (mShapeChangeListener != null) {
                mShapeChangeListener.onDeleteShape(shape);
            }
            try {
                // 重新绘制图案
                redrawOnBitmap();
            } catch (Exception e) {
                e.printStackTrace();
            }
        }

    }

    /**
     * 重做操作
     */
    public void redo() {
        if (mDeleteShapeList != null && mDeleteShapeList.size() >= 1) {
            Shape shape = mDeleteShapeList.get(mDeleteShapeList.size() - 1);
            mSaveShapeList.add(shape);
            mDeleteShapeList.remove(shape);
            if (mShapeChangeListener != null) {
                mShapeChangeListener.onDeleteShape(shape);
            }
            try {
                // 重新绘制图案
                redrawOnBitmap();
            } catch (Exception e) {
                e.printStackTrace();
            }
        }
    }

    /**
     * 加载之前的画
     */
    public void loadFile(String filename) {
        mIsFileExist = true;
        // 遍历文件夹下每个 xml 文件
        File scanFilePath = new File(filename);
        if (scanFilePath.isDirectory()) {
            for (File file : scanFilePath.listFiles()) {
                String fileAbsolutePath = file.getAbsolutePath();
                if (fileAbsolutePath.endsWith(".xml")) {
                    // 将 xml 解析放入 mGallery 中
                    try {
                        mGallery.AddPainting(XmlOperation.TransXmlToShape(fileAbsolutePath), mBitmap);
                    } catch (DocumentException e) {
                        e.printStackTrace();
                    }
                }
            }
            // 添加名字属性
            String name = filename.substring(filename.lastIndexOf("/") + 1);
            mGallery.setName(name);
            // 修改对应页数相关
            mCurrentPageNum = mGallery.getNum();
            // 载入当前第一页内容
            mSaveShapeList.clear();
            mSaveShapeList.addAll(mGallery.getPaintingList().get(0));
            try {
                redrawOnBitmap();
            } catch (Exception e) {
                e.printStackTrace();
            }
        }
    }

    public int getCurrentPageNum() {
        return mCurrentPageNum;
    }

    public void setCurrentPageNum(int currentPageNum) {
        mCurrentPageNum = currentPageNum;
    }

    public int getCurrentPageIndex() {
        return mCurrentPageIndex;
    }

    public void setCurrentPageIndex(int currentPageIndex) {
        mCurrentPageIndex = currentPageIndex;
    }

    /**
     * 获取笔迹宽度
     */
    public float getBrushSize() {
        return mPaint.getStrokeWidth();
    }

    /**
     * 设置笔迹宽度
     */
    public void setBrushSize(float brushsize) {
        mCurrentWidth = brushsize;
        mPaint.setStrokeWidth(brushsize);
    }

    /**
     * 获取笔迹颜色
     */
    public String getBrushColor() {
        return mCurrentColor;
    }

    /**
     * 设置笔迹颜色
     */
    public void setBrushColor(String color) {
        mCurrentColor = color;
        mPaint.setColor(Color.parseColor(mCurrentColor));
    }

    /**
     * 获取笔迹类型
     */
    public int getCurrentKind() {
        return mCurrentKind;
    }

    /**
     * 设置笔迹类型
     */
    public void setCurrentKind(int currentKind) {
        this.mCurrentKind = currentKind;
    }

    /**
     * 获取SaveList
     */
    public List<Shape> getSaveShapeList() {
        return mSaveShapeList;
    }

    /**
     * 设置 SaveList
     */
    public void setSaveShapeList(List<Shape> saveShapeList) {
        mSaveShapeList = saveShapeList;
    }

    /**
     * 获取Bitmap
     */
    public Bitmap getmBitmap() {
        return mBitmap;
    }

    /**
     * 设置 Bitmap
     */
    public void setBitmap(Bitmap bitmap) {
        this.mBitmap = bitmap;
    }

    /**
     * 设置画布宽度
     */
    public void setCanvasWidth(int canvasWidth) {
        mCanvasWidth = canvasWidth;
    }

    /**
     * 设置画布高度
     */
    public void setCanvasHeight(int canvasHeight) {
        mCanvasHeight = canvasHeight;
    }

    public Gallery getGallery() {
        return mGallery;
    }

    public void setGallery(Gallery mGallery) {
        this.mGallery = mGallery;
    }

    /**
     * 重新绘制 Bitmap 上的图案
     */
    public void redrawOnBitmap() {
        // 重新设置画布，相当于清空画布
        initCanvas();
        // 依次遍历，绘制对应图案
        for (int i = 0; i < mSaveShapeList.size(); i++) {
            mSaveShapeList.get(i).draw(mCanvas);
        }
        if (mMoveShapeState && mNeedMoveList != null) {
            // 设置虚线的间隔和点的长度
            PathEffect effects = new DashPathEffect(new float[]{8, 8, 8, 8}, 1);
            Paint newPaint = new Paint();
            newPaint.setPathEffect(effects);
            newPaint.setColor(Color.parseColor(Constants.colors[1]));
            newPaint.setStyle(Paint.Style.STROKE);
            mCanvas.drawRect(mNeedMoveRect.left, mNeedMoveRect.top, mNeedMoveRect.right, mNeedMoveRect.bottom, newPaint);
        }
        invalidate();
    }

    /**
     * 判断是否有内容
     */
    public boolean isEmpty() {
        return mSaveShapeList.size() == 0;
    }

    /**
     * 判断是否是在已存在文件上编辑
     */
    public boolean isFileExist() {
        return mIsFileExist;
    }

    /**
     * 判断是否发生了编辑
     */
    public boolean isEdited() {
        return mPreShapeListSize != mSaveShapeList.size();
    }

    /**
     * 绘制新的图形
     */
    public void drawNewImage() {
        if (mGallery.getNum() != mCurrentPageNum) {
            // 保存当前笔迹集合及 Bitmap
            mGallery.AddPainting(mSaveShapeList, mBitmap);
        }
        mSaveShapeList.clear();
        mDeleteShapeList.clear();
        // 清空画布及相关数据
        initCanvas();
        invalidate();
    }

    /**
     * 返回到上一页
     */
    public void turnToPrePage() {
        // 刚好处于最后一页要往前翻
        if (mGallery.getNum() == mCurrentPageNum - 1) {
            // 保存当前图形
            mGallery.AddPainting(mSaveShapeList, mBitmap);
        } else {
            // 覆盖当前图形
            mGallery.CoverPainting(mSaveShapeList, mBitmap, mCurrentPageIndex);
        }
        // 清空画布及相关数据
        initCanvas();
        mSaveShapeList.clear();
        mDeleteShapeList.clear();
        // 加载上一页内容
        mSaveShapeList.addAll(mGallery.getPaintingList().get(mCurrentPageIndex - 1));
        try {
            redrawOnBitmap();
        } catch (Exception e) {
            e.printStackTrace();
        }
    }


    /**
     * 跳转下一页
     */
    public void turnToNextPage() {
        // 覆盖当前图形
        mGallery.CoverPainting(mSaveShapeList, mBitmap, mCurrentPageIndex - 2);
        // 清空画布及相关数据
        initCanvas();
        mSaveShapeList.clear();
        mDeleteShapeList.clear();
        // 加载下一页内容
        mSaveShapeList.addAll(mGallery.getPaintingList().get(mCurrentPageIndex - 1));
        try {
            redrawOnBitmap();
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    /**
     * 判断是否需要保存最后一张并处理
     */
    public void saveTheLast() {
        if (mGallery.getNum() == mCurrentPageNum - 1) {
            mGallery.AddPainting(mSaveShapeList, mBitmap);
        } else {
            mGallery.CoverPainting(mSaveShapeList, mBitmap, mCurrentPageIndex - 1);
        }
    }

    /**
     * 修改橡皮擦状态
     */
    public void changeEraserState() {
        if (mMoveShapeState) {
            mMoveShapeState = false;
        }
        mEraserState = !mEraserState;
    }

    /**
     * 修改笔迹相关操作状态
     */
    public void changeCutState() {
        if (mEraserState) {
            mEraserState = false;
        }
        mMoveShapeState = !mMoveShapeState;
    }

    public void setShapeChangeListener(ShapeChangeListener listener) {
        mShapeChangeListener = listener;
    }

    public interface ShapeChangeListener {
        void onAddShape(Shape shape);

        void onMoveShape(Shape shape);

        void onDeleteShape(Shape shape);
    }
}
