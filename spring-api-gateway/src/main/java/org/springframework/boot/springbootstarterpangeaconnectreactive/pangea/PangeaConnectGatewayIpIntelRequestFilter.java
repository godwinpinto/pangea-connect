package org.springframework.boot.springbootstarterpangeaconnectreactive.pangea;

import cloud.pangeacyber.pangea.Config;
import cloud.pangeacyber.pangea.intel.IPIntelClient;
import cloud.pangeacyber.pangea.intel.requests.IPReputationRequest;
import cloud.pangeacyber.pangea.intel.responses.IPReputationResponse;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.cloud.gateway.filter.GatewayFilterChain;
import org.springframework.cloud.gateway.filter.GlobalFilter;
import org.springframework.core.Ordered;
import org.springframework.data.redis.core.ReactiveRedisTemplate;
import org.springframework.http.HttpStatus;
import org.springframework.http.server.reactive.ServerHttpRequest;
import org.springframework.stereotype.Component;
import org.springframework.web.server.ResponseStatusException;
import org.springframework.web.server.ServerWebExchange;
import reactor.core.publisher.Mono;

import java.util.Objects;
import java.util.concurrent.CompletableFuture;

@Component
public class PangeaConnectGatewayIpIntelRequestFilter implements GlobalFilter, Ordered {

    private final ReactiveRedisTemplate<String, String> redisTemplate;

    public PangeaConnectGatewayIpIntelRequestFilter(ReactiveRedisTemplate<String, String> redisTemplate) {
        this.redisTemplate = redisTemplate;
    }

    @Value("${pangea.token}")
    private String token;

    @Value("${pangea.domain}")
    private String domain;

    @Value("${pangea.ip-intel.enabled}")
    private boolean ipIntelEnabled;

    @Value("${pangea.ip-intel.score-threshold}")
    private int ipIntelThreshold;

    @Value("${pangea.ip-intel.type}")
    private String ipIntelType;

    @Override
    public Mono<Void> filter(ServerWebExchange exchange, GatewayFilterChain chain) {
        if(!ipIntelEnabled || ipIntelType==null || ipIntelType.isEmpty()){
            return chain.filter(exchange);
        }
        ServerHttpRequest request = exchange.getRequest();
        String ipAddress;
        if(ipIntelType!=null && ipIntelType.equals("header")){
            ipAddress= request.getHeaders().getFirst("X-Forwarded-For");
        }else if(ipIntelType!=null && ipIntelType.equals("request")){
            ipAddress = Objects.requireNonNull(request.getRemoteAddress()).getAddress().getHostAddress();
        } else {
            ipAddress = "";
        }
        System.out.println("ipAddress"+ipAddress);

        return redisTemplate.opsForValue().get(ipAddress)
                .flatMap(isBlocked ->
                        allowOrDenyRequest(exchange, chain, isBlocked))
                .switchIfEmpty(Mono.fromRunnable(()->{
                    processIpIntel(ipAddress);
                }))
                .then(chain.filter(exchange));
    }

    private void processIpIntel(String ipAddress) {
        CompletableFuture.runAsync(() -> {

            Config pangeaConfig = new Config.Builder(token, domain).build();
            IPIntelClient client = new IPIntelClient.Builder(pangeaConfig).build();
            IPReputationResponse response = null;
            try {
                response = client.reputation(
                        new IPReputationRequest.Builder(ipAddress).provider("crowdstrike").verbose(true).raw(true).build()
                );
                String blockedStatus;
                if(response.getStatus().toLowerCase().equals("success")){
                    if(response.getResult().getData().getScore()>=ipIntelThreshold){
                        blockedStatus="Y";
                    }else{
                        blockedStatus="N";
                    }
                    redisTemplate.opsForValue().set(ipAddress, blockedStatus).subscribe();
                }else{
                    redisTemplate.opsForValue().set(ipAddress, "Y").subscribe();
                }

            } catch (Exception ignored){

            }
        });
    }

    private static Mono<Void> allowOrDenyRequest(ServerWebExchange exchange, GatewayFilterChain chain, String isBlocked) {
        if (isBlocked != null && isBlocked.equals("Y")) {
            return Mono.error(new ResponseStatusException(HttpStatus.FORBIDDEN, "Forbidden"));
        }else{
            return chain.filter(exchange);
        }
    }

    @Override
    public int getOrder() {
        return -2;
    }
}
