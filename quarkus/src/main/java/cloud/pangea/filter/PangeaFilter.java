package cloud.pangea.filter;

import cloud.pangeacyber.pangea.Config;
import cloud.pangeacyber.pangea.intel.IPIntelClient;
import cloud.pangeacyber.pangea.intel.requests.IPReputationRequest;
import cloud.pangeacyber.pangea.intel.responses.IPReputationResponse;
import io.quarkus.redis.datasource.ReactiveRedisDataSource;
import io.quarkus.redis.datasource.value.ReactiveValueCommands;
import io.smallrye.mutiny.Uni;
import io.smallrye.mutiny.infrastructure.Infrastructure;
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.ws.rs.container.ContainerRequestContext;
import jakarta.ws.rs.core.Response;
import org.eclipse.microprofile.config.inject.ConfigProperty;
import org.jboss.resteasy.reactive.server.ServerRequestFilter;

@ApplicationScoped
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

    private ReactiveValueCommands<String,String> valCommands;
//    private ValueCommands<String, String> cmd;

    public PangeaFilter(ReactiveRedisDataSource ds) {
        valCommands=ds.value(String.class, String.class);
    }

    Uni<String> getFromCache(String key) {
        return valCommands.get(key);
    }

    void setToCache(String key, String value) {
        System.out.println("Setting cache for key: "+key+" value: "+value);
        valCommands.setex(key, 15 , value);
    }

    @ServerRequestFilter(preMatching = true)
    public Uni<Void> ipIntelFilter(ContainerRequestContext requestContext) {
        if (!ipIntelEnabled || ipIntelType == null || ipIntelType.isEmpty()) {
            return Uni.createFrom().voidItem();
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
                System.out.println(response.getStatus()+"::"+response.getResult().getRawData().toString());
                if ("success".equalsIgnoreCase(response.getStatus())) {
                    setToCache(ipAddress, response.getResult().getData().getScore() >= ipIntelThreshold ? "Y" : "N");
                } else {
                    setToCache(ipAddress, "N");
                }
            } catch (Exception ignored) {
            }
                return Uni.createFrom().voidItem();
        });
    }


    private Uni<Void> allowOrDenyRequest(ContainerRequestContext requestContext, String isBlocked, String ipAddress) {
        System.out.println("isBlocked"+isBlocked);
        if ("Y".equals(isBlocked)) {
            requestContext.abortWith(Response.status(Response.Status.FORBIDDEN).build());
            return Uni.createFrom().voidItem();
        } else if (isBlocked == null) {
            processIpIntel(ipAddress).runSubscriptionOn(Infrastructure.getDefaultWorkerPool())
                    .subscribeAsCompletionStage().whenComplete((i, e) -> {
                        System.out.println("IP Intel completed");
                        if(e != null) {
                            e.printStackTrace();
                        }
                    });;
            return Uni.createFrom().voidItem();
        } else {
            return Uni.createFrom().nullItem();
        }
    }
}
