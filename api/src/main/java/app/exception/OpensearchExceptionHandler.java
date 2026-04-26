package app.exception;

import java.util.NoSuchElementException;

import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.ControllerAdvice;
import org.springframework.web.bind.annotation.ExceptionHandler;

@ControllerAdvice
public class OpensearchExceptionHandler {
    @ExceptionHandler(NoSuchElementException.class)
    public ResponseEntity<String> noElementFound(NoSuchElementException e) {
        return ResponseEntity.status(404).body(e.getMessage());
    }
}
