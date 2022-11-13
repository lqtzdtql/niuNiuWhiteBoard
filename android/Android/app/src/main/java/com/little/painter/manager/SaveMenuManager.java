package com.little.painter.manager;

import android.view.View;
import android.widget.TextView;

import com.little.painter.Constants;
import com.little.painter.R;

public class SaveMenuManager implements View.OnClickListener {
    private View mView;
    private TextView mTvSavePng;
    private TextView mTvSaveSvg;
    private SaveBtnClickListener mSaveBtnClickListener;

    public SaveMenuManager(View mView) {
        this.mView = mView;
        initView();
        initEvent();
    }

    private void initEvent() {
        mTvSaveSvg.setOnClickListener(this);
        mTvSavePng.setOnClickListener(this);
    }

    private void initView() {
        mTvSavePng = (TextView) mView.findViewById(R.id.tv_save_png);
        mTvSaveSvg = (TextView) mView.findViewById(R.id.tv_save_svg);
    }

    @Override
    public void onClick(View v) {
        switch (v.getId()) {
            case R.id.tv_save_png:
                if (mSaveBtnClickListener != null) {
                    mSaveBtnClickListener.onSaveClick(Constants.PNG);
                }
                break;
            case R.id.tv_save_svg:
                if (mSaveBtnClickListener != null) {
                    mSaveBtnClickListener.onSaveClick(Constants.SVG);
                }
                break;
            default:
                break;
        }
    }

    public void setSaveBtnClickListener(SaveBtnClickListener listener) {
        this.mSaveBtnClickListener = listener;
    }

    public interface SaveBtnClickListener {
        void onSaveClick(int savekind);
    }
}
