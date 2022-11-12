package com.liuyue.painter.activity;

import android.graphics.Rect;
import android.os.Bundle;
import android.view.View;

import androidx.recyclerview.widget.LinearLayoutManager;
import androidx.recyclerview.widget.RecyclerView;

import com.blankj.utilcode.util.ThreadUtils;
import com.google.android.material.floatingactionbutton.FloatingActionButton;
import com.liuyue.painter.R;
import com.liuyue.painter.adapter.RoomListAdapter;
import com.liuyue.painter.model.HallBean;
import com.liuyue.painter.model.LoginBean;
import com.liuyue.painter.utils.AppServer;

import java.util.ArrayList;

public class HomeActivity extends BaseActivity {

    private LoginBean.UserInfoBean mUserInfoBean;

    private RecyclerView mRecyclerView;
    private RoomListAdapter mRoomListAdapter;
    private View mEmptyView;
    private FloatingActionButton mFloatingActionButton;

    @Override
    protected int getLayoutId() {
        return R.layout.activity_home;
    }

    @Override
    protected void initView() {
        mRecyclerView = findViewById(R.id.recycler_view);
        mEmptyView = findViewById(R.id.empty_view);
        mFloatingActionButton = findViewById(R.id.fab);

        mRoomListAdapter = new RoomListAdapter(new ArrayList<>());
        mRoomListAdapter.setOnClickListener(uuid -> {
            Bundle bundle = new Bundle();
            bundle.putString("roomUUID", uuid);
            startActivity(RoomActivity.class, bundle);
        });
        mRecyclerView.setLayoutManager(new LinearLayoutManager(this));
        mRecyclerView.addItemDecoration(new SpacesItemDecoration(15));
        mRecyclerView.setAdapter(mRoomListAdapter);
    }

    @Override
    protected void initEvent() {
        mFloatingActionButton.setOnClickListener(this);
    }

    @Override
    protected void initData() {
        mUserInfoBean = (LoginBean.UserInfoBean) getIntent().getSerializableExtra("userInfo");
        ThreadUtils.getSinglePool().execute(() -> {
            HallBean hallBean = AppServer.getInstance().getRoomList();
            if (hallBean == null) {
                mRecyclerView.setVisibility(View.INVISIBLE);
                mEmptyView.setVisibility(View.VISIBLE);
            } else {
                mEmptyView.setVisibility(View.INVISIBLE);
                mRecyclerView.setVisibility(View.VISIBLE);
                mRoomListAdapter.notifyData(hallBean.getRoomList());
            }
        });
    }

    @Override
    public void onClick(View view) {
        if (view.getId() == R.id.fab) {
            startActivity(CreateRoomActivity.class);
        }
    }

    @Override
    protected void onDestroy() {
        super.onDestroy();
        ThreadUtils.getSinglePool().execute(() -> AppServer.getInstance().exitRoom(mUserInfoBean.getUuid()));
    }

    public static class SpacesItemDecoration extends RecyclerView.ItemDecoration {
        private final int space;

        public SpacesItemDecoration(int space) {
            this.space = space;
        }

        @Override
        public void getItemOffsets(Rect outRect, View view, RecyclerView parent, RecyclerView.State state) {
            if (parent.getChildAdapterPosition(view) != 0) {
                outRect.top = space;
            }
        }
    }
}
