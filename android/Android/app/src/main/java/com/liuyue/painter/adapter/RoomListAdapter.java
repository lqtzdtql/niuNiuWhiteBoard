package com.liuyue.painter.adapter;

import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.ImageView;
import android.widget.TextView;

import androidx.annotation.NonNull;
import androidx.recyclerview.widget.RecyclerView;

import com.liuyue.painter.R;
import com.liuyue.painter.model.HallBean;

import java.util.List;

public class RoomListAdapter extends RecyclerView.Adapter<RecyclerView.ViewHolder> {

    private List<HallBean.RoomlistBean> mDataList;
    private OnClickListener mOnClickListener;

    public RoomListAdapter(List<HallBean.RoomlistBean> list) {
        mDataList = list;
    }

    public void setOnClickListener(OnClickListener listener) {
        mOnClickListener = listener;
    }

    public void notifyData(List<HallBean.RoomlistBean> list) {
        mDataList = list;
        notifyDataSetChanged();
    }

    @NonNull
    @Override
    public RecyclerView.ViewHolder onCreateViewHolder(@NonNull ViewGroup parent, int viewType) {
        View view = LayoutInflater.from(parent.getContext()).inflate(R.layout.item_room, null);
        return new MyViewHolder(view);
    }

    @Override
    public void onBindViewHolder(@NonNull RecyclerView.ViewHolder holder, int position) {
        MyViewHolder myViewHolder = (MyViewHolder) holder;
        HallBean.RoomlistBean roomInfo = mDataList.get(position);
        myViewHolder.mTvRoomName.setText(roomInfo.getName());
        myViewHolder.itemView.setOnClickListener(v -> {
            if (mOnClickListener != null) {
                mOnClickListener.onClick(roomInfo.getUuid());
            }
        });
    }

    @Override
    public int getItemCount() {
        return mDataList.size();
    }

    private static class MyViewHolder extends RecyclerView.ViewHolder {

        public ImageView mIvRoomCover;
        public TextView mTvRoomName;
        public TextView mTvHost;
        public TextView mTvPeopleNumber;

        public MyViewHolder(@NonNull View itemView) {
            super(itemView);
            mIvRoomCover = itemView.findViewById(R.id.iv_room_cover);
            mTvRoomName = itemView.findViewById(R.id.tv_room_name);
            mTvHost = itemView.findViewById(R.id.tv_host_name);
            mTvPeopleNumber = itemView.findViewById(R.id.tv_current_number);
        }
    }

    public interface OnClickListener {
        void onClick(String uuid);
    }
}
