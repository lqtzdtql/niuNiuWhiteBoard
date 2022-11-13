package com.little.painter.model;

import com.google.gson.Gson;

public class WhiteBoardAuthBean {

    /**
     * message : enter room success
     * user_uuid : participantuuid
     * user_name : username
     * room_uuid : roomuuid
     * code : 200
     */

    private String message;
    private String user_uuid;
    private String user_name;
    private String room_uuid;
    private int code;

    public static WhiteBoardAuthBean parse(String str) {
        return new Gson().fromJson(str, WhiteBoardAuthBean.class);
    }

    public String getMessage() {
        return message;
    }

    public void setMessage(String message) {
        this.message = message;
    }

    public String getUserUUID() {
        return user_uuid;
    }

    public void setUserUUID(String userUUID) {
        this.user_uuid = userUUID;
    }

    public String getUserName() {
        return user_name;
    }

    public void setUserName(String userName) {
        this.user_name = userName;
    }

    public String getRoomUUID() {
        return room_uuid;
    }

    public void setRoomUUID(String roomUUID) {
        this.room_uuid = roomUUID;
    }

    public int getCode() {
        return code;
    }

    public void setCode(int code) {
        this.code = code;
    }
}
