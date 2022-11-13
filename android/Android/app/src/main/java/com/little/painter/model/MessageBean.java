package com.little.painter.model;

import com.google.gson.Gson;

public class MessageBean {
    public static final int HEAT_BEAT = 1;
    public static final int UPDATE_BOARD = 2;
    public static final int OBJECT_NEW = 3;
    public static final int OBJECT_MODIFY = 4;
    public static final int OBJECT_DELETE = 5;
    public static final int SWITCH_BOARD = 6;
    public static final int DRAWING_LOCK = 7;
    public static final int CREATE_BOARD = 8;
    public static final int CAN_LOCK = 9;
    public static final int LEAVE_ROOM = 10;
    public static final int CANVAS_LIST = 11;

    private String from;
    private String to;
    private String toWhiteBoard;
    private String toUser;
    private String objectId;
    private int contentType;
    private String content;
    private long timestamp;
    private boolean isLock;
    private boolean readOnly;
    private String leaveUser;

    public static MessageBean parse(String str) {
        return new Gson().fromJson(str, MessageBean.class);
    }

    public String getFrom() {
        return from;
    }

    public void setFrom(String from) {
        this.from = from;
    }

    public String getTo() {
        return to;
    }

    public void setTo(String to) {
        this.to = to;
    }

    public String getToWhiteBoard() {
        return toWhiteBoard;
    }

    public void setToWhiteBoard(String toWhiteBoard) {
        this.toWhiteBoard = toWhiteBoard;
    }

    public String getToUser() {
        return toUser;
    }

    public void setToUser(String toUser) {
        this.toUser = toUser;
    }

    public String getObjectId() {
        return objectId;
    }

    public void setObjectId(String objectId) {
        this.objectId = objectId;
    }

    public int getContentType() {
        return contentType;
    }

    public void setContentType(int contentType) {
        this.contentType = contentType;
    }

    public String getContent() {
        return content;
    }

    public void setContent(String content) {
        this.content = content;
    }

    public long getTimestamp() {
        return timestamp;
    }

    public void setTimestamp(long timestamp) {
        this.timestamp = timestamp;
    }

    public boolean isLock() {
        return isLock;
    }

    public void setLock(boolean lock) {
        isLock = lock;
    }

    public boolean isReadOnly() {
        return readOnly;
    }

    public void setReadOnly(boolean readOnly) {
        this.readOnly = readOnly;
    }

    public String getLeaveUser() {
        return leaveUser;
    }

    public void setLeaveUser(String leaveUser) {
        this.leaveUser = leaveUser;
    }
}
