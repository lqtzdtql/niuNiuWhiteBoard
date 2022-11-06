package com.liuyue.painter.model;

import com.google.gson.Gson;

import java.util.List;

public class HallBean {
    /**
     * code : 200
     * message : get roomlist success
     * roomlist : [{"uuid":"01GH4CACYSK4DC9DJQEKS8E07J","name":"teatRoom","host_id":1,"type":"teaching_room"},{"uuid":"01GH4EV2T8A3RGW2NZ7S840SGE","name":"teatRoom","host_id":2,"type":"teaching_room"}]
     */

    private int code;
    private String message;
    private List<RoomlistBean> roomList;

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
        return roomList;
    }

    public void setRoomList(List<RoomlistBean> roomList) {
        this.roomList = roomList;
    }

    public static class RoomlistBean {
        /**
         * uuid : 01GH4CACYSK4DC9DJQEKS8E07J
         * name : teatRoom
         * host_id : 1
         * type : teaching_room
         */

        private String uuid;
        private String name;
        private int hostId;
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

        public int getHostId() {
            return hostId;
        }

        public void setHostId(int hostId) {
            this.hostId = hostId;
        }

        public String getType() {
            return type;
        }

        public void setType(String type) {
            this.type = type;
        }
    }
}
