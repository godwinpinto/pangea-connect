package cloud.pangea.filter;

import cloud.pangeacyber.pangea.Config;
import cloud.pangeacyber.pangea.intel.IPIntelClient;
import cloud.pangeacyber.pangea.intel.requests.IPReputationRequest;
import cloud.pangeacyber.pangea.intel.responses.IPReputationResponse;
import io.quarkus.redis.client.RedisClientName;
import io.quarkus.redis.datasource.ReactiveRedisDataSource;
import io.quarkus.redis.datasource.value.ReactiveValueCommands;
import io.smallrye.mutiny.Uni;
import io.smallrye.mutiny.infrastructure.Infrastructure;
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;
import jakarta.ws.rs.container.ContainerRequestContext;
import jakarta.ws.rs.core.Response;
import org.eclipse.microprofile.config.inject.ConfigProperty;
import org.jboss.resteasy.reactive.server.ServerRequestFilter;

public class PangeaFilter {


    @ConfigProperty(name = "pangea.token")
    private String token;

    @ConfigProperty(name = "pangea.domain")
    private String domain;

    @ConfigProperty(name = "pangea.ip-intel.enabled")
    private boolean ipIntelEnabled;

    @ConfigProperty(name = "pangea.ip-intel.score-threshold")
    private int ipIntelThreshold;

    @ConfigProperty(name = "pangea.ip-intel.type")
    private String ipIntelType;

    private final ReactiveValueCommands<String,String> valCommands;
    public PangeaFilter(@RedisClientName("pangea-cache")
                            ReactiveRedisDataSource ds) {
        valCommands=ds.value(String.class, String.class);
    }

    Uni<String> getFromCache(String key) {
        return valCommands.get(key);
    }

    Uni<Void> setToCache(String key, String value) {
        return valCommands.setex(key, 15*60 , value);
    }

    @ServerRequestFilter(preMatching = true)
    public Uni<Response> ipIntelFilter(ContainerRequestContext requestContext) {
        if (!ipIntelEnabled || ipIntelType == null || ipIntelType.isEmpty()) {
            return Uni.createFrom().nullItem();
        }

        String ipAddress;
        if (ipIntelType.equals("header")) {
            ipAddress = requestContext.getHeaderString("X-Forwarded-For");
        } else if (ipIntelType.equals("request")) {
            ipAddress = requestContext.getUriInfo().getRequestUri().getHost();
        } else {
            ipAddress = "";
        }

        return getFromCache(ipAddress)
                .onItem().transformToUni(isBlocked -> allowOrDenyRequest(requestContext, isBlocked, ipAddress));

    }

    private Uni<Void> processIpIntel(String ipAddress) {
        Config pangeaConfig = new Config.Builder(token, domain).build();
        IPIntelClient client = new IPIntelClient.Builder(pangeaConfig).build();

        return Uni.createFrom().voidItem().call(() -> {
            try {
                IPReputationResponse response = client.reputation(
                        new IPReputationRequest.Builder(ipAddress).provider("crowdstrike").verbose(true).raw(true).build()
                );
                if ("success".equalsIgnoreCase(response.getStatus())) {
                    return setToCache(ipAddress, response.getResult().getData().getScore() >= ipIntelThreshold ? "Y" : "N");
                } else {
                    return setToCache(ipAddress, "N");
                }
            } catch (Exception ignored) {
            }
            return Uni.createFrom().voidItem();
        });
    }


    private Uni<Response> allowOrDenyRequest(ContainerRequestContext requestContext, String isBlocked, String ipAddress) {
        if (isBlocked == null) {
            processIpIntel(ipAddress).runSubscriptionOn(Infrastructure.getDefaultWorkerPool())
                    .subscribeAsCompletionStage().isDone();
            return Uni.createFrom().nullItem();
        }else if ("Y".equals(isBlocked)) {
            return Uni.createFrom().item(Response.status(Response.Status.FORBIDDEN).build());
        } else {
            return Uni.createFrom().nullItem();
        }
    }
}
