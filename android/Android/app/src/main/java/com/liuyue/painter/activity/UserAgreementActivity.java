package com.liuyue.painter.activity;

import android.view.View;
import android.webkit.WebView;
import android.widget.ImageView;
import android.widget.TextView;

import com.blankj.utilcode.util.BarUtils;
import com.liuyue.painter.R;

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
        initStatusBar();
        mBackBtn = findViewById(R.id.back_btn);
        mTvBarTitle = findViewById(R.id.tv_bar_title);
        mTvBarTitle.setText("用户协议");
        mAgreementWebView = findViewById(R.id.agreement_webview);
        mAgreementWebView.loadUrl("file:///android_asset/user_agreement.html");
    }

    private void initStatusBar() {
        BarUtils.setStatusBarColor(this, getResources().getColor(R.color.blue, null));
        // 设置状态栏文字颜色及图标为浅色
        getWindow().getDecorView().setSystemUiVisibility(View.SYSTEM_UI_FLAG_LAYOUT_STABLE);
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
            case R.id.back_btn:
                finish();
                break;
            default:
                break;
        }
    }

}
