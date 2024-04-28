package org.springframework.boot.springbootpangea.controller;


import org.springframework.boot.springbootpangea.bean.HelloResponseBean;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class HelloWorldController {

    @GetMapping(value = "/get")
    public HelloResponseBean hello() {
        return new HelloResponseBean("Hello, World!");    //prints Hello, World!
    }
}
