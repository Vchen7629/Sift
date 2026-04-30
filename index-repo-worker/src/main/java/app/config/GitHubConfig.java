package app.config;

import java.io.IOException;

import org.kohsuke.github.GitHub;
import org.kohsuke.github.GitHubBuilder;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class GitHubConfig {

    @Value("${github.pat_token}")
    private String patToken;
    
    @Bean
    public GitHub gitHubClient() throws IOException {
        return new GitHubBuilder()
            .withOAuthToken(patToken)
            .build();
    }
}
