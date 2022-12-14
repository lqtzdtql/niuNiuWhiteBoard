package com.little.painter.activity;

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
import com.little.painter.Constants;
import com.little.painter.R;
import com.little.painter.callback.PageChangeCall;
import com.little.painter.manager.ChooseUIManager;
import com.little.painter.manager.FootUIManager;
import com.little.painter.manager.SaveMenuManager;
import com.little.painter.manager.WhiteBoardConn;
import com.little.painter.model.LoginBean;
import com.little.painter.model.SaveOperation;
import com.little.painter.model.SavePngOperation;
import com.little.painter.model.SaveSvgOperation;
import com.little.painter.model.WhiteBoardAuthBean;
import com.little.painter.shape.Shape;
import com.little.painter.utils.AppServer;
import com.little.painter.utils.SaveUtil;
import com.little.painter.view.ArtBoard;

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
        mSaveMenuManager = new SaveMenuManager(mSaveMenuLayout);
        mFootUIManager = new FootUIManager(footLayout);
        mChooseUIManager = new ChooseUIManager(this, mChooseLayout);
    }

    @Override
    protected void initEvent() {
        mExitBtn.setOnClickListener(this);
        mUndoBtn.setOnClickListener(this);
        mRedoBtn.setOnClickListener(this);
        mMoreBtn.setOnClickListener(this);

        mArtBoard.setShapeChangeListener(mShapeChangeListener);
        // ??????????????????????????????
        mFootUIManager.setSizeBtnListener(mSizeBtnOnclickListener);
        // ??????????????????????????????
        mFootUIManager.setColorBtnOnclickListener(mColorBtnOnclickListener);
        // ????????????????????????
        mFootUIManager.setShapeBtnClickListener(mShapeBtnClickListener);
        // ????????????
        mFootUIManager.setPageBtnOnClickListener(mPageBtnOnClickListener);
        // ???????????????
        mFootUIManager.setEraserBtnOnClickListener(mEraserBtnOnClickListener);
        // ??????????????????
        mFootUIManager.setShapeChooseBtnOnclickListener(mShapeChooseBtnOnclickListener);
        // ????????????
        mSaveMenuManager.setSaveBtnClickListener(mSaveBtnClickListener);
    }

    @Override
    protected void initData() {
        mRoomUUID = getIntent().getStringExtra("roomUUID");
        mUserInfoBean = (LoginBean.UserInfoBean) getIntent().getSerializableExtra("userInfo");
        ThreadUtils.getSinglePool().execute(() -> {
            String whiteBoardToken = AppServer.getInstance().enterRoom(mUserInfoBean.getUuid());
            if (whiteBoardToken == null) {
                ToastUtils.showShort("??????????????????");
                finish();
                return;
            }
            WhiteBoardAuthBean authBean = AppServer.getInstance().authWhiteBoard(whiteBoardToken);
            if (authBean == null) {
                ToastUtils.showShort("??????????????????");
                finish();
                return;
            }
            try {
                mWhiteBoardConn = new WhiteBoardConn(authBean.getUserUUID(), whiteBoardToken);
                mWhiteBoardConn.connect();
            } catch (URISyntaxException e) {
                ToastUtils.showShort("??????????????????");
                finish();
            }
        });
    }

    @Override
    public void onClick(View v) {
        switch (v.getId()) {
            case R.id.iv_back:
                finish();
                break;
            case R.id.iv_undo:
                // ????????????
                mArtBoard.undo();
                break;
            case R.id.iv_redo:
                // ????????????
                mArtBoard.redo();
                break;
            case R.id.iv_setting:
                // ?????????????????????
                showMenuLogic();
                break;
            default:
                break;
        }
    }

    /**
     * ??????????????????
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
     * ?????????????????????????????????
     */
    private void chooseCoverOrNot() {
        new AlertDialog.Builder(this)
                .setMessage("????????????????????????")
                .setNegativeButton("??????", (dialog, which) -> {
                    Message msg = new Message();
                    msg.what = Constants.MSG_EXIT;
                    mHandler.sendMessage(msg);
                    dialog.dismiss();
                })
                .setPositiveButton("??????", (dialog, which) -> {
                    String pathname = mArtBoard.getGallery().getName();
                    SaveUtil suThread = new SaveUtil(mHandler, pathname, mArtBoard);
                    suThread.start();
                    dialog.dismiss();
                }).show();
    }

    /**
     * ??????????????????????????????
     */
    private void chooseSaveOrNot() {
        final EditText et = new EditText(RoomActivity.this);
        new AlertDialog.Builder(this)
                .setMessage("???????????????")
                .setView(et)
                .setPositiveButton("??????", (dialog, which) -> {
                    if ("".equals(et.getText().toString())) {
                        Toast.makeText(RoomActivity.this, "??????????????????????????????", Toast.LENGTH_SHORT).show();
                    } else {
                        // ???????????????????????????????????????
                        String name = et.getText().toString();
                        SaveUtil suThread = new SaveUtil(mHandler, name, mArtBoard);
                        suThread.start();
                        dialog.dismiss();
                    }
                })
                .setNegativeButton("?????????", (dialog, which) -> {
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

    private final ArtBoard.ShapeChangeListener mShapeChangeListener = new ArtBoard.ShapeChangeListener() {
        @Override
        public void onAddShape(Shape shape) {
            // mWhiteBoardConn.addShape(shape);
        }

        @Override
        public void onMoveShape(Shape shape) {
            // mWhiteBoardConn.moveShape(shape);
        }

        @Override
        public void onDeleteShape(Shape shape) {
            // mWhiteBoardConn.deleteShape(shape);
        }
    };

    /**
     * ????????????????????????
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
     * ??????????????????
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
                mChooseUIManager.setColorChangeCall(color -> mArtBoard.setBrushColor(color));
            }
        }
    };

    /**
     * ??????????????????
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
     * ??????????????????
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
                // ????????????
                mChooseUIManager.setPageChangeCall(new PageChangeCall() {
                    @Override
                    public void onPageAddCall(int pagenum, int pageindex) {
                        mArtBoard.setCurrentPageNum(pagenum);
                        mArtBoard.setCurrentPageIndex(pageindex);
                        mArtBoard.drawNewImage();
                    }

                    @Override
                    public void onPagePreCall(int pageindex) {
                        mArtBoard.setCurrentPageIndex(pageindex);
                        mArtBoard.turnToPrePage();
                    }

                    @Override
                    public void onPageNextCall(int pageindex) {
                        mArtBoard.setCurrentPageIndex(pageindex);
                        mArtBoard.turnToNextPage();
                    }
                });
            }
        }
    };

    /**
     * ???????????????
     */
    private final FootUIManager.EraserBtnOnClickListener mEraserBtnOnClickListener = new FootUIManager.EraserBtnOnClickListener() {
        @Override
        public void EraserClicked() {
            mArtBoard.changeEraserState();
        }
    };

    /**
     * ??????????????????
     */
    private final FootUIManager.ShapeChooseBtnOnclickListener mShapeChooseBtnOnclickListener = new FootUIManager.ShapeChooseBtnOnclickListener() {
        @Override
        public void ChooseOnClicked() {
            mArtBoard.changeCutState();
        }
    };

    /**
     * ??????????????????
     */
    private final SaveMenuManager.SaveBtnClickListener mSaveBtnClickListener = new SaveMenuManager.SaveBtnClickListener() {
        @Override
        public void onSaveClick(int savekind) {
            if (savekind == Constants.PNG) {
                mSaveOperation = new SavePngOperation();
            } else {
                mSaveOperation = new SaveSvgOperation();
            }
            // ???????????????????????????
            mSaveOperation.GetContent(mArtBoard);
            final EditText et = new EditText(RoomActivity.this);
            new AlertDialog.Builder(RoomActivity.this)
                    .setMessage("?????????????????????")
                    .setView(et)
                    .setPositiveButton("??????", (dialog, which) -> {
                        if ("".equals(et.getText().toString())) {
                            Toast.makeText(RoomActivity.this, "?????????????????????", Toast.LENGTH_SHORT).show();
                        } else {
                            try {
                                mSaveFileName = et.getText().toString();
                                // ??????????????????
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
                    .setNegativeButton("??????", (dialog, which) -> dialog.dismiss()).show();
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
