package com.little.painter.model;

import com.google.gson.Gson;

import java.util.List;

public class HallBean {

    /**
     * code : 200
     * message : 获取房间列表成功
     * roomlist : [{"uuid":"01GHKVWDRSMRJYGK3JK2VD0Z6Q","name":"213213123","host_name":"1","type":"teaching_room"}]
     */

    private int code;
    private String message;
    private List<RoomlistBean> roomlist;

    public static HallBean parse(String str) {
        return new Gson().fromJson(str, HallBean.class);
    }

    public int getCode() {
        return code;
    }

    public void setCode(int code) {
        this.code = code;
    }

    public String getMessage() {
        return message;
    }

    public void setMessage(String message) {
        this.message = message;
    }

    public List<RoomlistBean> getRoomList() {
        return roomlist;
    }

    public void setRoomlist(List<RoomlistBean> roomlist) {
        this.roomlist = roomlist;
    }

    public static class RoomlistBean {
        /**
         * uuid : 01GHKVWDRSMRJYGK3JK2VD0Z6Q
         * name : 213213123
         * host_name : 1
         * type : teaching_room
         */

        private String uuid;
        private String name;
        private String host_name;
        private String type;

        public static RoomlistBean parse(String str) {
            return new Gson().fromJson(str, RoomlistBean.class);
        }

        public String getUuid() {
            return uuid;
        }

        public void setUuid(String uuid) {
            this.uuid = uuid;
        }

        public String getName() {
            return name;
        }

        public void setName(String name) {
            this.name = name;
        }

        public String getHostName() {
            return host_name;
        }

        public void setHostName(String hostName) {
            this.host_name = hostName;
        }

        public String getType() {
            return type;
        }

        public void setType(String type) {
            this.type = type;
        }
    }
}
