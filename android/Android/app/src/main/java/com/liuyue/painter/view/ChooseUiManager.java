package com.liuyue.painter.view;

import android.content.Context;
import android.widget.LinearLayout;
import android.widget.SeekBar;

import com.liuyue.painter.Constants;
import com.liuyue.painter.callback.ColorCallBack;
import com.liuyue.painter.callback.ColorChangeCall;
import com.liuyue.painter.callback.PageChangeCall;
import com.liuyue.painter.callback.ShapeChangeCall;
import com.liuyue.painter.callback.SizeChangeCall;

public class ChooseUiManager {
    Context mContext;
    LinearLayout mView;
    SizeChangeCall mCall;
    ColorCallBack colorCallback;
    ShapeChangeCall shapeChangeCall;
    PageChangeCall pageChangeCall;

    public ChooseUiManager(Context context, LinearLayout layout) {
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
        this.colorCallback = back;
    }

    /**
     * 显示画笔颜色选择组件
     */
    public void ShowColorUi() {
        ColorPalette cp = new ColorPalette(mContext);
        LinearLayout.LayoutParams lp = new LinearLayout.LayoutParams(LinearLayout.LayoutParams.MATCH_PARENT, LinearLayout.LayoutParams.WRAP_CONTENT);
        lp.setMargins(30, 30, 10, 40);
        mView.addView(cp, lp);
        cp.setColorChangeCall(new ColorChangeCall() {
            @Override
            public void callByColorChange(String color) {
                // 将选中的值作为画笔的颜色
                colorCallback.setChangeColor(color);
            }
        });
    }

    /**
     * 注册图形选择回调
     *
     * @param back
     */
    public void setShapeChangeCall(ShapeChangeCall back) {
        this.shapeChangeCall = back;
    }

    /**
     * 显示图形绘制选择组件
     */
    public void ShowShapeUi(int currentkind) {
        ShapeSelectView sv = new ShapeSelectView(mContext);
        LinearLayout.LayoutParams lp = new LinearLayout.LayoutParams(LinearLayout.LayoutParams.MATCH_PARENT, LinearLayout.LayoutParams.WRAP_CONTENT);
        lp.setMargins(30, 30, 10, 30);
        mView.addView(sv, lp);
        ShapeSelectView.setKind(currentkind);
        sv.setKindBtnClickedListener(new ShapeSelectView.KindBtnClickedListener() {
            @Override
            public void onKindBtnClicked(int kind) {
                shapeChangeCall.CallByShapeChange(kind);
            }
        });

    }

    public void setPageChangeCall(PageChangeCall back) {
        this.pageChangeCall = back;
    }

    /**
     * 显示多页添加组件
     */
    public void ShowPageUi(int currentpagenum, int currentpageindex) {
        PageSelectView pv = new PageSelectView(mContext);
        LinearLayout.LayoutParams lp = new LinearLayout.LayoutParams(LinearLayout.LayoutParams.MATCH_PARENT, LinearLayout.LayoutParams.WRAP_CONTENT);
        lp.setMargins(30, 30, 10, 30);
        mView.addView(pv, lp);
        pv.setPagenum(currentpagenum);
        pv.setPageindex(currentpageindex);
        pv.SetPageComponentClickListener(new PageSelectView.PageComponentClickListener() {
            @Override
            public void AddPageClicked(int pagenum, int pageindex) {
                pageChangeCall.PageAddCall(pagenum, pageindex);
            }

            @Override
            public void PrePageClicked(int pageindex) {

                pageChangeCall.PagePreCall(pageindex);
            }

            @Override
            public void NextPageClicked(int pageindex) {
                pageChangeCall.PageNextCall(pageindex);
            }
        });
    }
}
