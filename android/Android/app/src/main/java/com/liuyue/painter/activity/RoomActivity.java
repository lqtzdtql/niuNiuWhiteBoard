package com.liuyue.painter.activity;

import android.view.View;

import com.liuyue.painter.R;

public class RoomActivity extends BaseActivity {

    private String mRoomUUID;

    @Override
    protected int getLayoutId() {
        return R.layout.activity_room;
    }

    @Override
    protected void initView() {

    }

    @Override
    protected void initEvent() {

    }

    @Override
    protected void initData() {
        mRoomUUID = getIntent().getStringExtra("roomUUID");
    }

    @Override
    public void onClick(View v) {

    }
}
