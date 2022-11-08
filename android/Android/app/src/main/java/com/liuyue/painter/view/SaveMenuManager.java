package com.liuyue.painter.view;

import android.content.Context;
import android.view.View;
import android.widget.TextView;

import com.liuyue.painter.Constants;
import com.liuyue.painter.R;

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
        SavePngTv = (TextView) mView.findViewById(R.id.save_png_btn);
        SaveSvgTv = (TextView) mView.findViewById(R.id.save_svg_btn);
    }

    @Override
    public void onClick(View v) {
        switch (v.getId()) {
            case R.id.save_png_btn:
                if (mSaveBtnClickListener != null)
                    mSaveBtnClickListener.SaveClick(Constants.PNG);
                break;
            case R.id.save_svg_btn:
                if (mSaveBtnClickListener != null)
                    mSaveBtnClickListener.SaveClick(Constants.SVG);
                break;

        }
    }

    public void setSaveBtnClickListener(SaveBtnClickListener listener) {
        this.mSaveBtnClickListener = listener;
    }

    public interface SaveBtnClickListener {
        void SaveClick(int savekind);
    }
}
