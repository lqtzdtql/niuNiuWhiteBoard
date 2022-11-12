package com.liuyue.painter.activity;

import android.os.Environment;
import android.os.Handler;
import android.os.Looper;
import android.os.Message;
import android.view.View;
import android.widget.EditText;
import android.widget.ImageView;
import android.widget.LinearLayout;
import android.widget.Toast;

import androidx.appcompat.app.AlertDialog;

import com.liuyue.painter.Constants;
import com.liuyue.painter.R;
import com.liuyue.painter.callback.PageChangeCall;
import com.liuyue.painter.model.SaveOperation;
import com.liuyue.painter.model.SavePngOperation;
import com.liuyue.painter.model.SaveSvgOperation;
import com.liuyue.painter.utils.SaveUtil;
import com.liuyue.painter.view.ChooseUIManager;
import com.liuyue.painter.view.FootUIManager;
import com.liuyue.painter.view.ArtBoard;
import com.liuyue.painter.view.SaveMenuManager;

public class RoomActivity extends BaseActivity {

    private ArtBoard mArtBoard;
    private ImageView mExitBtn;
    private ImageView mUndoBtn;
    private ImageView mRedoBtn;
    private ImageView mMoreBtn;
    private LinearLayout mFootLayout;
    private FootUIManager mFootUIManager;
    private LinearLayout mSaveMenuLayout;
    private SaveMenuManager mSaveMenuManager;
    private LinearLayout mChooseLayout;
    private ChooseUIManager mChooseUIManager;
    private String saveFileName;
    private SaveOperation so;

    private boolean mIsMenuShow;

    private String mRoomUUID;

    @Override
    protected int getLayoutId() {
        return R.layout.activity_room;
    }

    @Override
    protected void initView() {
        mArtBoard = (ArtBoard) findViewById(R.id.id_paint_view);
        mExitBtn = (ImageView) findViewById(R.id.id_exit_btn);
        mUndoBtn = (ImageView) findViewById(R.id.id_undo_btn);
        mRedoBtn = (ImageView) findViewById(R.id.id_redo_btn);
        mMoreBtn = (ImageView) findViewById(R.id.id_save_btn);

        mFootLayout = (LinearLayout) findViewById(R.id.id_foot_layout);
        mChooseLayout = (LinearLayout) findViewById(R.id.id_choose_layout);
        mSaveMenuLayout = (LinearLayout) findViewById(R.id.save_menu_layout);
        mSaveMenuManager = new SaveMenuManager(this, mSaveMenuLayout);
        mFootUIManager = new FootUIManager(this, mFootLayout);
        mChooseUIManager = new ChooseUIManager(this, mChooseLayout);
    }

    @Override
    protected void initEvent() {
        mExitBtn.setOnClickListener(this);
        mUndoBtn.setOnClickListener(this);
        mRedoBtn.setOnClickListener(this);
        mMoreBtn.setOnClickListener(this);

        // 选择画笔大小事件监听
        mFootUIManager.setSizeBtnListener(mSizeBtnOnclickListener);
        // 选择画笔颜色事件监听
        mFootUIManager.setColorBtnOnclickListener(mColorBtnOnclickListener);
        // 选择形状事件监听
        mFootUIManager.setShapeBtnClickListener(mShapeBtnClickListener);
        // 加页监听
        mFootUIManager.setPageBtnOnClickListener(mPageBtnOnClickListener);
        // 橡皮擦监听
        mFootUIManager.setEraserBtnOnClickListener(mEraserBtnOnClickListener);
        // 笔迹操作监听
        mFootUIManager.setShapeChooseBtnOnclickListener(mShapeChooseBtnOnclickListener);
        // 保存监听
        mSaveMenuManager.setSaveBtnClickListener(mSaveBtnClickListener);
    }

    @Override
    protected void initData() {
        mRoomUUID = getIntent().getStringExtra("roomUUID");
    }

    @Override
    public void onClick(View v) {
        switch (v.getId()) {
            case R.id.id_exit_btn:
                // 退出保存操作
                exit();
                break;
            case R.id.id_undo_btn:
                // 撤销操作
                mArtBoard.undo();
                break;
            case R.id.id_redo_btn:
                // 重做操作
                mArtBoard.redo();
                break;
            case R.id.id_save_btn:
                // 另存为相关操作
                showMenuLogic();
                break;
            default:
                break;
        }
    }

    /**
     * 菜单显示逻辑
     */
    private void showMenuLogic() {
        if (mIsMenuShow) {
            mSaveMenuLayout.setVisibility(View.VISIBLE);
        } else {
            mSaveMenuLayout.setVisibility(View.GONE);
        }
        mIsMenuShow = !mIsMenuShow;
    }

    /**
     * 退出逻辑执行
     */
    private void exit() {
        // 判断是否有内容
        if (!mArtBoard.isEmpty()) {
            // 判断是否处于已存在的文件
            if (mArtBoard.isFileExist()) {
                // 是否发生了修改
                if (mArtBoard.isEdited()) {
                    // 保存最后一张
                    mArtBoard.saveTheLast();
                    // 进入覆盖保存
                    chooseCoverOrNot();
                } else {
                    // 直接退出
                    finish();
                }
            } else {
                // 保存最后一张
                mArtBoard.saveTheLast();
                // 进入选择名字保存
                chooseSaveOrNot();
            }
        } else {
            finish();
        }

    }

    /**
     * 选择是否覆盖已存在图片
     */
    private void chooseCoverOrNot() {
        new AlertDialog.Builder(this)
                .setMessage("是否覆盖当前图片")
                .setNegativeButton("取消", (dialog, which) -> {
                    Message msg = new Message();
                    msg.what = Constants.MSG_EXIT;
                    mHandler.sendMessage(msg);
                    dialog.dismiss();
                })
                .setPositiveButton("确定", (dialog, which) -> {

                    String pathname = mArtBoard.getGallery().getName();
                    SaveUtil suThread = new SaveUtil(mHandler, pathname, mArtBoard);
                    suThread.start();
                    dialog.dismiss();
                }).show();
    }

    /**
     * 是否选择保存绘制图案
     */
    private void chooseSaveOrNot() {
        final EditText et = new EditText(RoomActivity.this);
        new AlertDialog.Builder(this)
                .setMessage("是否保存？")
                .setView(et)
                .setPositiveButton("保存", (dialog, which) -> {
                    if ("".equals(et.getText().toString())) {
                        Toast.makeText(RoomActivity.this, "请输入保存文件的名字", Toast.LENGTH_SHORT).show();
                    } else {
                        // 输入的名字作为文件夹的名字
                        String name = et.getText().toString();
                        SaveUtil suThread = new SaveUtil(mHandler, name, mArtBoard);
                        suThread.start();
                        dialog.dismiss();
                    }
                })
                .setNegativeButton("不保存", (dialog, which) -> {
                    Message msg = new Message();
                    msg.what = Constants.MSG_EXIT;
                    mHandler.sendMessage(msg);
                    dialog.dismiss();
                }).show();
    }

    private final Handler mHandler = new Handler(Looper.getMainLooper()) {
        @Override
        public void handleMessage(Message msg) {
            super.handleMessage(msg);
            if (msg.what == Constants.MSG_EXIT) {
                finish();
            }
            if (msg.what == Constants.MSG_REDRAW) {
                try {
                    mArtBoard.redrawOnBitmap();
                } catch (Exception e) {
                    e.printStackTrace();
                }
            }
        }
    };

    /**
     * 画笔粗细点击事件
     */
    private final FootUIManager.SizeBtnOnclickListener mSizeBtnOnclickListener = new FootUIManager.SizeBtnOnclickListener() {
        @Override
        public void Clicked(boolean isShow) {
            if (isShow) {
                mChooseLayout.removeAllViews();
                mChooseLayout.setVisibility(View.GONE);
            } else {
                mChooseLayout.setVisibility(View.VISIBLE);
                mChooseUIManager.ShowSizeUi(mArtBoard.getBrushSize());
                mChooseUIManager.setSizeCallback(size -> mArtBoard.setBrushSize(size));
            }
        }
    };

    /**
     * 颜色点击事件
     *
     * @param isShow
     */
    private final FootUIManager.ColorBtnOnclickListener mColorBtnOnclickListener = new FootUIManager.ColorBtnOnclickListener() {
        @Override
        public void ColorClicked(boolean isShow) {
            if (isShow) {
                mChooseLayout.removeAllViews();
                mChooseLayout.setVisibility(View.GONE);
            } else {
                mChooseLayout.setVisibility(View.VISIBLE);
                mChooseUIManager.ShowColorUi();
                mChooseUIManager.setColorCallBack(color -> mArtBoard.setBrushColor(color));
            }
        }
    };

    /**
     * 图形点击事件
     */
    private final FootUIManager.ShapeBtnClickListener mShapeBtnClickListener = new FootUIManager.ShapeBtnClickListener() {
        @Override
        public void onShapeBtnClicked(boolean isShow) {
            if (isShow) {
                mChooseLayout.removeAllViews();
                mChooseLayout.setVisibility(View.GONE);
            } else {
                mChooseLayout.setVisibility(View.VISIBLE);
                mChooseUIManager.ShowShapeUi(mArtBoard.getCurrentKind());
                mChooseUIManager.setShapeChangeCall(kind -> mArtBoard.setCurrentKind(kind));
            }
        }
    };

    /**
     * 多页点击事件
     */
    private final FootUIManager.PageBtnOnClickListener mPageBtnOnClickListener = new FootUIManager.PageBtnOnClickListener() {
        @Override
        public void PageClicked(boolean isShow) {
            if (isShow) {
                mChooseLayout.removeAllViews();
                mChooseLayout.setVisibility(View.GONE);
            } else {
                mChooseLayout.setVisibility(View.VISIBLE);
                mChooseUIManager.ShowPageUi(mArtBoard.getCurrentPageNum(), mArtBoard.getCurrentPageIndex());
                // 接口回调
                mChooseUIManager.setPageChangeCall(new PageChangeCall() {
                    @Override
                    public void PageAddCall(int pagenum, int pageindex) {
                        mArtBoard.setCurrentPageNum(pagenum);
                        mArtBoard.setCurrentPageIndex(pageindex);
                        mArtBoard.drawNewImage();
                    }

                    @Override
                    public void PagePreCall(int pageindex) {
                        mArtBoard.setCurrentPageIndex(pageindex);
                        mArtBoard.turnToPrePage();
                    }

                    @Override
                    public void PageNextCall(int pageindex) {
                        mArtBoard.setCurrentPageIndex(pageindex);
                        mArtBoard.turnToNextPage();
                    }
                });
            }
        }
    };

    /**
     * 橡皮擦监听
     */
    private final FootUIManager.EraserBtnOnClickListener mEraserBtnOnClickListener = new FootUIManager.EraserBtnOnClickListener() {
        @Override
        public void EraserClicked() {
            mArtBoard.changeEraserState();
        }
    };

    /**
     * 笔迹操作监听
     */
    private final FootUIManager.ShapeChooseBtnOnclickListener mShapeChooseBtnOnclickListener = new FootUIManager.ShapeChooseBtnOnclickListener() {
        @Override
        public void ChooseOnClicked() {
            mArtBoard.changeCutState();
        }
    };

    /**
     * 另存点击事件
     */
    private final SaveMenuManager.SaveBtnClickListener mSaveBtnClickListener = new SaveMenuManager.SaveBtnClickListener() {
        @Override
        public void onSaveClick(int savekind) {
            if (savekind == Constants.PNG) {
                so = new SavePngOperation();
            } else {
                so = new SaveSvgOperation();
            }
            // 获取需要保存的内容
            so.GetContent(mArtBoard);
            final EditText et = new EditText(RoomActivity.this);
            new AlertDialog.Builder(RoomActivity.this)
                    .setMessage("输入保存文件名")
                    .setView(et)
                    .setPositiveButton("保存", (dialog, which) -> {
                        if ("".equals(et.getText().toString())) {
                            Toast.makeText(RoomActivity.this, "文件名不能为空", Toast.LENGTH_SHORT).show();
                        } else {
                            try {
                                saveFileName = et.getText().toString();
                                // 开启线程保存
                                new Thread(new MyPicSaveRunnable()).start();
                                Message msg = new Message();
                                msg.what = Constants.MSG_REDRAW;
                                mHandler.sendMessageDelayed(msg, 300);
                            } catch (Exception e) {
                                e.printStackTrace();
                            }
                        }

                        dialog.dismiss();
                    })
                    .setNegativeButton("取消", (dialog, which) -> dialog.dismiss()).show();
        }
    };

    private final class MyPicSaveRunnable implements Runnable {

        @Override
        public void run() {
            if (so != null) {
                so.setFilepath(Environment.getExternalStorageDirectory() + "/palette");
                so.setFilename(saveFileName);
                so.GetAbusoluteFileName();
                so.SavePainting();
            }
        }
    }
}
