package app.exception;

import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.ControllerAdvice;
import org.springframework.web.bind.annotation.ExceptionHandler;

@ControllerAdvice
public class GeneralExceptionHandler {
    @ExceptionHandler(IllegalArgumentException.class)
    public ResponseEntity<String> illegalArgumentHander(IllegalArgumentException e) {
        return ResponseEntity.status(400).body(e.getMessage());
    }
}
