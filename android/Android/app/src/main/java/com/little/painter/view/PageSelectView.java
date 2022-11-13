package com.little.painter.view;

import android.content.Context;
import android.util.AttributeSet;
import android.view.LayoutInflater;
import android.view.View;
import android.widget.ImageView;
import android.widget.LinearLayout;
import android.widget.TextView;

import com.little.painter.R;

public class PageSelectView extends LinearLayout implements View.OnClickListener {
    private int mPageNum = 0;
    private int mPageIndex = 0;

    private ImageView mAddPageBtn;
    private ImageView mTurnPrePageBtn;
    private ImageView mTurnNextPageBtn;
    private TextView mPageNumText;
    private PageComponentClickListener mPageComponentClickListener;

    public PageSelectView(Context context) {
        this(context, null);
    }

    public PageSelectView(Context context, AttributeSet attrs) {
        this(context, attrs, 0);
    }


    public PageSelectView(Context context, AttributeSet attrs, int defStyleAttr) {
        super(context, attrs, defStyleAttr);
        LayoutInflater.from(context).inflate(R.layout.add_page_layout, this);
        initView();
        initEvent();
    }

    public int getPageNum() {
        return mPageNum;
    }

    public void setPageNum(int pageNum) {
        this.mPageNum = pageNum;
    }

    public int getPageIndex() {
        return mPageIndex;
    }

    public void setPageIndex(int pageIndex) {
        this.mPageIndex = pageIndex;
        mPageNumText.setText(pageIndex + "");
    }

    private void initView() {
        mAddPageBtn = (ImageView) findViewById(R.id.id_add_page);
        mTurnPrePageBtn = (ImageView) findViewById(R.id.id_pre_page);
        mPageNumText = (TextView) findViewById(R.id.id_page_num);
        mTurnNextPageBtn = (ImageView) findViewById(R.id.id_next_page);
    }

    private void initEvent() {
        mAddPageBtn.setOnClickListener(this);
        mTurnPrePageBtn.setOnClickListener(this);
        mTurnNextPageBtn.setOnClickListener(this);
        mPageNumText.setOnClickListener(this);
    }

    @Override
    public void onClick(View v) {
        switch (v.getId()) {
            case R.id.id_add_page:
                // 加页，显示最新一页
                mPageNum++;
                mPageIndex = mPageNum;
                mPageNumText.setText(mPageIndex + "");
                if (mPageComponentClickListener != null) {
                    mPageComponentClickListener.AddPageClicked(mPageNum, mPageIndex);
                }
                break;
            case R.id.id_pre_page:
                // 回调前一页
                if (mPageIndex > 1) {
                    mPageIndex--;
                    mPageNumText.setText(mPageIndex + "");
                    if (mPageComponentClickListener != null) {
                        mPageComponentClickListener.PrePageClicked(mPageIndex);
                    }
                }
                break;
            case R.id.id_page_num:
                // TODO 显示浏览组件 暂时不写
                break;
            case R.id.id_next_page:
                // 调至下一页
                if (mPageIndex < mPageNum) {
                    mPageIndex++;
                    mPageNumText.setText(mPageIndex + "");
                    if (mPageComponentClickListener != null) {
                        mPageComponentClickListener.NextPageClicked(mPageIndex);
                    }
                }
                break;
            default:
                break;
        }
    }

    public void SetPageComponentClickListener(PageComponentClickListener listener) {
        this.mPageComponentClickListener = listener;
    }

    public interface PageComponentClickListener {

        void AddPageClicked(int pagenum, int pageindex);

        void PrePageClicked(int pageindex);

        void NextPageClicked(int pageindex);

    }
}
