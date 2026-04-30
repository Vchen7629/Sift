package app.config;

import java.io.IOException;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import io.nats.client.JetStreamApiException;
import io.nats.client.JetStreamManagement;
import io.nats.client.api.RetentionPolicy;
import io.nats.client.api.StorageType;
import io.nats.client.api.StreamConfiguration;
import io.nats.client.api.StreamInfo;

@Configuration
public class NatsStreamConfig {
    public static final String STREAM_NAME = "sift-job-processing-stream";

    private static final String SUBJECT_NAMES = "index-repo.subject.>";
    private static final int REPLICA_COUNT = 1;

    @Bean(destroyMethod = "close")
    public StreamInfo natsStreamConfig(JetStreamManagement jsm) throws IOException, JetStreamApiException {
        StreamConfiguration config = StreamConfiguration.builder()
            .name(STREAM_NAME)
            .subjects(SUBJECT_NAMES)
            .storageType(StorageType.File)
            .retentionPolicy(RetentionPolicy.WorkQueue)
            .replicas(REPLICA_COUNT)
            .build();

        try {
            return jsm.updateStream(config);
        } catch (JetStreamApiException e) {
            if (e.getErrorCode() == 404) {
                return jsm.addStream(config);
            }
            throw e;
        }
    }
}