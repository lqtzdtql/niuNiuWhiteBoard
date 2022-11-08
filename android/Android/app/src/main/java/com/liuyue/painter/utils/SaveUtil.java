package com.liuyue.painter.utils;

import android.os.Handler;
import android.os.Message;

import com.liuyue.painter.Constants;
import com.liuyue.painter.view.MyCanvas;

public class SaveUtil extends Thread {

    private final Handler mHandler;
    private final String mName;
    private final MyCanvas mMycanvas;

    public SaveUtil(Handler mHandler, String mName, MyCanvas mMycanvas) {
        this.mHandler = mHandler;
        this.mName = mName;
        this.mMycanvas = mMycanvas;
    }

    @Override
    public void run() {
        super.run();
        // 创建palette下的文件夹
        // 向文件夹存储xml
        for (int i = 0; i < mMycanvas.getmGallery().getPaintingList().size(); i++) {
            // 创建xml文件名
            String xmlfilename = StoreOperation.getXmlFileName(this.mName, i + "");
            // 创建xml文件
            XmlOperation.CreatXml(mMycanvas.getmGallery().getPaintingList().get(i), xmlfilename);
            // 向主线程发送信息
            Message msg = new Message();
            msg.what = Constants.MSG_EXIT;
            mHandler.sendMessage(msg);
        }
    }

}
