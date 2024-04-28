package org.springframework.boot.pangeareactive.controller;

import org.springframework.boot.pangeareactive.bean.HelloResponseBean;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;
import reactor.core.publisher.Mono;

@RestController
public class HelloWorldController {

    @GetMapping(value = "/get")
    public Mono<HelloResponseBean> hello() {
        HelloResponseBean helloResponseBean = new HelloResponseBean("Hello, World!");
        return Mono.just(helloResponseBean);    //prints Hello, World!
    }
}
