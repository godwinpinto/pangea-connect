package org.springframework.boot.pangeareactive.bean;

public class HelloResponseBean {
    private String message;

    public HelloResponseBean() {
    }

    public HelloResponseBean(String message) {
        this.message = message;
    }

    public String getMessage() {
        return message;
    }

    public void setMessage(String message) {
        this.message = message;
    }
}
