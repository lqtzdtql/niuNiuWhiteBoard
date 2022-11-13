package com.little.painter.activity;

import android.os.Bundle;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;
import android.widget.RadioButton;
import android.widget.TextView;

import com.blankj.utilcode.util.ThreadUtils;
import com.blankj.utilcode.util.ToastUtils;
import com.little.painter.R;
import com.little.painter.model.CreateRoomBean;
import com.little.painter.utils.AppServer;

public class CreateRoomActivity extends BaseActivity {

    private View mBtBack;
    private TextView mTvTitle;
    private EditText mEtRoomName;
    private RadioButton mRbNeedPassword;
    private TextView mRoomPassword;
    private EditText mEtRoomPassword;
    private Button mBtCreateRoom;

    @Override
    protected int getLayoutId() {
        return R.layout.activity_create_room;
    }

    @Override
    protected void initView() {
        mBtBack = findViewById(R.id.bt_back);
        mTvTitle = findViewById(R.id.tv_bar_title);
        mEtRoomName = findViewById(R.id.et_room_name);
        mRbNeedPassword = findViewById(R.id.rb_need_password);
        mRoomPassword = findViewById(R.id.room_password);
        mEtRoomPassword = findViewById(R.id.et_room_password);
        mBtCreateRoom = findViewById(R.id.bt_create_room);
    }

    @Override
    protected void initEvent() {
        mBtBack.setOnClickListener(this);
        mRbNeedPassword.setOnCheckedChangeListener((buttonView, isChecked) -> {
            if (isChecked) {
                mRoomPassword.setVisibility(View.VISIBLE);
                mEtRoomPassword.setVisibility(View.VISIBLE);
            } else {
                mRoomPassword.setVisibility(View.INVISIBLE);
                mEtRoomPassword.setVisibility(View.INVISIBLE);
            }
        });
        mBtCreateRoom.setOnClickListener(this);
    }

    @Override
    protected void initData() {
        mTvTitle.setText("创建房间");
    }

    @Override
    public void onClick(View view) {
        if (view.getId() == R.id.bt_back) {
            finish();
        } else if (view.getId() == R.id.bt_create_room) {
            ThreadUtils.getSinglePool().execute(() -> {
                String roomName = mEtRoomName.getText().toString();
                CreateRoomBean createRoomBean = AppServer.getInstance().createRoom(roomName, "teaching_room");
                if (createRoomBean != null) {
                    Bundle bundle = new Bundle();
                    bundle.putString("roomUUID", createRoomBean.getUuid());
                    startActivity(RoomActivity.class, bundle);
                    finish();
                } else {
                    ToastUtils.showShort("创建房间失败");
                }
            });
        }
    }
}
