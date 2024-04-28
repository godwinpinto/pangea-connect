package org.springframework.boot.springbootpangea.filter;


import cloud.pangeacyber.pangea.Config;
import cloud.pangeacyber.pangea.intel.IPIntelClient;
import cloud.pangeacyber.pangea.intel.requests.IPReputationRequest;
import cloud.pangeacyber.pangea.intel.responses.IPReputationResponse;
import jakarta.servlet.http.HttpServletResponse;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.core.annotation.Order;
import org.springframework.data.redis.core.RedisTemplate;
import org.springframework.stereotype.Component;

import jakarta.servlet.Filter;
import jakarta.servlet.FilterChain;
import jakarta.servlet.FilterConfig;
import jakarta.servlet.ServletException;
import jakarta.servlet.ServletRequest;
import jakarta.servlet.ServletResponse;
import jakarta.servlet.http.HttpServletRequest;
import org.springframework.web.filter.OncePerRequestFilter;

import java.io.IOException;
import java.util.Objects;
import java.util.concurrent.CompletableFuture;

/**
 * A filter to create transaction before and commit it once request completes
 * The current implemenatation is just for demo
 *
 */
@Component
@Order(1)
public class PangeaConnectIpIntelpFilter extends OncePerRequestFilter {

    private final RedisTemplate<String, String> redisTemplate;

    public PangeaConnectIpIntelpFilter(RedisTemplate<String, String> redisTemplate) {
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
    protected void doFilterInternal(HttpServletRequest request, HttpServletResponse response, FilterChain filterChain) throws ServletException, IOException {

        if (!ipIntelEnabled || ipIntelType == null || ipIntelType.isEmpty()) {
            filterChain.doFilter(request, response);
            return;
        }

        String ipAddress;
        if (ipIntelType.equals("header")) {
            ipAddress = request.getHeader("X-Forwarded-For");
        } else if (ipIntelType.equals("request")) {
            ipAddress = Objects.requireNonNull(request.getRemoteAddr());
        } else {
            ipAddress = "";
        }

        String isBlocked = redisTemplate.opsForValue().get(ipAddress);
        allowOrDenyRequest(request, response, filterChain, isBlocked,ipAddress);
    }

    private void processIpIntel(String ipAddress) {
        Config pangeaConfig = new Config.Builder(token, domain).build();
        IPIntelClient client = new IPIntelClient.Builder(pangeaConfig).build();
        IPReputationResponse response;
        try {
            response = client.reputation(
                    new IPReputationRequest.Builder(ipAddress).provider("crowdstrike").verbose(true).raw(true).build()
            );
            String blockedStatus;
            if (response.getStatus().equalsIgnoreCase("success")) {
                blockedStatus = response.getResult().getData().getScore() >= ipIntelThreshold ? "Y" : "N";
                redisTemplate.opsForValue().set(ipAddress, blockedStatus);
            } else {
                redisTemplate.opsForValue().set(ipAddress, "Y");
            }
        } catch (Exception ignored) {
        }
    }

    private void allowOrDenyRequest(HttpServletRequest request, HttpServletResponse response, FilterChain filterChain, String isBlocked,String ipAddress) throws IOException, ServletException {
        if (isBlocked != null && isBlocked.equals("Y")) {
            response.setStatus(HttpServletResponse.SC_FORBIDDEN);
        } else if(isBlocked==null){
            CompletableFuture.runAsync(() -> processIpIntel(ipAddress));
            filterChain.doFilter(request, response);
        } else{
            filterChain.doFilter(request, response);
        }
    }

}
