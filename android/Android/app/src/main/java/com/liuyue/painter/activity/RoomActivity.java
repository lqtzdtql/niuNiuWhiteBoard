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

import com.blankj.utilcode.util.ThreadUtils;
import com.blankj.utilcode.util.ToastUtils;
import com.liuyue.painter.Constants;
import com.liuyue.painter.R;
import com.liuyue.painter.callback.PageChangeCall;
import com.liuyue.painter.manager.WhiteBoardConn;
import com.liuyue.painter.model.LoginBean;
import com.liuyue.painter.model.SaveOperation;
import com.liuyue.painter.model.SavePngOperation;
import com.liuyue.painter.model.SaveSvgOperation;
import com.liuyue.painter.model.WhiteBoardAuthBean;
import com.liuyue.painter.utils.AppServer;
import com.liuyue.painter.utils.SaveUtil;
import com.liuyue.painter.view.ArtBoard;
import com.liuyue.painter.manager.ChooseUIManager;
import com.liuyue.painter.manager.FootUIManager;
import com.liuyue.painter.manager.SaveMenuManager;

import java.net.URISyntaxException;

public class RoomActivity extends BaseActivity {
    private ArtBoard mArtBoard;
    private ImageView mExitBtn;
    private ImageView mUndoBtn;
    private ImageView mRedoBtn;
    private ImageView mMoreBtn;
    private FootUIManager mFootUIManager;
    private LinearLayout mSaveMenuLayout;
    private SaveMenuManager mSaveMenuManager;
    private LinearLayout mChooseLayout;
    private ChooseUIManager mChooseUIManager;

    private boolean mIsMenuShow;
    private String mSaveFileName;
    private SaveOperation mSaveOperation;

    private String mRoomUUID;
    private LoginBean.UserInfoBean mUserInfoBean;
    private WhiteBoardConn mWhiteBoardConn;

    @Override
    protected int getLayoutId() {
        return R.layout.activity_room;
    }

    @Override
    protected void initView() {
        mArtBoard = (ArtBoard) findViewById(R.id.paint_view);
        mExitBtn = (ImageView) findViewById(R.id.iv_back);
        mUndoBtn = (ImageView) findViewById(R.id.iv_undo);
        mRedoBtn = (ImageView) findViewById(R.id.iv_redo);
        mMoreBtn = (ImageView) findViewById(R.id.iv_setting);

        LinearLayout footLayout = (LinearLayout) findViewById(R.id.ll_foot);
        mChooseLayout = (LinearLayout) findViewById(R.id.ll_choose_panel);
        mSaveMenuLayout = (LinearLayout) findViewById(R.id.ll_setting_menu);
        mSaveMenuManager = new SaveMenuManager(this, mSaveMenuLayout);
        mFootUIManager = new FootUIManager(this, footLayout);
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
        mUserInfoBean = (LoginBean.UserInfoBean) getIntent().getSerializableExtra("userInfo");
        ThreadUtils.getSinglePool().execute(() -> {
            String whiteBoardToken = AppServer.getInstance().enterRoom(mUserInfoBean.getUuid());
            if (whiteBoardToken == null) {
                ToastUtils.showShort("加入房间失败");
                finish();
                return;
            }
            WhiteBoardAuthBean authBean = AppServer.getInstance().authWhiteBoard(whiteBoardToken);
            if (authBean == null) {
                ToastUtils.showShort("加入房间失败");
                finish();
                return;
            }
            try {
                mWhiteBoardConn = new WhiteBoardConn(authBean.getUserUUID(),whiteBoardToken);
                mWhiteBoardConn.connect();
            } catch (URISyntaxException e) {
                ToastUtils.showShort("加入房间失败");
                finish();
            }
        });
    }

    @Override
    public void onClick(View v) {
        switch (v.getId()) {
            case R.id.iv_back:
                // 退出保存操作
                exit();
                break;
            case R.id.iv_undo:
                // 撤销操作
                mArtBoard.undo();
                break;
            case R.id.iv_redo:
                // 重做操作
                mArtBoard.redo();
                break;
            case R.id.iv_setting:
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
                mSaveOperation = new SavePngOperation();
            } else {
                mSaveOperation = new SaveSvgOperation();
            }
            // 获取需要保存的内容
            mSaveOperation.GetContent(mArtBoard);
            final EditText et = new EditText(RoomActivity.this);
            new AlertDialog.Builder(RoomActivity.this)
                    .setMessage("输入保存文件名")
                    .setView(et)
                    .setPositiveButton("保存", (dialog, which) -> {
                        if ("".equals(et.getText().toString())) {
                            Toast.makeText(RoomActivity.this, "文件名不能为空", Toast.LENGTH_SHORT).show();
                        } else {
                            try {
                                mSaveFileName = et.getText().toString();
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

    @Override
    protected void onDestroy() {
        super.onDestroy();
        ThreadUtils.getSinglePool().execute(() -> AppServer.getInstance().exitRoom(mRoomUUID));
    }

    private final class MyPicSaveRunnable implements Runnable {

        @Override
        public void run() {
            if (mSaveOperation != null) {
                mSaveOperation.setFilepath(Environment.getExternalStorageDirectory() + "/palette");
                mSaveOperation.setFilename(mSaveFileName);
                mSaveOperation.GetAbusoluteFileName();
                mSaveOperation.SavePainting();
            }
        }
    }
}
