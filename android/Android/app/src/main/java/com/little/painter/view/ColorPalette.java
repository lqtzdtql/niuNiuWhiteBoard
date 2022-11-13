package com.little.painter.view;

import android.content.Context;
import android.graphics.Color;
import android.util.AttributeSet;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.BaseAdapter;
import android.widget.GridView;
import android.widget.LinearLayout;

import com.little.painter.Constants;
import com.little.painter.R;
import com.little.painter.callback.ColorChangeCall;

import de.hdodenhof.circleimageview.CircleImageView;

public class ColorPalette extends LinearLayout {
    private GridView mGridview;
    private Context mContext;
    private ColorChangeCall mCall;

    public ColorPalette(Context context) {
        this(context, null);
    }

    public ColorPalette(Context context, AttributeSet attrs) {
        this(context, attrs, 0);
    }

    public ColorPalette(Context context, AttributeSet attrs, int defStyleAttr) {
        super(context, attrs, defStyleAttr);
        this.mContext = context;
        LayoutInflater.from(context).inflate(R.layout.color_palette, this);
        initView();
        initEvent();
    }

    public void setColorChangeCall(ColorChangeCall call) {
        this.mCall = call;
    }

    private void initEvent() {
        mGridview.setOnItemClickListener((parent, view, position, id) -> {
            mCall.callByColorChange(Constants.colors[position]);
        });
    }

    private void initView() {
        mGridview = (GridView) findViewById(R.id.id_color_gridview);
        mGridview.setAdapter(new BaseAdapter() {
            @Override
            public int getCount() {
                return Constants.colors.length;
            }

            @Override
            public Object getItem(int position) {
                return Constants.colors[position];
            }

            @Override
            public long getItemId(int position) {
                return position;
            }

            @Override
            public View getView(int position, View convertView, ViewGroup parent) {
                ViewHolder holder;
                if (convertView == null) {
                    convertView = LayoutInflater.from(mContext).inflate(R.layout.color_palette_item, null);
                    holder = new ViewHolder();
                    holder.mCircleImageView = (CircleImageView) convertView.findViewById(R.id.id_color_item);
                    convertView.setTag(holder);
                } else {
                    holder = (ViewHolder) convertView.getTag();
                }
                holder.mCircleImageView.setBorderColor(Color.parseColor(Constants.colors[position]));
                return convertView;
            }
        });

    }

    public static class ViewHolder {
        CircleImageView mCircleImageView;
    }
}

