package com.liuyue.painter.activity;

import android.os.Bundle;
import android.text.Editable;
import android.text.SpannableString;
import android.text.Spanned;
import android.text.TextWatcher;
import android.text.method.LinkMovementMethod;
import android.text.style.ClickableSpan;
import android.text.style.UnderlineSpan;
import android.view.MotionEvent;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;
import android.widget.RadioButton;
import android.widget.TextView;

import androidx.annotation.NonNull;

import com.blankj.utilcode.util.KeyboardUtils;
import com.blankj.utilcode.util.RegexUtils;
import com.blankj.utilcode.util.SPUtils;
import com.blankj.utilcode.util.ThreadUtils;
import com.blankj.utilcode.util.ToastUtils;
import com.google.android.material.textfield.TextInputLayout;
import com.liuyue.painter.Constants;
import com.liuyue.painter.R;
import com.liuyue.painter.model.LoginBean;
import com.liuyue.painter.model.SignupBean;
import com.liuyue.painter.utils.AppServer;
import com.liuyue.painter.utils.NickNameUtils;

public class LoginActivity extends BaseActivity {
    private EditText mEtPhone;
    private EditText mEtPassword;
    private TextInputLayout mTextInputLayout;
    private Button mBtLogin;
    private RadioButton mRbUserAgreement;
    private TextView mTvUserAgreement;

    private boolean mAgreeUserAgreement = true;

    @Override
    protected int getLayoutId() {
        return R.layout.activity_login;
    }

    @Override
    protected void initView() {
        mEtPhone = findViewById(R.id.et_login_phone);
        mEtPassword = findViewById(R.id.et_login_password);
        mTextInputLayout = findViewById(R.id.layout_login_phone);
        mBtLogin = findViewById(R.id.bt_login_login);
        mRbUserAgreement = findViewById(R.id.rb_login_user_agreement);
        mTvUserAgreement = findViewById(R.id.tv_login_user_agreement);
        initUserAgreement();
    }

    @Override
    protected void initEvent() {
        mRbUserAgreement.setOnClickListener(this);
        mBtLogin.setOnClickListener(this);
        mEtPhone.addTextChangedListener(new PhoneInputWatcher() {
            @Override
            public void afterTextChanged(Editable s) {
                if (RegexUtils.isMobileSimple(s)) {
                    mTextInputLayout.setError(null);
                    mTextInputLayout.setErrorEnabled(false);
                } else {
                    mTextInputLayout.setError("不正确的手机号码");
                    mTextInputLayout.setErrorEnabled(true);
                }
            }
        });
    }

    @Override
    protected void initData() {
        SPUtils spUtils = SPUtils.getInstance(Constants.SP_USER_INFO);
        String userId = spUtils.getString(Constants.KEY_USER_PHONE);
        String password = spUtils.getString(Constants.KEY_USER_PASSWORD);
        mEtPhone.setText(userId);
        mEtPassword.setText(password);
    }

    private void initUserAgreement() {
        SpannableString spannableString = new SpannableString("登录即表示同意用户协议");
        spannableString.setSpan(new UnderlineSpan(), 7, 11, Spanned.SPAN_EXCLUSIVE_EXCLUSIVE);
        ClickableSpan clickableSpan = new ClickableSpan() {
            @Override
            public void onClick(@NonNull View widget) {
                startActivity(UserAgreementActivity.class);
            }
        };
        spannableString.setSpan(clickableSpan, 7, 11, Spanned.SPAN_EXCLUSIVE_EXCLUSIVE);
        mTvUserAgreement.setMovementMethod(LinkMovementMethod.getInstance());
        mTvUserAgreement.setText(spannableString);
    }

    @Override
    public void onClick(View v) {
        switch (v.getId()) {
            case R.id.rb_login_user_agreement:
                // 用户协议
                if (mAgreeUserAgreement) {
                    mRbUserAgreement.setChecked(false);
                    mAgreeUserAgreement = false;
                } else {
                    mRbUserAgreement.setChecked(true);
                    mAgreeUserAgreement = true;
                }
                break;
            case R.id.bt_login_login:
                // 登录
                if (!mRbUserAgreement.isChecked()) {
                    ToastUtils.showShort("请先勾选同意用户协议");
                    return;
                }
                if (!RegexUtils.isMobileSimple(mEtPhone.getText().toString())) {
                    ToastUtils.showShort("手机号码格式不正确");
                    return;
                }
                ThreadUtils.getSinglePool().execute(() -> {
                    String phone = mEtPhone.getText().toString();
                    String password = mEtPassword.getText().toString();
                    LoginBean loginBean = AppServer.getInstance().login(phone, password);
                    Bundle bundle = new Bundle();
                    if (loginBean == null) {
                        String nickName = NickNameUtils.getNickName();
                        SignupBean signupBean = AppServer.getInstance().signup(phone, password, nickName);
                        if (signupBean == null) {
                            ToastUtils.showShort("登录失败");
                        } else {
                            loginBean = AppServer.getInstance().login(phone, password);
                            if (loginBean == null) {
                                ToastUtils.showShort("登录失败");
                            } else {
                                ToastUtils.showShort("登陆成功");
                                SPUtils spUtils = SPUtils.getInstance(Constants.SP_USER_INFO);
                                spUtils.put(Constants.KEY_USER_PHONE, phone);
                                spUtils.put(Constants.KEY_USER_PASSWORD, password);

                                bundle.putSerializable("userInfo", loginBean.getUserInfo());
                                startActivity(HomeActivity.class, bundle);
                                finish();
                            }
                        }
                    } else {
                        ToastUtils.showShort("登陆成功");
                        SPUtils spUtils = SPUtils.getInstance(Constants.SP_USER_INFO);
                        spUtils.put(Constants.KEY_USER_PHONE, phone);
                        spUtils.put(Constants.KEY_USER_PASSWORD, password);

                        bundle.putSerializable("userInfo", loginBean.getUserInfo());
                        startActivity(HomeActivity.class, bundle);
                        finish();
                    }
                });
                break;
            default:
                break;
        }
    }

    @Override
    public boolean dispatchTouchEvent(MotionEvent motionEvent) {
        if (motionEvent.getAction() == MotionEvent.ACTION_DOWN) {
            View view = getCurrentFocus();
            if (isShouldHideKeyboard(view, motionEvent)) {
                KeyboardUtils.hideSoftInput(this);
            }
        }
        return super.dispatchTouchEvent(motionEvent);
    }

    private boolean isShouldHideKeyboard(View view, MotionEvent event) {
        if ((view instanceof EditText)) {
            int[] location = {0, 0};
            view.getLocationOnScreen(location);
            int left = location[0],
                    top = location[1],
                    bottom = top + view.getHeight(),
                    right = left + view.getWidth();
            return !(event.getRawX() > left && event.getRawX() < right
                    && event.getRawY() > top && event.getRawY() < bottom);
        }
        return false;
    }

    private abstract static class PhoneInputWatcher implements TextWatcher {
        @Override
        public void beforeTextChanged(CharSequence s, int start, int count, int after) {
        }

        @Override
        public void onTextChanged(CharSequence s, int start, int before, int count) {
        }
    }
}