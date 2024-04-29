package cloud.pangea;

import cloud.pangea.bean.HelloResponseBean;
import jakarta.ws.rs.GET;
import jakarta.ws.rs.Path;
import jakarta.ws.rs.Produces;
import io.quarkus.security.UnauthorizedException;
import io.smallrye.mutiny.Uni;
import jakarta.ws.rs.*;
import jakarta.ws.rs.core.MediaType;
import jakarta.ws.rs.core.Response;

@Path("/get")
public class GreetingResource {

    @GET
    @Produces(MediaType.APPLICATION_JSON)
    @Consumes(MediaType.APPLICATION_JSON)
    public Uni<Response> get()   {
            return Uni.createFrom().item(Response.status(Response.Status.BAD_REQUEST)
                    .entity(new HelloResponseBean("Hello World")).build());
    }
}
