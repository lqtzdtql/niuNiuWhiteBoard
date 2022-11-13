package com.little.painter.callback;

public interface PageChangeCall {
    void onPageAddCall(int pagenum, int pageindex);
    void onPagePreCall(int pageindex);
    void onPageNextCall(int pageindex);
}
