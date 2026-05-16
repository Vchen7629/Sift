package app.controller;

import static net.logstash.logback.argument.StructuredArguments.kv;

import java.io.IOException;
import java.util.Collections;

import org.kohsuke.github.GitHub;
import org.springframework.http.ResponseEntity;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.context.SecurityContext;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.web.context.HttpSessionSecurityContextRepository;
import org.springframework.security.web.context.SecurityContextRepository;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestHeader;
import org.springframework.web.bind.annotation.RequestMapping;


import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import lombok.extern.slf4j.Slf4j;

@Controller
@RequestMapping("auth")
@Slf4j
public class AuthController {
    private final SecurityContextRepository securityContextRepository = new HttpSessionSecurityContextRepository();
    
    @PostMapping("token")
    public ResponseEntity<String> validateGithubPatToken(
        @RequestHeader("Authorization") String bearerToken, HttpServletRequest request, HttpServletResponse response
    ) {
        log.info("recieved new auth request", kv("bearer", bearerToken));        

        String patToken = bearerToken.split("\\s+")[1];
        String username;

        try {
            GitHub gh = GitHub.connectUsingOAuth(patToken);
            username = gh.getMyself().getLogin(); 
        } catch (IOException e) {
            return ResponseEntity.status(403).body("Invalid Github Token");
        }

        UsernamePasswordAuthenticationToken auth = UsernamePasswordAuthenticationToken.authenticated(
            username, null, Collections.emptyList()
        );

        SecurityContext context = SecurityContextHolder.createEmptyContext();
        context.setAuthentication(auth);
        securityContextRepository.saveContext(context, request, response);

        return ResponseEntity.ok("Authenticated as " + username);
    }
}
