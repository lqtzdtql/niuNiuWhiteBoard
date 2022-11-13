package com.liuyue.painter.utils;

import android.util.Log;

import com.google.gson.JsonObject;
import com.liuyue.painter.model.CreateRoomBean;
import com.liuyue.painter.model.HallBean;
import com.liuyue.painter.model.LoginBean;
import com.liuyue.painter.model.RoomBean;
import com.liuyue.painter.model.SignupBean;
import com.liuyue.painter.model.WhiteBoardAuthBean;

import org.json.JSONObject;

import java.io.IOException;

import okhttp3.MediaType;
import okhttp3.OkHttpClient;
import okhttp3.Request;
import okhttp3.RequestBody;
import okhttp3.Response;

public class AppServer {
    private static final String SERVER_ADDRESS = "http://81.68.68.216:8282";
    private static final String MEDIA_TYPE_JSON = "application/json";

    private final OkHttpClient mOkHttpClient;
    private String mToken;

    private static class AppServerHolder {
        private static final AppServer INSTANCE = new AppServer();
    }

    private AppServer() {
        mOkHttpClient = new OkHttpClient();
    }

    public static AppServer getInstance() {
        return AppServerHolder.INSTANCE;
    }

    public LoginBean login(String phone, String password) {
        String url = SERVER_ADDRESS + "/login";
        JsonObject jsonObject = new JsonObject();
        jsonObject.addProperty("mobile", phone);
        jsonObject.addProperty("passwd", password);
        RequestBody requestBody = RequestBody.create(jsonObject.toString(), MediaType.parse(MEDIA_TYPE_JSON));
        Request request = new Request.Builder().url(url).method("POST", requestBody).build();
        try {
            Response response = mOkHttpClient.newCall(request).execute();
            if (response.code() != 200) {
                return null;
            }
            mToken = response.header("Refresh-token");
            return LoginBean.parse(response.body().string());
        } catch (IOException e) {
            e.printStackTrace();
            return null;
        }
    }

    public SignupBean signup(String phone, String password, String nickName) {
        String url = SERVER_ADDRESS + "/signup";
        JsonObject jsonObject = new JsonObject();
        jsonObject.addProperty("mobile", phone);
        jsonObject.addProperty("passwd", password);
        jsonObject.addProperty("name", nickName);
        RequestBody requestBody = RequestBody.create(jsonObject.toString(), MediaType.parse(MEDIA_TYPE_JSON));
        Request request = new Request.Builder().url(url).method("POST", requestBody).build();
        try {
            Response response = mOkHttpClient.newCall(request).execute();
            if (response.code() != 200) {
                return null;
            }
            return SignupBean.parse(response.body().string());
        } catch (IOException e) {
            e.printStackTrace();
            return null;
        }
    }

    public LoginBean.UserInfoBean getUserInfo(String userUUID) {
        String url = SERVER_ADDRESS + "/v1/userinfo/" + userUUID;
        Request request = new Request.Builder()
                .url(url)
                .method("GET", null)
                .addHeader("Access-Token", mToken)
                .build();
        try {
            Response response = mOkHttpClient.newCall(request).execute();
            if (response.code() != 200) {
                return null;
            }
            return LoginBean.UserInfoBean.parse(response.body().string());
        } catch (IOException e) {
            return null;
        }
    }

    public HallBean getRoomList() {
        String url = SERVER_ADDRESS + "/v1/roomlist";
        Request request = new Request.Builder()
                .url(url)
                .method("GET", null)
                .addHeader("Access-Token", mToken)
                .build();
        try {
            Response response = mOkHttpClient.newCall(request).execute();
            if (response.code() != 200) {
                return null;
            }
            return HallBean.parse(response.body().string());
        } catch (IOException e) {
            return null;
        }
    }

    public CreateRoomBean createRoom(String roomName, String roomType) {
        String url = SERVER_ADDRESS + "/v1/rooms";
        JsonObject jsonObject = new JsonObject();
        jsonObject.addProperty("name", roomName);
        jsonObject.addProperty("type", roomType);
        RequestBody requestBody = RequestBody.create(jsonObject.toString(), MediaType.parse(MEDIA_TYPE_JSON));
        Request request = new Request.Builder()
                .url(url)
                .method("POST", requestBody)
                .addHeader("Access-Token", mToken)
                .build();
        try {
            Response response = mOkHttpClient.newCall(request).execute();
            if (response.code() != 200) {
                Log.e("飞", "createRoom: " + response.body().string());
                return null;
            }
            return CreateRoomBean.parse(response.body().string());
        } catch (IOException e) {
            return null;
        }
    }

    public String enterRoom(String userUUID) {
        String url = SERVER_ADDRESS + "/v1/rooms/" + userUUID + "/whiteboard";
        Request request = new Request.Builder()
                .url(url)
                .method("GET", null)
                .addHeader("ASSESS-Token", mToken)
                .build();
        try {
            Response response = mOkHttpClient.newCall(request).execute();
            if (response.code() != 200) {
                return null;
            }
            JSONObject jsonObject = new JSONObject(response.body().string());
            return jsonObject.getString("token");
        } catch (Exception e) {
            e.printStackTrace();
            return null;
        }
    }

    public WhiteBoardAuthBean authWhiteBoard(String roomToken) {
        String url = SERVER_ADDRESS + "/auth?token=" + roomToken;
        Request request = new Request.Builder()
                .url(url)
                .method("GET", null)
                .build();
        try {
            Response response = mOkHttpClient.newCall(request).execute();
            if (response.code() != 200) {
                return null;
            }
            return WhiteBoardAuthBean.parse(response.body().string());
        } catch (Exception e) {
            e.printStackTrace();
            return null;
        }
    }

    public RoomBean getRoomInfo(String roomUUID) {
        String url = SERVER_ADDRESS + "/v1/rooms/" + roomUUID;
        Request request = new Request.Builder()
                .url(url)
                .method("GET", null)
                .addHeader("Access-Token", mToken)
                .build();
        try {
            Response response = mOkHttpClient.newCall(request).execute();
            if (response.code() != 200) {
                return null;
            }
            return RoomBean.parse(response.body().string());
        } catch (IOException e) {
            e.printStackTrace();
            return null;
        }
    }

    public boolean exitRoom(String roomUUID) {
        String url = SERVER_ADDRESS + "/v1/rooms/" + roomUUID + "/exit";
        Request request = new Request.Builder()
                .url(url)
                .method("GET", null)
                .addHeader("Access-Token", mToken)
                .build();
        try {
            Response response = mOkHttpClient.newCall(request).execute();
            return response.code() == 200;
        } catch (IOException e) {
            return false;
        }
    }

    public boolean logout() {
        String url = SERVER_ADDRESS + "/v1/logout";
        Request request = new Request.Builder()
                .url(url)
                .method("GET", null)
                .addHeader("Access-Token", mToken)
                .build();
        try {
            Response response = mOkHttpClient.newCall(request).execute();
            return response.code() == 200;
        } catch (IOException e) {
            return false;
        }
    }

}
