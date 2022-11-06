package com.liuyue.painter.model;

import com.google.gson.Gson;

import java.util.List;

public class RoomBean {
    /**
     * uuid : 01GH4GZW1E9TQXTGK3HJ8W7JQC
     * name : teatRoom
     * host_id : 1
     * created_time : 2022-11-06T02:31:32+08:00
     * updated_time : 2022-11-06T02:31:32+08:00
     * type : teaching_room
     * participants : [{"name":"user1","user_uuid":"01GH4C5DCD2C10B56499RJP758","permission":"host"}]
     */

    private String uuid;
    private String name;
    private int hostId;
    private String createdTime;
    private String updatedTime;
    private String type;
    private List<ParticipantsBean> participants;

    public static RoomBean parse(String str) {
        return new Gson().fromJson(str, RoomBean.class);
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

    public String getCreatedTime() {
        return createdTime;
    }

    public void setCreatedTime(String createdTime) {
        this.createdTime = createdTime;
    }

    public String getUpdatedTime() {
        return updatedTime;
    }

    public void setUpdatedTime(String updatedTime) {
        this.updatedTime = updatedTime;
    }

    public String getType() {
        return type;
    }

    public void setType(String type) {
        this.type = type;
    }

    public List<ParticipantsBean> getParticipants() {
        return participants;
    }

    public void setParticipants(List<ParticipantsBean> participants) {
        this.participants = participants;
    }

    public static class ParticipantsBean {
        /**
         * name : user1
         * user_uuid : 01GH4C5DCD2C10B56499RJP758
         * permission : host
         */

        private String name;
        private String userUUID;
        private String permission;

        public static ParticipantsBean parse(String str) {
            return new Gson().fromJson(str, ParticipantsBean.class);
        }

        public String getName() {
            return name;
        }

        public void setName(String name) {
            this.name = name;
        }

        public String getUserUUID() {
            return userUUID;
        }

        public void setUserUUID(String userUUID) {
            this.userUUID = userUUID;
        }

        public String getPermission() {
            return permission;
        }

        public void setPermission(String permission) {
            this.permission = permission;
        }
    }
}
