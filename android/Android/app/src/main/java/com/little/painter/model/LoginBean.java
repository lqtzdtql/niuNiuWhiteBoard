package com.little.painter.model;

import com.google.gson.Gson;

import java.io.Serializable;

public class LoginBean {

    /**
     * code : 200
     * message : login success
     * user_info : {"id":3,"uuid":"01GH10GE835DW2W4X1A9S8DVCS","name":"userName","mobile":"13344443333"}
     */

    private int code;
    private String message;
    private UserInfoBean userInfo;

    public static LoginBean parse(String str) {
        return new Gson().fromJson(str, LoginBean.class);
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

    public UserInfoBean getUserInfo() {
        return userInfo;
    }

    public void setUserInfo(UserInfoBean userInfo) {
        this.userInfo = userInfo;
    }

    public static class UserInfoBean implements Serializable {
        /**
         * id : 3
         * uuid : 01GH10GE835DW2W4X1A9S8DVCS
         * name : userName
         * mobile : 13344443333
         */

        private int id;
        private String uuid;
        private String name;
        private String mobile;

        public static UserInfoBean parse(String str) {
            return new Gson().fromJson(str, UserInfoBean.class);
        }

        public int getId() {
            return id;
        }

        public void setId(int id) {
            this.id = id;
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

        public String getMobile() {
            return mobile;
        }

        public void setMobile(String mobile) {
            this.mobile = mobile;
        }
    }
}
