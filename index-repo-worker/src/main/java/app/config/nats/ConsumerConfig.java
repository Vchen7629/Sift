package app.config.nats;

import java.io.IOException;
import java.time.Duration;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import io.nats.client.JetStream;
import io.nats.client.JetStreamApiException;
import io.nats.client.JetStreamManagement;
import io.nats.client.JetStreamSubscription;
import io.nats.client.PullSubscribeOptions;
import io.nats.client.api.AckPolicy;
import io.nats.client.api.ConsumerConfiguration;
import io.nats.client.api.ConsumerInfo;
import io.nats.client.api.StreamInfo;

@Configuration
public class ConsumerConfig {
    private final static String CONSUMER_NAME = "jobProcessConsumer";
    private final static String STREAM_NAME = StreamConfig.STREAM_NAME;
    private final static String SUBJECT_NAME = "index-repo.subject.request";
    private final static int RETRY_COUNT_LIMIT = 5;
    private final static Duration ACK_WAIT = Duration.ofSeconds(30);
    
    @Bean
    public ConsumerInfo indexRepoConsumer(JetStreamManagement jsm, StreamInfo streamInfo) throws IOException, JetStreamApiException {
        ConsumerConfiguration config = ConsumerConfiguration.builder()
            .durable(CONSUMER_NAME)
            .filterSubject(SUBJECT_NAME)
            .ackPolicy(AckPolicy.Explicit)
            .maxDeliver(RETRY_COUNT_LIMIT)
            .ackWait(ACK_WAIT)
            .build();

        return jsm.addOrUpdateConsumer(STREAM_NAME, config);
    }

    @Bean(destroyMethod = "unsubscribe")
    public JetStreamSubscription pullSubscription(
        JetStream js, ConsumerInfo consumerInfo
    ) throws IOException, JetStreamApiException {
        PullSubscribeOptions options = PullSubscribeOptions.builder()
            .durable(consumerInfo.getName())
            .build();

        return js.subscribe(SUBJECT_NAME, options);
    }
}
