package com.little.painter.utils;

import android.os.Handler;
import android.os.Message;

import com.little.painter.Constants;
import com.little.painter.view.ArtBoard;

public class SaveUtil extends Thread {

    private final Handler mHandler;
    private final String mName;
    private final ArtBoard mMycanvas;

    public SaveUtil(Handler mHandler, String mName, ArtBoard mMycanvas) {
        this.mHandler = mHandler;
        this.mName = mName;
        this.mMycanvas = mMycanvas;
    }

    @Override
    public void run() {
        super.run();
        // 创建palette下的文件夹
        // 向文件夹存储xml
        for (int i = 0; i < mMycanvas.getGallery().getPaintingList().size(); i++) {
            // 创建xml文件名
            String xmlfilename = StoreOperation.getXmlFileName(this.mName, i + "");
            // 创建xml文件
            XmlOperation.CreatXml(mMycanvas.getGallery().getPaintingList().get(i), xmlfilename);
            // 向主线程发送信息
            Message msg = new Message();
            msg.what = Constants.MSG_EXIT;
            mHandler.sendMessage(msg);
        }
    }

}
