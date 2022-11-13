package com.little.painter.manager;

import android.view.View;
import android.widget.ImageButton;

import com.little.painter.R;

public class FootUIManager implements View.OnClickListener {
    private View mView;
    private ImageButton mBrushSizeBtn;
    private ImageButton mBrushColorBtn;
    private ImageButton mAddPageBtn;
    private ImageButton mCutPathBtn;
    private ImageButton mEraserBtn;
    private ImageButton mShapeBtn;

    private boolean mIsShow = false;
    private SizeBtnOnclickListener mSizeBtnOnclickListener;
    private ColorBtnOnclickListener mColorBtnOnclickListener;
    private ShapeBtnClickListener mShapeBtnClickListener;
    private PageBtnOnClickListener mPageBtnOnClickListener;
    private EraserBtnOnClickListener mEraserBtnOnClickListener;
    private ShapeChooseBtnOnclickListener mShapeChooseBtnOnclickListener;

    public FootUIManager(View mView) {
        this.mView = mView;
        initView();
        initEvent();
    }

    private void initEvent() {
        mBrushSizeBtn.setOnClickListener(this);
        mBrushColorBtn.setOnClickListener(this);
        mAddPageBtn.setOnClickListener(this);
        mCutPathBtn.setOnClickListener(this);
        mEraserBtn.setOnClickListener(this);
        mShapeBtn.setOnClickListener(this);
    }

    private void initView() {
        mBrushSizeBtn = (ImageButton) mView.findViewById(R.id.bt_brush_size_choose);
        mBrushColorBtn = (ImageButton) mView.findViewById(R.id.bt_color_choose);
        mAddPageBtn = (ImageButton) mView.findViewById(R.id.bt_add_canvas);
        mCutPathBtn = (ImageButton) mView.findViewById(R.id.bt_cut_path);
        mEraserBtn = (ImageButton) mView.findViewById(R.id.bt_eraser);
        mShapeBtn = (ImageButton) mView.findViewById(R.id.bt_select_shape);
    }

    @Override
    public void onClick(View v) {
        switch (v.getId()) {
            case R.id.bt_brush_size_choose:
                // 选择宽度功能
                if (mSizeBtnOnclickListener != null) {
                    mSizeBtnOnclickListener.Clicked(mIsShow);
                    // 状态取反
                    mIsShow = !mIsShow;
                }
                break;

            case R.id.bt_color_choose:
                // 选择颜色功能
                if (mColorBtnOnclickListener != null) {
                    mColorBtnOnclickListener.ColorClicked(mIsShow);
                    mIsShow = !mIsShow;
                }
                break;
            case R.id.bt_add_canvas:
                // 添加多页功能
                if (mPageBtnOnClickListener != null) {
                    mPageBtnOnClickListener.PageClicked(mIsShow);
                    mIsShow = !mIsShow;
                }
                break;
            case R.id.bt_cut_path:
                // TODO 笔迹获取及相关操作
                if (mShapeChooseBtnOnclickListener != null) {
                    mShapeChooseBtnOnclickListener.ChooseOnClicked();
                }
                break;
            case R.id.bt_eraser:
                // 橡皮擦功能
                if (mEraserBtnOnClickListener != null) {
                    mEraserBtnOnClickListener.EraserClicked();
                }
                break;
            case R.id.bt_select_shape:
                // 选择图形功能
                if (mShapeBtnClickListener != null) {
                    mShapeBtnClickListener.onShapeBtnClicked(mIsShow);
                    mIsShow = !mIsShow;
                }
                break;
            default:
                break;
        }
    }

    public void setSizeBtnListener(SizeBtnOnclickListener listener) {
        this.mSizeBtnOnclickListener = listener;
    }

    public void setColorBtnOnclickListener(ColorBtnOnclickListener listener) {
        this.mColorBtnOnclickListener = listener;
    }

    public void setShapeBtnClickListener(ShapeBtnClickListener listener) {
        this.mShapeBtnClickListener = listener;
    }

    public void setPageBtnOnClickListener(PageBtnOnClickListener listener) {
        this.mPageBtnOnClickListener = listener;
    }

    public void setEraserBtnOnClickListener(EraserBtnOnClickListener listener) {
        this.mEraserBtnOnClickListener = listener;
    }

    public void setShapeChooseBtnOnclickListener(ShapeChooseBtnOnclickListener listener) {

        this.mShapeChooseBtnOnclickListener = listener;
    }

    /**
     * 笔迹粗细调整按钮点击监听接口
     */
    public interface SizeBtnOnclickListener {
        void Clicked(boolean isShow);
    }

    /**
     * 笔迹颜色调整按钮点击监听接口
     */
    public interface ColorBtnOnclickListener {
        void ColorClicked(boolean isShow);
    }

    /**
     * 图形选择监听接口
     */
    public interface ShapeBtnClickListener {
        void onShapeBtnClicked(boolean isShow);
    }


    /**
     * 添加多页监听接口
     */
    public interface PageBtnOnClickListener {
        void PageClicked(boolean isShow);
    }

    /**
     * 橡皮擦监听接口
     */
    public interface EraserBtnOnClickListener {
        void EraserClicked();
    }

    /**
     * 笔迹移动监听接口
     */
    public interface ShapeChooseBtnOnclickListener {
        void ChooseOnClicked();
    }

}
