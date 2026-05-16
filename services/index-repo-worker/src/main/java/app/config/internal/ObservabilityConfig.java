package app.config.internal;

import io.micrometer.observation.ObservationPredicate;
import io.micrometer.observation.ObservationRegistry;
import io.micrometer.observation.aop.ObservedAspect;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.http.server.observation.ServerRequestObservationContext;

@Configuration
public class ObservabilityConfig {

    @Bean
    ObservedAspect observedAspect(ObservationRegistry observationRegistry) {
        return new ObservedAspect(observationRegistry);
    }

    @Bean
    ObservationPredicate noScheduledTaskObservations() {
        return (name, context) -> !name.equals("tasks.scheduled.execution");
    }

    @Bean
    ObservationPredicate noActuatorObservations() {
        return (name, context) -> !(context instanceof ServerRequestObservationContext c
            && c.getCarrier().getRequestURI().startsWith("/actuator"));
    }
}
