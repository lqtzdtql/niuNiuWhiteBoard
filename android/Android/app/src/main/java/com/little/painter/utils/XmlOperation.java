package com.little.painter.utils;

import android.graphics.Color;
import android.graphics.Paint;

import com.little.painter.Constants;
import com.little.painter.shape.Circle;
import com.little.painter.shape.Ink;
import com.little.painter.shape.Line;
import com.little.painter.shape.Point;
import com.little.painter.shape.Rectangle;
import com.little.painter.shape.Shape;

import org.dom4j.Attribute;
import org.dom4j.Document;
import org.dom4j.DocumentException;
import org.dom4j.DocumentHelper;
import org.dom4j.Element;
import org.dom4j.io.OutputFormat;
import org.dom4j.io.SAXReader;
import org.dom4j.io.XMLWriter;

import java.io.File;
import java.io.FileWriter;
import java.io.IOException;
import java.util.ArrayList;
import java.util.Iterator;
import java.util.List;

public class XmlOperation {

    /**
     * 将画图保存成XML格式文件
     *
     * @param shapeList
     */
    public static int CreatXml(List<Shape> shapeList, String filename) {

        int returnValue = 0; // 创建结果返回，0失败，1成功

        Document document = DocumentHelper.createDocument();// 创建Document对象
        Element root = document.addElement("Painting");// 创建根节点

        for (int i = 0; i < shapeList.size(); i++) {

            Element shape = root.addElement("PaintShape");// 创建Shape对象对应节点

            Element kind = shape.addElement("Kind");// shape对象的kind节点
            kind.setText(String.valueOf(shapeList.get(i).getKind()));

            Element color = shape.addElement("Color");// Shape对象的color节点
            color.setText(shapeList.get(i).getColor());

            Element width = shape.addElement("Width");// Shape对象的width节点
            width.setText(String.valueOf(shapeList.get(i).getWidth()));

            Element PointList = shape.addElement("Pointlist");// 创建Shape对象的PointList节点
            PointList.setText(changeListDateToStr(shapeList.get(i).getPointList()));
        }
        try {
            OutputFormat format = OutputFormat.createPrettyPrint();// 格式化输出
            format.setEncoding("UTF-8");// 制定xml编码
            XMLWriter writer = new XMLWriter(new FileWriter(filename), format);
            writer.write(document);
            writer.close();
            returnValue = 1;
        } catch (IOException e) {
            e.printStackTrace();
        }
        return returnValue;
    }


    /**
     * 将xml文件解析成对应的list<shape>
     * 采用DOM4j解析机制
     *
     * @param filename
     * @return
     */
    public static List<Shape> TransXmlToShape(String filename) throws DocumentException {

        List<Shape> lists = new ArrayList<>();
        // 得到Document对象
        SAXReader reader = new SAXReader();
        Document document = reader.read(new File(filename));
        // 获取文档的根节点.
        Element root = document.getRootElement();
        Iterator node = root.elementIterator();
        while (node.hasNext()) {
            Element shapeElement = (Element) node.next();
            Shape shape = null;
            switch (Integer.valueOf(shapeElement.elementText("Kind"))) {
                case Constants.INK:
                    shape = new Ink();
                    break;
                case Constants.LINE:
                    shape = new Line();
                    break;
                case Constants.RECT:
                    shape = new Rectangle();
                    break;
                case Constants.CIRCLE:
                    shape = new Circle();
                    break;
            }
            // 设置公共属性
            shape.setColor(shapeElement.elementText("Color"));
            shape.setWidth(Float.valueOf(shapeElement.elementText("Width")));
            shape.setPointList(ChangeStrToListData(shapeElement.elementText("Pointlist")));
            Paint newpaint = new Paint();
            newpaint.setStyle(Paint.Style.STROKE);
            newpaint.setStrokeWidth(shape.getWidth());
            newpaint.setColor(Color.parseColor(shape.getColor()));
            shape.setPaint(newpaint);
            // 设置各自专属属性
            shape.setOwnProperty();
            // 加到List当中
            lists.add(shape);
        }
        return lists;
    }

    /**
     * 遍历xml文档打印相关属性信息
     *
     * @param root 还未使用
     */
    private static void listNodes(Element root) {
        // 首先获取当前节点的所有属性节点
        List<Attribute> list = root.attributes();
        // 同时迭代当前节点下面的所有子节点
        // 使用递归
        Iterator<Element> iterator = root.elementIterator();
        while (iterator.hasNext()) {
            Element e = iterator.next();
            listNodes(e);
        }
    }


    /**
     * 将String数据转化成笔迹点集合
     */
    private static List<Point> ChangeStrToListData(String text) {
        List<Point> lists = new ArrayList<>();
        String[] strarray = text.split(";");
        for (String str : strarray) {
            Point newpoint = new Point();

            String str1 = str.substring(0, str.indexOf(","));
            newpoint.setX(Float.parseFloat(str1));

            String str2 = str.substring(str.indexOf(",") + 1);
            newpoint.setY(Float.parseFloat(str2));

            lists.add(newpoint);
        }
        return lists;
    }

    /**
     * 将笔迹点数据转成String格式
     */
    private static String changeListDateToStr(List<Point> pointList) {

        StringBuffer sb = new StringBuffer();
        for (int j = 0; j < pointList.size(); j++) {
            sb.append(pointList.get(j).getX() + "," + pointList.get(j).getY() + ";");
        }
        return sb.toString();
    }


}
