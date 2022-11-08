package com.liuyue.painter.view;

import android.content.Context;
import android.util.AttributeSet;
import android.view.LayoutInflater;
import android.view.View;
import android.widget.ImageView;
import android.widget.LinearLayout;

import com.liuyue.painter.Constants;
import com.liuyue.painter.R;

public class ShapeSelectView extends LinearLayout implements View.OnClickListener {

    // 默认选择的是曲线
    public static int kind;
    private final Context mContext;
    ImageView mSelectInkBtn;
    ImageView mSelectLineBtn;
    ImageView mSelectRectBtn;
    ImageView mSelectCircleBtn;
    private KindBtnClickedListener mKindBtnClickedListener;

    public ShapeSelectView(Context context) {
        this(context, null);
    }


    public ShapeSelectView(Context context, AttributeSet attrs) {
        this(context, attrs, 0);
    }

    public ShapeSelectView(Context context, AttributeSet attrs, int defStyleAttr) {
        super(context, attrs, defStyleAttr);
        this.mContext = context;
        LayoutInflater.from(mContext).inflate(R.layout.shape_palette, this);
        initView();
        initEvent();
    }

    public static int getKind() {
        return kind;
    }

    public static void setKind(int kind) {
        ShapeSelectView.kind = kind;
    }

    private void initView() {
        mSelectInkBtn = (ImageView) findViewById(R.id.id_select_ink);
        mSelectLineBtn = (ImageView) findViewById(R.id.id_select_line);
        mSelectRectBtn = (ImageView) findViewById(R.id.id_select_rect);
        mSelectCircleBtn = (ImageView) findViewById(R.id.id_select_circle);
    }

    private void initEvent() {
        mSelectInkBtn.setOnClickListener(this);
        mSelectLineBtn.setOnClickListener(this);
        mSelectRectBtn.setOnClickListener(this);
        mSelectCircleBtn.setOnClickListener(this);
    }

    @Override
    public void onClick(View v) {
        switch (v.getId()) {
            case R.id.id_select_ink:
                setKind(Constants.INK);
                break;
            case R.id.id_select_line:
                setKind(Constants.LINE);
                break;
            case R.id.id_select_rect:
                setKind(Constants.RECT);
                break;
            case R.id.id_select_circle:
                setKind(Constants.CIRCLE);
                break;
        }
        if (mKindBtnClickedListener != null) {
            mKindBtnClickedListener.onKindBtnClicked(kind);
        }
    }

    public void setKindBtnClickedListener(KindBtnClickedListener listener) {
        this.mKindBtnClickedListener = listener;
    }

    public interface KindBtnClickedListener {
        void onKindBtnClicked(int kind);
    }
}
