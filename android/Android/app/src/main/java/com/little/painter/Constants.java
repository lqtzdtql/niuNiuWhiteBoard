package com.little.painter;

public class Constants {
    public static final String SP_USER_INFO = "USER_INFO";

    public static final String KEY_USER_PHONE = "USER_PHONE";
    public static final String KEY_USER_PASSWORD = "USER_PASSWORD";

    /**
     * ======================================== 绘图相关 ============================================
     */

    /**
     * 画笔最大
     */
    public static final int maxBrushSize = 20;
    /**
     * 画笔最小
     */
    public static final int minBrushSize = 5;
    /**
     * 画笔预设颜色
     */
    public static final String[] colors = new String[]{
            "#242424",
            "#FF0000",
            "#9ACD32",
            "#473C8B",
            "#EEEE00",
            "#EE8262",
            "#EE3A8C",
            "#836FFF",
            "#CDCDB4",
            "#FF7F24"
    };

    /**
     * 曲线笔迹
     */
    public static final int INK = 1;
    /**
     * 直线
     */
    public static final int LINE = 2;
    /**
     * 矩阵
     */
    public static final int RECT = 3;
    /**
     * 圆
     */
    public static final int CIRCLE = 4;

    /**
     * 退出信号
     */
    public static final int MSG_EXIT = 0X111;

    /**
     * 保存格式
     */
    public static final int PNG = 5;
    public static final int SVG = 6;

    /**
     * 重绘信号
     */
    public static final int MSG_REDRAW = 0X112;
}
