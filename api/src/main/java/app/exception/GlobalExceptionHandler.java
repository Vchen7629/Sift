package app.exception;

import java.io.IOException;

import org.kohsuke.github.GHFileNotFoundException;
import org.kohsuke.github.HttpException;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.ControllerAdvice;
import org.springframework.web.bind.annotation.ExceptionHandler;

@ControllerAdvice
public class GlobalExceptionHandler {
    @ExceptionHandler(GHFileNotFoundException.class)
    public ResponseEntity<String> handleNotFound(GHFileNotFoundException e) {
        return ResponseEntity.status(404).body("Repository not found");
    }

    @ExceptionHandler(HttpException.class)
    public ResponseEntity<String> handleHttpException(HttpException e) {
        if (e.getResponseCode() == 403) {
            return ResponseEntity.status(403).body("GitHub rate limit exceeded");
        }
        return ResponseEntity.status(500).body("Github API error");
    }

    @ExceptionHandler(IOException.class)
    public ResponseEntity<String> handleIOException(IOException e) {
        return ResponseEntity.status(500).body("Github API error");
    }
}
