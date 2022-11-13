package com.little.painter.manager;

import android.content.Context;
import android.view.View;
import android.widget.TextView;

import com.little.painter.Constants;
import com.little.painter.R;

public class SaveMenuManager implements View.OnClickListener {

    Context mContext;
    View mView;
    TextView SavePngTv;
    TextView SaveSvgTv;
    SaveBtnClickListener mSaveBtnClickListener;

    public SaveMenuManager(Context mContext, View mView) {
        this.mContext = mContext;
        this.mView = mView;
        initView();
        initEvent();
    }

    private void initEvent() {
        SaveSvgTv.setOnClickListener(this);
        SavePngTv.setOnClickListener(this);
    }

    private void initView() {
        SavePngTv = (TextView) mView.findViewById(R.id.tv_save_png);
        SaveSvgTv = (TextView) mView.findViewById(R.id.tv_save_svg);
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
