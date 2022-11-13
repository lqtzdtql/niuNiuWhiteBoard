package com.little.painter.activity;

import android.view.View;
import android.webkit.WebView;
import android.widget.ImageView;
import android.widget.TextView;

import com.little.painter.R;

public class UserAgreementActivity extends BaseActivity {
    private TextView mTvBarTitle;
    private ImageView mBackBtn;
    private WebView mAgreementWebView;

    @Override
    protected int getLayoutId() {
        return R.layout.activity_user_agreement;
    }

    @Override
    protected void initView() {
        mBackBtn = findViewById(R.id.bt_back);
        mTvBarTitle = findViewById(R.id.tv_bar_title);
        mTvBarTitle.setText("用户协议");
        mAgreementWebView = findViewById(R.id.agreement_webview);
        mAgreementWebView.loadUrl("file:///android_asset/user_agreement.html");
    }

    @Override
    protected void initEvent() {
        mBackBtn.setOnClickListener(this);
    }

    @Override
    protected void initData() {

    }

    @Override
    public void onClick(View v) {
        switch (v.getId()) {
            case R.id.bt_back:
                finish();
                break;
            default:
                break;
        }
    }

}
