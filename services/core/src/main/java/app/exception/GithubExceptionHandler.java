package app.exception;

import org.kohsuke.github.GHFileNotFoundException;
import org.kohsuke.github.HttpException;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.ControllerAdvice;
import org.springframework.web.bind.annotation.ExceptionHandler;

@ControllerAdvice
public class GithubExceptionHandler {
    @ExceptionHandler(GHFileNotFoundException.class)
    public ResponseEntity<String> gHFileNotFoundException(GHFileNotFoundException e) {
        return ResponseEntity.status(404).body("Repository not found");
    }

    @ExceptionHandler(HttpException.class)
    public ResponseEntity<String> httpException(HttpException e) {
        if (e.getResponseCode() == 403) {
            return ResponseEntity.status(403).body("GitHub rate limit exceeded");
        }
        return ResponseEntity.status(500).body("Github API error");
    }
}
