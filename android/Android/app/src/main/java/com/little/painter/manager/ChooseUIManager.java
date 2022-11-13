package com.little.painter.manager;

import android.content.Context;
import android.widget.LinearLayout;
import android.widget.SeekBar;

import com.little.painter.Constants;
import com.little.painter.callback.ColorCallBack;
import com.little.painter.callback.PageChangeCall;
import com.little.painter.callback.ShapeChangeCall;
import com.little.painter.callback.SizeChangeCall;
import com.little.painter.view.ColorPalette;
import com.little.painter.view.PageSelectView;
import com.little.painter.view.ShapeSelectView;

public class ChooseUIManager {
    private final Context mContext;
    private final LinearLayout mView;
    private SizeChangeCall mCall;
    private ColorCallBack mColorCallBack;
    private ShapeChangeCall mShapeChangeCall;
    private PageChangeCall mPageChangeCall;

    public ChooseUIManager(Context context, LinearLayout layout) {
        this.mContext = context;
        this.mView = layout;
    }

    /**
     * 注册画笔大小设定回调
     *
     * @param call
     */
    public void setSizeCallback(SizeChangeCall call) {
        this.mCall = call;
    }

    /**
     * 显示画笔尺寸选择组件
     */
    public void ShowSizeUi(float currentsize) {
        SeekBar seekBar = new SeekBar(mContext);
        seekBar.setMax(Constants.maxBrushSize);
        seekBar.setMin(Constants.minBrushSize);
        seekBar.setMinimumHeight(10);
        seekBar.setProgress(Math.round(currentsize));
        LinearLayout.LayoutParams lp = new LinearLayout.LayoutParams(LinearLayout.LayoutParams.MATCH_PARENT, LinearLayout.LayoutParams.WRAP_CONTENT);
        lp.setMargins(100, 0, 100, 30);
        mView.addView(seekBar, lp);
        seekBar.setOnSeekBarChangeListener(new SeekBar.OnSeekBarChangeListener() {

            @Override
            public void onProgressChanged(SeekBar seekBar, int progress, boolean fromUser) {

            }

            @Override
            public void onStartTrackingTouch(SeekBar seekBar) {

            }

            @Override
            public void onStopTrackingTouch(SeekBar seekBar) {
                // 将停下来时候的值作为画笔的当前粗细大小
                int currentnum = seekBar.getProgress();
                mCall.callBySizeChange(currentnum);
            }
        });
    }

    /**
     * 注册颜色修改回调
     *
     * @param back
     */
    public void setColorCallBack(ColorCallBack back) {
        this.mColorCallBack = back;
    }

    /**
     * 显示画笔颜色选择组件
     */
    public void ShowColorUi() {
        ColorPalette cp = new ColorPalette(mContext);
        LinearLayout.LayoutParams lp = new LinearLayout.LayoutParams(LinearLayout.LayoutParams.MATCH_PARENT, LinearLayout.LayoutParams.WRAP_CONTENT);
        lp.setMargins(30, 30, 10, 40);
        mView.addView(cp, lp);
        cp.setColorChangeCall(color -> {
            // 将选中的值作为画笔的颜色
            mColorCallBack.setChangeColor(color);
        });
    }

    /**
     * 注册图形选择回调
     *
     * @param back
     */
    public void setShapeChangeCall(ShapeChangeCall back) {
        this.mShapeChangeCall = back;
    }

    /**
     * 显示图形绘制选择组件
     */
    public void ShowShapeUi(int currentkind) {
        ShapeSelectView shapeSelectView = new ShapeSelectView(mContext);
        LinearLayout.LayoutParams lp = new LinearLayout.LayoutParams(LinearLayout.LayoutParams.MATCH_PARENT, LinearLayout.LayoutParams.WRAP_CONTENT);
        lp.setMargins(30, 30, 10, 30);
        mView.addView(shapeSelectView, lp);
        shapeSelectView.setKind(currentkind);
        shapeSelectView.setKindBtnClickedListener(kind -> mShapeChangeCall.CallByShapeChange(kind));

    }

    public void setPageChangeCall(PageChangeCall back) {
        this.mPageChangeCall = back;
    }

    /**
     * 显示多页添加组件
     */
    public void ShowPageUi(int currentpagenum, int currentpageindex) {
        PageSelectView pv = new PageSelectView(mContext);
        LinearLayout.LayoutParams lp = new LinearLayout.LayoutParams(LinearLayout.LayoutParams.MATCH_PARENT, LinearLayout.LayoutParams.WRAP_CONTENT);
        lp.setMargins(30, 30, 10, 30);
        mView.addView(pv, lp);
        pv.setPageNum(currentpagenum);
        pv.setPageIndex(currentpageindex);
        pv.SetPageComponentClickListener(new PageSelectView.PageComponentClickListener() {
            @Override
            public void AddPageClicked(int pagenum, int pageindex) {
                mPageChangeCall.PageAddCall(pagenum, pageindex);
            }

            @Override
            public void PrePageClicked(int pageindex) {
                mPageChangeCall.PagePreCall(pageindex);
            }

            @Override
            public void NextPageClicked(int pageindex) {
                mPageChangeCall.PageNextCall(pageindex);
            }
        });
    }
}
