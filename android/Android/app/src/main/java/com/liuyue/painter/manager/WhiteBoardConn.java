package com.liuyue.painter.manager;

import com.google.gson.JsonObject;
import com.liuyue.painter.model.MessageBean;

import org.java_websocket.client.WebSocketClient;
import org.java_websocket.handshake.ServerHandshake;

import java.net.URI;
import java.net.URISyntaxException;
import java.util.concurrent.ScheduledThreadPoolExecutor;
import java.util.concurrent.TimeUnit;

/**
 * 白板连接管理
 */
public class WhiteBoardConn extends WebSocketClient {
    private static final String SERVER_ADDRESS = "http://81.68.68.216:8282";
    private static final long HEAT_INTERVAL = 3 * 1000L;
    private final ScheduledThreadPoolExecutor mHeatService = new ScheduledThreadPoolExecutor(1);
    private final String mUserUUID;
    private final JsonObject mHeatMessage;

    public WhiteBoardConn(String uuid, String roomToken) throws URISyntaxException {
        super(new URI(SERVER_ADDRESS + "/websocket?token=" + roomToken));
        mUserUUID = uuid;
        mHeatMessage = new JsonObject();
        mHeatMessage.addProperty("from", mUserUUID);
        mHeatMessage.addProperty("contentType", MessageBean.HEAT_BEAT);
    }


    @Override
    public void onOpen(ServerHandshake handshakedata) {
        mHeatService.scheduleWithFixedDelay(() -> send(mHeatMessage.toString()), HEAT_INTERVAL, HEAT_INTERVAL, TimeUnit.MILLISECONDS);
    }

    @Override
    public void onMessage(String message) {
        MessageBean messageBean = MessageBean.parse(message);
    }

    @Override
    public void onClose(int code, String reason, boolean remote) {

    }

    @Override
    public void onError(Exception ex) {

    }
}
